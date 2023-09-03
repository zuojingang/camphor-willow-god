package crawl

import (
	"camphor-willow-god/text"
	"fmt"
	"github.com/gocolly/colly"
	"net/url"
	"strings"
)

func BookCrawl() {
	// 定义一个异常
	var err error
	// 声明书实例
	book := text.NewBook()
	book.Name = "凡人修仙传"
	book.Author = "忘语"
	// 初始化书分类
	*book.Categories = append(*book.Categories,
		text.NewCategoryByName("仙侠"),
		text.NewCategoryByName("古典仙侠"))
	if len(*book.Categories) > 0 {
		// 持久化分类
		text.PersistentCategories(book.Categories)
		// 设置分类
		book.Category = (*book.Categories)[len(*book.Categories)-1].Id
	}
	// 持久化书
	_, *book = book.TryInsert()
	// 初始化书标签
	*book.Tags = append(*book.Tags,
		text.NewTagByName("凡人流"))
	// 持久化标签
	if len(*book.Tags) > 0 {
		// 持久化标签
		text.PersistentTags(book.Tags)

		bookTags := new([]*text.BookTag)
		for _, tag := range *book.Tags {
			*bookTags = append(*bookTags, text.NewBookTagWithArgs(book.Id, tag.Id))
		}
		text.PersistentBookTags(bookTags)
	}

	c := colly.NewCollector()

	c.OnHTML(".catalog-volume", func(e *colly.HTMLElement) {
		// 计算分卷索引位置
		volumeIndex := len(*book.Volumes)
		// 声明分卷
		volume := text.NewVolume()
		// 书ID
		volume.BookId = book.Id
		// 扩展书分卷
		*book.Volumes = append(*book.Volumes, volume)
		// 处理分卷
		volume.Name = e.ChildText("label > .volume-header > .volume-name")
		// 分卷索引位置
		volume.Index = int32(volumeIndex)
		// 分卷名称
		if strings.Contains(volume.Name, "VIP") {
			return
		}
		// 持久化分卷
		_, *volume = volume.TryInsert()
		// 处理章节名称
		e.ForEach(".volume-chapters > .chapter-item > .chapter-name", func(_ int, eChapter *colly.HTMLElement) {
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
			chapter.OriginId, err = ExtractChapterId(urlString)
			if err != nil {
				fmt.Println("some error")
				return
			}
			// 章节索引
			chapter.Index = int32(chapterIndex)
			// 章节名称
			chapter.Name = eChapter.Text
			// 章节
			*volume.Chapters = append(*volume.Chapters, chapter)
			// 持久化章节
			var tryInsert bool
			tryInsert, *chapter = chapter.TryInsert()
			if tryInsert {
				// 访问章节
				err = eChapter.Request.Visit(eChapter.Attr("href"))
				if err != nil {
					fmt.Println("some error")
					return
				}
			}
		})
	})
	// 处理章节
	c.OnHTML(".chapter-wrapper > .relative > .print > .content", func(e *colly.HTMLElement) {
		// 提取章节ID
		chapterOriginId, _ := ExtractChapterId(e.Request.URL.Path)
		// 章节
		for _, volume := range *book.Volumes {

			for index, chapter := range *volume.Chapters {

				if chapter == nil || chapter.OriginId != chapterOriginId {
					continue
				}
				e.ForEach("p", func(_ int, eParagraph *colly.HTMLElement) {
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
				// 先清理掉之前存储的段落
				chapter.ClearParagraphs()
				// 持久化段落
				text.PersistentParagraphs(chapter.Paragraph)
				// 释放空间
				(*volume.Chapters)[index] = nil
			}
		}
	})
	// 访问入口地址
	_ = c.Visit("https://www.qidian.com/book/107580/")
}

// ExtractChapterId 提取章节ID
func ExtractChapterId(urlString string) (string, error) {
	// 章节链接
	var urlObj, err = new(url.URL).Parse(urlString)
	if err != nil {
		fmt.Println("some error")
		return "", err
	}
	// 章节链接相对路径
	var path = urlObj.Path
	// 章节链接相对路径拆解元素
	var pathElements = strings.Split(path, "/")
	// 章节ID
	return pathElements[3], nil
}
