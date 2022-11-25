package handler

import (
	"net/http"
)

func HTTPInterceptor(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")
			if len(username) < 3 || IsTokenValid(token) {
				http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
				return
			}
			handler(w, r)
		})
}
