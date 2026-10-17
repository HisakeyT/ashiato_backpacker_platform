import { Box, Button, Link, TextField, Typography } from "@mui/material";
import { useState } from "react";
import { Link as RouterLink, useLocation, useNavigate } from "react-router-dom";

import { PageContainer } from "../../components/layout/PageContainer";

import { useAuth } from "../../app/providers/useAuth";
import type { LoginLocationState } from "./types";

type LoginFormProps = {
  email: string;
  password: string;
  onEmailChange: (value: string) => void;
  onPasswordChange: (value: string) => void;
  onSubmit: () => void;
};

const LoginForm = ({
  email,
  password,
  onEmailChange,
  onPasswordChange,
  onSubmit,
}: LoginFormProps) => (
  <>
    <TextField
      label="Email"
      type="email"
      value={email}
      onChange={(event) => onEmailChange(event.target.value)}
      fullWidth
    />
    <TextField
      label="Password"
      type="password"
      value={password}
      onChange={(event) => onPasswordChange(event.target.value)}
      fullWidth
    />
    <Button variant="contained" type="button" fullWidth onClick={onSubmit}>
      Login
    </Button>
  </>
);

const RegistrationLink = () => (
  <Typography variant="body2" color="text.secondary" align="center">
    アカウントをお持ちでない方は{" "}
    <Link component={RouterLink} to="/register">
      新規登録
    </Link>
  </Typography>
);

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

        <LoginForm
          email={email}
          password={password}
          onEmailChange={setEmail}
          onPasswordChange={setPassword}
          onSubmit={handleSubmit}
        />

        <RegistrationLink />
      </Box>
    </PageContainer>
  );
}
