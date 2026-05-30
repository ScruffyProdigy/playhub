-- Mark existing auto-generated display names as provisional.
UPDATE users
SET display_name = display_name || ' (new)',
    updated_at = NOW()
WHERE display_name = split_part(email, '@', 1)
  AND display_name NOT LIKE '% (new)';
