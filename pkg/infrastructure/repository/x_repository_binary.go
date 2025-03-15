package repository

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

// parseCompressedBinaryXFile は、圧縮バイナリ形式の X ファイルを解凍し、
// 解凍後のデータをバイナリパーサーへ渡します。
func (rep *XRepository) parseCompressedBinaryXFile(model *pmx.PmxModel) error {
	var decompressedBuffer []byte
	var err error

	// 圧縮バイナリ形式の X ファイルを解凍します。
	if decompressedBuffer, err = rep.decompressedBinaryXFile(); err != nil {
		return fmt.Errorf("failed to decompress binary X file: %w", err)
	}

	// 解凍後のデータをバイナリパーサーへ渡します。
	if err := rep.parseBinaryXFile(model, decompressedBuffer); err != nil {
		return fmt.Errorf("failed to parse binary X file: %w", err)
	}

	return nil
}

// parseBinaryXFile は、解凍済みのバイナリデータ buffer から各セクション（頂点、面、材質、テクスチャ）をパースします。
func (rep *XRepository) parseBinaryXFile(model *pmx.PmxModel, buffer []byte) error {
	r := bytes.NewReader(buffer)

	// ヘッダー読み出し（例："xof 0303bin 0032" を想定）
	header := make([]byte, 16)
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}
	headerStr := string(header)
	if len(headerStr) < 12 || headerStr[8:12] != "bin " {
		return fmt.Errorf("invalid binary X file header: %s", headerStr)
	}

	// 各セクションを別関数でパース
	if err := parseVertices(r, model); err != nil {
		return fmt.Errorf("failed to parse vertices: %w", err)
	}
	if err := parseFaces(r, model); err != nil {
		return fmt.Errorf("failed to parse faces: %w", err)
	}
	if err := parseMaterials(r, model); err != nil {
		return fmt.Errorf("failed to parse materials: %w", err)
	}
	if err := parseTextures(r, model); err != nil {
		return fmt.Errorf("failed to parse textures: %w", err)
	}

	model.Comment = fmt.Sprintf("Binary X File parsed: %d vertices, %d faces, %d materials, %d textures",
		model.Vertices.Length(), model.Faces.Length(), model.Materials.Length(), model.Textures.Length())

	return nil
}

// parseVertices は、バイナリデータから頂点数と各頂点の座標を読み取り、model.Vertices に登録します。
// フォーマット例:
//
//	int32 vertexCount
//	各頂点: float32 x, float32 y, float32 z
func parseVertices(r *bytes.Reader, model *pmx.PmxModel) error {
	var vertexCount int32
	if err := binary.Read(r, binary.LittleEndian, &vertexCount); err != nil {
		return fmt.Errorf("failed to read vertex count: %w", err)
	}

	for i := 0; i < int(vertexCount); i++ {
		var x, y, z float32
		if err := binary.Read(r, binary.LittleEndian, &x); err != nil {
			return fmt.Errorf("failed to read vertex x: %w", err)
		}
		if err := binary.Read(r, binary.LittleEndian, &y); err != nil {
			return fmt.Errorf("failed to read vertex y: %w", err)
		}
		if err := binary.Read(r, binary.LittleEndian, &z); err != nil {
			return fmt.Errorf("failed to read vertex z: %w", err)
		}

		// pmx の頂点作成（例：座標を10倍して登録）
		v := pmx.NewVertex()
		v.Position = &mmath.MVec3{
			X: float64(x) * 10,
			Y: float64(y) * 10,
			Z: float64(z) * 10,
		}
		// BDEF1、エッジ倍率、法線はデフォルト値を設定
		v.Deform = pmx.NewBdef1(0)
		v.EdgeFactor = 1.0
		v.Normal = &mmath.MVec3{X: 0, Y: 1, Z: 0}
		model.Vertices.Append(v)
	}
	return nil
}

// parseFaces は、バイナリデータから面数と各面の頂点インデックスを読み取り、model.Faces に登録します。
// フォーマット例:
//
//	int32 faceCount
//	各面:
//	  int32 vertexCount (通常3または4)
//	  その数の int32 インデックス
//
// なお、4頂点の場合はファン状に分割して2つの三角形とします。
func parseFaces(r *bytes.Reader, model *pmx.PmxModel) error {
	var faceCount int32
	if err := binary.Read(r, binary.LittleEndian, &faceCount); err != nil {
		return fmt.Errorf("failed to read face count: %w", err)
	}

	for i := 0; i < int(faceCount); i++ {
		var count int32
		if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
			return fmt.Errorf("failed to read vertex count for face %d: %w", i, err)
		}

		indices := make([]int, count)
		for j := 0; j < int(count); j++ {
			var idx int32
			if err := binary.Read(r, binary.LittleEndian, &idx); err != nil {
				return fmt.Errorf("failed to read index for face %d: %w", i, err)
			}
			indices[j] = int(idx)
		}

		// 三角形または四角形の場合の処理
		switch count {
		case 3:
			face := pmx.NewFace()
			face.VertexIndexes[0] = indices[0]
			face.VertexIndexes[1] = indices[1]
			face.VertexIndexes[2] = indices[2]
			model.Faces.Append(face)
		case 4:
			// 4頂点の場合、2つの三角形に分割
			face1 := pmx.NewFace()
			face1.VertexIndexes[0] = indices[0]
			face1.VertexIndexes[1] = indices[1]
			face1.VertexIndexes[2] = indices[2]
			model.Faces.Append(face1)

			face2 := pmx.NewFace()
			face2.VertexIndexes[0] = indices[0]
			face2.VertexIndexes[1] = indices[2]
			face2.VertexIndexes[2] = indices[3]
			model.Faces.Append(face2)
		default:
			// 4頂点以上の場合は、先頭頂点を基準としたファン分割
			for j := 1; j < int(count)-1; j++ {
				face := pmx.NewFace()
				face.VertexIndexes[0] = indices[0]
				face.VertexIndexes[1] = indices[j]
				face.VertexIndexes[2] = indices[j+1]
				model.Faces.Append(face)
			}
		}
	}
	return nil
}

// parseMaterials は、バイナリデータから材質情報を読み取り、model.Materials に登録します。
// フォーマット例:
//
//	int32 materialCount
//	各材質について、以下の順で読み出し:
//	  Diffuse: float32 r, g, b, a
//	  Power: float32
//	  Specular: float32 r, g, b
//	  Ambient: float32 r, g, b
func parseMaterials(r *bytes.Reader, model *pmx.PmxModel) error {
	var materialCount int32
	if err := binary.Read(r, binary.LittleEndian, &materialCount); err != nil {
		return fmt.Errorf("failed to read material count: %w", err)
	}

	for i := 0; i < int(materialCount); i++ {
		var dr, dg, db, da float32
		if err := binary.Read(r, binary.LittleEndian, &dr); err != nil {
			return fmt.Errorf("failed to read diffuse R for material %d: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &dg); err != nil {
			return fmt.Errorf("failed to read diffuse G for material %d: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &db); err != nil {
			return fmt.Errorf("failed to read diffuse B for material %d: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &da); err != nil {
			return fmt.Errorf("failed to read diffuse A for material %d: %w", i, err)
		}

		var power float32
		if err := binary.Read(r, binary.LittleEndian, &power); err != nil {
			return fmt.Errorf("failed to read power for material %d: %w", i, err)
		}

		var sr, sg, sb float32
		if err := binary.Read(r, binary.LittleEndian, &sr); err != nil {
			return fmt.Errorf("failed to read specular R for material %d: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &sg); err != nil {
			return fmt.Errorf("failed to read specular G for material %d: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &sb); err != nil {
			return fmt.Errorf("failed to read specular B for material %d: %w", i, err)
		}

		var ar, ag, ab float32
		if err := binary.Read(r, binary.LittleEndian, &ar); err != nil {
			return fmt.Errorf("failed to read ambient R for material %d: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &ag); err != nil {
			return fmt.Errorf("failed to read ambient G for material %d: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &ab); err != nil {
			return fmt.Errorf("failed to read ambient B for material %d: %w", i, err)
		}

		m := pmx.NewMaterial()
		m.Diffuse = &mmath.MVec4{X: float64(dr), Y: float64(dg), Z: float64(db), W: float64(da)}
		m.Specular = &mmath.MVec4{X: float64(sr), Y: float64(sg), Z: float64(sb), W: float64(power)}
		m.Ambient = &mmath.MVec3{X: float64(ar), Y: float64(ag), Z: float64(ab)}
		m.SetName(fmt.Sprintf("材質%02d", i+1))
		model.Materials.Append(m)
	}
	return nil
}

// parseTextures は、バイナリデータからテクスチャ情報（テクスチャ名）を読み取り、model.Textures に登録します。
// フォーマット例:
//
//	int32 textureCount
//	各テクスチャについて:
//	  int32 stringLength
//	  [stringLength] バイトのテクスチャ名（Shift-JIS等の場合は必要に応じて変換してください）
func parseTextures(r *bytes.Reader, model *pmx.PmxModel) error {
	var textureCount int32
	if err := binary.Read(r, binary.LittleEndian, &textureCount); err != nil {
		return fmt.Errorf("failed to read texture count: %w", err)
	}

	for i := 0; i < int(textureCount); i++ {
		var strLen int32
		if err := binary.Read(r, binary.LittleEndian, &strLen); err != nil {
			return fmt.Errorf("failed to read texture string length for texture %d: %w", i, err)
		}
		strBytes := make([]byte, strLen)
		if _, err := r.Read(strBytes); err != nil {
			return fmt.Errorf("failed to read texture name for texture %d: %w", i, err)
		}
		textureName := string(strBytes)
		tex := pmx.NewTexture()
		tex.SetName(textureName)
		model.Textures.Append(tex)
	}
	return nil
}
