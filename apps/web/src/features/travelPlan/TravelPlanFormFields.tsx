import {
  Alert,
  Box,
  Button,
  List,
  ListItemButton,
  ListItemText,
  Stack,
  TextField,
} from "@mui/material";

import { FormField } from "../../components/FormField";

import type { PlaceSearchState, PlanDetails } from "./useTravelPlanForm";

export const PlanDates = ({
  details,
  minDate,
  maxDate,
  mode,
}: {
  details: PlanDetails;
  minDate?: string;
  maxDate?: string;
  mode: "planning" | "recorded";
}) => (
  <>
    <FormField
      label={mode === "planning" ? "開始日（任意）" : "開始日"}
      htmlFor="plan-start-date"
      required={mode === "recorded"}
    >
      <TextField
        id="plan-start-date"
        type="date"
        value={details.startDate}
        onChange={(e) => {
          details.setStartDate(e.target.value);
          if (!e.target.value) details.setEndDate("");
        }}
        required={mode === "recorded"}
        fullWidth
        slotProps={{ htmlInput: { min: minDate, max: maxDate } }}
      />
    </FormField>
    <FormField label="終了日(任意)" htmlFor="plan-end-date">
      <TextField
        id="plan-end-date"
        type="date"
        value={details.endDate}
        onChange={(e) => details.setEndDate(e.target.value)}
        disabled={!details.startDate}
        fullWidth
        slotProps={{ htmlInput: { min: details.startDate || minDate, max: maxDate } }}
      />
    </FormField>
  </>
);

const PlaceSearchInput = ({
  details,
  search,
}: {
  details: PlanDetails;
  search: PlaceSearchState;
}) => (
  <Box sx={{ display: "flex", alignItems: "flex-start", gap: 1 }}>
    <TextField
      id="plan-place"
      value={details.place}
      onChange={(e) => details.setPlace(e.target.value)}
      onKeyDown={(e) => {
        if (e.key === "Enter" && !e.nativeEvent.isComposing) {
          e.preventDefault();
          void search.handleSearch();
        }
      }}
      placeholder="例：アンコール・ワット"
      helperText="国名や市名を足す、または別の表記(カタカナ・英語名など)で検索できます"
      required
      fullWidth
    />
    <Button
      variant="outlined"
      onClick={() => void search.handleSearch()}
      disabled={search.isSearching || !details.place.trim()}
      sx={{ whiteSpace: "nowrap", height: 56, flexShrink: 0 }}
    >
      {search.isSearching ? "検索中..." : "検索"}
    </Button>
  </Box>
);

const PlaceResults = ({ search }: { search: PlaceSearchState }) => (
  <>
    {search.searchError && <Alert severity="error">{search.searchError}</Alert>}
    {search.searched && search.candidates.length === 0 && !search.searchError && (
      <Alert severity="info">
        候補が見つかりませんでした。国名や市名を足す、または別の表記(カタカナ・英語名など)で試せます。
        検索せずに、このまま手入力で保存することもできます。
      </Alert>
    )}
    {search.candidates.length > 0 && (
      <List dense disablePadding sx={{ border: 1, borderColor: "divider", borderRadius: 1 }}>
        {search.candidates.map((candidate) => (
          <ListItemButton
            key={`${candidate.latitude},${candidate.longitude}`}
            onClick={() => search.handleSelect(candidate)}
          >
            <ListItemText primary={candidate.name} secondary={candidate.label} />
          </ListItemButton>
        ))}
      </List>
    )}
  </>
);

export const PlanPlace = ({
  details,
  search,
}: {
  details: PlanDetails;
  search: PlaceSearchState;
}) => (
  <FormField label="場所" htmlFor="plan-place" required>
    <Stack spacing={1}>
      <PlaceSearchInput details={details} search={search} />
      <PlaceResults search={search} />
    </Stack>
  </FormField>
);

export const PlanLocation = ({ details }: { details: PlanDetails }) => (
  <Stack direction={{ xs: "column", sm: "row" }} spacing={2}>
    <FormField label="国(任意)" htmlFor="plan-country">
      <TextField
        id="plan-country"
        value={details.country}
        onChange={(e) => details.setCountry(e.target.value)}
        fullWidth
      />
    </FormField>
    <FormField label="州・県(任意)" htmlFor="plan-state">
      <TextField
        id="plan-state"
        value={details.state}
        onChange={(e) => details.setState(e.target.value)}
        fullWidth
      />
    </FormField>
    <FormField label="都市(任意)" htmlFor="plan-city">
      <TextField
        id="plan-city"
        value={details.city}
        onChange={(e) => details.setCity(e.target.value)}
        fullWidth
      />
    </FormField>
  </Stack>
);

export const PlanContent = ({ details }: { details: PlanDetails }) => (
  <FormField label="内容" htmlFor="plan-content" required>
    <TextField
      id="plan-content"
      value={details.content}
      onChange={(e) => details.setContent(e.target.value)}
      placeholder="例：朝日を見てから遺跡を回る"
      required
      fullWidth
      multiline
      minRows={3}
    />
  </FormField>
);

export const PlanActions = ({
  onCancel,
  isSubmitting,
  submitLabel,
}: {
  onCancel: () => void;
  isSubmitting: boolean;
  submitLabel: string;
}) => (
  <Box sx={{ display: "flex", justifyContent: "flex-end", gap: 1 }}>
    <Button onClick={onCancel} color="inherit" disabled={isSubmitting}>
      キャンセル
    </Button>
    <Button type="submit" variant="contained" disabled={isSubmitting}>
      {isSubmitting ? "送信中..." : submitLabel}
    </Button>
  </Box>
);
