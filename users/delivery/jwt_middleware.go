package delivery

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func JWTAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]
		pubKeyBytes, err := ioutil.ReadFile("public_key.pem")
		if err != nil {
			http.Error(w, "Error reading public key", http.StatusInternalServerError)
			return
		}

		pubKeyPEM, _ := pem.Decode(pubKeyBytes)
		publicKey, err := x509.ParsePKIXPublicKey(pubKeyPEM.Bytes)
		if err != nil {
			http.Error(w, "Error parsing public key", http.StatusInternalServerError)
			return
		}

		rsaPublicKey := publicKey.(*rsa.PublicKey)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return rsaPublicKey, nil
		})

		if err != nil {
			http.Error(w, "Error parsing token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
