package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/about"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/account"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/analytics"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/article"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/discussion"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/handcraft"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/kuliner"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/lodging"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/media"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/search"
	"github.com/weeber-id/desatanjungbunga-backend/src/controllers/travel"
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
	root := router.Group("/backend")
	{
		root.GET("/", controllers.HealthCheck)
		root.POST("/login", account.AdminLogin)
		root.POST("/logout", account.AdminLogut)

		root.GET("/about", about.Get)

		root.GET("/search", search.GetSearch)

		root.GET("/article", article.GetOne)
		root.GET("/articles", article.GetMultiple)

		root.GET("/travel", travel.GetOne)
		root.GET("/travels", travel.GetMultiple)

		root.GET("/handcraft", handcraft.GetOne)
		root.GET("/handcrafts", handcraft.GetMultiple)

		root.GET("/culinary", kuliner.GetOne)
		root.GET("/culinaries", kuliner.GetMultiple)

		root.GET("/lodging", lodging.GetOne)
		root.GET("/lodgings", lodging.GetMultiple)

		root.GET("/discussion", discussion.GetMultiple)
		root.POST("/discussion/create", discussion.Create)
		root.POST("/discussion/delete", discussion.Delete)

		admin := root.Group("/admin")
		admin.Use(middlewares.AdminAuthorization())
		{
			admin.GET("/", account.AdminInformation)
			admin.GET("/list", account.AdminList)
			admin.POST("/register", account.AdminCreate)
			admin.POST("/delete", account.AdminDelete)
			admin.POST("/update/password", account.AdminChangePassword)
			admin.POST("/update/profile-picture", account.AdminChangeProfilePicture)
			admin.POST("/update/profile-picture/delete", account.AdminDeleteProfilePicture)

			admin.GET("/about", about.AdminGet)
			admin.POST("/about/update", about.AdminUpdate)
			admin.POST("/about/delete", about.AdminDelete)

			admin.GET("/analytics/content-count", analytics.ContentCount)

			admin.POST("/media/upload/public", media.UploadPublicFile)

			admin.GET("/article", article.GetOne)
			admin.GET("/articles", article.GetMultiple)
			admin.POST("/article/create", article.Create)
			admin.POST("/article/update", article.Update)
			admin.POST("/article/update/active", article.ChangeActiveDeactive)
			admin.POST("/article/update/recommendation", article.ChangeRecommendation)
			admin.POST("/article/delete", article.Delete)

			admin.GET("/travel", travel.GetOne)
			admin.GET("/travels", travel.GetMultiple)
			admin.POST("/travel/create", travel.Create)
			admin.POST("/travel/update", travel.Update)
			admin.POST("/travel/update/active", travel.ChangeActiveDeactive)
			admin.POST("/travel/update/recommendation", travel.ChangeRecommendation)
			admin.POST("/travel/delete", travel.Delete)

			admin.GET("/handcraft", handcraft.GetOne)
			admin.GET("/handcrafts", handcraft.GetMultiple)
			admin.POST("/handcraft/create", handcraft.Create)
			admin.POST("/handcraft/update", handcraft.Update)
			admin.POST("/handcraft/update/active", handcraft.ChangeActiveDeactive)
			admin.POST("/handcraft/update/recommendation", handcraft.ChangeRecommendation)
			admin.POST("/handcraft/delete", handcraft.Delete)

			admin.GET("/culinary", kuliner.GetOne)
			admin.GET("/culinaries", kuliner.GetMultiple)
			admin.POST("/culinary/create", kuliner.Create)
			admin.POST("/culinary/update", kuliner.Update)
			admin.POST("/culinary/update/active", kuliner.ChangeActiveDeactive)
			admin.POST("/culinary/update/recommendation", kuliner.ChangeRecommendation)
			admin.POST("/culinary/delete", kuliner.Delete)

			admin.GET("/lodging", lodging.GetOne)
			admin.GET("/lodgings", lodging.GetMultiple)
			admin.POST("/lodging/create", lodging.Create)
			admin.POST("/lodging/update", lodging.Update)
			admin.POST("/lodging/update/active", lodging.ChangeActiveDeactive)
			admin.POST("/lodging/update/recommendation", lodging.ChangeRecommendation)
			admin.POST("/lodging/delete", lodging.Delete)
			admin.GET("/lodging/facilities", lodging.GetFacilities)
			admin.POST("/lodging/facility/create", lodging.CreateFacilities)

			admin.GET("/discussion", discussion.GetMultiple)
			admin.POST("/discussion/create", discussion.Create)
			admin.POST("/discussion/delete", discussion.Delete)
		}
	}

	router.Run(fmt.Sprintf(":%d", port))
}
