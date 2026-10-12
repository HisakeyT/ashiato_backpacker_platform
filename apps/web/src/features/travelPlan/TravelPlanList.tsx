import { Box, Typography } from "@mui/material";
import type { TravelPlan } from "./types";
import { TravelPlanCard } from "./TravelPlanCard";

type TravelPlanListProps = {
  plans: TravelPlan[];
  onEdit: (plan: TravelPlan) => void;
  onDelete: (plan: TravelPlan) => void;
};

const sortPlans = (plans: TravelPlan[]): TravelPlan[] => {
  return [...plans].sort((a, b) => {
    if (a.startDate !== b.startDate) return a.startDate < b.startDate ? -1 : 1;
    return a.sortOrder - b.sortOrder;
  });
};

const groupByDate = (plans: TravelPlan[]): [string, TravelPlan[]][] => {
  const map = new Map<string, TravelPlan[]>();
  for (const plan of plans) {
    const list = map.get(plan.startDate) ?? [];
    list.push(plan);
    map.set(plan.startDate, list);
  }
  return [...map.entries()];
};

export const TravelPlanList = ({ plans, onEdit, onDelete }: TravelPlanListProps) => {
  if (plans.length === 0) {
    return <Typography color="text.secondary">まだ予定がありません</Typography>;
  }

  const grouped = groupByDate(sortPlans(plans));

  return (
    <Box className="TravelPlanList" sx={{ display: "grid", gap: 4 }}>
      {grouped.map(([date, dayPlans]) => (
        <Box key={date}>
          <Typography variant="h6" gutterBottom>
            {date}
          </Typography>
          <Box
            sx={{
              display: "grid",
              gap: 2,
              justifyContent: "center",
              gridTemplateColumns: "repeat(auto-fill, minmax(min(100%, 220px), 1fr))",
            }}
          >
            {dayPlans.map((plan) => (
              <TravelPlanCard key={plan.id} plan={plan} onEdit={onEdit} onDelete={onDelete} />
            ))}
          </Box>
        </Box>
      ))}
    </Box>
  );
};
