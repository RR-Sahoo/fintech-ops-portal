package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/RR-Sahoo/fintech-starter-kit/backend/handlers"
	"github.com/RR-Sahoo/fintech-starter-kit/backend/repository"
)

// corsMiddleware configures CORS headers for frontend integration (e.g. Next.js on port 3000).
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	// Enable CORS middleware
	r.Use(corsMiddleware())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "fintech-backend",
		})
	})

	// In-memory transaction repository pre-seeded with mock fintech data
	transactionRepo := repository.NewMemoryTransactionRepository()
	transactionHandler := handlers.NewTransactionHandler(transactionRepo)

	// API v1 route group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/transactions", transactionHandler.GetTransactions)
		v1.POST("/transactions/reconcile", transactionHandler.ReconcileTransaction)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
