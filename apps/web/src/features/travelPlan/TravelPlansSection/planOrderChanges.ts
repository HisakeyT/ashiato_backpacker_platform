import { sortPlans } from "../sortPlans.ts";
import type { TravelPlan } from "../types";

export const movePlan = (plans: TravelPlan[], plan: TravelPlan, target: number): TravelPlan[] => {
  const next = [...plans];
  next.splice(
    next.findIndex((item) => item.id === plan.id),
    1,
  );
  next.splice(target, 0, plan);
  return next;
};

export const persistChangedPlanOrder = async (
  current: TravelPlan[],
  ordered: TravelPlan[],
  persist: (planIds: number[]) => Promise<void>,
): Promise<boolean> => {
  const currentIds = sortPlans(current, "planning").map((plan) => plan.id);
  const orderedIds = ordered.map((plan) => plan.id);
  if (
    orderedIds.length === currentIds.length &&
    orderedIds.every((id, index) => id === currentIds[index])
  ) {
    return false;
  }

  await persist(orderedIds);
  return true;
};
