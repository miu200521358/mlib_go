// 指示: miu200521358
package x

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

type xFormat int

const (
	xFormatText xFormat = iota
	xFormatBinary
	xFormatCompressedBinary
	xFormatInvalid
)

type meshContext struct {
	vertexOffset      int
	vertexCount       int
	faceGroups        [][]*model.Face
	faceIndexGroups   [][]int
	normals           []mmath.Vec3
	normalFaceIndexes [][]int
}

// XRepository はX形式の読み取りを表す。
type XRepository struct{}

// NewXRepository はXRepositoryを生成する。
func NewXRepository() *XRepository {
	return &XRepository{}
}

// CanLoad は拡張子に応じて読み込み可否を判定する。
func (r *XRepository) CanLoad(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".x")
}

// InferName はパスから表示名を推定する。
func (r *XRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Load はX形式の読み込みを行う。
func (r *XRepository) Load(path string) (hashable.IHashable, error) {
	if !r.CanLoad(path) {
		return nil, io_common.NewIoExtInvalid(path, nil)
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, io_common.NewIoFileNotFound(path, err)
		}
		return nil, io_common.NewIoParseFailed("Xファイルのオープンに失敗しました", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, io_common.NewIoParseFailed("Xファイルの読み取りに失敗しました", err)
	}
	if len(data) < 16 {
		return nil, io_common.NewIoParseFailed("Xヘッダが不足しています", nil)
	}

	format, err := detectFormat(data[:16])
	if err != nil {
		return nil, err
	}

	modelData := model.NewPmxModel()
	modelData.SetPath(path)
	modelData.SetName(r.InferName(path))

	switch format {
	case xFormatText:
		parser, err := newTextParser(bytes.NewReader(data))
		if err != nil {
			return nil, io_common.NewIoParseFailed("Xテキストの解析準備に失敗しました", err)
		}
		if err := parser.Parse(modelData); err != nil {
			return nil, io_common.NewIoParseFailed("Xテキストの解析に失敗しました", err)
		}
	case xFormatBinary:
		parser := newBinaryParser(data, modelData)
		if err := parser.Parse(); err != nil {
			return nil, io_common.NewIoParseFailed("Xバイナリの解析に失敗しました", err)
		}
	case xFormatCompressedBinary:
		decompressed, err := decompressMSZip(data)
		if err != nil {
			return nil, io_common.NewIoParseFailed("X圧縮バイナリの解凍に失敗しました", err)
		}
		parser := newBinaryParser(decompressed, modelData)
		if err := parser.Parse(); err != nil {
			return nil, io_common.NewIoParseFailed("X圧縮バイナリの解析に失敗しました", err)
		}
	default:
		return nil, io_common.NewIoFormatNotSupported("X形式の判定に失敗しました", nil)
	}

	addCenterBone(modelData)
	addDisplaySlot(modelData)

	info, err := file.Stat()
	if err != nil {
		return nil, io_common.NewIoParseFailed("Xファイル情報の取得に失敗しました", err)
	}
	modelData.SetFileModTime(info.ModTime().UnixNano())
	modelData.UpdateHash()
	return modelData, nil
}

// detectFormat はヘッダからX形式を判定する。
func detectFormat(header []byte) (xFormat, error) {
	if len(header) < 16 {
		return xFormatInvalid, io_common.NewIoParseFailed("Xヘッダが不足しています", nil)
	}
	indicator := string(header[8:12])
	switch indicator {
	case "txt ":
		return xFormatText, nil
	case "tzip":
		return xFormatBinary, nil
	case "bzip":
		return xFormatCompressedBinary, nil
	default:
		return xFormatInvalid, io_common.NewIoFormatNotSupported("X形式が未対応です", nil)
	}
}

// addCenterBone はセンターボーンを追加する。
func addCenterBone(modelData *model.PmxModel) {
	if modelData == nil || modelData.Bones == nil {
		return
	}
	if _, err := modelData.Bones.GetByName(model.CENTER.String()); err == nil {
		return
	}
	bone := &model.Bone{}
	bone.SetName(model.CENTER.String())
	bone.ParentIndex = -1
	bone.TailIndex = -1
	bone.EffectIndex = -1
	bone.BoneFlag = model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_CAN_TRANSLATE | model.BONE_FLAG_IS_VISIBLE
	bone.IsSystem = false
	bone.Position = mmath.Vec3{}
	modelData.Bones.AppendRaw(bone)
}

// addDisplaySlot は既定の表示枠を追加する。
func addDisplaySlot(modelData *model.PmxModel) {
	if modelData == nil || modelData.Bones == nil {
		return
	}
	if _, err := modelData.Bones.GetByName(model.CENTER.String()); err != nil {
		return
	}
	modelData.CreateDefaultDisplaySlots()
}

// applyMeshNormals は読み込んだメッシュ情報から頂点法線を確定する。
func applyMeshNormals(modelData *model.PmxModel, meshCtx *meshContext) error {
	if modelData == nil || meshCtx == nil {
		return nil
	}
	vertexCount := meshCtx.vertexCount
	if vertexCount <= 0 {
		vertexCount = modelData.Vertices.Len() - meshCtx.vertexOffset
	}
	if vertexCount <= 0 {
		return nil
	}

	switch {
	case len(meshCtx.normals) > 0 && len(meshCtx.normalFaceIndexes) > 0:
		return applyMeshNormalsByMapping(modelData, meshCtx, vertexCount)
	case len(meshCtx.normals) > 0:
		return applyMeshNormalsByVertexOrder(modelData, meshCtx, vertexCount)
	default:
		return applyMeshNormalsFromFaces(modelData, meshCtx, vertexCount)
	}
}

// applyMeshNormalsByMapping はMeshNormalsの面インデックスを使って法線を算出する。
func applyMeshNormalsByMapping(modelData *model.PmxModel, meshCtx *meshContext, vertexCount int) error {
	if len(meshCtx.faceIndexGroups) == 0 {
		return newParseFailed("MeshNormalsの面情報が不足しています")
	}
	if len(meshCtx.normalFaceIndexes) != len(meshCtx.faceIndexGroups) {
		return newParseFailed("MeshNormalsの面数が不正です")
	}

	accumulated := make([]mmath.Vec3, vertexCount)
	for faceIdx, faceIndexes := range meshCtx.faceIndexGroups {
		normalIndexes := meshCtx.normalFaceIndexes[faceIdx]
		if len(normalIndexes) < len(faceIndexes) {
			return newParseFailed("MeshNormalsの面頂点数が不足しています")
		}
		// 取り込んだ面頂点より多い法線が来た場合は、幾何の簡略化に合わせて切り詰める。
		if len(normalIndexes) > len(faceIndexes) {
			normalIndexes = normalIndexes[:len(faceIndexes)]
		}
		for i, vertexIndex := range faceIndexes {
			local := vertexIndex - meshCtx.vertexOffset
			if local < 0 || local >= vertexCount {
				return newParseFailed("MeshNormalsの頂点番号が不正です")
			}
			normalIndex := normalIndexes[i]
			if normalIndex < 0 || normalIndex >= len(meshCtx.normals) {
				return newParseFailed("MeshNormalsの法線番号が不正です")
			}
			accumulated[local].Add(meshCtx.normals[normalIndex])
		}
	}

	return applyAccumulatedNormals(modelData, meshCtx.vertexOffset, accumulated)
}

// applyMeshNormalsByVertexOrder はMeshNormalsの順序を頂点順とみなして法線を設定する。
func applyMeshNormalsByVertexOrder(modelData *model.PmxModel, meshCtx *meshContext, vertexCount int) error {
	limit := vertexCount
	if len(meshCtx.normals) < limit {
		limit = len(meshCtx.normals)
	}
	accumulated := make([]mmath.Vec3, vertexCount)
	for i := 0; i < limit; i++ {
		accumulated[i] = meshCtx.normals[i]
	}
	return applyAccumulatedNormals(modelData, meshCtx.vertexOffset, accumulated)
}

// applyMeshNormalsFromFaces は面の幾何情報から法線を算出する。
func applyMeshNormalsFromFaces(modelData *model.PmxModel, meshCtx *meshContext, vertexCount int) error {
	if len(meshCtx.faceGroups) == 0 {
		return nil
	}

	accumulated := make([]mmath.Vec3, vertexCount)
	for _, faces := range meshCtx.faceGroups {
		for _, face := range faces {
			v0, err := modelData.Vertices.Get(face.VertexIndexes[0])
			if err != nil {
				return err
			}
			v1, err := modelData.Vertices.Get(face.VertexIndexes[1])
			if err != nil {
				return err
			}
			v2, err := modelData.Vertices.Get(face.VertexIndexes[2])
			if err != nil {
				return err
			}
			edge1 := v1.Position.Subed(v0.Position)
			edge2 := v2.Position.Subed(v0.Position)
			faceNormal := edge1.Cross(edge2)
			if faceNormal.IsZero() {
				continue
			}
			for _, vertexIndex := range face.VertexIndexes {
				local := vertexIndex - meshCtx.vertexOffset
				if local < 0 || local >= vertexCount {
					return newParseFailed("Meshの頂点番号が不正です")
				}
				accumulated[local].Add(faceNormal)
			}
		}
	}

	return applyAccumulatedNormals(modelData, meshCtx.vertexOffset, accumulated)
}

// applyAccumulatedNormals は積算済み法線を正規化して頂点へ設定する。
func applyAccumulatedNormals(modelData *model.PmxModel, vertexOffset int, accumulated []mmath.Vec3) error {
	for i, normal := range accumulated {
		if normal.IsZero() {
			continue
		}
		vertex, err := modelData.Vertices.Get(vertexOffset + i)
		if err != nil {
			return err
		}
		vertex.Normal = normal.Normalized()
	}
	return nil
}
