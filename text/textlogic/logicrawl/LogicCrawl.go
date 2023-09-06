package logicrawl

import (
	"camphor-willow-god/text/textlogic"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"unicode"
)

func LogicCrawl(entryPath string) {
	c := colly.NewCollector(
	//colly.Debugger(&debug.LogDebugger{}),
	)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	// 初始化词组精炼器
	logicWordRefine := textlogic.NewLogicWordRefine()
	// 监控段落，并惊醒词组精炼
	c.OnHTML("p", func(eP *colly.HTMLElement) {
		sentenceProcess(eP, logicWordRefine)
	})
	// 页面地址集合，存放所有访问过的页面地址，防止重复访问
	urlSet := mapset.NewSet[string]()
	// 监控链接，并执行调用
	c.OnHTML("a[href]", func(eHref *colly.HTMLElement) {
		href := eHref.Attr("href")
		// 为避免重复访问页面（或许是相互之间引用形成的环）
		if urlSet.Contains(href) {
			//fmt.Println("重复页面，跳过访问")
			return
		}
		urlSet.Add(href)
		_ = eHref.Request.Visit(href)
		fmt.Println("发现新页面，走，我们进去逛逛")
	})
	_ = c.Visit(entryPath)
}

func LogicCrawlAsync(entryPath string) {
	c := colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
	)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	// 初始化词组精炼器
	logicWordRefine := textlogic.NewLogicWordRefineSync()
	// 监控段落，并惊醒词组精炼
	c.OnHTML("p", func(eP *colly.HTMLElement) {
		sentenceProcessSyncMap(eP, logicWordRefine)
	})
	// 页面地址集合，存放所有访问过的页面地址，防止重复访问
	urlSet := mapset.NewSet[string]()
	// 监控链接，并执行调用
	c.OnHTML("a[href]", func(eHref *colly.HTMLElement) {
		href := eHref.Attr("href")
		// 为避免重复访问页面（或许是相互之间引用形成的环）
		if urlSet.Contains(href) {
			//fmt.Println("重复页面，跳过访问")
			return
		}
		urlSet.Add(href)
		_ = eHref.Request.Visit(href)
		fmt.Println("发现新页面，走，我们进去逛逛")
	})
	_ = c.Visit(entryPath)
	//c.Wait()
	// 信号接收通道
	signals := make(chan os.Signal)
	// 订阅信号
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range signals {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("receive exit signal ", s.String(), ",exit...")
			logicWordRefine.PersistentLogicWords()
			os.Exit(0)
		}
	}
}

// 语句处理map
func sentenceProcess(eP *colly.HTMLElement, logicWordRefine *textlogic.LogicWordRefine) {
	wordStr := strings.TrimSpace(eP.Text)
	runeWordStr := []rune(wordStr)
	// 只处理完整的句子
	if len(runeWordStr) < 2 || !strings.Contains("。！？", string(runeWordStr[len(runeWordStr)-1])) {
		return
	}
	for wordIdx, char := range runeWordStr {
		if len(runeWordStr)-1 == wordIdx {
			continue
		}
		if unicode.IsSpace(char) {
			continue
		}
		nextChar := runeWordStr[wordIdx+1]
		if unicode.IsSpace(nextChar) {
			continue
		}
		word := string(char)
		nextWord := string(nextChar)
		// 联想词组
		associationWord := textlogic.LogicAssociationWord{End: nextWord}
		if logicWordRefine.WordsFrequency[word] == nil {
			logicWordRefine.WordsFrequency[word] = make(map[textlogic.LogicAssociationWord]int32)
		}
		logicWordRefine.WordsFrequency[word][associationWord]++
		// 精炼扩展词组
		logicWordRefine.RefineExtension()
	}
}

// 语句处理sync.Map
func sentenceProcessSyncMap(eP *colly.HTMLElement, logicWordRefine *textlogic.LogicWordRefineSync) {
	wordStr := strings.TrimSpace(eP.Text)
	runeWordStr := []rune(wordStr)
	// 只处理完整的句子
	if len(runeWordStr) < 2 || !strings.Contains("。！？", string(runeWordStr[len(runeWordStr)-1])) {
		return
	}
	for wordIdx, char := range runeWordStr {
		if len(runeWordStr)-1 == wordIdx {
			continue
		}
		if unicode.IsSpace(char) {
			continue
		}
		nextChar := runeWordStr[wordIdx+1]
		if unicode.IsSpace(nextChar) {
			continue
		}
		word := string(char)
		nextWord := string(nextChar)
		// 联想词组
		associationWord := textlogic.LogicAssociationWord{End: nextWord}
		associationWordFrequencyMapValue, _ := logicWordRefine.WordsFrequency.LoadOrStore(word, new(sync.Map))
		associationWordFrequencyMap := associationWordFrequencyMapValue.(*sync.Map)
		frequencyValue, _ := associationWordFrequencyMap.LoadOrStore(associationWord, 0)
		frequency := frequencyValue.(int) + 1
		associationWordFrequencyMap.Store(associationWord, frequency)
		// 精炼扩展词组
		logicWordRefine.RefineExtension()
	}
}
