// Command workflow-plugin-turnio is a workflow engine external plugin that
// provides turn.io WhatsApp API integration.
// It runs as a subprocess and communicates with the host workflow engine via
// the go-plugin protocol.
package main

import (
	"github.com/GoCodeAlone/workflow-plugin-turnio/internal"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

func main() {
	sdk.Serve(internal.NewTurnIOPlugin())
}
