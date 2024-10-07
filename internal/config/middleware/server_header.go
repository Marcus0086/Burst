package middleware

import "net/http"

func ServerHeaderMiddleware(writer http.ResponseWriter) {
	writer.Header().Set("Server", "Burst")
}
