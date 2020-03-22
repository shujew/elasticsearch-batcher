package ingestv1

import (
	"fmt"
	"github.com/shujew/elasticsearch-batcher/config"
	"net/http"
	"time"
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
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Server-Name, X-Server-Unix-Timestamp")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST OPTIONS")
	(*w).Header().Set("X-Server-Name", "Elasticsearch-Batcher/1.0")
	(*w).Header().Set("X-Server-Unix-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
}
