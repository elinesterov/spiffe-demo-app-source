package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/spirl/spiffe-demo-app/api"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

const (
	// Workload API unix socket path
	WorkloadAPIPath = "unix:///tmp/spirl/spiffe.sock"
)

func main() {

	// Initialize SPIFFE Workload API client
	ctx := context.Background()
	//  use SPIFFE_ENDPOINT_SOCKET environment variable that points
	// to the Workload API socket instead of specifying it directly
	// clientOptions := workloadapi.WithAddr(WorkloadAPIPath)
	// client, err := workloadapi.New(ctx, clientOptions)
	client, err := workloadapi.New(ctx)
	if err != nil {
		log.Fatalf("Unable to create client: %v", err)
	}
	defer client.Close()

	// Initialize API
	api, err := api.NewAPI(ctx, client)
	if err != nil {
		log.Fatalf("Unable to create API: %v", err)
	}

	mux := http.NewServeMux()

	// Serve index.html file at root path
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexPath := path.Join(os.Getenv("KO_DATA_PATH"), "index.html")
		http.ServeFile(w, r, indexPath)
	})

	// Serve static files from the "public" directory
	mux.Handle("/kodata/", http.StripPrefix("/kodata/", http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH")))))

	// Serve the API endpoints
	mux.HandleFunc("/api/gettrustbundle", api.GetTrustBundleHandler)
	mux.HandleFunc("/api/getjwtsvid", api.GetJwtHandler)
	mux.HandleFunc("/api/getx509svid", api.GetX509Handler)

	handler := loggingMiddleware(mux)
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Printf("Server listening on %s", server.Addr)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe: %v", err)
	}

}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
