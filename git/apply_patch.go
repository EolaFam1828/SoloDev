// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package git

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/harness/gitness/errors"
	"github.com/harness/gitness/git/api"
	"github.com/harness/gitness/git/hook"
	"github.com/harness/gitness/git/parser"
	"github.com/harness/gitness/git/sha"
	"github.com/harness/gitness/git/sharedrepo"

	"github.com/rs/zerolog/log"
)

type ApplyPatchParams struct {
	WriteParams
	Message   string
	Branch    string
	NewBranch string
	Patch     []byte

	Committer     *Identity
	CommitterDate *time.Time
	Author        *Identity
	AuthorDate    *time.Time
}

func (p *ApplyPatchParams) Validate() error {
	if err := p.WriteParams.Validate(); err != nil {
		return err
	}
	if len(p.Patch) == 0 {
		return errors.InvalidArgument("patch cannot be empty")
	}

	return nil
}

type ApplyPatchOutput struct {
	CommitID sha.SHA
}

func (s *Service) ApplyPatch(ctx context.Context, params *ApplyPatchParams) (ApplyPatchOutput, error) {
	if params == nil {
		return ApplyPatchOutput{}, ErrNoParamsProvided
	}
	if err := params.Validate(); err != nil {
		return ApplyPatchOutput{}, err
	}

	committer := params.Actor
	if params.Committer != nil {
		committer = *params.Committer
	}
	committerDate := time.Now().UTC()
	if params.CommitterDate != nil {
		committerDate = *params.CommitterDate
	}

	author := committer
	if params.Author != nil {
		author = *params.Author
	}
	authorDate := committerDate
	if params.AuthorDate != nil {
		authorDate = *params.AuthorDate
	}

	repoPath := getFullPathForRepo(s.reposRoot, params.RepoUID)

	isEmpty, err := s.git.HasBranches(ctx, repoPath)
	if err != nil {
		return ApplyPatchOutput{}, fmt.Errorf("ApplyPatch: failed to determine if repository is empty: %w", err)
	}
	if isEmpty {
		return ApplyPatchOutput{}, errors.PreconditionFailed("can't apply a patch to an empty repository")
	}

	headerParams := &CommitFilesParams{
		WriteParams: params.WriteParams,
		Message:     params.Message,
		Branch:      params.Branch,
		NewBranch:   params.NewBranch,
	}
	commit, err := s.validateAndPrepareCommitFilesHeader(ctx, repoPath, isEmpty, headerParams)
	if err != nil {
		return ApplyPatchOutput{}, err
	}
	if commit == nil {
		return ApplyPatchOutput{}, errors.PreconditionFailed("cannot resolve base commit for patch application")
	}

	refOldSHA := commit.SHA
	branchRef := api.GetReferenceFromBranchName(headerParams.Branch)
	if headerParams.Branch != headerParams.NewBranch {
		refOldSHA = sha.Nil
		branchRef = api.GetReferenceFromBranchName(headerParams.NewBranch)
	}

	refUpdater, err := hook.CreateRefUpdater(s.hookClientFactory, params.EnvVars, repoPath)
	if err != nil {
		return ApplyPatchOutput{}, fmt.Errorf("failed to create reference updater: %w", err)
	}

	patchFile, err := os.CreateTemp(s.sharedRepoRoot, "apply-patch-*.patch")
	if err != nil {
		return ApplyPatchOutput{}, fmt.Errorf("failed to create temporary patch file: %w", err)
	}
	patchFileName := patchFile.Name()
	defer func() {
		if rmErr := os.Remove(patchFileName); rmErr != nil {
			log.Ctx(ctx).Warn().Err(rmErr).Str("path", patchFileName).Msg("failed to remove temporary patch file")
		}
	}()

	if _, err := patchFile.Write(params.Patch); err != nil {
		_ = patchFile.Close()
		return ApplyPatchOutput{}, fmt.Errorf("failed to write temporary patch file: %w", err)
	}
	if err := patchFile.Close(); err != nil {
		return ApplyPatchOutput{}, fmt.Errorf("failed to close temporary patch file: %w", err)
	}

	var commitSHA sha.SHA
	err = sharedrepo.Run(ctx, refUpdater, s.sharedRepoRoot, repoPath, func(r *sharedrepo.SharedRepo) error {
		if err := r.SetIndex(ctx, commit.SHA); err != nil {
			return fmt.Errorf("failed to set base commit index: %w", err)
		}

		if err := r.ApplyToIndex(ctx, patchFileName); err != nil {
			return fmt.Errorf("failed to apply patch: %w", err)
		}

		treeSHA, err := r.WriteTree(ctx)
		if err != nil {
			return fmt.Errorf("failed to write patched tree: %w", err)
		}
		if commit.TreeSHA.Equal(treeSHA) {
			return errors.InvalidArgument("patch produced no effective changes")
		}

		authorSig := &api.Signature{
			Identity: api.Identity{
				Name:  author.Name,
				Email: author.Email,
			},
			When: authorDate,
		}
		committerSig := &api.Signature{
			Identity: api.Identity{
				Name:  committer.Name,
				Email: committer.Email,
			},
			When: committerDate,
		}

		commitSHA, err = r.CommitTree(
			ctx,
			authorSig,
			committerSig,
			treeSHA,
			parser.CleanUpWhitespace(params.Message),
			false,
			commit.SHA,
		)
		if err != nil {
			return fmt.Errorf("failed to create patch commit: %w", err)
		}

		ref := hook.ReferenceUpdate{
			Ref: branchRef,
			Old: refOldSHA,
			New: commitSHA,
		}
		if err := refUpdater.Init(ctx, []hook.ReferenceUpdate{ref}); err != nil {
			return fmt.Errorf("failed to update patch reference old=%s new=%s: %w", refOldSHA, commitSHA, err)
		}

		return nil
	})
	if err != nil {
		return ApplyPatchOutput{}, fmt.Errorf("ApplyPatch: failed to apply patch: %w", err)
	}

	return ApplyPatchOutput{CommitID: commitSHA}, nil
}
