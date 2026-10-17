import { useState } from "react";
import { Link as RouterLink, useNavigate, useLocation } from "react-router-dom";
import { Link, Box, Button, TextField, Typography } from "@mui/material";
import { PageContainer } from "../../components/layout/PageContainer";
import { useAuth } from "../../app/providers/useAuth";
import type { LoginLocationState } from "./types";

export function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const state = location.state as LoginLocationState;

  const handleSubmit = async () => {
    try {
      await login(email, password);

      navigate(state?.from || "/", { replace: true });
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <PageContainer>
      <Box
        sx={{
          maxWidth: 480,
          mx: "auto",
          display: "flex",
          flexDirection: "column",
          gap: 2,
        }}
      >
        <Typography variant="h4" component="h1">
          Login
        </Typography>

        <TextField
          label="Email"
          type="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          fullWidth
        />

        <TextField
          label="Password"
          type="password"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          fullWidth
        />

        <Button
          variant="contained"
          type="button"
          fullWidth
          onClick={handleSubmit}
        >
          Login
        </Button>

        <Typography variant="body2" color="text.secondary" align="center">
          アカウントをお持ちでない方は{" "}
          <Link component={RouterLink} to="/register">
            新規登録
          </Link>
        </Typography>
      </Box>
    </PageContainer>
  );
}
