package workitem

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
)

const (
	outputJSON   = "json"
	outputYAML   = "yaml"
	outputSimple = "simple"
)

//go:embed list_simple.tpl
var listSimpleTpl string

//go:embed list.md
var listDoc string

// ListConfig holds configuration for the workitem list command.
// These values can be set in the config file under "workitem.list".
type ListConfig struct {
	// Default output format to use if not specified.
	DefaultOutput string `yaml:"default_output" json:"default_output"`
	// Custom output templates is a map of output format names to their templates.
	CustomOutputTemplates map[string]string `yaml:"custom_output_templates" json:"custom_output_templates"`
	// Maximum number of work items to return.
	Top int `yaml:"top" json:"top"`

	filterConfig `yaml:"-"`
	output       *util.EnumFlag[string] `yaml:"-"`
	top          int                    `yaml:"-"`
	all          bool                   `yaml:"-"` // show all work items, not just mine
	wiType       string                 `yaml:"-"` // filter by work item type
	state        string                 `yaml:"-"` // filter by state
	assignee     string                 `yaml:"-"` // filter by assignee (substring match)
}

func (l *ListConfig) OnResolved(c *cobra.Command) error {
	// Add custom output formats from config
	for name := range l.CustomOutputTemplates {
		l.output.AddAllowed(name)
	}

	// Update default value if configured
	if l.DefaultOutput != "" {
		flag := c.PersistentFlags().Lookup("output")
		if flag != nil {
			flag.DefValue = l.DefaultOutput
		}
		// Set the value if not explicitly changed by the user
		if !c.Flags().Changed("output") {
			_ = l.output.Set(l.DefaultOutput)
		}
	}

	// Use config top if not set via flag
	if l.top == 0 && l.Top > 0 {
		l.top = l.Top
	}
	if l.top == 0 {
		l.top = 50 // Default to 50
	}

	// Validate after all allowed values have been added
	return l.output.Validate()
}

func listCmd() *cobra.Command {
	opts := defaultListConfig()

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List work items",
		Long:    listDoc,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.keywords = args
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}
			return newListProcessor(c).process()
		},
	}

	opts.RegisterFlags(cmd)
	flags := cmd.PersistentFlags()

	// output format
	flags.VarP(opts.output, "output", "o", "output format")
	opts.output.RegisterCompletion(cmd, "output")

	// additional filters
	flags.IntVarP(&opts.top, "top", "n", 50, "maximum number of work items to return")
	flags.BoolVarP(&opts.all, "all", "a", false, "show all work items (not just mine)")
	flags.StringVarP(
		&opts.wiType,
		"type",
		"t",
		"",
		"filter by work item type (e.g., Bug, Task, \"User Story\")",
	)
	flags.StringVarP(&opts.state, "state", "s", "", "filter by state (e.g., New, Active, Closed)")
	flags.StringVarP(&opts.assignee, "assignee", "A", "", "filter by assignee alias or email (substring match); implies --all")

	return cmd
}

func newListProcessor(c *common[*ListConfig]) listProcessor {
	return listProcessor{c}
}

func defaultListConfig() *ListConfig {
	opts := &ListConfig{
		DefaultOutput:         outputSimple,
		CustomOutputTemplates: make(map[string]string),
		output: util.NewEnumFlag(outputSimple, outputJSON, outputYAML).
			Default(outputSimple),
	}

	config.Register(config.CommandConfig{
		Path:   "workitem.list",
		Target: opts,
	})

	return opts
}

type listProcessor struct {
	*common[*ListConfig]
}

func (l listProcessor) process() error {
	wis, err := l.find()
	if err != nil {
		return err
	}

	return l.render(wis)
}

func (l listProcessor) find() ([]models.WorkItem, error) {
	// Execute WIQL query to get work item IDs
	result, err := l.query()
	if err != nil {
		return nil, err
	}

	if len(result.WorkItems) == 0 {
		return nil, nil
	}

	// Extract IDs from WIQL result
	ids := make([]int, len(result.WorkItems))
	for i, ref := range result.WorkItems {
		ids[i] = ref.ID
	}

	// Fetch full work item details
	wis, err := l.client.WorkItems(l.cfg.Repository).List(
		l.ctx,
		ids,
		models.WorkItemExpandNone,
	)
	if err != nil {
		return nil, err
	}

	return l.filter(wis)
}

func (l listProcessor) query() (*models.WIQLResult, error) {
	wiql := l.buildWIQL()
	log.Debug("executing WIQL", "query", wiql)

	return l.client.WIQL(l.cfg.Repository).Query(l.ctx, wiql, l.opts.top)
}

// wiqlEscape escapes a string value for safe embedding inside WIQL single-quoted literals.
func wiqlEscape(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func (l listProcessor) buildWIQL() string {
	var sb strings.Builder
	sb.WriteString("SELECT [System.Id], [System.Title], [System.State], ")
	sb.WriteString("[System.WorkItemType], [System.AssignedTo], [System.ChangedDate] ")
	sb.WriteString("FROM WorkItems WHERE ")

	conditions := []string{}

	// Filter by assignee
	if l.opts.assignee != "" {
		// Explicit assignee: substring match on display name / email
		conditions = append(conditions, fmt.Sprintf("[System.AssignedTo] = '%s'", wiqlEscape(l.opts.assignee)))
	} else if !l.opts.all {
		// Default (including --mine): show only work items assigned to me
		conditions = append(conditions, "[System.AssignedTo] = @Me")
	}

	// Filter by work item type
	if l.opts.wiType != "" {
		conditions = append(conditions, fmt.Sprintf("[System.WorkItemType] = '%s'", wiqlEscape(l.opts.wiType)))
	}

	// Filter by state
	if l.opts.state != "" {
		conditions = append(conditions, fmt.Sprintf("[System.State] = '%s'", wiqlEscape(l.opts.state)))
	}

	// Default: exclude closed/done items
	if l.opts.state == "" {
		conditions = append(conditions, "[System.State] <> 'Closed'")
		conditions = append(conditions, "[System.State] <> 'Done'")
		conditions = append(conditions, "[System.State] <> 'Removed'")
	}

	if len(conditions) == 0 {
		// Always need at least one condition
		conditions = append(conditions, "[System.Id] > 0")
	}

	sb.WriteString(strings.Join(conditions, " AND "))
	sb.WriteString(" ORDER BY [System.ChangedDate] DESC")

	return sb.String()
}

func (l listProcessor) filter(all []models.WorkItem) ([]models.WorkItem, error) {
	if len(l.opts.keywords) == 0 {
		return all, nil
	}

	// Filter by keywords in title
	filtered := make([]models.WorkItem, 0, len(all))
	for _, wi := range all {
		title := strings.ToLower(getStringField(wi, models.FieldTitle))
		matches := true
		for _, kw := range l.opts.keywords {
			if !strings.Contains(title, strings.ToLower(kw)) {
				matches = false
				break
			}
		}
		if matches {
			filtered = append(filtered, wi)
		}
	}

	return filtered, nil
}

func (l listProcessor) render(all []models.WorkItem) error {
	log.Debug("found work items", "count", len(all))

	output := strings.ToLower(l.opts.output.Value())
	switch output {
	case outputYAML:
		return styles.DumpYAML(all)
	case outputJSON:
		return styles.DumpJSON(all)
	case outputSimple:
		return l.renderTemplate(listSimpleTpl, all)
	default:
		if tpl, ok := l.opts.CustomOutputTemplates[output]; ok {
			return l.renderTemplate(tpl, all)
		}
	}

	return util.StrErr("unknown output format: " + l.opts.output.Value())
}

func (l listProcessor) renderTemplate(tpl string, all []models.WorkItem) error {
	items := make([]WorkItemView, len(all))
	for i, wi := range all {
		items[i] = toWorkItemView(wi, l.baseURL)
	}
	return styles.RenderOut(tpl, items)
}

// WorkItemView is a simplified view of a work item for template rendering.
type WorkItemView struct {
	ID          int
	Title       string
	State       string
	Type        string
	AssignedTo  string
	ChangedDate string
	WebURL      string
}

func toWorkItemView(wi models.WorkItem, baseURL string) WorkItemView {
	return WorkItemView{
		ID:          wi.ID,
		Title:       getStringField(wi, models.FieldTitle),
		State:       getStringField(wi, models.FieldState),
		Type:        getStringField(wi, models.FieldWorkItemType),
		AssignedTo:  getAssignedTo(wi),
		ChangedDate: getStringField(wi, models.FieldChangedDate),
		WebURL:      fmt.Sprintf("%s/%d", baseURL, wi.ID),
	}
}

func getStringField(wi models.WorkItem, field string) string {
	if wi.Fields == nil {
		return ""
	}
	v, ok := wi.Fields[field]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func getAssignedTo(wi models.WorkItem) string {
	if wi.Fields == nil {
		return ""
	}
	v, ok := wi.Fields[models.FieldAssignedTo]
	if !ok {
		return ""
	}
	// AssignedTo is an IdentityRef object with displayName field
	if m, ok := v.(map[string]any); ok {
		if name, ok := m["displayName"].(string); ok {
			return name
		}
	}
	return ""
}
