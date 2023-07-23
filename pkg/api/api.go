package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"module36a/pkg/storage"

	"github.com/gorilla/mux"
)

type API struct {
	r  *mux.Router
	db storage.Interfase
}

// Конструктор API
func New(db storage.Interfase) *API {
	api := API{}
	api.db = db
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	// получить n последних новостей
	api.r.HandleFunc("/news/{n}", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	// веб-приложение
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// Router возвращает маршрутизатор запросов
func (api *API) Router() *mux.Router {
	return api.r
}

// Получение последних n новостей
func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["n"]
	n, _ := strconv.Atoi(s)
	if n == 0 {
		n = 10
	}
	posts, err := api.db.News(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
	w.WriteHeader(http.StatusOK)
}
