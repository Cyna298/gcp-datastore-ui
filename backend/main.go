package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func frontendRouting(server *gin.Engine) {
	server.Static("/_next", "./out/_next")
	server.StaticFile("/", "./out/index.html")
	server.StaticFile("/next.svg", "./out/next.svg")
	server.StaticFile("/vercel.svg", "./out/vercel.svg")
}

func GetAllKindsRoute(client *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		kinds, err := GetAllKinds(ctx, client)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(200, gin.H{
			"kinds": kinds,
		})
	}
}

func GetAllEntitiesRoute(client *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		kind := c.Param("kind")
		//limit and cursor from query params
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		}
		cursor := c.Query("cursor")

		//sortKey and sortDirection from query params
		sortKey := c.Query("sortKey")
		sortDirection := c.Query("sortDirection")
		entities, nextCursor, err := GetAllEntities(ctx, client, kind, sortKey, sortDirection, limit, cursor)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(200, gin.H{
			"entities":   entities,
			"nextCursor": nextCursor,
		})
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	files, err := content.ReadDir("out")
	if err != nil {
		fmt.Println("Error reading embedded directory:", err)
	} else {
		fmt.Println("Files in embedded directory:")
		for _, file := range files {
			fmt.Println(" - ", file.Name())
		}
	}
	port := flag.String("port", "8080", "Port for the emulator")
	projectId := flag.String("project", "my-project", "Project ID for the emulator")
	emulatorHost := flag.String("emuHost", "localhost:8081", "Host for the emulator")
	emulatorHostPath := flag.String("emuHostPath", "localhost:8081/datastore", "Host path for the emulator")
	datastoreHost := flag.String("dsHost", "http://localhost:8081", "Host for the datastore")

	flag.Usage = usage

	flag.Parse()

	os.Setenv("DATASTORE_DATASET", *projectId)
	os.Setenv("DATASTORE_EMULATOR_HOST", *emulatorHost)
	os.Setenv("DATASTORE_EMULATOR_HOST_PATH", *emulatorHostPath)
	os.Setenv("DATASTORE_HOST", *datastoreHost)
	os.Setenv("DATASTORE_PROJECT_ID", *projectId)

	fmt.Println("Starting server on port:", port)

	ctx := context.Background()
	client, err := NewDatastoreClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	server := gin.Default()
	server.Use(CORSMiddleware())

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	frontendRouting(server)
	server.GET("/api/kinds", GetAllKindsRoute(client))
	server.GET("/api/entities/:kind/", GetAllEntitiesRoute(client))

	if err := server.Run(":" + *port); err != nil {
		log.Fatal("Error starting server: ", err)
	}

}
