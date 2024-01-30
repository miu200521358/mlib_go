package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type BoneFrames struct {
	Data map[string]*BoneNameFrames
}

func NewBoneFrames() *BoneFrames {
	return &BoneFrames{
		Data: make(map[string]*BoneNameFrames, 0),
	}
}

func (bfs *BoneFrames) Contains(boneName string) bool {
	_, ok := bfs.Data[boneName]
	return ok
}

func (bfs *BoneFrames) Append(bnfs *BoneNameFrames) {
	bfs.Data[bnfs.Name] = bnfs
}

func (bfs *BoneFrames) GetItem(boneName string) *BoneNameFrames {
	if !bfs.Contains(boneName) {
		bfs.Append(NewBoneNameFrames(boneName))
	}
	return bfs.Data[boneName]
}

func (bfs *BoneFrames) Animate(
	fnos []int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	isOutLog bool,
	description string,
) BoneTrees {
	// 処理対象ボーン一覧取得
	targetBoneNames := bfs.getAnimatedBoneNames(model, boneNames)

	// ボーン行列作成
	boneNameIndexes, boneOffsetMatrixes, boneRevertOffsetMatrixes, bonePositionMatrixes :=
		bfs.createBoneMatrixes(model, targetBoneNames)

	// ボーン変形行列操作
	positions, rotations, scales := bfs.getBoneMatrixes(fnos, model, targetBoneNames, isOutLog, description)

	// ボーン行列計算
	return bfs.calc(
		fnos,
		model,
		targetBoneNames,
		boneNameIndexes,
		boneOffsetMatrixes,
		boneRevertOffsetMatrixes,
		bonePositionMatrixes,
		positions,
		rotations,
		scales,
		isOutLog,
		description,
	)
}

func (bfs *BoneFrames) calc(
	fnos []int,
	model *pmx.PmxModel,
	targetBoneNames []string,
	boneNameIndexes map[string]int,
	boneOffsetMatrixes, boneRevertOffsetMatrixes, bonePositionMatrixes []*mmath.MMat4,
	positions, rotations, scales [][]*mmath.MMat4,
	isOutLog bool,
	description string,
) BoneTrees {
	// 各ボーンの座標変換行列×逆BOf行列
	matrixes := make([][]*mmath.MMat4, 0, len(fnos))
	for i := range fnos {
		matrixes = append(matrixes, make([]*mmath.MMat4, 0, len(targetBoneNames)))
		for j := range model.Bones.GetSortedData() {
			matrixes[i] = append(matrixes[i], mmath.NewMMat4())
			// 逆BOf行列(初期姿勢行列)
			matrixes[i][j].Mul(boneRevertOffsetMatrixes[j])
			// 位置
			matrixes[i][j].Mul(positions[i][j])
			// 回転
			matrixes[i][j].Mul(rotations[i][j])
			// スケール
			matrixes[i][j].Mul(scales[i][j])
		}
	}

	boneTrees := NewBoneTrees()

	resultMatrixes := make([][]*mmath.MMat4, 0, len(fnos))
	for i, fno := range fnos {
		resultMatrixes = append(resultMatrixes, make([]*mmath.MMat4, 0, len(targetBoneNames)))
		for j, bone := range model.Bones.GetSortedData() {
			resultMatrixes[i] = append(resultMatrixes[i], matrixes[i][j])
			for k := range bone.ParentBoneIndexes {
				// 親ボーンの変形行列を掛ける
				v := matrixes[i][k].Muled(resultMatrixes[i][j])
				resultMatrixes[i][j] = &v
			}
			// BOf行列を掛けてローカル行列を作成
			localMatrix := resultMatrixes[i][j].Muled(boneOffsetMatrixes[j])
			// 初期位置行列を掛けてグローバル行列を作成
			p := positions[i][j].Translation()
			r := rotations[i][j].Quaternion()
			s := scales[i][j].Scaling()
			boneTrees.SetItem(bone.Name, fno, NewBoneTree(
				bone.Name,
				fno,
				resultMatrixes[i][j], // グローバル行列はそのまま
				&localMatrix,         // ローカル行列
				&p,                   // 移動
				&r,                   // 回転
				&mmath.MVec3{s.GetX(), s.GetY(), s.GetZ()}, // 拡大率
			))
		}
	}

	return *boneTrees
}

// ボーン行列を作成する
func (bfs *BoneFrames) createBoneMatrixes(
	model *pmx.PmxModel,
	targetBoneNames []string,
) (map[string]int, []*mmath.MMat4, []*mmath.MMat4, []*mmath.MMat4) {
	boneNameIndexes := make(map[string]int, 0)
	boneOffsetMatrixes := make([]*mmath.MMat4, 0)
	boneRevertOffsetMatrixes := make([]*mmath.MMat4, 0)
	bonePositionMatrixes := make([]*mmath.MMat4, 0)

	for _, bone := range model.Bones.GetSortedData() {
		// ボーン名:インデックス
		boneNameIndexes[bone.Name] = bone.GetIndex()
		// ボーンのBOf行列
		boneOffsetMatrixes = append(boneOffsetMatrixes, bone.OffsetMatrix.Copy())
		// ボーンの逆BOf行列
		boneRevertOffsetMatrixes = append(boneRevertOffsetMatrixes, bone.RevertOffsetMatrix.Copy())
		// ボーンの初期位置行列
		posMat := mmath.NewMMat4()
		posMat.Translate(bone.Position)
		bonePositionMatrixes = append(bonePositionMatrixes, posMat)
	}

	return boneNameIndexes, boneOffsetMatrixes, boneRevertOffsetMatrixes, bonePositionMatrixes
}

// アニメーション対象ボーン一覧取得
func (bfs *BoneFrames) getAnimatedBoneNames(
	model *pmx.PmxModel,
	boneNames []string,
) []string {
	if len(boneNames) == 0 {
		return model.Bones.GetNames()
	} else {
		targetBoneNames := make([]string, 0)
		for _, boneName := range boneNames {
			if !slices.Contains(targetBoneNames, boneName) {
				targetBoneNames = append(targetBoneNames, boneName)
			}
			relativeBoneIndexes := model.Bones.GetItemByName(boneName).RelativeBoneIndexes
			for _, index := range relativeBoneIndexes {
				boneName := model.Bones.GetItem(index).Name
				if !slices.Contains(targetBoneNames, boneName) {
					targetBoneNames = append(targetBoneNames, boneName)
				}
			}
		}

		resultBoneNames := make([]string, 0)

		// ボーンINDEXでソート
		for _, bone := range model.Bones.GetSortedData() {
			if slices.Contains(targetBoneNames, bone.Name) {
				resultBoneNames = append(resultBoneNames, bone.Name)
			}
		}

		return resultBoneNames
	}
}

// ボーン変形行列を求める
func (bfs *BoneFrames) getBoneMatrixes(
	fnos []int,
	model *pmx.PmxModel,
	targetBoneNames []string,
	isOutLog bool,
	description string,
) ([][]*mmath.MMat4, [][]*mmath.MMat4, [][]*mmath.MMat4) {
	positions := make([][]*mmath.MMat4, 0, len(fnos))
	rotations := make([][]*mmath.MMat4, 0, len(fnos))
	scales := make([][]*mmath.MMat4, 0, len(fnos))

	for i, fno := range fnos {
		positions = append(positions, make([]*mmath.MMat4, 0, len(targetBoneNames)))
		rotations = append(rotations, make([]*mmath.MMat4, 0, len(targetBoneNames)))
		scales = append(scales, make([]*mmath.MMat4, 0, len(targetBoneNames)))
		for _, bone := range model.Bones.GetSortedData() {
			if slices.Contains(targetBoneNames, bone.Name) {
				// ボーンが対象の場合、そのボーンの移動位置、回転角度、拡大率を取得
				positions[i] = append(positions[i], bfs.getPosition(fno, bone.Name))
				rotations[i] = append(rotations[i], bfs.getRotation(fno, bone.Name))
				scales[i] = append(scales[i], bfs.getScale(fno, bone.Name))
			} else {
				// ボーンが対象外の場合、空の行列を追加
				positions[i] = append(positions[i], mmath.NewMMat4())
				rotations[i] = append(rotations[i], mmath.NewMMat4())
				scales[i] = append(scales[i], mmath.NewMMat4())
			}
		}
	}

	return positions, rotations, scales
}

// 該当キーフレにおけるボーンの移動位置
func (bfs *BoneFrames) getPosition(fno int, boneName string) *mmath.MMat4 {
	bf := bfs.Data[boneName].GetItem(fno)
	mat := mmath.NewMMat4()
	mat[0][3] = bf.Position.GetX()
	mat[1][3] = bf.Position.GetY()
	mat[2][3] = bf.Position.GetZ()
	return mat
}

// 該当キーフレにおけるボーンの回転角度
func (bfs *BoneFrames) getRotation(fno int, boneName string) *mmath.MMat4 {
	bf := bfs.Data[boneName].GetItem(fno)
	rot := bf.Rotation.GetQuaternion().ToMat4()
	return rot
}

// 該当キーフレにおけるボーンの拡大率
func (bfs *BoneFrames) getScale(fno int, boneName string) *mmath.MMat4 {
	bf := bfs.Data[boneName].GetItem(fno)
	mat := mmath.NewMMat4()
	mat[0][0] = bf.Scale.GetX() + 1.0
	mat[1][1] = bf.Scale.GetY() + 1.0
	mat[2][2] = bf.Scale.GetZ() + 1.0
	return mat
}
