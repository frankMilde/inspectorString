package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/frankMilde/strtrans"
	"github.com/go-zoo/bone" // faster than httprouter
)

func runServer() error {
	mux := bone.New()
	mux.HandleFunc("/", inputs)
	mux.HandleFunc("/api/", serveAnalysis)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(*listenAddr, httpLogger(mux))

	return err
}

var getTempl = template.Must(template.ParseFiles(getTemplate))

func inputs(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		getTempl.Execute(w, nil)

		return
	}

	s := req.FormValue("string")

	// create NewClient
	// send request to backend
	// handle response from backend
	// display result
	backendQuery := "/api/?string=" + s

	http.Redirect(w, req, backendQuery, http.StatusFound)
}

func serveAnalysis(w http.ResponseWriter, req *http.Request) {
	s := req.URL.Query().Get("string")

	writeHtml(&w, inspectString(s))
}

// writeMethods set the proper MIME type to w.Header and displays the
// content as a string
// To see MIME type:  curl -i localhost:8080/json/url
func writeHtml(w *http.ResponseWriter, s string) {
	(*w).Header().Set("Content-Type", "text/html; charset=utf-8")
	(*w).Header().Set("Encoding", "utf-8")

	fmt.Fprint(*w, "<!DOCTYPE html>\n")
	fmt.Fprint(*w, "<html>\n")
	fmt.Fprint(*w, "\t<head>\n")
	fmt.Fprint(*w, "\t\t<title> IS - Inspector String</title>\n")
	fmt.Fprint(*w, "\t</head>\n")
	fmt.Fprint(*w, "\t<body>\n")
	io.WriteString(*w, s)
	fmt.Fprint(*w, "\t</body>\n")
	fmt.Fprint(*w, "</html>\n")
}
func writeText(w *http.ResponseWriter, s string) {
	(*w).Header().Set("Content-Type", "text/plain; charset=utf-8")
	(*w).Header().Set("Encoding", "utf-8")
	io.WriteString(*w, s)
}
func writeErr(w *http.ResponseWriter, err error) {
	(*w).Header().Set("Content-Type", "text/plain; charset=utf-8")
	(*w).Header().Set("Encoding", "utf-8")

	//TODO: Log this error in separate permanent logfile or DB
	errMsg := strtrans.LinebreaksToTwoLinebreaks(err.Error())

	io.WriteString(*w, errMsg)
}

// httpLogger cleanly logs all HTTP requests by wrapping the handler created
// by httprouter
func httpLogger(router http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		router.ServeHTTP(w, req)
		finishTime := time.Now()
		elapsedTime := finishTime.Sub(startTime)
		log.Println(req.Host, req.Method, req.URL, elapsedTime)
	})
}
