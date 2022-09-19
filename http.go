package main

import (
	"io"
	"net/http"
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
			l.Debug("[%s]%s %s (elapsed %v ms)", r.RemoteAddr, r.Method, r.URL, time.Since(t).Milliseconds())
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
		req, err := http.NewRequest(r.Method, host+r.URL.Path, r.Body)
		assert(err)
		req.Header = r.Header
		cli := http.Client{Timeout: time.Minute}
		res, err := cli.Do(req)
		assert(err)
		defer res.Body.Close()
		_, err = io.Copy(w, res.Body)
		assert(err)
	}
}
