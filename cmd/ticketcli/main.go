package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"

	app "github.com/jaybkim1/sdk-nameservice-example"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	ticketservicecmd "github.com/workspace/sdk-nameservice-example/x/ticketservice/client/cli"
)

const storeAcc = "acc"
const storeTCnames = "tc_tickets"
const storeTCowners = "tc_owners"
const storeTCprices = "tc_prices"

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
	rpc.AddCommands(rootCmd)

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
	queryCmd.AddCommand(client.LineBreak)
	queryCmd.AddCommand(client.GetCommands(
		authcmd.GetAccountCmd(storeAcc, cdc, authcmd.GetAccountDecoder(cdc)),
		ticketservicecmd.GetCmdResolveTicket("ticketservice", cdc),
		ticketservicecmd.GetCmdWhois("ticketservice", cdc),
	)...)

	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		ticketservicecmd.GetCmdBuyTicket(cdc),
		ticketservicecmd.GetCmdSetTicket(cdc),
	)...)

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
