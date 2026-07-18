-- Seed data: default brand, categories, sample hosts and published
-- experiences so the platform has real content on first boot.

INSERT INTO brands (slug, name, tagline, description, logo_url, primary_color, secondary_color, accent_color, default_locale, supported_locales, default_currency, contact_phone, contact_email, domain)
VALUES (
    'didar',
    'دیدار',
    'تجربه‌های واقعی، از دل ایران',
    'دیدار پلتفرمی برای کشف و رزرو تجربه‌های محلی و اصیل ایرانی است؛ از آشپزی خانگی در تهران تا سفالگری در اصفهان.',
    '/brand/didar/logo.svg',
    '#708D23',
    '#3D4F14',
    '#95A844',
    'fa',
    ARRAY['fa', 'en'],
    'IRR',
    '+982100000000',
    'hello@didar.travel',
    'localhost'
);

INSERT INTO experience_categories (slug, title_fa, title_en, icon, sort_order) VALUES
    ('cooking',    'آشپزی و غذا',        'Food & Cooking',      'utensils', 1),
    ('handicraft', 'صنایع‌دستی و هنر',    'Art & Handicraft',    'palette',  2),
    ('nature',     'طبیعت‌گردی',          'Nature & Outdoors',   'mountain', 3),
    ('history',    'تاریخ و فرهنگ',       'History & Culture',   'landmark', 4),
    ('music',      'موسیقی و اجرا',       'Music & Performance', 'music',    5),
    ('adventure',  'ماجراجویی',           'Adventure',           'compass',  6);

-- Hosts
INSERT INTO users (phone, full_name, avatar_url, is_host, status) VALUES
    ('+989121110001', 'آرش کریمی',    'https://i.pravatar.cc/200?u=arash-karimi',    true, 'active'),
    ('+989121110002', 'مریم حسینی',   'https://i.pravatar.cc/200?u=maryam-hosseini', true, 'active'),
    ('+989121110003', 'رضا صادقی',    'https://i.pravatar.cc/200?u=reza-sadeghi',    true, 'active'),
    ('+989121110004', 'نگار احمدی',   'https://i.pravatar.cc/200?u=negar-ahmadi',    true, 'active'),
    ('+989121110005', 'سینا مرادی',   'https://i.pravatar.cc/200?u=sina-moradi',     true, 'active'),
    ('+989121110006', 'فاطمه رستمی',  'https://i.pravatar.cc/200?u=fatemeh-rostami', true, 'active');

INSERT INTO host_profiles (user_id, brand_id, display_name, bio, city, is_verified, rating_avg, rating_count)
SELECT u.id, b.id, u.full_name, bio, city, true, rating_avg, rating_count
FROM (VALUES
    ('+989121110001', 'سفال‌گری در کوچه‌باغ‌های اصفهان را از پدربزرگم یاد گرفتم و بیش از ۱۵ سال است در کارگاه خودم مهمان‌نوازم.', 'اصفهان', 4.90, 128),
    ('+989121110002', 'آشپز خانگی با عشق به غذاهای اصیل ایرانی؛ میزبان کلاس‌های آشپزی خانوادگی در آپارتمانم در تهران.', 'تهران', 4.80, 95),
    ('+989121110003', 'راهنمای گردشگری و عاشق شعر و باغ‌های شیراز، هر عصر در باغ ارم قدم می‌زنم و از تاریخ فارس می‌گویم.', 'شیراز', 4.95, 210),
    ('+989121110004', 'معمار و پژوهشگر بافت تاریخی یزد، عاشق بادگیرها و معماری خاک ایران.', 'یزد', 4.85, 76),
    ('+989121110005', 'طبیعت‌گرد بومی گیلان، تورهای پیاده‌روی در جنگل‌های ماسوله و روستاهای ابری برگزار می‌کنم.', 'رشت', 4.75, 64),
    ('+989121110006', 'قالی‌باف نسل سوم خانواده‌ای هنرمند در کاشان، کارگاه سنتی قالی‌بافی را به مهمانان نشان می‌دهم.', 'کاشان', 4.92, 143)
) AS host_data(phone, bio, city, rating_avg, rating_count)
JOIN users u ON u.phone = host_data.phone
JOIN brands b ON b.slug = 'didar';

-- Experiences
INSERT INTO experiences (
    brand_id, host_id, category_id, slug, title, summary, description, language,
    price_amount, currency, duration_minutes, max_guests, min_guests, city, address, location, status,
    rating_avg, rating_count
)
SELECT
    b.id, u.id, c.id, e.slug, e.title, e.summary, e.description, 'fa',
    e.price_amount, 'IRR', e.duration_minutes, e.max_guests, e.min_guests, e.city, e.address,
    ST_SetSRID(ST_MakePoint(e.lon, e.lat), 4326)::geography, 'published',
    e.rating_avg, e.rating_count
FROM (VALUES
    (
        'pottery-workshop-isfahan',
        'کارگاه سفالگری سنتی در اصفهان',
        'یک بعدازظهر در کارگاه خانوادگی سفالگری، چرخ‌کار شوید و اثر خودتان را با خود ببرید.',
        'در این تجربه در کوچه‌باغ‌های محله جوباره اصفهان، با تاریخ سفالگری ایرانی آشنا می‌شوید و زیر نظر آرش، خودتان روی چرخ سفالگری کار می‌کنید. در پایان، اثر پخته‌شده خودتان را با پست برایتان ارسال می‌کنیم. چای و شیرینی محلی هم پذیرای شما خواهد بود.',
        1850000, 180, 6, 1, 'اصفهان', 'محله جوباره، اصفهان',
        51.6650, 32.6470, '+989121110001', 'handicraft', 4.90, 128
    ),
    (
        'home-cooking-class-tehran',
        'کلاس آشپزی خانگی؛ طعم‌های اصیل ایرانی',
        'یاد بگیرید چطور قرمه‌سبزی، تهچین و آش‌رشته اصیل بپزید، همراه با یک خانواده تهرانی.',
        'در آشپزخانه گرم و خانگی مریم در محله‌ی پاسداران تهران، دست‌به‌دست هم غذای سنتی ایرانی می‌پزیم؛ از پاک‌کردن سبزی معطر تا دم‌کردن برنج زعفرانی. در پایان همگی سر سفره می‌نشینیم و از دستپخت خودمان لذت می‌بریم.',
        1450000, 210, 8, 2, 'تهران', 'محله پاسداران، تهران',
        51.4530, 35.7719, '+989121110002', 'cooking', 4.80, 95
    ),
    (
        'poetry-garden-walk-shiraz',
        'قدم‌زنی شاعرانه در باغ ارم و آرامگاه حافظ',
        'عصرگاهی آرام در باغ‌های شیراز، همراه با شعرخوانی و روایت تاریخ فارس.',
        'با رضا در باغ ارم قدم می‌زنیم، از تاریخ باغ‌سازی ایرانی می‌شنویم و در مسیر به آرامگاه حافظ می‌رویم تا زیر سایه سروها فال حافظ بگیریم. این تجربه برای دوستداران شعر و تاریخ فارس طراحی شده است.',
        980000, 150, 10, 1, 'شیراز', 'باغ ارم، شیراز',
        52.5211, 29.6428, '+989121110003', 'history', 4.95, 210
    ),
    (
        'windcatcher-architecture-yazd',
        'کشف بادگیرها و بافت تاریخی یزد',
        'تور پیاده در کوچه‌های خشتی یزد و بازدید از بادگیرها، آب‌انبارها و خانه‌های تاریخی.',
        'نگار به‌عنوان معمار و پژوهشگر بافت یزد، رازهای معماری همساز با کویر را برایتان می‌گوید: از سازوکار بادگیرها تا زندگی در خانه‌های خشتی. این پیاده‌روی از میدان امیرچخماق شروع می‌شود.',
        1200000, 150, 8, 1, 'یزد', 'میدان امیرچخماق، یزد',
        54.3685, 31.8974, '+989121110004', 'history', 4.85, 76
    ),
    (
        'masuleh-hiking-tour',
        'پیاده‌روی در جنگل‌های ابری ماسوله',
        'یک روز کامل طبیعت‌گردی در روستای پله‌ای ماسوله و جنگل‌های اطراف آن.',
        'همراه سینا از روستای ماسوله شروع می‌کنیم و در مسیرهای جنگلی گیلان قدم می‌زنیم؛ از چشمه‌ها و آبشارهای کوچک عبور می‌کنیم و با گیاهان و پرندگان محلی آشنا می‌شویم. ناهار محلی گیلانی هم بخشی از این تجربه است.',
        1650000, 360, 6, 2, 'رشت', 'روستای ماسوله، فومن',
        48.9890, 37.1533, '+989121110005', 'nature', 4.75, 64
    ),
    (
        'traditional-carpet-weaving-kashan',
        'کارگاه قالی‌بافی سنتی در کاشان',
        'دستان خود را روی دار قالی بگذارید و گره‌های اول قالی‌بافی سنتی ایرانی را یاد بگیرید.',
        'فاطمه، قالی‌باف نسل سومِ خانواده‌ای هنرمند در کاشان، شما را با ابزار و رنگ‌های طبیعی قالی‌بافی آشنا می‌کند و زیر نظر او چند گره اول قالی خودتان را می‌بافید. داستان قالی ایرانی و نمادهای آن را هم خواهید شنید.',
        1350000, 150, 5, 1, 'کاشان', 'بازار سنتی کاشان',
        51.4100, 33.9850, '+989121110006', 'handicraft', 4.92, 143
    ),
    (
        'persian-tea-and-sweets-tehran',
        'مراسم چای‌دم‌کردن و شیرینی‌پزی ایرانی',
        'دم‌کردن چای اصیل ایرانی به‌سبک سماور و پختن شیرینی‌های سنتی در پاتوقی قدیمی تهران.',
        'در این تجربه با مریم، رسم چای‌دم‌کردن با سماور زغالی، تعارفات ایرانی و پختن چند نوع شیرینی خانگی مثل نان برنجی و باقلوا را یاد می‌گیرید. مناسب برای علاقه‌مندان به فرهنگ پذیرایی ایرانی.',
        890000, 120, 6, 1, 'تهران', 'محله پاسداران، تهران',
        51.4460, 35.7750, '+989121110002', 'cooking', 4.70, 52
    ),
    (
        'starlight-desert-safari-yazd',
        'سفاری کویر و شب‌نجوم در کوره‌های یزد',
        'غروب در کویر مرکزی ایران، آتش‌نشینی دور آتش و رصد ستارگان دور از آلودگی نور شهر.',
        'همراه نگار از یزد به سمت کویر مرکزی حرکت می‌کنیم، غروب کویر را تجربه می‌کنیم، دور آتش می‌نشینیم و با تلسکوپ کوچک همراه، آسمان پرستاره کویر را رصد می‌کنیم. یک شب چای و خرمای محلی هم داریم.',
        2100000, 300, 8, 2, 'یزد', 'کویر مرکزی، اطراف یزد',
        54.4200, 31.8300, '+989121110004', 'adventure', 4.88, 39
    )
) AS e(slug, title, summary, description, price_amount, duration_minutes, max_guests, min_guests, city, address, lon, lat, host_phone, category_slug, rating_avg, rating_count)
JOIN brands b ON b.slug = 'didar'
JOIN users u ON u.phone = e.host_phone
JOIN experience_categories c ON c.slug = e.category_slug;

-- Cover photo + gallery for every experience (deterministic placeholder images).
INSERT INTO experience_photos (experience_id, url, alt, sort_order, is_cover)
SELECT x.id, photo.url, x.title, photo.sort_order, photo.is_cover
FROM experiences x
JOIN LATERAL (
    VALUES
        ('https://picsum.photos/seed/' || x.slug || '-1/1200/800', 0, true),
        ('https://picsum.photos/seed/' || x.slug || '-2/1200/800', 1, false),
        ('https://picsum.photos/seed/' || x.slug || '-3/1200/800', 2, false)
) AS photo(url, sort_order, is_cover) ON true;

-- Upcoming bookable slots (next 3 Fridays/Saturdays) for every experience.
INSERT INTO experience_slots (experience_id, start_at, end_at, capacity, booked_count)
SELECT
    x.id,
    d.start_at,
    d.start_at + (x.duration_minutes || ' minutes')::interval,
    x.max_guests,
    0
FROM experiences x
JOIN LATERAL (
    VALUES
        (date_trunc('day', now()) + interval '3 days' + interval '17 hours'),
        (date_trunc('day', now()) + interval '7 days' + interval '17 hours'),
        (date_trunc('day', now()) + interval '10 days' + interval '10 hours'),
        (date_trunc('day', now()) + interval '14 days' + interval '10 hours')
) AS d(start_at) ON true;
