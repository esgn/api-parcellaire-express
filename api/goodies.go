package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// cPARAM is the default string formatter parameter.
const cPARAM string = "{}"

// -----------------
// Logger Goodies 🥲
// -----------------

// formatLog simplifies the use of a string formatter.
func formatLog(message string, ps ...interface{}) string {
	pl := len(ps)
	arr := strings.Split(message, cPARAM)
	al := len(arr)
	buf := []string{}
	i := 0
	for key := range arr {
		buf = append(buf, arr[key])
		if key < pl {
			buf = append(buf, fmt.Sprintf("%+v", ps[key]))
			// I add '{}' only if has less parameters
		} else if pl < (al - 1) {
			buf = append(buf, cPARAM)
		}
		i++
	}
	return strings.Join(buf, "")
}

// --------------------
// Http Tmux Goodies 😜
// --------------------

// statusWriter wraps an http response.
type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

// WriteHeader wraps the header writer
func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// WriteHeader wraps the response writer
func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// Middleware : from https://golang.io/fr/tutoriels/les-middlewares-avec-go/
type Middleware func(http.Handler) http.Handler

// Controller :
type Controller func(http.ResponseWriter, *http.Request)

// ThenFunc : syntaxic sugar to add middleware
func (mw Middleware) ThenFunc(controller Controller) http.Handler {
	return mw(http.HandlerFunc(controller))
}

// Use : create middleware handler.
func Use(mws ...Middleware) Middleware {
	return func(endPoint http.Handler) http.Handler {
		for _, mw := range mws {
			endPoint = mw(endPoint)
		}
		return endPoint
	}
}

// LogMw : Add logging to each controller
func LogMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_w := statusWriter{ResponseWriter: w}
		next.ServeHTTP(&_w, r)
		if _w.status >= 400 {
			log.Println(formatLog("{} {} {} : {} > ❌ {}", r.RemoteAddr, r.Method, r.URL.String(), r.Form.Encode(), _w.status))

			return
		}
		log.Println(formatLog("{} {} {} : {} > OK {}", r.RemoteAddr, r.Method, r.URL.String(), r.Form.Encode(), _w.status))
	})
}

// Global connecion pool
var DB *sql.DB

// TrMw is usefull for endpoints with transactional operation like update/insert/delete
func TrMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_tx, _err := DB.Begin()
		if _err != nil {
			fmt.Fprintf(w, "KO : %q", _err)
			return
		}

		next.ServeHTTP(w, r)

		// I expect the defered CloseDatabase to not commit if an error happen
		_tx.Commit()
	})
}

// --------------------
