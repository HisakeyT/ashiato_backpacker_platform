import { Box, Button, MenuItem, Stack, TextField } from "@mui/material";

import { mealLabels } from "./types";
import type { RecordInput, RecordKind } from "./types";

// eslint-disable-next-line eslint/max-lines-per-function
export const RecordForm = ({
  kind,
  value,
  onChange,
  onSubmit,
  onCancel,
  busy,
  editing,
}: {
  kind: RecordKind;
  value: RecordInput;
  onChange: (value: RecordInput) => void;
  onSubmit: () => void;
  onCancel: () => void;
  busy: boolean;
  editing: boolean;
}) => (
  <Stack
    component="form"
    spacing={1}
    sx={{ my: 2 }}
    onSubmit={(event) => {
      event.preventDefault();
      onSubmit();
    }}
  >
    <TextField
      label="名前"
      required
      size="small"
      value={value.name}
      onChange={(event) => onChange({ ...value, name: event.target.value })}
    />
    <TextField
      label="場所・関連情報"
      multiline
      minRows={2}
      size="small"
      value={value.locationText ?? ""}
      onChange={(event) => onChange({ ...value, locationText: event.target.value || null })}
    />
    <TextField
      label="地図リンク"
      type="url"
      size="small"
      value={value.mapUrl ?? ""}
      onChange={(event) => onChange({ ...value, mapUrl: event.target.value || null })}
    />
    {kind === "meals" && (
      <TextField
        select
        label="食事区分"
        size="small"
        value={value.mealType ?? ""}
        onChange={(event) => onChange({ ...value, mealType: event.target.value || null })}
      >
        <MenuItem value="">未選択</MenuItem>
        {Object.entries(mealLabels).map(([key, label]) => (
          <MenuItem key={key} value={key}>
            {label}
          </MenuItem>
        ))}
      </TextField>
    )}
    <Box sx={{ display: "flex", gap: 1 }}>
      <Button type="submit" variant="contained" disabled={busy || !value.name.trim()}>
        {editing ? "更新" : "追加"}
      </Button>
      {editing && (
        <Button onClick={onCancel} disabled={busy}>
          キャンセル
        </Button>
      )}
    </Box>
  </Stack>
);
