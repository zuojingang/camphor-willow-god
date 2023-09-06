package textlogic

import (
	"sync"
	"sync/atomic"
)

// LogicWordRefineSync 词组精炼
type LogicWordRefineSync struct {
	WordsFrequency *sync.Map            // 词组=>出现频次
	Count          *atomic.Int32        // 当前层级词组数量
	Level          int32                // 当前层级
	MaxLevel       int32                // 最大层级
	Next           *LogicWordRefineSync // 下一层精炼
}

func NewLogicWordRefineSync() *LogicWordRefineSync {
	return &LogicWordRefineSync{
		WordsFrequency: new(sync.Map),
		Count:          new(atomic.Int32),
		Level:          1,
	}
}

func NewLogicWordRefineSyncWithParams(level int32, maxLevel int32) *LogicWordRefineSync {
	return &LogicWordRefineSync{
		WordsFrequency: new(sync.Map),
		Count:          new(atomic.Int32),
		Level:          level,
		MaxLevel:       maxLevel,
	}
}

// Refine 精炼词组
func (lwr *LogicWordRefineSync) Refine() {
	if lwr == nil {
		return
	}
	// 每一千个字符处理一次
	if lwr.Count.Load() < 1000 {
		return
	}
	// 默认精炼1000000（一百万）次
	if lwr.MaxLevel == 0 {
		lwr.MaxLevel = 1000000
	}
	if lwr.Level >= lwr.MaxLevel {
		return
	}
	wordsFrequency := lwr.WordsFrequency
	wordsFrequency.Range(func(key, value any) bool {
		wordHead := key.(string)
		associationWordMap := value.(*sync.Map)
		associationWordMap.Range(func(key, value any) bool {
			associationWord := key.(LogicAssociationWord)
			frequency := value.(int)
			// 词频达标，则向上晋位
			if frequency < 2 {
				// 重置词组频次
				associationWordMap.Store(associationWord, 0)
				return true
			}
			if lwr.Next == nil {
				lwr.Next = NewLogicWordRefineSyncWithParams(lwr.Level+1, lwr.MaxLevel)
			}
			nextLevelAssociationWordMapValue, _ := lwr.Next.WordsFrequency.LoadOrStore(wordHead, new(sync.Map))
			nextLevelAssociationWordMap := nextLevelAssociationWordMapValue.(*sync.Map)
			nextLevelAssociationWordFrequencyValue, _ := nextLevelAssociationWordMap.LoadOrStore(associationWord, 0)
			nextLevelAssociationWordFrequency := nextLevelAssociationWordFrequencyValue.(int) + 1
			nextLevelAssociationWordMap.Store(associationWord, nextLevelAssociationWordFrequency)
			lwr.Next.Count.Add(1)
			// 经过足够过滤层级的词汇需要入库（创建或修改）
			if lwr.Level >= 100 {
				logicWord := LogicWord{
					Head:                 wordHead,
					LogicAssociationWord: associationWord,
					Level:                lwr.Level,
				}
				logicWord.InsertOrUpdate()
			}
			return true
		})
		return true
	})
	lwr.WordsFrequency = new(sync.Map)
	lwr.Count.Store(0)
	lwr.Next.Refine()
}
