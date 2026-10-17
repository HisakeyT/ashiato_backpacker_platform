import { useCallback, useMemo, useState } from "react";
import type { ReactNode } from "react";
import type { AlertColor } from "@mui/material";
import { FlashMessage } from "../../components/FlashMessage";
import { FlashContext } from "./FlashContext";

type Flash = { message: string; severity: AlertColor };

export const FlashProvider = ({ children }: { children: ReactNode }) => {
  const [flash, setFlash] = useState<Flash | null>(null);
  const [open, setOpen] = useState(false);

  const showFlash = useCallback(
    (message: string, severity: AlertColor = "success") => {
      setFlash({ message, severity });
      setOpen(true);
    },
    [],
  );

  const value = useMemo(() => ({ showFlash }), [showFlash]);

  return (
    <FlashContext.Provider value={value}>
      {children}
      <FlashMessage
        open={open}
        severity={flash?.severity}
        onClose={() => setOpen(false)}
      >
        {flash?.message}
      </FlashMessage>
    </FlashContext.Provider>
  );
};
