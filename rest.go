package knox

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

/**
 * Error Format for all Action
 */
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

/**
 * Format and send error
 */
func (e Error) Send(w http.ResponseWriter, code int) {
	log.Warning(e.Message)

	w.Header().Add("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(e)
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

/**
 * Authentification protocol
 */
type SecureMethod struct {
	RequireAuthentification bool
	Protocol                string
	Values                  map[string]string
	ReturnFormat            string
}

/**
 * Custom router to store security of route
 */
type SecureRouter struct {
	mux.Router
	LogQueryMode   bool // Log basic query and responde (HTTP code, URL, time ...)
	LogContentMode bool // Log query and response body and headers witch can be hudge
	security       map[string]SecureMethod
	Protocol       map[string]func(http.Handler, bool) http.Handler
	ServeCallback  func(*ResponseWriterLogger, *http.Request, time.Time)
}

func (sr *SecureRouter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request), secureMethod ...SecureMethod) *mux.Route {
	if len(secureMethod) > 0 {
		sr.security[path] = secureMethod[0]
	} else {
		sr.security[path] = SecureMethod{RequireAuthentification: false, ReturnFormat: "json"}
	}
	return sr.Router.HandleFunc(path, f)
}

func (sr *SecureRouter) AddProtocol(protocol string, f func(http.Handler, bool) http.Handler) {
	sr.Protocol[protocol] = f
}

//Wrapper to log query execution time and print it
func (rl SecureRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if rl.LogQueryMode {
		logQuery(r, rl.LogContentMode)
	}

	// Start custom respons Writer to store response and status
	var sw = &ResponseWriterLogger{ResponseWriter: w}

	var match mux.RouteMatch
	if rl.Match(r, &match) {
		// In case of error, use default mux serve for NotFoundHandler, MethodNotAllowedHandler etc ...
		if match.MatchErr != nil {
			rl.Router.ServeHTTP(sw, r)
			if rl.LogQueryMode {
				logResponse(sw, r, start, rl.LogContentMode)
			}
			if rl.ServeCallback != nil {
				rl.ServeCallback(sw, r, start)
			}
			return
		} else {
			pathTemplate, _ := match.Route.GetPathTemplate()
			if secureMethod, ok := rl.security[pathTemplate]; ok {
				switch secureMethod.ReturnFormat {
				case "json":
					w.Header().Add("Content-type", "application/json; charset=utf-8")
				case "html":
					w.Header().Add("Content-type", "text/html; charset=utf-8")
				default:
					w.Header().Add("Content-type", "application/json; charset=utf-8")
				}

				// Security only use path and ignore method
				if secureMethod.Protocol != "" {
					if protocol, ok := rl.Protocol[secureMethod.Protocol]; ok {
						protocol(&rl.Router, secureMethod.RequireAuthentification).ServeHTTP(sw, r)
					} else {
						http.HandlerFunc(error401).ServeHTTP(sw, r)
					}
					if rl.LogQueryMode {
						logResponse(sw, r, start, rl.LogContentMode)
					}
					if rl.ServeCallback != nil {
						rl.ServeCallback(sw, r, start)
					}
					return
				} else {
					// Not secure route
					rl.Router.ServeHTTP(sw, r)
					if rl.LogQueryMode {
						logResponse(sw, r, start, rl.LogContentMode)
					}
					if rl.ServeCallback != nil {
						rl.ServeCallback(sw, r, start)
					}
					return
				}
			}
		}
	}
	rl.Router.ServeHTTP(sw, r)
	if rl.ServeCallback != nil {
		rl.ServeCallback(sw, r, start)
	}
	if rl.LogQueryMode {
		logResponse(sw, r, start, rl.LogContentMode)
	}
}

func NewRouter() SecureRouter {
	var router = mux.NewRouter()

	// Error handler
	router.NotFoundHandler = http.HandlerFunc(error404) // This handler is call in RestHandler.ServeHttp
	router.MethodNotAllowedHandler = http.HandlerFunc(error405)

	var sr = SecureRouter{
		*router,
		false,
		false,
		map[string]SecureMethod{},
		map[string]func(http.Handler, bool) http.Handler{},
		nil,
	}
	// Defaut middleware
	sr.Protocol["jwt"] = jwtMiddleware

	// Default handler
	sr.HandleFunc("/robots.txt", robot)
	sr.HandleFunc("/favicon.ico", favicon)
	if os.Getenv("PROB_URL") != "" {
		sr.HandleFunc(os.Getenv("PROB_URL"), probe)
	}

	return sr
}
