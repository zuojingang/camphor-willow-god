package logicrawl

import (
	"camphor-willow-god/text/textlogic"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"strings"
	"sync"
	"unicode"
)

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
		fmt.Println("发现第", len(urlSet.ToSlice()), "个新页面", href, "，扫荡一圈")
	})
	_ = c.Visit(entryPath)
	c.Wait()
}

// 语句处理sync.Map
func sentenceProcessSyncMap(eP *colly.HTMLElement, logicWordRefine *textlogic.LogicWordRefineSync) {
	wordStr := strings.TrimSpace(eP.Text)
	runeWordStr := []rune(wordStr)
	// 只处理完整的句子
	stopCode := "。！？"
	if len(runeWordStr) < 2 || !strings.Contains(stopCode, string(runeWordStr[len(runeWordStr)-1])) {
		return
	}
	for headIdx, headChar := range runeWordStr {
		if len(runeWordStr)-1 == headIdx {
			continue
		}
		if unicode.IsSpace(headChar) {
			continue
		}
		for endDealtIdx, endChar := range runeWordStr[headIdx+1:] {
			if unicode.IsSpace(endChar) {
				break
			}
			// 联想词组
			head := string(headChar)
			distance := 1 + endDealtIdx
			end := string(endChar)
			associationWordFrequencyMapValue, _ := logicWordRefine.WordsFrequency.LoadOrStore(head, new(sync.Map))
			associationWord := textlogic.LogicAssociationWord{
				Distance: int32(distance),
				End:      end,
			}
			associationWordFrequencyMap := associationWordFrequencyMapValue.(*sync.Map)
			frequencyValue, _ := associationWordFrequencyMap.LoadOrStore(associationWord, 0)
			frequency := frequencyValue.(int) + 1
			associationWordFrequencyMap.Store(associationWord, frequency)
			logicWordRefine.Count.Add(1)
			// 词组最长100字符
			if distance >= 100 {
				break
			}
		}
	}
	// 精炼词组
	logicWordRefine.Refine()
}
