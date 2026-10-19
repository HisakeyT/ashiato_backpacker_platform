import { Box, Button, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { Link as RouterLink } from "react-router-dom";

import { useAuth } from "../../../app/providers/useAuth";
import { getTravels } from "../repository";
import { TravelCard } from "../TravelCard";
import type { Travel } from "../types";

const TravelListHeader = ({ nickname }: { nickname?: string }) => (
  <Box
    sx={{
      display: "flex",
      flexDirection: { xs: "column", sm: "row" },
      alignItems: { xs: "flex-start", sm: "center" },
      justifyContent: "space-between",
      gap: 2,
      mb: 3,
    }}
  >
    <Typography
      variant="h4"
      component="h1"
      sx={{ fontSize: { xs: "1.5rem", sm: "2.125rem" }, wordBreak: "break-word" }}
    >
      {nickname}の旅行記
    </Typography>
    <Button
      component={RouterLink}
      to="/travels/new"
      variant="contained"
      color="primary"
      sx={{ flexShrink: 0 }}
    >
      新しい旅行を作る
    </Button>
  </Box>
);

const TravelCards = ({ travels }: { travels: Travel[] }) =>
  travels.length === 0 ? (
    <Typography color="text.secondary">まだ旅の記録がありません。最初の一歩を記録しよう</Typography>
  ) : (
    <Box
      sx={{ display: "grid", gap: 2, gridTemplateColumns: "repeat(auto-fill, minmax(260px, 1fr))" }}
    >
      {travels.map((travel) => (
        <TravelCard key={travel.id} travel={travel} />
      ))}
    </Box>
  );

export const TravelListPage = () => {
  const { token, user } = useAuth();
  const [travels, setTravels] = useState<Travel[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!token) return;

    const fetchTravels = async () => {
      try {
        const data = await getTravels(token);
        setTravels(data);
      } catch {
        setError("旅行の取得に失敗しました");
      } finally {
        setIsLoading(false);
      }
    };

    fetchTravels();
  }, [token]);

  if (isLoading) return <p>Loading...</p>;
  if (error) return <p>{error}</p>;

  return (
    <Box className="TravelListPage" sx={{ p: 2 }}>
      <TravelListHeader nickname={user?.nickname} />
      <TravelCards travels={travels} />
    </Box>
  );
};
