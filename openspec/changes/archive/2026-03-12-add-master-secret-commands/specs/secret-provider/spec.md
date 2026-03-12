## MODIFIED Requirements

### Requirement: Secret Provider インターフェース
Secret management SHALL be abstracted behind a common interface so providers are swappable. The interface SHALL expose `GetMasterKey() ([]byte, error)`, `CreateMasterKey() error`, and `DeleteMasterKey() error`.

#### Scenario: インターフェースの実装
- **WHEN** `dotenv.yaml` の `cloud` フィールドに基づいてプロバイダを初期化する
- **THEN** 対応する Secret Provider 実装が返される

### Requirement: GCP Secret Manager プロバイダ
The GCP provider SHALL retrieve secrets from the `projects/{project}/secrets/{secret}/versions/latest` path and SHALL manage the configured secret lifecycle in the same project.

#### Scenario: 正常なキー取得
- **WHEN** 有効な GCP 認証情報と正しい project_id / secret_id が設定されている
- **THEN** Secret Manager から32バイトのマスターキーを取得して返す

#### Scenario: 認証エラー
- **WHEN** GCP 認証情報が無効またはアプリケーションデフォルト認証が未設定
- **THEN** 認証エラーメッセージと `gcloud auth application-default login` の案内を表示する

#### Scenario: シークレットが存在しない
- **WHEN** 指定された secret_id が Secret Manager に存在しない
- **THEN** エラーメッセージでシークレット名を明示して終了する

#### Scenario: マスター鍵を作成
- **WHEN** 指定された secret_id が存在せず `CreateMasterKey()` を呼び出す
- **THEN** 32 バイトのランダム鍵を持つシークレットを作成する

#### Scenario: 既存シークレットの作成失敗
- **WHEN** 指定された secret_id が既に存在する状態で `CreateMasterKey()` を呼び出す
- **THEN** シークレット名を含む already exists エラーを返す

#### Scenario: マスター鍵を削除
- **WHEN** 指定された secret_id が存在する状態で `DeleteMasterKey()` を呼び出す
- **THEN** 対応するシークレットを削除する

### Requirement: AWS Secrets Manager プロバイダ
The AWS provider SHALL retrieve secrets from Secrets Manager using the configured `region` and `secret_id`, and SHALL manage the lifecycle of that configured secret.

#### Scenario: 正常なキー取得
- **WHEN** 有効な AWS 認証情報と正しい region / secret_id が設定されている
- **THEN** Secrets Manager から32バイトのマスターキーを取得して返す

#### Scenario: 認証エラー
- **WHEN** AWS 認証情報が無効または未設定
- **THEN** 認証エラーメッセージと AWS 認証設定の案内を表示する

#### Scenario: シークレットが存在しない
- **WHEN** 指定された secret_id が Secrets Manager に存在しない
- **THEN** エラーメッセージでシークレット名を明示して終了する

#### Scenario: マスター鍵を作成
- **WHEN** 指定された secret_id が存在せず `CreateMasterKey()` を呼び出す
- **THEN** 32 バイトのランダム鍵を持つシークレットを作成する

#### Scenario: 既存シークレットの作成失敗
- **WHEN** 指定された secret_id が既に存在する状態で `CreateMasterKey()` を呼び出す
- **THEN** シークレット名を含む already exists エラーを返す

#### Scenario: マスター鍵を削除
- **WHEN** 指定された secret_id が存在する状態で `DeleteMasterKey()` を呼び出す
- **THEN** 対応するシークレットを削除する
