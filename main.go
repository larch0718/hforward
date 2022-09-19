package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/xrfang/logging/v2"
)

func main() {
	ver := flag.Bool("version", false, "show version")
	conf := flag.String("conf", "", "configuration file")
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	if *conf == "" {
		fmt.Printf("ERROR: missing configuration file (-conf)")
		os.Exit(1)
	}
	loadConfig(*conf)
	if len(cf.Redirects) == 0 {
		fmt.Printf("nothing need to redirect")
		os.Exit(0)
	}
	assert(logging.Init(cf.Logging.Path, cf.Logging.Level, &logging.Options{
		Split: cf.Logging.Split,
		Keep:  cf.Logging.Keep,
	}))
	L = logging.Open("log")
	var svr http.Server
	for idx, it := range cf.Redirects {
		h := http.NewServeMux()
		h.HandleFunc("/", CORS(NewHttpHandFunc(it)))
		svr = http.Server{Addr: fmt.Sprintf(":%d", it.Port), Handler: h}
		if idx == len(cf.Redirects)-1 {
			break
		}
		go func(svr http.Server) {
			L.Catch(nil)
			assert(svr.ListenAndServe())
		}(svr)
	}
	assert(svr.ListenAndServe())
}
