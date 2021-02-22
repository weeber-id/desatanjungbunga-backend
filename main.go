package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/account"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/article"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/belanja"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/discussion"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/kuliner"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/media"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/wisata"
	"github.com/weeber-id/desatanjungbunga-backend/src/middlewares"
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

	services.InitializationMinio()

	port := 8080
	log.Printf(
		"Service Version %s run on port %d",
		variables.Version, port,
	)

	router := gin.Default()
	router.Use(middlewares.CORS())
	root := router.Group("/api")
	{
		root.GET("/", controllers.HealthCheck)
		root.POST("/login", account.AdminLogin)
		root.POST("/logout", account.AdminLogut)

		admin := root.Group("/admin")
		admin.Use(middlewares.AdminAuthorization())
		{
			admin.POST("/register", account.AdminCreate)

			admin.POST("/media/upload/public", media.UploadPublicFile)

			admin.GET("/article", article.GetOne)
			admin.GET("/articles", article.GetMultiple)
			admin.POST("/article/create", article.Create)
			admin.POST("/article/update", article.Update)
			admin.POST("/article/delete", article.Delete)

			admin.GET("/travel", wisata.GetOne)
			admin.GET("/travels", wisata.GetMultiple)
			admin.POST("/travel/create", wisata.Create)
			admin.POST("/travel/update", wisata.Update)
			admin.POST("/travel/delete", wisata.Delete)

			admin.GET("/shopping", belanja.GetOne)
			admin.GET("/shoppings", belanja.GetMultiple)
			admin.POST("/shopping/create", belanja.Create)
			admin.POST("/shopping/update", belanja.Update)
			admin.POST("/shopping/delete", belanja.Delete)

			admin.GET("/culinary", kuliner.GetOne)
			admin.GET("/culinaries", kuliner.GetMultiple)
			admin.POST("/culinary/create", kuliner.Create)
			admin.POST("/culinary/update", kuliner.Update)
			admin.POST("/culinary/delete", kuliner.Delete)

			admin.GET("/discussion", discussion.GetMultiple)
			admin.POST("/discussion/create", discussion.Create)
			admin.POST("/discussion/delete", discussion.Delete)
		}
	}

	router.Run(fmt.Sprintf(":%d", port))
}
