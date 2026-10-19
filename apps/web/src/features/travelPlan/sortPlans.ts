import type { TravelPlan } from "./types";

export const sortPlans = (plans: TravelPlan[]): TravelPlan[] => {
  return [...plans].sort((a, b) => {
    return a.startDate < b.startDate ? -1 : 1;
  });
};
