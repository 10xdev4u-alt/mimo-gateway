package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleBenchmark(c *gin.Context) {
	start := time.Now()

	iterations := 10
	for i := 0; i < iterations; i++ {
		time.Sleep(10 * time.Millisecond)
	}

	elapsed := time.Since(start)
	avgMs := elapsed.Milliseconds() / int64(iterations)

	c.JSON(http.StatusOK, gin.H{
		"iterations": iterations,
		"total_ms":   elapsed.Milliseconds(),
		"avg_ms":     avgMs,
		"ops_per_sec": float64(iterations) / elapsed.Seconds(),
	})
}
