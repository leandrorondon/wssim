package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// printHelp
func printHelp() {
	fmt.Println("wssim is a generic webservice simulator.")
	fmt.Println("Usage: wssim [OPTIONS]")
	fmt.Println("\nwssim listens in port 8099 by default.")
	fmt.Println("\nOptions:")
	fmt.Println("    -p,  --port      start wssim in a custom port")
	fmt.Println("    -h,  --help      print this help")
}

// main
func main() {
	// Load parameter: port
	var port int
	flag.Usage = printHelp

	flag.IntVar(&port, "p", 8099, "set wssim to listen in a custom port")
	flag.Parse()

	// Configure common handlers
	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler)

	router := httprouter.New()

	// Static webserver
	fs := http.FileServer(http.Dir("web"))
	router.NotFound = commonHandlers.Then(fs)

	// API
	router.GET("/api/:function", wrapHandler(commonHandlers.ThenFunc(apiHandler)))
	router.GET("/api/:function/:id", wrapHandler(commonHandlers.ThenFunc(apiHandler)))
	router.HEAD("/api/:function", wrapHandler(commonHandlers.ThenFunc(apiHandler)))
	router.POST("/api/:function", wrapHandler(commonHandlers.ThenFunc(apiHandler)))
	router.POST("/api/:function/:id", wrapHandler(commonHandlers.ThenFunc(apiHandler)))
	router.DELETE("/api/:function/:id", wrapHandler(commonHandlers.ThenFunc(apiHandler)))
	router.PUT("/api/:function", wrapHandler(commonHandlers.ThenFunc(apiHandler)))
	router.PUT("/api/:function/:id", wrapHandler(commonHandlers.ThenFunc(apiHandler)))

	sport := fmt.Sprintf(":%d", port)
	log.Println("Listening port", sport)
	log.Fatal(http.ListenAndServe(sport, router))
}

// loggingHandler
func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
		if r.Method == "POST" || r.Method == "PUT" {
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("error reading body, ", err)
			} else {
				log.Println("      ", string(bodyBytes))
			}
		}
	}

	return http.HandlerFunc(fn)
}

// recoverHandler
func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				log.Printf("internal error: %s\n", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// wrapHandler
func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	var st int
	var resp string

	s, err := ioutil.ReadFile("responses/statuscode.txt")
	if err == nil {
		st, err = strconv.Atoi(string(s))
	} else {
		fmt.Println(err)
	}

	if err != nil || st == 0 {
		st = http.StatusOK
	}

	accepts := strings.Split(r.Header.Get("Accept"), ",")

	contentType := "application/json"
	if len(accepts) > 0 {
		if accepts[0] != "*/*" {
			contentType = accepts[0]
		}
	}

	fmt.Println(contentType)

	w.Header().Set("Content-Type", contentType)

	if st == http.StatusOK {
		params := context.Get(r, "params").(httprouter.Params)
		s := params.ByName("function")

		if s != "" {
			data, err := loadResponse(r.Method, s)

			if err != nil {
				st = http.StatusBadRequest
				resp = fmt.Sprintf("{\"result\": \"%s\"}", err)
			} else {
				resp = data
			}
		}
	} else {
		resp = fmt.Sprintf("{\"result\": \"%s\"}", http.StatusText(st))
	}

	w.WriteHeader(st)
	fmt.Fprintf(w, resp)
}

// loadResponse
func loadResponse(m string, f string) (string, error) {
	filePath := fmt.Sprintf("responses/%s/%s.json", m, f)

	// Carrega resposta do arquivo
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("response file not found")
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", errors.New("failed to open response file")
	}

	return string(data), nil
}
