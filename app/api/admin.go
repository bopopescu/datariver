package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetMyRoleResp struct {
	Role string `json:"role"`
}

func GetMyRole(c *gin.Context) {
	resp := GetMyRoleResp{Role: "only test, not work"}
	c.JSON(http.StatusOK, resp)
}

func SetLogLevel(c *gin.Context) {
	c.JSON(http.StatusOK, struct{}{})
	/*
		level := c.Query("level")
		log.Info("Start SetLogLevel %+v", level)
		if level != "info" && level != "debug" && level != "error" && level != "warn" {
			c.JSON(http.StatusOK, "levle should be info debug error warn")
			return
		}
		log.SetLevelByString(level)
		tlog.SetLevelByString(level)
		c.JSON(http.StatusOK, fmt.Sprintf("set log level to %+v success", level))
	*/
}
