import { Box } from "@mui/material";

import { RecordSection } from "./RecordSection";
import { recordKinds } from "./types";

export const PlanRecords = ({
  travelId,
  planId,
  publicView,
  readOnly = false,
  openByDefault = false,
}: {
  travelId: number;
  planId: number;
  publicView: boolean;
  readOnly?: boolean;
  openByDefault?: boolean;
}) => (
  <Box sx={{ mt: 2 }}>
    {recordKinds.map((kind) => (
      <RecordSection
        key={kind}
        ids={{ travelId, planId, kind }}
        publicView={publicView}
        readOnly={readOnly}
        openByDefault={openByDefault}
      />
    ))}
  </Box>
);
