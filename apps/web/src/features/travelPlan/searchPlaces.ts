export type PlaceCandidate = {
  label: string; // display_name。候補リストの表示専用(保存しない)
  name: string;  // place の初期値
  country: string | null;
  state: string | null;
  city: string | null;
  latitude: number;
  longitude: number;
};

type NominatimAddress = Record<string, string | undefined>;

type NominatimResult = {
  lat: string;
  lon: string;
  name?: string;
  display_name: string;
  address?: NominatimAddress;
};

const ENDPOINT = "https://nominatim.openstreetmap.org/search";

const MIN_INTERVAL_MS = 1000;
let lastRequestAt = 0; // 直近に確保した送信時刻
const cache = new Map<string, PlaceCandidate[]>();

const sleep = (ms: number, signal?: AbortSignal) =>
  new Promise<void>((resolve, reject) => {
    const id = setTimeout(resolve, ms);
    signal?.addEventListener(
      "abort",
      () => {
        clearTimeout(id);
        reject(new DOMException("Aborted", "AbortError"));
      },
      { once: true },
    );
  });

const pick = (address: NominatimAddress, keys: string[]): string | null => {
  for (const key of keys) {
    const v = address[key];
    if (v) return v;
  }
  return null;
};

const CITY_KEYS = [
  "city",
  "town",
  "village",
  "municipality",
  "city_district",
  "suburb",
  "county",
];
const STATE_KEYS = ["state", "province", "region"];

const toCandidate = (r: NominatimResult): PlaceCandidate => {
  const address = r.address ?? {};
  return {
    label: r.display_name,
    name: r.name || r.display_name.split(",")[0].trim(),
    country: address.country ?? null,
    state: pick(address, STATE_KEYS),
    city: pick(address, CITY_KEYS),
    latitude: Number(r.lat),
    longitude: Number(r.lon),
  };
};

export const searchPlaces = async (
  query: string,
  signal?: AbortSignal,
): Promise<PlaceCandidate[]> => {
  const q = query.trim();
  if (!q) return [];

  const key = q.toLowerCase();
  const cached = cache.get(key);
  if (cached) return cached;

  // 前回の送信から 1 秒以上空ける(同時に呼ばれても順番に枠を確保する)
  const startAt = Math.max(Date.now(), lastRequestAt + MIN_INTERVAL_MS);
  lastRequestAt = startAt;
  const wait = startAt - Date.now();
  if (wait > 0) await sleep(wait, signal);

  const params = new URLSearchParams({
    q,
    format: "jsonv2",
    addressdetails: "1",
    limit: "5",
    "accept-language": "ja",
  });

  const res = await fetch(`${ENDPOINT}?${params}`, { signal });
  if (!res.ok) throw new Error(`Nominatim error: ${res.status}`);

  const data: NominatimResult[] = await res.json();
  const result = data.map(toCandidate);
  cache.set(key, result); // 成功したものだけキャッシュ
  return result;
};
