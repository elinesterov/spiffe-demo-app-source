package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	bundles, err := a.client.FetchX509Bundles(a.ctx)
	if err != nil {
		str := "Error fetching bundles: " + err.Error()
		log.Printf("%v", str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}

	bundleMap := make(map[string][]string)
	for _, bundle := range bundles.Bundles() {
		trustDomain := bundle.TrustDomain().IDString()
		for _, authority := range bundle.X509Authorities() {
			encoded := base64.StdEncoding.EncodeToString(authority.Raw)
			bundleMap[trustDomain] = append(bundleMap[trustDomain], encoded)
		}
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

func (a *API) JWTProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Log the request headers and parameters and body
	log.Printf("Request headers: %v", r.Header)
	log.Printf("Request parameters: %v", r.URL.Query())

	// Log body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	// ensure the reuquest body can be read again later
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	log.Printf("Request body: %v", string(body))

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

	// Extract `iat` and `exp` values from `Claims` map and calculate `expires_in`.
	iatFloat, okIat := svid.Claims["iat"].(float64)
	expFloat, okExp := svid.Claims["exp"].(float64)
	if !okIat || !okExp {
		log.Printf("failed to parse `iat` or `exp` from claims")
		// handle error
	}

	exp := int64(expFloat - iatFloat)

	// Convert the JWT-SVID to a JSON response
	// following the oauth2 token response format
	response := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}{
		AccessToken: svid.Marshal(),
		TokenType:   "Bearer",
		ExpiresIn:   exp,
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
