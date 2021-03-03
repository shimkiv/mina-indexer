INSERT INTO validators (
  public_key,
  start_height,
  start_time,
  last_height,
  last_time,
  blocks_proposed,
  blocks_created,
  created_at,
  updated_at
)
VALUES
  @values
ON CONFLICT (public_key) DO UPDATE
SET
  last_height        = excluded.last_height,
  last_time          = excluded.last_time,
  blocks_proposed    = excluded.blocks_proposed,
  blocks_created     = COALESCE((SELECT COUNT(1) FROM blocks WHERE creator = excluded.public_key), 0),
  updated_at         = excluded.updated_at