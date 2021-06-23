# Fund Aggregator

This Go application reads a genesis JSON file and aggregate total funds of each account.

Funds can exist as the following formats:
- Available
- Delegated
- Unbonding
- ~~Rewards~~

Thus, this application aggregates all funds by inspecting the genesis file.

**NOTE: To reduce complexity of this application, the genesis file must be exported with `--for-zero-height`
which distributes rewards to proper delegators automatically.
In other words, this application doesn't inspect rewards of each account.**


## Building

```bash
go build ./...
```


## Running

```bash
GENESIS_PATH=<json-file-path> ./fund-aggregator
```

The result is printed to stdout as CSV format:
```csv
panacea1h2k9m0s5qwpnxrwumscn0hs3jmvhxhy2m05yjj,100.0
panacea1fpvuwt4krlmzaq6tyvtwy0w3h8yhp2cv59da88,20.31
panacea1w3ze2ulad0jq7zcps7kdwsadhlh9mc275zt83z,777.1
```