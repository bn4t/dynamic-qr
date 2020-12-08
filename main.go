package main

import (
	"context"
	"flag"
	"git.bn4t.me/bn4t/dynamic-qr/internal/qrcode"
	"git.bn4t.me/bn4t/dynamic-qr/internal/router"
	_ "git.bn4t.me/bn4t/dynamic-qr/internal/statik"
	"git.bn4t.me/bn4t/dynamic-qr/internal/template"
	"github.com/rakyll/statik/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var dbPath string
var bindAddr string
var csrfKey string

func main() {
	flag.StringVar(&dbPath, "db", "./qrcodes.db", "Specify the location of the database.")
	flag.StringVar(&bindAddr, "bind", "127.0.0.1:8080", "Specify the address to bind to.")
	flag.StringVar(&csrfKey, "csrf", "foo", "The csrf key to be used.")

	store, err := qrcode.NewSqliteQrcodeStore(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// load up statik fs
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// load the templates
	tmpl, err := template.LoadTemplates(statikFS)
	if err != nil {
		log.Fatal(err)
	}

	r := router.NewRouter(qrcode.NewQrcodeHandler(store, tmpl), []byte(csrfKey), statikFS)

	srv := &http.Server{
		Handler: r,
		Addr:    bindAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Print("started listener on " + bindAddr)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-sigc

	log.Print("stopping gracefully..")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Print(err)
	}
	if err := store.Close(); err != nil {
		log.Print(err)
	}
}
