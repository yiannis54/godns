package middleware

import (
	"fmt"
	"net/http"
)

func BasicAuthHandler(next http.HandlerFunc, user, pass, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if checkBasicAuth(r, user, pass) {
			next(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func checkBasicAuth(r *http.Request, user, pass string) bool {
	u, p, ok := r.BasicAuth()
	if !ok {
		return false
	}

	return u == user && p == pass
}
