CREATE TABLE experience_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        TEXT NOT NULL UNIQUE,
    title_fa    TEXT NOT NULL,
    title_en    TEXT NOT NULL,
    icon        TEXT,
    sort_order  INT NOT NULL DEFAULT 0
);
