## ADDED Requirements

### Requirement: AES-256-GCM によるデータ暗号化
The encryption engine SHALL use AES-256-GCM to encrypt plaintext data. The nonce SHALL be 12 bytes of cryptographically secure randomness.

#### Scenario: 正常な暗号化
- **WHEN** 平文データと32バイトのデータキーが与えられる
- **THEN** 12バイトのランダムノンスを生成し、AES-256-GCM で暗号化して nonce + ciphertext(GCMタグ付き）を返す

#### Scenario: 不正なキー長
- **WHEN** 32バイト以外のキーで暗号化を試みる
- **THEN** エラーを返す

### Requirement: AES-256-GCM によるデータ復号
The encryption engine SHALL use AES-256-GCM to decrypt ciphertext.

#### Scenario: 正常な復号
- **WHEN** 正しいデータキーと暗号文（nonce + ciphertext）が与えられる
- **THEN** 元の平文データを返す

#### Scenario: 不正なキーでの復号
- **WHEN** 異なるデータキーで暗号文の復号を試みる
- **THEN** GCM認証エラーを返す

#### Scenario: 改ざんされた暗号文の復号
- **WHEN** 暗号文の一部が改ざんされた状態で復号を試みる
- **THEN** GCM認証エラーを返す

### Requirement: データキーの生成
A random 32-byte data key SHALL be generated for each file during encryption.

#### Scenario: データキー生成
- **WHEN** 暗号化処理が開始される
- **THEN** crypto/rand を使用して32バイトのランダムデータキーを生成する

### Requirement: データキーのラップ/アンラップ
Data keys SHALL be wrapped with the master key using AES-256-GCM and unwrapped during decryption.

#### Scenario: データキーのラップ
- **WHEN** 32バイトのデータキーと32バイトのマスターキーが与えられる
- **THEN** マスターキーで AES-256-GCM を使ってデータキーを暗号化し、wrapped key を返す

#### Scenario: データキーのアンラップ
- **WHEN** wrapped key と正しいマスターキーが与えられる
- **THEN** 元のデータキーを復元する

#### Scenario: 不正なマスターキーでのアンラップ
- **WHEN** 異なるマスターキーで wrapped key のアンラップを試みる
- **THEN** GCM認証エラーを返す
