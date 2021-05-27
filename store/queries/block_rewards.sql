SELECT
  to_char(block_time, $INTERVAL) AS interval,
  SUM(reward) AS amount
FROM
  block_rewards
WHERE
  public_key = ?
  AND delegate = ?
  AND block_time BETWEEN ? AND ?
  AND reward_owner_type = ?
GROUP BY
  to_char(block_time, $INTERVAL)
