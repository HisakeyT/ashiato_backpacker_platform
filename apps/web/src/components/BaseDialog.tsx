// src/components/BaseDialog.tsx
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
} from "@mui/material";
import type { DialogProps } from "@mui/material";
import type { ReactNode } from "react";

type Props = {
  className?: string;
  open: boolean;
  onClose: () => void;
  title: ReactNode;
  children: ReactNode;
  actions?: ReactNode;
  busy?: boolean; // 処理中は背景クリック・Esc で閉じない
  maxWidth?: DialogProps["maxWidth"];
};

export const BaseDialog = ({
  className,
  open,
  onClose,
  title,
  children,
  actions,
  busy = false,
  maxWidth = "xs",
}: Props) => (
  <Dialog
    className={className}
    open={open}
    onClose={busy ? undefined : onClose}
    fullWidth
    maxWidth={maxWidth}
  >
    <DialogTitle>{title}</DialogTitle>
    <DialogContent>{children}</DialogContent>
    {actions && <DialogActions>{actions}</DialogActions>}
  </Dialog>
);
