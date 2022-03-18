DELETE FROM chain_stats
WHERE
  time = DATE_TRUNC('@bucket', ?::CHAIN_TIME)
  AND bucket = '@bucket'
