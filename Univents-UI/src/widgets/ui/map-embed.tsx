import { useEffect, useState } from 'react'
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet'
import { Navigation, Map as MapIcon } from 'lucide-react'
import 'leaflet/dist/leaflet.css'

import L from 'leaflet'
import icon from 'leaflet/dist/images/marker-icon.png'
import iconShadow from 'leaflet/dist/images/marker-shadow.png'

const DefaultIcon = L.icon({
  iconUrl: icon,
  shadowUrl: iconShadow,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
})

L.Marker.prototype.options.icon = DefaultIcon

interface NominatimResult {
  lat: string
  lon: string
  display_name: string
}

interface MapEmbedProps {
  name: string
  address: string
}

type MapStyle = 'standard' | 'satellite' | 'terrain'

const tileLayers: Record<MapStyle, { url: string; name: string }> = {
  standard: {
    url: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
    name: 'Padrão',
  },
  satellite: {
    url: 'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}',
    name: 'Satélite',
  },
  terrain: {
    url: 'https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png',
    name: 'Terreno',
  },
}

export default function MapEmbed({ name, address }: MapEmbedProps) {
  const [coords, setCoords] = useState<[number, number] | null>(null)
  const [displayName, setDisplayName] = useState('')
  const [style, setStyle] = useState<MapStyle>('standard')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)

  const fullAddress = `${name}, ${address}`

  useEffect(() => {
    setLoading(true)
    fetch(
      `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(
        fullAddress
      )}&limit=1`
    )
      .then(async (res) => {
        const data: NominatimResult[] = await res.json()

        if (!Array.isArray(data) || !data[0]) {
          throw new Error('Invalid geocode response')
        }

        return data
      })
      .then((data) => {
        setCoords([Number(data[0].lat), Number(data[0].lon)])
        setDisplayName(data[0].display_name)
      })
      .catch(() => {
        setError(true)
        setLoading(false)
      })
  }, [fullAddress])

  const openInGoogleMaps = () => {
    if (!coords) return
    const url = `https://www.google.com/maps/dir/?api=1&destination=${coords[0]},${coords[1]}`
    window.open(url, '_blank')
  }

  const openInWaze = () => {
    if (!coords) return
    const url = `https://waze.com/ul?ll=${coords[0]},${coords[1]}&navigate=yes`
    window.open(url, '_blank')
  }

  if (loading) {
    return (
      <div className="h-48 rounded-xl bg-muted animate-pulse flex items-center justify-center">
        <span className="text-xs text-muted-foreground">Carregando mapa...</span>
      </div>
    )
  }

  if (error || !coords) {
    return (
      <div className="h-48 rounded-xl bg-muted border border-dashed border-border flex flex-col items-center justify-center gap-2 p-4 text-center">
        <MapIcon className="w-8 h-8 text-muted-foreground/50" />
        <p className="text-xs text-muted-foreground">Não foi possível carregar o mapa</p>
        <a
          href={`https://www.google.com/maps/search/${encodeURIComponent(fullAddress)}`}
          target="_blank"
          rel="noopener noreferrer"
          className="text-xs text-primary hover:underline"
        >
          Buscar no Google Maps
        </a>
      </div>
    )
  }

  return (
    <div className="space-y-3">
      {/* Map */}
      <div className="relative">
        <MapContainer
          center={coords}
          zoom={16}
          scrollWheelZoom={false}
          className="h-48 rounded-xl z-0"
        >
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            url={tileLayers[style].url}
          />
          <Marker position={coords}>
            <Popup>
              <div className="text-sm">
                <p className="font-semibold">{name}</p>
                <p className="text-muted-foreground">{address}</p>
              </div>
            </Popup>
          </Marker>
        </MapContainer>

        {/* Layers */}
        <div className="absolute top-2 right-2 z-400">
          <div className="bg-background/90 backdrop-blur-sm rounded-lg border border-border shadow-sm p-1 flex gap-1">
            {(Object.keys(tileLayers) as MapStyle[]).map((s) => (
              <button
                key={s}
                onClick={() => { setStyle(s) }}
                className={`px-2 py-1 text-[10px] font-medium rounded-md transition-colors ${style === s
                  ? 'bg-primary text-primary-foreground'
                  : 'hover:bg-muted text-muted-foreground'
                  }`}
              >
                {tileLayers[s].name}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Info and Actions */}
      <div className="space-y-2">
        {displayName && (
          <p className="text-xs text-muted-foreground line-clamp-2">{displayName}</p>
        )}

        <div className="flex gap-2">
          <button
            onClick={openInGoogleMaps}
            className="flex-1 flex items-center justify-center gap-1.5 py-2.5 rounded-xl bg-primary text-primary-foreground text-xs font-semibold hover:brightness-110 active:scale-[0.98] transition-all"
          >
            <Navigation className="w-3.5 h-3.5" />
            Google Maps
          </button>
          <button
            onClick={openInWaze}
            className="flex-1 flex items-center justify-center gap-1.5 py-2.5 rounded-xl bg-blue-500/10 text-blue-600 border border-blue-500/20 text-xs font-semibold hover:bg-blue-500/15 active:scale-[0.98] transition-all"
          >
            <Navigation className="w-3.5 h-3.5" />
            Waze
          </button>
        </div>
      </div>
    </div>
  )
}