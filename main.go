package main

import (
	"backend/service"
	"backend/view"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

func frontendRouting(server *gin.Engine) {
	server.Static("/_next", "./out/_next")
	server.StaticFile("/", "./out/index.html")
	server.StaticFile("/next.svg", "./out/next.svg")
	server.StaticFile("/vercel.svg", "./out/vercel.svg")
}

// func GetAllKindsRoute(client *datastore.Client) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := context.Background()
// 		kinds, err := GetAllKinds(ctx, client)
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"error": err,
// 			})
// 			return
// 		}
// 		c.JSON(200, gin.H{
// 			"kinds": kinds,
// 		})
// 	}
// }

// func GetAllEntitiesRoute(client *datastore.Client) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := context.Background()
// 		kind := c.Param("kind")
// 		//limit and cursor from query params
// 		limit, err := strconv.Atoi(c.Query("limit"))
// 		if err != nil {
// 			limit = 10
// 		}
// 		cursor := c.Query("cursor")
//
// 		//sortKey and sortDirection from query params
// 		sortKey := c.Query("sortKey")
// 		sortDirection := c.Query("sortDirection")
// 		entities, nextCursor, err := GetAllEntities(ctx, client, kind, sortKey, sortDirection, limit, cursor)
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"error": err,
// 			})
// 			return
// 		}
// 		c.JSON(200, gin.H{
// 			"entities":   entities,
// 			"nextCursor": nextCursor,
// 		})
// 	}
// }

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

type APIServer struct {
	listenAddr string
	client     *datastore.Client
}

func (as *APIServer) ServeTempl(w http.ResponseWriter, r *http.Request) error {

	kinds, err := service.GetAllKinds(r.Context(), as.client)
	if err != nil {
		log.Fatal(err)
		return err
	}
	view.Show(kinds).Render(r.Context(), w)
	return nil
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error
type HttpError struct {
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}
func makeHttpHandler(f ApiFunc) http.HandlerFunc {
	// log the endpoint
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)

		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, HttpError{
				Message: err.Error(),
			})
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func (as *APIServer) Detail(w http.ResponseWriter, r *http.Request) error {
	entity := r.URL.Query().Get("entity")

	entities, _, err := service.GetAllEntities(r.Context(), as.client, entity, "", "", 10, "")

	if err != nil {
		return err
	}

	props, err := service.GetAttrs(r.Context(), as.client, entity)
	if err != nil {
		return err
	}

	fmt.Println("ENTITIES")

	for _, e := range entities {
		fmt.Println(e)
		fmt.Println("\n\n")

	}

	view.Entities(props, entities).Render(r.Context(), w)
	return nil
}

func main() {
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
	client, err := service.NewDatastoreClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	as := APIServer{
		client: client,
	}

	router := http.NewServeMux()

	router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	router.HandleFunc("/", makeHttpHandler(as.ServeTempl))
	router.HandleFunc("/detail", makeHttpHandler(as.Detail))
	http.ListenAndServe("localhost:8080", router)

	// if err := server.Run(":" + *port); err != nil {
	// 	log.Fatal("Error starting server: ", err)
	// }

}
