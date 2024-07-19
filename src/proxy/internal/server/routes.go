package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.StripSlashes)

	s.router.Use(s.proxyMiddleware)

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
	})
}

func (s *Server) proxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		r = r.WithContext(ctx)

		hostname := r.Host
		subdomain := s.getSubdomain(hostname)

		resolvesTo := fmt.Sprintf("%s/%s/build", s.basePath, subdomain)
		target, err := url.Parse(resolvesTo)
		if err != nil {
			log.Printf("Error parsing target URL: %v", err)
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Director = s.modifyRequest(target)
		proxy.ModifyResponse = s.modifyResponse
		proxy.ErrorHandler = s.errorHandler

		proxy.ServeHTTP(w, r)
	})
}

func (s *Server) getSubdomain(hostname string) string {
	if strings.Count(hostname, ".") >= 1 {
		return strings.Split(hostname, ".")[0]
	}
	return "default" // Handle localhost case
}

func (s *Server) modifyRequest(target *url.URL) func(*http.Request) {
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.Host = target.Host

		if strings.HasSuffix(req.URL.Path, "build/") || req.URL.Path == "/" {
			req.URL.Path += "index.html"
		}

		// Remove sensitive headers
		req.Header.Del("Authorization")
		req.Header.Del("Cookie")
	}
}

func (s *Server) modifyResponse(res *http.Response) error {
	// Add security headers
	res.Header.Set("X-Frame-Options", "DENY")
	res.Header.Set("X-Content-Type-Options", "nosniff")
	res.Header.Set("X-XSS-Protection", "1; mode=block")
	res.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	res.Header.Set("Content-Security-Policy", "default-src 'self'")

	return nil
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
