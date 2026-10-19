import type { TravelPlan } from "./types";

export const sortPlans = (
  plans: TravelPlan[],
  mode: "planning" | "recorded" = "recorded",
): TravelPlan[] => {
  return [...plans].sort((a, b) => {
    if (mode === "planning") return a.sortOrder - b.sortOrder || a.id - b.id;
    return (
      (a.startDate ?? "").localeCompare(b.startDate ?? "") ||
      (a.endDate ?? "").localeCompare(b.endDate ?? "") ||
      a.sortOrder - b.sortOrder ||
      a.id - b.id
    );
  });
};
