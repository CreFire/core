package main

import (
	"core/tools/log"
	"io"
	"net/http"
)

func main() {
	runWeb()
}
func runWeb() {
	context := http.NewServeMux()
	context.Handle("/404", http.NotFoundHandler())
	context.Handle("/", GetHttp())

	err := http.ListenAndServe(":81", context)
	if err != nil {
		log.Error("listen", log.Err(err))
		return
	}
}

func GetHttp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			log.InfoF("POST ", r.Method)
		}
		input, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("err", log.Err(err))
			return
		}
		InputStr := string(input)
		log.Info("GET", log.String("input", InputStr))
		_, err = w.Write(input)
		if err != nil {
			log.Error("err", log.Err(err))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
