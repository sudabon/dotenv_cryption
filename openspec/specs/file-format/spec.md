# file-format Specification

## Purpose
TBD - created by archiving change envcrypt-cli. Update Purpose after archive.
## Requirements
### Requirement: ENVC バイナリフォーマットの書き込み
Encrypted files SHALL be written with the following binary structure:
- MAGIC: 4バイト (`ENVC`)
- VERSION: 1バイト (現行: `0x01`)
- NONCE_LEN: 1バイト
- WRAPPED_KEY_LEN: 2バイト (ビッグエンディアン)
- NONCE: NONCE_LEN バイト
- WRAPPED_KEY: WRAPPED_KEY_LEN バイト
- CIPHERTEXT: 残り全バイト

#### Scenario: 暗号化ファイルの書き込み
- **WHEN** nonce, wrapped key, ciphertext が与えられる
- **THEN** ENVC ヘッダー付きのバイナリファイルを生成する

#### Scenario: ヘッダーのバージョン
- **WHEN** 暗号化ファイルを書き込む
- **THEN** VERSION フィールドは `0x01` が設定される

### Requirement: ENVC バイナリフォーマットの読み込み
The file reader SHALL parse the encrypted file header and extract each field.

#### Scenario: 正常なファイルの読み込み
- **WHEN** 有効な ENVC フォーマットのファイルを読み込む
- **THEN** nonce, wrapped key, ciphertext を正しく抽出する

#### Scenario: マジックバイトの検証
- **WHEN** 先頭4バイトが `ENVC` でないファイルを読み込む
- **THEN** エラー「invalid file format: missing ENVC header」を返す

#### Scenario: サポート外バージョン
- **WHEN** VERSION が `0x01` 以外のファイルを読み込む
- **THEN** エラー「unsupported file format version」を返す

#### Scenario: ファイルサイズの検証
- **WHEN** ヘッダーで宣言された NONCE_LEN + WRAPPED_KEY_LEN よりファイルの残りデータが少ない
- **THEN** エラー「corrupted file: unexpected end of data」を返す

