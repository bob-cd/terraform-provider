package resource_provider

import (
	"net/http"
	"time"
)

var (
	Host   = "http://localhost:7777"
	Client = &http.Client{Timeout: 10 * time.Second}
)
