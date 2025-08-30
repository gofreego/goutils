package debug

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/gofreego/goutils/logger"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version,omitempty"`
}

// RuntimeInfo represents runtime information
type RuntimeInfo struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"num_goroutine"`
	NumCPU       int    `json:"num_cpu"`
	GOMAXPROCS   int    `json:"gomaxprocs"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
}

// ServiceInfo represents service information
type ServiceInfo struct {
	ServiceName string                 `json:"service_name"`
	Version     string                 `json:"version"`
	BuildTime   string                 `json:"build_time,omitempty"`
	GitCommit   string                 `json:"git_commit,omitempty"`
	Environment string                 `json:"environment,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// MemoryStats represents memory statistics
type MemoryStats struct {
	Alloc        uint64 `json:"alloc"`
	TotalAlloc   uint64 `json:"total_alloc"`
	Sys          uint64 `json:"sys"`
	NumGC        uint32 `json:"num_gc"`
	PauseTotalNs uint64 `json:"pause_total_ns"`
}

// HealthHandler handles health check requests
func HealthHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status:    "ok",
			Timestamp: time.Now(),
			Service:   serviceName,
			Version:   "1.0.0", // You can make this configurable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// ReadinessHandler handles readiness probe requests
func ReadinessHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add any readiness checks here (database connections, external services, etc.)
		response := HealthResponse{
			Status:    "ready",
			Timestamp: time.Now(),
			Service:   serviceName,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// LivenessHandler handles liveness probe requests
func LivenessHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status:    "alive",
			Timestamp: time.Now(),
			Service:   serviceName,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// RuntimeHandler provides runtime information
func RuntimeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := RuntimeInfo{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
			GOMAXPROCS:   runtime.GOMAXPROCS(0),
			GOOS:         runtime.GOOS,
			GOARCH:       runtime.GOARCH,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}

// InfoHandler provides service information
func InfoHandler(serviceName, environment string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := ServiceInfo{
			ServiceName: serviceName,
			Version:     "1.0.0", // Make this configurable
			Environment: environment,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}

// MemoryHandler provides memory statistics
func MemoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		stats := MemoryStats{
			Alloc:        m.Alloc,
			TotalAlloc:   m.TotalAlloc,
			Sys:          m.Sys,
			NumGC:        m.NumGC,
			PauseTotalNs: m.PauseTotalNs,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}

// VarsHandler provides runtime variables in expvar format
func VarsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		vars := map[string]interface{}{
			"cmdline":       os.Args,
			"memstats":      m,
			"num_goroutine": runtime.NumGoroutine(),
			"version":       runtime.Version(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vars)
	}
}

// EnvHandler provides filtered environment variables
func EnvHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Filter out sensitive environment variables
		sensitiveKeys := map[string]bool{
			"PASSWORD": true,
			"SECRET":   true,
			"TOKEN":    true,
			"KEY":      true,
			"PRIVATE":  true,
		}

		env := make(map[string]string)
		for _, e := range os.Environ() {
			if len(e) == 0 {
				continue
			}

			pair := make([]string, 2)
			idx := 0
			for i, c := range e {
				if c == '=' {
					pair[0] = e[:i]
					pair[1] = e[i+1:]
					idx = 1
					break
				}
			}

			if idx == 0 {
				pair[0] = e
				pair[1] = ""
			}

			key := pair[0]
			value := pair[1]

			// Check if the key contains sensitive information
			isSensitive := false
			for sensitiveKey := range sensitiveKeys {
				if contains(key, sensitiveKey) {
					isSensitive = true
					break
				}
			}

			if isSensitive {
				env[key] = "***REDACTED***"
			} else {
				env[key] = value
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(env)
	}
}

// PProfIndexHandler provides a custom pprof index page
func PProfIndexHandler(basePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, PProfIndexTemplate, basePath, basePath, basePath, basePath, basePath, basePath, basePath)
	}
}

// DebugIndexHandler provides a main debug page with links to all debug endpoints
func DebugIndexHandler(serviceName, environment, basePath string, enablePprof bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		pprofSection := ""
		if enablePprof {
			pprofSection = fmt.Sprintf(PProfSectionTemplate, basePath, basePath, basePath, basePath, basePath, basePath, basePath)
		}

		fmt.Fprintf(w, DebugIndexTemplate, serviceName, serviceName, environment, basePath, basePath, basePath, basePath, basePath, basePath, basePath, basePath, pprofSection)
	}
}

// RegisterDebugHandlers registers all debug handlers with the given mux
func RegisterDebugHandlers(ctx context.Context, cfg *Config, mux *http.ServeMux, serviceName, environment string) {
	if !cfg.Enabled {
		return
	}
	// Health check endpoints
	mux.HandleFunc("/health", HealthHandler(serviceName))
	mux.HandleFunc("/health/ready", ReadinessHandler(serviceName))
	mux.HandleFunc("/health/live", LivenessHandler(serviceName))

	// Debug endpoints
	mux.HandleFunc("/debug/info", InfoHandler(serviceName, environment))
	mux.HandleFunc("/debug/runtime", RuntimeHandler())
	mux.HandleFunc("/debug/memory", MemoryHandler())
	mux.HandleFunc("/debug/vars", VarsHandler())
	mux.HandleFunc("/debug/env", EnvHandler())

	// Profiling endpoints (only if enabled)
	if cfg.EnablePprof {
		mux.HandleFunc("/debug/pprof/", PProfIndexHandler(""))
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		mux.Handle("/debug/pprof/block", pprof.Handler("block"))
		mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
		mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))

		logger.Info(ctx, "pprof debug endpoints enabled at /debug/pprof/")
	}

	logger.Info(ctx, "Debug endpoints registered: /health, /health/ready, /health/live, /debug/*")
}

// RegisterDebugHandlersWithGateway registers debug handlers with grpc-gateway ServeMux
func RegisterDebugHandlersWithGateway(ctx context.Context, cfg *Config, mux *gwruntime.ServeMux, serviceName, environment, basePath string) {

	// Health check endpoints
	mux.HandlePath("GET", basePath+"/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		HealthHandler(serviceName)(w, r)
	})
	mux.HandlePath("GET", basePath+"/health/ready", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		ReadinessHandler(serviceName)(w, r)
	})
	mux.HandlePath("GET", basePath+"/health/live", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		LivenessHandler(serviceName)(w, r)
	})

	// Debug endpoints
	mux.HandlePath("GET", basePath+"/debug", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		DebugIndexHandler(serviceName, environment, basePath, cfg.EnablePprof)(w, r)
	})
	mux.HandlePath("GET", basePath+"/debug/info", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		InfoHandler(serviceName, environment)(w, r)
	})
	mux.HandlePath("GET", basePath+"/debug/runtime", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		RuntimeHandler()(w, r)
	})
	mux.HandlePath("GET", basePath+"/debug/memory", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		MemoryHandler()(w, r)
	})
	mux.HandlePath("GET", basePath+"/debug/vars", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		VarsHandler()(w, r)
	})
	mux.HandlePath("GET", basePath+"/debug/env", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		EnvHandler()(w, r)
	})

	// Profiling endpoints (only if enabled)
	if cfg.EnablePprof {
		mux.HandlePath("GET", basePath+"/debug/pprof/", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			PProfIndexHandler(basePath)(w, r)
		})
		mux.HandlePath("GET", basePath+"/debug/pprof/cmdline", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			pprof.Cmdline(w, r)
		})
		mux.HandlePath("GET", basePath+"/debug/pprof/profile", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			pprof.Profile(w, r)
		})
		mux.HandlePath("GET", basePath+"/debug/pprof/symbol", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			pprof.Symbol(w, r)
		})
		mux.HandlePath("GET", basePath+"/debug/pprof/trace", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			pprof.Trace(w, r)
		})
		mux.HandlePath("GET", basePath+"/debug/pprof/goroutine", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			pprof.Handler("goroutine").ServeHTTP(w, r)
		})
		mux.HandlePath("GET", basePath+"/debug/pprof/heap", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			pprof.Handler("heap").ServeHTTP(w, r)
		})
		mux.HandlePath("GET", basePath+"/debug/pprof/threadcreate", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			pprof.Handler("threadcreate").ServeHTTP(w, r)
		})
		mux.HandlePath("GET", basePath+"/debug/pprof/block", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			pprof.Handler("block").ServeHTTP(w, r)
		})
		mux.HandlePath("GET", basePath+"/debug/pprof/mutex", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			pprof.Handler("mutex").ServeHTTP(w, r)
		})
		mux.HandlePath("GET", basePath+"/debug/pprof/allocs", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			pprof.Handler("allocs").ServeHTTP(w, r)
		})

		logger.Info(ctx, "pprof debug endpoints enabled at %s/debug/pprof/", basePath)
	}

	logger.Info(ctx, "Debug endpoints registered: %s/health, %s/health/ready, %s/health/live, %s/debug, %s/debug/*", basePath, basePath, basePath, basePath, basePath)
} // Helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return indexOf(s, substr) >= 0
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i, b := range []byte(s) {
		if b >= 'A' && b <= 'Z' {
			result[i] = b + ('a' - 'A')
		} else {
			result[i] = b
		}
	}
	return string(result)
}

func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
