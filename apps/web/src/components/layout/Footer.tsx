import { Box, Typography } from "@mui/material";
import { SERVICE_NAME } from "../../constants/site.ts";

export function Footer() {
  return (
    <Box
      className="Footer"
      component="footer"
      sx={{
        py: 3,
        textAlign: "center",
      }}
    >
      <Typography variant="body2" color="text.secondary">
        © {SERVICE_NAME}
      </Typography>
    </Box>
  );
}
