package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
)

func (app *App) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

func (app *App) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

func (app *App) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

func (app *App) GetBankKeeper() bankkeeper.Keeper {
	return app.BankKeeper
}

func (app *App) GetStakingKeeper() *stakingkeeper.Keeper {
	return &app.StakingKeeper
}

func (app *App) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.AccountKeeper
}
