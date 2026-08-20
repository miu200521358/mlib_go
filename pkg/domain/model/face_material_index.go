// 指示: miu200521358
package model

import "fmt"

// FaceIndexRange は、材質に属する面の半開区間 [Start, End) を表す。
//
// PMX の面 index は材質の並び順に連続しているため、各材質は一つの範囲で
// 表現できる。End は範囲に含まれないので、面数は End - Start で求められる。
type FaceIndexRange struct {
	Start int
	End   int
}

// FaceMaterialIndex は面 index と材質 index の相互参照を保持する。
//
// 内部スライスは構築後に変更しない。参照メソッドが返す値は整数または
// 値型のため、呼び出し側が索引を壊すことはない。
type FaceMaterialIndex struct {
	faceToMaterial      []int
	materialToFaceRange []FaceIndexRange
}

// NewFaceMaterialIndex はモデルの材質面範囲から面・材質索引を構築する。
//
// 材質の VerticesCount は三角形の頂点数（面数の 3 倍）として解釈する。
// nil モデル、負数または 3 の倍数でない VerticesCount、面数との合計不一致は
// 不完全な索引を返さずエラーにする。
func NewFaceMaterialIndex(modelData *PmxModel) (*FaceMaterialIndex, error) {
	if modelData == nil {
		return nil, fmt.Errorf("モデルが未設定です")
	}
	if modelData.Materials == nil {
		return nil, fmt.Errorf("材質コレクションが未設定です")
	}
	if modelData.Faces == nil {
		return nil, fmt.Errorf("面コレクションが未設定です")
	}

	faceCount := modelData.Faces.Len()
	index := &FaceMaterialIndex{
		faceToMaterial:      make([]int, faceCount),
		materialToFaceRange: make([]FaceIndexRange, modelData.Materials.Len()),
	}
	faceOffset := 0
	for materialIndex, materialData := range modelData.Materials.Values() {
		if materialData == nil {
			return nil, fmt.Errorf("材質が未設定です: index=%d", materialIndex)
		}
		verticesCount := materialData.VerticesCount
		if verticesCount < 0 || verticesCount%3 != 0 {
			return nil, fmt.Errorf("材質頂点数が不正です: index=%d verticesCount=%d", materialIndex, verticesCount)
		}
		materialFaceCount := verticesCount / 3
		if faceOffset+materialFaceCount > faceCount {
			return nil, fmt.Errorf(
				"材質面範囲が面数を超えています: index=%d start=%d count=%d faces=%d",
				materialIndex,
				faceOffset,
				materialFaceCount,
				faceCount,
			)
		}
		faceRange := FaceIndexRange{Start: faceOffset, End: faceOffset + materialFaceCount}
		index.materialToFaceRange[materialIndex] = faceRange
		for faceIndex := faceRange.Start; faceIndex < faceRange.End; faceIndex++ {
			index.faceToMaterial[faceIndex] = materialIndex
		}
		faceOffset = faceRange.End
	}
	if faceOffset != faceCount {
		return nil, fmt.Errorf("材質頂点数と面数が一致しません: mappedFaces=%d totalFaces=%d", faceOffset, faceCount)
	}
	return index, nil
}

// BuildFaceMaterialIndex は NewFaceMaterialIndex の別名である。
//
// 構築処理を動詞始まりで呼びたい利用箇所のために提供する。
func BuildFaceMaterialIndex(modelData *PmxModel) (*FaceMaterialIndex, error) {
	return NewFaceMaterialIndex(modelData)
}

// MaterialIndex は面 index に対応する材質 index を返す。
// 範囲外の面 index は false を返す。
func (i *FaceMaterialIndex) MaterialIndex(faceIndex int) (int, bool) {
	if i == nil || faceIndex < 0 || faceIndex >= len(i.faceToMaterial) {
		return -1, false
	}
	return i.faceToMaterial[faceIndex], true
}

// FaceRange は材質 index に対応する面の半開区間を返す。
// 範囲外の材質 index は false を返す。
func (i *FaceMaterialIndex) FaceRange(materialIndex int) (FaceIndexRange, bool) {
	if i == nil || materialIndex < 0 || materialIndex >= len(i.materialToFaceRange) {
		return FaceIndexRange{}, false
	}
	return i.materialToFaceRange[materialIndex], true
}

// FaceCount は索引が保持する面数を返す。
func (i *FaceMaterialIndex) FaceCount() int {
	if i == nil {
		return 0
	}
	return len(i.faceToMaterial)
}

// MaterialCount は索引が保持する材質数を返す。
func (i *FaceMaterialIndex) MaterialCount() int {
	if i == nil {
		return 0
	}
	return len(i.materialToFaceRange)
}
