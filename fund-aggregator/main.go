package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/kelseyhightower/envconfig"
	"github.com/medibloc/panacea-core/app"
	"github.com/medibloc/panacea-core/types/util"
	tmtypes "github.com/tendermint/tendermint/types"
	"log"
	"os"
)

type Config struct {
	GenesisPath string `envconfig:"GENESIS_PATH" required:"true"`
}

func main() {
	sdkConfig := sdk.GetConfig()
	sdkConfig.SetCoinType(371)
	sdkConfig.SetFullFundraiserPath("44'/371'/0'/0/0")
	sdkConfig.SetBech32PrefixForAccount(util.Bech32PrefixAccAddr, util.Bech32PrefixAccPub)
	sdkConfig.SetBech32PrefixForValidator(util.Bech32PrefixValAddr, util.Bech32PrefixValPub)
	sdkConfig.SetBech32PrefixForConsensusNode(util.Bech32PrefixConsAddr, util.Bech32PrefixConsPub)
	sdkConfig.Seal()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Panic(err)
	}

	genState, err := readAndValidateGenState(&config)
	if err != nil {
		log.Panic(err)
	}

	availableCoins, err := getAllAvailableCoins(genState)
	if err != nil {
		log.Panic(err)
	}

	stakedCoins, err := getAllStakedCoins(genState)
	if err != nil {
		log.Panic(err)
	}

	totalCoins := aggregate(availableCoins, stakedCoins)
	if err := printAsCSV(totalCoins); err != nil {
		log.Panic(err)
	}
}


func readAndValidateGenState(config *Config) (map[string]json.RawMessage, error) {
	cdc := app.MakeCodec()

	genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read genesis file: %w", err)
	}

	var genState map[string]json.RawMessage
	if err := cdc.UnmarshalJSON(genDoc.AppState, &genState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}

	if err := app.ModuleBasics.ValidateGenesis(genState); err != nil {
		return nil, fmt.Errorf("invalid genesis state: %w", err)
	}

	return genState, nil
}

func getAllAvailableCoins(genState map[string]json.RawMessage) (map[string]sdk.Dec, error) {
	var data genaccounts.GenesisState
	if err := genaccounts.ModuleCdc.UnmarshalJSON(genState[genaccounts.ModuleName], &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal genaccounts genesis state: %w", err)
	}

	ret := make(map[string]sdk.Dec, 0)
	for _, account := range genaccounts.GenesisAccounts(data) {
		ret[account.Address.String()] = sdk.NewDecFromInt(account.Coins.AmountOf("umed"))
	}
	return ret, nil
}

func getAllStakedCoins(genState map[string]json.RawMessage) (map[string]sdk.Dec, error) {
	var data staking.GenesisState
	if err := staking.ModuleCdc.UnmarshalJSON(genState[staking.ModuleName], &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal staking genesis state: %w", err)
	}

	ret := make(map[string]sdk.Dec, 0)

	for _, delegation := range data.Delegations {
		delegator := delegation.DelegatorAddress.String()
		addToCoinMap(ret, delegator, delegation.Shares)
	}

	for _, unbonding := range data.UnbondingDelegations {
		delegator := unbonding.DelegatorAddress.String()
		for _, entry := range unbonding.Entries {
			addToCoinMap(ret, delegator, sdk.NewDecFromInt(entry.Balance))
		}
	}

	return ret, nil
}

func aggregate(coinMaps ...map[string]sdk.Dec) map[string]sdk.Dec {
	ret := make(map[string]sdk.Dec, 0)

	for _, coinMap := range coinMaps {
		for address, amount := range coinMap {
			addToCoinMap(ret, address, amount)
		}
	}

	return ret
}

func addToCoinMap(coinMap map[string]sdk.Dec, address string, amount sdk.Dec) {
	if cur, ok := coinMap[address]; ok {
		coinMap[address] = cur.Add(amount)
	} else {
		coinMap[address] = amount
	}
}

func printAsCSV(coinMap map[string]sdk.Dec) error {
	w := csv.NewWriter(os.Stdout)
	for address, amount := range coinMap {
		if err := w.Write([]string{address, amount.String()}); err != nil {
			return fmt.Errorf("failed to print coinMap as CSV: %w", err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		return fmt.Errorf("failed to print coinMap as CSV: %w", err)
	}
	return nil
}
