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

### domain/mmath（数学ライブラリ）

**ファイル構成:**
- `vector.go` - Vec2, Vec3, Vec4
- `matrix.go` - Mat4
- `quaternion.go` - Quaternion
- `curve.go` - Curve（ベジェ補間）
- `bounding.go` - BoundingBox, BoundingSphere, BoundingCapsule
- `scalar.go` - スカラー演算ユーティリティ
- `llrb.go` - LLRB木（インデックス管理）

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
- [ ] `mmath` - 数学ライブラリ（依存なし、最初に移行）
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

## リファクタリング完了後

リファクタリングが完了したら、このAGENTS.mdを通常の開発ガイドに更新する:

1. 移行チェックリストを削除
2. 移行時の注意事項を通常の開発ルールに変更
3. 新機能追加のベストプラクティスを追記
4. トラブルシューティングセクションを追加
