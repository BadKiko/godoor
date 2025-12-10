package routes

import (
	"godoor/serializers"
	"github.com/gin-gonic/gin"
)

// ServerSettingsViewSet handles server settings endpoints
func ServerSettingsCurrent(c *gin.Context) {
	// No authentication required - public endpoint like in Tandoor
	// Attention: No login required, do not return sensitive data

	settings := serializers.SerializeServerSettings()
	c.JSON(200, settings)
}
