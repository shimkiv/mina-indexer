INSERT INTO block_rewards (
  public_key,
  delegate,
  block_height,
  block_time,
  reward,
  reward_owner_type
)
VALUES @values

ON CONFLICT (public_key, delegate, block_height, reward_owner_type) DO UPDATE
SET
  reward       = excluded.reward
