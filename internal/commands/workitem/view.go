package workitem

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/sh"
	"github.com/spf13/cobra"
)

//go:embed view.tpl
var viewTpl string

//go:embed view.md
var viewDoc string

type ViewConfig struct {
	filterConfig
	browse    bool
	relations bool
	output    *util.EnumFlag[string]
}

func viewCmd() *cobra.Command {
	opts := &ViewConfig{
		output: util.NewEnumFlag(outputSimple, outputJSON, outputYAML).
			Default(outputSimple),
	}

	cmd := &cobra.Command{
		Use:     "view <id|text>",
		Aliases: []string{"v"},
		Short:   "View detail of a work item",
		Long:    viewDoc,
		Args:    cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.keywords = args
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}

			return newViewProcessor(c).process(args)
		},
	}
	opts.RegisterFlags(cmd)
	flags := cmd.Flags()
	flags.BoolVarP(&opts.browse, "browse", "b", false, "open work item in browser")
	flags.BoolVarP(
		&opts.relations,
		"relations",
		"r",
		false,
		"include relations (links to other items)",
	)
	flags.VarP(opts.output, "output", "o", "output format")
	opts.output.RegisterCompletion(cmd, "output")
	return cmd
}

func newViewProcessor(c *common[*ViewConfig]) *viewProcessor {
	lp := newListProcessor(copyCommon(c, func(b *common[*ListConfig]) *common[*ListConfig] {
		b.opts = &ListConfig{filterConfig: c.opts.filterConfig}
		return b
	}))
	return &viewProcessor{common: c, lp: lp}
}

type viewProcessor struct {
	*common[*ViewConfig]
	lp listProcessor
}

func (v viewProcessor) process(args []string) error {
	wiID, err := v.findWorkItemID(args)
	if err != nil || wiID == 0 {
		return err
	}

	return v.renderByID(wiID)
}

func (v viewProcessor) findWorkItemID(args []string) (int, error) {
	// 1. Try if the first arg is a work item ID
	if len(args) == 1 {
		if id, err := strconv.Atoi(args[0]); err == nil {
			// Verify the work item exists
			_, err = v.client.WorkItems(v.cfg.Repository).ByID(v.ctx, id, models.WorkItemExpandNone)
			if err == nil {
				return id, nil
			}
			// If error, treat the numeric arg as a keyword
		}
	}

	// 2. Fallback to list/filter logic
	wis, err := v.lp.find()
	if err != nil {
		return 0, err
	}
	if len(wis) == 0 {
		return 0, errors.New("no work item found matching the criteria")
	}

	if len(wis) == 1 {
		return wis[0].ID, nil
	}

	if wi, ok := pickWorkItem(wis); ok {
		return wi.ID, nil
	}

	return 0, nil
}

const wiPickTpl = `{{.ID}} {{.Type | faint}} {{.State}} - {{.Title}}`

func pickWorkItem(wis []models.WorkItem) (models.WorkItem, bool) {
	selected := ui.Pick(wis, ui.PickConfig[models.WorkItem]{
		Render: func(w io.Writer, wi models.WorkItem, matches []int) {
			view := toWorkItemView(wi, "")
			view.Title = styles.HighlightMatch(view.Title, matches)
			util.PanicIf(styles.Render(w, wiPickTpl, view))
		},
		FilterValue: func(wi models.WorkItem) string {
			return strings.ToLower(getStringField(wi, models.FieldTitle))
		},
	})

	if selected.IsNil() {
		return models.WorkItem{}, false
	}

	return selected.Get(), true
}

func (v viewProcessor) renderByID(id int) error {
	expand := models.WorkItemExpandNone
	if v.opts.relations {
		expand = models.WorkItemExpandRelations
	}

	wi, err := v.client.WorkItems(v.cfg.Repository).ByID(v.ctx, id, expand)
	if err != nil {
		return err
	}

	return v.renderOne(*wi)
}

func (v viewProcessor) renderOne(wi models.WorkItem) error {
	detail := toWorkItemDetail(wi, v.baseURL)

	if v.opts.browse {
		fmt.Println(detail.WebURL)
		return sh.Browse(detail.WebURL)
	}

	output := strings.ToLower(v.opts.output.Value())
	switch output {
	case outputYAML:
		return styles.DumpYAML(wi)
	case outputJSON:
		return styles.DumpJSON(wi)
	default:
		return styles.RenderOut(viewTpl, detail)
	}
}

// WorkItemDetail is a detailed view of a work item for template rendering.
type WorkItemDetail struct {
	ID            int
	Rev           int
	Title         string
	State         string
	Reason        string
	Type          string
	AssignedTo    string
	CreatedBy     string
	CreatedDate   string
	ChangedBy     string
	ChangedDate   string
	AreaPath      string
	IterationPath string
	Priority      string
	Tags          []string
	Description   string
	ParentID      int
	CommentCount  int
	WebURL        string
	Relations     []RelationView
}

// RelationView is a simplified view of a work item relation.
type RelationView struct {
	Type string
	Name string
	URL  string
}

func toWorkItemDetail(wi models.WorkItem, baseURL string) WorkItemDetail {
	detail := WorkItemDetail{
		ID:            wi.ID,
		Rev:           wi.Rev,
		Title:         getStringField(wi, models.FieldTitle),
		State:         getStringField(wi, models.FieldState),
		Reason:        getStringField(wi, models.FieldReason),
		Type:          getStringField(wi, models.FieldWorkItemType),
		AssignedTo:    getAssignedTo(wi),
		CreatedBy:     getIdentityName(wi, models.FieldCreatedBy),
		CreatedDate:   formatDate(getStringField(wi, models.FieldCreatedDate)),
		ChangedBy:     getIdentityName(wi, models.FieldChangedBy),
		ChangedDate:   formatDate(getStringField(wi, models.FieldChangedDate)),
		AreaPath:      getStringField(wi, models.FieldAreaPath),
		IterationPath: getStringField(wi, models.FieldIterationPath),
		Priority:      getPriority(wi),
		Tags:          getTags(wi),
		Description:   getDescription(wi),
		ParentID:      getParentID(wi),
		CommentCount:  getIntField(wi, models.FieldCommentCount),
		WebURL:        fmt.Sprintf("%s/%d", baseURL, wi.ID),
	}

	// Process relations
	for _, rel := range wi.Relations {
		rv := RelationView{
			Type: getRelationType(rel.Rel),
			URL:  rel.URL,
		}
		if rel.Attributes != nil {
			if name, ok := rel.Attributes["name"].(string); ok {
				rv.Name = name
			}
		}
		detail.Relations = append(detail.Relations, rv)
	}

	return detail
}

func getIdentityName(wi models.WorkItem, field string) string {
	if wi.Fields == nil {
		return ""
	}
	v, ok := wi.Fields[field]
	if !ok {
		return ""
	}
	if m, ok := v.(map[string]any); ok {
		if name, ok := m["displayName"].(string); ok {
			return name
		}
	}
	return ""
}

func getPriority(wi models.WorkItem) string {
	if wi.Fields == nil {
		return ""
	}
	v, ok := wi.Fields[models.FieldPriority]
	if !ok {
		return ""
	}
	switch p := v.(type) {
	case float64:
		return strconv.Itoa(int(p))
	case int:
		return strconv.Itoa(p)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func getTags(wi models.WorkItem) []string {
	tags := getStringField(wi, models.FieldTags)
	if tags == "" {
		return nil
	}
	// Tags are semicolon-delimited
	parts := strings.Split(tags, ";")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}

func getDescription(wi models.WorkItem) string {
	desc := getStringField(wi, models.FieldDescription)
	// Remove HTML tags for display
	return stripHTML(desc)
}

func getParentID(wi models.WorkItem) int {
	if wi.Fields == nil {
		return 0
	}
	v, ok := wi.Fields[models.FieldParent]
	if !ok {
		return 0
	}
	switch p := v.(type) {
	case float64:
		return int(p)
	case int:
		return p
	default:
		return 0
	}
}

func getIntField(wi models.WorkItem, field string) int {
	if wi.Fields == nil {
		return 0
	}
	v, ok := wi.Fields[field]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	default:
		return 0
	}
}

func formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	// ADO returns dates like "2024-01-15T10:30:00Z"
	// Extract just the date part for display
	if len(dateStr) >= 10 {
		return dateStr[:10]
	}
	return dateStr
}

func getRelationType(rel string) string {
	// Common relation types
	switch rel {
	case "System.LinkTypes.Hierarchy-Forward":
		return "Parent"
	case "System.LinkTypes.Hierarchy-Reverse":
		return "Child"
	case "System.LinkTypes.Related":
		return "Related"
	case "System.LinkTypes.Dependency-Forward":
		return "Predecessor"
	case "System.LinkTypes.Dependency-Reverse":
		return "Successor"
	case "ArtifactLink":
		return "Artifact"
	default:
		// Try to extract a readable name from the relation type
		parts := strings.Split(rel, ".")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
		return rel
	}
}

func stripHTML(html string) string {
	if html == "" {
		return ""
	}
	// Simple HTML stripping - remove tags
	var result strings.Builder
	inTag := false
	for _, r := range html {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			result.WriteRune(r)
		}
	}
	// Clean up extra whitespace
	s := result.String()
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	return strings.TrimSpace(s)
}
