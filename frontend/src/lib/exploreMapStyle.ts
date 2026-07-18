import type { StyleSpecification } from 'maplibre-gl'

// Raster basemap avoids OpenFreeMap vector-tile null-property parse errors in MapLibre.
export const exploreMapStyle: StyleSpecification = {
  version: 8,
  sources: {
    carto: {
      type: 'raster',
      tiles: [
        'https://a.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}@2x.png',
        'https://b.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}@2x.png',
        'https://c.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}@2x.png',
      ],
      tileSize: 256,
      attribution: '&copy; OpenStreetMap &copy; CARTO',
    },
  },
  layers: [
    {
      id: 'carto',
      type: 'raster',
      source: 'carto',
      minzoom: 0,
      maxzoom: 19,
    },
  ],
}
