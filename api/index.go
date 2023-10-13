package index

import (
	"net/http"

	handler "github.com/tom-sherman/git-gateway-vercel"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handler.Handler(w, r)
}
