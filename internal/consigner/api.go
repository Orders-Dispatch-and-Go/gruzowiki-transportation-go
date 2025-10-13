package consigner

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"net/http"
)

const CONSIGNER = "/consigners"

func RegisterHandler(r *routing.Router) http.Handler {
	r.To("POST", CONSIGNER)
}
