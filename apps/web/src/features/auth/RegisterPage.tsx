import { type FormEvent, useState } from "react";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import { Link, Alert, Box, Button, TextField, Typography } from "@mui/material";
import { useFlash } from "../../app/providers/useFlash";
import { ApiError } from "../../lib/apiFetch";
import { PageContainer } from "../../components/layout/PageContainer";
import { register } from "./repository";

export function RegisterPage() {
  const { showFlash } = useFlash();
  const [nickname, setNickname] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordConfirmation, setPasswordConfirmation] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const navigate = useNavigate();

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setError(null);

    if (password !== passwordConfirmation) {
      setError("パスワードが一致しません");
      return;
    }

    setIsSubmitting(true);
    try {
      await register({ nickname, email, password, passwordConfirmation });
      showFlash("登録が完了しました。ログインしてください");
      navigate("/login");
    } catch (e) {
      if (e instanceof ApiError && e.status === 409) {
        setError("このメールアドレスは既に使われています");
      } else if (e instanceof ApiError && e.status === 400) {
        setError("入力内容を確認してください");
      } else {
        setError("登録に失敗しました");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <PageContainer>
      <Box
        component="form"
        onSubmit={handleSubmit}
        sx={{
          maxWidth: 480,
          mx: "auto",
          display: "flex",
          flexDirection: "column",
          gap: 2,
        }}
      >
        <Typography variant="h4" component="h1">
          Register
        </Typography>

        {error && <Alert severity="error">{error}</Alert>}

        <TextField
          label="Nickname"
          value={nickname}
          onChange={(e) => setNickname(e.target.value)}
          required
          fullWidth
        />
        <TextField
          label="Email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          fullWidth
        />
        <TextField
          label="Password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          fullWidth
        />
        <TextField
          label="Password (confirm)"
          type="password"
          value={passwordConfirmation}
          onChange={(e) => setPasswordConfirmation(e.target.value)}
          required
          fullWidth
        />

        <Button
          type="submit"
          variant="contained"
          fullWidth
          disabled={isSubmitting}
        >
          登録する
        </Button>

        <Typography variant="body2" color="text.secondary" align="center">
          アカウントをお持ちの方は{" "}
          <Link component={RouterLink} to="/login">
            ログイン
          </Link>
        </Typography>
      </Box>
    </PageContainer>
  );
}
