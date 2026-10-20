import { Box, Button, Card, CardActions, CardContent, Chip, Typography } from "@mui/material";

import { formatPeriod } from "./format";
import type { TravelPlan } from "./types";

type TravelPlanCardProps = {
  plan: TravelPlan;
  onEdit: (plan: TravelPlan) => void;
  onDelete: (plan: TravelPlan) => void;
  orderNumber?: number;
};

export const TravelPlanCard = ({ plan, onEdit, onDelete, orderNumber }: TravelPlanCardProps) => {
  return (
    <Card
      className="TravelPlanCard"
      sx={{ height: "100%", display: "flex", flexDirection: "column" }}
    >
      <CardContent sx={{ flexGrow: 1 }}>
        <Box
          sx={{
            display: "flex",
            alignItems: "flex-start",
            justifyContent: "space-between",
            gap: 1,
            mb: 1,
          }}
        >
          <Typography variant="h6" sx={{ overflowWrap: "anywhere" }}>
            {plan.place}
          </Typography>
          {orderNumber !== undefined && (
            <Chip
              size="small"
              variant="outlined"
              label={orderNumber}
              aria-label={`訪問順 ${orderNumber}`}
            />
          )}
        </Box>
        <Typography variant="body2" color="text.secondary" gutterBottom>
          {formatPeriod(plan.startDate, plan.endDate)}
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ whiteSpace: "pre-wrap" }}>
          {plan.content}
        </Typography>
      </CardContent>
      <CardActions sx={{ flexWrap: "wrap", justifyContent: "flex-end" }}>
        <Button size="small" onClick={() => onEdit(plan)}>
          編集
        </Button>
        <Button size="small" color="error" onClick={() => onDelete(plan)}>
          削除
        </Button>
      </CardActions>
    </Card>
  );
};
