# Web

React・TypeScript・Vite で構成したフロントエンドです。

## 開発

以下のコマンドはリポジトリのルートで実行します。初回と `apps/web/package*.json` を変更した後は、web イメージをビルドしてください。

```sh
docker compose build web
docker compose up web
```

## コードの確認と修正

| コマンド | 用途 |
| --- | --- |
| `docker compose run --rm web npm run lint` | Oxlint で問題を確認する |
| `docker compose run --rm web npm run lint -- --fix` | Oxlint が安全に自動修正できる問題を修正する |
| `docker compose run --rm web npm run format` | Oxfmt で `src` と Vite の設定ファイルを整形する |
| `docker compose run --rm web npm run format:check` | 整形が必要なファイルがないか確認する |
| `docker compose run --rm web npm run build` | TypeScript の型チェックと本番ビルドを行う |

Oxlint の設定は [`.oxlintrc.json`](./.oxlintrc.json) にあります。正しさ・疑わしいコードに関するルールと、React、TypeScript、import、JSX のアクセシビリティに関するルールを有効にしています。`max-lines-per-function`、`max-statements`、`max-params`、`max-lines` は Oxlint の既定上限を使用します。

import 文の並び順は [`.oxfmtrc.json`](./.oxfmtrc.json) の `sortImports` で「外部ライブラリ → `src/components` の共通コンポーネント → その他のプロジェクト内 import」に統一し、`docker compose run --rm web npm run format:check` で確認します。Oxlint の `sort-imports` は、各 import 文の `{ ... }` 内の名前順を検査します。

通常の `--fix` で残る指摘は内容を確認して手動で修正してください。Oxlint には `--fix-suggestions` と `--fix-dangerously` もありますが、これらはプログラムの動作を変える可能性があります。必要な場合だけ変更差分を確認しながら使ってください。

```sh
docker compose run --rm web npm run lint -- --fix-suggestions
docker compose run --rm web npm run lint -- --fix-dangerously
```

修正後は `docker compose run --rm web npm run lint`、`docker compose run --rm web npm run format:check`、`docker compose run --rm web npm run build` を実行します。
