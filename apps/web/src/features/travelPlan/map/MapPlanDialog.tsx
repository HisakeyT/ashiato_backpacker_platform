import { Box, Button, Typography } from "@mui/material";

import { BaseDialog } from "../../../components/BaseDialog";

import { PlanRecords } from "../../planRecord/PlanRecords";
import { formatPeriod } from "../format";
import type { MapPlan } from "./mapPlans";

export const MapPlanDialog = ({
  selected,
  titleId,
  onClose,
  travelId,
  publicView,
}: {
  selected: MapPlan[];
  titleId: string;
  onClose: () => void;
  travelId: number;
  publicView: boolean;
}) => (
  <BaseDialog
    open={selected.length > 0}
    onClose={onClose}
    title="地点のプラン"
    titleId={titleId}
    maxWidth="sm"
    actions={<Button onClick={onClose}>閉じる</Button>}
  >
    {selected.map(({ plan, number }) => (
      <Box key={plan.id} sx={{ mb: 2, overflowWrap: "anywhere" }}>
        <Typography component="h3" variant="h6">
          {number}. {plan.place}
        </Typography>
        <Typography color="text.secondary">{formatPeriod(plan.startDate, plan.endDate)}</Typography>
        <Typography sx={{ whiteSpace: "pre-wrap", mt: 1 }}>{plan.content}</Typography>
        <PlanRecords
          key={plan.id}
          travelId={travelId}
          planId={plan.id}
          publicView={publicView}
          readOnly
          openByDefault
        />
      </Box>
    ))}
  </BaseDialog>
);
