import { apiFetch } from "../../lib/apiFetch";
import type { PlanRecord, RecordInput, RecordKind } from "./types";

export type RecordIDs = { travelId: number; planId: number; kind: RecordKind };

const path = ({ travelId, planId, kind }: RecordIDs) =>
  `/api/travels/${travelId}/plans/${planId}/records/${kind}`;

export const listRecords = ({
  ids,
  publicView,
  token,
}: {
  ids: RecordIDs;
  publicView: boolean;
  token?: string | null;
}): Promise<PlanRecord[]> =>
  apiFetch<PlanRecord[]>(
    publicView
      ? `/api/travels/public/${ids.travelId}/plans/${ids.planId}/records/${ids.kind}`
      : path(ids),
    { method: "GET", token: token ?? undefined },
  );

export const createRecord = (token: string, ids: RecordIDs, input: RecordInput) =>
  apiFetch<PlanRecord>(path(ids), { method: "POST", token, body: input });

export const updateRecord = (
  token: string,
  target: { ids: RecordIDs; recordId: number; input: RecordInput },
) =>
  apiFetch<PlanRecord>(`${path(target.ids)}/${target.recordId}`, {
    method: "PATCH",
    token,
    body: target.input,
  });

export const deleteRecord = (token: string, ids: RecordIDs, recordId: number) =>
  apiFetch<void>(`${path(ids)}/${recordId}`, { method: "DELETE", token });
