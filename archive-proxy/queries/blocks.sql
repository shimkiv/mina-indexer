SELECT
  blocks.height,
  blocks.state_hash,
  blocks.parent_hash,
  blocks.ledger_hash,
  snarked_ledger_hashes.value AS snarked_ledger_hash,
  creator_keys.value AS creator,
  winnder_keys.value AS winner,
  blocks.timestamp,
  blocks.global_slot_since_genesis,
  blocks.global_slot
FROM
  blocks
INNER JOIN public_keys creator_keys
  ON creator_keys.id = blocks.creator_id
INNER JOIN public_keys winnder_keys
  ON winnder_keys.id = blocks.block_winner_id
INNER JOIN snarked_ledger_hashes
  ON snarked_ledger_hashes.id = blocks.snarked_ledger_hash_id
WHERE
  blocks.height >= $1
ORDER BY
  blocks.height ASC
LIMIT
  $2