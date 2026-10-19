import { Button, Card, CardActions, CardContent, Menu, MenuItem, Typography } from "@mui/material";
import { useState } from "react";

import { formatPeriod } from "./format";
import type { TravelPlan } from "./types";

type TravelPlanCardProps = {
  plan: TravelPlan;
  onEdit: (plan: TravelPlan) => void;
  onDelete: (plan: TravelPlan) => void;
  position?: number;
  total?: number;
  onMove?: (plan: TravelPlan, target: number) => void;
};

type MoveProps = {
  plan: TravelPlan;
  position: number;
  total: number;
  onMove: (plan: TravelPlan, target: number) => void;
};

const PlanJumpMenu = ({ plan, position, total, onMove }: MoveProps) => {
  const [menuAnchor, setMenuAnchor] = useState<HTMLElement | null>(null);
  const jump = (target: number) => {
    setMenuAnchor(null);
    onMove(plan, target);
  };
  return (
    <>
      <Button size="small" onClick={(event) => setMenuAnchor(event.currentTarget)}>
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
    </>
  );
};

const PlanMoveControls = ({ plan, position, total, onMove }: MoveProps) => (
  <>
    <Button size="small" disabled={position === 0} onClick={() => onMove(plan, position - 1)}>
      上へ
    </Button>
    <Button
      size="small"
      disabled={position === total - 1}
      onClick={() => onMove(plan, position + 1)}
    >
      下へ
    </Button>
    <PlanJumpMenu plan={plan} position={position} total={total} onMove={onMove} />
  </>
);

export const TravelPlanCard = ({
  plan,
  onEdit,
  onDelete,
  position,
  total,
  onMove,
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
        {position !== undefined && total !== undefined && onMove && (
          <PlanMoveControls plan={plan} position={position} total={total} onMove={onMove} />
        )}
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
