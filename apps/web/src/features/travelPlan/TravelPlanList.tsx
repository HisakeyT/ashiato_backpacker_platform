import { Box, Typography } from "@mui/material";

import { sortPlans } from "./sortPlans";
import { TravelPlanCard } from "./TravelPlanCard";
import type { TravelPlan } from "./types";

type TravelPlanListProps = {
  plans: TravelPlan[];
  onEdit: (plan: TravelPlan) => void;
  onDelete: (plan: TravelPlan) => void;
  mode: "planning" | "recorded";
  onMove?: (plan: TravelPlan, target: number) => void;
};

export const TravelPlanList = ({ plans, onEdit, onDelete, mode, onMove }: TravelPlanListProps) => {
  if (plans.length === 0) {
    return <Typography color="text.secondary">まだ予定がありません</Typography>;
  }

  return (
    <Box
      className="TravelPlanList"
      sx={{
        display: "grid",
        gap: 2,
        gridTemplateColumns: "1fr",
      }}
    >
      {sortPlans(plans, mode).map((plan, index) => (
        <TravelPlanCard
          key={plan.id}
          plan={plan}
          onEdit={onEdit}
          onDelete={onDelete}
          position={mode === "planning" ? index : undefined}
          total={mode === "planning" ? plans.length : undefined}
          onMove={onMove}
        />
      ))}
    </Box>
  );
};
