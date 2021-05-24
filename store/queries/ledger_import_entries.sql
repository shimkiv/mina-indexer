INSERT INTO ledger_entries (
  ledger_id,
  public_key,
  delegate,
  delegation,
  balance,
  weight,
  timing_initial_minimum_balance,
  timing_cliff_time,
  timing_cliff_amount,
  timing_vesting_period,
  timing_vesting_increment
)
VALUES @values