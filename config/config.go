package config

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

// env var definition
//var logLevelEnvVar = "ESB_LOG_LEVEL"
var httpPortEnvVar = "ESB_HTTP_PORT"
//var allowedHostsEnvVar = "ESB_ALLOWED_HOSTS"
var esHostEnvVar = "ESB_ES_HOST"
var esUsernameEnvVar = "ESB_ES_USERNAME"
var esPasswordEnvVar = "ESB_ES_PASSWORD"
var esTimeoutEnvVar = "ESB_ES_TIMEOUT_SECONDS"
var esFlushIntervalEnvVar = "ESB_FLUSH_INTERVAL_SECONDS"

func GetLogLevel() log.Level {
	//TODO: implement
	//value := getEnvValue(logLevelEnvVar, "INFO")
	return log.TraceLevel
}

func GetHttpPort() string {
	return getEnvValue(httpPortEnvVar, "8889")
}

func GetAllowedHosts() map[string]bool {
	//TODO: implement
	//value := getEnvValue(allowedHostsEnvVar, "")
	return map[string]bool{}
}

func GetESHost() string {
	return getEnvValue(esHostEnvVar, "http://localhost:9200")
}

func GetESUsername() string {
	return getEnvValue(esUsernameEnvVar, "")
}

func GetESPassword() string {
	return getEnvValue(esPasswordEnvVar, "")
}

func GetESTimeout() time.Duration {
	value := getEnvValue(esTimeoutEnvVar, "2")
	if i, err := strconv.Atoi(value); err == nil {
		return time.Duration(i) * time.Second
	}
	return 2 * time.Second
}

func GetFlushInterval() time.Duration {
	value := getEnvValue(esFlushIntervalEnvVar, "10")
	if i, err := strconv.Atoi(value); err == nil {
		return time.Duration(i) * time.Second
	}
	return 10 * time.Second
}

func getEnvValue(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) > 0 {
		return value
	} else {
		return defaultValue
	}
}