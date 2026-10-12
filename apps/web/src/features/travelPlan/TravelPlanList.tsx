import { useEffect, useState } from "react";
import { Box, Typography } from "@mui/material";
import { useAuth } from "../../app/providers/useAuth";
import { useFlash } from "../../app/providers/useFlash";
import { getTravelPlans, deleteTravelPlan } from "./repository";
import type { TravelPlan } from "./types";
import { TravelPlanCard } from "./TravelPlanCard";
import { TravelPlanDialog } from "./TravelPlanDialog";
import { ConfirmDialog } from "../../components/ConfirmDialog";

type TravelPlanListProps = {
  travelId: number;
  startDate: string;
  endDate: string;
  reloadKey?: number;
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

export const TravelPlanList = ({ travelId, startDate, endDate, reloadKey = 0 }: TravelPlanListProps) => {
  const { token } = useAuth();
  const { showFlash } = useFlash();
  const [plans, setPlans] = useState<TravelPlan[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editingPlan, setEditingPlan] = useState<TravelPlan | null>(null);
  const [innerReloadKey, setInnerReloadKey] = useState(0);
  const [deletingPlan, setDeletingPlan] = useState<TravelPlan | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const handleCloseDelete = () => {
    if (isDeleting) return;

    setDeleteError(null);
    setDeletingPlan(null);
  }

  const handleConfirmDelete = async () => {
    if (!token || !deletingPlan) return;
    setIsDeleting(true);
    setDeleteError(null);

    try {
      await deleteTravelPlan(token, travelId, deletingPlan.id);
      setDeletingPlan(null);
      setInnerReloadKey((k) => k + 1);
      showFlash("プランを削除しました");
    } catch {
      setDeleteError("プランの削除に失敗しました");
    } finally {
      setIsDeleting(false);
    }
  }

  useEffect(() => {
    if (!token) return;

    const fetchPlans = async () => {
      try {
        const data = await getTravelPlans(token, travelId);
        setPlans(data);
      } catch {
        setError("予定の取得に失敗しました");
      } finally {
        setIsLoading(false);
      }
    };

    fetchPlans();
  }, [token, travelId, reloadKey, innerReloadKey]);

  if (isLoading) return <Typography>Loading...</Typography>;
  if (error) return <Typography color="error">{error}</Typography>;
  if (plans.length === 0) {
    return (
      <Typography color="text.secondary">まだ予定がありません</Typography>
    );
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
              <TravelPlanCard key={plan.id} plan={plan} onEdit={setEditingPlan} onDelete={setDeletingPlan} />
            ))}
          </Box>
        </Box>
      ))}
      {
        editingPlan && (
          <TravelPlanDialog
            open
            travelId={travelId}
            startDate={startDate}
            endDate={endDate}
            plan={editingPlan}
            onClose={() => setEditingPlan(null)}
            onSaved={() => {
              setInnerReloadKey((k) => k + 1);
              showFlash("プランを更新しました");
            }}
          />
        )
      }

      {deletingPlan && (
        <ConfirmDialog
          open
          heading="プランを削除しますか？"
          description={`「${deletingPlan.place}」のプランが削除されます。この操作は取り消せません。`}
          error={deleteError}
          busy={isDeleting}
          confirmLabel="削除する"
          busyLabel="削除中..."
          confirmColor="error"
          onClose={handleCloseDelete}
          onConfirm={handleConfirmDelete}
        />
      )}
    </Box >
  );
};
