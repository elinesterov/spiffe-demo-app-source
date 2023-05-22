package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
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
	start := time.Now()
	svid, err := a.client.FetchJWTSVID(a.ctx, jwtsvid.Params{Audience: "example.org"})
	if err != nil {
		str := "Error fetching JWT SVID: " + err.Error()
		log.Printf("%v", str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}

	elapsed := time.Since(start)
	log.Printf("JWT SVID fetched in %s", elapsed)

	// Convert the JWT-SVID to a JSON response
	response := struct {
		Token string `json:"token"`
	}{
		Token: svid.Marshal(),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response as JSON: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("JWT: %s", svid.Marshal())
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (a *API) GetX509Handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	svid, err := a.client.FetchX509SVID(a.ctx)
	if err != nil {
		str := "Error fetching X509 SVID: " + err.Error()
		log.Printf("%v", str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}
	elapsed := time.Since(start)
	log.Printf("X509 SVID fetched in %s", elapsed)

	// Convert the X509-SVID to a JSON response
	response := struct {
		Cert string `json:"cert"`
	}{
		Cert: base64.StdEncoding.EncodeToString((svid.Certificates[0].Raw)),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response as JSON: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("X.509-SVID: %s", jsonResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (a *API) GetTrustBundleHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	bundles, err := a.client.FetchX509Bundles(a.ctx)
	if err != nil {
		str := "Error fetching bundles: " + err.Error()
		log.Printf("%v", str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}
	elapsed := time.Since(start)
	log.Printf("Bundles fetched in %s", elapsed)

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

	log.Printf("Trust bundle: %s", jsonResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
