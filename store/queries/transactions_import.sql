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
  failure_reason,
  sequence_number,
  secondary_sequence_number,
  created_at,
  updated_at
)
VALUES @values

ON CONFLICT (hash) DO UPDATE
SET
  updated_at = excluded.updated_at