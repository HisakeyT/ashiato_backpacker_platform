import { useState } from "react";
import { Box, Button, IconButton, Menu, MenuItem, Typography } from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
import { Link } from "react-router-dom";

type NavLink = { label: string; to: string };

type HeaderNavProps = {
  links: NavLink[];
  isAuthenticated: boolean;
  nickname?: string;
  onLogout: () => void;
};

const navButtonSx = { fontSize: "1rem", fontWeight: 500, color: "inherit" };

export const DesktopNav = ({ links, isAuthenticated, nickname, onLogout }: HeaderNavProps) => (
  <>
    <Box sx={{ display: { xs: "none", md: "flex" }, flexGrow: 1, gap: 1 }}>
      {links.map((link) => (
        <Button key={link.to} component={Link} to={link.to} sx={navButtonSx}>
          {link.label}
        </Button>
      ))}
    </Box>
    <Box sx={{ display: { xs: "none", md: "flex" }, alignItems: "center" }}>
      {isAuthenticated ? (
        <>
          <Typography sx={{ mr: 2, fontSize: "1rem" }}>Hello, {nickname}</Typography>
          <Button onClick={onLogout} sx={navButtonSx}>
            Logout
          </Button>
        </>
      ) : (
        <Button component={Link} to="/login" sx={navButtonSx}>
          ログイン
        </Button>
      )}
    </Box>
  </>
);

export const MobileNav = ({ links, isAuthenticated, onLogout }: HeaderNavProps) => {
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
  const closeMenu = () => setAnchorEl(null);
  const handleLogout = () => {
    closeMenu();
    onLogout();
  };

  return (
    <Box sx={{ display: { xs: "flex", md: "none" }, ml: "auto" }}>
      <IconButton
        color="inherit"
        aria-label="メニューを開く"
        onClick={(event) => setAnchorEl(event.currentTarget)}
      >
        <MenuIcon />
      </IconButton>
      <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={closeMenu}>
        {links.map((link) => (
          <MenuItem key={link.to} component={Link} to={link.to} sx={{ py: 2 }} onClick={closeMenu}>
            {link.label}
          </MenuItem>
        ))}
        {isAuthenticated ? (
          <MenuItem sx={{ py: 2 }} onClick={handleLogout}>
            ログアウト
          </MenuItem>
        ) : (
          <MenuItem component={Link} to="/login" sx={{ py: 2 }} onClick={closeMenu}>
            ログイン
          </MenuItem>
        )}
      </Menu>
    </Box>
  );
};
