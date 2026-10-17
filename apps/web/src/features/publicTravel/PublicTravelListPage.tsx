import { useEffect, useState } from "react";
import { Box, Typography } from "@mui/material";
import { getPublicTravels } from "../travel/repository";
import type { PublicTravel } from "../travel/types";
import { PublicTravelCard } from "./PublicTravelCard";

export const PublicTravelListPage = () => {
  const [travels, setTravels] = useState<PublicTravel[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPublicTravels = async () => {
      try {
        const data = await getPublicTravels();
        setTravels(data);
      } catch {
        setError("旅行の取得に失敗しました");
      } finally {
        setIsLoading(false);
      }
    };

    fetchPublicTravels();
  }, []);

  if (isLoading) return <p>Loading...</p>;
  if (error) return <p>{error}</p>;

  return (
    <Box className="PublicTravelListPage" sx={{ p: 2 }}>
      <Typography
        variant="h4"
        component="h1"
        gutterBottom
        sx={{ fontSize: { xs: "1.5rem", sm: "2.125rem" }, wordBreak: "break-word" }}
      >
        みんなの旅行記
      </Typography>

      {travels.length === 0 ? (
        <Typography color="text.secondary">
          公開されている旅行はまだありません
        </Typography>
      ) : (
        <Box
          sx={{
            display: "grid",
            gap: 2,
            gridTemplateColumns: "repeat(auto-fill, minmax(260px, 1fr))",
          }}
        >
          {travels.map((travel) => (
            <PublicTravelCard key={travel.id} publicTravel={travel} />
          ))}
        </Box>
      )}
    </Box>
  );
};
