package main

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/ihsw/go-home/app/DefaultManager"
	"github.com/ihsw/go-home/app/PostManager"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
)

type customReader struct {
	*bytes.Buffer
}

func (r customReader) Close() error { return nil }

func getHandler(db *sqlx.DB) *mux.Router {
	r := mux.NewRouter()

	// route handlers
	r = PostManager.Init(r, db)
	r = DefaultManager.Init(r)
	return r
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// read the request body for logging
		var (
			body []byte
			err  error
		)
		if body, err = ioutil.ReadAll(req.Body); err != nil {
			log.WithFields(log.Fields{
				"url": req.URL,
				"err": err.Error(),
			}).Error("Could not read request body")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not read request body")
			return
		}

		// re-adding the request body
		req.Body = customReader{bytes.NewBuffer(body)}

		// logging the request body
		log.WithFields(log.Fields{
			"url":  req.URL,
			"body": string(body),
		}).Info("Url hit")

		// passing onto the next middleware
		h.ServeHTTP(w, req)
	})
}

func main() {
	log.Info("Starting up")

	db, err := sqlx.Connect(
		"postgres",
		"postgres://postgres@db/postgres?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = http.ListenAndServe(":80", loggingMiddleware(getHandler(db)))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Could not listen on 80")
	}
}
