package textlogic

import (
	"fmt"
	"sync"
)

// LogicWordRefineSync 词组精炼
type LogicWordRefineSync struct {
	WordsFrequency *sync.Map            // 词组=>出现频次
	Level          int32                // 当前层级
	MaxLevel       int32                // 最大层级
	Next           *LogicWordRefineSync // 下一层精炼
}

func NewLogicWordRefineSync() *LogicWordRefineSync {
	return &LogicWordRefineSync{
		WordsFrequency: new(sync.Map),
		Level:          1,
	}
}

func NewLogicWordRefineSyncWithParams(level int32, maxLevel int32) *LogicWordRefineSync {
	return &LogicWordRefineSync{
		WordsFrequency: new(sync.Map),
		Level:          level,
		MaxLevel:       maxLevel,
	}
}

// RefineExtension 精炼扩展词汇
func (lwr *LogicWordRefineSync) RefineExtension() {
	if lwr == nil {
		return
	}
	// 默认精炼扩展1000000（一百万）次
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
			// 重置词组频次
			associationWordMap.Store(associationWord, 0)
			//// 经过十层以上过滤的词汇都需要入库（创建或修改）
			//if lwr.Level >= 10 {
			//	logicWord := LogicWord{
			//		Head:                 wordHead,
			//		LogicAssociationWord: associationWord,
			//		Level:                lwr.Level,
			//	}
			//	logicWord.InsertOrUpdate()
			//}
			nextAssociationWordMapValue, _ := wordsFrequency.Load(associationWord.End)
			if nextAssociationWordMapValue == nil {
				return true
			}
			nextAssociationWordMap := nextAssociationWordMapValue.(*sync.Map)
			nextAssociationWordMap.Range(func(key, value any) bool {
				nextAssociationWord := key.(LogicAssociationWord)
				nextAssociationWordFrequency := value.(int)
				// 词频达标，则向上晋位
				if nextAssociationWordFrequency < 2 {
					return true
				}
				combineAssociationWord := LogicAssociationWord{
					Middle: associationWord.Middle + associationWord.End + nextAssociationWord.Middle,
					End:    nextAssociationWord.End,
				}
				nextLevelCombineAssociationWordFrequencyValue, _ := nextLevelAssociationWordMap.LoadOrStore(combineAssociationWord, 0)
				nextLevelCombineAssociationWordFrequency := nextLevelCombineAssociationWordFrequencyValue.(int) + 1
				nextLevelAssociationWordMap.Store(combineAssociationWord, nextLevelCombineAssociationWordFrequency)
				//// 经过十层以上过滤的词汇都需要入库（创建或修改）
				//if lwr.Level >= 10 {
				//	combineWord := LogicWord{
				//		Head:                 wordHead,
				//		LogicAssociationWord: combineAssociationWord,
				//		Level:                lwr.Level,
				//	}
				//	combineWord.InsertOrUpdate()
				//}
				return true
			})
			return true
		})
		return true
	})
	lwr.Next.RefineExtension()
}

func (lwr *LogicWordRefineSync) PersistentLogicWords() {
	for lwr != nil {
		lwr.WordsFrequency.Range(func(key, value any) bool {
			head := key.(string)
			associationWordFrequencyMap := value.(*sync.Map)
			associationWordFrequencyMap.Range(func(key, value any) bool {
				associationWord := key.(LogicAssociationWord)
				// 经过十层以上过滤的词汇都需要入库（创建或修改）
				if lwr.Level > 10 {
					logicWord := LogicWord{
						Head:                 head,
						LogicAssociationWord: associationWord,
						Level:                lwr.Level,
					}
					logicWord.InsertOrUpdate()
					fmt.Println("write db before exit...")
				}
				return true
			})
			*lwr = *lwr.Next
			return true
		})
	}
}
