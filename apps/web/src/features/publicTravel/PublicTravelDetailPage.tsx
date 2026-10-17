import { Box, Button, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { ConfirmDialog } from "../../components/ConfirmDialog";
import { PageBreadcrumbs } from "../../components/PageBreadcrumbs";

import { useAuth } from "../../app/providers/useAuth.ts";
import { useFlash } from "../../app/providers/useFlash.ts";
import { copyTravel, getPublicTravel } from "../travel/repository";
import type { PublicTravel } from "../travel/types";
import { getPublicTravelPlans } from "../travelPlan/repository";
import type { TravelPlan } from "../travelPlan/types";
import { PublicTravelPlanList } from "./PublicTravelPlanList";

const usePublicTravel = (travelId?: string) => {
  const [travel, setTravel] = useState<PublicTravel | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    if (!travelId) return;
    const fetchTravel = async () => {
      try {
        setTravel(await getPublicTravel(Number(travelId)));
      } catch {
        setError("旅行が見つかりませんでした");
      } finally {
        setIsLoading(false);
      }
    };
    void fetchTravel();
  }, [travelId]);
  return { travel, isLoading, error };
};

const usePublicPlans = (travelId?: string) => {
  const [plans, setPlans] = useState<TravelPlan[]>([]);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    if (!travelId) return;
    const fetchPlans = async () => {
      try {
        setPlans(await getPublicTravelPlans(Number(travelId)));
      } catch {
        setError("予定の取得に失敗しました");
      }
    };
    void fetchPlans();
  }, [travelId]);
  return { plans, error };
};

const PublicTravelCopy = ({ travel }: { travel: PublicTravel }) => {
  const navigate = useNavigate();
  const { token } = useAuth();
  const { showFlash } = useFlash();
  const [isCopying, setIsCopying] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const handleClick = () => {
    if (!token) {
      navigate("/login", { state: { from: `/public/travels/${travel.id}` } });
      return;
    }
    setError(null);
    setIsOpen(true);
  };
  const handleConfirm = async () => {
    if (!token) return;
    setIsCopying(true);
    try {
      const { id } = await copyTravel(token, travel.id);
      showFlash("旅行をコピーしました。日付を確認して、必要に応じて編集してください。");
      navigate(`/travels/${id}`);
    } catch {
      setError("コピーに失敗しました。もう一度お試しください");
      setIsCopying(false);
    }
  };
  return (
    <>
      <Button variant="contained" size="large" onClick={handleClick} sx={{ flexShrink: 0 }}>
        {token ? "自分の旅行にコピー" : "ログインしてコピー"}
      </Button>
      <ConfirmDialog
        open={isOpen}
        heading="この旅を自分の旅行にコピーしますか？"
        description={`「${travel.title}」のプランを、あなたの旅行としてコピーします。元の旅行は変わりません。コピーした旅行は非公開で、あとから自由に編集できます。日付は元のままなので、旅行の予定に合わせて直してください。`}
        error={error}
        busy={isCopying}
        confirmLabel="コピーする"
        busyLabel="コピー中..."
        onClose={() => setIsOpen(false)}
        onConfirm={handleConfirm}
      />
    </>
  );
};

const PublicTravelHeader = ({ travel }: { travel: PublicTravel }) => (
  <Box
    sx={{
      display: "flex",
      flexDirection: { xs: "column", sm: "row" },
      alignItems: { xs: "flex-start", sm: "flex-start" },
      justifyContent: "space-between",
      gap: 2,
    }}
  >
    <Box sx={{ minWidth: 0 }}>
      <Typography
        variant="h4"
        component="h1"
        sx={{ fontSize: { xs: "1.5rem", sm: "2.125rem" }, wordBreak: "break-word", mb: 1 }}
      >
        {travel.title}
      </Typography>
      <Typography color="text.secondary" gutterBottom>
        {travel.startDate} 〜 {travel.endDate}
      </Typography>
      <Typography variant="body2" color="text.secondary">
        by {travel.authorNickname}
      </Typography>
    </Box>
    <PublicTravelCopy travel={travel} />
  </Box>
);

export const PublicTravelDetailPage = () => {
  const { travelId } = useParams();
  const { travel, isLoading, error } = usePublicTravel(travelId);
  const { plans, error: plansError } = usePublicPlans(travelId);
  if (isLoading) return <Typography>Loading...</Typography>;
  if (error || !travel)
    return <Typography color="error">{error ?? "旅行が見つかりませんでした"}</Typography>;
  return (
    <Box className="PublicTravelDetailPage" sx={{ p: 2, maxWidth: 960, mx: "auto" }}>
      <PageBreadcrumbs
        items={[{ label: "みんなの旅行記", to: "/public/travels" }, { label: travel.title }]}
      />
      <PublicTravelHeader travel={travel} />
      <Typography variant="h5" component="h2" sx={{ mt: 4, mb: 2 }}>
        旅のプラン
      </Typography>
      {plansError ? (
        <Typography color="error">{plansError}</Typography>
      ) : (
        <PublicTravelPlanList plans={plans} />
      )}
    </Box>
  );
};
