import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Box, Typography } from "@mui/material";
import { useAuth } from "../../app/providers/useAuth";
import { useFlash } from "../../app/providers/useFlash";
import { PageBreadcrumbs } from "../../components/PageBreadcrumbs";
import { getTravel, updateTravel } from "./repository";
import { TravelForm } from "./TravelForm";
import type { Travel, UpdateTravelInput } from "./types";

export const TravelEditPage = () => {
  const { travelId } = useParams();
  const { token } = useAuth();
  const { showFlash } = useFlash();
  const navigate = useNavigate();

  const [travel, setTravel] = useState<Travel | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (!token || !travelId) return;

    const fetchTravel = async () => {
      try {
        const data = await getTravel(token, Number(travelId));
        setTravel(data);
      } catch {
        setLoadError("旅行の取得に失敗しました");
      } finally {
        setIsLoading(false);
      }
    };

    fetchTravel();
  }, [token, travelId]);

  const handleSubmit = async (input: UpdateTravelInput) => {
    if (!token || !travel) return;
    setIsSubmitting(true);
    setError(null);
    try {
      await updateTravel(token, travel.id, input);
      showFlash("旅行を更新しました");
      navigate(`/travels/${travel.id}`);
    } catch {
      setError("旅行の更新に失敗しました");
      setIsSubmitting(false);
    }
  };

  if (isLoading) return <p>Loading...</p>;
  if (loadError) return <p>{loadError}</p>;
  if (!travel) return <p>旅行が見つかりません</p>;

  return (
    <Box className="TravelEditPage">
      <PageBreadcrumbs
        items={[
          { label: "旅行一覧", to: "/travels" },
          { label: travel.title, to: `/travels/${travel.id}` },
          { label: "編集" },
        ]}
      />
      <Typography
        variant="h4"
        component="h1"
        sx={{
          mb: 3,
          fontSize: { xs: "1.5rem", sm: "2.125rem" },
          wordBreak: "break-word",
        }}
      >
        旅行を編集する
      </Typography>
      <TravelForm
        initialValue={{
          title: travel.title,
          startDate: travel.startDate,
          endDate: travel.endDate,
          isPublic: travel.isPublic,
        }}
        submitLabel="更新する"
        error={error}
        isSubmitting={isSubmitting}
        onSubmit={handleSubmit}
      />
    </Box>
  );
};
