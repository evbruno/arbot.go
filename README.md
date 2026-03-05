# arbot.go

A Go-based foundation for an automated Ethereum arbitrage bot. The current codebase focuses on ABI bindings and on-chain reads for Uniswap V2-style factories and pairs. More strategy, execution, and risk controls will be added as the project evolves.

## Status

Early stage. The repository is intended to grow into a full, automated arbitrage system, but it is not production-ready.

## Structure

- `main.go` - Entry point.
- `abi/` - Contract ABIs used by generated bindings.
- `factory/` - Factory contract bindings and helpers.
- `pair/` - Pair contract bindings and helpers.

## Getting Started

### Prerequisites

- Go 1.20+ (1.22 recommended)
- An Ethereum JSON-RPC endpoint (for example, Alchemy or Infura)

### Install

```bash
go mod download
```

### Run

```bash
export MAINNET_RPC="https://eth-mainnet.g.alchemy.com/v2/<KEY>"
go run .
```

## Configuration

- `MAINNET_RPC`: Required. JSON-RPC endpoint used for chain access.

## Roadmap

- Price discovery and pool graph construction
- Path finding and profitability checks (gas-aware)
- Transaction simulation and bundle submission
- Execution with safety checks and circuit breakers
- Observability (metrics, logs, alerts)

## Disclaimer

This software is for research and educational use only. It is not financial advice. Use at your own risk.
