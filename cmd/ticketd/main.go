package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/server"

	app "github.com/jaybkim1/sdk-nameservice-example"
)

// Create "ticketd" executable file under $HOME directory
var (
	DefaultNodeHome = os.ExpandEnv("$HOME/.ticketd")
)

// Initialize this app
var appInit = server.AppInit{
	AppGenState: server.SimpleAppGenState,
	AppGenTx:    server.SimpleAppGenTx,
}

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	// Define root command
	rootCmd := &cobra.Command{
		Use:               "ticketd",
		Short:             "Ticket App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	// tendermint show-node-id, show_node-validator, show-address
	server.AddCommands(ctx, cdc, rootCmd, appInit,
		server.ConstructAppCreator(newApp, "ticket"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "ticket"))

	// prepare and add flags (config and data folders)
	executor := cli.PrepareBaseCmd(rootCmd, "TC", DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewTicketApp(logger, db)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	return nil, nil, nil
}
