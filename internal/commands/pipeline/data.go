package pipeline

import (
	"strconv"

	"github.com/letientai299/ado/internal/models"
)

// Pipeline is the DTO for pipeline display.
type Pipeline struct {
	Id           int32                        `yaml:"id"            json:"id"`
	Name         string                       `yaml:"name"          json:"name"`
	Path         string                       `yaml:"path"          json:"path"`
	YamlFilename string                       `yaml:"yaml_filename" json:"yaml_filename,omitempty"`
	QueueStatus  models.DefinitionQueueStatus `yaml:"queue_status"  json:"queue_status,omitempty"`
	RepoName     string                       `yaml:"repo_name"     json:"repo_name,omitempty"`
	WebURL       string                       `yaml:"web_url"       json:"web_url,omitempty"`
}

func toPipeline(m models.BuildDefinition, baseURL string) Pipeline {
	p := Pipeline{
		Id:          m.Id,
		Name:        m.Name,
		Path:        m.Path,
		QueueStatus: m.QueueStatus,
	}

	if m.Process != nil {
		p.YamlFilename = m.Process.YamlFilename
	}

	if m.Repository != nil {
		p.RepoName = m.Repository.Name
	}

	p.WebURL = baseURL + strconv.FormatInt(int64(m.Id), 10)

	return p
}
