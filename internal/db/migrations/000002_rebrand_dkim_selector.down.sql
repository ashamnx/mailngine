-- Revert DKIM selector default back to hm1
ALTER TABLE domains ALTER COLUMN dkim_selector SET DEFAULT 'hm1';
