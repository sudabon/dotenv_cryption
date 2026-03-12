## Context

既存の CLI は `encrypt` / `decrypt` のみを提供し、マスター鍵の登録は README の手動手順に依存している。実装上も Secret provider は `GetMasterKey()` だけを公開しており、Secret Manager に対する書き込み系操作は持っていない。

この変更では `dotenv.yaml` の `cloud` 設定を single source of truth とし、同じ provider 選択ロジックの上にマスター鍵の作成・削除を追加する。作成時の鍵サイズは暗号化処理で前提としている 32 バイトを維持する。

## Goals / Non-Goals

**Goals:**

- `envcrypt create master` で設定済みの Secret Manager に 32 バイトのマスター鍵を新規登録できる
- `envcrypt delete master` で設定済みのシークレットを削除できる
- 既存の `encrypt` / `decrypt` の provider 選択とエラーハンドリングの流れを極力再利用する
- 既存の認証ガイダンスに加え、存在しない/既に存在するシークレットを分かりやすく報告する

**Non-Goals:**

- 任意名のシークレット作成や `dotenv.yaml` を無視した上書き先指定
- 既存シークレットのローテーションや値更新
- 対話的な確認プロンプトやリカバリーワークフロー

## Decisions

### CLI 階層は `create master` / `delete master` にする

ユーザー要求どおり、ルート配下に `create` と `delete` の親コマンドを追加し、その配下に `master` サブコマンドをぶら下げる。これにより将来 `create config` や `delete cache` のような管理系コマンドを追加しやすい。

### SecretProvider インターフェースを拡張する

`create master` / `delete master` は `dotenv.yaml` に基づく provider 選択が必要であり、既存の factory をそのまま使える構成が最も小さい。別の admin provider を新設するよりも、`SecretProvider` に `CreateMasterKey()` と `DeleteMasterKey()` を追加する方が影響範囲が限定され、CLI 側の依存注入も維持できる。

### 生成するマスター鍵は 32 バイトの生バイト列とする

既存の暗号化エンジンは 32 バイト鍵を前提に検証しているため、作成コマンドでも同じ長さのランダムバイト列を生成する。GCP は secret payload としてそのまま version に格納し、AWS は `SecretBinary` で保存する。GCP の README で行っていた hex 文字列保存は 64 バイトになり既存実装と整合しないため、この変更では CLI 生成値を正とする。

### 削除は provider ごとに「即時利用不可」を優先する

GCP は secret リソースを削除する。AWS は recovery window 付き削除だと同名シークレットが長期間再作成できず運用上扱いづらいため、`ForceDeleteWithoutRecovery` を使って即時削除する。

## Risks / Trade-offs

- [誤削除] → `delete master` は破壊的操作なので README とコマンド説明で対象が `dotenv.yaml` の設定先であることを明示する
- [既存 README との不整合] → 手動作成手順を CLI 優先に更新し、必要ならクラウド CLI の代替手順として残す
- [クラウド API のエラー差異] → 既存の認証エラー分類を維持しつつ、already exists / not found を provider ごとに吸収する
