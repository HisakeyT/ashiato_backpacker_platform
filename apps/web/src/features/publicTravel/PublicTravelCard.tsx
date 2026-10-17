import { Card, CardActionArea, CardContent, Typography } from "@mui/material";
import { Link as RouterLink } from "react-router-dom";

import type { PublicTravel } from "../travel/types";
import { formatPeriod } from "../travelPlan/format";

type Props = {
  publicTravel: PublicTravel;
};

export const PublicTravelCard = ({ publicTravel }: Props) => {
  return (
    <Card className="PublicTravelCard">
      <CardActionArea component={RouterLink} to={`/public/travels/${publicTravel.id}`}>
        <CardContent>
          <Typography variant="h6" component="h2" sx={{ wordBreak: "break-word" }}>
            {publicTravel.title}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {formatPeriod(publicTravel.startDate, publicTravel.endDate)}
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            by {publicTravel.authorNickname}
          </Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );
};
