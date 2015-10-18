package main

import (
	"compress/gzip"
	"errors"
	"flag"
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

var fListen = flag.String("listen", "localhost:8080", "host and port to listen on")
var fDB = flag.String("db", "purchases.db", "sqlite3 database file")
var fUsername = flag.String("username", "test", "username for basic auth")
var fPassword = flag.String("password", "password", "password for basic auth")

func backup() error {
	src, err := os.Open(*fDB)
	defer src.Close()
	if err != nil {
		return errors.New("could not open database to backup")
	}

	backupPath := path.Join(filepath.Dir(*fDB), "backup")
	err = os.MkdirAll(backupPath, 0755)
	if err != nil {
		return errors.New("could not create backup")
	}

	destFile := path.Join(backupPath, filepath.Base(*fDB)+".gz")
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

func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !(username == *fUsername &&
			password == *fPassword) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
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

	http.Handle("/", logger(basicAuth(httpgzip.NewHandler(r))))

	log.Printf("Serving on %s\n", *fListen)
	log.Fatal(http.ListenAndServe(*fListen, nil))
}
