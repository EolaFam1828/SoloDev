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

package aiworker

import (
	"bytes"
	"context"
	"encoding/json"
	stdliberrors "errors"
	"fmt"
	"io"
	"strings"

	"github.com/harness/gitness/git"
	"github.com/harness/gitness/git/parser"
	"github.com/harness/gitness/types"
)

const maxSourceCodeBytes int64 = 64 * 1024

func (h *remWorkerHandler) enrichSourceCode(ctx context.Context, rem *types.Remediation) error {
	if rem == nil || rem.SourceCode != "" || rem.FilePath == "" || rem.RepoID <= 0 {
		return nil
	}

	gitRef := strings.TrimSpace(rem.CommitSHA)
	if gitRef == "" {
		gitRef = strings.TrimSpace(rem.Branch)
	}
	if gitRef == "" || h.repoStore == nil || h.git == nil {
		return nil
	}

	repo, err := h.repoStore.Find(ctx, rem.RepoID)
	if err != nil {
		return fmt.Errorf("find remediation repo: %w", err)
	}

	node, err := h.git.GetTreeNode(ctx, &git.GetTreeNodeParams{
		ReadParams: git.ReadParams{RepoUID: repo.GitUID},
		GitREF:     gitRef,
		Path:       rem.FilePath,
	})
	if err != nil {
		return fmt.Errorf("resolve remediation file: %w", err)
	}
	if node.Node.Type != git.TreeNodeTypeBlob {
		return nil
	}

	blob, err := h.git.GetBlob(ctx, &git.GetBlobParams{
		ReadParams: git.ReadParams{RepoUID: repo.GitUID},
		SHA:        node.Node.SHA,
		SizeLimit:  maxSourceCodeBytes + 1,
	})
	if err != nil {
		return fmt.Errorf("load remediation blob: %w", err)
	}
	defer blob.Content.Close()

	if blob.Size > maxSourceCodeBytes {
		return nil
	}

	content, err := io.ReadAll(blob.Content)
	if err != nil {
		return fmt.Errorf("read remediation blob: %w", err)
	}
	if _, ok := parser.IsLFSPointer(ctx, content, blob.Size); ok {
		return nil
	}

	scanner, _, err := parser.ReadTextFile(bytes.NewReader(content), nil)
	if err != nil {
		if stdliberrors.Is(err, parser.ErrBinaryFile) {
			return nil
		}
		return fmt.Errorf("inspect remediation source text: %w", err)
	}

	var text strings.Builder
	for scanner.Scan() {
		text.Write(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		if stdliberrors.Is(err, parser.ErrBinaryFile) {
			return nil
		}
		return fmt.Errorf("scan remediation source text: %w", err)
	}

	rem.SourceCode = text.String()
	if rem.SourceCode == "" {
		return nil
	}

	return h.remStore.Update(ctx, rem)
}

// setContextMetadata stores the context bundle provenance in remediation metadata.
func setContextMetadata(existing json.RawMessage, contextJSON json.RawMessage) json.RawMessage {
	var meta map[string]json.RawMessage
	if len(existing) > 0 {
		if err := json.Unmarshal(existing, &meta); err != nil {
			meta = make(map[string]json.RawMessage)
		}
	} else {
		meta = make(map[string]json.RawMessage)
	}
	meta["context_bundle"] = contextJSON
	out, err := json.Marshal(meta)
	if err != nil {
		return existing
	}
	return out
}

func setDeliveryMetadata(
	rem *types.Remediation,
	mode types.RemediationDeliveryMode,
	state types.RemediationDeliveryState,
	lastError string,
	prNumber int64,
	attemptedAt int64,
) error {
	delivery, err := types.GetRemediationDeliveryMetadata(rem.Metadata, mode)
	if err != nil {
		return err
	}
	if mode != "" {
		delivery.Mode = mode
	}
	if state != "" {
		delivery.State = state
	}
	delivery.LastError = lastError
	delivery.PRNumber = prNumber
	delivery.AttemptedAt = attemptedAt

	metadata, err := types.SetRemediationDeliveryMetadata(rem.Metadata, delivery)
	if err != nil {
		return err
	}
	rem.Metadata = metadata

	return nil
}
