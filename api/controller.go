package api

import (
	"GuthiNetwork/core"
	"GuthiNetwork/platform"

	"github.com/gin-gonic/gin"
)

type AddressModel struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

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
		var address AddressModel
		if err := c.BindJSON(&address); err != nil {
			c.AbortWithError(400, err)
			return
		}
		c.JSON(200, 1)

		network_platform.ConnectToNode(address.IP + ":" + address.Port)
	}

	return fn
}

/*
Memory Information
*/
func GetMemoryInfo(network_platform *platform.NetworkPlatform) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		memory_info := make(map[string]core.MemoryStatus)
		for _, cache := range network_platform.Connection_caches {
			memory_info[cache.Node_ref.Name] = cache.Memory_info
		}
		c.JSON(200, memory_info)
	}

	return fn
}

func GetCpuInfo(network_platform *platform.NetworkPlatform) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		memory_info := make(map[string]core.ProcessorStatus)
		for _, cache := range network_platform.Connection_caches {
			memory_info[cache.Node_ref.Name] = cache.Cpu_info
		}
		c.JSON(200, memory_info)
	}

	return fn
}
