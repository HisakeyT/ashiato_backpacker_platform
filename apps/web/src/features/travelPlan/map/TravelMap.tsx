import { Alert, Box, Button, Typography } from "@mui/material";
import { Component, type ReactNode, Suspense, lazy, useId, useMemo, useState } from "react";

import type { TravelPlan } from "../types";
import { MapPlanDialog } from "./MapPlanDialog";
import { mapPlans } from "./mapPlans";

// Leaflet を使う Canvas は、座標のあるプランを実際に描画するときだけ読み込む。
const Canvas = lazy(() => import("./TravelMapCanvas"));

/** 地図の読み込みや描画に失敗したとき、プラン一覧へ誘導する表示に切り替える。 */
class MapErrorBoundary extends Component<{ children: ReactNode }, { failed: boolean }> {
  state = { failed: false };
  static getDerivedStateFromError() {
    return { failed: true };
  }
  render() {
    return this.state.failed ? (
      <Alert severity="warning">地図を表示できませんでした。下のプラン一覧をご利用ください。</Alert>
    ) : (
      this.props.children
    );
  }
}

const MissingPlans = ({ missing }: { missing: TravelPlan[] }) => (
  <>
    {missing.length > 0 && (
      <Alert severity="info" sx={{ mt: 2 }}>
        位置情報がない、または不正なプラン {missing.length}
        件は地図に表示していません。プラン一覧では引き続き確認できます。
        <Box component="details" sx={{ overflowWrap: "anywhere" }}>
          <summary>対象のプランを見る</summary>
          <ul>
            {missing.map((plan) => (
              <li key={plan.id}>{plan.place}</li>
            ))}
          </ul>
        </Box>
      </Alert>
    )}
  </>
);

const MapContent = ({
  located,
  missing,
  onSelect,
}: ReturnType<typeof mapPlans> & { onSelect: (ids: number[]) => void }) => (
  <>
    {located.length ? (
      <>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          番号は登録した訪問順です。線は地点を直線で結んだもので、実際の道路・交通経路ではありません。位置情報のないプランは飛ばして結びます。
        </Typography>
        <MapErrorBoundary>
          <Suspense fallback={<Typography component="output">地図を読み込んでいます…</Typography>}>
            <Canvas located={located} onSelect={onSelect} />
          </Suspense>
        </MapErrorBoundary>
        <Box
          sx={{ display: "flex", gap: 1, flexWrap: "wrap", mt: 1 }}
          aria-label="地図のプラン詳細"
        >
          {located.map(({ plan, number }) => (
            <Button
              key={plan.id}
              size="small"
              onClick={() => onSelect([plan.id])}
              sx={{ overflowWrap: "anywhere" }}
            >
              {number}. {plan.place}
            </Button>
          ))}
        </Box>
      </>
    ) : (
      <Alert severity="info">
        位置情報のあるプランはありません。プランに位置情報を登録すると地図に表示されます。
      </Alert>
    )}
    <MissingPlans missing={missing} />
  </>
);

export const TravelMap = ({ plans }: { plans: TravelPlan[] }) => {
  const { located, missing } = useMemo(() => mapPlans(plans), [plans]);
  const [selectedIds, setSelectedIds] = useState<number[]>([]);
  const selected = located.filter(({ plan }) => selectedIds.includes(plan.id));
  const titleId = useId();
  // ボタンの移動先は地図の直後。ここから続くページ内のプラン一覧へ進める。
  const listId = useId();
  return (
    <Box component="section" aria-label="旅の地図" sx={{ mt: 4 }}>
      <Typography variant="h5" component="h2" sx={{ mb: 1 }}>
        旅の地図
      </Typography>
      <Button component="a" href={`#${listId}`}>
        プラン一覧へ進む
      </Button>
      <MapContent located={located} missing={missing} onSelect={setSelectedIds} />
      <MapPlanDialog selected={selected} titleId={titleId} onClose={() => setSelectedIds([])} />
      <Box id={listId} tabIndex={-1} aria-label="プラン一覧の開始" />
    </Box>
  );
};
