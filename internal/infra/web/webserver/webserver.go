package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Route struct {
	Path   string
	Method string
	Action http.HandlerFunc
}

type WebServer struct {
	Router        chi.Router
	Handlers      []*Route
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	var handlers []*Route

	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      handlers,
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(route *Route) {
	s.Handlers = append(s.Handlers, route)
}

// loop through the handlers and add them to the router
// register middleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	s.Router.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	s.Router.Route("/api", func(r chi.Router) {
		for _, route := range s.Handlers {
			switch route.Method {
			case http.MethodGet:
				r.Get(route.Path, route.Action)
			case http.MethodPost:
				r.Post(route.Path, route.Action)
			default:
				panic("not found method")
			}
		}
	})
	http.ListenAndServe(s.WebServerPort, s.Router)
}
