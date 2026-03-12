## Why

`envcrypt` の利用開始時に、利用者が別途 `gcloud` / `aws` CLI でマスター鍵を作成・削除する必要があり、セットアップと運用が分断されている。`dotenv.yaml` の設定先に対して `envcrypt` 自体がマスター鍵を作成・削除できれば、初期導入と鍵のライフサイクル管理を一貫した操作で行える。

## What Changes

- `envcrypt create master` コマンドを追加し、`dotenv.yaml` で設定された GCP Secret Manager または AWS Secrets Manager に 32 バイトの新しいマスター鍵を登録する
- `envcrypt delete master` コマンドを追加し、`dotenv.yaml` で設定されたマスター鍵シークレットを削除する
- Secret provider にマスター鍵の作成・削除操作を追加し、クラウドごとのエラーを CLI で扱えるようにする
- README のセットアップ手順を CLI ベースの運用に更新する

## Capabilities

### New Capabilities

### Modified Capabilities

- `cli-commands`: `create master` と `delete master` の管理コマンドを追加する
- `secret-provider`: Secret provider インターフェースと GCP/AWS 実装にマスター鍵の作成・削除を追加する

## Impact

- `cmd/` の Cobra コマンド構成
- `internal/provider/` のインターフェースと GCP/AWS 実装
- provider と CLI のテスト
- `README.md` のセットアップ/運用手順
