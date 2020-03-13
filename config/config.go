// Package config provides a simple interface for the app
// to access values defined by environment variables
package config

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

// Constants below define the environment variable key
// used by the specified function defined in var name
const (
	debugEnvVar           string = "ESB_DEBUG"
	httpPortEnvVar        string = "ESB_HTTP_PORT"
	allowAllOriginsEnvVar string = "ESB_ALLOW_ALL_ORIGINS"
	allowedOriginsEnvVar  string = "ESB_ALLOWED_ORIGINS"
	esHostEnvVar          string = "ESB_ES_HOST"
	esUsernameEnvVar      string = "ESB_ES_USERNAME"
	esPasswordEnvVar      string = "ESB_ES_PASSWORD"
	esTimeoutEnvVar       string = "ESB_ES_TIMEOUT_SECONDS"
	esFlushIntervalEnvVar string = "ESB_FLUSH_INTERVAL_SECONDS"
)

// GetLogLevel returns the log level to be used throughout
// the app
func GetLogLevel() log.Level {
	value := getEnvValue(debugEnvVar, "OFF")
	if value == "ON" {
		return log.TraceLevel
	}
	return log.InfoLevel
}

// GetHTTPPort returns the http port on which the app
// should run on
func GetHTTPPort() string {
	return getEnvValue(httpPortEnvVar, "8889")
}

// GetAllowAllOrigins returns whether the app should allow all
// origins or not (CORS)
func GetAllowAllOrigins() bool {
	value := getEnvValue(allowAllOriginsEnvVar, "true")
	return value == "true"
}

// GetAllowedOrigins returns all origins which the app should
// allow assuming GetAllowAllOrigins() == false  (CORS)
func GetAllowedOrigins() map[string]bool {
	allowedHosts := map[string]bool{}
	value := getEnvValue(allowedOriginsEnvVar, "")
	for _, allowedHost := range strings.Split(value, ",") {
		if len(allowedHost) > 0 {
			allowedHosts[allowedHost] = true
		}
	}
	return allowedHosts
}

// GetESHost returns the protocol + hostname of Elasticsearch cluster
func GetESHost() string {
	return getEnvValue(esHostEnvVar, "http://localhost:9200")
}

// GetESUsername returns the username to be used, if any, when pushing
// data to the _bulk endpoint
func GetESUsername() string {
	return getEnvValue(esUsernameEnvVar, "")
}

// GetESPassword returns the password to be used, if any, when pushing
// data to the _bulk endpoint
func GetESPassword() string {
	return getEnvValue(esPasswordEnvVar, "")
}

// GetESTimeout returns the duration for which the httpClient should
// wait before closing the connection with the ES _bulk endpoint
func GetESTimeout() time.Duration {
	value := getEnvValue(esTimeoutEnvVar, "60")
	if i, err := strconv.Atoi(value); err == nil {
		return time.Duration(i) * time.Second
	}
	return 60 * time.Second
}

// GetFlushInterval returns the duration for which the httpClient should
// wait before closing the connection with the ES _bulk endpoint
func GetFlushInterval() time.Duration {
	value := getEnvValue(esFlushIntervalEnvVar, "60")
	if i, err := strconv.Atoi(value); err == nil {
		return time.Duration(i) * time.Second
	}
	return 60 * time.Second
}

// getEnvValue returns the value (or default value if key not
// present) for an environment variable
func getEnvValue(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) > 0 {
		return value
	}
	return defaultValue
}
