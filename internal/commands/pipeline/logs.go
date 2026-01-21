package pipeline

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/spf13/cobra"
)

//go:embed logs.md
var logsDoc string

// LogsConfig holds configuration for the pipeline logs command.
type LogsConfig struct {
	filterConfig
	build string
	stage string
	job   string
	tail  int
}

func logsCmd() *cobra.Command {
	opts := &LogsConfig{}

	cmd := &cobra.Command{
		Use:     "logs [keywords...]",
		Aliases: []string{"log"},
		Short:   "View build logs",
		Long:    logsDoc,
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.keywords = args
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}
			return newLogsProcessor(c).process(args)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.build, "build", "b", "", "build number or ID (skip build picker)")
	flags.StringVarP(&opts.stage, "stage", "s", "", "stage name pattern (filter stages)")
	flags.StringVarP(&opts.job, "job", "j", "", "job name pattern (filter jobs)")
	flags.IntVarP(&opts.tail, "tail", "n", 0, "show only last N lines")

	return cmd
}

func newLogsProcessor(c *common[*LogsConfig]) *logsProcessor {
	return &logsProcessor{common: c}
}

type logsProcessor struct {
	*common[*LogsConfig]
}

func (l *logsProcessor) process(args []string) error {
	pipeline, err := l.selectPipeline(args)
	if err != nil {
		return err
	}

	build, timeline, err := l.selectBuild(pipeline.Id)
	if err != nil {
		return err
	}

	stage, err := l.selectStage(timeline)
	if err != nil {
		return err
	}

	job, err := l.selectJob(timeline, stage)
	if err != nil {
		return err
	}

	return l.displayLogs(build.Id, job)
}

func (l *logsProcessor) selectPipeline(args []string) (*models.BuildDefinition, error) {
	// Try if the first arg is a pipeline ID
	if len(args) == 1 {
		if id, err := strconv.ParseInt(args[0], 10, 32); err == nil {
			m, err := l.client.Pipelines().Definitions(l.cfg.Repository).ByID(l.ctx, int32(id))
			if err == nil {
				return m, nil
			}
		}
	}

	// Fallback to list/filter logic
	pipelines, err := l.list()
	if err != nil {
		return nil, err
	}

	pipelines = l.filter(pipelines)

	switch len(pipelines) {
	case 0:
		return nil, errors.New("no pipeline found matching the criteria")
	case 1:
		return &pipelines[0], nil
	default:
		return l.pickPipeline(pipelines)
	}
}

func (l *logsProcessor) pickPipeline(
	pipelines []models.BuildDefinition,
) (*models.BuildDefinition, error) {
	selected := pick(pipelines)
	if selected.IsSome() {
		p := selected.Get()
		return &p, nil
	}
	return nil, errors.New("no pipeline selected")
}

// buildDisplay is a display-friendly wrapper for builds.
type buildDisplay struct {
	models.Build
	Stages   []models.TimelineRecord
	Timeline *models.Timeline
}

func (l *logsProcessor) selectBuild(pipelineID int32) (*models.Build, *models.Timeline, error) {
	// If the build flag is provided, try to use it
	if l.opts.build != "" {
		build, err := l.findBuild(pipelineID, l.opts.build)
		if err != nil {
			return nil, nil, err
		}
		timeline, err := l.client.Builds().ForProject(l.cfg.Repository).Timeline(l.ctx, build.Id)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get build timeline: %w", err)
		}
		return build, timeline, nil
	}

	// List recent builds
	builds, err := l.client.Builds().ForProject(l.cfg.Repository).List(l.ctx, rest.BuildListOptions{
		DefinitionID: pipelineID,
		Top:          20,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list builds: %w", err)
	}

	if len(builds) == 0 {
		return nil, nil, errors.New("no builds found for this pipeline")
	}

	// Fetch timelines for all builds in parallel to show stage status
	displayBuilds := make([]buildDisplay, len(builds))
	type result struct {
		idx      int
		timeline *models.Timeline
	}
	results := make(chan result, len(builds))

	for i, b := range builds {
		go func(idx int, buildID int32) {
			timeline, _ := l.client.Builds().ForProject(l.cfg.Repository).Timeline(l.ctx, buildID)
			results <- result{idx: idx, timeline: timeline}
		}(i, b.Id)
	}

	for range builds {
		r := <-results
		b := builds[r.idx]
		if r.timeline != nil {
			stages := extractStages(r.timeline)
			displayBuilds[r.idx] = buildDisplay{Build: b, Stages: stages, Timeline: r.timeline}
		} else {
			displayBuilds[r.idx] = buildDisplay{Build: b}
		}
	}

	if len(builds) == 1 {
		return &builds[0], displayBuilds[0].Timeline, nil
	}

	return l.pickBuild(displayBuilds)
}

func (l *logsProcessor) findBuild(pipelineID int32, buildRef string) (*models.Build, error) {
	// Try parsing as build ID first
	if id, err := strconv.ParseInt(buildRef, 10, 32); err == nil {
		build, err := l.client.Builds().ForProject(l.cfg.Repository).ByID(l.ctx, int32(id))
		if err == nil {
			return build, nil
		}
	}

	// Otherwise search by build number
	builds, err := l.client.Builds().ForProject(l.cfg.Repository).List(l.ctx, rest.BuildListOptions{
		DefinitionID: pipelineID,
		Top:          50,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list builds: %w", err)
	}

	for i := range builds {
		if builds[i].BuildNumber == buildRef {
			return &builds[i], nil
		}
	}

	return nil, fmt.Errorf("build not found: %s", buildRef)
}

func extractStages(timeline *models.Timeline) []models.TimelineRecord {
	var stages []models.TimelineRecord
	for _, r := range timeline.Records {
		if r.Type == "Stage" {
			stages = append(stages, r)
		}
	}
	sort.Slice(stages, func(i, j int) bool {
		return stages[i].Order < stages[j].Order
	})
	return stages
}

func (l *logsProcessor) pickBuild(builds []buildDisplay) (*models.Build, *models.Timeline, error) {
	selected := ui.Pick(builds, ui.PickConfig[buildDisplay]{
		Title:      "Select a build",
		ItemHeight: 2,
		Render:     renderBuildItem,
		FilterValue: func(b buildDisplay) string {
			return strings.ToLower(b.BuildNumber + " " + b.SourceVersionMessage)
		},
	})

	if selected.IsSome() {
		b := selected.Get()
		return &b.Build, b.Timeline, nil
	}
	return nil, nil, errors.New("no build selected")
}

func renderBuildItem(w io.Writer, b buildDisplay, matches []int) {
	// Line 1: #<build number> • <commit message or branch>
	line1 := formatBuildLine1(b.Build)
	line1 = styles.HighlightMatch(line1, matches)

	// Line 2: <person>, <trigger>  <stage list>
	person := formatPerson(b.RequestedFor)
	trigger := formatTrigger(b.Build)
	stagesViz := formatStagesViz(b.Stages)

	line2 := fmt.Sprintf("%s, %s  %s", person, trigger, stagesViz)

	_, _ = fmt.Fprintf(w, "%s\n%s", line1, line2)
}

func formatBuildLine1(b models.Build) string {
	// #<build number> • <commit message or source info>
	var info string
	if b.SourceVersionMessage != "" {
		info = firstLine(b.SourceVersionMessage)
	} else {
		// Fallback to branch info
		info = formatBranch(b.SourceBranch)
	}
	return fmt.Sprintf("#%s • %s", b.BuildNumber, info)
}

func firstLine(s string) string {
	if idx := strings.Index(s, "\n"); idx > 0 {
		return s[:idx]
	}
	return s
}

func formatBranch(branch string) string {
	branch = strings.TrimPrefix(branch, "refs/heads/")
	// For PR branches like "refs/pull/1335520/merge", extract the PR number
	if strings.HasPrefix(branch, "refs/pull/") {
		parts := strings.Split(branch, "/")
		if len(parts) >= 3 {
			return "PR #" + parts[2]
		}
	}
	// Truncate long branch names
	if len(branch) > 40 {
		return branch[:37] + "..."
	}
	return branch
}

func formatPerson(ref *models.IdentityRef) string {
	if ref == nil {
		return "unknown"
	}
	return ref.DisplayName
}

func formatTrigger(b models.Build) string {
	switch b.Reason {
	case "manual":
		return "by manual"
	case "individualCI", "batchedCI":
		return "by CI"
	case "pullRequest":
		if b.TriggerInfo != nil && b.TriggerInfo.PrNumber != "" {
			return "by PR #" + b.TriggerInfo.PrNumber
		}
		return "by PR"
	case "buildCompletion":
		if b.TriggeredByBuild != nil {
			return "by " + b.TriggeredByBuild.BuildNumber
		}
		return "by trigger"
	case "schedule":
		return "by schedule"
	case "resourceTrigger":
		// For resource triggers, extract trigger info from the build number if available
		// Build numbers like "20260120.6_Buddy20260120.7" contain trigger info after "_"
		if idx := strings.Index(b.BuildNumber, "_"); idx > 0 {
			return "by " + b.BuildNumber[idx+1:]
		}
		return "by resource"
	default:
		if b.Reason != "" {
			return "by " + b.Reason
		}
		return ""
	}
}

func formatStagesViz(stages []models.TimelineRecord) string {
	if len(stages) == 0 {
		return ""
	}

	var parts []string
	for _, s := range stages {
		parts = append(parts, stageIcon(s.Result, s.State))
	}
	return strings.Join(parts, "-")
}

func stageIcon(result, state string) string {
	switch result {
	case "succeeded":
		return styles.Success("☑")
	case "failed":
		return styles.Error("☒")
	case "canceled", "skipped":
		return styles.Warn("☐")
	case "partiallySucceeded":
		return styles.Warn("☑")
	default:
		// Check state for in-progress or pending
		if state == "inProgress" {
			return "▶"
		}
		return "☐"
	}
}

func resultIcon(result string) string {
	switch result {
	case "succeeded":
		return styles.Success("✓")
	case "failed":
		return styles.Error("✗")
	case "canceled":
		return styles.Warn("○")
	case "partiallySucceeded":
		return styles.Warn("◐")
	default:
		return styles.Time("●")
	}
}

// stageDisplay is a display-friendly wrapper for stage timeline records.
type stageDisplay struct {
	models.TimelineRecord
}

func (l *logsProcessor) selectStage(timeline *models.Timeline) (*models.TimelineRecord, error) {
	stages := extractStages(timeline)

	if len(stages) == 0 {
		// No stages, return nil to indicate skipping stage selection
		return nil, nil
	}

	// Apply stage filter if provided
	if l.opts.stage != "" {
		var filtered []models.TimelineRecord
		pattern := strings.ToLower(l.opts.stage)
		for _, s := range stages {
			if strings.Contains(strings.ToLower(s.Name), pattern) {
				filtered = append(filtered, s)
			}
		}
		stages = filtered
	}

	if len(stages) == 0 {
		return nil, errors.New("no stages found matching the filter")
	}

	if len(stages) == 1 {
		return &stages[0], nil
	}

	// Convert to display items
	displayStages := make([]stageDisplay, len(stages))
	for i, s := range stages {
		displayStages[i] = stageDisplay{TimelineRecord: s}
	}

	return l.pickStage(displayStages)
}

func (l *logsProcessor) pickStage(stages []stageDisplay) (*models.TimelineRecord, error) {
	selected := ui.Pick(stages, ui.PickConfig[stageDisplay]{
		Title: "Select a stage",
		Render: func(w io.Writer, s stageDisplay, matches []int) {
			icon := stageIcon(s.Result, s.State)
			name := styles.HighlightMatch(s.Name, matches)
			_, _ = fmt.Fprintf(w, "%s %s", icon, name)
		},
		FilterValue: func(s stageDisplay) string { return strings.ToLower(s.Name) },
	})

	if selected.IsSome() {
		s := selected.Get()
		return &s.TimelineRecord, nil
	}
	return nil, errors.New("no stage selected")
}

// jobDisplay is a display-friendly wrapper for timeline records.
type jobDisplay struct {
	models.TimelineRecord
}

func (l *logsProcessor) selectJob(
	timeline *models.Timeline,
	stage *models.TimelineRecord,
) (*models.TimelineRecord, error) {
	// Build a set of phase IDs that belong to the selected stage
	// Timeline hierarchy: Stage → Phase → Job
	stagePhases := make(map[string]bool)
	if stage != nil {
		for _, r := range timeline.Records {
			if r.Type == "Phase" && r.ParentId == stage.Id {
				stagePhases[r.Id] = true
			}
		}
	}

	// Filter to only jobs that have logs
	var jobs []models.TimelineRecord
	for _, r := range timeline.Records {
		if r.Type != "Job" || r.Log == nil {
			continue
		}
		// If the stage is specified, only include jobs from phases in that stage
		if stage != nil && !stagePhases[r.ParentId] {
			continue
		}
		jobs = append(jobs, r)
	}

	if len(jobs) == 0 {
		return nil, errors.New("no jobs with logs found")
	}

	// Sort by order
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Order < jobs[j].Order
	})

	// Apply job filter if provided
	if l.opts.job != "" {
		var filtered []models.TimelineRecord
		pattern := strings.ToLower(l.opts.job)
		for _, j := range jobs {
			if strings.Contains(strings.ToLower(j.Name), pattern) {
				filtered = append(filtered, j)
			}
		}
		jobs = filtered
	}

	if len(jobs) == 0 {
		return nil, errors.New("no jobs found matching the filter")
	}

	if len(jobs) == 1 {
		return &jobs[0], nil
	}

	// Convert to display items
	displayJobs := make([]jobDisplay, len(jobs))
	for i, j := range jobs {
		displayJobs[i] = jobDisplay{TimelineRecord: j}
	}

	return l.pickJob(displayJobs)
}

func (l *logsProcessor) pickJob(jobs []jobDisplay) (*models.TimelineRecord, error) {
	selected := ui.Pick(jobs, ui.PickConfig[jobDisplay]{
		Title: "Select a job",
		Render: func(w io.Writer, j jobDisplay, matches []int) {
			icon := resultIcon(j.Result)
			name := styles.HighlightMatch(j.Name, matches)
			_, _ = fmt.Fprintf(w, "%s %s", icon, name)
		},
		FilterValue: func(j jobDisplay) string { return strings.ToLower(j.Name) },
	})

	if selected.IsSome() {
		j := selected.Get()
		return &j.TimelineRecord, nil
	}
	return nil, errors.New("no job selected")
}

func (l *logsProcessor) displayLogs(buildID int32, job *models.TimelineRecord) error {
	if job.Log == nil {
		return errors.New("job has no logs")
	}

	// Fetch the entire log (tail is handled after fetching since API needs line count)
	content, err := l.client.Builds().
		ForProject(l.cfg.Repository).
		LogContent(l.ctx, buildID, job.Log.Id, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to get log content: %w", err)
	}

	if l.opts.tail > 0 {
		lines := strings.Split(content, "\n")
		if len(lines) > l.opts.tail {
			lines = lines[len(lines)-l.opts.tail:]
		}
		content = strings.Join(lines, "\n")
	}

	fmt.Print(content)
	return nil
}
