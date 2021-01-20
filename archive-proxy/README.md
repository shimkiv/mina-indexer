# mina-archive-api

Archive data API for Mina

## API Reference

| Method | Path             | Description
|--------|------------------|---------------------------------------------------
| GET    | /                | Blockchain stats
| GET    | /status          | Archive db status
| GET    | /blocks          | Get all blocks
| GET    | /blocks/:hash    | Get block details by hash
| GET    | /public_keys     | Get all public keys
| GET    | /public_keys/:id | Get public keys details by value
| GET    | /staking_ledger  | Get staking ledger dump