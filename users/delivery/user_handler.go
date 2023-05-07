package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/wys1203/go-gorilla-example/users/entity"
	"github.com/wys1203/go-gorilla-example/users/usecase"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type userHandler struct {
	userUsecase usecase.UserUsecase

	clients     map[*websocket.Conn]bool
	clientsLock sync.Mutex
}

func NewUserHandler(userUsecase usecase.UserUsecase) *userHandler {
	return &userHandler{
		userUsecase: userUsecase,
		clients:     make(map[*websocket.Conn]bool),
	}
}

func (h *userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	sortBy := r.URL.Query().Get("sortBy")
	order := r.URL.Query().Get("order")
	users, err := h.userUsecase.GetAll(page, size, sortBy, order)
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
		h.broadcastFailedSignIn(creds.Acct)
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

func (h *userHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	acct := vars["acct"]

	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.userUsecase.Update(acct, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *userHandler) updateUserFullname(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	acct := vars["acct"]

	var req struct {
		Fullname string `json:"fullname"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.userUsecase.UpdateFullname(acct, req.Fullname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *userHandler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer func() {
		conn.Close()
		h.removeClient(conn)
	}()

	h.addClient(conn)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}
		log.Printf("Received message: %s\n", msg)
	}
}

func (h *userHandler) addClient(conn *websocket.Conn) {
	h.clientsLock.Lock()
	defer h.clientsLock.Unlock()
	h.clients[conn] = true
}

func (h *userHandler) removeClient(conn *websocket.Conn) {
	h.clientsLock.Lock()
	defer h.clientsLock.Unlock()
	delete(h.clients, conn)
}

func (h *userHandler) broadcastFailedSignIn(acct string) {
	h.clientsLock.Lock()
	defer h.clientsLock.Unlock()

	message := struct {
		Type string `json:"type"`
		Acct string `json:"acct"`
	}{
		Type: "failed_sign_in",
		Acct: acct,
	}

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to marshal message:", err)
		return
	}

	for client := range h.clients {
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Failed to send message to client:", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}

func (h *userHandler) RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/signup", h.signUp).Methods(http.MethodPost)
	router.HandleFunc("/signin", h.signIn).Methods(http.MethodPost)
	router.HandleFunc("/ws", h.wsHandler).Methods(http.MethodGet)

	protectedRouter := router.PathPrefix("/users").Subrouter()
	protectedRouter.Use(JWTAuthenticationMiddleware)
	protectedRouter.HandleFunc("", h.GetAllUsers).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/search", h.SearchUsers).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/{acct}", h.getUserDetails).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/{acct}", h.deleteUser).Methods(http.MethodDelete)
	protectedRouter.HandleFunc("/{acct}", h.updateUser).Methods(http.MethodPatch)
	protectedRouter.HandleFunc("/{acct}/fullname", h.updateUserFullname).Methods(http.MethodPut)

}
