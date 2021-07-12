SELECT
 public_key
FROM
 validators
WHERE
  id NOT IN (SELECT validator_id FROM validator_stats WHERE bucket = ? AND time = ?)