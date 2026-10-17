import { useState, type SubmitEvent } from "react";
import {
  Alert,
  Box,
  Button,
  FormControlLabel,
  Stack,
  Switch,
  TextField,
} from "@mui/material";
import type { CreateTravelInput } from "./types";
import { FormField } from "../../components/FormField";

type TravelFormProps = {
  initialValue?: CreateTravelInput;
  submitLabel: string;
  error: string | null;
  isSubmitting: boolean;
  onSubmit: (input: CreateTravelInput) => Promise<void>;
};

const emptyValue: CreateTravelInput = {
  title: "",
  startDate: "",
  endDate: "",
  isPublic: false,
};

export const TravelForm = ({
  initialValue = emptyValue,
  submitLabel,
  error,
  isSubmitting,
  onSubmit,
}: TravelFormProps) => {
  const [title, setTitle] = useState(initialValue.title);
  const [startDate, setStartDate] = useState(initialValue.startDate);
  const [endDate, setEndDate] = useState(initialValue.endDate);
  const [isPublic, setIsPublic] = useState(initialValue.isPublic);

  const handleSubmit = async (e: SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (isSubmitting) return;
    await onSubmit({ title, startDate, endDate, isPublic });
  };

  return (
    <Box
      className="TravelForm"
      component="form"
      onSubmit={handleSubmit}
      sx={{ maxWidth: 560 }}
    >
      <Stack spacing={3}>
        {error && <Alert severity="error">{error}</Alert>}

        <FormField label="タイトル" htmlFor="travel-title" required>
          <TextField
            id="travel-title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="例：東南アジア3週間の旅"
            required
            fullWidth
          />
        </FormField>

        <Stack direction={{ xs: "column", sm: "row" }} spacing={2}>
          <FormField label="開始日" htmlFor="travel-start-date" required>
            <TextField
              id="travel-start-date"
              type="date"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
              required
              fullWidth
            />
          </FormField>
          <FormField label="終了日" htmlFor="travel-end-date" required>
            <TextField
              id="travel-end-date"
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
              required
              fullWidth
              slotProps={{ htmlInput: { min: startDate || undefined } }}
            />
          </FormField>
        </Stack>

        <FormField label="公開設定" htmlFor="travel-is-public">
          <FormControlLabel
            control={
              <Switch
                id="travel-is-public"
                checked={isPublic}
                onChange={(e) => setIsPublic(e.target.checked)}
              />
            }
            label={isPublic ? "公開（誰でも見られます）" : "非公開（自分だけが見られます）"}
          />
        </FormField>

        <Box>
          <Button type="submit" variant="contained" color="primary" disabled={isSubmitting}>
            {isSubmitting ? "送信中..." : submitLabel}
          </Button>
        </Box>
      </Stack>
    </Box>
  );
};
