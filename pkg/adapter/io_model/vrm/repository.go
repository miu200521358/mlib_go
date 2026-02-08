// 指示: miu200521358
package vrm

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/vrm"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
	"gonum.org/v1/gonum/spatial/r3"
)

const (
	glbHeaderLength   = 12
	glbChunkHeadSize  = 8
	glbMagic          = 0x46546C67
	glbJSONChunkType  = 0x4E4F534A
	glbMinValidLength = glbHeaderLength + glbChunkHeadSize
)

// VrmRepository はVRM入力の読み込み契約を表す。
type VrmRepository struct{}

// NewVrmRepository はVrmRepositoryを生成する。
func NewVrmRepository() *VrmRepository {
	return &VrmRepository{}
}

// CanLoad は拡張子に応じて読み込み可否を判定する。
func (r *VrmRepository) CanLoad(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".vrm")
}

// InferName はパスから表示名を推定する。
func (r *VrmRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Load はVRMを読み込む。
func (r *VrmRepository) Load(path string) (hashable.IHashable, error) {
	if !r.CanLoad(path) {
		return nil, io_common.NewIoExtInvalid(path, nil)
	}
	loadTargetName := filepath.Base(path)
	logVrmStep("VRM読込開始: file=%s", loadTargetName)

	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, io_common.NewIoFileNotFound(path, err)
		}
		return nil, io_common.NewIoParseFailed("VRMファイルの読み取りに失敗しました", err)
	}
	logVrmStep("VRM読込ステップ: ファイル読み取り完了 bytes=%d", len(b))

	jsonChunk, binChunk, err := parseGLBChunks(b)
	if err != nil {
		return nil, io_common.NewIoParseFailed("VRM GLBチャンクの解析に失敗しました", err)
	}
	logVrmStep("VRM読込ステップ: GLBチャンク解析完了 jsonBytes=%d binBytes=%d", len(jsonChunk), len(binChunk))

	doc := gltfDocument{}
	if err := json.Unmarshal(jsonChunk, &doc); err != nil {
		return nil, io_common.NewIoParseFailed("VRM JSONチャンクの解析に失敗しました", err)
	}
	logVrmStep(
		"VRM読込ステップ: JSON解析完了 nodes=%d meshes=%d primitives=%d accessors=%d",
		len(doc.Nodes),
		len(doc.Meshes),
		countGltfPrimitives(doc.Meshes),
		len(doc.Accessors),
	)

	parentIndexes, err := buildNodeParentIndexes(doc.Nodes)
	if err != nil {
		return nil, err
	}
	logVrmStep("VRM読込ステップ: ノード親子解析完了")

	worldPositions, err := buildNodeWorldPositions(doc.Nodes, parentIndexes)
	if err != nil {
		return nil, err
	}
	logVrmStep("VRM読込ステップ: ノードワールド座標計算完了")

	vrmData, humanBoneNames, err := buildVrmData(&doc, parentIndexes)
	if err != nil {
		return nil, err
	}
	version := ""
	profile := ""
	if vrmData != nil {
		version = string(vrmData.Version)
		profile = string(vrmData.Profile)
	}
	logVrmStep(
		"VRM読込ステップ: VRM拡張解析完了 version=%s profile=%s humanBones=%d",
		version,
		profile,
		len(humanBoneNames),
	)

	logVrmStep("VRM読込ステップ: PMX構築開始")
	modelData, err := buildPmxModel(path, &doc, binChunk, worldPositions, parentIndexes, humanBoneNames, vrmData, r.InferName(path))
	if err != nil {
		return nil, err
	}
	logVrmStep(
		"VRM読込ステップ: PMX構築完了 bones=%d vertices=%d faces=%d materials=%d",
		modelData.Bones.Len(),
		modelData.Vertices.Len(),
		modelData.Faces.Len(),
		modelData.Materials.Len(),
	)

	modelData.CreateDefaultDisplaySlots()
	info, err := os.Stat(path)
	if err != nil {
		return nil, io_common.NewIoParseFailed("VRMファイル情報の取得に失敗しました", err)
	}
	modelData.SetFileModTime(info.ModTime().UnixNano())
	modelData.UpdateHash()
	logVrmStep("VRM読込完了: file=%s hash=%s", loadTargetName, modelData.Hash())
	return modelData, nil
}

// countGltfPrimitives はglTF内のprimitive総数を返す。
func countGltfPrimitives(meshes []gltfMesh) int {
	total := 0
	for _, mesh := range meshes {
		total += len(mesh.Primitives)
	}
	return total
}

// logVrmStep はVRM変換の進捗ログを出力する。
func logVrmStep(format string, params ...any) {
	logger := logging.DefaultLogger()
	if logger == nil {
		return
	}
	logger.Info(format, params...)
}

// logVrmDebug はVRM変換のデバッグログを出力する。
func logVrmDebug(format string, params ...any) {
	logger := logging.DefaultLogger()
	if logger == nil {
		return
	}
	logger.Debug(format, params...)
}

// logVrmWarn はVRM変換の警告ログを出力する。
func logVrmWarn(format string, params ...any) {
	logger := logging.DefaultLogger()
	if logger == nil {
		return
	}
	logger.Warn(format, params...)
}

// gltfDocument はVRM読込時に必要なglTFトップレベル要素を表す。
type gltfDocument struct {
	Asset          gltfAsset                  `json:"asset"`
	Buffers        []gltfBuffer               `json:"buffers"`
	BufferViews    []gltfBufferView           `json:"bufferViews"`
	Accessors      []gltfAccessor             `json:"accessors"`
	Meshes         []gltfMesh                 `json:"meshes"`
	Skins          []gltfSkin                 `json:"skins"`
	Materials      []gltfMaterial             `json:"materials"`
	Textures       []gltfTexture              `json:"textures"`
	Images         []gltfImage                `json:"images"`
	ExtensionsUsed []string                   `json:"extensionsUsed"`
	Nodes          []gltfNode                 `json:"nodes"`
	Extensions     map[string]json.RawMessage `json:"extensions"`
	Scenes         []gltfScene                `json:"scenes"`
	Scene          int                        `json:"scene"`
}

// gltfAsset はglTF asset要素を表す。
type gltfAsset struct {
	Version   string `json:"version"`
	Generator string `json:"generator"`
}

// gltfScene はglTF scene要素を表す。
type gltfScene struct {
	Nodes []int `json:"nodes"`
}

// gltfNode はglTF node要素を表す。
type gltfNode struct {
	Name        string    `json:"name"`
	Mesh        *int      `json:"mesh"`
	Skin        *int      `json:"skin"`
	Children    []int     `json:"children"`
	Matrix      []float64 `json:"matrix"`
	Translation []float64 `json:"translation"`
	Rotation    []float64 `json:"rotation"`
	Scale       []float64 `json:"scale"`
}

// gltfBuffer はglTF buffer要素を表す。
type gltfBuffer struct {
	ByteLength int `json:"byteLength"`
}

// gltfBufferView はglTF bufferView要素を表す。
type gltfBufferView struct {
	Buffer     int `json:"buffer"`
	ByteOffset int `json:"byteOffset"`
	ByteLength int `json:"byteLength"`
	ByteStride int `json:"byteStride"`
}

// gltfAccessor はglTF accessor要素を表す。
type gltfAccessor struct {
	BufferView    *int   `json:"bufferView"`
	ByteOffset    int    `json:"byteOffset"`
	ComponentType int    `json:"componentType"`
	Count         int    `json:"count"`
	Type          string `json:"type"`
	Normalized    bool   `json:"normalized"`
}

// gltfMesh はglTF mesh要素を表す。
type gltfMesh struct {
	Name       string          `json:"name"`
	Primitives []gltfPrimitive `json:"primitives"`
}

// gltfPrimitive はglTF mesh primitive要素を表す。
type gltfPrimitive struct {
	Attributes map[string]int   `json:"attributes"`
	Indices    *int             `json:"indices"`
	Material   *int             `json:"material"`
	Mode       *int             `json:"mode"`
	Targets    []map[string]int `json:"targets"`
}

// gltfSkin はglTF skin要素を表す。
type gltfSkin struct {
	Joints []int `json:"joints"`
}

// gltfMaterial はglTF material要素を表す。
type gltfMaterial struct {
	Name                 string                   `json:"name"`
	DoubleSided          bool                     `json:"doubleSided"`
	PbrMetallicRoughness gltfPbrMetallicRoughness `json:"pbrMetallicRoughness"`
}

// gltfPbrMetallicRoughness はPBR基本材質情報を表す。
type gltfPbrMetallicRoughness struct {
	BaseColorFactor  []float64       `json:"baseColorFactor"`
	BaseColorTexture *gltfTextureRef `json:"baseColorTexture"`
}

// gltfTextureRef は材質から参照されるテクスチャ参照を表す。
type gltfTextureRef struct {
	Index int `json:"index"`
}

// gltfTexture はglTF texture要素を表す。
type gltfTexture struct {
	Source *int `json:"source"`
}

// gltfImage はglTF image要素を表す。
type gltfImage struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

// vrm0Extension はVRM0拡張の必要要素を表す。
type vrm0Extension struct {
	ExporterVersion string       `json:"exporterVersion"`
	Meta            vrm0Meta     `json:"meta"`
	Humanoid        vrm0Humanoid `json:"humanoid"`
}

// vrm0Meta はVRM0 meta要素を表す。
type vrm0Meta struct {
	Title   string `json:"title"`
	Version string `json:"version"`
	Author  string `json:"author"`
}

// vrm0Humanoid はVRM0 humanoid要素を表す。
type vrm0Humanoid struct {
	HumanBones []vrm0HumanBone `json:"humanBones"`
}

// vrm0HumanBone はVRM0 humanBones要素を表す。
type vrm0HumanBone struct {
	Bone string `json:"bone"`
	Node int    `json:"node"`
}

// vrm1Extension はVRM1拡張の必要要素を表す。
type vrm1Extension struct {
	SpecVersion string       `json:"specVersion"`
	Meta        vrm1Meta     `json:"meta"`
	Humanoid    vrm1Humanoid `json:"humanoid"`
}

// vrm1Meta はVRM1 meta要素を表す。
type vrm1Meta struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Authors []string `json:"authors"`
}

// vrm1Humanoid はVRM1 humanoid要素を表す。
type vrm1Humanoid struct {
	HumanBones map[string]vrm1HumanBone `json:"humanBones"`
}

// vrm1HumanBone はVRM1 humanBones要素を表す。
type vrm1HumanBone struct {
	Node *int `json:"node"`
}

// selectedHumanBone はPMX標準名へ採用するhumanoid情報を表す。
type selectedHumanBone struct {
	NodeIndex int
	Priority  int
}

// parseGLBJSONChunk はGLBバイナリからJSONチャンクを取り出す。
func parseGLBJSONChunk(b []byte) ([]byte, error) {
	if len(b) < glbMinValidLength {
		return nil, io_common.NewIoParseFailed("VRMヘッダが不足しています", nil)
	}
	magic := binary.LittleEndian.Uint32(b[0:4])
	if magic != glbMagic {
		return nil, io_common.NewIoParseFailed("GLBマジックが不正です", nil)
	}
	version := binary.LittleEndian.Uint32(b[4:8])
	if version != 2 {
		return nil, io_common.NewIoFormatNotSupported("GLBバージョンが未対応です: %d", nil, version)
	}
	totalLength := binary.LittleEndian.Uint32(b[8:12])
	if totalLength > uint32(len(b)) {
		return nil, io_common.NewIoParseFailed("GLB全体長が不正です", nil)
	}

	offset := glbHeaderLength
	for offset+glbChunkHeadSize <= len(b) {
		chunkLength := int(binary.LittleEndian.Uint32(b[offset : offset+4]))
		chunkType := binary.LittleEndian.Uint32(b[offset+4 : offset+8])
		chunkStart := offset + glbChunkHeadSize
		chunkEnd := chunkStart + chunkLength
		if chunkLength < 0 || chunkEnd > len(b) {
			return nil, io_common.NewIoParseFailed("GLBチャンク長が不正です", nil)
		}
		if chunkType == glbJSONChunkType {
			return b[chunkStart:chunkEnd], nil
		}
		offset = chunkEnd
	}
	return nil, io_common.NewIoParseFailed("GLB JSONチャンクが見つかりません", nil)
}

// buildNodeParentIndexes はnode配列から親インデックス配列を生成する。
func buildNodeParentIndexes(nodes []gltfNode) ([]int, error) {
	parentIndexes := make([]int, len(nodes))
	for i := range parentIndexes {
		parentIndexes[i] = -1
	}
	for parentIndex, node := range nodes {
		for _, childIndex := range node.Children {
			if childIndex < 0 || childIndex >= len(nodes) {
				return nil, io_common.NewIoParseFailed("node.children のindexが不正です: %d", nil, childIndex)
			}
			if parentIndexes[childIndex] == -1 {
				parentIndexes[childIndex] = parentIndex
			}
		}
	}
	return parentIndexes, nil
}

// buildNodeWorldPositions はnodeのローカル変換からワールド座標を算出する。
func buildNodeWorldPositions(nodes []gltfNode, parents []int) ([]mmath.Vec3, error) {
	worldMats := make([]mmath.Mat4, len(nodes))
	worldPositions := make([]mmath.Vec3, len(nodes))
	state := make([]int, len(nodes))

	for i := range nodes {
		if err := resolveNodeWorldMatrix(nodes, parents, i, state, worldMats, worldPositions); err != nil {
			return nil, err
		}
	}
	return worldPositions, nil
}

// resolveNodeWorldMatrix はnodeのワールド行列を再帰的に解決する。
func resolveNodeWorldMatrix(
	nodes []gltfNode,
	parents []int,
	nodeIndex int,
	state []int,
	worldMats []mmath.Mat4,
	worldPositions []mmath.Vec3,
) error {
	if nodeIndex < 0 || nodeIndex >= len(nodes) {
		return io_common.NewIoParseFailed("node index が不正です: %d", nil, nodeIndex)
	}
	if state[nodeIndex] == 2 {
		return nil
	}
	if state[nodeIndex] == 1 {
		return io_common.NewIoParseFailed("node親子関係に循環があります: %d", nil, nodeIndex)
	}
	state[nodeIndex] = 1
	local, err := nodeLocalMatrix(nodes[nodeIndex])
	if err != nil {
		return err
	}
	parentIndex := parents[nodeIndex]
	if parentIndex >= 0 {
		if err := resolveNodeWorldMatrix(nodes, parents, parentIndex, state, worldMats, worldPositions); err != nil {
			return err
		}
		worldMats[nodeIndex] = worldMats[parentIndex].Muled(local)
	} else {
		worldMats[nodeIndex] = local
	}
	worldPositions[nodeIndex] = worldMats[nodeIndex].Translation()
	state[nodeIndex] = 2
	return nil
}

// nodeLocalMatrix はnode要素からローカル行列を生成する。
func nodeLocalMatrix(node gltfNode) (mmath.Mat4, error) {
	if len(node.Matrix) > 0 {
		if len(node.Matrix) != 16 {
			return mmath.NewMat4(), io_common.NewIoParseFailed("node.matrix の要素数が不正です: %d", nil, len(node.Matrix))
		}
		mat := mmath.NewMat4()
		for i := 0; i < 16; i++ {
			mat[i] = node.Matrix[i]
		}
		return mat, nil
	}

	translation, err := parseVec3(node.Translation, mmath.ZERO_VEC3, "node.translation")
	if err != nil {
		return mmath.NewMat4(), err
	}
	scale, err := parseVec3(node.Scale, mmath.ONE_VEC3, "node.scale")
	if err != nil {
		return mmath.NewMat4(), err
	}
	rotation, err := parseQuaternion(node.Rotation)
	if err != nil {
		return mmath.NewMat4(), err
	}

	return translation.ToMat4().Muled(rotation.ToMat4()).Muled(scale.ToScaleMat4()), nil
}

// parseVec3 はスライスをVec3へ変換する。
func parseVec3(values []float64, defaultValue mmath.Vec3, label string) (mmath.Vec3, error) {
	if len(values) == 0 {
		return defaultValue, nil
	}
	if len(values) != 3 {
		return mmath.ZERO_VEC3, io_common.NewIoParseFailed("%s の要素数が不正です: %d", nil, label, len(values))
	}
	return mmath.Vec3{Vec: r3.Vec{X: values[0], Y: values[1], Z: values[2]}}, nil
}

// parseQuaternion はスライスをQuaternionへ変換する。
func parseQuaternion(values []float64) (mmath.Quaternion, error) {
	if len(values) == 0 {
		return mmath.NewQuaternion(), nil
	}
	if len(values) != 4 {
		return mmath.NewQuaternion(), io_common.NewIoParseFailed("node.rotation の要素数が不正です: %d", nil, len(values))
	}
	return mmath.NewQuaternionByValues(values[0], values[1], values[2], values[3]).Normalized(), nil
}

// buildVrmData はglTF文書からVrmDataとHumanoid採用名を構築する。
func buildVrmData(doc *gltfDocument, parents []int) (*vrm.VrmData, map[int]string, error) {
	version := detectVrmVersion(doc)
	if version == "" {
		return nil, nil, io_common.NewIoFormatNotSupported("VRM拡張が見つかりません", nil)
	}

	vrmData := vrm.NewVrmData()
	vrmData.Version = version
	vrmData.AssetGenerator = doc.Asset.Generator
	if doc.Extensions != nil {
		for k, raw := range doc.Extensions {
			vrmData.RawExtensions[k] = raw
		}
	}

	for i, node := range doc.Nodes {
		nodeData := vrm.NewNode(i)
		nodeData.Name = node.Name
		nodeData.ParentIndex = parents[i]
		nodeData.Children = append(nodeData.Children, node.Children...)
		if tr, err := parseVec3(node.Translation, mmath.ZERO_VEC3, "node.translation"); err == nil {
			nodeData.Translation = tr
		}
		vrmData.Nodes = append(vrmData.Nodes, *nodeData)
	}

	exporterVersion := ""
	humanoidMap := map[string]int{}
	if version == vrm.VRM_VERSION_1 {
		ext, err := parseVRM1Extension(doc.Extensions)
		if err != nil {
			return nil, nil, err
		}
		vrmData.Vrm1 = ext
		humanoidMap = extractHumanoidMapFromVRM1(ext)

		// VRM0拡張が同居している場合、作成元判定に exporterVersion を利用する。
		if ext0, err := parseVRM0Extension(doc.Extensions); err == nil && ext0 != nil {
			exporterVersion = ext0.ExporterVersion
		}
	} else {
		ext, err := parseVRM0Extension(doc.Extensions)
		if err != nil {
			return nil, nil, err
		}
		if ext == nil {
			return nil, nil, io_common.NewIoFormatNotSupported("VRM0拡張の解析に失敗しました", nil)
		}
		vrmData.Vrm0 = ext
		exporterVersion = ext.ExporterVersion
		humanoidMap = extractHumanoidMapFromVRM0(ext)
	}

	vrmData.Profile = detectProfile(doc.Asset.Generator, exporterVersion)
	return vrmData, selectHumanoidBoneNames(humanoidMap), nil
}

// parseVRM0Extension はextensionsからVRM0情報を抽出する。
func parseVRM0Extension(extensions map[string]json.RawMessage) (*vrm.Vrm0Data, error) {
	if extensions == nil {
		return nil, nil
	}
	raw, ok := extensions["VRM"]
	if !ok {
		return nil, nil
	}
	ext := vrm0Extension{}
	if err := json.Unmarshal(raw, &ext); err != nil {
		return nil, io_common.NewIoParseFailed("VRM0拡張のJSON解析に失敗しました", err)
	}
	data := vrm.NewVrm0Data()
	data.ExporterVersion = ext.ExporterVersion
	data.Meta = &vrm.Vrm0Meta{
		Title:   ext.Meta.Title,
		Version: ext.Meta.Version,
		Author:  ext.Meta.Author,
	}
	data.Humanoid = &vrm.Vrm0Humanoid{
		HumanBones: make([]vrm.Vrm0HumanBone, 0, len(ext.Humanoid.HumanBones)),
	}
	for _, b := range ext.Humanoid.HumanBones {
		data.Humanoid.HumanBones = append(data.Humanoid.HumanBones, vrm.Vrm0HumanBone{
			Bone: b.Bone,
			Node: b.Node,
		})
	}
	return data, nil
}

// parseVRM1Extension はextensionsからVRM1情報を抽出する。
func parseVRM1Extension(extensions map[string]json.RawMessage) (*vrm.Vrm1Data, error) {
	if extensions == nil {
		return nil, io_common.NewIoFormatNotSupported("VRM1拡張が存在しません", nil)
	}
	raw, ok := extensions["VRMC_vrm"]
	if !ok {
		return nil, io_common.NewIoFormatNotSupported("VRM1拡張が存在しません", nil)
	}
	ext := vrm1Extension{}
	if err := json.Unmarshal(raw, &ext); err != nil {
		return nil, io_common.NewIoParseFailed("VRM1拡張のJSON解析に失敗しました", err)
	}
	data := vrm.NewVrm1Data()
	data.SpecVersion = ext.SpecVersion
	data.Meta = &vrm.Vrm1Meta{
		Name:    ext.Meta.Name,
		Version: ext.Meta.Version,
		Authors: append([]string{}, ext.Meta.Authors...),
	}
	data.Humanoid = &vrm.Vrm1Humanoid{
		HumanBones: map[string]vrm.Vrm1HumanBone{},
	}
	for key, bone := range ext.Humanoid.HumanBones {
		if bone.Node == nil {
			continue
		}
		data.Humanoid.HumanBones[key] = vrm.Vrm1HumanBone{Node: *bone.Node}
	}
	return data, nil
}

// detectVrmVersion は拡張宣言から優先バージョンを判定する。
func detectVrmVersion(doc *gltfDocument) vrm.VrmVersion {
	hasVrm1 := containsIgnoreCase(doc.ExtensionsUsed, "VRMC_vrm")
	hasVrm0 := containsIgnoreCase(doc.ExtensionsUsed, "VRM")
	if doc.Extensions != nil {
		if _, ok := doc.Extensions["VRMC_vrm"]; ok {
			hasVrm1 = true
		}
		if _, ok := doc.Extensions["VRM"]; ok {
			hasVrm0 = true
		}
	}

	// VRM0/1 同時宣言時は VRM1 を優先する。
	if hasVrm1 {
		return vrm.VRM_VERSION_1
	}
	if hasVrm0 {
		return vrm.VRM_VERSION_0
	}
	return ""
}

// containsIgnoreCase は大文字小文字を無視して要素を検索する。
func containsIgnoreCase(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(value, target) {
			return true
		}
	}
	return false
}

// detectProfile は作成元情報からVRMプロファイルを判定する。
func detectProfile(assetGenerator string, exporterVersion string) vrm.VrmProfile {
	generatorLower := strings.ToLower(assetGenerator)
	exporterLower := strings.ToLower(exporterVersion)
	if strings.Contains(generatorLower, "vroid") || strings.Contains(exporterLower, "vroid") {
		return vrm.VRM_PROFILE_VROID
	}
	return vrm.VRM_PROFILE_STANDARD
}

// extractHumanoidMapFromVRM0 はVRM0 Humanoid情報をmapへ展開する。
func extractHumanoidMapFromVRM0(data *vrm.Vrm0Data) map[string]int {
	out := map[string]int{}
	if data == nil || data.Humanoid == nil {
		return out
	}
	for _, bone := range data.Humanoid.HumanBones {
		if bone.Node < 0 {
			continue
		}
		out[bone.Bone] = bone.Node
	}
	return out
}

// extractHumanoidMapFromVRM1 はVRM1 Humanoid情報をmapへ展開する。
func extractHumanoidMapFromVRM1(data *vrm.Vrm1Data) map[string]int {
	out := map[string]int{}
	if data == nil || data.Humanoid == nil {
		return out
	}
	for key, bone := range data.Humanoid.HumanBones {
		if bone.Node < 0 {
			continue
		}
		out[key] = bone.Node
	}
	return out
}

// selectHumanoidBoneNames はHumanoid情報からPMX標準名への採用ノードを決定する。
func selectHumanoidBoneNames(humanoid map[string]int) map[int]string {
	selectedByName := map[string]selectedHumanBone{}
	boneKeys := make([]string, 0, len(humanoid))
	for key := range humanoid {
		boneKeys = append(boneKeys, key)
	}
	sort.Strings(boneKeys)

	for _, boneKey := range boneKeys {
		nodeIndex := humanoid[boneKey]
		pmxName, ok := humanoidToPmxName(boneKey)
		if !ok {
			continue
		}
		priority := humanBonePriority(boneKey)
		if current, exists := selectedByName[pmxName]; exists {
			if priority < current.Priority {
				continue
			}
			if priority == current.Priority && nodeIndex >= current.NodeIndex {
				continue
			}
		}
		selectedByName[pmxName] = selectedHumanBone{
			NodeIndex: nodeIndex,
			Priority:  priority,
		}
	}

	selectedByNode := map[int]string{}
	for pmxName, selected := range selectedByName {
		selectedByNode[selected.NodeIndex] = pmxName
	}
	return selectedByNode
}

// humanoidToPmxName はHumanoid名をPMX標準名へ変換する。
func humanoidToPmxName(humanoidName string) (string, bool) {
	switch humanoidName {
	case "hips":
		return "下半身", true
	case "spine":
		return "上半身", true
	case "chest", "upperChest":
		return "上半身2", true
	case "neck":
		return "首", true
	case "head":
		return "頭", true
	case "leftShoulder":
		return "左肩", true
	case "rightShoulder":
		return "右肩", true
	case "leftUpperArm":
		return "左腕", true
	case "rightUpperArm":
		return "右腕", true
	case "leftLowerArm":
		return "左ひじ", true
	case "rightLowerArm":
		return "右ひじ", true
	case "leftHand":
		return "左手首", true
	case "rightHand":
		return "右手首", true
	case "leftUpperLeg":
		return "左足", true
	case "rightUpperLeg":
		return "右足", true
	case "leftLowerLeg":
		return "左ひざ", true
	case "rightLowerLeg":
		return "右ひざ", true
	case "leftFoot":
		return "左足首", true
	case "rightFoot":
		return "右足首", true
	case "leftToes":
		return "左つま先", true
	case "rightToes":
		return "右つま先", true
	case "leftEye":
		return "左目", true
	case "rightEye":
		return "右目", true
	case "jaw":
		return "あご", true
	default:
		return "", false
	}
}

// humanBonePriority は同一PMX名の競合解決優先度を返す。
func humanBonePriority(humanBoneName string) int {
	switch humanBoneName {
	case "upperChest":
		return 10
	case "chest":
		return 5
	default:
		return 0
	}
}

// buildPmxModel はVRM解析結果からPMXモデルを構築する。
func buildPmxModel(
	path string,
	doc *gltfDocument,
	binChunk []byte,
	worldPositions []mmath.Vec3,
	parentIndexes []int,
	humanBoneNames map[int]string,
	vrmData *vrm.VrmData,
	inferredName string,
) (*model.PmxModel, error) {
	conversion := buildVrmConversion(vrmData)
	modelData := model.NewPmxModel()
	modelData.SetPath(path)
	modelData.SetName(inferredName)
	modelData.VrmData = vrmData

	nodeToBoneIndex := map[int]int{}
	usedNames := map[string]int{}
	for nodeIndex, node := range doc.Nodes {
		boneName := resolveNodeBoneName(nodeIndex, node.Name, humanBoneNames)
		boneName = ensureUniqueBoneName(boneName, usedNames)
		bone := &model.Bone{
			Position:         convertVrmPositionToPmx(worldPositions[nodeIndex], conversion),
			ParentIndex:      -1,
			TailIndex:        -1,
			EffectIndex:      -1,
			BoneFlag:         model.BONE_FLAG_CAN_ROTATE | model.BONE_FLAG_IS_VISIBLE | model.BONE_FLAG_CAN_MANIPULATE,
			TailPosition:     mmath.ZERO_VEC3,
			DisplaySlotIndex: 0,
			IsSystem:         false,
		}
		bone.SetName(boneName)
		nodeToBoneIndex[nodeIndex] = modelData.Bones.AppendRaw(bone)
	}

	for nodeIndex := range doc.Nodes {
		boneIndex := nodeToBoneIndex[nodeIndex]
		bone, err := modelData.Bones.Get(boneIndex)
		if err != nil {
			return nil, io_common.NewIoParseFailed("PMXボーンの親設定に失敗しました", err)
		}
		parentNodeIndex := parentIndexes[nodeIndex]
		if parentNodeIndex >= 0 {
			if parentBoneIndex, ok := nodeToBoneIndex[parentNodeIndex]; ok {
				bone.ParentIndex = parentBoneIndex
			}
		}
		if bone.ParentIndex < 0 {
			bone.BoneFlag |= model.BONE_FLAG_CAN_TRANSLATE
		}
		bone.Layer = boneLayer(modelData, bone.Index())
	}

	for nodeIndex, node := range doc.Nodes {
		boneIndex := nodeToBoneIndex[nodeIndex]
		bone, err := modelData.Bones.Get(boneIndex)
		if err != nil {
			return nil, io_common.NewIoParseFailed("PMXボーンの末端設定に失敗しました", err)
		}

		tailBoneIndex := findTailBoneIndex(node.Children, nodeToBoneIndex)
		if tailBoneIndex >= 0 {
			bone.TailIndex = tailBoneIndex
			bone.BoneFlag |= model.BONE_FLAG_TAIL_IS_BONE
			continue
		}

		bone.BoneFlag &^= model.BONE_FLAG_TAIL_IS_BONE
		parentBone := getParentBone(modelData, bone.ParentIndex)
		bone.TailPosition = generateTailOffset(bone, parentBone)
	}

	if err := appendMeshData(modelData, doc, binChunk, nodeToBoneIndex, conversion); err != nil {
		return nil, err
	}

	return modelData, nil
}

// resolveNodeBoneName はnode名とHumanoid対応からPMXボーン名を決定する。
func resolveNodeBoneName(nodeIndex int, nodeName string, humanBoneNames map[int]string) string {
	if name, ok := humanBoneNames[nodeIndex]; ok && name != "" {
		return name
	}
	trimmed := strings.TrimSpace(nodeName)
	if trimmed != "" {
		return trimmed
	}
	return fmt.Sprintf("node_%03d", nodeIndex)
}

// ensureUniqueBoneName は同名ボーンの重複を回避する。
func ensureUniqueBoneName(name string, used map[string]int) string {
	if used == nil {
		return name
	}
	if _, ok := used[name]; !ok {
		used[name] = 1
		return name
	}
	index := used[name]
	used[name] = index + 1
	return fmt.Sprintf("%s_%d", name, index)
}

// convertVrmPositionToPmx はOpenGL系の座標をPMX向けへ変換する。
func convertVrmPositionToPmx(v mmath.Vec3, conversion vrmConversion) mmath.Vec3 {
	return mmath.Vec3{
		Vec: r3.Vec{
			X: v.X * conversion.Scale * conversion.Axis.X,
			Y: v.Y * conversion.Scale * conversion.Axis.Y,
			Z: v.Z * conversion.Scale * conversion.Axis.Z,
		},
	}
}

// findTailBoneIndex は子node一覧から末端接続先のボーンindexを返す。
func findTailBoneIndex(children []int, nodeToBoneIndex map[int]int) int {
	for _, childNodeIndex := range children {
		if idx, ok := nodeToBoneIndex[childNodeIndex]; ok {
			return idx
		}
	}
	return -1
}

// getParentBone は親indexから親ボーンを取得する。
func getParentBone(modelData *model.PmxModel, parentIndex int) *model.Bone {
	if modelData == nil || modelData.Bones == nil || parentIndex < 0 {
		return nil
	}
	bone, err := modelData.Bones.Get(parentIndex)
	if err != nil {
		return nil
	}
	return bone
}

// generateTailOffset は子無しボーン向けにテールオフセットを算出する。
func generateTailOffset(bone *model.Bone, parentBone *model.Bone) mmath.Vec3 {
	if bone == nil || parentBone == nil {
		return mmath.Vec3{Vec: r3.Vec{Y: 0.1}}
	}
	direction := bone.Position.Subed(parentBone.Position)
	length := direction.Length()
	if length <= 0 {
		return mmath.Vec3{Vec: r3.Vec{Y: 0.1}}
	}
	return direction.Normalized().MuledScalar(length * 0.5)
}

// boneLayer は親を遡ってレイヤー深度を算出する。
func boneLayer(modelData *model.PmxModel, boneIndex int) int {
	if modelData == nil || modelData.Bones == nil || boneIndex < 0 {
		return 0
	}
	layer := 0
	currentIndex := boneIndex
	for {
		bone, err := modelData.Bones.Get(currentIndex)
		if err != nil || bone == nil || bone.ParentIndex < 0 {
			return layer
		}
		layer++
		currentIndex = bone.ParentIndex
		if layer > modelData.Bones.Len() {
			return layer
		}
	}
}
