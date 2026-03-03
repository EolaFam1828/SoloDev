// Copyright 2026 EolaFam1828. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/harness/gitness/app/api/controller"
	"github.com/harness/gitness/app/api/usererror"
	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/app/bootstrap"
	"github.com/harness/gitness/app/paths"
	"github.com/harness/gitness/app/services/protection"
	"github.com/harness/gitness/audit"
	"github.com/harness/gitness/git"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"

	"github.com/rs/zerolog/log"
)

type ApplyPatchOptions struct {
	Title     string
	Message   string
	Branch    string
	NewBranch string
	Patch     string
	Author    *git.Identity

	DryRunRules bool
	BypassRules bool
}

func (in *ApplyPatchOptions) Sanitize() error {
	in.Title = strings.TrimSpace(in.Title)
	in.Message = strings.TrimSpace(in.Message)
	in.Branch = strings.TrimSpace(in.Branch)
	in.NewBranch = strings.TrimSpace(in.NewBranch)

	if in.Author != nil {
		in.Author.Name = strings.TrimSpace(in.Author.Name)
		in.Author.Email = strings.TrimSpace(in.Author.Email)
	}

	if len(in.Title) > 1024 {
		return usererror.BadRequest("Commit title is too long.")
	}
	if len(in.Message) > 65536 {
		return usererror.BadRequest("Commit message is too long.")
	}
	if strings.TrimSpace(in.Patch) == "" {
		return usererror.BadRequest("Patch cannot be empty.")
	}

	return nil
}

func (c *Controller) ApplyPatch(
	ctx context.Context,
	session *auth.Session,
	repoRef string,
	in *ApplyPatchOptions,
) (git.ApplyPatchOutput, []types.RuleViolations, error) {
	repo, err := c.getRepoCheckAccess(ctx, session, repoRef, enum.PermissionRepoPush)
	if err != nil {
		return git.ApplyPatchOutput{}, nil, err
	}
	if err := in.Sanitize(); err != nil {
		return git.ApplyPatchOutput{}, nil, err
	}
	if in.Branch == "" {
		in.Branch = repo.DefaultBranch
	}

	rules, isRepoOwner, err := c.fetchBranchRules(ctx, session, repo)
	if err != nil {
		return git.ApplyPatchOutput{}, nil, err
	}

	refAction := protection.RefActionUpdate
	branchName := in.Branch
	if in.NewBranch != "" {
		refAction = protection.RefActionCreate
		branchName = in.NewBranch
	}

	violations, err := rules.RefChangeVerify(ctx, protection.RefChangeVerifyInput{
		ResolveUserGroupID: c.userGroupService.ListUserIDsByGroupIDs,
		Actor:              &session.Principal,
		AllowBypass:        in.BypassRules,
		IsRepoOwner:        isRepoOwner,
		Repo:               repo,
		RefAction:          refAction,
		RefType:            protection.RefTypeBranch,
		RefNames:           []string{branchName},
	})
	if err != nil {
		return git.ApplyPatchOutput{}, nil, fmt.Errorf("failed to verify protection rules: %w", err)
	}

	if in.DryRunRules {
		return git.ApplyPatchOutput{}, violations, nil
	}
	if protection.IsCritical(violations) {
		return git.ApplyPatchOutput{}, violations, nil
	}

	writeParams, err := controller.CreateRPCInternalWriteParams(ctx, c.urlProvider, session, repo)
	if err != nil {
		return git.ApplyPatchOutput{}, nil, fmt.Errorf("failed to create RPC write params: %w", err)
	}

	now := time.Now()
	author := identityFromPrincipal(session.Principal)
	if in.Author != nil {
		author = in.Author
	}

	out, err := c.git.ApplyPatch(ctx, &git.ApplyPatchParams{
		WriteParams:   writeParams,
		Message:       git.CommitMessage(in.Title, in.Message),
		Branch:        in.Branch,
		NewBranch:     in.NewBranch,
		Patch:         []byte(in.Patch),
		Committer:     identityFromPrincipal(bootstrap.NewSystemServiceSession().Principal),
		CommitterDate: &now,
		Author:        author,
		AuthorDate:    &now,
	})
	if err != nil {
		return git.ApplyPatchOutput{}, nil, err
	}

	if protection.IsBypassed(violations) {
		err = c.auditService.Log(ctx,
			session.Principal,
			audit.NewResource(
				audit.ResourceTypeRepository,
				repo.Identifier,
				audit.RepoPath,
				repo.Path,
				audit.BypassAction,
				audit.BypassActionCommitted,
				audit.BypassedResourceType,
				audit.BypassedResourceTypeCommit,
				audit.BypassedResourceName,
				out.CommitID.String(),
				audit.ResourceName,
				fmt.Sprintf(
					audit.BypassSHALabelFormat,
					repo.Identifier,
					out.CommitID.String()[0:6],
				),
			),
			audit.ActionBypassed,
			paths.Parent(repo.Path),
			audit.WithNewObject(audit.CommitObject{
				CommitSHA:      out.CommitID.String(),
				RepoPath:       repo.Path,
				RuleViolations: violations,
			}),
		)
	}
	if err != nil {
		log.Ctx(ctx).Warn().Msgf("failed to insert audit log for apply patch operation: %s", err)
	}

	return out, violations, nil
}
