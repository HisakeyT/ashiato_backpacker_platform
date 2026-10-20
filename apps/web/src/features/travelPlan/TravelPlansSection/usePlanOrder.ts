import { useAuth } from "../../../app/providers/useAuth";
import { useFlash } from "../../../app/providers/useFlash";
import { reorderTravelPlans } from "../repository";
import type { TravelPlan } from "../types";
import { persistChangedPlanOrder } from "./planOrderChanges";

export const usePlanOrder = (
  travelId: number,
  plans: TravelPlan[],
  setPlans: (plans: TravelPlan[]) => void,
) => {
  const { token } = useAuth();
  const { showFlash } = useFlash();

  const save = async (ordered: TravelPlan[]) => {
    const changed = await persistChangedPlanOrder(plans, ordered, async (ids) => {
      if (!token) throw new Error("ログインが必要です");
      await reorderTravelPlans(token, travelId, ids);
    });
    if (!changed) return;
    setPlans(ordered.map((plan, index) => ({ ...plan, sortOrder: index + 1 })));
    showFlash("プランの順序を更新しました");
  };

  return { save };
};
