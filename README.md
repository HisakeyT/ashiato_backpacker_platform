# Backpacker Travel Sharing Platform(Ashiato)

バックパッカー向けの旅行共有プラットフォーム。

自分の旅を記録して公開し、他のバックパッカーの旅を参考にできます。公開された旅はコピーして、自分の旅行計画として編集できます。旅行前の計画にも、旅行後の振り返りにも使えます。

## 主な機能

- ユーザー登録・ログイン・プロフィール
- 旅の作成・編集
  - 滞在(どの国・都市・街に、いつからいつまでいたか)
  - 宿泊・食事・観光・その他・移動の記録
  - 費用の記録と合計表示
  - 旅行記
- 旅の公開・非公開
- 公開された旅の閲覧
- 他人の旅のコピー・編集(費用は引き継がない)

## 技術スタック

| 分類 | 技術 |
| --- | --- |
| Backend | Go / Gin / GORM / MySQL |
| 認証 | JWT / bcrypt |
| Migration | goose |
| Frontend | React / TypeScript / Vite |
| 開発環境 | Docker / Docker Compose / Air |
| その他 | GitHub |

## 構成

Backend は Handler → UseCase → Repository の層構造で、API サーバーと SPA を分離した構成です。

## 設計

### データモデル

```
user
 └ travel(一つの旅行)
    ├ travel_plan(どこにいたか: 期間 date〜end_date)
    │   ├ stays          宿泊
    │   ├ meals          食事(朝昼晩は meal_type 列)
    │   ├ sightseeings   観光
    │   └ others         その他
    └ transfers(plan と plan を結ぶ移動)
```

- `travel_plan` は「予定」ではなく「どこにいたか」の単位。期間は `date`〜`end_date`(NULL は同日)で表す
- 宿泊・食事・観光・その他は plan の子テーブル(`travel_plan_id` が外部キー)
- 移動は plan 同士を結ぶ `transfers`(出発・到着の plan を持つ。NULL は旅行の外)
- 各テーブルは `cost`(NULL 可)を持ち、合計は全テーブルを合算して出す

### 設計判断

- plan に `type` は持たせない(plan と宿泊・食事は粒度が違うため)
- 種別ごとにテーブルを分ける(固有の列が後から増えても影響しないようにするため)
- 予実管理はしない(plan を編集すればそれが実績になる)
- `CopyTravel` は `cost` を引き継がない。plan の子と transfers はコピーし、transfers は新旧 plan id の対応表で付け替える

### 開発の進め方

1機能ごとに DB → Go(model → repository → usecase → handler → response、テスト)→ フロントの順で進め、各ステップで build / test と commit を挟みます。DB の変更は goose の新規マイグレーションで行い、既存のマイグレーションは編集しません。

### 今後の予定

`end_date` の追加 → 観光 → 宿泊・食事・その他 → 移動 → 費用の合計表示 → まとめ表示 → 写真・地図、の順に進めます。

## 開発

起動手順(Docker Compose、マイグレーション、環境変数)をここに書く。
