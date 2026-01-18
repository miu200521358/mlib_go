// 指示: miu200521358
package x

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
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
	vertexOffset int
	faceGroups   [][]*model.Face
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
	if _, err := modelData.Bones.GetByName("センター"); err == nil {
		return
	}
	bone := &model.Bone{}
	bone.SetName("センター")
	bone.ParentIndex = -1
	bone.TailIndex = -1
	bone.EffectIndex = -1
	bone.BoneFlag = model.BONE_FLAG_CAN_MANIPULATE | model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_CAN_TRANSLATE | model.BONE_FLAG_IS_VISIBLE
	bone.IsSystem = false
	bone.Position = mmath.Vec3{}
	modelData.Bones.AppendRaw(bone)
}
