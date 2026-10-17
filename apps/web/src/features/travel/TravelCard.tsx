import React from "react";
import { Card, CardActionArea, CardContent, Chip, Typography } from "@mui/material";
import { Link as RouterLink } from "react-router-dom";
import type { Travel } from "./types";

type Props = {
  travel: Travel;
};

export const TravelCard = ({ travel }: Props) => {
  return (
    <Card className="TravelCard" sx={{ height: "100%" }}>
      <CardActionArea
        component={RouterLink}
        to={`/travels/${travel.id}`}
        sx={{ height: "100%" }}
      >
        <CardContent>
          <Typography variant="h6" gutterBottom>
            {travel.title}
          </Typography>
          <Typography variant="body2" color="text.secondary" gutterBottom>
            {travel.startDate} 〜 {travel.endDate}
          </Typography>
          <Chip
            size="small"
            label={travel.isPublic ? "公開" : "非公開"}
            color={travel.isPublic ? "primary" : "default"}
          />
        </CardContent>
      </CardActionArea>
    </Card>
  );
};
