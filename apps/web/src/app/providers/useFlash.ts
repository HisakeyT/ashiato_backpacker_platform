import { useContext } from "react";

import { FlashContext } from "./FlashContext";

export const useFlash = () => {
  const ctx = useContext(FlashContext);
  if (!ctx) throw new Error("useFlash must be used within FlashProvider");

  return ctx;
};
