import type { Photo } from './types'

export function demoPhotosFor(slug: string, title: string): Photo[] {
  return Array.from({ length: 5 }, (_, index) => ({
    id: `${slug}-${index + 1}`,
    url: `https://picsum.photos/seed/${slug}-${index + 1}/1200/800`,
    alt: title,
    sort_order: index,
    is_cover: index === 0,
  }))
}
