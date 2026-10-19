import { useEffect, useState } from "react";

import { useAuth } from "../../app/providers/useAuth";
import { getTravelPlans } from "./repository";
import type { TravelPlan } from "./types";

export const usePlanData = (travelId: number) => {
  const { token } = useAuth();
  const [plans, setPlans] = useState<TravelPlan[]>([]);
  const [plansError, setPlansError] = useState<string | null>(null);
  useEffect(() => {
    if (!token) return;
    let cancelled = false;
    void getTravelPlans(token, travelId).then(
      (nextPlans) => {
        if (cancelled) return;
        setPlans(nextPlans);
        setPlansError(null);
      },
      () => {
        if (!cancelled) setPlansError("予定の取得に失敗しました");
      },
    );
    return () => {
      cancelled = true;
    };
  }, [token, travelId]);
  const refreshPlans = async () => {
    if (!token) return;
    try {
      setPlans(await getTravelPlans(token, travelId));
      setPlansError(null);
    } catch {
      setPlansError("予定の取得に失敗しました");
    }
  };
  return { plans, setPlans, plansError, refreshPlans };
};
