package server

import (
	"net/http"
	"strings"
)

func GetIPAndUserAgent(r *http.Request) (userAgent string, ip string) {
	ip = r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	userAgent = r.UserAgent()
	return userAgent, ip
}
