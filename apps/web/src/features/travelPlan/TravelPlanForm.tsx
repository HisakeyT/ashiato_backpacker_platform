import { useState, type SubmitEvent } from "react";
import {
  Alert,
  Box,
  Button,
  List,
  ListItemButton,
  ListItemText,
  Stack,
  TextField,
} from "@mui/material"
import { FormField } from "../../components/FormField";
import { searchPlaces, type PlaceCandidate } from "./searchPlaces";

export type TravelPlanFormValue = {
  startDate: string;
  endDate: string; // 空文字 = 未指定(同日)
  place: string;
  country: string; // 空文字 = 未指定
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
  const [endDate, setEndDate] = useState(initialValue?.endDate ?? "");
  const [place, setPlace] = useState(initialValue?.place ?? "");
  const [country, setCountry] = useState(initialValue?.country ?? "");
  const [state, setState] = useState(initialValue?.state ?? "");
  const [city, setCity] = useState(initialValue?.city ?? "");
  const [latitude, setLatitude] = useState<number | null>(initialValue?.latitude ?? null);
  const [longitude, setLongitude] = useState<number | null>(initialValue?.longitude ?? null);
  const [content, setContent] = useState(initialValue?.content ?? "");

  const [candidates, setCandidates] = useState<PlaceCandidate[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchError, setSearchError] = useState<string | null>(null);
  const [searched, setSearched] = useState(false);

  const handleSearch = async () => {
    if (isSearching || !place.trim()) return;
    setIsSearching(true);
    setSearchError(null);
    try {
      const result = await searchPlaces(place);
      setCandidates(result);
      setSearched(true);
    } catch {
      setSearchError("場所の検索に失敗しました");
    } finally {
      setIsSearching(false);
    }
  };

  const handleSelect = (c: PlaceCandidate) => {
    setPlace(c.name);
    setCountry(c.country ?? "");
    setState(c.state ?? "");
    setCity(c.city ?? "");
    setLatitude(c.latitude);
    setLongitude(c.longitude);
    setCandidates([]);
    setSearched(false);
  };

  const handleSubmit = async (e: SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (isSubmitting) return;
    await onSubmit({
      startDate,
      endDate,
      place,
      country,
      state,
      city,
      latitude,
      longitude,
      content,
    });
  };

  return (
    <Box className="TravelPlanForm" component="form" onSubmit={handleSubmit}>
      <Stack spacing={3}>
        {error && <Alert severity="error">{error}</Alert>}

        <FormField label="開始日" htmlFor="plan-start-date" required>
          <TextField
            id="plan-start-date"
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
            required
            fullWidth
            slotProps={{ htmlInput: { min: minDate, max: maxDate } }}
          />
        </FormField>

        <FormField label="終了日(任意)" htmlFor="plan-end-date">
          <TextField
            id="plan-end-date"
            type="date"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
            fullWidth
            slotProps={{
              htmlInput: { min: startDate || minDate, max: maxDate },
            }}
          />
        </FormField>

        <FormField label="場所" htmlFor="plan-place" required>
          <Stack spacing={1}>
            <Box sx={{ display: "flex", gap: 1 }}>
              <TextField
                id="plan-place"
                value={place}
                onChange={(e) => setPlace(e.target.value)}
                onKeyDown={(e) => {
                  // IME 変換確定の Enter では検索しない
                  if (e.key === "Enter" && !e.nativeEvent.isComposing) {
                    e.preventDefault();
                    void handleSearch();
                  }
                }}
                placeholder="例：アンコールワット"
                required
                fullWidth
              />
              <Button
                variant="outlined"
                onClick={() => void handleSearch()}
                disabled={isSearching || !place.trim()}
                sx={{ whiteSpace: "nowrap" }}
              >
                {isSearching ? "検索中..." : "検索"}
              </Button>
            </Box>

            {searchError && <Alert severity="error">{searchError}</Alert>}
            {searched && candidates.length === 0 && !searchError && (
              <Alert severity="info">候補が見つかりませんでした。手入力のまま保存できます。</Alert>
            )}
            {candidates.length > 0 && (
              <List dense disablePadding sx={{ border: 1, borderColor: "divider", borderRadius: 1 }}>
                {candidates.map((c) => (
                  <ListItemButton key={`${c.latitude},${c.longitude}`} onClick={() => handleSelect(c)}>
                    <ListItemText primary={c.name} secondary={c.label} />
                  </ListItemButton>
                ))}
              </List>
            )}
          </Stack>
        </FormField>

        <Stack direction={{ xs: "column", sm: "row" }} spacing={2}>
          <FormField label="国(任意)" htmlFor="plan-country">
            <TextField
              id="plan-country"
              value={country}
              onChange={(e) => setCountry(e.target.value)}
              fullWidth
            />
          </FormField>
          <FormField label="州・県(任意)" htmlFor="plan-state">
            <TextField
              id="plan-state"
              value={state}
              onChange={(e) => setState(e.target.value)}
              fullWidth
            />
          </FormField>
          <FormField label="都市(任意)" htmlFor="plan-city">
            <TextField
              id="plan-city"
              value={city}
              onChange={(e) => setCity(e.target.value)}
              fullWidth
            />
          </FormField>
        </Stack>

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
