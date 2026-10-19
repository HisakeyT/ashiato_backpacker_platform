import { Alert, Box } from "@mui/material";
import L from "leaflet";
import { useEffect, useRef, useState } from "react";

import "leaflet/dist/leaflet.css";
import "./travelMap.css";
import { type MapPlan, groupMapPlans } from "./mapPlans";

/** 同じ座標の訪問番号を一つのピンに表示し、プラン名をピンの説明に設定する。 */
const addMarker = (map: L.Map, group: MapPlan[]) => {
  const label = group.map((item) => item.number).join("・");
  const marker = L.marker(group[0].position, {
    icon: L.divIcon({
      className: "travel-map-pin",
      html: `<span>${label}</span>`,
      iconSize: [44, 44],
      iconAnchor: [22, 22],
    }),
    title: group.map((item) => `${item.number}. ${item.plan.place}`).join("、"),
    // 地図上のピンは操作対象にせず、プランの詳細は下の一覧から確認できる。
    keyboard: false,
  }).addTo(map);

  const element = marker.getElement();
  element?.setAttribute("role", "img");
  element?.setAttribute("aria-label", marker.options.title ?? "訪問地点");
};

/** 背景タイル、ピン、訪問順の直線を追加し、全地点が入る表示範囲に合わせる。 */
const populateMap = (map: L.Map, located: MapPlan[], onTileError: () => void) => {
  L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution:
      '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
  })
    .on("tileerror", onTileError)
    .addTo(map);

  for (const group of groupMapPlans(located)) addMarker(map, group);

  const positions = located.map((item) => item.position);
  // 座標のないプランは線から除く。線は実際の移動経路ではなく地点間の直線。
  if (positions.length > 1)
    L.polyline(positions, { color: "#1565c0", weight: 3, dashArray: "8 6" }).addTo(map);

  map.fitBounds(L.latLngBounds(positions), { padding: [32, 32], maxZoom: 12 });
};

/** Leaflet の地図を DOM に作成し、再描画・破棄に合わせて関連リソースを管理する。 */
export default function TravelMapCanvas({ located }: { located: MapPlan[] }) {
  const container = useRef<HTMLDivElement>(null);
  const [tileError, setTileError] = useState(false);

  useEffect(() => {
    if (!container.current) return;
    const map = L.map(container.current, { scrollWheelZoom: false });

    // コンテナのサイズが変わったとき、Leaflet に表示領域を再計算させる。
    const observer = new ResizeObserver(() => map.invalidateSize());
    try {
      populateMap(map, located, () => setTileError(true));
      observer.observe(container.current);
    } catch (error) {
      // 初期化途中で失敗しても、作成済みの監視と地図を残さない。
      observer.disconnect();
      map.remove();
      throw error;
    }
    return () => {
      // 地点の更新やアンマウント時に、古い地図と監視を破棄する。
      observer.disconnect();
      map.remove();
    };
  }, [located]);

  return (
    <>
      {tileError && (
        <Alert severity="warning">
          背景地図を読み込めませんでした。プランの情報は下の一覧から確認できます。
        </Alert>
      )}
      <Box
        ref={container}
        component="section"
        aria-label="プランの訪問順の地図"
        sx={{ height: { xs: 320, sm: 420 }, zIndex: 0, borderRadius: 2 }}
      />
    </>
  );
}
