package ticketservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetTicket{}, "ticketservice/SetTicket", nil)
	cdc.RegisterConcrete(MsgBuyTicket{}, "ticketservice/BuyTicket", nil)
}
