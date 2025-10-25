package middleware

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
)


type HTTPStats struct {
	activeConnections int64
}

var globalStats = &HTTPStats{}


func StatsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		atomic.AddInt64(&globalStats.activeConnections, 1)

		
		defer func() {
			atomic.AddInt64(&globalStats.activeConnections, -1)
		}()

		c.Next()
	}
}


type StatsInfo struct {
	ActiveConnections int64 `json:"active_connections"`
}


func GetStats() StatsInfo {
	return StatsInfo{
		ActiveConnections: atomic.LoadInt64(&globalStats.activeConnections),
	}
}
