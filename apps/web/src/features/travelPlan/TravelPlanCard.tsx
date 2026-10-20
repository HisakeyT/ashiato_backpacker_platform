import { Button, Card, CardActions, CardContent, Typography } from "@mui/material";

import { formatPeriod } from "./format";
import type { TravelPlan } from "./types";

type TravelPlanCardProps = {
  plan: TravelPlan;
  onEdit: (plan: TravelPlan) => void;
  onDelete: (plan: TravelPlan) => void;
  position?: number;
  total?: number;
};

export const TravelPlanCard = ({
  plan,
  onEdit,
  onDelete,
  position,
  total,
}: TravelPlanCardProps) => {
  return (
    <Card
      id={`plan-card-${plan.id}`}
      tabIndex={-1}
      className="TravelPlanCard"
      sx={{ display: "flex", flexDirection: "column" }}
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
        {position !== undefined && total !== undefined && (
          <Typography variant="caption" color="text.secondary">
            {position + 1} / {total}
          </Typography>
        )}
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
