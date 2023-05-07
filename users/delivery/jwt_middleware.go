package delivery

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/wys1203/go-gorilla-example/errors"
)

func JWTAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			errors.JSONHandleError(w, errors.NewErrorWrapper(http.StatusUnauthorized, nil, "Authorization header is missing"))
			return
		}

		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			errors.JSONHandleError(w, errors.NewErrorWrapper(http.StatusUnauthorized, nil, "Invalid token format"))
			return
		}

		tokenString := bearerToken[1]
		pubKeyBytes, err := ioutil.ReadFile("public_key.pem")
		if err != nil {
			errors.JSONHandleError(w, errors.NewErrorWrapper(http.StatusInternalServerError, err, "Error reading public key"))
			return
		}

		pubKeyPEM, _ := pem.Decode(pubKeyBytes)
		publicKey, err := x509.ParsePKIXPublicKey(pubKeyPEM.Bytes)
		if err != nil {
			errors.JSONHandleError(w, errors.NewErrorWrapper(http.StatusInternalServerError, err, "Error parsing public key"))
			return
		}

		rsaPublicKey := publicKey.(*rsa.PublicKey)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return rsaPublicKey, nil
		})

		if err != nil {
			errors.JSONHandleError(w, errors.NewErrorWrapper(http.StatusUnauthorized, err, "Error parsing token"))
			return
		}

		if !token.Valid {
			errors.JSONHandleError(w, errors.NewErrorWrapper(http.StatusUnauthorized, err, "Invalid token"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
