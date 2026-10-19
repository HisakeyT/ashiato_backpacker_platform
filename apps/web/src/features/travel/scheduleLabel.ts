import type { Travel } from "./types";

export const scheduleLabel = (mode: Travel["scheduleMode"]): string =>
  mode === "planning" ? "旅行前の計画" : "旅行後の記録";
