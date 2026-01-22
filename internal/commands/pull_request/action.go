package pull_request

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest"
	"github.com/letientai299/ado/internal/util"
)

type action string

const (
	actionApprove    action = "approve"
	actionReject     action = "reject"
	actionResetVote  action = "resetVote"
	actionBuild      action = "build"
	actionComplete   action = "complete"
	actionPublish    action = "publish"
	actionDraft      action = "draft"
	actionAbandon    action = "abandon"
	actionReactivate action = "reactivate"
)

var allActions = []action{
	actionApprove,
	actionReject,
	actionResetVote,
	actionBuild,
	actionComplete,
	actionPublish,
	actionDraft,
	actionAbandon,
	actionReactivate,
}

func (a action) applicable(cur *models.GitPullRequest) bool {
	switch a {
	case actionApprove, actionReject, actionResetVote:
		return cur.Status != nil && *cur.Status == models.PullRequestStatusActive
	case actionDraft:
		return !cur.IsDraft
	case actionPublish:
		return cur.IsDraft
	case actionComplete:
		return !cur.IsDraft && cur.Status != nil && *cur.Status == models.PullRequestStatusActive
	case actionAbandon:
		return cur.Status != nil && *cur.Status == models.PullRequestStatusActive
	case actionReactivate:
		return cur.Status != nil && *cur.Status == models.PullRequestStatusAbandoned
	}
	return false
}

// exec performs the action. Returns true if the PR model was updated and needs
// to be sent to ADO via the Update API.
func (a action) exec(
	ctx context.Context,
	client *rest.Client,
	org string,
	cur, next *models.GitPullRequest,
) (bool, error) {
	if !a.applicable(cur) {
		log.Warnf("Action '%s' is not applicable for current PR status", a)
		return false, nil
	}

	repo := models.ToRepo(cur.Repository, org)
	gitPRs := client.Git().PRs(repo)

	// Handle vote actions via the reviewers API
	if vote, ok := a.voteValue(); ok {
		log.Infof("%s PR", a.displayName())
		_, err := gitPRs.Vote(ctx, cur.PullRequestId, vote)
		return false, err
	}

	if a == actionBuild {
		return build(ctx, client, cur, repo)
	}

	switch a {
	case actionDraft:
		log.Infof("Marking PR as draft")
		next.IsDraft = true

	case actionPublish:
		log.Infof("Marking PR as active")
		next.IsDraft = false

	case actionComplete:
		log.Infof("Marking PR as completed")
		next.Status = util.Ptr(models.PullRequestStatusCompleted)

	case actionAbandon:
		log.Infof("Marking PR as abandoned")
		next.Status = util.Ptr(models.PullRequestStatusAbandoned)

	case actionReactivate:
		log.Infof("Marking PR as active")
		next.Status = util.Ptr(models.PullRequestStatusActive)

	default:
		log.Warnf("unsupported action: %s, ignoring", a)
		return false, nil
	}

	return true, nil
}

func build(
	ctx context.Context,
	client *rest.Client,
	cur *models.GitPullRequest,
	repo config.Repository,
) (bool, error) {
	log.Infof("Re-queueing PR build pipelines")
	projectID := cur.Repository.Project.Id
	evals, err := client.Policy().Evaluations(ctx, repo, projectID, cur.PullRequestId)
	if err != nil {
		return false, err
	}
	for _, e := range evals {
		if e.Configuration.Type.Id != models.PolicyTypeBuildValidation {
			continue
		}

		if determineBuildStatus(e) == PRBuildStatusExpired {
			log.Infof("Re-queueing build: %s", e.Configuration.Type.DisplayName)
			_, err = client.Policy().
				Requeue(ctx, repo, projectID, e.EvaluationId)
			if err != nil {
				log.Errorf("Failed to re-queue build %s: %v", e.EvaluationId, err)
				return false, err
			}
		}
	}
	return false, nil
}

// voteValue returns the vote value if this is a vote action.
func (a action) voteValue() (models.PrVote, bool) {
	switch a {
	case actionApprove:
		return models.VoteApproved, true
	case actionReject:
		return models.VoteRejected, true
	case actionResetVote:
		return models.VoteNone, true
	}
	return 0, false
}

func (a action) displayName() string {
	switch a {
	case actionApprove:
		return "Approving"
	case actionReject:
		return "Rejecting"
	case actionResetVote:
		return "Resetting vote on"
	case actionBuild:
		return "Queueing builds on"
	}
	return string(a)
}

func (a action) hasVoted(
	ctx context.Context,
	prs rest.GitPRs,
	userID string,
	pr *models.GitPullRequest,
) bool {
	vote, isVote := a.voteValue()
	if !isVote {
		return false
	}

	reviewers, err := prs.Reviewers(ctx, pr.PullRequestId)
	if err != nil {
		return false
	}

	for _, r := range reviewers {
		if r.Id == userID {
			return models.PrVote(r.Vote) == vote
		}
	}
	return false
}
