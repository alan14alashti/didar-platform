import type { Experience } from './types'
import { demoPhotosFor } from './demoPhotos'

export const demoNow = new Date().toISOString()

export const demoHosts = {
  arash: {
    id: 'h1',
    full_name: 'آرش کریمی',
    avatar_url: 'https://i.pravatar.cc/200?u=arash-karimi',
    display_name: 'آرش کریمی',
    bio: 'سفالگر و هنرمند اصفهانی.',
    is_verified: true,
    rating_avg: 4.9,
    rating_count: 128,
  },
  maryam: {
    id: 'h2',
    full_name: 'مریم حسینی',
    avatar_url: 'https://i.pravatar.cc/200?u=maryam-hosseini',
    display_name: 'مریم حسینی',
    bio: 'آشپز خانگی تهرانی.',
    is_verified: true,
    rating_avg: 4.8,
    rating_count: 95,
  },
  reza: {
    id: 'h3',
    full_name: 'رضا صادقی',
    avatar_url: 'https://i.pravatar.cc/200?u=reza-sadeghi',
    display_name: 'رضا صادقی',
    bio: 'راهنمای فرهنگی شیراز.',
    is_verified: true,
    rating_avg: 4.95,
    rating_count: 210,
  },
  negar: {
    id: 'h4',
    full_name: 'نگار احمدی',
    avatar_url: 'https://i.pravatar.cc/200?u=negar-ahmadi',
    display_name: 'نگار احمدی',
    bio: 'معمار و پژوهشگر بافت تاریخی.',
    is_verified: true,
    rating_avg: 4.85,
    rating_count: 76,
  },
  sina: {
    id: 'h5',
    full_name: 'سینا مرادی',
    avatar_url: 'https://i.pravatar.cc/200?u=sina-moradi',
    display_name: 'سینا مرادی',
    bio: 'طبیعت‌گرد شمال ایران.',
    is_verified: true,
    rating_avg: 4.75,
    rating_count: 64,
  },
  fatemeh: {
    id: 'h6',
    full_name: 'فاطمه رستمی',
    avatar_url: 'https://i.pravatar.cc/200?u=fatemeh-rostami',
    display_name: 'فاطمه رستمی',
    bio: 'صنعت‌گر و قالی‌باف کاشانی.',
    is_verified: true,
    rating_avg: 4.92,
    rating_count: 143,
  },
} satisfies Record<string, NonNullable<Experience['host']>>

function photosFor(slug: string, title: string) {
  return demoPhotosFor(slug, title)
}

function slotsFor(experienceId: string, durationMinutes: number, maxGuests: number) {
  const offsets = [3, 7, 10, 14]
  return offsets.map((days, index) => {
    const start = new Date()
    start.setDate(start.getDate() + days)
    start.setHours(index % 2 === 0 ? 17 : 10, 0, 0, 0)
    const end = new Date(start.getTime() + durationMinutes * 60_000)
    return {
      id: `${experienceId}-slot-${index + 1}`,
      experience_id: experienceId,
      start_at: start.toISOString(),
      end_at: end.toISOString(),
      capacity: maxGuests,
      booked_count: 0,
      status: 'open',
    }
  })
}

export function buildExperience(
  id: string,
  slug: string,
  title: string,
  summary: string,
  description: string,
  price_amount: number,
  duration_minutes: number,
  max_guests: number,
  min_guests: number,
  city: string,
  address: string,
  latitude: number,
  longitude: number,
  category_id: string,
  category_slug: string,
  rating_avg: number,
  rating_count: number,
  host: Experience['host'],
): Experience {
  const photos = photosFor(slug, title)
  return {
    id,
    brand_id: '1',
    host_id: host!.id,
    category_id,
    category_slug,
    slug,
    title,
    summary,
    description,
    language: 'fa',
    price_amount,
    currency: 'IRR',
    duration_minutes,
    max_guests,
    min_guests,
    city,
    address,
    latitude,
    longitude,
    status: 'published',
    rating_avg,
    rating_count,
    created_at: demoNow,
    updated_at: demoNow,
    cover_photo_url: photos.find((photo) => photo.is_cover)?.url ?? photos[0].url,
    photos,
    slots: slotsFor(id, duration_minutes, max_guests),
    host,
  }
}

const e = buildExperience
const H = demoHosts

export const demoExperiencesMore: Experience[] = [
  e('19', 'shiraz-hafez-poetry-evening', 'شب شعر و فال حافظ در آرامگاه', 'شبی در آرامگاه حافظ با شعرخوانی و فال حافظ.', 'تجربه‌ای فرهنگی در شیراز.', 920000, 120, 12, 1, 'شیراز', 'آرامگاه حافظ', 29.624, 52.558, '4', 'history', 4.96, 178, H.reza),
  e('20', 'tehran-street-food-night', 'گشت شبانه غذاهای خیابانی تهران', 'از ساندویچ بندری تا آش رشته در کوچه‌های قدیمی.', 'تور غذای خیابانی تهران.', 850000, 180, 8, 2, 'تهران', 'خیابان ۳۰ تیر', 35.689, 51.422, '1', 'cooking', 4.84, 142, H.maryam),
  e('21', 'lahijan-tea-garden-tour', 'پیاده‌روی در باغات چای لاهیجان', 'از مزرعه تا فنجان؛ مسیر چای شمال.', 'بازدید از باغ چای گیلان.', 1050000, 150, 10, 2, 'لاهیجان', 'باغات چای لاهیجان', 37.207, 50.004, '3', 'nature', 4.81, 55, H.sina),
  e('22', 'hormoz-rainbow-island', 'گشت جزیره رنگارنگ هرمز', 'خاک سرخ، سواحل نمکی و جزیره‌ای از رنگ‌ها.', 'تور یک‌روزه هرمز.', 1950000, 360, 8, 2, 'هرمز', 'جزیره هرمز', 27.094, 56.452, '6', 'adventure', 4.90, 88, H.negar),
  e('23', 'qeshm-star-valley-sunset', 'غروب در دره ستارگان قشم', 'پیاده‌روی در دره‌ای با سنگ‌های کهکشانی.', 'طبیعت‌گردی دره ستارگان.', 1380000, 180, 10, 2, 'قشم', 'دره ستارگان', 26.750, 55.770, '6', 'adventure', 4.87, 73, H.sina),
  e('24', 'rasht-gilani-feast', 'پخت و چشیدن غذاهای گیلانی', 'میرza قاسمی، باقالا قاتوق و ماهی دودی.', 'کلاس آشپزی رشت.', 1320000, 210, 8, 2, 'رشت', 'محله گلسار', 37.268, 49.589, '1', 'cooking', 4.92, 104, H.sina),
  e('25', 'mashhad-breakfast-pilgrimage', 'صبحانه سنتی مشهد پس از زیارت', 'شله‌زرد، نان سنگک و چای زعفرانی در حجره قدیمی.', 'تجربه صبحانه مشهد.', 680000, 90, 10, 2, 'مشهد', 'حرم مطهر', 36.288, 59.615, '1', 'cooking', 4.78, 156, H.maryam),
  e('26', 'hamedan-alavian-tour', 'کشف همدان؛ الوند و گنجنامه', 'تور تاریخی هگمتانه و صخره‌های گنجنامه.', 'تاریخ و طبیعت همدان.', 990000, 240, 10, 1, 'همدان', 'گنجنامه', 34.798, 48.515, '4', 'history', 4.83, 62, H.reza),
  e('27', 'sanandaj-kurdish-music', 'شب موسیقی کردی در سنندج', 'اجرای زنده دهل و تمبور با چای محلی.', 'موسیقی سنتی کردستان.', 1150000, 120, 14, 2, 'سنندج', 'خانه تاریخی سنندج', 35.314, 46.998, '5', 'music', 4.88, 47, H.arash),
  e('28', 'isfahan-miniatu-workshop', 'کارگاه نقاشی مینیاتور اصفهان', 'طراحی گل و مرغ با قلم مویی.', 'هنر مینیاتور ایرانی.', 1480000, 150, 6, 1, 'اصفهان', 'میدان نقش جهان', 32.657, 51.677, '2', 'handicraft', 4.94, 91, H.arash),
  e('29', 'chabahar-mud-volcanoes', 'کوه‌های گلی و سواحل چابهار', 'سفر به یکی از شگفت‌انگیزترین سواحل ایران.', 'ماجراجویی در سیستان.', 2400000, 420, 6, 2, 'چابهار', 'کوه‌های گلی', 25.292, 60.643, '6', 'adventure', 4.85, 34, H.negar),
  e('30', 'urmia-lake-birdwatch', 'تماشای پرندگان دریاچه ارومیه', 'تور پرنده‌نگری با راهنمای محلی.', 'طبیعت‌گردی دریاچه ارومیه.', 1100000, 180, 8, 2, 'ارومیه', 'سواحل دریاچه ارومیه', 37.552, 45.076, '3', 'nature', 4.76, 41, H.fatemeh),
  e('31', 'naein-desert-citadel', 'قلعه نارین و کویر نائین', 'بازدید از قلعه ۲۰۰۰ساله و غروب کویر.', 'تاریخ و کویر نائین.', 1250000, 210, 8, 2, 'نائین', 'قلعه نارین', 32.860, 53.088, '4', 'history', 4.86, 38, H.negar),
  e('32', 'abadeh-wood-inlay', 'هنر خاتم‌کاری و منبت در آباده', 'ساخت یک جعبه خاتم کوچک.', 'صنایع‌دستی فارس.', 1190000, 150, 5, 1, 'آباده', 'بازار آباده', 31.160, 52.650, '2', 'handicraft', 4.80, 33, H.arash),
  e('33', 'amol-hyrcanian-forest', 'پیاده‌روی جنگل‌های هیرکانی آمل', 'مسیر آبشار و درختان کهنسال.', 'طبیعت‌گردی جنگل شمال.', 1420000, 300, 8, 2, 'آمل', 'جنگل‌های اطراف آمل', 36.469, 52.351, '3', 'nature', 4.89, 57, H.sina),
  e('34', 'kish-coral-reef-snorkel', 'غواصی سطحی در آب‌های مرجانی کیش', 'مشاهده مرجان و ماهی‌های رنگی.', 'تجربه دریایی کیش.', 1850000, 180, 6, 2, 'کیش', 'ساحل مرجان', 26.557, 54.019, '6', 'adventure', 4.91, 69, H.sina),
  e('35', 'ahvaz-karun-river-walk', 'غروب کنار کارون و بازار عربی', 'پیاده‌روی کنار رود و چشیدن ماهی‌شک.', 'فرهنگ و طبیعت اهواز.', 870000, 150, 10, 2, 'اهواز', 'پل سفید', 31.320, 48.669, '3', 'nature', 4.74, 48, H.reza),
  e('36', 'bojnurd-saffron-tea-ceremony', 'مراسم چای زعفران و شیرینی بجنورد', 'آموزش دم‌کردن چای زعفران خراسانی.', 'تجربه چای شمال خراسان.', 760000, 90, 8, 2, 'بجنورد', 'بازار قدیم بجنورد', 37.475, 57.329, '1', 'cooking', 4.72, 29, H.fatemeh),
  e('37', 'zanjan-soltaniyeh-dome', 'بازدید از گنبد سلطانیه', 'شاهکار معماری ایلخانی با راهنمای تاریخ.', 'تور فرهنگی زنجان.', 940000, 180, 12, 1, 'زنجان', 'گنبد سلطانیه', 36.434, 48.796, '4', 'history', 4.87, 54, H.negar),
  e('38', 'tehran-indie-music-jam', 'جَم موسیقی مستقل در کافه تهران', 'اجرای زنده بندهای مستقل در فضایی صمیمی.', 'شب موسیقی معاصر تهران.', 990000, 120, 20, 2, 'تهران', 'خیابان ایرانشهر', 35.707, 51.414, '5', 'music', 4.82, 118, H.reza),
]
