package main

import (
	"camphor-willow-god/text/textlogic/logicrawl"
)

func main() {
	// 爬取书
	//bookcrawl.BookCrawl()
	//bookcrawl.BookSearchCrawl("凡人修仙传", "忘语")
	//bookcrawl.BookCrawl("星界", "千里送一血")

	// 词汇分析
	logicrawl.LogicCrawlAsync("https://www.hao123.com/")
}
