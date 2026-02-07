// 指示: miu200521358
package vrm

import (
	"encoding/binary"
	"fmt"
	"math"
	"path/filepath"
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/vrm"
	"gonum.org/v1/gonum/spatial/r3"
)

const (
	gltfComponentTypeByte          = 5120
	gltfComponentTypeUnsignedByte  = 5121
	gltfComponentTypeShort         = 5122
	gltfComponentTypeUnsignedShort = 5123
	gltfComponentTypeUnsignedInt   = 5125
	gltfComponentTypeFloat         = 5126

	gltfPrimitiveModePoints        = 0
	gltfPrimitiveModeLines         = 1
	gltfPrimitiveModeLineLoop      = 2
	gltfPrimitiveModeLineStrip     = 3
	gltfPrimitiveModeTriangles     = 4
	gltfPrimitiveModeTriangleStrip = 5
	gltfPrimitiveModeTriangleFan   = 6

	vroidMeterScale = 12.5
)

// vrmConversion はVRM->PMX変換時の座標設定を表す。
type vrmConversion struct {
	Scale          float64
	Axis           mmath.Vec3
	ReverseWinding bool
}

// accessorReadPlan はaccessorの読み取り計画を表す。
type accessorReadPlan struct {
	Accessor      gltfAccessor
	ComponentSize int
	ComponentNum  int
	Stride        int
	BaseOffset    int
	ViewStart     int
	ViewEnd       int
}

// accessorValueCache はaccessor値の再読込みを抑止するキャッシュ。
type accessorValueCache struct {
	floatValues  map[int][][]float64
	intValues    map[int][][]int
	vertexRanges map[string]primitiveVertexRange
}

// primitiveVertexRange はprimitiveで生成した頂点範囲を表す。
type primitiveVertexRange struct {
	Start int
	Count int
}

// weightedBone はボーンウェイト計算用の一時構造体。
type weightedBone struct {
	BoneIndex int
	Weight    float64
}

// buildVrmConversion はVRMプロファイルに応じた座標変換設定を返す。
func buildVrmConversion(vrmData *vrm.VrmData) vrmConversion {
	conversion := vrmConversion{
		Scale: vroidMeterScale,
		Axis: mmath.Vec3{
			Vec: r3.Vec{X: -1.0, Y: 1.0, Z: 1.0},
		},
	}
	if vrmData != nil && vrmData.Profile == vrm.VRM_PROFILE_VROID {
		// VRoidはVRM0/VRM1で座標配置が異なるため、バージョンごとに軸変換を切り替える。
		// VRM0: 既存規約どおり X 反転（-1, 1, 1）
		// VRM1: VRoid Studio 1.x 出力の前向きを保つため標準軸（1, 1, -1）
		if vrmData.Version == vrm.VRM_VERSION_1 {
			conversion.Axis = mmath.Vec3{
				Vec: r3.Vec{X: 1.0, Y: 1.0, Z: -1.0},
			}
		}
	}
	conversion.ReverseWinding = conversion.Axis.X*conversion.Axis.Y*conversion.Axis.Z < 0
	return conversion
}

// convertVrmNormalToPmx は法線をVRM座標系からPMX座標系へ変換する。
func convertVrmNormalToPmx(v mmath.Vec3, conversion vrmConversion) mmath.Vec3 {
	converted := mmath.Vec3{
		Vec: r3.Vec{
			X: v.X * conversion.Axis.X,
			Y: v.Y * conversion.Axis.Y,
			Z: v.Z * conversion.Axis.Z,
		},
	}
	return converted.Normalized()
}

// appendMeshData はglTF mesh/primitives からPMXの頂点・面・材質を生成する。
func appendMeshData(
	modelData *model.PmxModel,
	doc *gltfDocument,
	binChunk []byte,
	nodeToBoneIndex map[int]int,
	conversion vrmConversion,
) error {
	if modelData == nil || doc == nil || len(doc.Meshes) == 0 {
		return nil
	}

	textureIndexesByImage := appendImageTextures(modelData, doc.Images)
	totalPrimitives := countGltfPrimitives(doc.Meshes)
	cache := newAccessorValueCache()
	primitiveStep := 0
	meshUniquePrimitiveIndex := map[int]map[string]int{}
	logVrmStep(
		"VRMメッシュ変換開始: nodes=%d meshes=%d primitives=%d textures=%d",
		len(doc.Nodes),
		len(doc.Meshes),
		totalPrimitives,
		len(textureIndexesByImage),
	)

	for nodeIndex, node := range doc.Nodes {
		if node.Mesh == nil {
			continue
		}
		meshIndex := *node.Mesh
		if meshIndex < 0 || meshIndex >= len(doc.Meshes) {
			return io_common.NewIoParseFailed("node.mesh のindexが不正です: %d", nil, meshIndex)
		}

		mesh := doc.Meshes[meshIndex]
		logVrmStep(
			"VRMメッシュ変換ステップ: node=%d mesh=%d primitives=%d",
			nodeIndex,
			meshIndex,
			len(mesh.Primitives),
		)
		seenByKey := meshUniquePrimitiveIndex[meshIndex]
		if seenByKey == nil {
			seenByKey = map[string]int{}
			meshUniquePrimitiveIndex[meshIndex] = seenByKey
		}
		for primitiveIndex, primitive := range mesh.Primitives {
			primitiveStep++
			primitiveName := resolvePrimitiveName(mesh, meshIndex, primitiveIndex)
			if shouldSkipPrimitiveForUnsupportedTargets(primitive, primitiveIndex, seenByKey) {
				logVrmStep(
					"VRMプリミティブ変換スキップ: step=%d/%d node=%d mesh=%d primitive=%d name=%s reason=%s",
					primitiveStep,
					totalPrimitives,
					nodeIndex,
					meshIndex,
					primitiveIndex,
					primitiveName,
					"morph targets未対応のため重複base primitiveを省略",
				)
				continue
			}
			beforeVertexCount := modelData.Vertices.Len()
			beforeFaceCount := modelData.Faces.Len()
			beforeMaterialCount := modelData.Materials.Len()
			logVrmDebug(
				"VRMプリミティブ変換開始: step=%d/%d node=%d mesh=%d primitive=%d name=%s",
				primitiveStep,
				totalPrimitives,
				nodeIndex,
				meshIndex,
				primitiveIndex,
				primitiveName,
			)
			if primitiveStep%100 == 0 {
				logVrmStep(
					"VRMプリミティブ変換開始: step=%d/%d node=%d mesh=%d primitive=%d name=%s",
					primitiveStep,
					totalPrimitives,
					nodeIndex,
					meshIndex,
					primitiveIndex,
					primitiveName,
				)
			}
			if err := appendPrimitiveMeshData(
				modelData,
				doc,
				binChunk,
				nodeIndex,
				node,
				primitive,
				primitiveName,
				nodeToBoneIndex,
				textureIndexesByImage,
				conversion,
				cache,
			); err != nil {
				return err
			}
			logVrmStep(
				"VRMプリミティブ変換完了: step=%d/%d node=%d mesh=%d primitive=%d name=%s addVertices=%d addFaces=%d addMaterials=%d",
				primitiveStep,
				totalPrimitives,
				nodeIndex,
				meshIndex,
				primitiveIndex,
				primitiveName,
				modelData.Vertices.Len()-beforeVertexCount,
				modelData.Faces.Len()-beforeFaceCount,
				modelData.Materials.Len()-beforeMaterialCount,
			)
		}
	}
	logVrmStep(
		"VRMメッシュ変換完了: vertices=%d faces=%d materials=%d textures=%d",
		modelData.Vertices.Len(),
		modelData.Faces.Len(),
		modelData.Materials.Len(),
		modelData.Textures.Len(),
	)
	return nil
}

// appendImageTextures はglTF image配列に対応するPMX Textureを作成する。
func appendImageTextures(modelData *model.PmxModel, images []gltfImage) []int {
	textureIndexes := make([]int, len(images))
	for imageIndex, image := range images {
		texture := model.NewTexture()
		texture.SetName(resolveImageTextureName(image, imageIndex))
		texture.SetValid(true)
		textureIndexes[imageIndex] = modelData.Textures.AppendRaw(texture)
	}
	return textureIndexes
}

// resolveImageTextureName はimage要素からPMX用テクスチャ名を決定する。
func resolveImageTextureName(image gltfImage, imageIndex int) string {
	uri := strings.TrimSpace(image.URI)
	if uri != "" && !strings.HasPrefix(uri, "data:") {
		return filepath.Base(filepath.FromSlash(uri))
	}
	name := strings.TrimSpace(image.Name)
	if name != "" {
		return name
	}
	return fmt.Sprintf("image_%03d", imageIndex)
}

// resolvePrimitiveName はメッシュ名とindexから材質名の基本名を生成する。
func resolvePrimitiveName(mesh gltfMesh, meshIndex int, primitiveIndex int) string {
	meshName := strings.TrimSpace(mesh.Name)
	if meshName == "" {
		meshName = fmt.Sprintf("mesh_%03d", meshIndex)
	}
	return fmt.Sprintf("%s_%03d", meshName, primitiveIndex)
}

// shouldSkipPrimitiveForUnsupportedTargets はmorph targets未対応時の重複primitiveを判定する。
func shouldSkipPrimitiveForUnsupportedTargets(
	primitive gltfPrimitive,
	primitiveIndex int,
	seenByKey map[string]int,
) bool {
	if seenByKey == nil {
		return false
	}
	key := primitiveBaseKey(primitive)
	if key == "" {
		return false
	}
	if firstIndex, exists := seenByKey[key]; exists {
		// targets未対応の間は同一base primitiveの重複展開を抑止する。
		return len(primitive.Targets) > 0 && firstIndex != primitiveIndex
	}
	seenByKey[key] = primitiveIndex
	return false
}

// primitiveBaseKey はtargetsを除くprimitiveの識別キーを返す。
func primitiveBaseKey(primitive gltfPrimitive) string {
	attributesKey := primitiveAttributesKey(primitive.Attributes)
	if attributesKey == "" {
		return ""
	}
	indices := -1
	if primitive.Indices != nil {
		indices = *primitive.Indices
	}
	material := -1
	if primitive.Material != nil {
		material = *primitive.Material
	}
	mode := gltfPrimitiveModeTriangles
	if primitive.Mode != nil {
		mode = *primitive.Mode
	}
	return fmt.Sprintf("attr=%s|idx=%d|mat=%d|mode=%d", attributesKey, indices, material, mode)
}

// primitiveAttributesKey はattribute mapの安定キーを生成する。
func primitiveAttributesKey(attributes map[string]int) string {
	if len(attributes) == 0 {
		return ""
	}
	keys := make([]string, 0, len(attributes))
	for key := range attributes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s:%d", key, attributes[key]))
	}
	return strings.Join(parts, ",")
}

// primitiveVertexKey は頂点配列を再利用するためのキーを返す。
func primitiveVertexKey(nodeIndex int, primitive gltfPrimitive) string {
	return fmt.Sprintf("node=%d|attrs=%s", nodeIndex, primitiveAttributesKey(primitive.Attributes))
}

// newAccessorValueCache はaccessorキャッシュを初期化する。
func newAccessorValueCache() *accessorValueCache {
	return &accessorValueCache{
		floatValues:  map[int][][]float64{},
		intValues:    map[int][][]int{},
		vertexRanges: map[string]primitiveVertexRange{},
	}
}

// readFloatValues はfloat accessor値をキャッシュ付きで返す。
func (c *accessorValueCache) readFloatValues(doc *gltfDocument, accessorIndex int, binChunk []byte) ([][]float64, error) {
	if c == nil {
		return readAccessorFloatValues(doc, accessorIndex, binChunk)
	}
	if values, ok := c.floatValues[accessorIndex]; ok {
		return values, nil
	}
	values, err := readAccessorFloatValues(doc, accessorIndex, binChunk)
	if err != nil {
		return nil, err
	}
	c.floatValues[accessorIndex] = values
	return values, nil
}

// readIntValues はint accessor値をキャッシュ付きで返す。
func (c *accessorValueCache) readIntValues(doc *gltfDocument, accessorIndex int, binChunk []byte) ([][]int, error) {
	if c == nil {
		return readAccessorIntValues(doc, accessorIndex, binChunk)
	}
	if values, ok := c.intValues[accessorIndex]; ok {
		return values, nil
	}
	values, err := readAccessorIntValues(doc, accessorIndex, binChunk)
	if err != nil {
		return nil, err
	}
	c.intValues[accessorIndex] = values
	return values, nil
}

// getVertexRange は頂点範囲キャッシュから開始indexを取得する。
func (c *accessorValueCache) getVertexRange(key string, expectedCount int) (int, bool) {
	if c == nil || strings.TrimSpace(key) == "" {
		return 0, false
	}
	r, ok := c.vertexRanges[key]
	if !ok {
		return 0, false
	}
	if r.Count != expectedCount {
		return 0, false
	}
	return r.Start, true
}

// setVertexRange は頂点範囲キャッシュを登録する。
func (c *accessorValueCache) setVertexRange(key string, start int, count int) {
	if c == nil || strings.TrimSpace(key) == "" || start < 0 || count <= 0 {
		return
	}
	c.vertexRanges[key] = primitiveVertexRange{
		Start: start,
		Count: count,
	}
}

// appendPrimitiveMeshData はprimitiveをPMX頂点・面・材質へ変換する。
func appendPrimitiveMeshData(
	modelData *model.PmxModel,
	doc *gltfDocument,
	binChunk []byte,
	nodeIndex int,
	node gltfNode,
	primitive gltfPrimitive,
	primitiveName string,
	nodeToBoneIndex map[int]int,
	textureIndexesByImage []int,
	conversion vrmConversion,
	cache *accessorValueCache,
) error {
	positionAccessor, ok := primitive.Attributes["POSITION"]
	if !ok {
		return io_common.NewIoParseFailed("mesh.primitive に POSITION がありません", nil)
	}

	positions, err := cache.readFloatValues(doc, positionAccessor, binChunk)
	if err != nil {
		return io_common.NewIoParseFailed("POSITION属性の読み取りに失敗しました(accessor=%d)", err, positionAccessor)
	}
	if len(positions) == 0 {
		return nil
	}

	normals, err := readOptionalFloatAttribute(doc, primitive.Attributes, "NORMAL", binChunk, cache)
	if err != nil {
		logVrmWarn(
			"VRM法線の読み取りに失敗したため既定法線で継続します: node=%d primitive=%s err=%s",
			nodeIndex,
			primitiveName,
			err.Error(),
		)
		normals = nil
	}
	uvs, err := readOptionalFloatAttribute(doc, primitive.Attributes, "TEXCOORD_0", binChunk, cache)
	if err != nil {
		return err
	}
	joints, err := readOptionalIntAttribute(doc, primitive.Attributes, "JOINTS_0", binChunk, cache)
	if err != nil {
		return err
	}
	weights, err := readOptionalFloatAttribute(doc, primitive.Attributes, "WEIGHTS_0", binChunk, cache)
	if err != nil {
		return err
	}

	indices, err := readPrimitiveIndices(doc, primitive, len(positions), binChunk, cache)
	if err != nil {
		return err
	}
	mode := gltfPrimitiveModeTriangles
	if primitive.Mode != nil {
		mode = *primitive.Mode
	}
	triangles := triangulateIndices(indices, mode)
	if len(triangles) == 0 {
		return nil
	}

	vertexStart, appendedVertices := appendOrReusePrimitiveVertices(
		modelData,
		doc,
		nodeIndex,
		node,
		primitive,
		positions,
		normals,
		uvs,
		joints,
		weights,
		nodeToBoneIndex,
		conversion,
		cache,
	)
	if appendedVertices == 0 {
		logVrmDebug(
			"VRM頂点再利用: node=%d primitive=%s vertexCount=%d",
			nodeIndex,
			primitiveName,
			len(positions),
		)
	}

	materialIndex := appendPrimitiveMaterial(modelData, doc, primitive, primitiveName, textureIndexesByImage, len(triangles)*3)
	for _, tri := range triangles {
		if tri[0] < 0 || tri[1] < 0 || tri[2] < 0 {
			continue
		}
		if tri[0] >= len(positions) || tri[1] >= len(positions) || tri[2] >= len(positions) {
			return io_common.NewIoParseFailed("indices が頂点数を超えています", nil)
		}

		v0 := vertexStart + tri[0]
		v1 := vertexStart + tri[1]
		v2 := vertexStart + tri[2]
		if conversion.ReverseWinding {
			v0, v2 = v2, v0
		}
		face := &model.Face{
			VertexIndexes: [3]int{v0, v1, v2},
		}
		modelData.Faces.AppendRaw(face)
		appendVertexMaterialIndex(modelData, v0, materialIndex)
		appendVertexMaterialIndex(modelData, v1, materialIndex)
		appendVertexMaterialIndex(modelData, v2, materialIndex)
	}
	return nil
}

// appendOrReusePrimitiveVertices はprimitive頂点を生成または既存頂点を再利用する。
func appendOrReusePrimitiveVertices(
	modelData *model.PmxModel,
	doc *gltfDocument,
	nodeIndex int,
	node gltfNode,
	primitive gltfPrimitive,
	positions [][]float64,
	normals [][]float64,
	uvs [][]float64,
	joints [][]int,
	weights [][]float64,
	nodeToBoneIndex map[int]int,
	conversion vrmConversion,
	cache *accessorValueCache,
) (int, int) {
	vertexKey := primitiveVertexKey(nodeIndex, primitive)
	if cachedStart, ok := cache.getVertexRange(vertexKey, len(positions)); ok {
		return cachedStart, 0
	}

	defaultBoneIndex := resolveDefaultBoneIndex(modelData, nodeToBoneIndex, nodeIndex)
	skinJoints := resolveSkinJoints(doc, node)
	vertexStart := modelData.Vertices.Len()
	for vertexIndex := 0; vertexIndex < len(positions); vertexIndex++ {
		position := toVec3(positions[vertexIndex], mmath.ZERO_VEC3)
		normal := mmath.Vec3{Vec: r3.Vec{X: 0, Y: 1, Z: 0}}
		if vertexIndex < len(normals) {
			normal = toVec3(normals[vertexIndex], normal)
		}
		uv := mmath.ZERO_VEC2
		if vertexIndex < len(uvs) {
			uv = toVec2(uvs[vertexIndex], mmath.ZERO_VEC2)
		}

		jointValues := []int{}
		if vertexIndex < len(joints) {
			jointValues = joints[vertexIndex]
		}
		weightValues := []float64{}
		if vertexIndex < len(weights) {
			weightValues = weights[vertexIndex]
		}
		deform := buildVertexDeform(defaultBoneIndex, jointValues, weightValues, skinJoints, nodeToBoneIndex)
		vertex := &model.Vertex{
			Position:        convertVrmPositionToPmx(position, conversion),
			Normal:          convertVrmNormalToPmx(normal, conversion),
			Uv:              uv,
			ExtendedUvs:     []mmath.Vec4{},
			DeformType:      deform.DeformType(),
			Deform:          deform,
			EdgeFactor:      1.0,
			MaterialIndexes: []int{},
		}
		modelData.Vertices.AppendRaw(vertex)
	}
	cache.setVertexRange(vertexKey, vertexStart, len(positions))
	return vertexStart, len(positions)
}

// appendVertexMaterialIndex は頂点に材質indexを追加する。
func appendVertexMaterialIndex(modelData *model.PmxModel, vertexIndex int, materialIndex int) {
	if modelData == nil || modelData.Vertices == nil || materialIndex < 0 {
		return
	}
	vertex, err := modelData.Vertices.Get(vertexIndex)
	if err != nil || vertex == nil {
		return
	}
	vertex.MaterialIndexes = append(vertex.MaterialIndexes, materialIndex)
}

// resolveDefaultBoneIndex は頂点デフォーム未解決時に使う既定ボーンindexを返す。
func resolveDefaultBoneIndex(modelData *model.PmxModel, nodeToBoneIndex map[int]int, nodeIndex int) int {
	if idx, ok := nodeToBoneIndex[nodeIndex]; ok && idx >= 0 {
		return idx
	}
	if modelData != nil && modelData.Bones != nil && modelData.Bones.Len() > 0 {
		return 0
	}
	return 0
}

// resolveSkinJoints はnode.skin からjoint node index配列を返す。
func resolveSkinJoints(doc *gltfDocument, node gltfNode) []int {
	if node.Skin == nil || doc == nil {
		return nil
	}
	skinIndex := *node.Skin
	if skinIndex < 0 || skinIndex >= len(doc.Skins) {
		return nil
	}
	return doc.Skins[skinIndex].Joints
}

// buildVertexDeform はJOINTS/WEIGHTSからPMXデフォームを生成する。
func buildVertexDeform(
	defaultBoneIndex int,
	joints []int,
	weights []float64,
	skinJoints []int,
	nodeToBoneIndex map[int]int,
) model.IDeform {
	weightByBone := map[int]float64{}
	maxCount := len(joints)
	if len(weights) < maxCount {
		maxCount = len(weights)
	}
	for i := 0; i < maxCount; i++ {
		weight := weights[i]
		if weight <= 0 {
			continue
		}
		jointIndex := joints[i]
		if jointIndex < 0 {
			continue
		}

		nodeIndex := jointIndex
		if len(skinJoints) > 0 {
			if jointIndex >= len(skinJoints) {
				continue
			}
			nodeIndex = skinJoints[jointIndex]
		}
		boneIndex, ok := nodeToBoneIndex[nodeIndex]
		if !ok || boneIndex < 0 {
			continue
		}
		weightByBone[boneIndex] += weight
	}

	if len(weightByBone) == 0 {
		return model.NewBdef1(defaultBoneIndex)
	}

	weightedBones := make([]weightedBone, 0, len(weightByBone))
	totalWeight := 0.0
	for boneIndex, weight := range weightByBone {
		if weight <= 0 {
			continue
		}
		weightedBones = append(weightedBones, weightedBone{
			BoneIndex: boneIndex,
			Weight:    weight,
		})
		totalWeight += weight
	}
	if totalWeight <= 0 || len(weightedBones) == 0 {
		return model.NewBdef1(defaultBoneIndex)
	}

	sort.Slice(weightedBones, func(i int, j int) bool {
		if weightedBones[i].Weight == weightedBones[j].Weight {
			return weightedBones[i].BoneIndex < weightedBones[j].BoneIndex
		}
		return weightedBones[i].Weight > weightedBones[j].Weight
	})

	if len(weightedBones) == 1 {
		return model.NewBdef1(weightedBones[0].BoneIndex)
	}

	if len(weightedBones) == 2 {
		weight0 := weightedBones[0].Weight / (weightedBones[0].Weight + weightedBones[1].Weight)
		return model.NewBdef2(weightedBones[0].BoneIndex, weightedBones[1].BoneIndex, weight0)
	}

	if len(weightedBones) > 4 {
		weightedBones = weightedBones[:4]
	}
	totalTopWeight := 0.0
	for _, wb := range weightedBones {
		totalTopWeight += wb.Weight
	}
	if totalTopWeight <= 0 {
		return model.NewBdef1(defaultBoneIndex)
	}

	indexes := [4]int{defaultBoneIndex, defaultBoneIndex, defaultBoneIndex, defaultBoneIndex}
	values := [4]float64{0, 0, 0, 0}
	for i := 0; i < len(weightedBones) && i < 4; i++ {
		indexes[i] = weightedBones[i].BoneIndex
		values[i] = weightedBones[i].Weight / totalTopWeight
	}
	return model.NewBdef4(indexes, values)
}

// appendPrimitiveMaterial はprimitive用材質を追加し、そのindexを返す。
func appendPrimitiveMaterial(
	modelData *model.PmxModel,
	doc *gltfDocument,
	primitive gltfPrimitive,
	primitiveName string,
	textureIndexesByImage []int,
	verticesCount int,
) int {
	material := model.NewMaterial()
	material.SetName(primitiveName)
	material.EnglishName = primitiveName
	material.Memo = "VRM primitive"
	material.Diffuse = mmath.Vec4{X: 1.0, Y: 1.0, Z: 1.0, W: 1.0}
	material.Specular = mmath.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 1.0}
	material.Ambient = mmath.Vec3{Vec: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}}
	material.Edge = mmath.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 1.0}
	material.EdgeSize = 1.0
	material.TextureFactor = mmath.ONE_VEC4
	material.SphereTextureFactor = mmath.ONE_VEC4
	material.ToonTextureFactor = mmath.ONE_VEC4
	material.DrawFlag = model.DRAW_FLAG_GROUND_SHADOW | model.DRAW_FLAG_DRAWING_ON_SELF_SHADOW_MAPS | model.DRAW_FLAG_DRAWING_SELF_SHADOWS
	material.VerticesCount = verticesCount
	material.TextureIndex = resolveMaterialTextureIndex(doc, primitive, textureIndexesByImage)

	if primitive.Material != nil {
		materialIndex := *primitive.Material
		if materialIndex >= 0 && materialIndex < len(doc.Materials) {
			sourceMaterial := doc.Materials[materialIndex]
			if sourceMaterial.Name != "" {
				material.SetName(sourceMaterial.Name)
				material.EnglishName = sourceMaterial.Name
			}
			if sourceMaterial.DoubleSided {
				material.DrawFlag |= model.DRAW_FLAG_DOUBLE_SIDED_DRAWING
			}
			if len(sourceMaterial.PbrMetallicRoughness.BaseColorFactor) == 4 {
				material.Diffuse = mmath.Vec4{
					X: sourceMaterial.PbrMetallicRoughness.BaseColorFactor[0],
					Y: sourceMaterial.PbrMetallicRoughness.BaseColorFactor[1],
					Z: sourceMaterial.PbrMetallicRoughness.BaseColorFactor[2],
					W: sourceMaterial.PbrMetallicRoughness.BaseColorFactor[3],
				}
			}
		}
	}

	return modelData.Materials.AppendRaw(material)
}

// resolveMaterialTextureIndex はbaseColorTextureからPMXテクスチャindexを解決する。
func resolveMaterialTextureIndex(doc *gltfDocument, primitive gltfPrimitive, textureIndexesByImage []int) int {
	if doc == nil || primitive.Material == nil {
		return -1
	}
	materialIndex := *primitive.Material
	if materialIndex < 0 || materialIndex >= len(doc.Materials) {
		return -1
	}
	baseTexture := doc.Materials[materialIndex].PbrMetallicRoughness.BaseColorTexture
	if baseTexture == nil {
		return -1
	}
	if baseTexture.Index < 0 || baseTexture.Index >= len(doc.Textures) {
		return -1
	}
	texture := doc.Textures[baseTexture.Index]
	if texture.Source == nil {
		return -1
	}
	imageIndex := *texture.Source
	if imageIndex < 0 || imageIndex >= len(textureIndexesByImage) {
		return -1
	}
	return textureIndexesByImage[imageIndex]
}

// triangulateIndices はprimitive modeに従って三角形indexへ変換する。
func triangulateIndices(indexes []int, mode int) [][3]int {
	switch mode {
	case gltfPrimitiveModeTriangles:
		out := make([][3]int, 0, len(indexes)/3)
		for i := 0; i+2 < len(indexes); i += 3 {
			out = append(out, [3]int{indexes[i], indexes[i+1], indexes[i+2]})
		}
		return out
	case gltfPrimitiveModeTriangleStrip:
		out := make([][3]int, 0, len(indexes)-2)
		for i := 0; i+2 < len(indexes); i++ {
			if i%2 == 0 {
				out = append(out, [3]int{indexes[i], indexes[i+1], indexes[i+2]})
			} else {
				out = append(out, [3]int{indexes[i+1], indexes[i], indexes[i+2]})
			}
		}
		return out
	case gltfPrimitiveModeTriangleFan:
		out := make([][3]int, 0, len(indexes)-2)
		for i := 1; i+1 < len(indexes); i++ {
			out = append(out, [3]int{indexes[0], indexes[i], indexes[i+1]})
		}
		return out
	case gltfPrimitiveModePoints, gltfPrimitiveModeLines, gltfPrimitiveModeLineLoop, gltfPrimitiveModeLineStrip:
		return [][3]int{}
	default:
		return [][3]int{}
	}
}

// readPrimitiveIndices はprimitiveのindex配列を返す。
func readPrimitiveIndices(
	doc *gltfDocument,
	primitive gltfPrimitive,
	vertexCount int,
	binChunk []byte,
	cache *accessorValueCache,
) ([]int, error) {
	if primitive.Indices == nil {
		indices := make([]int, vertexCount)
		for i := 0; i < vertexCount; i++ {
			indices[i] = i
		}
		return indices, nil
	}
	accessorIndex := *primitive.Indices
	values, err := cache.readIntValues(doc, accessorIndex, binChunk)
	if err != nil {
		return nil, io_common.NewIoParseFailed("indices の読み取りに失敗しました(accessor=%d)", err, accessorIndex)
	}
	indices := make([]int, 0, len(values))
	for _, row := range values {
		if len(row) == 0 {
			continue
		}
		indices = append(indices, row[0])
	}
	return indices, nil
}

// readOptionalFloatAttribute は任意属性accessorを読み取る。
func readOptionalFloatAttribute(
	doc *gltfDocument,
	attributes map[string]int,
	key string,
	binChunk []byte,
	cache *accessorValueCache,
) ([][]float64, error) {
	if attributes == nil {
		return nil, nil
	}
	accessorIndex, ok := attributes[key]
	if !ok {
		return nil, nil
	}
	values, err := cache.readFloatValues(doc, accessorIndex, binChunk)
	if err != nil {
		return nil, io_common.NewIoParseFailed("%s属性の読み取りに失敗しました(accessor=%d)", err, key, accessorIndex)
	}
	return values, nil
}

// readOptionalIntAttribute は任意属性accessorを読み取る。
func readOptionalIntAttribute(
	doc *gltfDocument,
	attributes map[string]int,
	key string,
	binChunk []byte,
	cache *accessorValueCache,
) ([][]int, error) {
	if attributes == nil {
		return nil, nil
	}
	accessorIndex, ok := attributes[key]
	if !ok {
		return nil, nil
	}
	values, err := cache.readIntValues(doc, accessorIndex, binChunk)
	if err != nil {
		return nil, io_common.NewIoParseFailed("%s属性の読み取りに失敗しました(accessor=%d)", err, key, accessorIndex)
	}
	return values, nil
}

// readAccessorFloatValues はaccessorをfloat値配列として読み取る。
func readAccessorFloatValues(doc *gltfDocument, accessorIndex int, binChunk []byte) ([][]float64, error) {
	plan, err := prepareAccessorRead(doc, accessorIndex, binChunk)
	if err != nil {
		return nil, err
	}
	values := make([][]float64, plan.Accessor.Count)
	for i := 0; i < plan.Accessor.Count; i++ {
		row := make([]float64, plan.ComponentNum)
		elementBase := plan.BaseOffset + i*plan.Stride
		for c := 0; c < plan.ComponentNum; c++ {
			value, readErr := readComponentAsFloat(plan.Accessor, binChunk, elementBase+c*plan.ComponentSize)
			if readErr != nil {
				return nil, readErr
			}
			row[c] = value
		}
		values[i] = row
	}
	return values, nil
}

// readAccessorIntValues はaccessorをint値配列として読み取る。
func readAccessorIntValues(doc *gltfDocument, accessorIndex int, binChunk []byte) ([][]int, error) {
	plan, err := prepareAccessorRead(doc, accessorIndex, binChunk)
	if err != nil {
		return nil, err
	}
	values := make([][]int, plan.Accessor.Count)
	for i := 0; i < plan.Accessor.Count; i++ {
		row := make([]int, plan.ComponentNum)
		elementBase := plan.BaseOffset + i*plan.Stride
		for c := 0; c < plan.ComponentNum; c++ {
			value, readErr := readComponentAsInt(plan.Accessor.ComponentType, binChunk, elementBase+c*plan.ComponentSize)
			if readErr != nil {
				return nil, readErr
			}
			row[c] = value
		}
		values[i] = row
	}
	return values, nil
}

// prepareAccessorRead はaccessor読み取りに必要な情報を検証して返す。
func prepareAccessorRead(doc *gltfDocument, accessorIndex int, binChunk []byte) (accessorReadPlan, error) {
	if doc == nil {
		return accessorReadPlan{}, io_common.NewIoParseFailed("gltf document が未設定です", nil)
	}
	if accessorIndex < 0 || accessorIndex >= len(doc.Accessors) {
		return accessorReadPlan{}, io_common.NewIoParseFailed("accessor index が不正です: %d", nil, accessorIndex)
	}
	accessor := doc.Accessors[accessorIndex]
	if accessor.BufferView == nil {
		return accessorReadPlan{}, io_common.NewIoParseFailed("sparse accessor は未対応です", nil)
	}
	if accessor.Count < 0 {
		return accessorReadPlan{}, io_common.NewIoParseFailed("accessor.count が不正です: %d", nil, accessor.Count)
	}

	viewIndex := *accessor.BufferView
	if viewIndex < 0 || viewIndex >= len(doc.BufferViews) {
		return accessorReadPlan{}, io_common.NewIoParseFailed("bufferView index が不正です: %d", nil, viewIndex)
	}
	view := doc.BufferViews[viewIndex]
	if view.Buffer != 0 {
		return accessorReadPlan{}, io_common.NewIoParseFailed("bufferView.buffer が未対応です: %d", nil, view.Buffer)
	}
	if view.ByteLength < 0 || view.ByteOffset < 0 {
		return accessorReadPlan{}, io_common.NewIoParseFailed("bufferView の byteOffset/byteLength が不正です", nil)
	}
	if view.ByteOffset+view.ByteLength > len(binChunk) {
		return accessorReadPlan{}, io_common.NewIoParseFailed("bufferView 範囲がBINチャンク外です", nil)
	}

	componentNum, err := accessorComponentNum(accessor.Type)
	if err != nil {
		return accessorReadPlan{}, err
	}
	componentSize, err := accessorComponentSize(accessor.ComponentType)
	if err != nil {
		return accessorReadPlan{}, err
	}
	elementSize := componentNum * componentSize
	stride := view.ByteStride
	if stride <= 0 {
		stride = elementSize
	}
	if stride < elementSize {
		return accessorReadPlan{}, io_common.NewIoParseFailed("bufferView.byteStride が要素サイズより小さいです", nil)
	}
	baseOffset := view.ByteOffset + accessor.ByteOffset
	if baseOffset < view.ByteOffset || baseOffset > view.ByteOffset+view.ByteLength {
		return accessorReadPlan{}, io_common.NewIoParseFailed("accessor.byteOffset が不正です", nil)
	}
	if accessor.Count > 0 {
		lastEnd := baseOffset + (accessor.Count-1)*stride + elementSize
		if lastEnd > view.ByteOffset+view.ByteLength || lastEnd > len(binChunk) {
			return accessorReadPlan{}, io_common.NewIoParseFailed("accessor 範囲がbufferViewを超えています", nil)
		}
	}

	return accessorReadPlan{
		Accessor:      accessor,
		ComponentSize: componentSize,
		ComponentNum:  componentNum,
		Stride:        stride,
		BaseOffset:    baseOffset,
		ViewStart:     view.ByteOffset,
		ViewEnd:       view.ByteOffset + view.ByteLength,
	}, nil
}

// accessorComponentNum はaccessor.typeから要素次元数を返す。
func accessorComponentNum(typeName string) (int, error) {
	switch typeName {
	case "SCALAR":
		return 1, nil
	case "VEC2":
		return 2, nil
	case "VEC3":
		return 3, nil
	case "VEC4":
		return 4, nil
	default:
		return 0, io_common.NewIoFormatNotSupported("accessor.type が未対応です: %s", nil, typeName)
	}
}

// accessorComponentSize はcomponentTypeのバイト幅を返す。
func accessorComponentSize(componentType int) (int, error) {
	switch componentType {
	case gltfComponentTypeByte, gltfComponentTypeUnsignedByte:
		return 1, nil
	case gltfComponentTypeShort, gltfComponentTypeUnsignedShort:
		return 2, nil
	case gltfComponentTypeUnsignedInt, gltfComponentTypeFloat:
		return 4, nil
	default:
		return 0, io_common.NewIoFormatNotSupported("accessor.componentType が未対応です: %d", nil, componentType)
	}
}

// readComponentAsFloat はcomponentTypeをfloat64へ変換する。
func readComponentAsFloat(accessor gltfAccessor, data []byte, offset int) (float64, error) {
	switch accessor.ComponentType {
	case gltfComponentTypeByte:
		value := float64(int8(data[offset]))
		if accessor.Normalized {
			return math.Max(value/127.0, -1.0), nil
		}
		return value, nil
	case gltfComponentTypeUnsignedByte:
		value := float64(data[offset])
		if accessor.Normalized {
			return value / 255.0, nil
		}
		return value, nil
	case gltfComponentTypeShort:
		value := float64(int16(binary.LittleEndian.Uint16(data[offset : offset+2])))
		if accessor.Normalized {
			return math.Max(value/32767.0, -1.0), nil
		}
		return value, nil
	case gltfComponentTypeUnsignedShort:
		value := float64(binary.LittleEndian.Uint16(data[offset : offset+2]))
		if accessor.Normalized {
			return value / 65535.0, nil
		}
		return value, nil
	case gltfComponentTypeUnsignedInt:
		value := float64(binary.LittleEndian.Uint32(data[offset : offset+4]))
		if accessor.Normalized {
			return value / 4294967295.0, nil
		}
		return value, nil
	case gltfComponentTypeFloat:
		return float64(math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))), nil
	default:
		return 0, io_common.NewIoFormatNotSupported("float componentType が未対応です: %d", nil, accessor.ComponentType)
	}
}

// readComponentAsInt はcomponentTypeをintへ変換する。
func readComponentAsInt(componentType int, data []byte, offset int) (int, error) {
	switch componentType {
	case gltfComponentTypeByte:
		return int(int8(data[offset])), nil
	case gltfComponentTypeUnsignedByte:
		return int(data[offset]), nil
	case gltfComponentTypeShort:
		return int(int16(binary.LittleEndian.Uint16(data[offset : offset+2]))), nil
	case gltfComponentTypeUnsignedShort:
		return int(binary.LittleEndian.Uint16(data[offset : offset+2])), nil
	case gltfComponentTypeUnsignedInt:
		return int(binary.LittleEndian.Uint32(data[offset : offset+4])), nil
	default:
		return 0, io_common.NewIoFormatNotSupported("int componentType が未対応です: %d", nil, componentType)
	}
}

// toVec3 はfloat配列をVec3へ変換する。
func toVec3(values []float64, defaultValue mmath.Vec3) mmath.Vec3 {
	if len(values) < 3 {
		return defaultValue
	}
	return mmath.Vec3{
		Vec: r3.Vec{
			X: values[0],
			Y: values[1],
			Z: values[2],
		},
	}
}

// toVec2 はfloat配列をVec2へ変換する。
func toVec2(values []float64, defaultValue mmath.Vec2) mmath.Vec2 {
	if len(values) < 2 {
		return defaultValue
	}
	return mmath.Vec2{
		X: values[0],
		Y: values[1],
	}
}
