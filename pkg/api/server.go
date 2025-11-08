package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dm0275/mcrun/pkg/minecraft"
)

// Server exposes an HTTP API for provisioning Minecraft servers via mcrun.
type Server struct {
	addr       string
	httpServer *http.Server
	logger     *log.Logger
}

// NewServer returns a configured API server listening on host:port.
func NewServer(host string, port int) *Server {
	addr := fmt.Sprintf("%s:%d", host, port)
	s := &Server{
		addr:   addr,
		logger: log.New(os.Stdout, "[api] ", log.LstdFlags),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/servers", s.handleServers)
	mux.HandleFunc("/servers/", s.handleServerByName)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.logRequests(mux),
	}

	return s
}

// Addr returns the server listening address.
func (s *Server) Addr() string {
	return s.addr
}

// Start begins serving requests. If the provided context is cancelled the server
// attempts a graceful shutdown.
func (s *Server) Start(ctx context.Context) error {
	go func() {
		if ctx == nil {
			return
		}

		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			s.logger.Printf("error shutting down API server: %v", err)
		}
	}()

	s.logger.Printf("HTTP API listening on %s", s.addr)
	err := s.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) || err == nil {
		return nil
	}

	return err
}

func (s *Server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.logger.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleServers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleListServers(w, r)
	case http.MethodPost:
		s.handleCreateServer(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleServerByName(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/servers/")
	path = strings.Trim(path, "/")
	if path == "" {
		writeError(w, http.StatusNotFound, "world name missing")
		return
	}

	parts := strings.Split(path, "/")
	worldName := parts[0]
	if worldName == "" {
		writeError(w, http.StatusNotFound, "world name missing")
		return
	}

	if len(parts) == 1 {
		switch r.Method {
		case http.MethodDelete:
			s.handleDeleteServer(w, r, worldName)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	action := parts[1]
	if len(parts) == 2 && action == "stop" {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		s.handleStopServer(w, r, worldName)
		return
	}

	writeError(w, http.StatusNotFound, "unknown server action")
}

func (s *Server) handleCreateServer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateServerRequest
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	config, err := req.ToMinecraftConfig()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := minecraft.SetupDirectories(config); err != nil {
		s.logger.Printf("failed to setup directories: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to prepare server directories")
		return
	}

	if config.LocalServerConfig {
		if err := minecraft.GenerateServerConfig(config.ServerConfigFile); err != nil {
			s.logger.Printf("failed to create server.properties: %v", err)
			writeError(w, http.StatusInternalServerError, "failed to generate server configuration")
			return
		}
	}

	if err := minecraft.SyncMods(config); err != nil {
		s.logger.Printf("failed to sync mods: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := minecraft.GenerateComposeFile(config); err != nil {
		s.logger.Printf("failed to create compose file: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to generate docker compose file")
		return
	}

	if err := minecraft.StartServer(config); err != nil {
		s.logger.Printf("failed to start server: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to start minecraft server")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"worldName": config.WorldName,
		"status":    "starting",
		"type":      req.TypeOrDefault(),
	})
}

func (s *Server) handleListServers(w http.ResponseWriter, r *http.Request) {
	servers, err := minecraft.ListServers()
	if err != nil {
		s.logger.Printf("failed to list servers: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to list servers")
		return
	}

	writeJSON(w, http.StatusOK, servers)
}

func (s *Server) handleStopServer(w http.ResponseWriter, r *http.Request, worldName string) {
	defer r.Body.Close()

	cfg := minecraft.NewMinecraftConfig()
	cfg.WorldName = worldName
	composeFile, err := minecraft.GetComposeFile(cfg)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeError(w, http.StatusNotFound, fmt.Sprintf("server %s not found", worldName))
			return
		}
		s.logger.Printf("failed to locate compose file: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to locate server definition")
		return
	}

	if err := minecraft.StopServer(composeFile); err != nil {
		s.logger.Printf("failed to stop server: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to stop minecraft server")
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"worldName": worldName,
		"status":    "stopping",
	})
}

func (s *Server) handleDeleteServer(w http.ResponseWriter, r *http.Request, worldName string) {
	defer r.Body.Close()

	cfg := minecraft.NewMinecraftConfig()
	cfg.WorldName = worldName
	composeFile, err := minecraft.GetComposeFile(cfg)
	if err == nil {
		if err := minecraft.StopServer(composeFile); err != nil {
			s.logger.Printf("failed to stop server before delete: %v", err)
			writeError(w, http.StatusInternalServerError, "failed to stop minecraft server")
			return
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		s.logger.Printf("failed to locate compose file: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to locate server definition")
		return
	}

	if err := minecraft.DeleteServerResources(worldName); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeError(w, http.StatusNotFound, fmt.Sprintf("server %s not found", worldName))
			return
		}
		s.logger.Printf("failed to delete server resources: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to delete server resources")
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"worldName": worldName,
		"status":    "deleted",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
