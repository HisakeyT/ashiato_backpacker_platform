# Ashiato（Backpacker Travel Sharing Platform）

Ashiato は、バックパッカーが旅を記録・共有し、ほかの人の旅を自分の旅行計画に活かせるサービスです。旅行ごとに滞在期間や場所をまとめ、公開された旅を参考にしながら、自分の旅を作れます。

## 現在の機能

- ニックネーム・メールアドレス・パスワードによるユーザー登録とログイン
- 自分の旅行の作成・一覧表示・編集・削除（タイトル、期間、公開状態）
- 旅行内のプランの作成・編集・削除（滞在期間、場所、内容、表示順）
- Nominatim を使った場所検索と、国・州／県・都市・座標の登録
- 公開旅行の一覧・詳細の閲覧（ログイン不要）
- 公開旅行を自分の非公開旅行としてコピーし、編集する機能

## 技術スタック

| 領域 | 技術 |
| --- | --- |
| Web | React、TypeScript、Vite、React Router、MUI |
| API | Go、Gin、GORM |
| データベース・マイグレーション | MySQL、goose |
| 認証 | JWT、bcrypt |
| 場所検索 | Nominatim（OpenStreetMap） |
| ローカル開発 | Docker Compose、Air |

## リポジトリ構成

Web と API を同じリポジトリで管理するモノレポです。ローカル開発では Docker Compose で Web・API・MySQL を起動し、Web からの `/api/*` を Vite が API に転送します。

```text
.
├── apps/
│   ├── web/              # React の SPA
│   └── api/              # Go の API と DB マイグレーション
├── compose.yml           # ローカル開発用のサービス定義
├── ARCHITECTURE.md       # 各アプリの設計とコードの配置
└── PRODUCT_BACKLOG.md    # 実装状況と開発順序
```

各アプリの起動・開発手順は [Web README](apps/web/README.md) と [API README](apps/api/README.md) を参照してください。設計の詳細は [アーキテクチャ](ARCHITECTURE.md) にまとめています。

## 今後の開発予定

1. 旅行・プランの期間表示と検証を整え、滞在場所を地図に表示する
2. 観光・宿泊・食事などの記録と、滞在場所を結ぶ移動を追加する
3. 費用の入力・集計、写真、旅行記を追加する
4. アカウント設定や公開旅行の検索を整え、MVP を完成させる

詳細な作業順序と進捗は [プロダクトバックログ](PRODUCT_BACKLOG.md) を参照してください。
