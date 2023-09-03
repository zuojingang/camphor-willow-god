package identified

import (
	"camphor-willow-god/config"
	"github.com/bwmarrin/snowflake"
)

var SnowflakeNode *snowflake.Node

func init() {
	snowflakeConfig := config.ApplicationConfig.Snowflake
	snowflake.Epoch = snowflakeConfig.StartTime
	snowflakeNode, err := snowflake.NewNode(snowflakeConfig.NodeId)
	if err != nil {
		panic(err)
	}
	SnowflakeNode = snowflakeNode
}

// IdGenerate 生成ID
func IdGenerate() int64 {
	return SnowflakeNode.Generate().Int64()
}
