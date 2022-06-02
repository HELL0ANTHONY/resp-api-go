package main

// curl localhost:8080/admin -u admin:secret
// run in terminal: ADMIN_PASSWORD=secret go run main.go
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// define el tipo de dato para devolver un valor
type Coaster struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	InPark       string `json:"inpark"`
	Height       int    `json:"height"`
}

type coastersHanlders struct {
	sync.Mutex
	store map[string]Coaster
}

// h = handlers
func (h *coastersHanlders) get(w http.ResponseWriter, r *http.Request) {
	coasters := make([]Coaster, len(h.store))
	h.Lock()
	i := 0
	for _, coaster := range h.store {
		coasters[i] = coaster
		i++
	}
	h.Unlock()
	jsonBytes, err := json.Marshal(coasters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		panic(err.Error())
	}
}
func (h *coastersHanlders) getRandomCoaster(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, len(h.store))
	h.Lock()
	i := 0
	for id := range h.store {
		ids[i] = id
		i++
	}
	defer h.Unlock()

	var target string
	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if len(ids) == 1 {
		target = ids[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}
	// fmt.Println(target)
	w.Header().Add("location", fmt.Sprintf("/coasters/%s", target))
	w.WriteHeader(http.StatusFound)
}

func (h *coastersHanlders) getCoaster(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parts[2] == "random" {
		h.getRandomCoaster(w, r)
		return
	}
	// fmt.Println("lo que llega por la URL", parts[2])
	h.Lock()
	coaster, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(coaster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		panic(err.Error())
	}
}

func (h *coastersHanlders) coasters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("method not allowed"))
		if err != nil {
			panic(err.Error())
		}
		return
	}
}

func (h *coastersHanlders) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			panic(err.Error())
		}
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		_, _ = w.Write([]byte(fmt.Sprintf(
			"need content-type 'application/json', but got '%s'", ct,
		)))
		return
	}
	var coaster Coaster
	err = json.Unmarshal(bodyBytes, &coaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	}
	coaster.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[coaster.Id] = coaster
	defer h.Unlock()
}

func newCoasterHandlers() *coastersHanlders {
	return &coastersHanlders{
		store: map[string]Coaster{},
		// esto estaba a manera de prueba. Enviaba esto haciendo un GET
		// store: map[string]Coaster{
		// 	"id1": {
		// 		Id:           "id1",
		// 		Name:         "Fury 23",
		// 		Height:       89,
		// 		InPark:       "Carowinds",
		// 		Manufacturer: "B+M",
		// 	},
		// },
	}
}

type adminPortal struct {
	password string
}

func newAdminPortal() *adminPortal {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		panic("required env var ADMIN_PASSWORD not set")
	}
	return &adminPortal{password: password}
}

func (a adminPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("401 - unauthorized"))
		return
	}
	_, _ = w.Write([]byte("<html><h1>Super secret admin portal</h1></html>"))
}

func main() {
	admin := newAdminPortal()
	coastersHandlers := newCoasterHandlers()
	// la función route (ruta, controller)
	http.HandleFunc("/coasters", coastersHandlers.coasters)
	http.HandleFunc("/coaster/", coastersHandlers.getCoaster)
	http.HandleFunc("/admin", admin.handler)

	err := http.ListenAndServe(":8080", nil) // inicializa la app en el P:8080
	if err != nil {
		panic(err.Error())
	}
}
