import { Alert, Box, Button, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

import { PageBreadcrumbs } from "../../../components/PageBreadcrumbs";

import { useAuth } from "../../../app/providers/useAuth";
import { getTravelPlans } from "../../travelPlan/repository";
import { TravelPlansSection } from "../../travelPlan/TravelPlansSection/TravelPlansSection";
import { getTravel, updateTravel } from "../repository";
import { TravelDetailHeader } from "../TravelDetailHeader";
import type { TravelDetail } from "../types";

const genericModeError = "状態の切り替えに失敗しました。もう一度お試しください。";

const recordingError = async (token: string, travelId: number): Promise<string> => {
  try {
    const plans = await getTravelPlans(token, travelId);
    const missing = plans.filter((plan) => !plan.startDate).map((plan) => plan.place);
    return missing.length
      ? `開始日が未入力のプラン: ${missing.join("、")}。下のプラン一覧から編集してください。`
      : genericModeError;
  } catch {
    return genericModeError;
  }
};

const useTravelDetail = (token: string | null, travelId: string | undefined) => {
  const [travel, setTravel] = useState<TravelDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    if (!token || !travelId) return;
    const fetchTravel = async () => {
      try {
        setTravel(await getTravel(token, Number(travelId)));
      } catch {
        setError("旅行が見つかりませんでした");
      } finally {
        setIsLoading(false);
      }
    };
    void fetchTravel();
  }, [token, travelId]);
  return { travel, setTravel, isLoading, error };
};

const ModeControl = ({
  token,
  travel,
  onUpdated,
}: {
  token: string | null;
  travel: TravelDetail;
  onUpdated: (travel: TravelDetail) => void;
}) => {
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const changeMode = async () => {
    if (!token) return;
    setBusy(true);
    setError(null);
    const nextMode = travel.scheduleMode === "planning" ? "recorded" : "planning";
    try {
      const updated = await updateTravel(token, travel.id, { scheduleMode: nextMode });
      onUpdated({ ...travel, scheduleMode: updated.scheduleMode });
    } catch {
      setError(nextMode === "recorded" ? await recordingError(token, travel.id) : genericModeError);
    } finally {
      setBusy(false);
    }
  };
  return (
    <Box sx={{ mt: 2 }}>
      <Button variant="outlined" disabled={busy} onClick={() => void changeMode()}>
        {travel.scheduleMode === "planning" ? "記録に切り替える" : "計画に戻す"}
      </Button>
      {error && (
        <Alert severity="error" sx={{ mt: 1 }}>
          {error}
        </Alert>
      )}
    </Box>
  );
};

export const TravelDetailPage = () => {
  const { token } = useAuth();
  const { travelId } = useParams();
  const { travel, setTravel, isLoading, error } = useTravelDetail(token, travelId);
  if (isLoading) return <Typography>Loading...</Typography>;
  if (error || !travel) {
    return <Typography color="error">{error ?? "旅行が見つかりませんでした"}</Typography>;
  }
  return (
    <Box className="TravelDetailPage">
      <PageBreadcrumbs items={[{ label: "旅行一覧", to: "/travels" }, { label: travel.title }]} />
      <TravelDetailHeader travel={travel} />
      <ModeControl token={token} travel={travel} onUpdated={setTravel} />
      <TravelPlansSection key={travel.id} travel={travel} />
    </Box>
  );
};
