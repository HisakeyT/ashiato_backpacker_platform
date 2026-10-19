import type { TravelPlan } from "../types";

export type MapPlan = { plan: TravelPlan; number: number; position: [number, number] };

const normalizeLongitude = (longitude: number): number => (longitude + 360) % 360;

/** 地点間で最も広い経度の空白を、地図の表示範囲の切れ目に選ぶ。 */
const findRangeStart = (longitudes: number[]): number => {
  let start = longitudes[0] ?? 0;
  let largestGap = -1;

  longitudes.forEach((longitude, index) => {
    // 最後から最初へ戻る区間は、最初の経度に 360 度を足して比較する。
    const next = index === longitudes.length - 1 ? longitudes[0] + 360 : longitudes[index + 1];
    const gap = next - longitude;
    if (gap > largestGap) {
      largestGap = gap;
      start = next;
    }
  });

  return start;
};

/** 経度を一続きの範囲に並べ直し、日付変更線をまたぐ地点を地図上でも近くに表示する。 */
const unwrapLongitudes = (located: MapPlan[]) => {
  // 例: 東経179度と西経179度なら、その間を 358 度ではなく 2 度として扱える。
  const longitudes = [
    ...new Set(located.map(({ position }) => normalizeLongitude(position[1]))),
  ].sort((a, b) => a - b);
  const start = findRangeStart(longitudes);

  // 切れ目より西側の地点を次の 360 度側へ移す。変更するのは MapPlan の座標だけ。
  for (const item of located) {
    let longitude = normalizeLongitude(item.position[1]);
    if (longitude < start) longitude += 360;
    item.position[1] = longitude;
  }
};

/** null、非有限値、緯度経度の範囲外を除き、Leaflet に渡す座標を作る。 */
const toPosition = (plan: TravelPlan): [number, number] | null => {
  const { latitude, longitude } = plan;
  if (
    latitude == null ||
    longitude == null ||
    !Number.isFinite(latitude) ||
    !Number.isFinite(longitude) ||
    Math.abs(latitude) > 90 ||
    Math.abs(longitude) > 180
  ) {
    return null;
  }
  return [latitude, longitude];
};

/** 訪問順に番号を付け、座標が有効なプランと地図に載せられないプランに分ける。 */
export const mapPlans = (plans: TravelPlan[]) => {
  const ordered = [...plans].sort((a, b) => a.sortOrder - b.sortOrder || a.id - b.id);
  const located: MapPlan[] = [];
  const missing: TravelPlan[] = [];

  ordered.forEach((plan, index) => {
    // 座標のないプランも訪問番号には含めるため、番号は元の一覧での位置から決める。
    const position = toPosition(plan);
    if (position) located.push({ plan, number: index + 1, position });
    else missing.push(plan);
  });

  unwrapLongitudes(located);
  return { located, missing };
};

/** 同じ座標のプランを一つのピンにまとめる。グループ内は訪問順を保つ。 */
export const groupMapPlans = (located: MapPlan[]) => {
  const groups = new Map<string, MapPlan[]>();
  for (const item of located) {
    const key = item.position.join(",");
    const group = groups.get(key);
    if (group) group.push(item);
    else groups.set(key, [item]);
  }
  return [...groups.values()];
};
