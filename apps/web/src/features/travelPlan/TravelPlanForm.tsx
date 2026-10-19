import { Alert, Box, Stack } from "@mui/material";
import type { SubmitEvent } from "react";

import {
  PlanActions,
  PlanContent,
  PlanDates,
  PlanLocation,
  PlanPlace,
} from "./TravelPlanFormFields";
import { type PlanDetails, usePlaceSearch, usePlanDetails } from "./useTravelPlanForm";

export type TravelPlanFormValue = {
  startDate: string;
  endDate: string;
  place: string;
  country: string;
  state: string;
  city: string;
  latitude: number | null;
  longitude: number | null;
  content: string;
};

type TravelPlanFormProps = {
  initialValue?: TravelPlanFormValue;
  defaultDate?: string;
  minDate?: string;
  maxDate?: string;
  mode: "planning" | "recorded";
  submitLabel: string;
  error: string | null;
  isSubmitting: boolean;
  onCancel: () => void;
  onSubmit: (value: TravelPlanFormValue) => Promise<void>;
};

const toFormValue = (details: PlanDetails): TravelPlanFormValue => ({
  startDate: details.startDate,
  endDate: details.endDate,
  place: details.place,
  country: details.country,
  state: details.state,
  city: details.city,
  latitude: details.latitude,
  longitude: details.longitude,
  content: details.content,
});

export const TravelPlanForm = ({
  initialValue,
  defaultDate,
  minDate,
  maxDate,
  mode,
  submitLabel,
  error,
  isSubmitting,
  onCancel,
  onSubmit,
}: TravelPlanFormProps) => {
  const details = usePlanDetails(initialValue, defaultDate);
  const search = usePlaceSearch(details);
  const handleSubmit = async (event: SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (isSubmitting) return;
    await onSubmit(toFormValue(details));
  };
  return (
    <Box className="TravelPlanForm" component="form" onSubmit={handleSubmit}>
      <Stack spacing={3}>
        {error && <Alert severity="error">{error}</Alert>}
        <PlanDates details={details} minDate={minDate} maxDate={maxDate} mode={mode} />
        <PlanPlace details={details} search={search} />
        <PlanLocation details={details} />
        <PlanContent details={details} />
        <PlanActions onCancel={onCancel} isSubmitting={isSubmitting} submitLabel={submitLabel} />
      </Stack>
    </Box>
  );
};
