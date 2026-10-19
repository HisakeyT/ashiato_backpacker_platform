import { Box, Typography } from "@mui/material";

import { sortPlans } from "../travelPlan/sortPlans";
import type { TravelPlan } from "../travelPlan/types";
import { PublicTravelPlanCard } from "./PublicTravelPlanCard";

type Props = {
  plans: TravelPlan[];
};

export const PublicTravelPlanList = ({ plans }: Props) => {
  if (plans.length === 0) {
    return <Typography color="text.secondary">まだ予定がありません</Typography>;
  }

  return (
    <Box
      className="PublicTravelPlanList"
      sx={{
        display: "grid",
        gap: 2,
        justifyContent: "center",
        gridTemplateColumns: "repeat(auto-fill, minmax(min(100%, 220px), 1fr))",
      }}
    >
      {sortPlans(plans).map((plan) => (
        <PublicTravelPlanCard key={plan.id} plan={plan} />
      ))}
    </Box>
  );
};
