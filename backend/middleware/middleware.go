package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sahandhnj/apiclient/service"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func VersionLogging(reqservice *service.ReqService) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			reqID := requestIDFromContext(r.Context())
			start := time.Now()
			vars := mux.Vars(r)
			modelName := vars["modelname"]
			versionNumber, err := strconv.Atoi(vars["versionNumber"])
			if err != nil {
				log.Fatal(err)
			}

			defer func() {
				reqservice.Add(modelName, versionNumber, start, time.Since(start).Nanoseconds())
				log.Println(r.URL.Path, time.Since(start), reqID, modelName, strconv.Itoa(versionNumber))
			}()
			f(w, r)
		}
	}
}

func Logging() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			reqID := requestIDFromContext(r.Context())
			start := time.Now()

			defer func() {
				log.Println(r.URL.Path, time.Since(start), reqID)
			}()
			f(w, r)
		}
	}
}

func Method(m string) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			f(w, r)
		}
	}
}

func LogReq() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := newContextWithRequestID(r.Context(), r)
			f(w, r.WithContext(ctx))
		}
	}
}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
