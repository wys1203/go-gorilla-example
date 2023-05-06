package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/wys1203/go-gorilla-example/users/usecase"
)

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *userHandler {
	return &userHandler{userUsecase: userUsecase}
}

func (h *userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *userHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	fullname := r.URL.Query().Get("fullname")

	users, err := h.userUsecase.SearchUsers(fullname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *userHandler) getUserDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	acct := vars["acct"]

	user, err := h.userUsecase.GetUserByAcct(acct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func RegisterUserRoutes(router *mux.Router, h *userHandler) {
	router.HandleFunc("/users", h.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/search", h.SearchUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{acct}", h.getUserDetails).Methods(http.MethodGet)
}
