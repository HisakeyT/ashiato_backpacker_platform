import { Button, Card, CardActions, CardContent, Typography } from "@mui/material";

import { formatPeriod } from "./format";
import type { TravelPlan } from "./types";

type TravelPlanCardProps = {
  plan: TravelPlan;
  onEdit: (plan: TravelPlan) => void;
  onDelete: (plan: TravelPlan) => void;
};

export const TravelPlanCard = ({ plan, onEdit, onDelete }: TravelPlanCardProps) => {
  return (
    <Card
      className="TravelPlanCard"
      sx={{ height: "100%", display: "flex", flexDirection: "column" }}
    >
      <CardContent sx={{ flexGrow: 1 }}>
        <Typography variant="caption" color="text.secondary">
          {formatPeriod(plan.startDate, plan.endDate)}
        </Typography>
        <Typography variant="h6" gutterBottom>
          {plan.place}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {plan.content}
        </Typography>
      </CardContent>
      <CardActions sx={{ justifyContent: "flex-end" }}>
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
