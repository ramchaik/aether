package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/time/rate"
)

func (s *Server) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.StripSlashes)

	// Add rate limiting
	s.router.Use(s.rateLimiter)

	// Add CORS middleware
	s.router.Use(s.corsMiddleware)

	s.router.Use(s.proxyMiddleware)

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Root path accessed")
	})
}

// Rate limiter middleware
func (s *Server) rateLimiter(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Every(time.Second), 10) // 10 requests per second
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CORS middleware
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) proxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		r = r.WithContext(ctx)

		pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")

		var projectID string
		if len(pathParts) > 0 && isValidUUID(pathParts[0]) {
			projectID = pathParts[0]
			http.SetCookie(w, &http.Cookie{
				Name:     "projectID",
				Value:    projectID,
				Path:     "/",
				HttpOnly: true, // Prevents JavaScript access to the cookie
				SameSite: http.SameSiteLaxMode,
			})
			log.Printf("Set cookie with projectID: %s", projectID)
		} else {
			cookie, err := r.Cookie("projectID")
			if err != nil {
				log.Printf("No projectID cookie found and no valid UUID in path: %v", err)
				http.Error(w, "Project ID not found", http.StatusBadRequest)
				return
			}
			projectID = cookie.Value
		}

		resolvesTo := fmt.Sprintf("%s/%s/build", s.basePath, projectID)

		target, err := url.Parse(resolvesTo)
		if err != nil {
			log.Printf("Error parsing target URL: %v", err)
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Director = s.modifyRequest(target, projectID)
		proxy.ErrorHandler = s.errorHandler
		proxy.ServeHTTP(w, r)
	})
}

func isValidUUID(uuid string) bool {
	re := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return re.MatchString(uuid)
}

func (s *Server) modifyRequest(target *url.URL, projectID string) func(*http.Request) {
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		path := strings.TrimPrefix(req.URL.Path, "/"+projectID)
		if path == "" || path == "/" {
			path = "/index.html"
		}
		req.URL.Path = singleJoiningSlash(target.Path, path)

		req.Host = target.Host

		// Remove sensitive headers
		req.Header.Del("Cookie")

		// Add X-Forwarded headers
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Proto", "http")
	}
}

func (s *Server) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Proxy error: %v", err)
	http.Error(w, "Proxy Error", http.StatusBadGateway)
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
