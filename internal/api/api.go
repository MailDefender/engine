package api

import (
	"github.com/gin-gonic/gin"

	"maildefender/engine/internal/api/handlers"
	"maildefender/engine/internal/utils"
)

var engine *gin.Engine

func init() {
	if utils.GetEnvBool("DEV_MODE", false) {
		gin.SetMode(gin.ReleaseMode)
	}
	engine = gin.New()
}

func RegisterHandlers() {
	v1 := engine.Group("/v1/engine")

	emails := v1.Group("/rules")
	emails.GET("", handlers.GetRules)
	emails.GET("/:id", handlers.GetRuleByID)
	emails.DELETE("/:id", handlers.DeleteRule)

	mailboxes := v1.Group("/reputations")
	mailboxes.GET("", handlers.GetReputations)
	mailboxes.GET("/search", handlers.SearchReputation)

	tokens := v1.Group("/token")
	tokens.POST("/validate/:token", handlers.ValidateToken)
}

func Run() error {
	return engine.Run(":8080")
}
