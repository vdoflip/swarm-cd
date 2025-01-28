package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/m-adawi/swarm-cd/util"
	sloggin "github.com/samber/slog-gin"
)

var router *gin.Engine = gin.New()

func init() {
	router.Use(sloggin.New(util.Logger))
	router.GET("/stacks", getStacks)
	router.StaticFile("/ui", "ui/dist/index.html")
	router.Static("/assets", "ui/dist/assets")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/ui")
	})
}

func RunServer() {
	config := util.Configs.Web
	if config == nil {
		config = &util.WebConfig{
			Host: "localhost",
			Port: 8080,
		}
	}
	
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	util.Logger.Info("starting web server", "address", addr)
	router.Run(addr)
}
