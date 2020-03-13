package ingestv1

import (
	"net/http"
	"shujew/elasticsearch-batcher/config"
)

// AllowAllOrigins is a bool which defines whether the
// app should allow all origins for CORS
var AllowAllOrigins = config.GetAllowAllOrigins()

// AllowedHeaderOrigins is a map which defines which
// origins the app should allow for CORS. It is
// ignored if AllowAllOrigins=true
var AllowedHeaderOrigins = config.GetAllowedOrigins()

// setDefaultHeaders sets the CORS headers on the response
func setDefaultHeaders(w *http.ResponseWriter, req *http.Request) {
	requestOrigin := req.Header.Get("Origin")
	if AllowAllOrigins || AllowedHeaderOrigins[requestOrigin] {
		(*w).Header().Set("Access-Control-Allow-Origin", requestOrigin)
	}
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST OPTIONS")
	(*w).Header().Set("Server", "Elasticsearch-Batcher/1.0")
}
