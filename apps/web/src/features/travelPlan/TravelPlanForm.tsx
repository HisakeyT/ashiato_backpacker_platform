import { useState, type SubmitEvent } from "react";
import { Alert, Box, Button, Stack, TextField } from "@mui/material";
import { FormField } from "../../components/FormField";

export type TravelPlanFormValue = {
  startDate: string;
  place: string;
  content: string;
};

type TravelPlanFormProps = {
  initialValue?: TravelPlanFormValue;
  defaultDate?: string;
  minDate?: string;
  maxDate?: string;
  submitLabel: string;
  error: string | null;
  isSubmitting: boolean;
  onCancel: () => void;
  onSubmit: (value: TravelPlanFormValue) => Promise<void>;
};

export const TravelPlanForm = ({
  initialValue,
  defaultDate = "",
  minDate,
  maxDate,
  submitLabel,
  error,
  isSubmitting,
  onCancel,
  onSubmit,
}: TravelPlanFormProps) => {
  const [startDate, setStartDate] = useState(initialValue?.startDate ?? defaultDate);
  const [place, setPlace] = useState(initialValue?.place ?? "");
  const [content, setContent] = useState(initialValue?.content ?? "");

  const handleSubmit = async (e: SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (isSubmitting) return;
    await onSubmit({ startDate, place, content });
  };

  return (
    <Box className="TravelPlanForm" component="form" onSubmit={handleSubmit}>
      <Stack spacing={3}>
        {error && <Alert severity="error">{error}</Alert>}

        <FormField label="日付" htmlFor="plan-date" required>
          <TextField
            id="plan-date"
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
            required
            fullWidth
            slotProps={{ htmlInput: { min: minDate, max: maxDate } }}
          />
        </FormField>

        <FormField label="場所" htmlFor="plan-place" required>
          <TextField
            id="plan-place"
            value={place}
            onChange={(e) => setPlace(e.target.value)}
            placeholder="例：アンコールワット"
            required
            fullWidth
          />
        </FormField>

        <FormField label="内容" htmlFor="plan-content" required>
          <TextField
            id="plan-content"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="例：朝日を見てから遺跡を回る"
            required
            fullWidth
            multiline
            minRows={3}
          />
        </FormField>

        <Box sx={{ display: "flex", justifyContent: "flex-end", gap: 1 }}>
          <Button onClick={onCancel} color="inherit" disabled={isSubmitting}>
            キャンセル
          </Button>
          <Button type="submit" variant="contained" disabled={isSubmitting}>
            {isSubmitting ? "送信中..." : submitLabel}
          </Button>
        </Box>
      </Stack>
    </Box>
  );
};
