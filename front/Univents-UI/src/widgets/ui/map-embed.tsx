import { useEffect, useRef, useState, useCallback } from "react";
import type { Map as LeafletMap, Marker as LeafletMarker } from "leaflet";
import 'leaflet/dist/leaflet.css'

export interface LocationInfo {
  name: string;
  address: string;
}

interface GeoResult {
  lat: number;
  lon: number;
  displayName: string;
}

interface NominatimResult {
  lat: string
  lon: string
  display_name: string
}

interface LocationMapProps {
  location: LocationInfo;
  /** Initial Map Zoom (default: 15) */
  zoom?: number;
  /** Container height (default: "400px") */
  height?: string;
  className?: string;
}

const CACHE_KEY_PREFIX = "geo_cache_";
const CACHE_TTL_MS = 7 * 24 * 60 * 60 * 1000; // 7 days

interface CacheEntry {
  result: GeoResult;
  timestamp: number;
}

function getCached(query: string): GeoResult | null {
  try {
    const raw = localStorage.getItem(CACHE_KEY_PREFIX + btoa(query));
    if (!raw) return null;
    const entry = JSON.parse(raw) as CacheEntry;
    if (Date.now() - entry.timestamp > CACHE_TTL_MS) {
      localStorage.removeItem(CACHE_KEY_PREFIX + btoa(query));
      return null;
    }
    return entry.result;
  } catch {
    return null;
  }
}

function setCache(query: string, result: GeoResult): void {
  const entry: CacheEntry = { result, timestamp: Date.now() };
  localStorage.setItem(CACHE_KEY_PREFIX + btoa(query), JSON.stringify(entry));
}

const geoQueue: Array<() => void> = [];
let geoRunning = false;

function enqueueGeoRequest(fn: () => void) {
  geoQueue.push(fn);
  if (!geoRunning) drainQueue();
}

function drainQueue() {
  if (geoQueue.length === 0) {
    geoRunning = false;
    return;
  }
  geoRunning = true;
  const next = geoQueue.shift();
  if (!next) return;
  next();
  setTimeout(drainQueue, 1100);
}

async function geocode(location: LocationInfo): Promise<GeoResult> {
  const query = `${location.name}, ${location.address}`;
  const cached = getCached(query);
  if (cached) return cached;

  return new Promise((resolve, reject) => {
    enqueueGeoRequest(() => {
      void (async () => {
        try {
          const params = new URLSearchParams({
            q: query,
            format: "json",
            limit: "1",
            addressdetails: "0",
          });

          const res = await fetch(
            `https://nominatim.openstreetmap.org/search?${params}`,
            {
              headers: {
                "Accept-Language": "pt-BR,pt;q=0.9,en;q=0.8",
              },
            }
          );

          if (!res.ok) throw new Error(`Nominatim HTTP ${res.status}`);

          const data: NominatimResult[] = await res.json();
          if (!data.length) throw new Error("Endereço não encontrado");

          const result: GeoResult = {
            lat: parseFloat(data[0].lat),
            lon: parseFloat(data[0].lon),
            displayName: data[0].display_name,
          };

          setCache(query, result);
          resolve(result);
        } catch (err) {
          reject(err instanceof Error ? err : new Error(String(err)));
        }
      })();
    });
  });
}

type Status = "idle" | "loading" | "success" | "error";

export function LocationMap({
  location,
  zoom = 15,
  height = "400px",
  className = "",
}: LocationMapProps) {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const mapRef = useRef<LeafletMap | null>(null);
  const markerRef = useRef<LeafletMarker | null>(null);
  const [status, setStatus] = useState<Status>("idle");
  const [errorMsg, setErrorMsg] = useState("");

  const initMap = useCallback(async () => {
    if (!mapContainerRef.current) return;

    setStatus("loading");
    setErrorMsg("");

    const L = (await import("leaflet")).default;

    // @ts-expect-error _getIconUrl is not in the official type
    delete L.Icon.Default.prototype._getIconUrl;
    L.Icon.Default.mergeOptions({
      iconRetinaUrl:
        "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
      iconUrl:
        "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
      shadowUrl:
        "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
    });

    try {
      const geo = await geocode(location);

      if (mapRef.current) {
        mapRef.current.setView([geo.lat, geo.lon], zoom);
        if (markerRef.current) {
          markerRef.current.setLatLng([geo.lat, geo.lon]);
          markerRef.current
            .getPopup()
            ?.setContent(popupContent(location.name, location.address));
        }
        setStatus("success");
        return;
      }

      const map = L.map(mapContainerRef.current, {
        center: [geo.lat, geo.lon],
        zoom,
        zoomControl: true,
        scrollWheelZoom: false,
      });

      L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
        attribution:
          '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        maxZoom: 19,
      }).addTo(map);

      const marker = L.marker([geo.lat, geo.lon])
        .addTo(map)
        .bindPopup(popupContent(location.name, location.address))
        .openPopup();

      mapRef.current = map;
      markerRef.current = marker;
      setStatus("success");
    } catch (err) {
      setStatus("error");
      setErrorMsg(err instanceof Error ? err.message : "Erro desconhecido");
    }
  }, [location, zoom]);

  useEffect(() => {
    void initMap();

    return () => {
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
        markerRef.current = null;
      }
    };
  }, [initMap]);

  return (
    <div className={`location-map-wrapper ${className}`} style={{ position: "relative" }}>
      <div
        ref={mapContainerRef}
        style={{
          height,
          width: "100%",
          borderRadius: "12px",
          overflow: "hidden",
          background: "#e8e0d8",
        }}
      />

      {status === "loading" && (
        <div style={overlayStyle}>
          <div style={spinnerStyle} />
          <span style={{ marginTop: 12, fontSize: 14, color: "#555" }}>
            Carregando mapa…
          </span>
        </div>
      )}

      {status === "error" && (
        <div style={overlayStyle}>
          <span style={{ fontSize: 32 }}>📍</span>
          <p style={{ margin: "8px 0 4px", fontWeight: 600, color: "#333" }}>
            {location.name}
          </p>
          <p style={{ margin: 0, fontSize: 13, color: "#666", textAlign: "center" }}>
            {location.address}
          </p>
          <p style={{ margin: "12px 0 0", fontSize: 12, color: "#e55" }}>
            {errorMsg}
          </p>
          <button
            onClick={() => void initMap()}
            style={retryButtonStyle}
          >
            Tentar novamente
          </button>
        </div>
      )}
    </div>
  );
}

function popupContent(name: string, address: string) {
  return `
    <div style="min-width:160px; font-family: sans-serif;">
      <strong style="display:block; margin-bottom:4px;">${name}</strong>
      <span style="font-size:12px; color:#555;">${address}</span>
    </div>
  `;
}

const overlayStyle: React.CSSProperties = {
  position: "absolute",
  inset: 0,
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  justifyContent: "center",
  background: "rgba(255,255,255,0.88)",
  borderRadius: "12px",
  zIndex: 1000,
  padding: "24px",
};

const spinnerStyle: React.CSSProperties = {
  width: 36,
  height: 36,
  border: "3px solid #ddd",
  borderTopColor: "#555",
  borderRadius: "50%",
  animation: "spin 0.8s linear infinite",
};

const retryButtonStyle: React.CSSProperties = {
  marginTop: 12,
  padding: "6px 16px",
  border: "1px solid #ccc",
  borderRadius: 6,
  background: "#fff",
  cursor: "pointer",
  fontSize: 13,
};


interface MultiLocationMapProps {
  locations: LocationInfo[];
  height?: string;
  className?: string;
}

export function MultiLocationMap({
  locations,
  height = "500px",
  className = "",
}: MultiLocationMapProps) {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const mapRef = useRef<LeafletMap | null>(null);
  const [status, setStatus] = useState<Status>("idle");
  const [errorMsg, setErrorMsg] = useState("");
  const cancelledRef = useRef(false);

  useEffect(() => {
    if (!mapContainerRef.current || !locations.length) return;

    cancelledRef.current = false;
    setStatus("loading");

    void (async () => {
      const L = (await import("leaflet")).default;

      // @ts-expect-error _getIconUrl is not in the official type
      delete L.Icon.Default.prototype._getIconUrl;
      L.Icon.Default.mergeOptions({
        iconRetinaUrl:
          "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
        iconUrl:
          "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
        shadowUrl:
          "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
      });

      try {
        const results = await Promise.all(locations.map(geocode));
        if (cancelledRef.current) return;

        if (mapRef.current) {
          mapRef.current.remove();
          mapRef.current = null;
        }

        const bounds: [number, number][] = results.map((r) => [r.lat, r.lon]);

        if (!mapContainerRef.current) return;
        const map = L.map(mapContainerRef.current, {
          scrollWheelZoom: false,
        });

        L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
          attribution:
            '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
          maxZoom: 19,
        }).addTo(map);

        results.forEach((geo, i) => {
          L.marker([geo.lat, geo.lon])
            .addTo(map)
            .bindPopup(popupContent(locations[i].name, locations[i].address));
        });

        map.fitBounds(bounds, { padding: [40, 40] });
        mapRef.current = map;
        setStatus("success");
      } catch (err) {
        if (!cancelledRef.current) {
          setStatus("error");
          setErrorMsg(err instanceof Error ? err.message : "Erro desconhecido");
        }
      }
    })();

    return () => {
      cancelledRef.current = true;
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
      }
    };
  }, [locations]);

  return (
    <div className={`location-map-wrapper ${className}`} style={{ position: "relative" }}>
      <div
        ref={mapContainerRef}
        style={{ height, width: "100%", borderRadius: "12px", background: "#e8e0d8" }}
      />
      {status === "loading" && (
        <div style={overlayStyle}>
          <div style={spinnerStyle} />
          <span style={{ marginTop: 12, fontSize: 14, color: "#555" }}>
            Geocodificando {locations.length} endereço{locations.length > 1 ? "s" : ""}…
          </span>
        </div>
      )}
      {status === "error" && (
        <div style={overlayStyle}>
          <span style={{ color: "#e55", fontSize: 13 }}>{errorMsg}</span>
        </div>
      )}
    </div>
  );
}