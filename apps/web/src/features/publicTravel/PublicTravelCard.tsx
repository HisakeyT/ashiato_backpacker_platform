import { Card, CardActionArea, CardContent, Chip, Typography } from "@mui/material";
import { Link as RouterLink } from "react-router-dom";

import { scheduleLabel } from "../travel/scheduleLabel";
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
          <Chip size="small" label={scheduleLabel(publicTravel.scheduleMode)} sx={{ mb: 1 }} />
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
