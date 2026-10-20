import { Alert, Box, Button } from "@mui/material";
import { useState } from "react";

import { ConfirmDialog } from "../../components/ConfirmDialog";

import { useFlash } from "../../app/providers/useFlash";
import { RecordForm } from "./RecordForm";
import { createRecord, deleteRecord, updateRecord } from "./repository";
import type { RecordIDs } from "./repository";
import type { PlanRecord, RecordInput } from "./types";
import { emptyRecordInput } from "./types";
import { notifyRecordChanged } from "./useRecordData";

// eslint-disable-next-line eslint/max-lines-per-function
export const RecordEditor = ({
  ids,
  token,
  records,
  load,
}: {
  ids: RecordIDs;
  token: string;
  records: PlanRecord[];
  load: () => Promise<void>;
}) => {
  const { showFlash } = useFlash();
  const [draft, setDraft] = useState<RecordInput>(emptyRecordInput);
  const [editing, setEditing] = useState<PlanRecord | null>(null);
  const [deleting, setDeleting] = useState<PlanRecord | null>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const resetDraft = () => {
    setDraft(emptyRecordInput());
    setEditing(null);
  };

  const save = async () => {
    setBusy(true);
    setError(null);
    try {
      if (editing) await updateRecord(token, { ids, recordId: editing.id, input: draft });
      else await createRecord(token, ids, draft);
      resetDraft();
      await load();
      notifyRecordChanged(ids.planId);
      showFlash(editing ? "記録を更新しました" : "記録を追加しました");
    } catch {
      setError("記録を保存できませんでした。入力内容を確認してください");
    } finally {
      setBusy(false);
    }
  };

  const remove = async () => {
    if (!deleting) return;
    setBusy(true);
    try {
      await deleteRecord(token, ids, deleting.id);
      setDeleting(null);
      await load();
      notifyRecordChanged(ids.planId);
      showFlash("記録を削除しました");
    } catch {
      setError("記録を削除できませんでした");
    } finally {
      setBusy(false);
    }
  };

  return (
    <>
      {error && <Alert severity="error">{error}</Alert>}
      {records.map((record) => (
        <Box key={record.id} sx={{ display: "flex", gap: 1 }}>
          <Button
            size="small"
            onClick={() => {
              setEditing(record);
              setDraft({
                name: record.name,
                locationText: record.locationText,
                mapUrl: record.mapUrl,
                mealType: record.mealType,
              });
            }}
          >
            「{record.name}」を編集
          </Button>
          <Button size="small" color="error" onClick={() => setDeleting(record)}>
            削除
          </Button>
        </Box>
      ))}
      <RecordForm
        kind={ids.kind}
        value={draft}
        onChange={setDraft}
        onSubmit={() => void save()}
        onCancel={resetDraft}
        busy={busy}
        editing={editing !== null}
      />
      <ConfirmDialog
        open={deleting !== null}
        heading="記録を削除しますか？"
        description={`「${deleting?.name ?? ""}」を削除します。`}
        busy={busy}
        confirmLabel="削除する"
        confirmColor="error"
        onClose={() => setDeleting(null)}
        onConfirm={() => void remove()}
      />
    </>
  );
};
