package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	svid, err := a.client.FetchJWTSVID(a.ctx, jwtsvid.Params{Audience: "example.org"})
	if err != nil {
		str := "Error fetching JWT SVID: " + err.Error()
		log.Printf("%v", str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}

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

	// TODO: hardcoded respose for now
	// jsonResponse := []byte(`{"token":"eyJhbGciOiJSUzI1NiIsImtpZCI6IlY4cGZoUXliVUdUYlk3eXY5WmdzVEpjc3RKWjdxcG9zIiwidHlwIjoiSldUIn0.eyJhdWQiOlsiZXhhbXBsZS5vcmciXSwiZXhwIjoxNjc4NjczNTU1LCJpYXQiOjE2Nzg2Njk5NTUsImlzcyI6Im9pZGMtZGlzY292ZXJ5LmV4YW1wbGUub3JnIiwic3ViIjoic3BpZmZlOi8vZXhhbXBsZS5vcmcvbnMvc3BpZmZlLWRlbW8vc2EvZGVmYXVsdCJ9.CfJaKAH31qMO0so9ivRj2Qv9SplrOnuwG5Ar88VxokA8osLf6_-imKryWjYwkwDt2eoolIVQiz7kqDBYgIbeTPLjskNRrl2W2jw4aHJHSxXtELN_GDFtfCW_9U4_FcHQ2ORM26i2WnCDlgSuR2ALBJta-8uMlu-OIV_O_g9cJzrjjiZ11R8wAqHy02L1dBGcI1uzCqi9Judhstt5GLux9pOWtnRtdqoSL9bUa5acNaBajGhy1ZuISESYwDKmUNkFKDw1GKUFvQMAHCl69w49TSWrg3ZcBwMTPTaUdVJ6qIc3QNv-YRIpDGy69E0_1M3DwPldwHEN1JJR1-Sg30M2rQ"}`)
	log.Printf("JWT: %s", svid.Marshal())
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (a *API) GetX509Handler(w http.ResponseWriter, r *http.Request) {
	svid, err := a.client.FetchX509SVID(a.ctx)
	if err != nil {
		str := "Error fetching X509 SVID: " + err.Error()
		log.Printf("%v", str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}

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

	// TODO: hardcoded respose for now
	// jsonResponse := []byte(`{"cert":"MIIC+TCCAeGgAwIBAgIRAI65Uy+8yQzyvNIg7LGE7sEwDQYJKoZIhvcNAQELBQAwNTELMAkGA1UEBhMCTkwxEDAOBgNVBAoTB0V4YW1wbGUxFDASBgNVBAMTC2V4YW1wbGUub3JnMB4XDTIzMDMxMzAxNDQzNFoXDTIzMDMxMzA1NDQ0NFowSDELMAkGA1UEBhMCVVMxDjAMBgNVBAoTBVNQSVJFMSkwJwYDVQQtEyBjZWZmNzFmMjVhODFiNmZlZDhiZjRmN2YzY2UyMDJlMTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABMLafsUpyFMrOzhlxXclY4WshIqpLTz4KHLO4vViLsjR3t2DOkufQcrXl+Q67AzuZ3rbJeB0bshiqfeba3JyGqSjgbswgbgwDgYDVR0PAQH/BAQDAgOoMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBTZK8j0R1SEliGk93nWR+VZIpTTOTAfBgNVHSMEGDAWgBQdko53hZNVodUv6T1crLYSlxc3ATA5BgNVHREEMjAwhi5zcGlmZmU6Ly9leGFtcGxlLm9yZy9ucy9zcGlmZmUtZGVtby9zYS9kZWZhdWx0MA0GCSqGSIb3DQEBCwUAA4IBAQBeBd07W9WYYGlbA0U/KUKYxEgupFQN65ksEe5of14s2k0BwykKm0lyZoYeKJBDetE61nExHg+RJPuU5VQIM1tzJG46SORJFui3xlJ4jD0oHr2CDNso7NGsR1fIwSxb7j6RZU9nufC0KvarTTBeaP8HNYxB74SqDcKtmihzFPEpfndWMxzBjfs9owVHjA8+e49sLGQL5rZJTyI8Xx03WkxLa61v/EYEAdm6u8oaiEv4Iq4a2zcw4VaZmay75ukDmsCoVRg8sZHoFjor6wGqRyuKdFrXH07nImkHYFksqE7yvDrp1+dzdDf9Ac5oSLLSTvlgukLvhI34pV0biTYRuOrJ"}`)
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
	// TODO: hardcoded respose for now
	// jsonResponse := []byte(`{"bundles":{"spiffe://example.org":"MIIDWDCCAkCgAwIBAgIRAOHjXwKDm65fwZvVEsKHwWowDQYJKoZIhvcNAQELBQAwNTELMAkGA1UEBhMCTkwxEDAOBgNVBAoTB0V4YW1wbGUxFDASBgNVBAMTC2V4YW1wbGUub3JnMB4XDTIzMDMwOTIxNTQ1MFoXDTIzMDMxMDIxNTUwMFowNTELMAkGA1UEBhMCTkwxEDAOBgNVBAoTB0V4YW1wbGUxFDASBgNVBAMTC2V4YW1wbGUub3JnMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAm2X4CAYhjwrN+SjJOdL7bdtm1u94C9w/RCJ6vEJHMqXlRWWFEXCA7fQP2QPcVCG5OaQkFmFWWuoMYd/laXyWq2S8TYXvCkl/mC6pRwSIYEf2p4Wy6bLab0g5FcgYgMiO0y+lvWpPPtigMKB7ebmeFIFgXhwwY5F7bFHvNu+qY7bX5VnsHJyg9Nht70ubLNfXdIfe8Aemj2v+u2pQTap3Ttz74B1jPTUWcg/c9DoyxrPBzW+qSY7SfnRI8UXv2/88sVutIvAgcWUyW1U0OmbPMbdLPgb5rQi9m2JC8NDqDfROzOiVIkBBak8RoyCti01WaLTvnuzVTYhhONpjt0N12QIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUv3faet48gsuHXgYWKRLycTn5kbEwHwYDVR0RBBgwFoYUc3BpZmZlOi8vZXhhbXBsZS5vcmcwDQYJKoZIhvcNAQELBQADggEBAEzSotmqYcSiZBW279p3DKlXM2oFmaQjGxmo6t/e3NXDKM5RHVB0PDevTFVIcE0ph65J1+MEIagq3gDUS5OCUKQMHuDQ0GgViNJbO66LyrPkuavIg7sSHZMhTKAJDoF7LuebwcSuVlEEglKbTNWGYYzxNqckOkAVSnJPtFOlsxtdk+r3zOORAsXZ3+XxVMPQZ4WGIF8uPyFHvNm3noL0XZmhKLiPlYWIvHow59LG2Lz9zYwCK2+OjVeVVayxzRzHOXDvAs0CfTLU4Lx38m0CBEKPxYrr4swYWSeAj4BiiAOT4lZg4uyko0w5CfoBvKBBvfl6z46QHhNzjWDqUDFthAw="}}`)
	log.Printf("Trust bundle: %s", jsonResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
