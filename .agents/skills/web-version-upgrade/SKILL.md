---
name: web-version-upgrade
description: apps/web の npm 依存関係と Node.js 実行環境を定期的に更新するときに使う。人が実行できるコマンドと、更新時の確認事項をまとめる。
---

# Web の定期バージョンアップ

以下はリポジトリのルートで実行する。対象は `apps/web`。Node.js と npm が必要。ローカルで実行できない場合の Docker コマンドは末尾に記載する。

## 1. 更新候補を確認

```sh
git status --short
node --version
npm --version
npm --prefix apps/web outdated --long
```

`npm outdated` の `Current` はインストール済み、`Wanted` は現在の指定範囲内の最新版、`Latest` は npm の最新リリース。更新候補がある場合の終了コード 1 は正常。通信エラーの場合はエラー内容を確認して再実行する。既存の作業中の変更は上書きしない。

必要に応じて `npm --prefix apps/web audit` で脆弱性も確認する。

## 2. 少数の関連パッケージを更新

更新するパッケージ名を選び、次のコマンドのパッケージ名を置き換える。通常の依存は `dependencies`、開発用の依存は `devDependencies` に残す。`npm install` は `package.json` と `package-lock.json` を更新する。

```sh
# dependencies の例: React と React DOM を同時に更新
npm install --prefix apps/web react@latest react-dom@latest

# devDependencies の例: 対応する型定義を更新
npm install --prefix apps/web --save-dev @types/react@latest @types/react-dom@latest
```

上の例を毎回すべて実行する必要はない。React と型定義、MUI と Emotion、Vite とプラグインなど、関連するものを一組ずつ更新する。指定範囲内の更新だけなら `npm update --prefix apps/web <package>` を使える。この場合、`package-lock.json` だけが変わることがある。

TypeScript は現在 `~` 指定なので、更新時も次のように維持する。

```sh
npm install --prefix apps/web --save-dev --save-prefix='~' typescript@latest
```

メジャーバージョンが変わる場合は、対象パッケージの公式変更履歴・移行ガイドと Node.js の対応バージョンを確認してから適用する。必要なコード修正を行い、次の検証を通してから別の組に進む。`apps/web/Dockerfile` の Node.js イメージも対応要件とサポート状況を確認し、必要な場合だけ更新する。

## 3. 再インストールして検証

```sh
npm ci --prefix apps/web
sh .agents/skills/web-version-upgrade/scripts/verify.sh
```

検証スクリプトは lint、整形チェック、TypeScript を含む本番ビルド、直接依存の整合性確認を順番に実行する。失敗したら更新した組の移行対応を行い、再実行する。

## 4. 差分を確認

```sh
git diff -- apps/web/package.json apps/web/package-lock.json apps/web/Dockerfile apps/web/src
git status --short
```

更新したパッケージ、メジャー更新で対応した内容、検証結果、残った問題を記録する。更新と無関係の変更を混ぜない。

## ローカルの Node.js で検証できない場合

Docker が使えるなら、Web README と同じ環境で実行する。`apps/web/package*.json` または Dockerfile を変更した後はイメージをビルドし直す。

```sh
docker compose build web
docker compose run --rm web npm ci
docker compose run --rm web npm run lint
docker compose run --rm web npm run format:check
docker compose run --rm web npm run build
docker compose run --rm web npm ls --depth=0
```
