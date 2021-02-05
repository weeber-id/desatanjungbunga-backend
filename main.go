package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/article"
	"github.com/weeber-id/desatanjungbunga-backend/src/services"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
)

func main() {
	// Environment section
	godotenv.Load("./devel.env")
	variables.InitializationVariable()

	// MongoDB section
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := services.InitializationMongo(ctx)
	defer client.Disconnect(ctx)

	port := 8080
	log.Printf(
		"Service Version %s run on port %d",
		variables.Version, port,
	)

	router := gin.Default()
	root := router.Group("/api")
	{
		root.GET("/", controllers.HealthCheck)

		admin := root.Group("/admin")
		{
			admin.GET("/article", article.GetOne)
			admin.GET("/articles", article.GetMultiple)
			admin.POST("/article/create", article.Create)
			admin.POST("/article/update", article.Update)
			admin.POST("/article/delete", article.Delete)
		}
	}

	router.Run(fmt.Sprintf(":%d", port))
}
