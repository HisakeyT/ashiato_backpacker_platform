import { Box, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

import { PageBreadcrumbs } from "../../components/PageBreadcrumbs";

import { useAuth } from "../../app/providers/useAuth";
import { getTravel } from "./repository";
import { TravelDetailHeader } from "./TravelDetailHeader";
import { TravelPlansSection } from "./TravelPlansSection";
import type { TravelDetail } from "./types";

export const TravelDetailPage = () => {
  const { token } = useAuth();
  const { travelId } = useParams();
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

  if (isLoading) return <Typography>Loading...</Typography>;
  if (error || !travel) {
    return <Typography color="error">{error ?? "旅行が見つかりませんでした"}</Typography>;
  }
  return (
    <Box className="TravelDetailPage">
      <PageBreadcrumbs items={[{ label: "旅行一覧", to: "/travels" }, { label: travel.title }]} />
      <TravelDetailHeader travel={travel} />
      <TravelPlansSection key={travel.id} travel={travel} />
    </Box>
  );
};
