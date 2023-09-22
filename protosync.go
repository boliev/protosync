package protosync

import (
	"log"

	"github.com/boliev/protosync/src/source"
)

// App is the main struct
type App struct {
}

// Run the app
func (a *App) Run() {
	log.Println("App is running!")
	source, err := source.CreateSource("github")
	if err != nil {
		log.Fatal(err)
		return
	}

	protos, err := source.GetAllProtos()
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, proto := range protos {
		log.Println(proto.URL)
		source.DownloadProto(proto.URL)
	}

}
