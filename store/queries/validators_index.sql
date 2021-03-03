SELECT
  validators.public_key,
  validators.identity_name,
  validators.start_height,
  validators.start_time,
  validators.last_height,
  validators.last_time,
  COALESCE(validators.stake::TEXT) AS stake,
  validators.blocks_created,
  validators.blocks_proposed,
  validators.delegations,
  COALESCE(validators.stake, 0)::TEXT AS stake,
  COALESCE(accounts.balance, 0)::TEXT AS account_balance,
  COALESCE(accounts.balance_unknown, 0)::TEXT AS account_balance_unknown
FROM
  validators
LEFT JOIN accounts
  ON accounts.public_key = validators.public_key
ORDER BY
  blocks_created DESC