package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"

	app "github.com/jaybkim1/sdk-nameservice-example"
	faucetcmd "github.com/jaybkim1/sdk-nameservice-example/x/faucet/client/cli"
	ticketservicecmd "github.com/jaybkim1/sdk-nameservice-example/x/ticketservice/client/cli"
)

const storeAcc = "acc"
const storeTCnames = "tc_tickets"
const storeTCowners = "tc_owners"
const storeTCprices = "tc_prices"

// Create "ticketcli" executable file under $HOME directory
var (
	rootCmd = &cobra.Command{
		Use:   "ticketcli",
		Short: "Ticket Client",
	}
	DefaultCLIHome = os.ExpandEnv("$HOME/.ticketcli")
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	rootCmd.AddCommand(client.ConfigCmd())

	// Add standard rpc commands
	rpc.AddCommands(rootCmd)

	// Define query command
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(),
	)
	tx.AddCommands(queryCmd, cdc)

	// query > account, resolve, whois command
	queryCmd.AddCommand(client.LineBreak)
	queryCmd.AddCommand(client.GetCommands(
		authcmd.GetAccountCmd(storeAcc, cdc, authcmd.GetAccountDecoder(cdc)),
		ticketservicecmd.GetCmdResolveTicket("ticketservice", cdc),
		ticketservicecmd.GetCmdWhois("ticketservice", cdc),
	)...)

	// Define tx command
	txCmd := &cobra.Command{
		Use:     "tx",
		Aliases: []string{"t"},
		Short:   "Transactions subcommands",
	}

	// tx > get-ticket, set-ticket, request-coins
	txCmd.AddCommand(client.PostCommands(
		ticketservicecmd.GetCmdBuyTicket(cdc),
		ticketservicecmd.GetCmdSetTicket(cdc),
		faucetcmd.GetCmdRequestCoins(cdc),
	)...)

	// queryCmd, tmCmd, LineBreak
	rootCmd.AddCommand(
		queryCmd,
		txCmd,
		client.LineBreak,
	)

	rootCmd.AddCommand(
		keys.Commands(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "TC", DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
