import { useState } from "react";

import { BaseDialog } from "../../components/BaseDialog";

import { useAuth } from "../../app/providers/useAuth";
import { createTravelPlan, updateTravelPlan } from "./repository";
import { TravelPlanForm, type TravelPlanFormValue } from "./TravelPlanForm/TravelPlanForm";
import type { TravelPlan } from "./types";

type TravelPlanDialogProps = {
  open: boolean;
  travelId: number;
  defaultDate?: string;
  startDate: string;
  endDate: string;
  mode: "planning" | "recorded";
  plan?: TravelPlan | null;
  onClose: () => void;
  onSaved: (result: "created" | "updated") => void;
};

const toPlanInput = (value: TravelPlanFormValue) => {
  const { startDate, endDate, country, state, city, latitude, longitude, ...rest } = value;
  return {
    ...rest,
    startDate: startDate || null,
    endDate: endDate || null,
    ...(country ? { country } : {}),
    ...(state ? { state } : {}),
    ...(city ? { city } : {}),
    ...(latitude != null && longitude != null ? { latitude, longitude } : {}),
  };
};

const initialFormValue = (plan?: TravelPlan | null): TravelPlanFormValue | undefined =>
  plan
    ? {
        startDate: plan.startDate ?? "",
        endDate: plan.endDate ?? "",
        place: plan.place,
        country: plan.country ?? "",
        state: plan.state ?? "",
        city: plan.city ?? "",
        latitude: plan.latitude,
        longitude: plan.longitude,
        content: plan.content,
      }
    : undefined;

type DialogViewProps = TravelPlanDialogProps & {
  error: string | null;
  isSubmitting: boolean;
  onSubmit: (value: TravelPlanFormValue) => Promise<void>;
};

const DialogView = ({
  open,
  plan,
  defaultDate,
  startDate,
  endDate,
  mode,
  error,
  isSubmitting,
  onClose,
  onSubmit,
}: DialogViewProps) => (
  <BaseDialog
    className="TravelPlanDialog"
    open={open}
    onClose={onClose}
    title={plan ? "プランを編集" : "プランを追加"}
    busy={isSubmitting}
    maxWidth="sm"
  >
    <TravelPlanForm
      initialValue={initialFormValue(plan)}
      defaultDate={defaultDate ?? (mode === "recorded" ? startDate : "")}
      minDate={startDate}
      maxDate={endDate}
      mode={mode}
      submitLabel={plan ? "保存する" : "追加する"}
      error={error}
      isSubmitting={isSubmitting}
      onCancel={onClose}
      onSubmit={onSubmit}
    />
  </BaseDialog>
);

export const TravelPlanDialog = (props: TravelPlanDialogProps) => {
  const { travelId, plan, onClose, onSaved } = props;
  const { token } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
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
      if (plan) await updateTravelPlan(token, { travelId, planId: plan.id }, toPlanInput(value));
      else await createTravelPlan(token, travelId, { ...toPlanInput(value), sortOrder: 0 });
      onSaved(plan ? "updated" : "created");
      onClose();
    } catch {
      setError(plan ? "プランの更新に失敗しました" : "プランの追加に失敗しました");
    } finally {
      setIsSubmitting(false);
    }
  };
  return (
    <DialogView
      {...props}
      onClose={handleClose}
      error={error}
      isSubmitting={isSubmitting}
      onSubmit={handleSubmit}
    />
  );
};
