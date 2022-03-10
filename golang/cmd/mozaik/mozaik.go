package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gergelymark/mozaik/pkg/config"
	"github.com/gergelymark/mozaik/pkg/mozaik"
	"github.com/gergelymark/mozaik/pkg/spa"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Creating directory for mozaikz")
	os.MkdirAll(config.Config.BasePath, 0755)

	legoColors, err := mozaik.LoadLegoColors("assets/colors.json")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(func(h http.Handler) http.Handler {
		return handlers.CompressHandler(h)
	})
	r.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	})
	r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

	// Colors API
	r.HandleFunc("/api/colors/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		colors := map[string]mozaik.Color{}

		for _, color := range legoColors {
			hex := fmt.Sprintf("#%02x%02x%02x", color.R, color.G, color.B)
			colors[hex] = color
		}

		resp, err := json.Marshal(colors)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(resp)
	}).Methods("GET")

	// Mozaik Image API
	r.HandleFunc("/api/images/{name}", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		vars := mux.Vars(r)
		name := vars["name"]

		fileBytes, err := ioutil.ReadFile(path.Join(config.Config.BasePath, name, "original.png"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileBytes)
	}).Methods("GET")

	// Mozaik API
	// List of mozaiks
	r.HandleFunc("/api/mozaik/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		mozaiks := mozaik.Mozaiks{}
		files, err := ioutil.ReadDir(config.Config.BasePath)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if f.IsDir() {
				moz, err := mozaik.Load(path.Join(f.Name(), "mozaik.json"))
				if err != nil {
					log.Println(err)
					continue
				}
				moz.Image = []mozaik.Colors{}
				mozaiks = append(mozaiks, *moz)
			}
		}

		resp, err := json.Marshal(mozaiks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write(resp)
	}).Methods("GET")

	// Create mozaik
	r.HandleFunc("/api/mozaik/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var mozaikData mozaik.MozaikData
		err := json.NewDecoder(r.Body).Decode(&mozaikData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mozaik, err := mozaikData.ToMozaik(legoColors)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mozaikJSON, err := mozaik.ToJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = mozaik.Save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(201)
		w.Write(mozaikJSON)
	}).Methods("POST")

	// Update mozaik
	r.HandleFunc("/api/mozaik/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var mozaik mozaik.Mozaik
		err := json.NewDecoder(r.Body).Decode(&mozaik)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = mozaik.Save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(201)
	}).Methods("PUT")

	// Get mozaik
	r.HandleFunc("/api/mozaik/{name}", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		vars := mux.Vars(r)
		name := vars["name"]

		fileBytes, err := ioutil.ReadFile(path.Join(config.Config.BasePath, name, "mozaik.json"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileBytes)
	}).Methods("GET")

	// Delete mozaik
	r.HandleFunc("/api/mozaik/{name}", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		vars := mux.Vars(r)
		name := vars["name"]

		err := os.RemoveAll(path.Join(config.Config.BasePath, name))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")

	handler := spa.SpaHandler("static", "index.html")
	r.PathPrefix("/").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(rw, r)
	})

	webserver := &http.Server{
		Addr:         ":5001",
		Handler:      r,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	webserver.ListenAndServe()

}
