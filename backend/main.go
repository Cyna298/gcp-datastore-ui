package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

//	func frontendRouting(server *gin.Engine) {
//		server.Static("/assets", "./dist/assets")
//		server.StaticFile("/", "./dist/index.html")
//		server.StaticFile("/vite.svg", "./dist/vite.svg")
//	}
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

		entities, nextCursor, err := GetAllEntities(ctx, client, kind, limit, cursor)
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

func main() {

	ctx := context.Background()
	client, err := NewDatastoreClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	server := gin.Default()

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	frontendRouting(server)
	server.GET("/api/kinds", GetAllKindsRoute(client))
	server.GET("/api/entities/:kind/", GetAllEntitiesRoute(client))

	server.Run(":8080")

	kinds, err := GetAllKinds(ctx, client)
	if err != nil {
		log.Fatalf("Failed to get kinds: %v", err)
	}
	fmt.Println("Kinds:", kinds)

	entities, nextCursor, err := GetAllEntities(ctx, client, "User", 10, "")
	if err != nil {
		log.Fatalf("Failed to get entities: %v", err)
	}

	entitiesJSON, _ := json.Marshal(entities)
	fmt.Println("Entities:", string(entitiesJSON))
	fmt.Println("Next cursor:", nextCursor)

}
