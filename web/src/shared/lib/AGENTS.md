# Small shared mechanics

- Keep utilities narrowly named with clear behavior (time formatting and required
  text validation are existing examples).
- Do not import higher frontend layers or construct API/domain policy here.
- Share only genuinely repeated behavior; avoid a generic helper dumping ground.
- Keep formatting resilient to missing/invalid values where the caller contract
  allows them, and test meaningful boundary cases.
