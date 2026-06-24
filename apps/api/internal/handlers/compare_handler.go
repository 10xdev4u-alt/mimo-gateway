package handlers

import (

)

func HandleCompareModels(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"models": []gin.H{
			{"id": "mimo-auto", "name": "MiMo Auto", "speed": "fast", "cost": "free"},
		},
		"comparison": gin.H{
			"mimo-auto": gin.H{
				"tokens_per_second": 12,
				"avg_latency":       "14.2s",
				"quality":           "good",
			},
		},
	})
}
