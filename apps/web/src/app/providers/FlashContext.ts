import type { AlertColor } from "@mui/material";
import { createContext } from "react";

export type FlashContextValue = {
  showFlash: (message: string, severity?: AlertColor) => void;
};

export const FlashContext = createContext<FlashContextValue | null>(null);
