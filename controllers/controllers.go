package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Coaster struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	InPark       string `json:"inpark"`
	Height       int    `json:"height"`
}

type CoastersHandlers struct {
	sync.Mutex
	store map[string]Coaster
}

func NewCoasterHandlers() *CoastersHandlers {
	return &CoastersHandlers{
		store: map[string]Coaster{},
	}
}

func (h *CoastersHandlers) Get(w http.ResponseWriter, r *http.Request) {
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

func (h *CoastersHandlers) GetRandomCoaster(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Add("location", fmt.Sprintf("/coasters/%s", target))
	w.WriteHeader(http.StatusFound)
}

func (h *CoastersHandlers) GetCoaster(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if parts[2] == "random" {
		h.GetRandomCoaster(w, r)
		return
	}
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

func (h *CoastersHandlers) Post(w http.ResponseWriter, r *http.Request) {
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

func (h *CoastersHandlers) Coasters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.Get(w, r)
		return
	case "POST":
		h.Post(w, r)
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
