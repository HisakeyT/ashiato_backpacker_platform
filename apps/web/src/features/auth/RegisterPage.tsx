import { Alert, Box, Button, Link, TextField, Typography } from "@mui/material";
import { type FormEvent, useState } from "react";
import { Link as RouterLink, useNavigate } from "react-router-dom";

import { PageContainer } from "../../components/layout/PageContainer";

import { useFlash } from "../../app/providers/useFlash";
import { ApiError } from "../../lib/apiFetch";
import { register } from "./repository";

const registrationError = (error: unknown): string => {
  if (error instanceof ApiError && error.status === 409)
    return "このメールアドレスは既に使われています";
  if (error instanceof ApiError && error.status === 400) return "入力内容を確認してください";
  return "登録に失敗しました";
};

type RegistrationFieldsProps = {
  values: { nickname: string; email: string; password: string; passwordConfirmation: string };
  onChange: (field: keyof RegistrationFieldsProps["values"], value: string) => void;
};

const RegistrationFields = ({ values, onChange }: RegistrationFieldsProps) => (
  <>
    <TextField
      label="Nickname"
      value={values.nickname}
      onChange={(e) => onChange("nickname", e.target.value)}
      required
      fullWidth
    />
    <TextField
      label="Email"
      type="email"
      value={values.email}
      onChange={(e) => onChange("email", e.target.value)}
      required
      fullWidth
    />
    <TextField
      label="Password"
      type="password"
      value={values.password}
      onChange={(e) => onChange("password", e.target.value)}
      required
      fullWidth
    />
    <TextField
      label="Password (confirm)"
      type="password"
      value={values.passwordConfirmation}
      onChange={(e) => onChange("passwordConfirmation", e.target.value)}
      required
      fullWidth
    />
  </>
);

const LoginLink = () => (
  <Typography variant="body2" color="text.secondary" align="center">
    アカウントをお持ちの方は{" "}
    <Link component={RouterLink} to="/login">
      ログイン
    </Link>
  </Typography>
);

const validatePasswords = (
  values: RegistrationFieldsProps["values"],
  setError: (message: string) => void,
): boolean => {
  if (values.password === values.passwordConfirmation) return true;
  setError("パスワードが一致しません");
  return false;
};

type RegistrationFormProps = RegistrationFieldsProps & {
  error: string | null;
  isSubmitting: boolean;
  onSubmit: (event: FormEvent) => void;
};

const RegistrationForm = ({
  values,
  onChange,
  error,
  isSubmitting,
  onSubmit,
}: RegistrationFormProps) => (
  <PageContainer>
    <Box
      component="form"
      onSubmit={onSubmit}
      sx={{ maxWidth: 480, mx: "auto", display: "flex", flexDirection: "column", gap: 2 }}
    >
      <Typography variant="h4" component="h1">
        Register
      </Typography>
      {error && <Alert severity="error">{error}</Alert>}
      <RegistrationFields values={values} onChange={onChange} />
      <Button type="submit" variant="contained" fullWidth disabled={isSubmitting}>
        登録する
      </Button>
      <LoginLink />
    </Box>
  </PageContainer>
);

export function RegisterPage() {
  const { showFlash } = useFlash();
  const [values, setValues] = useState({
    nickname: "",
    email: "",
    password: "",
    passwordConfirmation: "",
  });
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const navigate = useNavigate();

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setError(null);

    if (!validatePasswords(values, setError)) return;

    setIsSubmitting(true);
    try {
      await register(values);
      showFlash("登録が完了しました。ログインしてください");
      navigate("/login");
    } catch (e) {
      setError(registrationError(e));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <RegistrationForm
      values={values}
      onChange={(field, value) => setValues((current) => ({ ...current, [field]: value }))}
      error={error}
      isSubmitting={isSubmitting}
      onSubmit={handleSubmit}
    />
  );
}
