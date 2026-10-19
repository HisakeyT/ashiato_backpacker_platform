import { Box, Button, Typography } from "@mui/material";
import { useState } from "react";

import { ConfirmDialog } from "../../components/ConfirmDialog";

import { useFlash } from "../../app/providers/useFlash";
import type { TravelDetail } from "../travel/types";
import { TravelMap } from "./map/TravelMap";
import { TravelPlanDialog } from "./TravelPlanDialog";
import { TravelPlanList } from "./TravelPlanList";
import type { TravelPlan } from "./types";
import { usePlanData } from "./usePlanData";
import { usePlanDeletion } from "./usePlanDeletion";
import { usePlanOrder } from "./usePlanOrder";

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
      {!plansError && <TravelMap plans={plans} />}
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
