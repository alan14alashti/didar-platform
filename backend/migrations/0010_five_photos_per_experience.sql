-- Ensure every experience has exactly 5 gallery photos; sort_order 0 is the cover.

INSERT INTO experience_photos (experience_id, url, alt, sort_order, is_cover)
SELECT
    x.id,
    'https://picsum.photos/seed/' || x.slug || '-' || (n.i + 1)::text || '/1200/800',
    x.title,
    n.i,
    n.i = 0
FROM experiences x
CROSS JOIN generate_series(0, 4) AS n(i)
WHERE NOT EXISTS (
    SELECT 1 FROM experience_photos p
    WHERE p.experience_id = x.id AND p.sort_order = n.i
);

UPDATE experience_photos SET is_cover = (sort_order = 0);
