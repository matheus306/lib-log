package reqclientlog

import "net/http"

type LogConfig struct {
	IgnoredHeaders []string `json:"ignored_headers"`
	IgnoredFields  []string `json:"ignored_fields"`
}

type LoggingTransportDTO struct {
	Transport http.RoundTripper
	LogConfig *LogConfig
}
