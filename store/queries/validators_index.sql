WITH validator_fees AS (
  SELECT DISTINCT ON(account_address) account_address, validator_fee
  FROM validator_epochs
  ORDER BY account_address, epoch DESC
)
SELECT
  validators.public_key,
  validators.identity_name,
  validators.start_height,
  validators.start_time,
  validators.last_height,
  validators.last_time,
  validators.blocks_created,
  validators.blocks_proposed,
  validators.delegations,
  COALESCE(validators.stake, 0)::TEXT AS stake,
  COALESCE(accounts.balance, 0)::TEXT AS account_balance,
  COALESCE(accounts.balance_unknown, 0)::TEXT AS account_balance_unknown,
  validator_fees.validator_fee AS fee
FROM
  validators
LEFT JOIN accounts
  ON accounts.public_key = validators.public_key
LEFT JOIN validator_fees
  ON validator_fees.account_address = validators.public_key
ORDER BY
  blocks_created DESC
