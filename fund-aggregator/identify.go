package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"os"
)

type IdentifiedAccount struct {
	moniker string
	address string
	amount sdk.Dec
}

func identify(coinMap map[string]sdk.Dec, genState map[string]json.RawMessage) ([]IdentifiedAccount, error) {
	var data staking.GenesisState
	if err := staking.ModuleCdc.UnmarshalJSON(genState[staking.ModuleName], &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal staking genesis state: %w", err)
	}

	ret := make([]IdentifiedAccount, 0)

	for address, amount := range coinMap {
		addr, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return nil, fmt.Errorf("invalid acc address: %w", err)
		}
		acc := IdentifiedAccount{
			moniker: findMoniker(addr, data),
			address: address,
			amount: amount,
		}
		ret = append(ret, acc)
	}

	return ret, nil
}

func findMoniker(address sdk.AccAddress, state staking.GenesisState) string {
	for _, validator := range state.Validators {
		if validator.OperatorAddress.Equals(address) {
			return validator.GetMoniker()
		}
	}
	return ""
}

func printAsCSV(accounts []IdentifiedAccount) error {
	w := csv.NewWriter(os.Stdout)
	for _, account := range accounts {
		if err := w.Write([]string{account.moniker, account.address, account.amount.String()}); err != nil {
			return fmt.Errorf("failed to print account as CSV: %w", err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		return fmt.Errorf("failed to print accounts as CSV: %w", err)
	}
	return nil
}
