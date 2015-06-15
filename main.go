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
	"time"

	"github.com/daaku/go.httpgzip"
	"github.com/gorilla/mux"
)

const databaseFile = "purchases.db"

var listen = flag.String("listen", "localhost:8080", "host and port to listen on")

func backup() error {
	src, err := os.Open(databaseFile)
	defer src.Close()
	if err != nil {
		return errors.New("could not open database to backup")
	}

	err = os.MkdirAll("backup", 0755)
	if err != nil {
		return errors.New("could not create backup")
	}

	destFile := path.Join("backup", databaseFile+".gz")
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
	r.HandleFunc("/purchases", app.handleAddPurchase).Methods("POST")

	http.Handle("/", logger(httpgzip.NewHandler(r)))

	log.Printf("Serving on %s\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
	// log.Fatal(http.ListenAndServeTLS(*listen, "cert.pem", "key.pem", nil))
}
