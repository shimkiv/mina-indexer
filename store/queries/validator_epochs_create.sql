INSERT INTO validator_epochs (account_address, epoch, validator_fee)
VALUES (
  ?,
  COALESCE((SELECT MAX(epoch) FROM ledgers), 0),
  ?
)
ON CONFLICT (account_address, epoch) DO UPDATE
SET validator_fee = excluded.validator_fee
