import { Alert, Box, Button, Typography } from "@mui/material";
import { useState } from "react";

import { useAuth } from "../../app/providers/useAuth";
import { RecordEditor } from "./RecordEditor";
import { RecordItem } from "./RecordItem";
import type { RecordIDs } from "./repository";
import { recordLabels } from "./types";
import { useRecordData } from "./useRecordData";

export const RecordSection = ({
  ids,
  publicView,
  readOnly = false,
  openByDefault = false,
}: {
  ids: RecordIDs;
  publicView: boolean;
  readOnly?: boolean;
  openByDefault?: boolean;
}) => {
  const { token } = useAuth();
  const [active, setActive] = useState(openByDefault);
  const { records, loading, error, load } = useRecordData({ ids, publicView, token, active });
  return (
    <Box
      component="details"
      open={active}
      onToggle={(event) => setActive(event.currentTarget.open)}
      sx={{ mt: 1 }}
    >
      <Box component="summary" sx={{ cursor: "pointer", fontWeight: "bold" }}>
        {recordLabels[ids.kind]}
      </Box>
      <Box sx={{ pl: 1 }}>
        {loading && <Typography variant="body2">読み込み中...</Typography>}
        {error && (
          <Alert severity="error" sx={{ my: 1 }}>
            {error}
            <Button onClick={() => void load()}>再試行</Button>
          </Alert>
        )}
        {!loading && !error && records.length === 0 && (
          <Typography variant="body2" color="text.secondary">
            記録はありません
          </Typography>
        )}
        {records.map((record) => (
          <RecordItem key={record.id} record={record} kind={ids.kind} />
        ))}
        {!publicView && !readOnly && token && (
          <RecordEditor ids={ids} token={token} records={records} load={load} />
        )}
      </Box>
    </Box>
  );
};
