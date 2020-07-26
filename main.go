package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"micron/micron"
	"micron/middleware"
	scraper "micron/scraper/runnable"
	"micron/tag"
	"micron/user"
	"time"
)

func main() {
	client, cancel := initializeDatabaseConnection()
	defer (*cancel)()

	user.MongoClient = client

	scraper.Start(false)

	router := gin.Default()

	registration(router)
	login(router)

	users(router)
	tags(router)

	_ = router.Run()
}

func initializeDatabaseConnection() (*mongo.Client, *context.CancelFunc) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://database:27017").SetAuth(options.Credential{
		Username: "root",
		Password: "root",
	}))

	if err != nil {
		panic("Could not connect to mongodb")
	} else {
		log.Println("Successfully connected to database")
	}
	return client, &cancel
}

func registration(router *gin.Engine) {
	registration := router.Group("/register")
	{
		registration.POST("/", user.HandleUserRegistration)
	}
}

func login(router *gin.Engine) {
	login := router.Group("/login")
	{
		login.POST("/", user.HandleUserAuthorization)
	}
}

func users(router *gin.Engine) {
	users := router.Group("/users/")

	users.Use(middleware.Authorizer())
	{
		users.GET("/me", user.HandleUserByTokenRetrieval)

		users.GET("/me/tags", user.HandleUserTagsRetrieval)
		users.POST("/me/tags", user.HandleUserTagsAddition)
		users.DELETE("/me/tags", user.HandleUserTagsRemoval)

		users.GET("/me/microns", micron.HandleMicronRetrieval)
	}
}

func tags(router *gin.Engine) gin.IRoutes {
	return router.GET("/tags", middleware.Authorizer(), tag.HandleTagsRetrieval)
}
