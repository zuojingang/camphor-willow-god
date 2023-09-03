package snowflake

import (
	"camphor-willow-god/config"
	"github.com/bwmarrin/snowflake"
)

func init() {
	snowflakeConfig := config.ApplicationConfig.Snowflake
	snowflake.Epoch = snowflakeConfig.StartTime
	snowflake.NewNode(snowflakeConfig.NodeId)
	
}
