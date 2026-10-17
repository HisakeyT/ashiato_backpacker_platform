// src/components/ConfirmDialog.tsx
import { Alert, Button, DialogContentText } from "@mui/material";
import { BaseDialog } from "./BaseDialog";

type Props = {
  open: boolean;
  heading: string;
  description: string;
  error?: string | null;
  busy?: boolean;
  confirmLabel: string;
  busyLabel?: string; // 例: "削除中..."
  cancelLabel?: string;
  confirmColor?: "primary" | "error";
  onClose: () => void;
  onConfirm: () => void;
};

export const ConfirmDialog = ({
  open,
  heading,
  description,
  error,
  busy = false,
  confirmLabel,
  busyLabel,
  cancelLabel = "キャンセル",
  confirmColor = "primary",
  onClose,
  onConfirm,
}: Props) => (
  <BaseDialog
    className="ConfirmDialog"
    open={open}
    onClose={onClose}
    title={heading}
    busy={busy}
    actions={
      <>
        <Button color="inherit" onClick={onClose} disabled={busy}>
          {cancelLabel}
        </Button>
        <Button
          variant="contained"
          color={confirmColor}
          onClick={onConfirm}
          disabled={busy}
        >
          {busy && busyLabel ? busyLabel : confirmLabel}
        </Button>
      </>
    }
  >
    {error && (
      <Alert severity="error" sx={{ mb: 2 }}>
        {error}
      </Alert>
    )}
    <DialogContentText
      sx={{
        overflowWrap: "anywhere",
        wordBreak: "auto-phrase", // chrome だけの対応
        lineBreak: "strict",
        textWrap: "pretty",
      }}
    >
      {description}
    </DialogContentText>
  </BaseDialog >
);
