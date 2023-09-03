package main

import (
	"camphor-willow-god/text/crawl"
)

func main() {
	// 爬取书
	crawl.BookCrawl("凡人修仙传", "忘语")
	crawl.BookCrawl("星界", "千里送一血")
}
