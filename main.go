package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Middleware pour logs formatés
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, statusCode: 200}
		start := time.Now()
		next.ServeHTTP(rec, r)
		log.Printf(`{"time":"%s", "status":%d, "method":"%s", "url":"%s", "duration_ms":%d}`,
			start.Format(time.RFC3339), rec.statusCode, r.Method, r.URL.String(), time.Since(start).Milliseconds())
	})
}

// Struct pour capturer les codes HTTP
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// Fonction pour formater les noms (camelCase, snake_case -> Capitalisé)
func formatName(name string) string {
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	name = re.ReplaceAllString(name, "$1 $2")
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.Title(name)
	return name
}

// Fonction pour récupérer le hash Git
func getGitHash() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

// Handlers
func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprint(w, "Hello World")
	} else {
		fmt.Fprintf(w, "Hello %s", formatName(name))
	}
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"project": "goApi",
		"gitHash": getGitHash(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	// Lire le port depuis ENV ou flag
	defaultPort := os.Getenv("APP_PORT")
	if defaultPort == "" {
		defaultPort = "8000"
	}
	port := flag.String("port", defaultPort, "Port to listen on")
	flag.Parse()

	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/helloworld", helloHandler)
	mux.HandleFunc("/versionz", versionHandler)

	log.Printf("🚀 Server starting on port %s...", *port)
	err := http.ListenAndServe(":"+*port, loggingMiddleware(mux))
	if err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

