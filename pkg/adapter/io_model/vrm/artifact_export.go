// 指示: miu200521358
package vrm

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	exportDirMode   = 0o755
	exportFileMode  = 0o644
	glbBINChunkType = 0x004E4942
)

// ArtifactExportResult はVRM由来の補助出力結果を表す。
type ArtifactExportResult struct {
	GltfPath     string
	BinPath      string
	TextureNames []string
}

// artifactExportDocument は補助出力に必要な glTF 要素を表す。
type artifactExportDocument struct {
	BufferViews []artifactBufferView `json:"bufferViews"`
	Images      []artifactImage      `json:"images"`
}

// artifactBufferView は glTF の bufferView 要素を表す。
type artifactBufferView struct {
	ByteOffset int `json:"byteOffset"`
	ByteLength int `json:"byteLength"`
}

// artifactImage は glTF の image 要素を表す。
type artifactImage struct {
	Name       string `json:"name"`
	URI        string `json:"uri"`
	BufferView *int   `json:"bufferView"`
	MimeType   string `json:"mimeType"`
}

// ExportArtifacts はVRMから glTF とテクスチャ補助出力を生成する。
func ExportArtifacts(vrmPath string, gltfDir string, textureDir string) (*ArtifactExportResult, error) {
	trimmedVrmPath := strings.TrimSpace(vrmPath)
	if trimmedVrmPath == "" {
		return nil, fmt.Errorf("VRMパスが未指定です")
	}
	if !strings.EqualFold(filepath.Ext(trimmedVrmPath), ".vrm") {
		return nil, fmt.Errorf("VRM拡張子ではありません: %s", trimmedVrmPath)
	}
	if strings.TrimSpace(gltfDir) == "" {
		return nil, fmt.Errorf("glTF出力先ディレクトリが未指定です")
	}
	if strings.TrimSpace(textureDir) == "" {
		return nil, fmt.Errorf("テクスチャ出力先ディレクトリが未指定です")
	}
	if err := os.MkdirAll(gltfDir, exportDirMode); err != nil {
		return nil, fmt.Errorf("glTF出力先ディレクトリの作成に失敗しました: %w", err)
	}
	if err := os.MkdirAll(textureDir, exportDirMode); err != nil {
		return nil, fmt.Errorf("テクスチャ出力先ディレクトリの作成に失敗しました: %w", err)
	}

	sourceBytes, err := os.ReadFile(trimmedVrmPath)
	if err != nil {
		return nil, fmt.Errorf("VRMファイルの読み取りに失敗しました: %w", err)
	}
	jsonChunk, binChunk, err := parseGLBChunks(sourceBytes)
	if err != nil {
		return nil, err
	}

	baseName := strings.TrimSpace(strings.TrimSuffix(filepath.Base(trimmedVrmPath), filepath.Ext(trimmedVrmPath)))
	if baseName == "" {
		baseName = "model"
	}

	result := &ArtifactExportResult{
		TextureNames: []string{},
	}
	result.GltfPath = filepath.Join(gltfDir, baseName+".gltf")
	if err := os.WriteFile(result.GltfPath, jsonChunk, exportFileMode); err != nil {
		return nil, fmt.Errorf("glTF JSON の保存に失敗しました: %w", err)
	}
	if len(binChunk) > 0 {
		result.BinPath = filepath.Join(gltfDir, baseName+".bin")
		if err := os.WriteFile(result.BinPath, binChunk, exportFileMode); err != nil {
			return nil, fmt.Errorf("glTF BIN の保存に失敗しました: %w", err)
		}
	}

	var doc artifactExportDocument
	if err := json.Unmarshal(jsonChunk, &doc); err != nil {
		return nil, fmt.Errorf("glTF JSON の解析に失敗しました: %w", err)
	}
	textureNames, err := exportTexturesFromDocument(&doc, binChunk, trimmedVrmPath, textureDir, baseName)
	if err != nil {
		return nil, err
	}
	result.TextureNames = textureNames
	return result, nil
}

// parseGLBChunks はGLBバイト列からJSON/BINチャンクを抽出する。
func parseGLBChunks(sourceBytes []byte) ([]byte, []byte, error) {
	if len(sourceBytes) < glbMinValidLength {
		return nil, nil, fmt.Errorf("VRMヘッダが不足しています")
	}
	if binary.LittleEndian.Uint32(sourceBytes[0:4]) != glbMagic {
		return nil, nil, fmt.Errorf("GLBマジックが不正です")
	}
	if binary.LittleEndian.Uint32(sourceBytes[4:8]) != 2 {
		return nil, nil, fmt.Errorf("GLBバージョンが未対応です")
	}

	totalLength := int(binary.LittleEndian.Uint32(sourceBytes[8:12]))
	if totalLength <= 0 || totalLength > len(sourceBytes) {
		return nil, nil, fmt.Errorf("GLB全体長が不正です")
	}

	var jsonChunk []byte
	var binChunk []byte
	offset := glbHeaderLength
	for offset+glbChunkHeadSize <= totalLength {
		chunkLength := int(binary.LittleEndian.Uint32(sourceBytes[offset : offset+4]))
		chunkType := binary.LittleEndian.Uint32(sourceBytes[offset+4 : offset+8])
		chunkStart := offset + glbChunkHeadSize
		chunkEnd := chunkStart + chunkLength
		if chunkLength < 0 || chunkEnd > totalLength {
			return nil, nil, fmt.Errorf("GLBチャンク長が不正です")
		}
		chunkBytes := sourceBytes[chunkStart:chunkEnd]
		switch chunkType {
		case glbJSONChunkType:
			jsonChunk = append([]byte(nil), chunkBytes...)
		case glbBINChunkType:
			if len(binChunk) == 0 {
				binChunk = append([]byte(nil), chunkBytes...)
			}
		}
		offset = chunkEnd
	}
	if len(jsonChunk) == 0 {
		return nil, nil, fmt.Errorf("GLB JSONチャンクが見つかりません")
	}
	return jsonChunk, binChunk, nil
}

// exportTexturesFromDocument は glTF document からテクスチャを抽出して保存する。
func exportTexturesFromDocument(
	doc *artifactExportDocument,
	binChunk []byte,
	vrmPath string,
	textureDir string,
	baseName string,
) ([]string, error) {
	if doc == nil || len(doc.Images) == 0 {
		return []string{}, nil
	}
	textureNames := make([]string, len(doc.Images))
	used := map[string]int{}
	for imageIndex, image := range doc.Images {
		imageBytes, ext, ok := resolveImageData(image, doc.BufferViews, binChunk, vrmPath)
		if !ok || len(imageBytes) == 0 {
			continue
		}
		if ext == "" {
			ext = detectImageExt(imageBytes)
		}
		if ext == "" {
			ext = ".bin"
		}
		nameBase := chooseTextureBaseName(image, imageIndex, baseName)
		fileName := buildUniqueTextureFileName(nameBase, ext, used)
		savePath := filepath.Join(textureDir, fileName)
		if err := os.WriteFile(savePath, imageBytes, exportFileMode); err != nil {
			return nil, fmt.Errorf("テクスチャ抽出ファイルの保存に失敗しました: %w", err)
		}
		textureNames[imageIndex] = fileName
	}
	return textureNames, nil
}

// resolveImageData は image 要素から画像バイト列を解決する。
func resolveImageData(
	image artifactImage,
	bufferViews []artifactBufferView,
	binChunk []byte,
	vrmPath string,
) ([]byte, string, bool) {
	if strings.TrimSpace(image.URI) != "" {
		uri := strings.TrimSpace(image.URI)
		if strings.HasPrefix(uri, "data:") {
			data, ext, err := decodeDataURI(uri)
			if err != nil {
				return nil, "", false
			}
			if ext == "" {
				ext = extByMimeType(image.MimeType)
			}
			return data, ext, true
		}

		sourcePath := uri
		if !filepath.IsAbs(sourcePath) {
			sourcePath = filepath.Join(filepath.Dir(vrmPath), filepath.FromSlash(sourcePath))
		}
		data, err := os.ReadFile(sourcePath)
		if err != nil {
			return nil, "", false
		}
		ext := strings.ToLower(filepath.Ext(uri))
		if ext == "" {
			ext = extByMimeType(image.MimeType)
		}
		return data, ext, true
	}

	if image.BufferView == nil {
		return nil, "", false
	}
	viewIndex := *image.BufferView
	if viewIndex < 0 || viewIndex >= len(bufferViews) {
		return nil, "", false
	}
	view := bufferViews[viewIndex]
	if view.ByteLength <= 0 || view.ByteOffset < 0 {
		return nil, "", false
	}
	end := view.ByteOffset + view.ByteLength
	if end > len(binChunk) {
		return nil, "", false
	}
	data := append([]byte(nil), binChunk[view.ByteOffset:end]...)
	return data, extByMimeType(image.MimeType), true
}

// decodeDataURI はdata URIをデコードする。
func decodeDataURI(uri string) ([]byte, string, error) {
	commaIndex := strings.Index(uri, ",")
	if commaIndex <= 0 {
		return nil, "", fmt.Errorf("data URI の形式が不正です")
	}
	meta := strings.TrimPrefix(uri[:commaIndex], "data:")
	payload := uri[commaIndex+1:]

	parts := strings.Split(meta, ";")
	mediaType := ""
	isBase64 := false
	for _, part := range parts {
		token := strings.TrimSpace(part)
		if token == "" {
			continue
		}
		if token == "base64" {
			isBase64 = true
			continue
		}
		if mediaType == "" {
			mediaType = token
		}
	}
	if !isBase64 {
		return []byte(payload), extByMimeType(mediaType), nil
	}
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, "", err
	}
	return data, extByMimeType(mediaType), nil
}

// extByMimeType はMIMEタイプから拡張子を返す。
func extByMimeType(mimeType string) string {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	case "image/bmp":
		return ".bmp"
	case "image/gif":
		return ".gif"
	case "image/tga":
		return ".tga"
	default:
		return ""
	}
}

// detectImageExt はシグネチャから画像拡張子を推定する。
func detectImageExt(data []byte) string {
	if len(data) >= 8 &&
		data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
		data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
		return ".png"
	}
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return ".jpg"
	}
	if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return ".webp"
	}
	if len(data) >= 2 && data[0] == 'B' && data[1] == 'M' {
		return ".bmp"
	}
	if len(data) >= 6 && (string(data[0:6]) == "GIF87a" || string(data[0:6]) == "GIF89a") {
		return ".gif"
	}
	return ""
}

// chooseTextureBaseName は画像出力ファイルのベース名を決定する。
func chooseTextureBaseName(image artifactImage, index int, fallbackBase string) string {
	if name := sanitizeFileName(image.Name); name != "" {
		return name
	}
	if uri := strings.TrimSpace(image.URI); uri != "" && !strings.HasPrefix(uri, "data:") {
		base := strings.TrimSuffix(filepath.Base(filepath.FromSlash(uri)), filepath.Ext(uri))
		if name := sanitizeFileName(base); name != "" {
			return name
		}
	}
	base := sanitizeFileName(fallbackBase)
	if base == "" {
		base = "texture"
	}
	return fmt.Sprintf("%s_tex_%03d", base, index+1)
}

// buildUniqueTextureFileName は重複しないテクスチャファイル名を生成する。
func buildUniqueTextureFileName(base string, ext string, used map[string]int) string {
	safeBase := sanitizeFileName(base)
	if safeBase == "" {
		safeBase = "texture"
	}
	if ext == "" {
		ext = ".bin"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ext = strings.ToLower(ext)

	key := strings.ToLower(safeBase + ext)
	if _, exists := used[key]; !exists {
		used[key] = 1
		return safeBase + ext
	}
	serial := used[key]
	for {
		candidate := fmt.Sprintf("%s_%d", safeBase, serial)
		candidateKey := strings.ToLower(candidate + ext)
		if _, exists := used[candidateKey]; !exists {
			used[candidateKey] = 1
			used[key] = serial + 1
			return candidate + ext
		}
		serial++
	}
}

// sanitizeFileName はファイル名に使えない文字を置換する。
func sanitizeFileName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ""
	}
	replacer := strings.NewReplacer(
		"\\", "_",
		"/", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	safe := strings.TrimSpace(replacer.Replace(trimmed))
	if safe == "" {
		return ""
	}
	return strings.Trim(safe, ".")
}
