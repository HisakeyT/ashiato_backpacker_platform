import { Box, Container } from "@mui/material";
import { Outlet, useLocation } from "react-router-dom";

import { Footer } from "../components/layout/Footer";
import { Header } from "../components/layout/Header";

export function AppLayout() {

  const pathname = useLocation();
  const isFullScreen = pathname.pathname === "/";


  return (
    <Box
      className="AppLayout"
      sx={{
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        background: `
          radial-gradient(
            circle at 10% 20%,
            rgba(76, 104, 78, 0.08),
            transparent 30%
          ),
          radial-gradient(
            circle at 90% 80%,
            rgba(180, 140, 90, 0.08),
            transparent 30%
          ),
          #f7f4ec
        `,
      }}
    >
      <Header />

      {isFullScreen ? (
        <Box component="main" sx={{ flex: 1 }}>
          <Outlet />
        </Box>
      ) : (
        <Container
          className="AppLayout__content"
          component="main"
          maxWidth="lg"
          sx={{ flex: 1, py: 4 }}
        >
          <Outlet />
        </Container>
      )}

      <Footer />
    </Box>
  );
}
