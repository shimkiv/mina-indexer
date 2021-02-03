INSERT INTO accounts (
  public_key,
  delegate,
  balance,
  balance_unknown,
  nonce,
  start_height,
  start_time,
  last_height,
  last_time,
  created_at,
  updated_at
)
VALUES @values
ON CONFLICT (public_key) DO UPDATE
SET
  delegate        = excluded.delegate,
  balance         = excluded.balance,
  balance_unknown = excluded.balance_unknown,
  nonce           = excluded.nonce,
  last_height     = excluded.last_height,
  last_time       = excluded.last_time,
  updated_at      = excluded.updated_at