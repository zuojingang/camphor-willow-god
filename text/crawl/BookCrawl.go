package crawl

import (
	"camphor-willow-god/text"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"net/url"
	"slices"
	"strings"
)

func BookCrawl(bookName string, author string) {
	trimBookName := strings.TrimSpace(bookName)
	if trimBookName == "" {
		return
	}
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	c.OnHTML("#result-list > .book-img-text .book-mid-info", func(eBookInfo *colly.HTMLElement) {
		filterThenVisitBook(eBookInfo, author)
	})
	books := new([]*text.Book)
	c.OnHTML(".book-info > .book-info-top", func(eBook *colly.HTMLElement) {
		processBookSummaryInfo(eBook, books)
	})
	c.OnHTML(".catalog-volume", func(eVolume *colly.HTMLElement) {
		processVolumeInfo(eVolume, books)
	})
	// 处理章节
	c.OnHTML(".chapter-wrapper > .relative > .print > .content", func(eChapter *colly.HTMLElement) {
		processParagraphs(eChapter, books)
	})
	// 检索页
	searchUrl := "https://www.qidian.com/so/" + trimBookName + ".html"
	// 执行检索
	//_ = c.Visit("https://www.qidian.com/ajax/UserInfo/GetUserInfo?_csrfToken=75e49589-cc95-46b0-a21c-fe02d96ee82b")
	_ = c.Visit(searchUrl)
}

// 处理段落
func processParagraphs(eChapter *colly.HTMLElement, books *[]*text.Book) {
	book := filterBookByOriginId(eChapter, books)
	// 提取章节ID
	chapterOriginId := extractOriginChapterId(eChapter.Request.URL.Path)
	// 章节
	for _, volume := range *book.Volumes {

		for index, chapter := range *volume.Chapters {

			if chapter == nil || chapter.OriginId != chapterOriginId {
				continue
			}
			eChapter.ForEach("p", func(_ int, eParagraph *colly.HTMLElement) {
				// 计算当前段落索引
				paragraphIndex := len(*chapter.Paragraph)
				// 声明段落
				paragraph := text.NewParagraph()
				// 书ID
				paragraph.BookId = book.Id
				// 分卷ID
				paragraph.VolumeId = volume.Id
				// 章节ID
				paragraph.ChapterId = chapter.Id
				// 段落索引
				paragraph.Index = int32(paragraphIndex)
				// 段落内容
				paragraph.Content = eParagraph.Text
				// 章节段落更新
				*chapter.Paragraph = append(*chapter.Paragraph, paragraph)
			})
			// 持久化段落
			text.PersistentParagraphs(chapter.Paragraph)
			//if strings.Contains(volume.Name, "VIP") {
			//	fmt.Println()
			//}
			// 释放空间
			(*volume.Chapters)[index] = nil
		}
	}
}

// 处理分卷信息
func processVolumeInfo(eVolume *colly.HTMLElement, books *[]*text.Book) {
	book := filterBookByOriginId(eVolume, books)
	// 计算分卷索引位置
	volumeIndex := len(*book.Volumes)
	// 声明分卷
	volume := text.NewVolume()
	// 书ID
	volume.BookId = book.Id
	// 扩展书分卷
	*book.Volumes = append(*book.Volumes, volume)
	// 处理分卷
	volume.Name = eVolume.ChildText("label > .volume-header > .volume-name")
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
func processChapter(eChapter *colly.HTMLElement, volume *text.Volume, book *text.Book) {
	// 章节链接字符串
	urlString := eChapter.Attr("href")
	// 计算章节索引位置
	chapterIndex := len(*volume.Chapters)
	// 声明章节
	chapter := text.NewChapter()
	// 书ID
	chapter.BookId = book.Id
	// 分卷ID
	chapter.VolumeId = volume.Id
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
func processBookSummaryInfo(eBook *colly.HTMLElement, books *[]*text.Book) {
	// 声明书实例
	book := text.NewBook()
	book.OriginId = extractOriginBookId(eBook.Request.URL.Path)
	book.Name = eBook.ChildText("#bookName")
	book.Author = eBook.ChildText(".author-name > .writer")
	// 处理分类
	categoryParent := new(int64)
	*categoryParent = -1
	eBook.ForEach(".author-name > a", func(index int, eCategory *colly.HTMLElement) {
		categoryByName := text.NewCategoryByName(eCategory.Text)
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
		tagByName := text.NewTagByName(eTag.Text)
		// 持久化标签
		tagByName.InsertOrUpdate()
		// 持久化书标签
		bookTagWithArgs := text.NewBookTagWithArgs(book.Id, tagByName.Id)
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
	bookUrl := eBookInfo.ChildAttr(".book-info-title > a", "href")
	_ = eBookInfo.Request.Visit(bookUrl)
}

// filterBookByOriginId 根据书原始ID过滤出书
func filterBookByOriginId(e *colly.HTMLElement, books *[]*text.Book) *text.Book {
	originBookId := extractOriginBookId(e.Request.URL.Path)
	for _, b := range *books {
		if b.OriginId == originBookId {
			return b
		}
	}
	panic(fmt.Errorf("book %s not exist", originBookId))
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
