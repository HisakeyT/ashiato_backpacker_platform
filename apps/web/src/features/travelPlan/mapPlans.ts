import type { TravelPlan } from "./types";

export type MapPlan = { plan: TravelPlan; number: number; position: [number, number] };

export const mapPlans = (plans: TravelPlan[]) => {
  const ordered = [...plans].sort((a, b) => a.sortOrder - b.sortOrder || a.id - b.id);
  const located: MapPlan[] = [];
  const missing: TravelPlan[] = [];
  ordered.forEach((plan, index) => {
    const { latitude, longitude } = plan;
    if (
      latitude == null ||
      longitude == null ||
      !Number.isFinite(latitude) ||
      !Number.isFinite(longitude) ||
      Math.abs(latitude) > 90 ||
      Math.abs(longitude) > 180
    ) {
      missing.push(plan);
    } else {
      located.push({ plan, number: index + 1, position: [latitude, longitude] });
    }
  });
  unwrapLongitudes(located);
  return { located, missing };
};

export const groupMapPlans = (located: MapPlan[]) => {
  const groups = new Map<string, MapPlan[]>();
  for (const item of located) {
    const key = item.position.join(",");
    groups.set(key, [...(groups.get(key) ?? []), item]);
  }
  return [...groups.values()];
};

const unwrapLongitudes = (located: MapPlan[]) => {
  // Unwrap around the largest empty longitude gap, keeping dateline crossings compact.
  const longitudes = [...new Set(located.map(({ position }) => (position[1] + 360) % 360))].sort(
    (a, b) => a - b,
  );
  let start = longitudes[0] ?? 0;
  let largestGap = -1;
  longitudes.forEach((longitude, index) => {
    const next = longitudes[(index + 1) % longitudes.length];
    const gap = (index === longitudes.length - 1 ? next + 360 : next) - longitude;
    if (gap > largestGap) {
      largestGap = gap;
      start = next;
    }
  });
  for (const item of located) {
    let longitude = (item.position[1] + 360) % 360;
    if (longitude < start) longitude += 360;
    item.position[1] = longitude;
  }
};
