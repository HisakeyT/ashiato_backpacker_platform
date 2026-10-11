import { Card, CardContent, Typography } from "@mui/material";
import type { TravelPlan } from "../travelPlan/types";

type PublicTravelPlanCardProps = {
  plan: TravelPlan;
};

export const PublicTravelPlanCard = ({ plan }: PublicTravelPlanCardProps) => {
  return (
    <Card
      className="TravelPlanCard"
      sx={{ height: "100%", display: "flex", flexDirection: "column" }}
    >
      <CardContent sx={{ flexGrow: 1 }}>
        <Typography variant="h6" gutterBottom>
          {plan.place}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {plan.content}
        </Typography>
      </CardContent>
    </Card>
  );
};
