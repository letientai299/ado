package pull_request

import (
	_ "embed"
	"fmt"
	"math"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest/git_prs"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
)

//go:embed analysis.md
var analysisDoc string

//go:embed analysis.tpl
var analysisTpl string

type AnalysisConfig struct {
	from   time.Time
	to     time.Time
	top    int
	output *util.EnumFlag[string]

	fromStr string
	toStr   string
}

func analysisCmd() *cobra.Command {
	opts := &AnalysisConfig{
		output: util.NewEnumFlag(outputSimple, outputJSON, outputYAML).
			Default(outputSimple),
	}

	cmd := &cobra.Command{
		Use:   "analysis",
		Short: "Show PR statistics for a date range [experimental]",
		Long:  analysisDoc,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.fromStr != "" {
				t, err := parseAnalysisDate(opts.fromStr)
				if err != nil {
					return fmt.Errorf("--from: %w", err)
				}
				opts.from = t
			} else {
				opts.from = time.Now().UTC().AddDate(0, 0, -30)
			}

			if opts.toStr != "" {
				t, err := parseAnalysisDate(opts.toStr)
				if err != nil {
					return fmt.Errorf("--to: %w", err)
				}
				opts.to = t
			} else {
				opts.to = time.Now().UTC()
			}

			if err := opts.output.Validate(); err != nil {
				return err
			}

			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}
			return newAnalysisProcessor(c).process()
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.fromStr, "from", "", "start of date range (YYYY-MM-DD or RFC3339, default: 30 days ago)")
	flags.StringVar(&opts.toStr, "to", "", "end of date range (YYYY-MM-DD or RFC3339, default: now)")
	flags.IntVar(&opts.top, "top", 0, "number of top contributors/reviewers to show (0 = all)")
	flags.VarP(opts.output, "output", "o", "output format")
	opts.output.RegisterCompletion(cmd, "output")

	return cmd
}

func parseAnalysisDate(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("cannot parse %q: use RFC3339 or YYYY-MM-DD", s)
}

// ContributorStat holds a single entry in the top-contributors or top-reviewers list.
type ContributorStat struct {
	Name  string `json:"name"  yaml:"name"`
	Email string `json:"email" yaml:"email"`
	Count int    `json:"count" yaml:"count"`
}

// AnalysisResult is the top-level data passed to all output renderers.
type AnalysisResult struct {
	From string `json:"from" yaml:"from"`
	To   string `json:"to"   yaml:"to"`

	Total      int `json:"total"       yaml:"total"`
	Active     int `json:"active"      yaml:"active"`
	Completed  int `json:"completed"   yaml:"completed"`
	Abandoned  int `json:"abandoned"   yaml:"abandoned"`
	DraftCount int `json:"draft_count" yaml:"draft_count"`

	DraftRatio string `json:"draft_ratio" yaml:"draft_ratio"`

	AvgActiveTime    string `json:"avg_active_time"    yaml:"avg_active_time"`
	MedianActiveTime string `json:"median_active_time" yaml:"median_active_time"`

	TopContributors []ContributorStat `json:"top_contributors" yaml:"top_contributors"`
	TopReviewers    []ContributorStat `json:"top_reviewers"    yaml:"top_reviewers"`
}

const analysisPageSize = 100

type analysisProcessor struct {
	*common[*AnalysisConfig]
}

func newAnalysisProcessor(c *common[*AnalysisConfig]) *analysisProcessor {
	return &analysisProcessor{common: c}
}

func (a *analysisProcessor) process() error {
	prs, err := a.fetchAll()
	if err != nil {
		return err
	}

	result := a.compute(prs)
	return a.render(result)
}

func (a *analysisProcessor) fetchAll() ([]models.GitPullRequest, error) {
	criteria := &git_prs.SearchCriteria{
		Status:             util.Ptr(models.PullRequestStatusAll),
		MinTime:            util.Ptr(a.opts.from),
		MaxTime:            util.Ptr(a.opts.to),
		QueryTimeRangeType: util.Ptr(models.PullRequestTimeRangeTypeCreated),
	}

	var all []models.GitPullRequest
	skip := 0

	for {
		page, err := a.client.Git().PRs(a.cfg.Repository).List(
			a.ctx,
			git_prs.ListQuery{
				Top:            util.Ptr(analysisPageSize),
				Skip:           util.Ptr(skip),
				SearchCriteria: criteria,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("fetching PRs (skip=%d): %w", skip, err)
		}

		all = append(all, page...)

		if len(page) < analysisPageSize {
			break
		}
		skip += analysisPageSize
	}

	return all, nil
}

func (a *analysisProcessor) compute(prs []models.GitPullRequest) AnalysisResult {
	result := AnalysisResult{
		From: a.opts.from.Format("2006-01-02"),
		To:   a.opts.to.Format("2006-01-02"),
	}

	result.Total = len(prs)
	if result.Total == 0 {
		result.DraftRatio = "0.0%"
		return result
	}

	var durations []float64
	contributors := make(map[string]*ContributorStat)
	reviewers := make(map[string]*ContributorStat)

	for _, pr := range prs {
		if pr.Status != nil {
			switch *pr.Status {
			case models.PullRequestStatusActive:
				result.Active++
			case models.PullRequestStatusCompleted:
				result.Completed++
			case models.PullRequestStatusAbandoned:
				result.Abandoned++
			}
		}

		if pr.IsDraft {
			result.DraftCount++
		}

		if pr.Status != nil && *pr.Status == models.PullRequestStatusCompleted &&
			pr.CreationDate != nil && pr.ClosedDate != nil {
			d := pr.ClosedDate.Sub(*pr.CreationDate).Hours()
			durations = append(durations, d)
		}

		if pr.CreatedBy != nil && pr.CreatedBy.Id != "" {
			c, ok := contributors[pr.CreatedBy.Id]
			if !ok {
				c = &ContributorStat{
					Name:  pr.CreatedBy.DisplayName,
					Email: pr.CreatedBy.UniqueName,
				}
				contributors[pr.CreatedBy.Id] = c
			}
			c.Count++
		}

		for _, rev := range pr.Reviewers {
			if rev.Vote == 0 || rev.IsContainer {
				continue
			}
			r, ok := reviewers[rev.Id]
			if !ok {
				r = &ContributorStat{
					Name:  rev.DisplayName,
					Email: rev.UniqueName,
				}
				reviewers[rev.Id] = r
			}
			r.Count++
		}
	}

	result.DraftRatio = fmt.Sprintf("%.1f%%", float64(result.DraftCount)/float64(result.Total)*100)

	if len(durations) > 0 {
		result.AvgActiveTime = formatAnalysisDuration(averageFloat(durations))
		result.MedianActiveTime = formatAnalysisDuration(medianFloat(durations))
	}

	result.TopContributors = topNStats(contributors, a.opts.top)
	result.TopReviewers = topNStats(reviewers, a.opts.top)

	return result
}

func (a *analysisProcessor) render(result AnalysisResult) error {
	output := strings.ToLower(a.opts.output.Value())
	switch output {
	case outputJSON:
		return styles.DumpJSON(result)
	case outputYAML:
		return styles.DumpYAML(result)
	default:
		return styles.RenderOut(analysisTpl, result, template.FuncMap{
			"add1": func(i int) int { return i + 1 },
		})
	}
}

func averageFloat(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

func medianFloat(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	sorted := make([]float64, len(xs))
	copy(sorted, xs)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

func formatAnalysisDuration(hours float64) string {
	if hours < 0 {
		hours = 0
	}
	total := time.Duration(math.Round(hours * float64(time.Hour)))

	days := int(total.Hours()) / 24
	remHrs := int(total.Hours()) % 24
	remMins := int(total.Minutes()) % 60

	switch {
	case days > 0 && remHrs > 0:
		return fmt.Sprintf("%dd %dh", days, remHrs)
	case days > 0:
		return fmt.Sprintf("%dd", days)
	case remHrs > 0 && remMins > 0:
		return fmt.Sprintf("%dh %dm", remHrs, remMins)
	case remHrs > 0:
		return fmt.Sprintf("%dh", remHrs)
	default:
		return fmt.Sprintf("%dm", remMins)
	}
}

func topNStats(m map[string]*ContributorStat, n int) []ContributorStat {
	list := make([]ContributorStat, 0, len(m))
	for _, v := range m {
		list = append(list, *v)
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Count != list[j].Count {
			return list[i].Count > list[j].Count
		}
		return list[i].Name < list[j].Name
	})
	if n > 0 && len(list) > n {
		list = list[:n]
	}
	return list
}
