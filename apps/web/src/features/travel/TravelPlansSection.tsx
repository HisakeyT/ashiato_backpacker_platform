import { Box, Button, Typography } from "@mui/material";
import { useEffect, useRef, useState } from "react";

import { ConfirmDialog } from "../../components/ConfirmDialog";

import { useAuth } from "../../app/providers/useAuth";
import { useFlash } from "../../app/providers/useFlash";
import { deleteTravelPlan, getTravelPlans, reorderTravelPlans } from "../travelPlan/repository";
import { sortPlans } from "../travelPlan/sortPlans";
import { TravelPlanDialog } from "../travelPlan/TravelPlanDialog";
import { TravelPlanList } from "../travelPlan/TravelPlanList";
import type { TravelPlan } from "../travelPlan/types";
import type { TravelDetail } from "./types";

const usePlanData = (travelId: number) => {
  const { token } = useAuth();
  const [plans, setPlans] = useState<TravelPlan[]>([]);
  const [plansError, setPlansError] = useState<string | null>(null);
  useEffect(() => {
    if (!token) return;
    let cancelled = false;
    void getTravelPlans(token, travelId).then(
      (nextPlans) => {
        if (cancelled) return;
        setPlans(nextPlans);
        setPlansError(null);
      },
      () => {
        if (!cancelled) setPlansError("予定の取得に失敗しました");
      },
    );
    return () => {
      cancelled = true;
    };
  }, [token, travelId]);
  const refreshPlans = async () => {
    if (!token) return;
    try {
      setPlans(await getTravelPlans(token, travelId));
      setPlansError(null);
    } catch {
      setPlansError("予定の取得に失敗しました");
    }
  };
  return { plans, setPlans, plansError, refreshPlans };
};

const usePlanEditor = (refreshPlans: () => Promise<void>) => {
  const { showFlash } = useFlash();
  const [isOpen, setIsOpen] = useState(false);
  const [editingPlan, setEditingPlan] = useState<TravelPlan | null>(null);
  const handleSaved = (result: "created" | "updated") => {
    void refreshPlans();
    showFlash(result === "created" ? "プランを追加しました" : "プランを更新しました");
  };
  return { isOpen, setIsOpen, editingPlan, setEditingPlan, handleSaved };
};

const usePlanDeletion = (travelId: number, refreshPlans: () => Promise<void>) => {
  const { token } = useAuth();
  const { showFlash } = useFlash();
  const [deletingPlan, setDeletingPlan] = useState<TravelPlan | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const handleClose = () => {
    if (isDeleting) return;
    setError(null);
    setDeletingPlan(null);
  };
  const handleConfirm = async () => {
    if (!token || !deletingPlan) return;
    setIsDeleting(true);
    setError(null);
    try {
      await deleteTravelPlan(token, travelId, deletingPlan.id);
      setDeletingPlan(null);
      void refreshPlans();
      showFlash("プランを削除しました");
    } catch {
      setError("プランの削除に失敗しました");
    } finally {
      setIsDeleting(false);
    }
  };
  return { deletingPlan, setDeletingPlan, isDeleting, error, handleClose, handleConfirm };
};

const PlanSectionHeader = ({ onAdd }: { onAdd: () => void }) => (
  <Box
    sx={{
      display: "flex",
      alignItems: "center",
      justifyContent: "space-between",
      flexWrap: "wrap",
      gap: 2,
      mt: 4,
      mb: 2,
    }}
  >
    <Typography variant="h5" component="h2">
      旅のプラン
    </Typography>
    <Button variant="contained" onClick={onAdd}>
      プランを追加
    </Button>
  </Box>
);

const PlanEditorDialogs = ({
  travel,
  defaultDate,
  editor,
}: {
  travel: TravelDetail;
  defaultDate: string;
  editor: ReturnType<typeof usePlanEditor>;
}) => (
  <>
    <TravelPlanDialog
      open={editor.isOpen}
      travelId={travel.id}
      defaultDate={defaultDate}
      startDate={travel.startDate}
      endDate={travel.endDate}
      mode={travel.scheduleMode}
      onClose={() => editor.setIsOpen(false)}
      onSaved={editor.handleSaved}
    />
    {editor.editingPlan && (
      <TravelPlanDialog
        open
        travelId={travel.id}
        startDate={travel.startDate}
        endDate={travel.endDate}
        mode={travel.scheduleMode}
        plan={editor.editingPlan}
        onClose={() => editor.setEditingPlan(null)}
        onSaved={editor.handleSaved}
      />
    )}
  </>
);

const PlanDeleteDialog = ({ deletion }: { deletion: ReturnType<typeof usePlanDeletion> }) =>
  deletion.deletingPlan && (
    <ConfirmDialog
      open
      heading="プランを削除しますか？"
      description={`「${deletion.deletingPlan.place}」のプランが削除されます。この操作は取り消せません。`}
      error={deletion.error}
      busy={deletion.isDeleting}
      confirmLabel="削除する"
      busyLabel="削除中..."
      confirmColor="error"
      onClose={deletion.handleClose}
      onConfirm={deletion.handleConfirm}
    />
  );

const movePlan = (plans: TravelPlan[], plan: TravelPlan, target: number): TravelPlan[] => {
  const next = [...plans];
  next.splice(
    next.findIndex((item) => item.id === plan.id),
    1,
  );
  next.splice(target, 0, plan);
  return next.map((item, index) => ({ ...item, sortOrder: index + 1 }));
};

const focusMovedPlan = (plan: TravelPlan) => {
  requestAnimationFrame(() => document.getElementById(`plan-card-${plan.id}`)?.focus());
};

const showMovedPlan = ({
  plan,
  target,
  next,
  setPlans,
  showFlash,
}: {
  plan: TravelPlan;
  target: number;
  next: TravelPlan[];
  setPlans: (plans: TravelPlan[]) => void;
  showFlash: (message: string) => void;
}) => {
  setPlans(next);
  showFlash(`${plan.place}を${target + 1}番目に移動しました`);
  focusMovedPlan(plan);
};

const usePlanOrder = (
  travelId: number,
  plans: TravelPlan[],
  setPlans: (plans: TravelPlan[]) => void,
) => {
  const { token } = useAuth();
  const { showFlash } = useFlash();
  const busy = useRef(false);
  const move = async (plan: TravelPlan, target: number) => {
    if (!token || busy.current) return;
    busy.current = true;
    const before = sortPlans(plans, "planning");
    const next = movePlan(before, plan, target);
    showMovedPlan({ plan, target, next, setPlans, showFlash });
    try {
      await reorderTravelPlans(
        token,
        travelId,
        next.map((item) => item.id),
      );
    } catch {
      setPlans(before);
      showFlash("並べ替えに失敗しました。元の順序に戻しました", "error");
    } finally {
      busy.current = false;
    }
  };
  return { move };
};

const defaultPlanDate = (travel: TravelDetail, plans: TravelPlan[]): string => {
  if (travel.scheduleMode === "planning") return "";
  const lastDate = plans
    .map((plan) => plan.endDate ?? plan.startDate ?? travel.startDate)
    .reduce((latest, date) => (latest > date ? latest : date), travel.startDate);
  return lastDate > travel.endDate ? travel.endDate : lastDate;
};

export const TravelPlansSection = ({ travel }: { travel: TravelDetail }) => {
  const { plans, setPlans, plansError, refreshPlans } = usePlanData(travel.id);
  const editor = usePlanEditor(refreshPlans);
  const deletion = usePlanDeletion(travel.id, refreshPlans);
  const order = usePlanOrder(travel.id, plans, setPlans);
  const defaultDate = defaultPlanDate(travel, plans);

  return (
    <>
      <PlanSectionHeader onAdd={() => editor.setIsOpen(true)} />
      {plansError ? (
        <Typography color="error">{plansError}</Typography>
      ) : (
        <TravelPlanList
          plans={plans}
          onEdit={editor.setEditingPlan}
          onDelete={deletion.setDeletingPlan}
          mode={travel.scheduleMode}
          onMove={travel.scheduleMode === "planning" ? order.move : undefined}
        />
      )}
      <PlanEditorDialogs travel={travel} defaultDate={defaultDate} editor={editor} />
      <PlanDeleteDialog deletion={deletion} />
    </>
  );
};
