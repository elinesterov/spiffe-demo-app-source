package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

type API struct {
	client *workloadapi.Client
	ctx    context.Context
}

func NewAPI(ctx context.Context, client *workloadapi.Client) (*API, error) {
	return &API{
		client: client,
		ctx:    ctx,
	}, nil
}

func (a *API) GetJwtHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get JWT button pressed"))
}

func (a *API) GetX509Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get X509 button pressed"))
}

func (a *API) GetTrustBundleHandler(w http.ResponseWriter, r *http.Request) {
	bundles, err := a.client.FetchX509Bundles(a.ctx)
	if err != nil {
		str := "Error fetching bundles: " + err.Error()
		w.Write([]byte(str))
		return
	}

	bundleMap := make(map[string]string)
	for _, bundle := range bundles.Bundles() {
		encoded := base64.StdEncoding.EncodeToString(bundle.X509Authorities()[0].Raw)
		bundleMap[bundle.TrustDomain().IDString()] = encoded
	}
	jsonResponse, err := json.MarshalIndent(map[string]interface{}{"bundles": bundleMap}, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// w.Write(jsonResponse)
	log.Printf("Trust bundle: %s", jsonResponse)

	w.Write([]byte("Get Trust Bundle button pressed"))
}
