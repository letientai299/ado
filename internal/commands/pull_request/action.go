package pull_request

import (
	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/util"
)

type action string

var allActions = []action{
	actionComplete,
	actionPublish,
	actionDraft,
	actionAbandon,
	actionReactivate,
}

func (a action) applicable(cur *models.GitPullRequest) bool {
	switch a {
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

func (a action) exec(cur, next *models.GitPullRequest) bool {
	if !a.applicable(cur) {
		log.Warnf("Action '%s' is not applicable for current PR status", a)
		return false
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
		return false
	}

	return true
}

const (
	actionComplete   action = "complete"
	actionPublish    action = "publish"
	actionDraft      action = "draft" // change status to draft
	actionAbandon    action = "abandon"
	actionReactivate action = "reactivate"

	// TODO (tai): implement these 2 actions.
	// actionApprove action = "approve"
	// actionReject  action = "reject"
)
