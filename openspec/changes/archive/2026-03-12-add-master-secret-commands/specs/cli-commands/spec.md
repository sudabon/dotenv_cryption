## ADDED Requirements

### Requirement: create master コマンド
CLI SHALL provide an `envcrypt create master` command that generates a new 32-byte master key and stores it in the Secret Manager configured by `dotenv.yaml`.

#### Scenario: GCP でマスター鍵を作成
- **WHEN** `cloud: gcp` が設定された `dotenv.yaml` で `envcrypt create master` を実行する
- **THEN** `gcp.project_id` と `gcp.secret_id` のシークレットに新しい 32 バイト鍵が登録される

#### Scenario: AWS でマスター鍵を作成
- **WHEN** `cloud: aws` が設定された `dotenv.yaml` で `envcrypt create master` を実行する
- **THEN** `aws.region` と `aws.secret_id` のシークレットに新しい 32 バイト鍵が登録される

#### Scenario: 既存シークレットがある
- **WHEN** 同名のシークレットが既に存在する状態で `envcrypt create master` を実行する
- **THEN** シークレットを上書きせず、既に存在することを示すエラーで終了する

### Requirement: delete master コマンド
CLI SHALL provide an `envcrypt delete master` command that deletes the master key secret configured by `dotenv.yaml`.

#### Scenario: 設定済みシークレットを削除
- **WHEN** `dotenv.yaml` に設定されたシークレットが存在する状態で `envcrypt delete master` を実行する
- **THEN** 対応する Secret Manager からそのシークレットが削除される

#### Scenario: シークレットが存在しない
- **WHEN** 設定されたシークレットが存在しない状態で `envcrypt delete master` を実行する
- **THEN** シークレット名を含む not found エラーで終了する
