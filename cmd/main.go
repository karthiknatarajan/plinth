package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/karthiknatarajan/plinth/cmd/server"
	"github.com/karthiknatarajan/plinth/internal/version"
)

type cmdType = string

const (
	defaultEntrypoint = "version"
	templateService   = "%s-%s"

	cmdVersion = "version"
	cmdServer  = "server"
)

var (
	availableCMDs = []cmdType{
		cmdVersion,
		cmdServer,
	}
)

func main() {
	var (
		flagCMD string
	)
	flag.StringVar(&flagCMD,
		"cmd",
		defaultEntrypoint,
		"entrypoint. Available cmd:"+strings.Join(availableCMDs, ","),
	)
	flag.Parse()

	switch flagCMD {
	case cmdServer:
		version.ServiceName = fmt.Sprintf(templateService, version.ProjectName, cmdServer)
		server.Run()

	case cmdVersion:
		b, err := json.Marshal(version.Version{
			Commit:  version.AppCommit,
			Version: version.AppVersion,
		})
		if err != nil {
			panic(err)
		}
		_, _ = fmt.Println(string(b))
	default:
		log.Fatalf("unknown entrypoint")
	}
}
