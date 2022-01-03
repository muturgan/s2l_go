package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/muturgan/s2l_go/src/config"
	"github.com/muturgan/s2l_go/src/dal"
	"github.com/muturgan/s2l_go/src/models"
)

const _UNDERSCORE = "_"
const _EMPTY = ""

type server struct {
	dal dal.IDal
}

func newServer(dal dal.IDal) *server {
	return &server{dal: dal}
}

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var parts [3]string
	copy(parts[:], strings.Split(r.URL.Path, "/"))
	underscoreOrHash := parts[1]
	hashOrEmpty := parts[2]

	var hash string
	if underscoreOrHash != _UNDERSCORE && underscoreOrHash != _EMPTY {
		hash = underscoreOrHash
	} else if hashOrEmpty != _EMPTY {
		hash = hashOrEmpty
	} else {
		fmt.Println("something wrong with the url")
		fmt.Println(r.URL.Path)
		http.Error(w, "Invalid URL (I don't know why)", http.StatusInternalServerError)
		return
	}

	link, err := s.dal.GetLinkByHash(hash)
	if err != nil {
		fmt.Println("GetLinkByHash error!")
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if link == nil {
		fmt.Println("Not found!")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", link.Link)
	w.WriteHeader(http.StatusMovedPermanently)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) compressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var cr models.CompressRequest
	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		errMessage := "Incorrect request body. It should be a valid json-serialized object with a \"link\" field which is a valid url"
		http.Error(w, errMessage, http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(cr.Link)
	if err != nil {
		errMessage := "Incorrect request body. It should be a valid json-serialized object with a \"link\" field which is a valid url"
		http.Error(w, errMessage, http.StatusBadRequest)
		return
	}

	newShortLink, err := s.dal.CreateNewLink(cr.Link)
	if err != nil {
		fmt.Println("createNewLink error!")
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newShortLink)
}

func Serve(config *config.Config, dal dal.IDal) {
	s := newServer(dal)

	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/compress", s.compressHandler)
	mux.HandleFunc("/", s.indexHandler)

	fmt.Println("ok let's try to start at http://localhost" + config.GetServingAddress())

	err := http.ListenAndServe(config.GetServingAddress(), mux)
	if err != nil {
		log.Fatal(err)
	}
}
