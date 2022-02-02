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
type responseWriterLogger struct {
	http.ResponseWriter
	status  int
	content string
}

func (w *responseWriterLogger) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriterLogger) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.content = string(b)

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

func logResponse(sw *responseWriterLogger, r *http.Request, start time.Time, logContentMode bool) {
	if r.RequestURI != os.Getenv("PROB_URL") {
		var fields = log.Fields{
			"method":         r.Method,
			"url":            r.URL.Path,
			"execution_time": time.Since(start),
			"status":         sw.status,
		}
		if logContentMode {
			fields["content"] = string(sw.content)
			fields["header"] = sw.ResponseWriter.Header()
		}
		log.WithFields(fields).Debug("Query response")
	}
}
