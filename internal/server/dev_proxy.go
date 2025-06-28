package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

// DevProxy handles proxying requests to the Vite development server
type DevProxy struct {
	proxy  *httputil.ReverseProxy
	logger *zap.Logger
	target string
}

// NewDevProxy creates a new development proxy instance
func NewDevProxy(viteURL string, logger *zap.Logger) (*DevProxy, error) {
	target, err := url.Parse(viteURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Customize the director to handle WebSocket upgrades and other headers
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Origin-Host", req.Header.Get("Host"))
		req.Header.Set("X-Forwarded-Proto", "http")
	}

	devProxy := &DevProxy{
		proxy:  proxy,
		logger: logger,
		target: viteURL,
	}

	// Handle errors
	proxy.ErrorHandler = devProxy.errorHandler

	logger.Info("Development proxy configured", zap.String("target", viteURL))

	return devProxy, nil
}

// ServeHTTP handles the proxy request
func (dp *DevProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dp.logger.Debug("Proxying request to Vite dev server",
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
	)

	dp.proxy.ServeHTTP(w, r)
}

// ShouldProxy determines if a request should be proxied to the dev server
func (dp *DevProxy) ShouldProxy(path string) bool {
	// Skip API and WebSocket routes - they should be handled by their respective handlers
	return !strings.HasPrefix(path, "/api/") && !strings.HasPrefix(path, "/ws")
}

// errorHandler handles proxy errors
func (dp *DevProxy) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	dp.logger.Error("Vite proxy error",
		zap.Error(err),
		zap.String("url", r.URL.String()),
		zap.String("target", dp.target),
	)
	http.Error(w, "Vite dev server not available", http.StatusBadGateway)
}
