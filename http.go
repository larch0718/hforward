package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/xrfang/logging/v2"
)

type (
	RedirectItem struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	}
)

func CORS(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		defer L.Catch(func(l logging.Logger, i any) {
			if i != nil {
				http.Error(w, i.(error).Error(), http.StatusInternalServerError)
			}
			l.Print("[%s]%s %s (elapsed %v ms)", r.RemoteAddr, r.Method, r.URL, time.Since(t).Milliseconds())
		})
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		f(w, r)
	}
}

func NewHttpHandFunc(ri RedirectItem) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		host := strings.TrimSuffix(ri.Host, "/")
		u, err := url.Parse(host + r.URL.String())
		assert(err)
		proxy := httputil.ReverseProxy{
			Director: func(r *http.Request) { r.URL = u },
		}
		proxy.ServeHTTP(w, r)
	}
}
