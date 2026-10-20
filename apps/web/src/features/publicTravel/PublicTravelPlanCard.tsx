import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { Card, CardContent, Typography } from "@mui/material";

import { PlanRecords } from "../planRecord/PlanRecords";
import { formatPeriod } from "../travelPlan/format";
import type { TravelPlan } from "../travelPlan/types";

type PublicTravelPlanCardProps = {
  plan: TravelPlan;
  travelId: number;
};

export const PublicTravelPlanCard = ({ plan, travelId }: PublicTravelPlanCardProps) => {
  return (
    <Card
      component="details"
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
        <ExpandMoreIcon className="plan-expand-icon" color="action" />
      </CardContent>
      <CardContent sx={{ pt: 0, whiteSpace: "pre-wrap", overflowWrap: "anywhere" }}>
        {plan.content}
        <PlanRecords travelId={travelId} planId={plan.id} publicView />
      </CardContent>
    </Card>
  );
};
