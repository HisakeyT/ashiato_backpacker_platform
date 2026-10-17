# アーキテクチャ

この文書は、現在の構成と、新しい機能を実装するときのコードの置き場所を示します。

## 全体構成

`apps/web` は React・TypeScript の SPA、`apps/api` は Go・Gin の API、データベースは MySQL です。開発環境では Docker Compose で起動し、Web の `/api/*` を Vite が API に転送します。場所の候補検索は Web から Nominatim に直接アクセスします。

### 技術選定

| 技術 | この構成での役割 |
| --- | --- |
| Go・Gin | API と認証ミドルウェアを実装する。 |
| GORM・MySQL・goose | データの読み書き、永続化、スキーマ変更をそれぞれ担当する。スキーマはマイグレーションで管理する。 |
| React・TypeScript・React Router・MUI | 画面、型、画面遷移、共通 UI を構成する。 |
| Vite・Docker Compose | Web の開発サーバーと、Web・API・DB のローカル実行環境を構成する。 |
| JWT・bcrypt | API の認証とパスワードハッシュの照合に使う。 |

## API (`apps/api`)

クリーンアーキテクチャの考え方に沿って、業務処理と HTTP・DB 操作を分けています。Usecase はデータ操作をインターフェース経由で利用します。既存の Model には GORM のマッピングタグもあります。

```text
Handler (HTTP) → Usecase (業務処理) → Model (業務上のデータ・ルール)
                      ↓
              Repository interface
                      ↑
              GORM implementation → MySQL
```

| ディレクトリ | 役割・新規実装の置き場所 |
| --- | --- |
| `cmd/server` | アプリの起動点。依存関係の組み立てとルート登録を行う。新しい Handler を公開するときはここで接続する。 |
| `internal/model` | エンティティ、値、業務上の検証・エラー。HTTP や DB の操作はここに持ち込まない。 |
| `internal/usecase` | 機能の処理手順、所有者確認などの業務ルール。必要なデータ操作は Repository のインターフェースを通す。 |
| `internal/repository` | Usecase が必要とするデータ操作のインターフェース。 |
| `internal/repository/gorm` | Repository の MySQL / GORM 実装。クエリ、関連データの読み込み、トランザクションを置く。 |
| `internal/handler` | リクエスト・レスポンスの DTO、入力の変換、HTTP ステータスへの変換。業務判断は Usecase に渡す。 |
| `internal/auth` | JWT の生成・検証と認証ミドルウェア。 |
| `internal/database` | DB 接続の設定。 |
| `db/migrations` | goose のスキーマ変更。テーブルや列を追加するときは新しいマイグレーションを作る。 |

新しい機能では、必要に応じて **Model → Repository インターフェース → GORM 実装 → Usecase → Handler → `main.go` の接続** の順に追加します。既存の機能に関わる変更は、対応するファイルに加えます。テスト用の Repository 実装は `internal/repository/mock` に置きます。

### データモデル

現在の永続化モデルは `users` → `travels` → `travel_plans` の親子関係です。旅行はユーザーに、プランは旅行に属します。

- `users`: ニックネーム、重複しないメールアドレス、パスワードハッシュを保持する。
- `travels`: 旅行の期間・タイトル・公開状態を保持する。コピーした旅行は `copied_from_travel_id` でコピー元を参照する。
- `travel_plans`: 旅行内の期間・場所・内容・表示順を保持する。終了日と位置情報は任意で、終了日がない場合は開始日と同日として扱う。

新しいデータを保存するときは、まず既存の親子関係に属するか、別のエンティティが必要かを決めてから、`db/migrations` と `internal/model` を追加・変更します。

## Web (`apps/web`)

画面と API 呼び出しは機能ごとに `src/features` にまとめます。`features` 直下には、旅行・旅行プラン・認証のように、ユーザーの利用目的と操作がまとまる単位でディレクトリを作ります。同じ機能の一覧・詳細・編集などは同じディレクトリに置き、公開旅行の閲覧と自分の旅行の管理のように利用目的が異なる場合は分けます。

| ディレクトリ | 役割・新規実装の置き場所 |
| --- | --- |
| `src/app` | ルート定義、認証による画面の保護、アプリ全体の Provider とレイアウト。画面を追加するときは `routes.tsx` に登録する。 |
| `src/features/<機能名>` | その機能の画面、機能内だけで使う UI 部品、型、API 呼び出し、補助処理。 |
| `src/components` | 複数の機能で使う UI 部品。 |
| `src/lib` | 機能に依存しない共通処理。`apiFetch.ts` は API 通信と JSON のキー変換を担当する。 |
| `src/constants` | サイト全体で共有する定数。 |
| `src/assets`、`public` | 画像などの静的ファイル。 |

機能ディレクトリ内では、画面を `*Page.tsx`、画面内で使う UI をコンポーネントの `.tsx` ファイル、API の入出力などの型を `types.ts`、API 呼び出しを `repository.ts` に置きます。整形・並べ替えなど、その機能に閉じた処理は用途に応じた名前の `.ts` ファイルに分けます。必要なファイルだけを作り、API 呼び出しには共通の `apiFetch` を使います。複数の機能で使う UI は `src/components` に置きます。
