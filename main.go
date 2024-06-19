package main

import (
	"backend/service"
	"backend/view"
	"backend/viewmodel"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type APIServer struct {
	listenAddr string
	vm         *viewmodel.TableViewModel
}

func (as *APIServer) ServeTempl(w http.ResponseWriter, r *http.Request) error {

	entity := r.URL.Query().Get("entity")
	cursor := r.URL.Query().Get("cursor")
	sortKey := r.URL.Query().Get("sortKey")

	as.vm.UpdateKinds(r.Context())

	if entity != "" {
		as.vm.SelectKind(entity)
	}
	as.vm.GetData(r.Context(), cursor, sortKey)

	view.Show(as.vm).Render(r.Context(), w)
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

	vm := viewmodel.NewTableViewModel(client)

	as := APIServer{vm: vm}

	router := http.NewServeMux()

	router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	router.HandleFunc("/", makeHttpHandler(as.ServeTempl))
	http.ListenAndServe("localhost:8080", router)

}
