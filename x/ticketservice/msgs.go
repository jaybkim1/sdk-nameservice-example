package ticketservice

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSetTicket struct {
	TicketName string
	Value  string
	Owner  sdk.AccAddress
}

// NewMsgSetTicket sets the new message for the ticket
func NewMsgSetTicket(ticket string, value string, owner sdk.AccAddress) MsgSetTicket {
	return MsgSetTicket{
		TicketName: ticket,
		Value:  value,
		Owner:  owner,
	}
}


// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
// Implements Msg.
func (msg MsgSetTicket) Type() string { return "ticketservice" }
func (msg MsgSetTicket) Name() string { return "set_ticket" }

// ValidateBasic validates basic exception
func (msg MsgSetTicket) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.TicketName) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Name and Value cannot be empty")
	}
	return nil
}

// Implements Msg.
func (msg MsgSetTicket) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("<<<여기 GetSignBytes() >>>")
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgSetTicket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgBuyTicket struct {
	TicketName string
	Bid    sdk.Coins
	Buyer  sdk.AccAddress
}

// NewMsgBuyTicket returns MsgBuyTicket object
func NewMsgBuyTicket(ticket string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyTicket {
	return MsgBuyTicket{
		TicketName: ticket,
		Bid:    bid,
		Buyer:  buyer,
	}
}

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
// Implements Msg.
func (msg MsgBuyTicket) Type() string { return "ticketservice" }
func (msg MsgBuyTicket) Name() string { return "buy_ticket" }

// ValidateBasic validates basic exception
func (msg MsgBuyTicket) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.TicketName) == 0 {
		return sdk.ErrUnknownRequest("Ticket and Value cannot be empty")
	}
	if !msg.Bid.IsPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// Implements Msg.
func (msg MsgBuyTicket) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgBuyTicket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}
