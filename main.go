package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
}

var client *mongo.Client
var collection *mongo.Collection

func main() {

	connectToMongoDB()
	r := gin.Default()
	baseURL := os.Getenv("baseURL")
	r.GET(baseURL+"/ping", handlePing)
	r.GET(baseURL+"/pong", handlePong)

	// r.GET(baseURL+"/getUsersFromDB", getUsersFromDB)
	// r.GET(baseURL+"/", handleUserEvent)

	r.GET(baseURL+"/getUsers", getUsers)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func connectToMongoDB() {
	// Set client options docker-compose -f docker-compose_myapi.yml up -d

	clientOptions := options.Client().ApplyURI("mongodb://root:example@mongo:27017/")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set the collection
	collection = client.Database("mydatabase").Collection("users")
	log.Println("Connected to MongoDB!")
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func handlePong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ping",
	})
}

func getUsers(c *gin.Context) {
	var users []User

	// Fetch users from MongoDB
	cur, err := collection.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, users)
}
