package v1

import (
	"net/http"
	"shujew/elasticsearch-batcher/config"
)

var AllowedHeaderOrigins = config.GetAllowedHosts()

func setDefaultHeaders(w *http.ResponseWriter, req *http.Request) {
	requestOrigin := req.Header.Get("Origin")
	if AllowedHeaderOrigins[requestOrigin] {
		(*w).Header().Set("Access-Control-Allow-Origin", requestOrigin)
	}
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST OPTIONS")
	(*w).Header().Set("Server", "Elasticsearch-Batcher/1.0")
}