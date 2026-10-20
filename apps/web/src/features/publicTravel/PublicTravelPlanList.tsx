import { Box, Typography } from "@mui/material";

import { sortPlans } from "../travelPlan/sortPlans";
import type { TravelPlan } from "../travelPlan/types";
import { PublicTravelPlanCard } from "./PublicTravelPlanCard";

type Props = {
  plans: TravelPlan[];
  mode: "planning" | "recorded";
  travelId: number;
};

export const PublicTravelPlanList = ({ plans, mode, travelId }: Props) => {
  if (plans.length === 0) {
    return <Typography color="text.secondary">まだ予定がありません</Typography>;
  }

  return (
    <Box
      className="PublicTravelPlanList"
      sx={{
        display: "grid",
        gap: 2,
        gridTemplateColumns: "1fr",
      }}
    >
      {sortPlans(plans, mode).map((plan) => (
        <PublicTravelPlanCard key={plan.id} plan={plan} travelId={travelId} />
      ))}
    </Box>
  );
};
