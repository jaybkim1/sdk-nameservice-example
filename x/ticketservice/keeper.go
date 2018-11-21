package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper - handlers sets/gets of custom variables for your module
type Keeper struct {
	coinKeeper bank.Keeper

	namesStoreKey  sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	ownersStoreKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	pricesStoreKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

func NewKeeper(coinKeeper bank.Keeper, namesStoreKey sdk.StoreKey, ownersStoreKey sdk.StoreKey, priceStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper:     coinKeeper,
		namesStoreKey:  namesStoreKey,
		ownersStoreKey: ownersStoreKey,
		pricesStoreKey: priceStoreKey,
		cdc:            cdc,
	}
}

// GetTicket - gets the ticket and its value
func (k Keeper) ResolveTicket(ctx sdk.Context, ticket string) string {
	store := ctx.KVStore(k.namesStoreKey)
	bz := store.Get([]byte(ticket))
	return string(bz)
}

// SetTicket - sets the ticket and its value
func (k Keeper) SetTicket(ctx sdk.Context, ticket string, value string) {
	store := ctx.KVStore(k.namesStoreKey)
	store.Set([]byte(ticket), []byte(value))
}

// GetOwner - gets the current owner of a ticket
func (k Keeper) GetOwner(ctx sdk.Context, ticket string) sdk.AccAddress {
	store := ctx.KVStore(k.ownersStoreKey)
	bz := store.Get([]byte(ticket))
	return bz
}

// HasOwner - returns whether or not the ticket already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, ticket string) bool {
	store := ctx.KVStore(k.ownersStoreKey)
	bz := store.Get([]byte(ticket))
	return bz != nil
}

// SetOwner - sets the current owner of a ticket
func (k Keeper) SetOwner(ctx sdk.Context, ticket string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.ownersStoreKey)
	store.Set([]byte(ticket), owner)
}

// GetPrice - gets the current price of a ticket.  If price doesn't exist yet, set to 1steak.
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	if !k.HasOwner(ctx, name) {
		return sdk.Coins{sdk.NewInt64Coin("mycoin", 1)}
	}
	store := ctx.KVStore(k.pricesStoreKey)
	bz := store.Get([]byte(name))
	var price sdk.Coins
	k.cdc.MustUnmarshalBinary(bz, &price)
	return price
}

// SetPrice - sets the current price of a ticket
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	store := ctx.KVStore(k.pricesStoreKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinary(price))
}
