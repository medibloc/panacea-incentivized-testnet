# The Finalization Process of the Incentivized Testnet

As announced, the Incentivized Testnet is going to be finished on **24th June 12:00 KST (03:00 UTC)**.

The Panacea team is responsible for calculating how many rewards each validator has earned,
so that the Mainnet incentives are paid in proportion.


## Process

### Export states to JSON

1. At 12:00 KST, stop one of full nodes that are operated by MediBloc.
	- The 2nd duplex node, possibly
2. Export states to JSON by:
	```bash
 	panacead export --height=<final_height> --for-zero-height | jq > final.json
	```
 	- Note that the `--for-zero-height` must be specified, so that all staking rewards can be distributed to each account automatically.
3. Upload the `final.json` to [Github](https://github.com/medibloc/panacea-opentestnet), so that all participants can see.

### Run the [fund-aggregator](https://github.com/medibloc/panacea-opentestnet/tree/main/fund-aggregator) script

Please see its [README](https://github.com/medibloc/panacea-opentestnet/tree/main/fund-aggregator/README.md).

The scripts will print final funds of each account.

### Calculate the Mainnet incentive portions

Please see https://medibloc.gitbook.io/panacea-core/resources/incentivized-testnet#mainnet-incentive.
