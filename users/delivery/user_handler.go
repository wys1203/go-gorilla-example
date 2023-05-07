package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/wys1203/go-gorilla-example/users/entity"
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

func (h *userHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.userUsecase.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *userHandler) signIn(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Acct string `json:"acct"`
		Pwd  string `json:"pwd"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.userUsecase.Login(creds.Acct, creds.Pwd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *userHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	acct := vars["acct"]

	err := h.userUsecase.Delete(acct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RegisterUserRoutes(router *mux.Router, h *userHandler) {
	router.HandleFunc("/users", h.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/search", h.SearchUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{acct}", h.getUserDetails).Methods(http.MethodGet)
	router.HandleFunc("/users/{acct}", h.deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/signup", h.signUp).Methods(http.MethodPost)
	router.HandleFunc("/signin", h.signIn).Methods(http.MethodPost)
}
