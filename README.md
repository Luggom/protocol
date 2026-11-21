# Luggom

**PT** — O Luggom é uma blockchain **Cosmos SDK** (`x/luggom`) onde se **minera LUG pagando imposto real**, provado por **TLSNotary** — sem KYC, código aberto e permissionless.

**EN** — Luggom is a **Cosmos SDK** blockchain (`x/luggom`) where you **mine LUG by paying real taxes**, proven by **TLSNotary** — no KYC, open source and permissionless.

## Links · Ligações

- Website / Site: [EN](https://luggom.com/) · [PT](https://luggom.com/pt/)
- Whitepaper: [EN](https://luggom.com/whitepaper.html) · [PT](https://luggom.com/pt/whitepaper.html)
- Developers: [EN](https://luggom.com/developers.html) · [PT](https://luggom.com/pt/developers.html) · [Endpoints](TESTNET/ENDPOINTS.en.md)
- Tokenomics: [`docs/TOKENOMICS.md`](docs/TOKENOMICS.md)
- Testnet: [`TESTNET/README.md`](TESTNET/README.md) · [EN](TESTNET/README.en.md)
- Roadmap / Roteiro: [`TESTNET/ROADMAP.md`](TESTNET/ROADMAP.md) · [EN](TESTNET/ROADMAP.en.md)
- Add a bank / Adicionar um banco: [`TESTNET/DEVELOPER_GUIDE.md`](TESTNET/DEVELOPER_GUIDE.md) · [EN](TESTNET/DEVELOPER_GUIDE.en.md)

## Status

**PT** — Hoje há **dois bancos** comprovados de ponta a ponta — **Viacredi** e **Banco do Brasil** (tx `F2D532…`, altura 144732, `code 0`), no Brasil.

**EN** — Today there are **two banks** proven end to end — **Viacredi** and **Banco do Brasil** (tx `F2D532…`, height 144732, `code 0`), in Brazil.

## Repo structure · Estrutura do repo

- `x/luggom/` — the core Cosmos SDK module / o módulo Cosmos SDK core.
- `app/` — the application wiring / a montagem da aplicação.
- `proto/` — protobuf definitions / definições protobuf.
- `api/taxproof/`, `api/reclaimverify/`, `api/taxnature/` — the anti-fraud gate (receipt → attestation) / o portão anti-fraude.
- `cmd/luggomd/` — the node / o nó.
- `cmd/luggom-extract/` — the notary extractor / o extrator do notário.
- `cmd/luggom-oracle-feeder/`, `cmd/gen-oracle-history/`, `cmd/sandbox-bank/` — oracle and rail tooling / ferramentas de oráculo e rails.
- `notary/` — notary attestation spec and sample / spec da atestação e amostra.
- `deps/cosmos-sdk/` — the vendored Cosmos SDK fork (export tally) / o fork do Cosmos SDK (export tally).
- `stubs/` — local Go stubs needed by the proto tooling / stubs Go locais do tooling de proto.
- `docs/TOKENOMICS.md` — the canonical mining split, pool and validator rules / o split canónico, a pool e as regras de validadores.
- `TESTNET/` — the public testnet package and endpoints (PT/EN) / o pacote público e endpoints da testnet (PT/EN).

The desktop app, the Android APK and the browser collector are proprietary
clients distributed as signed binaries; their source is not part of this
repository.
Os clientes desktop, APK Android e coletor de browser são proprietários e
distribuídos como binários assinados; o código-fonte não faz parte deste
repositório.

## Build · Compilar

```sh
# Go 1.23.8
go build -ldflags="-checklinkname=0" -o build/luggomd ./cmd/luggomd
go build -ldflags="-checklinkname=0" -o build/luggom-extract ./cmd/luggom-extract
```

The repository vendors `deps/cosmos-sdk` and `stubs/`, so the Go build works as-is.
O repositório inclui `deps/cosmos-sdk` e `stubs/`, por isso o build Go funciona tal como está.

See [`TESTNET/README.md`](TESTNET/README.md) · [EN](TESTNET/README.en.md) for the full
run/validate/notarize/mine instructions.
Veja [`TESTNET/README.md`](TESTNET/README.md) · [EN](TESTNET/README.en.md) para as
instruções completas de rodar/validar/notariar/minerar.
