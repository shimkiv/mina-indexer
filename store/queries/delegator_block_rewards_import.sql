INSERT INTO delegator_block_rewards (
  public_key,
  delegate,
  block_height,
  block_time,
  reward
)
VALUES @values

ON CONFLICT (public_key, delegate, block_height) DO UPDATE
SET
  block_time   = excluded.block_time,
  reward       = excluded.reward
