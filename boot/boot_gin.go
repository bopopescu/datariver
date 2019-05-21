package boot

import (
	"fmt"
	"time"

	"datariver/app/api"

	"github.com/gin-gonic/gin"
)

func HandleTimeLoger(c *gin.Context) {
	start_time := time.Now()
	defer func() {
		info := fmt.Sprintf("Handle[%s][cost:%v]",
			c.Request.URL.Path, time.Now().Sub(start_time))
		fmt.Println(info)
	}()
	c.Next()
}

func StartGinServer() error {
	router := gin.Default()
	router.Use(HandleTimeLoger)
	//router.RegisterLoggerInfo(ginLogger)

	admin := router.Group("/admin")
	{
		admin.GET("/print_role", api.GetMyRole)
		admin.POST("/set_log_level", api.SetLogLevel)
	}

	go func() {
		router.Run(GConfig.BrokerConfig.RPCListen)
	}()

	return nil
}
