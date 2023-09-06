package textlogic

// LogicWordRefine 词组精炼
type LogicWordRefine struct {
	WordsFrequency map[string]map[LogicAssociationWord]int32 // 词组=>出现频次
	Level          int32                                     // 当前层级
	MaxLevel       int32                                     // 最大层级
	Next           *LogicWordRefine                          // 下一层精炼
}

func NewLogicWordRefine() *LogicWordRefine {
	return &LogicWordRefine{
		WordsFrequency: make(map[string]map[LogicAssociationWord]int32),
	}
}

func NewLogicWordRefineWithParams(level int32, maxLevel int32) *LogicWordRefine {
	return &LogicWordRefine{
		WordsFrequency: make(map[string]map[LogicAssociationWord]int32),
		Level:          level,
		MaxLevel:       maxLevel,
	}
}

// RefineExtension 精炼扩展词汇
func (lwr *LogicWordRefine) RefineExtension() {
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
	for wordHead, associationWordMap := range wordsFrequency {
		for associationWord, frequency := range associationWordMap {
			// 词频达标，则向上晋位
			if frequency < 3 {
				continue
			}
			if lwr.Next == nil {
				lwr.Next = NewLogicWordRefineWithParams(lwr.Level+1, lwr.MaxLevel)
			}
			if lwr.Next.WordsFrequency[wordHead] == nil {
				lwr.Next.WordsFrequency[wordHead] = make(map[LogicAssociationWord]int32)
			}
			lwr.Next.WordsFrequency[wordHead][associationWord]++
			// 经过十层以上过滤的词汇都需要入库（创建或修改）
			if lwr.Level >= 10 {
				logicWord := LogicWord{
					Head:                 wordHead,
					LogicAssociationWord: associationWord,
					Level:                lwr.Level,
				}
				logicWord.InsertOrUpdate()
			}
			// 重置词组频次
			associationWordMap[associationWord] = 0

			nextAssociationWordMap := wordsFrequency[associationWord.End]
			if nextAssociationWordMap == nil || &nextAssociationWordMap == &associationWordMap {
				continue
			}
			for nextAssociationWord, nextAssociationWordFrequency := range nextAssociationWordMap {
				// 词频达标，则向上晋位
				if nextAssociationWordFrequency < 3 {
					continue
				}
				combineAssociationWord := LogicAssociationWord{
					Middle: associationWord.Middle + associationWord.End + nextAssociationWord.Middle,
					End:    nextAssociationWord.End,
				}
				lwr.Next.WordsFrequency[wordHead][combineAssociationWord]++
				// 经过十层以上过滤的词汇都需要入库（创建或修改）
				if lwr.Level >= 10 {
					combineWord := LogicWord{
						Head:                 wordHead,
						LogicAssociationWord: combineAssociationWord,
						Level:                lwr.Level,
					}
					combineWord.InsertOrUpdate()
				}
			}
		}
	}
	lwr.Next.RefineExtension()
}
