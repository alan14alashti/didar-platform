-- Additional published experiences inspired by popular Jabama / Airbnb-style
-- local tours across Iran (food, craft, music, nature, culture).

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
        'tabriz-bazaar-food-walk',
        'گشت غذایی در بازار تاریخی تبریز',
        'از کوفته تبریزی تا نوقا و چای قورباغی؛ طعم‌های اصیل آذربایجان را در بازار بزرگ بچشید.',
        'در این تجربه مانند تورهای محبوب جاباما، در بازار تاریخی تبریز قدم می‌زنیم و در چند غرفه و رستوران محلی، غذاهای سنتی آذربایجان را می‌چشیم. راهنمای محلی داستان هر غذا و سنت پذیرایی را تعریف می‌کند.',
        1180000, 150, 8, 2, 'تبریز', 'بازار بزرگ تبریز',
        46.2919, 38.0800, '+989121110002', 'cooking', 4.86, 67
    ),
    (
        'persian-calligraphy-workshop-tehran',
        'کارگاه خوشنویسی و نستعلیق در تهران',
        'قلم و مرکب بگیرید و خط نستعلیق را در کارگاهی دنج در دل تهران یاد بگیرید.',
        'در این کلاس نیم‌روزه، استاد خوشنویسی مفاهیم پایه نستعلیق، نحوه نشستن، گرفتن قلم و نوشتن نام خود را آموزش می‌دهد. مناسب برای علاقه‌مندان به هنرهای سنتی ایرانی.',
        1250000, 120, 6, 1, 'تهران', 'محله سنگلج، تهران',
        51.4180, 35.6785, '+989121110001', 'handicraft', 4.91, 84
    ),
    (
        'live-traditional-music-isfahan',
        'شب موسیقی سنتی در خانه تاریخی اصفهان',
        'اجرای زنده موسیقی ایرانی در حیاط خانه‌ای قدیمی نزدیک میدان نقش جهان.',
        'شبی شبیه تجربه‌های Airbnb Experiences: در حیاط خانه‌ای تاریخی اصفهان، نوازندگان محلی موسیقی دستگاهی اجرا می‌کنند و میزبان درباره سازها و تاریخ موسیقی اصفهان می‌گوید. چای و شیرینی محلی سرو می‌شود.',
        1580000, 120, 12, 2, 'اصفهان', 'خیابان چهارباغ، اصفهان',
        51.6746, 32.6539, '+989121110001', 'music', 4.93, 112
    ),
    (
        'caspian-sea-boat-tour-anzali',
        'قایق‌سواری و تماشای تالاب انزلی',
        'صبحی روی آب؛ پرندگان، نیزارها و آرامش ساحل شمال.',
        'مانند تورهای ساحلی جاباما، با قایق در تالاب انزلی حرکت می‌کنیم، پرندگان مهاجر را از نزدیک می‌بینیم و در یک رستوران محلی کنار آب، ماهی تازه شمال می‌خوریم.',
        1750000, 240, 6, 2, 'بندر انزلی', 'اسکله بندر انزلی',
        49.4627, 37.4724, '+989121110005', 'nature', 4.82, 58
    ),
    (
        'kaluts-sunrise-kerman',
        'طلوع آفتاب در کلوت‌های شهداد',
        'سفر سحرگاهی به بیابان‌های Mars-like کرمان با صبحانه محلی.',
        'قبل از طلوع از کرمان حرکت می‌کنیم و در کلوت‌های شهداد طلوع را تماشا می‌کنیم. راهنمای محلی زمین‌شناسی منطقه را توضیح می‌دهد و صبحانه سنتی کرمان سرو می‌شود.',
        2200000, 300, 8, 2, 'کرمان', 'کلوت‌های شهداد، کرمان',
        57.7490, 30.4410, '+989121110004', 'adventure', 4.89, 41
    ),
    (
        'lavash-bread-baking-qazvin',
        'پخت نان لواش در تنور روستایی قزوین',
        'در روستایی نزدیک قزوین، خمیر بزنید و نان تازه را از تنور بیرون بکشید.',
        'در این تجربه خانوادگی، مراحل آماده‌سازی خمیر، پهن کردن و پخت در تنور گلی را یاد می‌گیرید و با خانواده میزبان نان تازه و پنیر محلی می‌خورید.',
        950000, 150, 8, 2, 'قزوین', 'روستای ویر، قزوین',
        49.8760, 36.0540, '+989121110002', 'cooking', 4.77, 36
    ),
    (
        'meymand-rocky-village-tour',
        'بازدید از روستای صخره‌ای میمند کرمان',
        'گشت در خانه‌های کنده‌شده از سنگ و آشنایی با زندگی بومی کویر.',
        'در این تور فرهنگی، در روستای تاریخی میمند قدم می‌زنیم، یک خانه محلی را از داخل می‌بینیم و میزبان درباره معماری بومی و سنت‌های عشایری توضیح می‌دهد.',
        1400000, 180, 10, 1, 'کرمان', 'روستای میمند، کرمان',
        55.4410, 30.1680, '+989121110004', 'history', 4.88, 53
    ),
    (
        'saffron-harvest-torbat',
        'برداشت زعفران در مزارع تربت حیدریه',
        'در فصل برداشت، در مزرعه زعفران مشغول چیدن شکوفه‌ها شوید و چای زعفران بنوشید.',
        'یک تجربه فصلی منحصربه‌فرد: در مزارع اطراف تربت حیدریه با کشاورزان محلی زعفران می‌چینید، فرآیند خشک‌کردن را می‌بینید و درباره گران‌ترین ادویه جهان یاد می‌گیرید.',
        1680000, 210, 8, 2, 'تربت حیدریه', 'مزارع زعفران، تربت حیدریه',
        59.2190, 35.2740, '+989121110006', 'nature', 4.94, 29
    ),
    (
        'kashan-rosewater-distillery',
        'گلاب‌گیری سنتی در قمصر کاشان',
        'در فصل گل، کارگاه گلاب‌گیری را ببینید و عطر گل محمدی را حس کنید.',
        'در قمصر کاشان وارد کارگاه گلاب‌گیری سنتی می‌شویم، فرآیند تقطیر گل را از نزدیک می‌بینیم و بسته‌ای گلاب اصلیل همراه خود می‌بریم.',
        1120000, 120, 10, 1, 'کاشان', 'قمصر، کاشان',
        51.4350, 33.9590, '+989121110006', 'handicraft', 4.83, 47
    ),
    (
        'tehran-grand-bazaar-tour',
        'تور معماری و تاریخ در بازار بزرگ تهران',
        'از سرای مروی تا گنبد مینا؛ داستان بازار قدیمی پایتخت.',
        'راهنمای محلی شما را در بازار تهران همراهی می‌کند، سراهای تاریخی، حمام‌ها و کاروانسراها را نشان می‌دهد و در یک چایخانه قدیمی توقف می‌کنید.',
        780000, 150, 8, 1, 'تهران', 'بازار بزرگ تهران',
        51.421, 35.672, '+989121110003', 'history', 4.79, 91
    )
) AS e(slug, title, summary, description, price_amount, duration_minutes, max_guests, min_guests, city, address, lon, lat, host_phone, category_slug, rating_avg, rating_count)
JOIN brands b ON b.slug = 'didar'
JOIN users u ON u.phone = e.host_phone
JOIN experience_categories c ON c.slug = e.category_slug
ON CONFLICT (slug) DO NOTHING;

INSERT INTO experience_photos (experience_id, url, alt, sort_order, is_cover)
SELECT x.id, photo.url, x.title, photo.sort_order, photo.is_cover
FROM experiences x
JOIN LATERAL (
    VALUES
        ('https://picsum.photos/seed/' || x.slug || '-1/1200/800', 0, true),
        ('https://picsum.photos/seed/' || x.slug || '-2/1200/800', 1, false)
) AS photo(url, sort_order, is_cover) ON true
WHERE x.slug IN (
    'tabriz-bazaar-food-walk', 'persian-calligraphy-workshop-tehran', 'live-traditional-music-isfahan',
    'caspian-sea-boat-tour-anzali', 'kaluts-sunrise-kerman', 'lavash-bread-baking-qazvin',
    'meymand-rocky-village-tour', 'saffron-harvest-torbat', 'kashan-rosewater-distillery', 'tehran-grand-bazaar-tour'
)
AND NOT EXISTS (SELECT 1 FROM experience_photos p WHERE p.experience_id = x.id);

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
        (date_trunc('day', now()) + interval '2 days' + interval '10 hours'),
        (date_trunc('day', now()) + interval '5 days' + interval '17 hours'),
        (date_trunc('day', now()) + interval '9 days' + interval '10 hours')
) AS d(start_at) ON true
WHERE x.slug IN (
    'tabriz-bazaar-food-walk', 'persian-calligraphy-workshop-tehran', 'live-traditional-music-isfahan',
    'caspian-sea-boat-tour-anzali', 'kaluts-sunrise-kerman', 'lavash-bread-baking-qazvin',
    'meymand-rocky-village-tour', 'saffron-harvest-torbat', 'kashan-rosewater-distillery', 'tehran-grand-bazaar-tour'
)
AND NOT EXISTS (SELECT 1 FROM experience_slots s WHERE s.experience_id = x.id);
