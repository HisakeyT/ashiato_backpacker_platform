import { useState } from "react";

import { useAuth } from "../../../app/providers/useAuth";
import { useFlash } from "../../../app/providers/useFlash";
import { deleteTravelPlan } from "../repository";
import type { TravelPlan } from "../types";

export const usePlanDeletion = (travelId: number, refreshPlans: () => Promise<void>) => {
  const { token } = useAuth();
  const { showFlash } = useFlash();
  const [deletingPlan, setDeletingPlan] = useState<TravelPlan | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const handleClose = () => {
    if (isDeleting) return;
    setError(null);
    setDeletingPlan(null);
  };
  const handleConfirm = async () => {
    if (!token || !deletingPlan) return;
    setIsDeleting(true);
    setError(null);
    try {
      await deleteTravelPlan(token, travelId, deletingPlan.id);
      setDeletingPlan(null);
      void refreshPlans();
      showFlash("プランを削除しました");
    } catch {
      setError("プランの削除に失敗しました");
    } finally {
      setIsDeleting(false);
    }
  };
  return { deletingPlan, setDeletingPlan, isDeleting, error, handleClose, handleConfirm };
};
