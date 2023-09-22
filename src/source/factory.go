package source

import (
	"fmt"

	"github.com/boliev/protosync/src/domain"
)

// CreateSource creates source according the source type
func CreateSource(sourceType string) (domain.Source, error) {
	if sourceType == "github" {
		return NewGithub(), nil
	}

	return nil, fmt.Errorf("unsupported source type %s", sourceType)
}
