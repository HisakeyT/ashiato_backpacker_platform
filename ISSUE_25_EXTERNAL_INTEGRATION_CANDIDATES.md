# Issue #25: 外部連携の候補メモ

## 状態と参照

- 状態: MVP 後の詳細設計は保留。役に立ちそうな候補だけを整理した。実装する機能・提供元・順序は未決定。
- 対象: [#25 外部連携を追加](https://github.com/HisakeyT/ashiato_backpacker_platform/issues/25)
- 親: [#35 MVP 後](https://github.com/HisakeyT/ashiato_backpacker_platform/issues/35)
- 関連: [#2 地図](https://github.com/HisakeyT/ashiato_backpacker_platform/issues/2)、[#29 移動](https://github.com/HisakeyT/ashiato_backpacker_platform/issues/29)、[#30 費用](https://github.com/HisakeyT/ashiato_backpacker_platform/issues/30)

## 判断の前提

このサービスの中心は、自分の旅程・実際の移動・費用・写真・旅行記を残し、他の人の公開旅行から計画のヒントを得ること。外部連携は既存の記録や計画を補助するものとして考え、外部サービスが使えなくても旅行そのものを閲覧・編集できる形を優先する。

既存の場所検索には Nominatim を使用し、地図表示は別 Issue で先行する。したがって「地図を表示する」こと自体は #25 の新規候補に数えない。

## 役に立ちそうな候補

| 候補 | 利用者にとっての価値 | 既存データとの接点 | 今の位置付け |
| --- | --- | --- | --- |
| 地図アプリで経路を開く | 隣接 Plan 間や記録した移動の実際の道順・所要時間を、使い慣れた地図アプリで確認できる。 | Plan の座標・場所、移動の出発・到着。 | 有力。まず外部地図へのリンクを候補とし、アプリ内の経路計算は利用実態を見て検討する。保存した移動記録や #2/#11 の概略線を自動で書き換えない。 |
| 旅行日・場所の天気 | 近い日程の Plan に天気予報を添え、持ち物や屋外の予定を考えやすくする。過去の旅行には当時の気象データを参考表示できる。 | 日付のある Plan と座標。 | 有力。ただし日付未定・座標なしでは表示しない。遠い未来の予定を確定的な天気として見せず、過去データも本人が体験した天気とは区別する。 |
| 為替レートの参考表示 | 旅先の費用感を自国通貨などで把握しやすくする。 | 旅行の `currency_code` と費用。 | 条件付き。現在の #30 は旅行ごとに1通貨・換算なしの設計なので、複数通貨の記録や参考換算が必要になった段階で検討する。記録した原額と参考換算額を混同しない。 |
| 航空券の検索・比較 | 計画中の長距離移動で便や価格を探せる。 | 飛行機の移動、都市・空港、予定日。 | 優先度は低め。現行の移動は自分の旅程・実績の記録が主目的。予約・価格比較の需要が確認できてから、提供範囲と事業上の条件を調べる。 |

## 提供方法・サービスの例（未選定）

- 経路: [Google Maps URLs](https://developers.google.com/maps/documentation/urls/get-started) は出発地・目的地を渡して地図アプリを開け、API キーは不要。アプリ内で経路形状や所要時間を取得するなら、[openrouteservice Directions](https://openrouteservice.org/dev/) のような別の API を検討する。交通手段・地域の対応範囲を選定時に確認する。
- 天気: [Open-Meteo Forecast API](https://open-meteo.com/en/docs) は座標から近い日付の予報を取得でき、[Historical Weather API](https://open-meteo.com/en/docs/historical-weather-api) は過去の推定気象データを提供する。採用前に[利用条件](https://open-meteo.com/en/terms)、商用利用・帰属表示・呼び出し上限を再確認する。
- 為替: [Frankfurter](https://frankfurter.dev/) は現在・過去のレートを取得できる例。為替レートの参照日時・取得元と、利用者が入力した元の費用を分けて扱う必要がある。
- 航空券: [Amadeus Flight Offers Search](https://github.com/amadeus4dev/developer-guides/blob/master/docs/resources/flights.md) は候補例であり、導入可否・提供地域・利用条件・価格は着手時に確認する。

## 今は決めないこと

- どの候補を実装するか、順番、画面/API/DB 仕様、提供元、料金プラン。
- 天気・為替などの外部値を旅行データへ保存するか、表示時だけ取得するか。
- 航空券予約・決済、外部サービスとのアカウント連携。

## 既存連携で気を付けること

現在使う Nominatim 公開サーバーには[利用ポリシー](https://operations.osmfoundation.org/policies/nominatim/)がある。場所検索の利用が増える場合、リクエスト制限・アプリ識別・キャッシュなどを確認し、必要なら別の提供方法を検討する。これは #25 の新機能というより、既存の場所検索を安定して運用するための確認事項。
