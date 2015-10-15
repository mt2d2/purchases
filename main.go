package main

import (
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/daaku/go.httpgzip"
	"github.com/gorilla/mux"
)

var listen = flag.String("listen", "localhost:8080", "host and port to listen on")
var db = flag.String("db", "purchases.db", "sqlite3 database file")

func backup() error {
	src, err := os.Open(*db)
	defer src.Close()
	if err != nil {
		return errors.New("could not open database to backup")
	}

	backupPath := path.Join(filepath.Dir(*db), "backup")
	err = os.MkdirAll(backupPath, 0755)
	if err != nil {
		return errors.New("could not create backup")
	}

	destFile := path.Join(backupPath, filepath.Base(*db)+".gz")
	dest, err := os.Create(destFile)
	defer dest.Close()
	if err != nil {
		return err
	}

	gzipWriter := gzip.NewWriter(dest)
	_, err = io.Copy(gzipWriter, src)
	if err != nil {
		return err
	}

	return gzipWriter.Close()
}

func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func main() {
	flag.Parse()

	err := backup()
	if err != nil {
		log.Panicln(err)
	}
	log.Println("backup complete")

	app := newApp()
	defer app.destroy()
	log.Println("database opened")

	r := mux.NewRouter()
	r.HandleFunc("/purchases", app.handlePurchases).Methods("GET")
	r.HandleFunc("/purchases/{id:[0-9]+}", app.handlePurchase).Methods("GET")
	r.HandleFunc("/purchases/{id:[0-9]+}", app.handleDelete).Methods("DELETE")
	r.HandleFunc("/purchases", app.handleAddPurchase).Methods("POST")

	http.Handle("/", logger(httpgzip.NewHandler(r)))

	log.Printf("Serving on %s\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
