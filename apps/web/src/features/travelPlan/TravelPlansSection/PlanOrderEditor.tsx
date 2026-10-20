import { Box, Button, Menu, MenuItem, Paper, Typography } from "@mui/material";
import { useEffect, useRef, useState } from "react";

import { useFlash } from "../../../app/providers/useFlash";
import { sortPlans } from "../sortPlans";
import type { TravelPlan } from "../types";
import { movePlan } from "./planOrderChanges";

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
        <MenuItem disabled={isSaving || position === 0} onClick={() => jump(0)}>
          先頭へ移動
        </MenuItem>
        <MenuItem disabled={isSaving || position === total - 1} onClick={() => jump(total - 1)}>
          末尾へ移動
        </MenuItem>
      </Menu>
    </Box>
  );
};

type Props = {
  plans: TravelPlan[];
  onSave: (ordered: TravelPlan[]) => Promise<void>;
  onClose: () => void;
};

const useDraftOrder = ({ plans, onSave, onClose }: Props) => {
  const { showFlash } = useFlash();
  const busy = useRef(false);
  const [ordered, setOrdered] = useState(() => sortPlans(plans, "planning"));
  const [isSaving, setIsSaving] = useState(false);

  const handleMove = (plan: TravelPlan, target: number) => {
    if (busy.current) return;
    setOrdered((current) => movePlan(current, plan, target));
    showFlash(`${plan.place}を${target + 1}番目に移動しました。まだ保存されていません。`, "info");
    requestAnimationFrame(() => document.getElementById(`plan-card-${plan.id}`)?.focus());
  };

  const handleDone = async () => {
    if (busy.current) return;
    busy.current = true;
    setIsSaving(true);
    try {
      await onSave(ordered);
      onClose();
    } catch {
      showFlash("順序の保存に失敗しました。もう一度お試しください。", "error");
    } finally {
      busy.current = false;
      setIsSaving(false);
    }
  };

  return { ordered, isSaving, handleMove, handleDone };
};

type RowsProps = {
  ordered: TravelPlan[];
  isSaving: boolean;
  onMove: (plan: TravelPlan, target: number) => void;
};

const PlanOrderRows = ({ ordered, isSaving, onMove }: RowsProps) => (
  <Box sx={{ display: "grid", gap: 1 }}>
    {ordered.map((plan, position) => (
      <Paper
        key={plan.id}
        id={`plan-card-${plan.id}`}
        tabIndex={-1}
        sx={{ p: 1.5, display: "flex", alignItems: "center", flexWrap: "wrap", gap: 1 }}
      >
        <Typography color="text.secondary" sx={{ minWidth: 48 }}>
          {position + 1} / {ordered.length}
        </Typography>
        <Typography sx={{ flexGrow: 1, overflowWrap: "anywhere" }}>{plan.place}</Typography>
        <PlanMoveControls
          plan={plan}
          position={position}
          total={ordered.length}
          isSaving={isSaving}
          onMove={onMove}
        />
      </Paper>
    ))}
  </Box>
);

export const PlanOrderEditor = (props: Props) => {
  const headingRef = useRef<HTMLHeadingElement>(null);
  const { ordered, isSaving, handleMove, handleDone } = useDraftOrder(props);
  useEffect(() => headingRef.current?.focus(), []);

  return (
    <Box component="section" aria-label="プランの順序変更" sx={{ mt: 4 }}>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", gap: 2 }}>
        <Typography component="h2" variant="h5" ref={headingRef} tabIndex={-1}>
          プランの順序を変更
        </Typography>
        <Box sx={{ display: "flex", gap: 1 }}>
          <Button disabled={isSaving} onClick={props.onClose}>
            キャンセル
          </Button>
          <Button variant="contained" disabled={isSaving} onClick={handleDone}>
            {isSaving ? "保存中…" : "完了"}
          </Button>
        </Box>
      </Box>
      <Typography color="text.secondary" sx={{ mt: 1, mb: 2 }}>
        並び替えた後、「完了」を押すと保存されます。
      </Typography>
      <PlanOrderRows ordered={ordered} isSaving={isSaving} onMove={handleMove} />
    </Box>
  );
};
