INSERT INTO transactions (
  type,
  hash,
  block_hash,
  block_height,
  time,
  nonce,
  sender,
  receiver,
  amount,
  fee,
  memo,
  status,
  canonical,
  failure_reason,
  sequence_number,
  secondary_sequence_number,
  created_at,
  updated_at
)
VALUES @values

ON CONFLICT (hash) DO UPDATE
SET
  sender         = excluded.sender,
  receiver       = excluded.receiver,
  amount         = excluded.amount,
  fee            = excluded.fee,
  status         = excluded.status,
  canonical      = excluded.canonical,
  failure_reason = excluded.failure_reason,
  updated_at     = excluded.updated_at