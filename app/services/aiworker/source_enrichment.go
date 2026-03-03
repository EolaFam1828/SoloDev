// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package aiworker

import (
	"bytes"
	"context"
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
