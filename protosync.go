package protosync

import (
	"log"

	"github.com/boliev/protosync/src/config"
	"github.com/boliev/protosync/src/domain"
	"github.com/boliev/protosync/src/source"
)

// App is the main struct
type App struct {
}

// Run the app
func (a *App) Run() {
	log.Println("App is running!")

	config, err := config.Parse(".protosync")
	if err != nil {
		log.Fatal(err)
		return
	}

	sources, err := a.сreateSourcesFromConfig(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, src := range sources {
		err = src.SyncProtos()
		if err != nil {
			log.Fatal(err)
			return
		}
	}

}

func (a *App) сreateSourcesFromConfig(config *config.Config) ([]domain.Source, error) {
	sources := []domain.Source{}
	if len(config.Sources.Github) > 0 {
		for name, cfg := range config.Sources.Github {
			sources = append(sources, source.NewGithub(name, cfg.User, cfg.Repo, cfg.Path, cfg.Ref, cfg.SyncPath))
		}
	}

	return sources, nil
}
