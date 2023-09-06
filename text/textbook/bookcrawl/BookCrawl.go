package bookcrawl

import (
	"camphor-willow-god/text/textbook"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"net/url"
	"slices"
	"strings"
	"unicode"
)

func BookCrawl() {
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
	)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := strings.TrimSpace(e.Attr("href"))
		bookPathBase := "//www.qidian.com/textbook/"
		if !strings.Contains(href, bookPathBase) {
			return
		}
		bookPathBaseIndex := strings.Index(href, bookPathBase)
		remainPath := href[bookPathBaseIndex+len(bookPathBase) : len(href)-1]
		remainPathSplit := strings.Split(remainPath, "/")
		if len(remainPathSplit) != 1 {
			return
		}
		// 访问书籍页面
		_ = e.Request.Visit(e.Attr("href"))
	})
	books := new([]*textbook.Book)
	c.OnHTML(".textbook-info > .textbook-info-top", func(eBook *colly.HTMLElement) {
		processBookSummaryInfo(eBook, books)
	})
	c.OnHTML(".catalog-volume", func(eVolume *colly.HTMLElement) {
		processVolumeInfo(eVolume, books)
	})
	// 处理章节
	c.OnHTML(".chapter-wrapper > .relative > .print > .content", func(eChapter *colly.HTMLElement) {
		processChapterWords(eChapter, books)
	})
	// 起点首页
	_ = c.Visit("https://www.qidian.com/")
	c.Wait()

	for _, book := range *books {
		for _, v := range *book.Volumes {
			for _, c := range *v.Chapters {
				c.UpdateWordsBookIndex()
			}
		}
	}
}

func BookSearchCrawl(bookName string, author string) {
	trimBookName := strings.TrimSpace(bookName)
	if trimBookName == "" {
		return
	}
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
	)
	extensions.RandomUserAgent(c)
	c.OnHTML("#result-list > .textbook-img-text .textbook-mid-info", func(eBookInfo *colly.HTMLElement) {
		filterThenVisitBook(eBookInfo, author)
	})
	books := new([]*textbook.Book)
	c.OnHTML(".textbook-info > .textbook-info-top", func(eBook *colly.HTMLElement) {
		processBookSummaryInfo(eBook, books)
	})
	c.OnHTML(".catalog-volume", func(eVolume *colly.HTMLElement) {
		processVolumeInfo(eVolume, books)
	})
	// 处理章节
	c.OnHTML(".chapter-wrapper > .relative > .print > .content", func(eChapter *colly.HTMLElement) {
		processChapterWords(eChapter, books)
	})
	// 检索页
	searchUrl := "https://www.qidian.com/so/" + trimBookName + ".html"
	// 执行检索
	//_ = c.Visit("https://www.qidian.com/ajax/UserInfo/GetUserInfo?_csrfToken=75e49589-cc95-46b0-a21c-fe02d96ee82b")
	_ = c.Visit(searchUrl)
	c.Wait()

	for _, book := range *books {
		for _, v := range *book.Volumes {
			for _, c := range *v.Chapters {
				c.UpdateWordsBookIndex()
			}
		}
	}
}

// 处理章节字符
func processChapterWords(eChapter *colly.HTMLElement, books *[]*textbook.Book) {
	book := filterBookByOriginId(eChapter, books)
	// 提取章节ID
	chapterOriginId := extractOriginChapterId(eChapter.Request.URL.Path)
	// 章节
	for _, volume := range *book.Volumes {

		for _, chapter := range *volume.Chapters {

			if chapter == nil || chapter.OriginId != chapterOriginId {
				continue
			}
			chapterText := strings.TrimSpace(eChapter.Text)
			if len(chapterText) == 0 {
				continue
			}
			wIdx := 0
			for _, w := range chapterText {
				if unicode.IsSpace(w) {
					continue
				}
				word := textbook.NewBookWord()
				word.BookId = book.Id
				word.VolumeIndex = volume.Index
				word.ChapterIndex = chapter.Index
				word.Index = int32(wIdx)
				word.Word = string(w)
				*chapter.Words = append(*chapter.Words, word)
				wIdx++
			}
			// 持久化章节内容
			textbook.BatchCreateBookWords(chapter.Words)
			chapter.Words = nil
		}
	}
}

// 处理分卷信息
func processVolumeInfo(eVolume *colly.HTMLElement, books *[]*textbook.Book) {
	// 处理分卷
	volumeName := eVolume.ChildText("label > .volume-header > .volume-name")
	// 本程序主要分析正文内容文字之间的逻辑关系，所以忽略掉非正文部分
	if strings.Contains(volumeName, "作品相关") {
		return
	}
	// 暂时没搞明白怎么处理VIP订阅相关问题，跳过VIP分卷
	if strings.Contains(volumeName, "VIP") {
		return
	}
	book := filterBookByOriginId(eVolume, books)
	// 计算分卷索引位置
	volumeIndex := len(*book.Volumes)
	// 声明分卷
	volume := textbook.NewBookVolume()
	// 书ID
	volume.BookId = book.Id
	// 扩展书分卷
	*book.Volumes = append(*book.Volumes, volume)
	// 分卷名称
	volume.Name = volumeName
	// 分卷索引位置
	volume.Index = int32(volumeIndex)
	// 持久化分卷
	volume.InsertOrUpdate()
	// 处理章节名称
	eVolume.ForEach(".volume-chapters > .chapter-item > .chapter-name", func(_ int, eChapter *colly.HTMLElement) {
		processChapter(eChapter, volume, book)
	})
}

// 处理章节信息
func processChapter(eChapter *colly.HTMLElement, volume *textbook.BookVolume, book *textbook.Book) {
	// 章节链接字符串
	urlString := eChapter.Attr("href")
	// 计算章节索引位置
	chapterIndex := len(*volume.Chapters)
	// 声明章节
	chapter := textbook.NewBookChapter()
	// 书ID
	chapter.BookId = book.Id
	// 分卷ID
	chapter.VolumeIndex = volume.Index
	// 章节ID
	chapter.OriginId = extractOriginChapterId(urlString)
	// 章节索引
	chapter.Index = int32(chapterIndex)
	// 章节名称
	chapter.Name = eChapter.Text
	// 持久化章节
	chapter.InsertOrUpdate()
	// 章节
	*volume.Chapters = append(*volume.Chapters, chapter)
	// 访问章节
	_ = eChapter.Request.Visit(eChapter.Attr("href"))
}

// 处理书总览信息
func processBookSummaryInfo(eBook *colly.HTMLElement, books *[]*textbook.Book) {
	// 声明书实例
	book := textbook.NewBook()
	book.OriginId = extractOriginBookId(eBook.Request.URL.Path)
	book.Name = eBook.ChildText("#bookName")
	book.Author = eBook.ChildText(".author-name > .writer")
	// 处理分类
	categoryParent := new(int64)
	*categoryParent = -1
	eBook.ForEach(".author-name > a", func(index int, eCategory *colly.HTMLElement) {
		categoryByName := textbook.NewCategoryByName(eCategory.Text)
		categoryByName.Parent = *categoryParent
		categoryByName.Level = int8(index + 1)
		// 持久化分类
		categoryByName.InsertOrUpdate()
		*categoryParent = categoryByName.Id
		*book.Categories = append(*book.Categories, categoryByName)
	})
	if len(*book.Categories) > 0 {
		// 设置分类
		book.Category = (*book.Categories)[len(*book.Categories)-1].Id
	}
	book.Reset()
	book.InsertOrUpdate()
	*books = append(*books, book)
	// 处理标签
	eBook.ForEach(".author-name > span", func(_ int, eTag *colly.HTMLElement) {
		class := eTag.Attr("class")
		if slices.Contains([]string{"writer", "dot"}, class) {
			return
		}
		tagByName := textbook.NewTagByName(eTag.Text)
		// 持久化标签
		tagByName.InsertOrUpdate()
		// 持久化书标签
		bookTagWithArgs := textbook.NewBookTagWithArgs(book.Id, tagByName.Id)
		bookTagWithArgs.InsertOrUpdate()

		*book.Tags = append(*book.Tags, tagByName)
	})
}

// filterThenVisitBook 根据作者过滤出书并访问
func filterThenVisitBook(eBookInfo *colly.HTMLElement, author string) {
	trimAuthor := strings.TrimSpace(author)
	currentBookAuthor := eBookInfo.ChildText(".author > .name")
	if !strings.Contains(currentBookAuthor, trimAuthor) {
		return
	}
	bookUrl := eBookInfo.ChildAttr(".textbook-info-title > a", "href")
	_ = eBookInfo.Request.Visit(bookUrl)
}

// filterBookByOriginId 根据书原始ID过滤出书
func filterBookByOriginId(e *colly.HTMLElement, books *[]*textbook.Book) *textbook.Book {
	originBookId := extractOriginBookId(e.Request.URL.Path)
	for _, b := range *books {
		if b.OriginId == originBookId {
			return b
		}
	}
	panic(fmt.Errorf("textbook %s not exist", originBookId))
}

// extractOriginBookId 提取原始书ID
func extractOriginBookId(urlString string) string {
	var urlObj, err = new(url.URL).Parse(urlString)
	if err != nil {
		panic(err)
	}
	// 相对路径
	var path = urlObj.Path
	// 相对路径拆解元素
	var pathElements = strings.Split(path, "/")
	return pathElements[2]
}

// extractOriginChapterId 提取原始章节ID
func extractOriginChapterId(urlString string) string {
	// 链接
	var urlObj, err = new(url.URL).Parse(urlString)
	if err != nil {
		panic(err)
	}
	// 相对路径
	var path = urlObj.Path
	// 相对路径拆解元素
	var pathElements = strings.Split(path, "/")
	return pathElements[3]
}
