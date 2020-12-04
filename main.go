package main

import (
	"git.bn4t.me/bn4t/dynamic-qr/internal/qrcode"
	"git.bn4t.me/bn4t/dynamic-qr/internal/router"
	"git.bn4t.me/bn4t/dynamic-qr/internal/template"
	"log"
	"net/http"
	"time"
)

func main() {
	store, err := qrcode.NewSqliteQrcodeStore("./qrcode.db")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.LoadTemplates("./static/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	r := router.NewRouter(qrcode.NewQrcodeHandler(store, tmpl), []byte("test"), "./static/public")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
