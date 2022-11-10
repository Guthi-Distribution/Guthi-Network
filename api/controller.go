package api

import (
	"GuthiNetwork/platform"

	"github.com/gin-gonic/gin"
)

func GetAvailableNodes(network_platform *platform.NetworkPlatform) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(200, network_platform.Connected_nodes)
	}

	return fn
}

func GetSelfNode(network_platform *platform.NetworkPlatform) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(200, network_platform.Self_node)
	}
	return fn
}

func PostConnectNode(network_platform *platform.NetworkPlatform) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var address string
		if err := c.BindJSON(address); err != nil {
			c.AbortWithError(400, err)
			return
		}

		network_platform.ConnectToNode(address)
	}

	return fn
}
