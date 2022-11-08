package api

import (
	"GuthiNetwork/platform"

	"github.com/gin-gonic/gin"
)

// yo kina cha tha chaina, tara chahincha
// TODO: For web dev guys to explain
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

const PORT = ":8080"

func StartServer(network_platform *platform.NetworkPlatform) {
	router := gin.Default()
	router.Use(CORSMiddleware())

	// nodes
	router.GET("/nodes", GetAvailableNodes(network_platform))
	router.GET("/self", GetSelfNode(network_platform))
	router.Run(PORT)
}
