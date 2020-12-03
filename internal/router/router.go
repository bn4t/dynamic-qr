package router

import (
	"git.bn4t.me/bn4t/dynamic-qr/internal/qrcode"
	"github.com/gorilla/mux"
)

func NewRouter(qrHandler *qrcode.QrcodeHandler) *mux.Router {
	r := mux.NewRouter()

	// TODO: csrf

	r.HandleFunc("/create-qr", qrHandler.Create).Methods("POST")
	r.HandleFunc("/manage/{password}", qrHandler.Manage).Methods("GET")
	r.HandleFunc("/manage/{password}", qrHandler.Update).Methods("POST")
	r.HandleFunc("/link/{id}", qrHandler.Redirect).Methods("GET")

	// TODO: index page, assets and templates

	return r
}
