package router

import (
	"git.bn4t.me/bn4t/dynamic-qr/internal/qrcode"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(qrHandler *qrcode.QrcodeHandler, csrfKey []byte, staticDir string) *mux.Router {
	r := mux.NewRouter()

	r.Use(csrf.Protect(csrfKey))

	r.HandleFunc("/new-qr", qrHandler.Store).Methods("POST")
	r.HandleFunc("/manage/{password}", qrHandler.Manage).Methods("GET")
	r.HandleFunc("/manage/{password}", qrHandler.Update).Methods("POST")
	r.HandleFunc("/link/{id}", qrHandler.Redirect).Methods("GET")
	r.HandleFunc("/", qrHandler.Create).Methods("GET")

	// static files
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir(staticDir)))
	r.PathPrefix("/fonts/").Handler(http.FileServer(http.Dir(staticDir)))

	return r
}
