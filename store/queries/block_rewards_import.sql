INSERT INTO block_rewards (
  owner_account,
  delegate,
  block_height,
  block_time,
  reward,
  owner_type
)
VALUES @values

ON CONFLICT (owner_account, delegate, block_height, owner_type) DO UPDATE
SET
  reward       = excluded.reward
