export const recordKinds = ["sightseeings", "stays", "meals", "others"] as const;
export type RecordKind = (typeof recordKinds)[number];

export const recordLabels: Record<RecordKind, string> = {
  sightseeings: "観光",
  stays: "宿泊",
  meals: "食事",
  others: "その他",
};

export const mealLabels: Record<string, string> = {
  breakfast: "朝食",
  lunch: "昼食",
  dinner: "夕食",
  snack: "軽食",
  other: "その他",
};

export type PlanRecord = {
  id: number;
  travelPlanId: number;
  name: string;
  mealType: string | null;
  locationText: string | null;
  mapUrl: string | null;
  cost: number | null;
  sortOrder: number;
};

export type RecordInput = {
  name: string;
  locationText: string | null;
  mapUrl: string | null;
  mealType?: string | null;
};

export const emptyRecordInput = (): RecordInput => ({ name: "", locationText: null, mapUrl: null });
