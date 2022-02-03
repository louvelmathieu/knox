package knox

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

/**
 * Custom status writter to log response
 */
type ResponseWriterLogger struct {
	http.ResponseWriter
	Status  int
	Content string
}

func (w *ResponseWriterLogger) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *ResponseWriterLogger) Write(b []byte) (int, error) {
	if w.Status == 0 {
		w.Status = http.StatusOK
	}
	w.Content = string(b)

	return w.ResponseWriter.Write(b)
}

func logQuery(r *http.Request, logContentMode bool) {
	if r.RequestURI != os.Getenv("PROB_URL") {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close() //  must close
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		var fields = log.Fields{
			"method": r.Method,
			"url":    r.RequestURI,
		}
		if logContentMode {
			fields["content"] = string(bodyBytes)
			fields["header"] = r.Header
		}
		log.WithFields(fields).Debug("Query content")
	}
}

func logResponse(sw *ResponseWriterLogger, r *http.Request, start time.Time, logContentMode bool) {
	if r.RequestURI != os.Getenv("PROB_URL") {
		var fields = log.Fields{
			"method":         r.Method,
			"url":            r.URL.Path,
			"execution_time": time.Since(start),
			"status":         sw.Status,
		}
		if logContentMode {
			fields["content"] = string(sw.Content)
			fields["header"] = sw.ResponseWriter.Header()
		}
		log.WithFields(fields).Debug("Query response")
	}
}
