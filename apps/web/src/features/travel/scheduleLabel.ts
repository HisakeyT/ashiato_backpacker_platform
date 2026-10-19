import type { Travel } from "./types";

export const scheduleLabel = (mode: Travel["scheduleMode"]): string =>
  mode === "planning" ? "計画" : "記録";
