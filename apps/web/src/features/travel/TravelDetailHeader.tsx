import { Box, Button, Chip, Link, Typography } from "@mui/material";
import { useState } from "react";
import { Link as RouterLink, useNavigate } from "react-router-dom";

import { ConfirmDialog } from "../../components/ConfirmDialog";

import { useAuth } from "../../app/providers/useAuth";
import { useFlash } from "../../app/providers/useFlash";
import { deleteTravel } from "./repository";
import type { TravelDetail } from "./types";

const DeleteTravelButton = ({ travel }: { travel: TravelDetail }) => {
  const { token } = useAuth();
  const { showFlash } = useFlash();
  const navigate = useNavigate();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const handleDelete = async () => {
    if (!token) return;
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
  return (
    <>
      <Button
        variant="outlined"
        color="error"
        onClick={() => {
          setDeleteError(null);
          setIsDialogOpen(true);
        }}
      >
        削除する
      </Button>
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
    </>
  );
};

const TravelHeading = ({ travel }: { travel: TravelDetail }) => (
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
      sx={{ fontSize: { xs: "1.5rem", sm: "2.125rem" }, wordBreak: "break-word" }}
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
      <DeleteTravelButton travel={travel} />
    </Box>
  </Box>
);

export const TravelDetailHeader = ({ travel }: { travel: TravelDetail }) => (
  <>
    <TravelHeading travel={travel} />
    <Typography color="text.secondary" gutterBottom>
      {travel.startDate} 〜 {travel.endDate}
    </Typography>
    <Chip
      sx={{ mr: 1 }}
      size="small"
      label={travel.scheduleMode === "planning" ? "旅行前の計画" : "旅行後の記録"}
    />
    <Chip
      size="small"
      label={travel.isPublic ? "公開" : "非公開"}
      color={travel.isPublic ? "primary" : "default"}
    />
    {travel.copiedFrom && (
      <Typography variant="body2" color="text.secondary" sx={{ mt: 1.5 }}>
        <Link component={RouterLink} to={`/public/travels/${travel.copiedFrom.travelId}`}>
          {travel.copiedFrom.authorNickname}さんの旅行
        </Link>
        を参考にしました
      </Typography>
    )}
  </>
);
