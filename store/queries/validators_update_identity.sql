UPDATE validators
SET
  identity_name = (CASE LENGTH($2) WHEN 0 THEN NULL ELSE $2 END)
WHERE
  public_key = $1
