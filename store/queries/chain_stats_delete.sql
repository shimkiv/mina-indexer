DELETE FROM chain_stats
WHERE
  time = DATE_TRUNC('@bucket', ?::timestamp)
  AND bucket = '@bucket'
