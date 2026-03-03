// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/harness/gitness/app/auth"

	"github.com/rs/zerolog/log"
)

// SSEHandler returns an http.Handler that serves the MCP protocol over HTTP/SSE.
// This is mounted on the existing chi router for remote clients (VSCode, Antigravity).
func (s *Server) SSEHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/sse", s.handleSSE)
	mux.HandleFunc("/message", s.handleHTTPMessage)
	// Root handles both SSE connection and message posting
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleSSE(w, r)
		case http.MethodPost:
			s.handleHTTPMessage(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return mux
}

// handleSSE establishes an SSE connection for the MCP client.
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	session, err := s.auth.AuthenticateHTTP(r)
	if err != nil {
		http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	log.Info().
		Str("principal", session.Principal.UID).
		Str("remote", r.RemoteAddr).
		Msg("MCP SSE client connected")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Send the endpoint URL for the client to post messages to
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	messageURL := fmt.Sprintf("%s://%s/mcp/message", scheme, r.Host)
	fmt.Fprintf(w, "event: endpoint\ndata: %s\n\n", messageURL)
	flusher.Flush()

	// Keep connection alive
	<-r.Context().Done()
	log.Info().Str("principal", session.Principal.UID).Msg("MCP SSE client disconnected")
}

// handleHTTPMessage handles a single JSON-RPC message over HTTP POST.
func (s *Server) handleHTTPMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := s.auth.AuthenticateHTTP(r)
	if err != nil {
		http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 10*1024*1024)) // 10MB limit
	if err != nil {
		http.Error(w, "read body failed", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := s.HandleMessage(r.Context(), session, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

// ServeSSE starts a standalone HTTP server for MCP SSE transport.
func (s *Server) ServeSSE(ctx context.Context, addr string) error {
	handler := s.SSEHandler()
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		_ = server.Close()
	}()

	log.Info().Str("addr", addr).Msg("MCP SSE server starting")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("MCP SSE server: %w", err)
	}

	wg.Wait()
	return nil
}

// StreamableHTTPHandler returns an http.Handler for the Streamable HTTP transport.
// This follows the MCP Streamable HTTP spec — POST for messages, GET for SSE stream.
func (s *Server) StreamableHTTPHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.handleStreamablePost(w, r)
		case http.MethodGet:
			s.handleSSE(w, r)
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (s *Server) handleStreamablePost(w http.ResponseWriter, r *http.Request) {
	session, err := s.auth.AuthenticateHTTP(r)
	if err != nil {
		// For initialize, allow anonymous
		session = &auth.Session{Principal: auth.AnonymousPrincipal}
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 10*1024*1024))
	if err != nil {
		http.Error(w, "read body failed", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate JSON before processing
	if !json.Valid(body) {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	response, err := s.HandleMessage(r.Context(), session, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
