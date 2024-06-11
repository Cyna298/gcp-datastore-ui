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
)

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

		fmt.Println("\n\n---------------------------------------------")
		fmt.Println(r.Method, r.URL.Path)
		fmt.Println("---------------------------------------------]\n")

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

	headers := make(map[string]service.TableHeader)

	for _, e := range entities {
		for _, x := range e {
			fmt.Println(x)
			if _, ok := headers[x.Name]; !ok {
				headers[x.Name] = service.TableHeader{
					Name: x.Name,
					Type: x.TypeOf,
				}

			}

		}
	}
	headerValues := make([]service.TableHeader, len(headers))
	i := 0
	for _, e := range headers {
		if e.Name == "key" {
			headerValues[0] = e
		} else {
			headerValues[i+1] = e
			i += 1
		}
	}

	view.Entities(headerValues, entities).Render(r.Context(), w)
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

}
