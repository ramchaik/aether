package server

import (
	"encoding/json"
	"logify/internal/database"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	api := "/api/v1"
	r.Route(api, func(r chi.Router) {
		r.Post("/logs", s.bulkInsertLogs)
		r.Get("/{project-id}/logs", s.getProjectLogs)
	})

	return r
}

func (s *Server) bulkInsertLogs(w http.ResponseWriter, r *http.Request) {
	var logs []database.Log
	if err := json.NewDecoder(r.Body).Decode(&logs); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.db.BulkInsertLogs(r.Context(), logs); err != nil {
		http.Error(w, "Failed to insert logs", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type LogResponse struct {
	ProjectID string `json:"projectId"`
	Timestamp int64  `json:"timestamp"`
	Log       string `json:"log"`
}

func (s *Server) getProjectLogs(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "project-id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	limit := 10 // Default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := s.db.GetProjectLogs(r.Context(), projectID, limit)
	if err != nil {
		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
		return
	}

	var logResponses []LogResponse
	for _, log := range logs {
		logResponses = append(logResponses, LogResponse{
			ProjectID: log.ProjectID,
			Timestamp: log.Timestamp,
			Log:       log.Log,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logResponses)
}
