# AGENTS ガイド - リファクタリング版

## このドキュメントについて

- **目的**: mlib_go を クリーンアーキテクチャ に基づいてリファクタリングするためのガイド
- **対象**: このリポジトリで作業するエージェントおよび開発者
- **ステータス**: リファクタリング進行中（完了後は通常の開発ガイドに更新予定）

---

## リファクタリングの目標

### 1. クリーンアーキテクチャの採用

依存関係は内側から外側への一方向のみ許可:

```
┌─────────────────────────────────────────────────────────────────┐
│                    INFRASTRUCTURE (infra/)                       │
│   mbullet, mgl, mfile, mconfig - 外部ライブラリ依存              │
└─────────────────────────────────┬───────────────────────────────┘
                                  │ implements
                                  ↓
┌─────────────────────────────────────────────────────────────────┐
│                  INTERFACE ADAPTER (adapter/)                    │
│   mgateway, mpresenter, mcontroller - データ変換                │
└─────────────────────────────────┬───────────────────────────────┘
                                  │ implements
                                  ↓
┌─────────────────────────────────────────────────────────────────┐
│                      USE CASE (usecase/)                         │
│   minput (IF), moutput (IF), minteractor - ビジネスロジック      │
└─────────────────────────────────┬───────────────────────────────┘
                                  │ depends on
                                  ↓
┌─────────────────────────────────────────────────────────────────┐
│                       DOMAIN (domain/)                           │
│   mmath, mmodel, mmotion, mdelta - 純粋なエンティティ            │
│   ※ 外部依存なし                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 2. 命名規則

| ルール | 説明 | 例 |
|--------|------|-----|
| フォルダ（中間） | プレフィックスなし | `domain/`, `usecase/port/` |
| フォルダ（最終） | `m` プレフィックス | `mmath/`, `mmodel/`, `mgateway/` |
| ファイル | プレフィックスなし | `vector.go`, `bone.go` |
| 型名 | パッケージ名で区別 | `mmath.Vec3`, `mmodel.Bone` |

### 3. 空パッケージの削除

用途のない空のパッケージは作成しない。将来必要になった時点で追加する。

---

## 新しいパッケージ構成

```
pkg/
│
├── domain/                      # ドメイン層（外部依存なし）
│   ├── mmath/                   # 数学ライブラリ
│   ├── mmodel/                  # PMXモデルエンティティ
│   ├── mmotion/                 # VMDモーションエンティティ
│   └── mdelta/                  # 変形差分
│
├── usecase/                     # ユースケース層
│   ├── port/
│   │   ├── minput/              # 入力ポートインターフェース
│   │   └── moutput/             # 出力ポートインターフェース（リポジトリIF等）
│   └── minteractor/             # ユースケース実装
│
├── adapter/                     # インターフェースアダプター層
│   ├── mgateway/                # リポジトリ実装（ファイルI/O）
│   ├── mpresenter/              # 出力変換
│   └── mcontroller/             # 入力処理
│
└── infra/                       # インフラストラクチャ層
    ├── physics/
    │   ├── mbt/                 # Bullet SWIG
    │   └── mbullet/             # 物理エンジン実装
    ├── render/
    │   └── mgl/                 # OpenGLレンダリング
    ├── file/
    │   └── mfile/               # ファイルユーティリティ
    └── config/
        └── mconfig/             # 設定・ログ・i18n
```

---

## パッケージ詳細

### domain/mmath（数学ライブラリ）✅ 実装完了

**ファイル構成:**
- `vector2.go` - Vec2（2次元ベクトル）
- `vector3.go` - Vec3（3次元ベクトル）
- `vector4.go` - Vec4（4次元ベクトル）
- `matrix.go` - Mat4（4x4行列）
- `quaternion.go` - Quaternion（クォータニオン）
- `curve.go` - Curve（ベジェ補間曲線）
- `scalar.go` - スカラー演算ユーティリティ、汎用関数
- `number.go` - Number, SignedNumber, Float（ジェネリクス型制約）
- `bounding.go` - BoundingBox, BoundingSphere, BoundingCapsule（※未実装）
- `llrb.go` - LLRB木（インデックス管理）（※未実装）

**依存**: なし（標準ライブラリのみ）

---

### domain/mmodel（PMXモデル）

**ファイル構成:**
- `pmx_model.go` - PmxModel（モデル全体）
- `bone.go` - Bone, Bones, IkLink, Ik
- `vertex.go` - Vertex, Vertices, Deform (Bdef1, Bdef2, Bdef4, Sdef)
- `face.go` - Face, Faces
- `material.go` - Material, Materials
- `morph.go` - Morph, Morphs, VertexMorphOffset, BoneMorphOffset, etc.
- `rigid_body.go` - RigidBody, RigidBodies
- `joint.go` - Joint, Joints
- `texture.go` - Texture, Textures
- `display_slot.go` - DisplaySlot, DisplaySlots
- `bone_config.go` - BoneConfig, StandardBoneName
- `collection.go` - IndexModels[T], IndexNameModels[T]

**依存**: `mmath` のみ

---

### domain/mmotion（VMDモーション）

**ファイル構成:**
- `vmd_motion.go` - VmdMotion（モーション全体）
- `base_frame.go` - IBaseFrame, BaseFrame, BaseFrames[T]
- `bone_frame.go` - BoneFrame, BoneNameFrames, BoneFrames
- `morph_frame.go` - MorphFrame, MorphNameFrames, MorphFrames
- `camera_frame.go` - CameraFrame, CameraFrames
- `light_frame.go` - LightFrame, LightFrames
- `shadow_frame.go` - ShadowFrame, ShadowFrames
- `ik_frame.go` - IkFrame, IkFrames
- `physics_frame.go` - GravityFrame, PhysicsResetFrame, Wind関連Frame
- `rigid_body_frame.go` - RigidBodyFrame, RigidBodyNameFrames, RigidBodyFrames
- `joint_frame.go` - JointFrame, JointNameFrames, JointFrames
- `curve.go` - BoneCurves, CameraCurves

**依存**: `mmath` のみ

---

### domain/mdelta（変形差分）

**ファイル構成:**
- `vmd_deltas.go` - VmdDeltas（変形結果全体）
- `bone_delta.go` - BoneDelta, BoneDeltas
- `morph_delta.go` - VertexMorphDelta, BoneMorphDelta, MaterialMorphDelta, MorphDeltas
- `physics_delta.go` - PhysicsDeltas, RigidBodyDelta, JointDelta
- `mesh_delta.go` - MeshDelta

**依存**: `mmath`, `mmodel`

---

### usecase/port/minput（入力ポート）

**ファイル構成:**
- `model_usecase.go` - IModelUsecase interface
- `motion_usecase.go` - IMotionUsecase interface
- `deform_usecase.go` - IDeformUsecase interface

**例:**
```go
type IDeformUsecase interface {
    DeformModel(model *mmodel.PmxModel, motion *mmotion.VmdMotion, frame float32) *mdelta.VmdDeltas
    DeformBone(model *mmodel.PmxModel, motion *mmotion.VmdMotion, boneIndex int, frame float32) *mdelta.BoneDelta
}
```

---

### usecase/port/moutput（出力ポート）

**ファイル構成:**
- `model_repository.go` - IModelRepository interface
- `motion_repository.go` - IMotionRepository interface
- `image_repository.go` - IImageRepository interface
- `physics_engine.go` - IPhysicsEngine interface
- `renderer.go` - IRenderer interface

**例:**
```go
type IModelRepository interface {
    Load(path string) (*mmodel.PmxModel, error)
    Save(path string, model *mmodel.PmxModel) error
}

type IPhysicsEngine interface {
    AddModel(modelIndex int, model *mmodel.PmxModel) error
    DeleteModel(modelIndex int) error
    StepSimulation(timeStep float32, maxSubSteps int, fixedTimeStep float32) error
    GetBoneMatrix(modelIndex int, boneIndex int) *mmath.Mat4
}
```

---

### usecase/minteractor（ユースケース実装）

**ファイル構成:**
- `model_interactor.go` - ModelInteractor（IModelUsecase実装）
- `motion_interactor.go` - MotionInteractor（IMotionUsecase実装）
- `deform_interactor.go` - DeformInteractor（IDeformUsecase実装）
- `bone_deform.go` - ボーン変形ロジック
- `ik_deform.go` - IK計算ロジック
- `morph_deform.go` - モーフ変形ロジック
- `physics_deform.go` - 物理変形ロジック

**依存**: `domain/*`, `usecase/port/*`（インターフェースのみ）

---

### adapter/mgateway（リポジトリ実装）

**ファイル構成:**
- `pmx_reader.go` - PmxReader（IModelRepository実装）
- `pmx_writer.go` - PmxWriter
- `vmd_reader.go` - VmdReader（IMotionRepository実装）
- `vmd_writer.go` - VmdWriter
- `vpd_reader.go` - VpdReader
- `image_loader.go` - ImageLoader（IImageRepository実装）

**依存**: `domain/*`, `usecase/port/moutput`, `infra/file/mfile`

---

### infra/physics/mbullet（物理エンジン実装）

**ファイル構成:**
- `physics_engine.go` - BulletEngine（IPhysicsEngine実装）
- `rigid_body.go` - 剛体管理
- `joint.go` - ジョイント管理
- `wind.go` - 風シミュレーション
- `debug_view.go` - デバッグ描画

**依存**: `domain/*`, `usecase/port/moutput`, `infra/physics/mbt`

---

## 移行チェックリスト

### Phase 1: Domain層の移行
- [x] `mmath` - 数学ライブラリ（依存なし、最初に移行）
- [ ] `mmodel` - PMXモデルエンティティ
- [ ] `mmotion` - VMDモーションエンティティ  
- [ ] `mdelta` - 変形差分

### Phase 2: UseCase層の構築
- [ ] `minput` - 入力ポートインターフェース定義
- [ ] `moutput` - 出力ポートインターフェース定義
- [ ] `minteractor` - ユースケース実装

### Phase 3: Adapter層の実装
- [ ] `mgateway` - ファイルI/O（PMX/VMDリーダー・ライター）
- [ ] `mpresenter` - 描画用データ変換
- [ ] `mcontroller` - UI入力処理

### Phase 4: Infrastructure層の移行
- [ ] `mbt` - Bullet SWIG（既存をコピー）
- [ ] `mbullet` - 物理エンジン実装
- [ ] `mgl` - OpenGLレンダリング
- [ ] `mfile` - ファイルユーティリティ
- [ ] `mconfig` - 設定・ログ・i18n

### Phase 5: 統合とテスト
- [ ] 全体の結合テスト
- [ ] 既存テストの移行
- [ ] パフォーマンス検証

---

## 移行時の注意事項

### 1. 依存方向の厳守

```go
// ❌ NG: domain から usecase への依存
package mmodel
import "github.com/miu200521358/mlib_go/pkg/usecase/minteractor"

// ✅ OK: usecase から domain への依存
package minteractor
import "github.com/miu200521358/mlib_go/pkg/domain/mmodel"
```

### 2. インターフェースによる依存性逆転

```go
// usecase/port/moutput/physics_engine.go
type IPhysicsEngine interface {
    StepSimulation(timeStep float32) error
}

// usecase/minteractor/physics_deform.go
type PhysicsDeformer struct {
    engine moutput.IPhysicsEngine  // インターフェースに依存
}

// infra/physics/mbullet/physics_engine.go
type BulletEngine struct { ... }
func (e *BulletEngine) StepSimulation(timeStep float32) error { ... }
```

### 3. 既存機能の維持

- 既存の public API は可能な限り維持
- 破壊的変更が必要な場合は移行ガイドを追記

### 4. テストの同時移行

- 各パッケージ移行時にテストも一緒に移動
- テストが通ることを確認してから次へ進む

---

## ビルド手順

```bash
# 環境設定
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

# 依存関係の取得
go mod tidy

# フォーマット
go fmt ./...

# テスト実行
go test ./...

# ビルド
go build -o build/mlib.exe cmd/main.go
```

---

## コーディングルール

### Goコーディングスタイル

#### 命名規則
| 対象 | 規則 | 例 |
|------|------|-----|
| 変数 | キャメルケース | `boneIndex`, `modelPath` |
| エクスポート関数/型 | パスカルケース | `NewBone()`, `PmxModel` |
| インターフェース | `I` プレフィックス | `IBaseFrame`, `IModelRepository` |
| メソッドレシーバー | 1〜2文字 | `b *Bone`, `v *Vec3` |
| 定数（const） | パスカルケース | `MaxBones`, `DefaultFPS` |
| 擬似定数（var） | 全大文字スネークケース | `IDENTITY_MATRIX`, `ZERO_VECTOR` |

#### 定数の使い分け
```go
// const で宣言できるもの（プリミティブ型）: パスカルケース
const MaxBones = 512
const DefaultFPS = 30

// var だが変更しない値（構造体等）: 全大文字スネークケースで「不変」を明示
var IDENTITY_MATRIX = NewMat4()
var ZERO_VECTOR = NewVec3(0, 0, 0)
```

#### メソッド命名規則（破壊的/非破壊）

| タイプ | 命名 | 例 | 説明 |
|--------|------|-----|------|
| 破壊的 | 動詞（現在形） | `Add`, `Sub`, `Mul`, `Normalize` | レシーバを変更し、自身を返す |
| 非破壊 | 動詞＋ed | `Added`, `Subed`, `Muled`, `Normalized` | 新しいオブジェクトを返す |
| 非破壊 | 形容詞 | `Copy`, `Clamped`, `Inverted` | 新しいオブジェクトを返す |

```go
// ✅ 破壊的メソッド: レシーバを変更し、チェーン呼び出し可能
func (v *Vec3) Add(other *Vec3) *Vec3 {
    v.X += other.X
    v.Y += other.Y
    v.Z += other.Z
    return v  // 自身を返す
}

// ✅ 非破壊メソッド: 新しいオブジェクトを返す
func (v *Vec3) Added(other *Vec3) *Vec3 {
    return NewVec3ByValues(v.X+other.X, v.Y+other.Y, v.Z+other.Z)
}
```

#### コメント
- エクスポートされる関数・型には必ずドキュメントコメントを記述
- 複雑なロジックには説明コメントを追加
- TODOコメントには担当者と日付を記載: `// TODO(name): 説明 (2024-01-01)`

---

## エラー処理

### 基本方針
- **panic は絶対に使用しない**
- エラーは必ず `error` として返り値で上位に伝播させる
- 最終的にアプリケーションの `main.go` でキャッチし、UIダイアログで表示

### エラーの流れ

```
[domain/usecase層]                    [adapter層]                    [main.go]
     │                                     │                             │
     │  return err                         │  return err                 │
     ├──────────────────────────────────────┼─────────────────────────────┤
     │                                     │                             │
     │  NewReadError()                     │  errors.As() で判定         │
     │  NewTerminateError()                │                             │
     └─────────────────────────────────────────────────────────────────────┘
                                                                          │
                                                         merr.ShowErrorDialog()
                                                         merr.ShowFatalErrorDialog()
```

### カスタムエラー型

適度な粒度でカスタムエラー型を定義（細かすぎない）:

| エラー型 | 用途 | ダイアログ |
|----------|------|-----------|
| `ReadError` | PMX/VMD等のファイル読み取りエラー | ShowErrorDialog |
| `WriteError` | ファイル書き込みエラー | ShowErrorDialog |
| `ValidationError` | データ検証エラー | ShowErrorDialog |
| `NotFoundError` | 要素が見つからない | ShowErrorDialog |
| `TerminateError` | 致命的エラー（アプリ終了） | ShowFatalErrorDialog |

### 実装ルール

**エラーを生成する側:**
```go
// ✅ OK: カスタムエラーを返す
func (bones *Bones) GetByName(name string) (*Bone, error) {
    if bone, ok := bones.nameMap[name]; ok {
        return bone, nil
    }
    return nil, merr.NewNotFoundError("bone", name)
}
```

**エラーを伝播する側:**
```go
// ✅ OK: エラーをラップして返す
func LoadModel(path string) (*PmxModel, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read file %s: %w", path, err)
    }
    // ...
}
```

**コンストラクタ・初期化関数:**
```go
// ✅ OK: すべての初期化関数はerrorを返す
func NewControlWindow(config *Config) (*ControlWindow, error) {
    if err := validate(config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    // ...
}
```

**エラーを最終処理する側（main.go）:**
```go
func main() {
    if err := run(); err != nil {
        if merr.IsTerminateError(err) {
            merr.ShowFatalErrorDialog(appConfig, err)
        } else {
            merr.ShowErrorDialog(appConfig, err)
        }
    }
}
```

### 禁止事項
- ❌ `panic()` の使用
- ❌ エラーを握りつぶす（`_ = someFunc()` でエラーを無視）
- ❌ エラーメッセージに機密情報を含める

---

## テスト

### 基本ルール
- **各処理には必ずテストを書く**
- テストファイルは `*_test.go` として同一パッケージに配置
- テーブル駆動テスト（Table-Driven Tests）を使用

### テーブル駆動テストの形式

```go
func TestBone_GetByName(t *testing.T) {
    tests := []struct {
        name      string
        boneName  string
        wantErr   bool
        errType   error
    }{
        {
            name:     "存在するボーン",
            boneName: "センター",
            wantErr:  false,
        },
        {
            name:     "存在しないボーン",
            boneName: "unknown",
            wantErr:  true,
            errType:  &merr.NotFoundError{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実行
            result, err := bones.GetByName(tt.boneName)
            
            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error but got nil")
                }
                if tt.errType != nil && !errors.As(err, &tt.errType) {
                    t.Errorf("expected error type %T but got %T", tt.errType, err)
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
            }
        })
    }
}
```

### 描画系テストの指針

描画系（OpenGL等）は自動テストが困難なため、以下の戦略を採用:

| レイヤ | テスト方法 |
|--------|-----------|
| **ドメイン層** | ユニットテスト（完全自動化） |
| **ユースケース層** | モック注入でユニットテスト |
| **レンダラーロジック** | 計算部分のみユニットテスト |
| **実際の描画** | 目視確認 + スクリーンショット比較（手動/半自動） |

**テスト可能な部分:**
```go
// ✅ テスト可能: シェーダーのユニフォーム計算
func CalculateMVPMatrix(model, view, projection *mmath.Mat4) *mmath.Mat4

// ✅ テスト可能: バッファデータの生成
func BuildVertexBuffer(vertices *Vertices) []float32

// ⚠️ 目視確認: 実際のレンダリング結果
func Render(scene *Scene)
```

**描画テスト用ユーティリティ:**
- `test_resources/` に期待される出力画像を配置
- 必要に応じてスクリーンショット比較ツールを活用
- CI では描画テストをスキップ（`//go:build !ci` タグ）

### 破壊的/非破壊メソッドのテスト

**破壊的メソッド（元オブジェクトを変更）:**
- 元のオブジェクトが変更されていることを確認
- 戻り値がレシーバ自身であることを確認（チェーン呼び出し対応）

```go
func TestVec3_Add(t *testing.T) {
    tests := []struct {
        name     string
        v1X, v1Y, v1Z float64
        v2       *Vec3
        expected *Vec3
    }{ ... }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            v := NewVec3ByValues(tt.v1X, tt.v1Y, tt.v1Z)
            result := v.Add(tt.v2)
            // 元のオブジェクトが変更されていることを確認
            if v.X != tt.expected.X || v.Y != tt.expected.Y { ... }
            // 戻り値がレシーバ自身であることを確認
            if result != v { ... }
        })
    }
}
```

**非破壊メソッド（新しいオブジェクトを返す）:**
- 元のオブジェクトが変更されていないことを確認

```go
func TestVec3_Subed(t *testing.T) {
    tests := []struct {
        name     string
        v1X, v1Y, v1Z float64
        v2       *Vec3
        expected *Vec3
    }{ ... }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            v1 := NewVec3ByValues(tt.v1X, tt.v1Y, tt.v1Z)
            result := v1.Subed(tt.v2)
            // 結果の確認
            if !result.NearEquals(tt.expected, 1e-10) { ... }
            // 元のベクトルが変更されていないことを確認
            if v1.X != tt.v1X || v1.Y != tt.v1Y || v1.Z != tt.v1Z { ... }
        })
    }
}
```

---

## 設計原則

### SOLID原則
- **S (単一責任)**: 1つのパッケージ/型は1つの責任のみ
- **O (オープン・クローズド)**: 拡張に開き、修正に閉じる
- **L (リスコフの置換)**: インターフェースを満たす実装は置換可能に
- **I (インターフェース分離)**: 大きなインターフェースより小さな専用インターフェース
- **D (依存性逆転)**: 具象ではなく抽象（インターフェース）に依存

### クリーンアーキテクチャのルール
- 依存は内側から外側への一方向のみ
- 外側の層はインターフェースを実装
- DIコンテナは使用せず、main.goで依存性を組み立て

### 禁止事項
- ❌ domain層から外側の層への依存
- ❌ 循環参照（import cycle）
- ❌ グローバル変数の使用
- ❌ init()での副作用のある処理
- ❌ インターフェースを返す関数（具象型を返すべき）

---

## SWIG/Bullet

### 基本方針
- `infra/physics/mbt/` の SWIG 生成コードは既存をそのまま流用
- SWIG の再生成が必要な場合は手動で実行（リポジトリ管理者が実施）

### 変更時の注意
- `bullet.i` を変更した場合のみ SWIG 再実行が必要
- 再生成後は `bt.go`, `bt.cxx`, `bt.h` が更新される

---

## リファクタリング完了後

リファクタリングが完了したら、このAGENTS.mdを通常の開発ガイドに更新する:

1. 移行チェックリストを削除
2. 移行時の注意事項を通常の開発ルールに変更
3. 新機能追加のベストプラクティスを追記
4. トラブルシューティングセクションを追加
