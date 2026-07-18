-- Twenty additional published experiences across Iran (food, craft, music, nature, culture).

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
        'shiraz-hafez-poetry-evening',
        'شب شعر و فال حافظ در آرامگاه',
        'شبی در آرامگاه حافظ با شعرخوانی و فال حافظ.',
        'تجربه‌ای فرهنگی در شیراز: در فضای آرامگاه حافظ، شعرخوانی زنده، فال حافظ و گفت‌وگو درباره شعر و فرهنگ فارسی.',
        920000, 120, 12, 1, 'شیراز', 'آرامگاه حافظ',
        52.558, 29.624, '+989121110003', 'history', 4.96, 178
    ),
    (
        'tehran-street-food-night',
        'گشت شبانه غذاهای خیابانی تهران',
        'از ساندویچ بندری تا آش رشته در کوچه‌های قدیمی.',
        'تور غذای خیابانی تهران: در محله‌های تاریخی از غذاهای محبوب خیابانی بچشید و داستان هر غذا را بشنوید.',
        850000, 180, 8, 2, 'تهران', 'خیابان ۳۰ تیر',
        51.422, 35.689, '+989121110002', 'cooking', 4.84, 142
    ),
    (
        'lahijan-tea-garden-tour',
        'پیاده‌روی در باغات چای لاهیجان',
        'از مزرعه تا فنجان؛ مسیر چای شمال.',
        'بازدید از باغ چای گیلان: پیاده‌روی در مزارع چای، آشنایی با فرآیند چیدن برگ و چای‌نوشی در خانه محلی.',
        1050000, 150, 10, 2, 'لاهیجان', 'باغات چای لاهیجان',
        50.004, 37.207, '+989121110005', 'nature', 4.81, 55
    ),
    (
        'hormoz-rainbow-island',
        'گشت جزیره رنگارنگ هرمز',
        'خاک سرخ، سواحل نمکی و جزیره‌ای از رنگ‌ها.',
        'تور یک‌روزه هرمز: کشف خاک سرخ، سواحل نمکی و روستای رنگارنگ جزیره هرمز با راهنمای محلی.',
        1950000, 360, 8, 2, 'هرمز', 'جزیره هرمز',
        56.452, 27.094, '+989121110004', 'adventure', 4.90, 88
    ),
    (
        'qeshm-star-valley-sunset',
        'غروب در دره ستارگان قشم',
        'پیاده‌روی در دره‌ای با سنگ‌های کهکشانی.',
        'طبیعت‌گردی دره ستارگان: پیاده‌روی در میان سنگ‌های فرسایش‌یافته و تماشای غروب در جزیره قشم.',
        1380000, 180, 10, 2, 'قشم', 'دره ستارگان',
        55.770, 26.750, '+989121110005', 'adventure', 4.87, 73
    ),
    (
        'rasht-gilani-feast',
        'پخت و چشیدن غذاهای گیلانی',
        'میرزا قاسمی، باقالا قاتوق و ماهی دودی.',
        'کلاس آشپزی رشت: پخت غذاهای محلی گیلان در آشپزخانه خانگی و سفره‌نشینی با میزبان.',
        1320000, 210, 8, 2, 'رشت', 'محله گلسار',
        49.589, 37.268, '+989121110005', 'cooking', 4.92, 104
    ),
    (
        'mashhad-breakfast-pilgrimage',
        'صبحانه سنتی مشهد پس از زیارت',
        'شله‌زرد، نان سنگک و چای زعفرانی در حجره قدیمی.',
        'تجربه صبحانه مشهد: صبحانه سنتی خراسانی در فضایی نزدیک حرم با داستان آداب و رسوم محلی.',
        680000, 90, 10, 2, 'مشهد', 'حرم مطهر',
        59.615, 36.288, '+989121110002', 'cooking', 4.78, 156
    ),
    (
        'hamedan-alavian-tour',
        'کشف همدان؛ الوند و گنجنامه',
        'تور تاریخی هگمتانه و صخره‌های گنجنامه.',
        'تاریخ و طبیعت همدان: بازدید از گنجنامه، داستان هگمتانه و پیاده‌روی در دامنه الوند.',
        990000, 240, 10, 1, 'همدان', 'گنجنامه',
        48.515, 34.798, '+989121110003', 'history', 4.83, 62
    ),
    (
        'sanandaj-kurdish-music',
        'شب موسیقی کردی در سنندج',
        'اجرای زنده دهل و تمبور با چای محلی.',
        'موسیقی سنتی کردستان: اجرای زنده در خانه تاریخی سنندج همراه با چای و شیرینی محلی.',
        1150000, 120, 14, 2, 'سنندج', 'خانه تاریخی سنندج',
        46.998, 35.314, '+989121110001', 'music', 4.88, 47
    ),
    (
        'isfahan-miniatu-workshop',
        'کارگاه نقاشی مینیاتور اصفهان',
        'طراحی گل و مرغ با قلم مویی.',
        'هنر مینیاتور ایرانی: آموزش پایه طراحی و رنگ‌آمیزی مینیاتور در کارگاهی نزدیک میدان نقش جهان.',
        1480000, 150, 6, 1, 'اصفهان', 'میدان نقش جهان',
        51.677, 32.657, '+989121110001', 'handicraft', 4.94, 91
    ),
    (
        'chabahar-mud-volcanoes',
        'کوه‌های گلی و سواحل چابهار',
        'سفر به یکی از شگفت‌انگیزترین سواحل ایران.',
        'ماجراجویی در سیستان: بازدید از کوه‌های گلی، سواحل ماسه‌ای و فرهنگ بلوچستان.',
        2400000, 420, 6, 2, 'چابهار', 'کوه‌های گلی',
        60.643, 25.292, '+989121110004', 'adventure', 4.85, 34
    ),
    (
        'urmia-lake-birdwatch',
        'تماشای پرندگان دریاچه ارومیه',
        'تور پرنده‌نگری با راهنمای محلی.',
        'طبیعت‌گردی دریاچه ارومیه: مشاهده پرندگان مهاجر و آشنایی با اکوسیستم دریاچه.',
        1100000, 180, 8, 2, 'ارومیه', 'سواحل دریاچه ارومیه',
        45.076, 37.552, '+989121110006', 'nature', 4.76, 41
    ),
    (
        'naein-desert-citadel',
        'قلعه نارین و کویر نائین',
        'بازدید از قلعه ۲۰۰۰ساله و غروب کویر.',
        'تاریخ و کویر نائین: گشت در قلعه نارین و تماشای غروب در دشت کویر.',
        1250000, 210, 8, 2, 'نائین', 'قلعه نارین',
        53.088, 32.860, '+989121110004', 'history', 4.86, 38
    ),
    (
        'abadeh-wood-inlay',
        'هنر خاتم‌کاری و منبت در آباده',
        'ساخت یک جعبه خاتم کوچک.',
        'صنایع‌دستی فارس: آشنایی با خاتم‌کاری و ساخت یک اثر کوچک در کارگاه آباده.',
        1190000, 150, 5, 1, 'آباده', 'بازار آباده',
        52.650, 31.160, '+989121110001', 'handicraft', 4.80, 33
    ),
    (
        'amol-hyrcanian-forest',
        'پیاده‌روی جنگل‌های هیرکانی آمل',
        'مسیر آبشار و درختان کهنسال.',
        'طبیعت‌گردی جنگل شمال: پیاده‌روی در جنگل‌های هیرکانی و توقف کنار آبشار.',
        1420000, 300, 8, 2, 'آمل', 'جنگل‌های اطراف آمل',
        52.351, 36.469, '+989121110005', 'nature', 4.89, 57
    ),
    (
        'kish-coral-reef-snorkel',
        'غواصی سطحی در آب‌های مرجانی کیش',
        'مشاهده مرجان و ماهی‌های رنگی.',
        'تجربه دریایی کیش: غواصی سطحی در آب‌های شفاف و مشاهده مرجان‌های رنگارنگ.',
        1850000, 180, 6, 2, 'کیش', 'ساحل مرجان',
        54.019, 26.557, '+989121110005', 'adventure', 4.91, 69
    ),
    (
        'ahvaz-karun-river-walk',
        'غروب کنار کارون و بازار عربی',
        'پیاده‌روی کنار رود و چشیدن ماهی‌شک.',
        'فرهنگ و طبیعت اهواز: پیاده‌روی کنار کارون، بازدید از بازار محلی و چشیدن غذاهای جنوبی.',
        870000, 150, 10, 2, 'اهواز', 'پل سفید',
        48.669, 31.320, '+989121110003', 'nature', 4.74, 48
    ),
    (
        'bojnurd-saffron-tea-ceremony',
        'مراسم چای زعفران و شیرینی بجنورد',
        'آموزش دم‌کردن چای زعفران خراسانی.',
        'تجربه چای شمال خراسان: مراسم دم‌آوری چای زعفران و شیرینی‌های محلی بجنورد.',
        760000, 90, 8, 2, 'بجنورد', 'بازار قدیم بجنورد',
        57.329, 37.475, '+989121110006', 'cooking', 4.72, 29
    ),
    (
        'zanjan-soltaniyeh-dome',
        'بازدید از گنبد سلطانیه',
        'شاهکار معماری ایلخانی با راهنمای تاریخ.',
        'تور فرهنگی زنجان: بازدید از گنبد سلطانیه و آشنایی با معماری ایلخانی.',
        940000, 180, 12, 1, 'زنجان', 'گنبد سلطانیه',
        48.796, 36.434, '+989121110004', 'history', 4.87, 54
    ),
    (
        'tehran-indie-music-jam',
        'جَم موسیقی مستقل در کافه تهران',
        'اجرای زنده بندهای مستقل در فضایی صمیمی.',
        'شب موسیقی معاصر تهران: اجرای زنده بندهای مستقل در کافه‌ای در مرکز تهران.',
        990000, 120, 20, 2, 'تهران', 'خیابان ایرانشهر',
        51.414, 35.707, '+989121110003', 'music', 4.82, 118
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
    'shiraz-hafez-poetry-evening', 'tehran-street-food-night', 'lahijan-tea-garden-tour',
    'hormoz-rainbow-island', 'qeshm-star-valley-sunset', 'rasht-gilani-feast',
    'mashhad-breakfast-pilgrimage', 'hamedan-alavian-tour', 'sanandaj-kurdish-music',
    'isfahan-miniatu-workshop', 'chabahar-mud-volcanoes', 'urmia-lake-birdwatch',
    'naein-desert-citadel', 'abadeh-wood-inlay', 'amol-hyrcanian-forest',
    'kish-coral-reef-snorkel', 'ahvaz-karun-river-walk', 'bojnurd-saffron-tea-ceremony',
    'zanjan-soltaniyeh-dome', 'tehran-indie-music-jam'
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
        (date_trunc('day', now()) + interval '3 days' + interval '10 hours'),
        (date_trunc('day', now()) + interval '7 days' + interval '17 hours'),
        (date_trunc('day', now()) + interval '10 days' + interval '10 hours'),
        (date_trunc('day', now()) + interval '14 days' + interval '17 hours')
) AS d(start_at) ON true
WHERE x.slug IN (
    'shiraz-hafez-poetry-evening', 'tehran-street-food-night', 'lahijan-tea-garden-tour',
    'hormoz-rainbow-island', 'qeshm-star-valley-sunset', 'rasht-gilani-feast',
    'mashhad-breakfast-pilgrimage', 'hamedan-alavian-tour', 'sanandaj-kurdish-music',
    'isfahan-miniatu-workshop', 'chabahar-mud-volcanoes', 'urmia-lake-birdwatch',
    'naein-desert-citadel', 'abadeh-wood-inlay', 'amol-hyrcanian-forest',
    'kish-coral-reef-snorkel', 'ahvaz-karun-river-walk', 'bojnurd-saffron-tea-ceremony',
    'zanjan-soltaniyeh-dome', 'tehran-indie-music-jam'
)
AND NOT EXISTS (SELECT 1 FROM experience_slots s WHERE s.experience_id = x.id);
