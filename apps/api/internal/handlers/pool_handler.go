package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type ConnectionPool struct {
	mu        sync.Mutex
	active    int
	maxActive int
	total     int64
}

var pool = &ConnectionPool{maxActive: 10}

func (p *ConnectionPool) Acquire() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.active >= p.maxActive {
		return false
	}
	p.active++
	p.total++
	return true
}

func (p *ConnectionPool) Release() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.active > 0 {
		p.active--
	}
}

func HandlePoolStats(c *gin.Context) {
	pool.mu.Lock()
	stats := gin.H{
		"active_connections": pool.active,
		"max_connections":    pool.maxActive,
		"total_connections":  pool.total,
	}
	pool.mu.Unlock()

	c.JSON(http.StatusOK, stats)
}
