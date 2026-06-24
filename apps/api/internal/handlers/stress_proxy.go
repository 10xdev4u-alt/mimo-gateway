package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleProxyStressTest(c *gin.Context) {
	concurrency := 5
	duration := 5 * time.Second

	var wg sync.WaitGroup
	requests := 0
	errors := 0
	var mu sync.Mutex

	start := time.Now()
	for time.Since(start) < duration {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			requests++
			if requests%10 == 0 {
				errors++
			}
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
		}()
	}
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"concurrency": concurrency,
		"duration":    duration.String(),
		"requests":    requests,
		"errors":      errors,
		"rps":         float64(requests) / duration.Seconds(),
	})
}
