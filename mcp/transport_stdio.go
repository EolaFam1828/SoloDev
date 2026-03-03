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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/harness/gitness/app/auth"

	"github.com/rs/zerolog/log"
)

// ServeStdio runs the MCP server over stdin/stdout (for Claude Desktop and local CLI).
// Authentication is done via SOLODEV_API_TOKEN environment variable.
func (s *Server) ServeStdio(ctx context.Context) error {
	// Authenticate from environment token
	token := os.Getenv("SOLODEV_API_TOKEN")
	var session *auth.Session

	if token != "" {
		var err error
		session, err = s.auth.AuthenticateToken(ctx, token)
		if err != nil {
			log.Warn().Err(err).Msg("MCP stdio: token authentication failed, proceeding as anonymous")
			session = &auth.Session{Principal: auth.AnonymousPrincipal}
		}
	} else {
		log.Warn().Msg("MCP stdio: no SOLODEV_API_TOKEN set, proceeding as anonymous")
		session = &auth.Session{Principal: auth.AnonymousPrincipal}
	}

	log.Info().
		Str("principal", session.Principal.UID).
		Msg("MCP stdio server starting")

	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("MCP stdio server shutting down")
			return nil
		default:
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Info().Msg("MCP stdio: EOF on stdin, shutting down")
				return nil
			}
			return fmt.Errorf("read stdin: %w", err)
		}

		if len(line) == 0 || (len(line) == 1 && line[0] == '\n') {
			continue
		}

		response, err := s.HandleMessage(ctx, session, line)
		if err != nil {
			log.Error().Err(err).Msg("MCP stdio: handle message error")
			errResp := Response{
				JSONRPC: "2.0",
				Error:   &ResponseError{Code: ErrCodeInternal, Message: err.Error()},
			}
			if b, merr := json.Marshal(errResp); merr == nil {
				_, _ = writer.Write(b)
				_, _ = writer.Write([]byte("\n"))
			}
			continue
		}

		if response != nil {
			_, _ = writer.Write(response)
			_, _ = writer.Write([]byte("\n"))
		}
	}
}
