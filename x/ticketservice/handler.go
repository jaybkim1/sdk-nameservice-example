package ticketservice

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)
  
// NewHandler returns a handler for "ticketservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetTicket:
			return handleMsgSetTicket(ctx, keeper, msg)
		case MsgBuyTicket:
			return handleMsgBuyTicket(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSetTicket
func handleMsgSetTicket(ctx sdk.Context, keeper Keeper, msg MsgSetTicket) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.TicketName)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	keeper.SetTicket(ctx, msg.TicketName, msg.Value) // If so, set the name to the value specified in the msg.
	return sdk.Result{}                        // return
}

// Handle MsgBuyTicket
func handleMsgBuyTicket(ctx sdk.Context, keeper Keeper, msg MsgBuyTicket) sdk.Result {
	if keeper.GetPrice(ctx, msg.TicketName).IsGTE(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.TicketName) {
		_, err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.TicketName), msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}
	keeper.SetOwner(ctx, msg.TicketName, msg.Buyer)
	keeper.SetPrice(ctx, msg.TicketName, msg.Bid)
	return sdk.Result{}
}
