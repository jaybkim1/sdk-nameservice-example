package app

import (
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/jaybkim1/sdk-nameservice-example/x/ticketservice"
	"github.com/jaybkim1/sdk-nameservice-example/x/faucet"
)

const (
	appName = "TicketApp"
)

type TicketApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain     *sdk.KVStoreKey
	keyAccount  *sdk.KVStoreKey
	keyTCtickets  *sdk.KVStoreKey
	keyTCowners *sdk.KVStoreKey
	keyTCprices *sdk.KVStoreKey

	accountMapper auth.AccountMapper
	bankKeeper    bank.Keeper
	ticketKeeper  ticketservice.Keeper
}

func NewTicketApp(logger log.Logger, db dbm.DB) *TicketApp {
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	var app = &TicketApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:     sdk.NewKVStoreKey("main"),
		keyAccount:  sdk.NewKVStoreKey("acc"),
		keyTCtickets:  sdk.NewKVStoreKey("tc_names"),
		keyTCowners: sdk.NewKVStoreKey("tc_owners"),
		keyTCprices: sdk.NewKVStoreKey("tc_prices"),
	}

	app.accountMapper = auth.NewAccountMapper(
		app.cdc,
		app.keyAccount,
		auth.ProtoBaseAccount,
	)

	app.bankKeeper = bank.NewBaseKeeper(app.accountMapper)

	app.ticketKeeper = ticketservice.NewKeeper(
		app.bankKeeper,
		app.keyTCtickets,
		app.keyTCowners,
		app.keyTCprices,
		app.cdc,
	)

	app.Router().
		AddRoute("ticketservice", ticketservice.NewHandler(app.ticketKeeper)).
		AddRoute("faucet", faucet.NewHandler(app.bankKeeper)).
		AddRoute("bank", bank.NewHandler(app.bankKeeper))

	app.QueryRouter().
		AddRoute("ticketservice", ticketservice.NewQuerier(app.ticketKeeper))

	app.SetInitChainer(app.initChainer)

	app.MountStoresIAVL(
		app.keyMain,
		app.keyAccount,
		app.keyTCtickets,
		app.keyTCowners,
		app.keyTCprices,
	)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

type GenesisState struct {
	Accounts []auth.BaseAccount `json:"accounts"`
}

func (app *TicketApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, &acc)
	}

	return abci.ResponseInitChain{}
}

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	ticketservice.RegisterCodec(cdc)
	faucet.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
