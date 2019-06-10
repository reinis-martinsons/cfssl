package xxx

import (
	"encoding/json"
	"net/http"

	"github.com/cloudflare/cfssl/api"
)

// Response contains the response to the /xxx API
type Response struct {
	Xxx bool `json:"xxx"`
}

func xxxHandler(w http.ResponseWriter, r *http.Request) error {
	response := api.NewSuccessResponseWithMessage(&Response{Xxx: true}, "Ziņojums", 66213)
	return json.NewEncoder(w).Encode(response)
}

// NewHandler creates a new handler to serve xxx checks.
func NewHandler() http.Handler {
	return api.HTTPHandler{
		Handler: api.HandlerFunc(xxxHandler),
		Methods: []string{"GET"},
	}
}
