import type { TravelPlan } from "./types";

/**
 * 計画中は仮の日付に左右されず、手動の訪問順で表示する。
 * 記録済みは開始日・終了日順。同じ日付なら手動順とIDで並びを固定する。
 * 元の配列は変更しない。
 */
export const sortPlans = (
  plans: TravelPlan[],
  mode: "planning" | "recorded" = "recorded",
): TravelPlan[] => {
  return [...plans].sort((a, b) => {
    if (mode === "planning") return a.sortOrder - b.sortOrder || a.id - b.id;
    return (
      (a.startDate ?? "").localeCompare(b.startDate ?? "") ||
      // 終了日がない単日のプランは、同じ開始日で終了日を指定したプランより先に置く。
      (a.endDate ?? "").localeCompare(b.endDate ?? "") ||
      a.sortOrder - b.sortOrder ||
      a.id - b.id
    );
  });
};
