package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"micron/commons"
	"micron/micron"
	"micron/middleware"
	tagDatabase "micron/scraper"
	scraper "micron/scraper/runnable"
	"micron/tag"
	"micron/user"
	"time"
)

func main() {
	// does the scraping stuff
	scraper.Start(false)

	client, cancel := initDb()
	defer (*cancel)()

	scraperService := tagDatabase.NewScraper()
	tagRepository := tag.NewRepository(scraperService)
	micronService := micron.NewService(scraperService)

	jwtService := commons.NewJwtService(client)

	userService := user.NewService(user.ServiceConfig{
		Store:   user.NewRepository(client),
		Tags:    tag.NewService(client, tagRepository),
		Jwt:     jwtService,
		Encrypt: commons.NewEncryptService(),
	})

	serve(userService, micronService, jwtService, tagRepository)
}

func serve(userService *user.Service, micronService micron.Service, jwtService commons.JwtService, tagRepository *tag.Repository) {
	router := gin.Default()

	registration(router, userService)
	login(router, userService)

	users(router, userService, micronService, jwtService)
	tags(router, tagRepository, userService, jwtService)

	_ = router.Run()
}

func initDb() (*mongo.Client, *context.CancelFunc) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	settings := options.Client().
		ApplyURI(commons.GetEnv("DATABASE_CONNECTION_STRING")).
		SetAuth(
			options.Credential{
				Username: commons.GetEnv("DATABASE_CONNECTION_USERNAME"),
				Password: commons.GetEnv("DATABASE_CONNECTION_PASSWORD"),
			})

	client, err := mongo.Connect(ctx, settings)

	if err != nil {
		panic("Could not connect to database")
	} else {
		log.Println("Successfully connected to database")
	}
	return client, &cancel
}

func registration(router *gin.Engine, userService *user.Service) {
	registration := router.Group("/register")
	{
		registration.POST("/", func(c *gin.Context) { user.HandleUserRegistration(c, userService) })
	}
}

func login(router *gin.Engine, userService *user.Service) {
	login := router.Group("/login")
	{
		login.POST("/", func(c *gin.Context) { user.HandleUserAuthorization(c, userService) })
	}
}

func users(router *gin.Engine, userService *user.Service, micronService micron.Service, jwtService commons.JwtService) {
	users := router.Group("/users/")

	users.Use(middleware.Authorizer(userService, jwtService))
	{
		users.POST("/me/logout", func(c *gin.Context) { user.HandleUserLogout(c, userService) })
		users.GET("/me", func(c *gin.Context) { user.HandleUserByTokenRetrieval(c, userService) })
		users.GET("/me/tags", func(c *gin.Context) { user.HandleUserTagsRetrieval(c, userService) })
		users.POST("/me/tags", func(c *gin.Context) { user.HandleUserTagsAddition(c, userService) })
		users.DELETE("/me/tags", func(c *gin.Context) { user.HandleUserTagsRemoval(c, userService) })
		users.GET("/me/microns", func(c *gin.Context) { micron.HandleMicronRetrieval(c, micronService, userService) })
	}
}

func tags(router *gin.Engine, service *tag.Repository, userService *user.Service, jwtService commons.JwtService) gin.IRoutes {
	return router.GET("/tags", middleware.Authorizer(userService, jwtService), func(c *gin.Context) {
		tag.HandleTagsRetrieval(c, service)
	})
}
