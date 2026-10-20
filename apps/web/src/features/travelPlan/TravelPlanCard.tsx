import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
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
      component="details"
      id={`plan-card-${plan.id}`}
      tabIndex={-1}
      className="TravelPlanCard"
      sx={{ "&[open] .plan-expand-icon": { transform: "rotate(180deg)" } }}
    >
      <CardContent
        component="summary"
        sx={{ cursor: "pointer", display: "flex", alignItems: "center", gap: 2, flexWrap: "wrap" }}
      >
        <Typography variant="body1" sx={{ flexGrow: 1, overflowWrap: "anywhere" }}>
          {plan.place}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {formatPeriod(plan.startDate, plan.endDate)}
        </Typography>
        {position !== undefined && total !== undefined && (
          <Typography variant="caption" color="text.secondary">
            {position + 1} / {total}
          </Typography>
        )}
        <ExpandMoreIcon className="plan-expand-icon" color="action" />
      </CardContent>
      <CardContent sx={{ pt: 0, whiteSpace: "pre-wrap", overflowWrap: "anywhere" }}>
        {plan.content}
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
