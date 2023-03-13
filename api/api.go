package api

import (
	"context"
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
	w.Write([]byte("Get JWT button pressed Get JWT button pressedGet JWT button pressedGet JWT button pressedGet JWT button pressedGet JWT button pressedGet JWT button pressedGet JWT button pressed Get JWT button pressedGet JWT button pressedGet JWT button pressedGet JWT button pressedGet JWT button pressedGet JWT button pressed"))
}

func (a *API) GetX509Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get X509 button pressedGet X509 button pressedGet X509 button pressedGet X509 button pressedGet X509 button pressedGet X509 button pressed Get X509 button pressedGet X509 button pressedGet X509 button pressedGet X509 button pressedGet X509 button pressedGet X509 button pressed"))
}

func (a *API) GetTrustBundleHandler(w http.ResponseWriter, r *http.Request) {
	// bundles, err := a.client.FetchX509Bundles(a.ctx)
	// if err != nil {
	// 	str := "Error fetching bundles: " + err.Error()
	// 	log.Printf("%v", str)
	// 	http.Error(w, str, http.StatusInternalServerError)
	// 	return
	// }

	// bundleMap := make(map[string]string)
	// for _, bundle := range bundles.Bundles() {
	// 	encoded := base64.StdEncoding.EncodeToString(bundle.X509Authorities()[0].Raw)
	// 	bundleMap[bundle.TrustDomain().IDString()] = encoded
	// }
	// jsonResponse, err := json.MarshalIndent(map[string]interface{}{"bundles": bundleMap}, "", "  ")
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
	// 	return
	// }
	// TODO: hardcoded respose for now
	jsonResponse := []byte(`{"bundles":{"spiffe://example.org":"MIIDWDCCAkCgAwIBAgIRAOHjXwKDm65fwZvVEsKHwWowDQYJKoZIhvcNAQELBQAwNTELMAkGA1UEBhMCTkwxEDAOBgNVBAoTB0V4YW1wbGUxFDASBgNVBAMTC2V4YW1wbGUub3JnMB4XDTIzMDMwOTIxNTQ1MFoXDTIzMDMxMDIxNTUwMFowNTELMAkGA1UEBhMCTkwxEDAOBgNVBAoTB0V4YW1wbGUxFDASBgNVBAMTC2V4YW1wbGUub3JnMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAm2X4CAYhjwrN+SjJOdL7bdtm1u94C9w/RCJ6vEJHMqXlRWWFEXCA7fQP2QPcVCG5OaQkFmFWWuoMYd/laXyWq2S8TYXvCkl/mC6pRwSIYEf2p4Wy6bLab0g5FcgYgMiO0y+lvWpPPtigMKB7ebmeFIFgXhwwY5F7bFHvNu+qY7bX5VnsHJyg9Nht70ubLNfXdIfe8Aemj2v+u2pQTap3Ttz74B1jPTUWcg/c9DoyxrPBzW+qSY7SfnRI8UXv2/88sVutIvAgcWUyW1U0OmbPMbdLPgb5rQi9m2JC8NDqDfROzOiVIkBBak8RoyCti01WaLTvnuzVTYhhONpjt0N12QIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUv3faet48gsuHXgYWKRLycTn5kbEwHwYDVR0RBBgwFoYUc3BpZmZlOi8vZXhhbXBsZS5vcmcwDQYJKoZIhvcNAQELBQADggEBAEzSotmqYcSiZBW279p3DKlXM2oFmaQjGxmo6t/e3NXDKM5RHVB0PDevTFVIcE0ph65J1+MEIagq3gDUS5OCUKQMHuDQ0GgViNJbO66LyrPkuavIg7sSHZMhTKAJDoF7LuebwcSuVlEEglKbTNWGYYzxNqckOkAVSnJPtFOlsxtdk+r3zOORAsXZ3+XxVMPQZ4WGIF8uPyFHvNm3noL0XZmhKLiPlYWIvHow59LG2Lz9zYwCK2+OjVeVVayxzRzHOXDvAs0CfTLU4Lx38m0CBEKPxYrr4swYWSeAj4BiiAOT4lZg4uyko0w5CfoBvKBBvfl6z46QHhNzjWDqUDFthAw="}}`)
	log.Printf("Trust bundle: %s", jsonResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
