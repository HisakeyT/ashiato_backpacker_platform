import { useState } from "react";
import { useAuth } from "../../app/providers/useAuth";
import { BaseDialog } from "../../components/BaseDialog";
import { createTravelPlan, updateTravelPlan } from "./repository";
import { TravelPlanForm, type TravelPlanFormValue } from "./TravelPlanForm";
import type { TravelPlan } from "./types";

type TravelPlanDialogProps = {
  open: boolean;
  travelId: number;
  startDate: string;
  endDate: string;
  plan?: TravelPlan | null;
  onClose: () => void;
  onSaved: (result: "created" | "updated") => void;
};

export const TravelPlanDialog = ({
  open,
  travelId,
  startDate,
  endDate,
  plan,
  onClose,
  onSaved,
}: TravelPlanDialogProps) => {
  const { token } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const isEdit = plan != null;

  const handleClose = () => {
    if (isSubmitting) return;
    setError(null);
    onClose();
  };

  const handleSubmit = async (value: TravelPlanFormValue) => {
    if (!token) return;
    setIsSubmitting(true);
    setError(null);
    try {
      const { endDate: formEndDate, ...rest } = value;
      const periodInput = formEndDate ? { endDate: formEndDate } : {};

      if (plan) {
        await updateTravelPlan(token, travelId, plan.id, { ...rest, ...periodInput });
      } else {
        await createTravelPlan(token, travelId, { ...rest, ...periodInput, sortOrder: 0 });
      }
      onSaved(isEdit ? "updated" : "created");
      onClose();
    } catch {
      setError(isEdit ? "プランの更新に失敗しました" : "プランの追加に失敗しました");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <BaseDialog
      className="TravelPlanDialog"
      open={open}
      onClose={handleClose}
      title={isEdit ? "プランを編集" : "プランを追加"}
      busy={isSubmitting}
      maxWidth="sm"
    >
      <TravelPlanForm
        initialValue={
          plan
            ? {
              startDate: plan.startDate,
              endDate: plan.endDate ?? "",
              place: plan.place,
              content: plan.content
            }
            : undefined
        }
        defaultDate={startDate}
        minDate={startDate}
        maxDate={endDate}
        submitLabel={isEdit ? "保存する" : "追加する"}
        error={error}
        isSubmitting={isSubmitting}
        onCancel={handleClose}
        onSubmit={handleSubmit}
      />
    </BaseDialog>
  );
};
