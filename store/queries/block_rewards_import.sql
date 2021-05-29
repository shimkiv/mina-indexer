INSERT INTO block_rewards (
  public_key,
  delegate,
  block_height,
  block_time,
  reward,
  owner_type
)
VALUES @values

ON CONFLICT (public_key, delegate, block_height, owner_type) DO UPDATE
SET
  reward       = excluded.reward
