import { Container } from "@mui/material";
import type { ReactNode } from "react";

type Props = {
  children: ReactNode;
};

export function PageContainer({ children }: Props) {
  return (
    <Container
      className="PageContainer"
      maxWidth="md"
      sx={{
        py: { xs: 4, sm: 6, md: 8 },
        px: { xs: 2, sm: 3 },
      }}
    >
      {children}
    </Container>
  );
}
