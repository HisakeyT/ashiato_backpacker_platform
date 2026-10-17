# API

Go・Gin・GORM・MySQL で構成したバックエンドです。

## 開発

以下のコマンドはリポジトリのルートで実行します。先にルートの `.env` で Compose が参照する `DB_*`、`MYSQL_*`、`JWT_SECRET` を設定してください。API から MySQL への接続には `DB_HOST=mysql`、`DB_PORT=3306` を使用します。`DB_USER`・`DB_PASSWORD`・`DB_NAME` は、それぞれ MySQL 側のユーザー・パスワード・データベース名と一致させます。

初回と `apps/api/go.mod`・`apps/api/go.sum`・`apps/api/Dockerfile` を変更した後は、api イメージをビルドしてください。

```sh
docker compose build api
docker compose up -d api
```

起動後、マイグレーションを適用します。

```sh
docker compose exec api sh -c 'goose -dir /app/db/migrations mysql "$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true" up'
```

API は `http://localhost:8080` で待ち受けます。`/health` で起動状態を確認できます。開発中は Air が Go ファイルの変更を検知して再ビルドします。

## コードの確認と修正

以下のコマンドは `apps/api` で実行します。最初に `mise install` で Go と golangci-lint を用意してください。

| コマンド | 用途 |
| --- | --- |
| `mise exec -- go test ./...` | テストを実行する |
| `mise exec -- golangci-lint run ./...` | lint で問題を確認する |
| `mise exec -- golangci-lint fmt ./...` | Go コードを整形する |
| `mise exec -- golangci-lint fmt --diff ./...` | 整形が必要な箇所を確認する |

golangci-lint の設定は [`.golangci.yml`](./.golangci.yml) にあります。関数長は `funlen` の既定値（60 行・40 文）、行長は 140 文字、ファイル長はコメントと空行を除いて 400 行を上限としています。`funlen` はテストファイルを対象外にしています。整形には `gci` と `gofmt` を使用します。

修正後はテスト、lint、整形差分の確認を実行してください。DB スキーマを変更するときは `db/migrations` に新しい goose マイグレーションを追加します。
