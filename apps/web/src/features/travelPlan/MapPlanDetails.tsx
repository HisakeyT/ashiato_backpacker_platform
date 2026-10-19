import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from "@mui/material";

import { formatPeriod } from "./format";
import type { MapPlan } from "./mapPlans";

export const MapPlanDetails = ({
  selected,
  titleId,
  onClose,
}: {
  selected: MapPlan[];
  titleId: string;
  onClose: () => void;
}) => (
  <Dialog
    open={selected.length > 0}
    onClose={onClose}
    aria-labelledby={titleId}
    fullWidth
    maxWidth="sm"
  >
    <DialogTitle id={titleId}>地点のプラン</DialogTitle>
    <DialogContent dividers>
      {selected.map(({ plan, number }) => (
        <Box key={plan.id} sx={{ mb: 2, overflowWrap: "anywhere" }}>
          <Typography component="h3" variant="h6">
            {number}. {plan.place}
          </Typography>
          <Typography color="text.secondary">
            {formatPeriod(plan.startDate, plan.endDate)}
          </Typography>
          <Typography sx={{ whiteSpace: "pre-wrap", mt: 1 }}>
            {plan.content || "内容は未登録です"}
          </Typography>
        </Box>
      ))}
    </DialogContent>
    <DialogActions>
      <Button onClick={onClose}>閉じる</Button>
    </DialogActions>
  </Dialog>
);
