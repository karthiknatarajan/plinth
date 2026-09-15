package logger

import (
	"os"
	"strings"
	"time"

	"github.com/karthiknatarajan/plinth/internal/types"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	zerologlog "github.com/rs/zerolog/log"
)

func SetupLogger(config *types.Config) {
	switch strings.ToLower(config.Log.Level) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	// if the terminal is a tty we should output the
	// logs in pretty format
	if isatty.IsTerminal(os.Stdout.Fd()) {
		zerologlog.Logger = zerologlog.Output(
			zerolog.ConsoleWriter{
				Out:        os.Stderr,
				NoColor:    false,
				TimeFormat: "15:04:05.999",
			},
		)
	}
}
