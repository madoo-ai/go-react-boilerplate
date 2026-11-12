// Package webui serves the Vite React build from go:embed. In development
// Vite serves the app on :5173 and proxies /api here; in production this
// handler owns the whole origin.
package webui

import (
	"embed"
	"html"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed dist
var distFS embed.FS

type Server struct {
	files     fs.FS
	appName   string
	appLocale string
}

func New(appName, appLocale string) *Server {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	if appLocale != "en" && appLocale != "id" {
		appLocale = "id"
	}
	if appName == "" {
		appName = "Go React Boilerplate"
	}
	return &Server{files: sub, appName: appName, appLocale: appLocale}
}

func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		requestPath := path.Clean("/" + r.URL.Path)
		if strings.HasPrefix(requestPath, "/api/") || requestPath == "/api" {
			http.NotFound(w, r)
			return
		}

		if strings.HasPrefix(requestPath, "/assets/") {
			s.serveFile(w, r, strings.TrimPrefix(requestPath, "/"), "public, max-age=31536000, immutable")
			return
		}

		// Prerendered marketing and auth routes: /login/index.html.
		if requestPath != "/" {
			candidate := strings.TrimPrefix(requestPath, "/") + "/index.html"
			if s.exists(candidate) {
				s.serveFile(w, r, candidate, "no-store")
				return
			}
			// A public file such as /favicon.svg.
			if filepathIsFile(requestPath) {
				candidate = strings.TrimPrefix(requestPath, "/")
				if s.exists(candidate) {
					s.serveFile(w, r, candidate, "public, max-age=3600")
					return
				}
			}
		}

		// The SPA shell. Hashed assets are immutable; the shell must not be
		// cached because it names those assets.
		s.serveFile(w, r, "index.html", "no-store")
	})
}

func (s *Server) exists(name string) bool {
	info, err := fs.Stat(s.files, name)
	return err == nil && !info.IsDir()
}

func (s *Server) serveFile(w http.ResponseWriter, r *http.Request, name, cacheControl string) {
	data, err := fs.ReadFile(s.files, name)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("web UI is not built; run `make webui`\n"))
		return
	}
	w.Header().Set("Cache-Control", cacheControl)
	if contentType := contentTypeFor(name); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	if strings.HasSuffix(name, ".html") {
		document := string(data)
		document = strings.ReplaceAll(document, `"__APP_LOCALE__"`, `"`+s.appLocale+`"`)
		document = strings.ReplaceAll(document, `"__APP_NAME__"`, `"`+html.EscapeString(s.appName)+`"`)
		document = strings.ReplaceAll(document, "__APP_NAME_TOKEN__", html.EscapeString(s.appName))
		data = []byte(document)
	}
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodHead {
		return
	}
	//nolint:gosec // data is trusted build output from go:embed, not user input.
	_, _ = w.Write(data)
}

func filepathIsFile(requestPath string) bool {
	last := requestPath[strings.LastIndex(requestPath, "/")+1:]
	return strings.Contains(last, ".")
}

func contentTypeFor(name string) string {
	switch {
	case strings.HasSuffix(name, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(name, ".js"):
		return "text/javascript; charset=utf-8"
	case strings.HasSuffix(name, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(name, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(name, ".png"):
		return "image/png"
	case strings.HasSuffix(name, ".ico"):
		return "image/x-icon"
	case strings.HasSuffix(name, ".woff2"):
		return "font/woff2"
	case strings.HasSuffix(name, ".woff"):
		return "font/woff"
	case strings.HasSuffix(name, ".json"):
		return "application/json"
	default:
		return ""
	}
}
