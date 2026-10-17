import React from "react";
import { Box, Typography } from "@mui/material";
import type { TravelPlan } from "./types";
import { sortPlans } from "./sortPlans";
import { TravelPlanCard } from "./TravelPlanCard";

type TravelPlanListProps = {
  plans: TravelPlan[];
  onEdit: (plan: TravelPlan) => void;
  onDelete: (plan: TravelPlan) => void;
};


export const TravelPlanList = ({ plans, onEdit, onDelete }: TravelPlanListProps) => {
  if (plans.length === 0) {
    return <Typography color="text.secondary">まだ予定がありません</Typography>;
  }

  return (
    <Box
      className="TravelPlanList"
      sx={{
        display: "grid",
        gap: 2,
        justifyContent: "center",
        gridTemplateColumns: "repeat(auto-fill, minmax(min(100%, 220px), 1fr))",
      }}
    >
      {sortPlans(plans).map((plan) => (
        <TravelPlanCard key={plan.id} plan={plan} onEdit={onEdit} onDelete={onDelete} />
      ))}
    </Box>
  );
};
