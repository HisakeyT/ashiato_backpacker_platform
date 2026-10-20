import { useCallback, useEffect, useRef, useState } from "react";

import { listRecords } from "./repository";
import type { RecordIDs } from "./repository";
import { recordLabels } from "./types";
import type { PlanRecord } from "./types";

export const notifyRecordChanged = (planId: number) =>
  window.dispatchEvent(new CustomEvent("plan-record-changed", { detail: planId }));

// eslint-disable-next-line eslint/max-lines-per-function
export const useRecordData = ({
  ids,
  publicView,
  token,
  active,
}: {
  ids: RecordIDs;
  publicView: boolean;
  token: string | null;
  active: boolean;
}) => {
  const [records, setRecords] = useState<PlanRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const sequence = useRef(0);
  const { travelId, planId, kind } = ids;

  const load = useCallback(async () => {
    const current = ++sequence.current;
    setLoading(true);
    setError(null);
    try {
      const result = await listRecords({ ids: { travelId, planId, kind }, publicView, token });
      if (current === sequence.current) setRecords(result);
    } catch {
      if (current === sequence.current) setError(`${recordLabels[kind]}の取得に失敗しました`);
    } finally {
      if (current === sequence.current) setLoading(false);
    }
  }, [travelId, planId, kind, publicView, token]);

  useEffect(() => {
    const timer = active ? window.setTimeout(() => void load(), 0) : undefined;
    const requestSequence = sequence;
    return () => {
      if (timer !== undefined) window.clearTimeout(timer);
      requestSequence.current++;
    };
  }, [active, load]);

  useEffect(() => {
    if (!active) return;
    const onChanged = (event: Event) => {
      if ((event as CustomEvent<number>).detail === planId) void load();
    };
    window.addEventListener("plan-record-changed", onChanged);
    return () => window.removeEventListener("plan-record-changed", onChanged);
  }, [active, planId, load]);

  return { records, loading, error, load };
};
