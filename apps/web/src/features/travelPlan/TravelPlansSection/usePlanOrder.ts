import { useRef } from "react";

import { useAuth } from "../../../app/providers/useAuth";
import { useFlash } from "../../../app/providers/useFlash";
import { reorderTravelPlans } from "../repository";
import { sortPlans } from "../sortPlans";
import type { TravelPlan } from "../types";

const movePlan = (plans: TravelPlan[], plan: TravelPlan, target: number): TravelPlan[] => {
  const next = [...plans];
  next.splice(
    next.findIndex((item) => item.id === plan.id),
    1,
  );
  next.splice(target, 0, plan);
  return next.map((item, index) => ({ ...item, sortOrder: index + 1 }));
};

const focusMovedPlan = (plan: TravelPlan) => {
  requestAnimationFrame(() => document.getElementById(`plan-card-${plan.id}`)?.focus());
};

const showMovedPlan = ({
  plan,
  target,
  next,
  setPlans,
  showFlash,
}: {
  plan: TravelPlan;
  target: number;
  next: TravelPlan[];
  setPlans: (plans: TravelPlan[]) => void;
  showFlash: (message: string) => void;
}) => {
  setPlans(next);
  showFlash(`${plan.place}を${target + 1}番目に移動しました`);
  focusMovedPlan(plan);
};

export const usePlanOrder = (
  travelId: number,
  plans: TravelPlan[],
  setPlans: (plans: TravelPlan[]) => void,
) => {
  const { token } = useAuth();
  const { showFlash } = useFlash();
  const busy = useRef(false);
  const move = async (plan: TravelPlan, target: number) => {
    if (!token || busy.current) return;
    busy.current = true;
    const before = sortPlans(plans, "planning");
    const next = movePlan(before, plan, target);
    showMovedPlan({ plan, target, next, setPlans, showFlash });
    try {
      await reorderTravelPlans(
        token,
        travelId,
        next.map((item) => item.id),
      );
    } catch {
      setPlans(before);
      showFlash("並べ替えに失敗しました。元の順序に戻しました", "error");
    } finally {
      busy.current = false;
    }
  };
  return { move };
};
