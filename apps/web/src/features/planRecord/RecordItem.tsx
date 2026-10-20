import { Box, Button, Typography } from "@mui/material";

import { RecordText } from "./RecordText";
import { mealLabels } from "./types";
import type { PlanRecord, RecordKind } from "./types";

export const RecordItem = ({ record, kind }: { record: PlanRecord; kind: RecordKind }) => (
  <Box sx={{ py: 1, overflowWrap: "anywhere" }}>
    <Typography sx={{ fontWeight: "bold" }}>{record.name}</Typography>
    {kind === "meals" && record.mealType && (
      <Typography variant="body2" color="text.secondary">
        {mealLabels[record.mealType]}
      </Typography>
    )}
    {record.locationText && <RecordText text={record.locationText} />}
    {record.mapUrl && (
      <Button
        component="a"
        size="small"
        href={record.mapUrl}
        target="_blank"
        rel="noopener noreferrer"
      >
        地図を開く
      </Button>
    )}
  </Box>
);
