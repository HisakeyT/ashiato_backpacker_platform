import { Alert, Snackbar } from "@mui/material";
import type { AlertColor } from "@mui/material";
import type { ReactNode, SyntheticEvent } from "react";

type Props = {
  open: boolean;
  onClose: () => void;
  children: ReactNode;
  severity?: AlertColor;
  autoHideDuration?: number;
};

export const FlashMessage = ({
  open,
  onClose,
  children,
  severity = "success",
  autoHideDuration = 5000,
}: Props) => {
  const handleClose = (_: SyntheticEvent | Event, reason?: string) => {
    if (reason === "clickaway") return;
    onClose();
  };

  return (
    <Snackbar
      open={open}
      autoHideDuration={autoHideDuration}
      onClose={handleClose}
      anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
    >
      <Alert severity={severity} variant="filled" onClose={onClose}>
        {children}
      </Alert>
    </Snackbar>
  );
};
