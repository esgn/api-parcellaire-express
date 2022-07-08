package main

// added by github.com/dnclain

import (
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
// As it is not possible to retrieve natively  the status and the length
// this interface ensure
type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

// WriteHeader wraps the header writer
// Traps the status code
func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// WriteHeader wraps the response writer
// Traps the length written
// Forces status 200-OK if no status is defined yet.
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

// Use : create middleware handlers chain.
func Use(mws ...Middleware) Middleware {
	_nmw := len(mws)

	if _nmw == 0 {
		log.Fatal("At least one middleware should be passed to Use()")
	}

	// reverse to execute middlware in declaration order.
	if _nmw > 1 {
		for i, j := 0, _nmw-1; i < j; i, j = i+1, j-1 {
			mws[i], mws[j] = mws[j], mws[i]
		}
	}

	return func(endPoint http.Handler) http.Handler {
		for _, mw := range mws {
			endPoint = mw(endPoint)
		}
		return endPoint
	}
}

// LogMw : Add logging to each controller
// Should be the FISRT middleware.
func LogMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ensure response contains status.
		_w, okType := w.(interface{}).(statusWriter)
		if !okType {
			_w = statusWriter{ResponseWriter: w}
		}

		next.ServeHTTP(&_w, r)
		if _w.status >= 400 {
			log.Println(formatLog("{} {} {} : {} > ❌ {}", r.RemoteAddr, r.Method, r.URL.String(), r.Form.Encode(), _w.status))
			return
		}
		log.Println(formatLog("{} {} {} : {} > OK {}", r.RemoteAddr, r.Method, r.URL.String(), r.Form.Encode(), _w.status))
	})
}

// --------------------
