import { AppBar, Toolbar, Typography } from "@mui/material";
import { Link, useNavigate } from "react-router-dom";

import { useAuth } from "../../app/providers/useAuth";
import { SERVICE_NAME } from "../../constants/site.ts";
import { DesktopNav, MobileNav } from "./HeaderNav";

export function Header() {
  const { isAuthenticated, user, logout } = useAuth();
  const navigate = useNavigate();
  const links = [
    { label: "みんなの旅行記", to: "/public/travels" },
    ...(isAuthenticated ? [{ label: "自分の旅行記", to: "/travels" }] : []),
  ];
  const handleLogout = () => {
    logout();
    navigate("/login");
  };
  const navProps = { links, isAuthenticated, nickname: user?.nickname, onLogout: handleLogout };

  return (
    <AppBar className="Header" position="sticky" elevation={0} sx={{ bgcolor: "#1e3a5f" }}>
      <Toolbar sx={{ minHeight: { xs: 56, md: 72 } }}>
        <Typography
          variant="h6"
          component={Link}
          to="/"
          sx={{
            fontSize: { xs: "1.5rem", md: "1.6rem" },
            fontWeight: 700,
            color: "inherit",
            textDecoration: "none",
            mr: { md: 4 },
          }}
        >
          {SERVICE_NAME}
        </Typography>
        <DesktopNav {...navProps} />
        <MobileNav {...navProps} />
      </Toolbar>
    </AppBar>
  );
}
