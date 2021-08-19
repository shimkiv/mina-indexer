INSERT INTO validator_epochs (
  staketab_id,
  account_address,
  epoch,
  validator_fee
)
VALUES
  @values
ON CONFLICT (staketab_id, epoch) DO UPDATE
SET
  validator_fee   = excluded.validator_fee