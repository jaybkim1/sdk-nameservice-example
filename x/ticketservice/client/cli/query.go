package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QueryResult struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// GetCmdResolveTicket queries information about a ticket
func GetCmdResolveTicket(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve [ticket]",
		Short: "resolve ticket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ticket := args[0]
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", queryRoute, ticket), nil)
			if err != nil {
				fmt.Printf("could not resolve ticket - %s \n", string(ticket))
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}

// GetCmdWhois queries information about a ticket
func GetCmdWhois(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whois [ticket]",
		Short: "Query whois info of ticket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ticket := args[0]
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", queryRoute, ticket), nil)
			if err != nil {
				fmt.Printf("could not resolve whois - %s \n", string(ticket))
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}
