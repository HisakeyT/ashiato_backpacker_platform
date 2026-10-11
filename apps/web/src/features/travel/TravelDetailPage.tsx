import { useEffect, useState } from "react";
import { useNavigate, Link as RouterLink, useParams } from "react-router-dom";
import { Link, Box, Button, Chip, Typography } from "@mui/material";
import { useAuth } from "../../app/providers/useAuth";
import { useFlash } from "../../app/providers/useFlash";
import { ConfirmDialog } from "../../components/ConfirmDialog";
import { TravelPlanList } from "../travelPlan/TravelPlanList";
import { TravelPlanDialog } from "../travelPlan/TravelPlanDialog";
import { getTravel, deleteTravel } from "./repository";
import type { TravelDetail } from "./types";

export const TravelDetailPage = () => {
  const { token } = useAuth();
  const { showFlash } = useFlash();
  const { travelId } = useParams();
  const [travel, setTravel] = useState<TravelDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const navigate = useNavigate();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const [isPlanDialogOpen, setIsPlanDialogOpen] = useState(false);

  useEffect(() => {
    if (!token || !travelId) return;

    const fetchTravel = async () => {
      try {
        const data = await getTravel(token, Number(travelId));
        setTravel(data);
      } catch {
        setError("旅行が見つかりませんでした");
      } finally {
        setIsLoading(false);
      }
    };

    fetchTravel();
  }, [token, travelId]);

  const [reloadKey, setReloadKey] = useState(0);

  const handlePlanSaved = (result: "created" | "updated") => {
    setReloadKey((k) => k + 1);
    showFlash(
      result === "created" ? "プランを追加しました" : "プランを更新しました",
    );
  };

  const handleOpenDialog = () => {
    setDeleteError(null);
    setIsDialogOpen(true);
  };

  const handleDelete = async () => {
    if (!token || !travel) return;
    setIsDeleting(true);
    setDeleteError(null);
    try {
      await deleteTravel(token, travel.id);
      showFlash("旅行を削除しました");
      navigate("/travels");
    } catch {
      setDeleteError("旅行の削除に失敗しました");
      setIsDeleting(false);
    }
  };

  if (isLoading) return <Typography>Loading...</Typography>;
  if (error || !travel) {
    return (
      <Typography color="error">
        {error ?? "旅行が見つかりませんでした"}
      </Typography>
    );
  }

  return (
    <Box className="TravelDetailPage">
      <Link
        component={RouterLink}
        to="/travels"
        underline="hover"
        sx={{ display: "inline-block", mb: 2 }}
      >
        ← 旅行一覧へ
      </Link>

      <Box
        className="TravelDetailPage__header"
        sx={{
          display: "flex",
          alignItems: "flex-start",
          justifyContent: "space-between",
          flexWrap: "wrap",
          gap: 2,
          mb: 1,
        }}
      >
        <Typography
          variant="h4"
          component="h1"
          sx={{
            fontSize: { xs: "1.5rem", sm: "2.125rem" },
            wordBreak: "break-word",
          }}
        >
          {travel.title}
        </Typography>

        <Box sx={{ display: "flex", gap: 1 }}>
          <Button
            component={RouterLink}
            to={`/travels/${travel.id}/edit`}
            variant="outlined"
            color="primary"
          >
            編集する
          </Button>
          <Button variant="outlined" color="error" onClick={handleOpenDialog}>
            削除する
          </Button>
        </Box>
      </Box>

      <Typography color="text.secondary" gutterBottom>
        {travel.startDate} 〜 {travel.endDate}
      </Typography>
      <Chip
        size="small"
        label={travel.isPublic ? "公開" : "非公開"}
        color={travel.isPublic ? "primary" : "default"}
      />

      {travel.copiedFrom && (
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1.5 }}>
          <Link
            component={RouterLink}
            to={`/public/travels/${travel.copiedFrom.travelId}`}
          >
            {travel.copiedFrom.authorNickname}さんの旅行
          </Link>
          を参考にしました
        </Typography>
      )}

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
        <Button variant="contained" onClick={() => setIsPlanDialogOpen(true)}>
          プランを追加
        </Button>
      </Box>
      <TravelPlanList
        travelId={travel.id}
        startDate={travel.startDate}
        endDate={travel.endDate}
        reloadKey={reloadKey}
      />

      <ConfirmDialog
        open={isDialogOpen}
        heading="旅行を削除しますか？"
        description={`「${travel.title}」と、その旅のプランがすべて削除されます。この操作は取り消せません。`}
        error={deleteError}
        busy={isDeleting}
        confirmLabel="削除する"
        busyLabel="削除中..."
        confirmColor="error"
        onClose={() => setIsDialogOpen(false)}
        onConfirm={handleDelete}
      />

      <TravelPlanDialog
        open={isPlanDialogOpen}
        travelId={travel.id}
        startDate={travel.startDate}
        endDate={travel.endDate}
        onClose={() => setIsPlanDialogOpen(false)}
        onSaved={handlePlanSaved}
      />
    </Box>
  );
};
