# Metrika Tech Test

## Running

`go run cmd/agent/main.go`

Runtime configuration is possible via environment variables as per the 12-factor app pattern.

If you are developing on a system without a proc filesystem available, you can set overrides to use the test fixtures
instead.

```bash
export PROCFS_PATH=$(pwd)/testdata/fixtures/proc
export ALGOD_PID_OVERRIDE=26231
```

| Variable           | Description                                           | Default               |
|--------------------|-------------------------------------------------------|-----------------------|
| SCAN_INTERVAL      | Sampling frequency                                    | 1s                    |
| OUTPUT_FILE        | Output file                                           | ./metrics.log         |
| ALGOD_ADDRESS      | URL of the algod service                              | http://localhost:8080 |
| ALGOD_API_TOKEN    | API token for algod service. _Mandatory_              | N/A                   |
| PROCFS_PATH        | Path to the proc filesystem                           | /proc                 |
| ALGOD_PID_OVERRIDE | Stop auto-detection of algod PID and hardcode to this | N/A                   |

## Tests

Simply run `make test` to run the entire test suite with the race detector enabled. On first run, this will also unpack
the test fixtures archive into `testdata`.

## Discussion

See [docs/discussion.md](./docs/discussion.md)
