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
	model pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	isOutLog bool,
	description string,
) BoneTrees {
	// 処理対象ボーン一覧取得
	targetBoneNames := bfs.getAnimatedBoneNames(model, boneNames)

	// ボーン行列作成
	boneNameIndexes, boneOffsetMatrixes, bonePositionMatrixes := bfs.createBoneMatrixes(model, targetBoneNames)

	// ボーン変形行列操作
	positions, rotations, scales := bfs.getBoneMatrixes(fnos, model, targetBoneNames, isOutLog, description)

	// ボーン行列計算
	return bfs.calc(
		fnos,
		model,
		targetBoneNames,
		boneNameIndexes,
		boneOffsetMatrixes,
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
	model pmx.PmxModel,
	targetBoneNames []string,
	boneNameIndexes map[string]int,
	boneOffsetMatrixes, bonePositionMatrixes []*mmath.MMat4,
	positions, rotations, scales [][]*mmath.MMat4,
	isOutLog bool,
	description string,
) BoneTrees {
	// 各ボーンの座標変換行列×逆BOf行列
	matrixes := make([][]*mmath.MMat4, 0, len(fnos))
	for i := range fnos {
		matrixes = append(matrixes, make([]*mmath.MMat4, 0, len(targetBoneNames)))
		for j := range targetBoneNames {
			matrixes[j] = append(matrixes[j], mmath.NewMMat4())
			// 逆BOf行列(初期姿勢行列)
			matrixes[j][i].Mul(boneOffsetMatrixes[j])
			// 位置
			matrixes[j][i].Mul(positions[j][i])
			// 回転
			matrixes[j][i].Mul(rotations[j][i])
			// スケール
			matrixes[j][i].Mul(scales[j][i])
		}
	}

	boneTrees := NewBoneTrees()

	resultMatrixes := make([][]*mmath.MMat4, 0, len(fnos))
	for i, fno := range fnos {
		resultMatrixes = append(resultMatrixes, make([]*mmath.MMat4, 0, len(targetBoneNames)))
		for j, boneName := range targetBoneNames {
			resultMatrixes[j] = append(resultMatrixes[j], matrixes[j][i])
			for k := range model.Bones.GetItemByName(boneName).ParentBoneIndexes {
				// 親ボーンの変形行列を掛ける
				resultMatrixes[j][i].Mul(matrixes[k][i])
			}
			globalMatrix := resultMatrixes[j][i].Muled(bonePositionMatrixes[j])
			p := positions[j][i].Translation()
			r := rotations[j][i].Quaternion()
			s := scales[j][i].Scaling()
			boneTrees.SetItem(boneName, fno, NewBoneTree(
				boneName,
				fno,
				&globalMatrix,        // ボーンの初期位置行列を掛ける
				resultMatrixes[j][i], // ローカル行列はそのまま
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
	model pmx.PmxModel,
	targetBoneNames []string,
) (map[string]int, []*mmath.MMat4, []*mmath.MMat4) {
	boneNameIndexes := make(map[string]int, 0)
	boneOffsetMatrixes := make([]*mmath.MMat4, 0)
	bonePositionMatrixes := make([]*mmath.MMat4, 0)

	for _, boneName := range targetBoneNames {
		bone := model.Bones.GetItemByName(boneName)
		// ボーン名:インデックス
		boneNameIndexes[boneName] = bone.GetIndex()
		// ボーンの初期姿勢行列
		boneOffsetMatrixes = append(boneOffsetMatrixes, bone.OffsetMatrix.Copy())
		// ボーンの初期位置行列
		posMat := mmath.NewMMat4()
		posMat[0][3] = bone.Position.GetX()
		posMat[1][3] = bone.Position.GetY()
		posMat[2][3] = bone.Position.GetZ()
		bonePositionMatrixes = append(bonePositionMatrixes, posMat)
	}

	return boneNameIndexes, boneOffsetMatrixes, bonePositionMatrixes
}

// アニメーション対象ボーン一覧取得
func (bfs *BoneFrames) getAnimatedBoneNames(
	model pmx.PmxModel,
	boneNames []string,
) []string {
	if len(boneNames) == 0 {
		return model.Bones.GetNames()
	} else {
		boneNames := make([]string, 0)
		for _, boneName := range boneNames {
			relativeBoneIndexes := model.Bones.GetItemByName(boneName).RelativeBoneIndexes
			for _, index := range relativeBoneIndexes {
				boneName := model.Bones.GetItem(index).Name
				if !slices.Contains(boneNames, boneName) {
					boneNames = append(boneNames, boneName)
				}
			}
		}
		return boneNames
	}
}

// ボーン変形行列を求める
func (bfs *BoneFrames) getBoneMatrixes(
	fnos []int,
	model pmx.PmxModel,
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
		for _, boneName := range targetBoneNames {
			positions[i] = append(positions[i], bfs.getPosition(fno, boneName))
			rotations[i] = append(rotations[i], bfs.getRotation(fno, boneName))
			scales[i] = append(scales[i], bfs.getScale(fno, boneName))
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
	mat[0][0] = bf.Scale.GetX()
	mat[1][1] = bf.Scale.GetY()
	mat[2][2] = bf.Scale.GetZ()
	return mat
}
