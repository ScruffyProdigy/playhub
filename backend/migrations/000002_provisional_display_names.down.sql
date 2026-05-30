UPDATE users
SET display_name = regexp_replace(display_name, ' \(new\)$', ''),
    updated_at = NOW()
WHERE display_name LIKE '% (new)';
