package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	defaultPort = 8080
)

type Catalog struct {
	Services []Service `json:"services"`
}

type Service struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Bindable bool   `json:"bindable"`
	Plans    []Plan `json:"plans"`
}

type Plan struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var (
	catalog = Catalog{
		Services: []Service{
			Service{
				Id:       "b7fb93d8-3d7f-4509-9730-9ee61acf14a5",
				Name:     "fake-service",
				Bindable: false,
				Plans:    []Plan{},
			},
		},
	}
)

func CatalogHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	catalogData, _ := json.Marshal(catalog)
	w.Write(catalogData)
}

func main() {
	var port int
	if os.Getenv("PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	} else {
		port = defaultPort
	}

	router := httprouter.New()
	router.GET("/v2/catalog", CatalogHandler)

	log.Printf("Starting fake broker on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), contentTypeHandler{router}))
}

type contentTypeHandler struct {
	handler http.Handler
}

func (cth contentTypeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cth.handler.ServeHTTP(w, r)
}
