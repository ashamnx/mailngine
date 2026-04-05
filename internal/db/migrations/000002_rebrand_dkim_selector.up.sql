-- Change default DKIM selector from hm1 to mn1 (rebrand)
ALTER TABLE domains ALTER COLUMN dkim_selector SET DEFAULT 'mn1';
