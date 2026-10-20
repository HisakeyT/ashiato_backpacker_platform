import { Box, Button, Menu, MenuItem, Paper, Typography } from "@mui/material";
import { useState } from "react";

import { sortPlans } from "../sortPlans";
import type { TravelPlan } from "../types";

type MoveProps = {
  plan: TravelPlan;
  position: number;
  total: number;
  isSaving: boolean;
  onMove: (plan: TravelPlan, target: number) => void;
};

const PlanMoveControls = ({ plan, position, total, isSaving, onMove }: MoveProps) => {
  const [menuAnchor, setMenuAnchor] = useState<HTMLElement | null>(null);
  const jump = (target: number) => {
    setMenuAnchor(null);
    onMove(plan, target);
  };

  return (
    <Box sx={{ display: "flex", flexWrap: "wrap", justifyContent: "flex-end" }}>
      <Button size="small" disabled={isSaving || position === 0} onClick={() => jump(position - 1)}>
        上へ
      </Button>
      <Button
        size="small"
        disabled={isSaving || position === total - 1}
        onClick={() => jump(position + 1)}
      >
        下へ
      </Button>
      <Button
        size="small"
        disabled={isSaving}
        onClick={(event) => setMenuAnchor(event.currentTarget)}
      >
        移動
      </Button>
      <Menu anchorEl={menuAnchor} open={Boolean(menuAnchor)} onClose={() => setMenuAnchor(null)}>
        <MenuItem disabled={position === 0} onClick={() => jump(0)}>
          先頭へ移動
        </MenuItem>
        <MenuItem disabled={position === total - 1} onClick={() => jump(total - 1)}>
          末尾へ移動
        </MenuItem>
      </Menu>
    </Box>
  );
};

type Props = {
  plans: TravelPlan[];
  isSaving: boolean;
  error: string | null;
  onMove: (plan: TravelPlan, target: number) => void;
  onDone: () => void;
};

export const PlanOrderEditor = ({ plans, isSaving, error, onMove, onDone }: Props) => (
  <Box component="section" aria-label="プランの順序変更" sx={{ mt: 4 }}>
    <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", gap: 2 }}>
      <Typography component="h2" variant="h5">
        プランの順序を変更
      </Typography>
      <Button variant="contained" disabled={isSaving} onClick={onDone}>
        完了
      </Button>
    </Box>
    <Typography color="text.secondary" sx={{ mt: 1, mb: 2 }}>
      変更するたびに自動保存されます。
    </Typography>
    {error ? (
      <Typography color="error">{error}</Typography>
    ) : (
      <Box sx={{ display: "grid", gap: 1 }}>
        {sortPlans(plans, "planning").map((plan, position) => (
          <Paper
            key={plan.id}
            id={`plan-card-${plan.id}`}
            tabIndex={-1}
            sx={{ p: 1.5, display: "flex", alignItems: "center", flexWrap: "wrap", gap: 1 }}
          >
            <Typography color="text.secondary" sx={{ minWidth: 48 }}>
              {position + 1} / {plans.length}
            </Typography>
            <Typography sx={{ flexGrow: 1, overflowWrap: "anywhere" }}>{plan.place}</Typography>
            <PlanMoveControls
              plan={plan}
              position={position}
              total={plans.length}
              isSaving={isSaving}
              onMove={onMove}
            />
          </Paper>
        ))}
      </Box>
    )}
  </Box>
);
