//go:build windows
// +build windows

// 指示: miu200521358
package render

import (
	"context"
	"runtime"
	"sync"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mgl"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

// ModelRenderer (旧 RenderModel) のバッファ初期化処理をまとめたファイルです。
// このファイルでは、頂点情報、法線、ボーンライン・ボーンポイント、選択頂点、カーソル位置、SSBO など
// モデル描画に必要な OpenGL バッファの生成・データ転送処理を行います。

// ModelRendererBufferData はOpenGLバッファ生成前のデータをまとめた構造体です。
type ModelRendererBufferData struct {
	Vertices               []float32
	NormalVertices         []float32
	SelectedVertexVertices []float32
	Faces                  []uint32
	NormalIndexes          []uint32
	BoneLines              []float32
	BoneLineFaces          []uint32
	BoneLineIndexes        []int
	BonePoints             []float32
	BonePointFaces         []uint32
	BonePointIndexes       []int
	SelectedVertexIndexes  []uint32
}

// ModelRendererBufferOptions はバッファ生成時に必要なデータ種別を指定します。
type ModelRendererBufferOptions struct {
	NeedNormal         bool
	NeedSelectedVertex bool
	NeedBoneDebug      bool
}

// PrepareModelRendererBufferData はバッファ生成前の頂点/面/ボーン情報を作成する。
func PrepareModelRendererBufferData(ctx context.Context, modelData *model.PmxModel, workerCount int) (*ModelRendererBufferData, error) {
	return PrepareModelRendererBufferDataWithOptions(ctx, modelData, workerCount, ModelRendererBufferOptions{
		NeedNormal:         true,
		NeedSelectedVertex: true,
		NeedBoneDebug:      true,
	})
}

// PrepareModelRendererBufferDataWithOptions はオプションに応じてバッファ生成前のデータを作成する。
func PrepareModelRendererBufferDataWithOptions(ctx context.Context, modelData *model.PmxModel, workerCount int, options ModelRendererBufferOptions) (*ModelRendererBufferData, error) {
	if modelData == nil {
		return &ModelRendererBufferData{}, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var wg sync.WaitGroup
	var loadErr error
	var errMu sync.Mutex

	setErr := func(err error) {
		if err == nil {
			return
		}
		errMu.Lock()
		if loadErr == nil {
			loadErr = err
		}
		errMu.Unlock()
	}

	var vertices []float32
	var normalVertices []float32
	var selectedVertexVertices []float32
	var faces []uint32
	var boneLines []float32
	var boneLineFaces []uint32
	var boneLineIndexes []int
	var bonePoints []float32
	var bonePointFaces []uint32
	var bonePointIndexes []int

	wg.Add(1)
	go func() {
		defer wg.Done()
		targets := vertexDataTargets{
			NeedVertices: true,
			NeedNormal:   options.NeedNormal,
			NeedSelected: options.NeedSelectedVertex,
		}
		v, n, s, err := createAllVertexData(ctx, modelData, workerCount, targets)
		if err != nil {
			setErr(err)
			return
		}
		vertices = v
		normalVertices = n
		selectedVertexVertices = s
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		f, err := createIndexesData(ctx, modelData, workerCount)
		if err != nil {
			setErr(err)
			return
		}
		faces = f
	}()

	if options.NeedBoneDebug {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ctx.Err() != nil {
				setErr(ctx.Err())
				return
			}
			boneLines, boneLineFaces, boneLineIndexes,
				bonePoints, bonePointFaces, bonePointIndexes = createAllBoneDebugData(modelData)
		}()
	}

	wg.Wait()

	if loadErr != nil {
		return nil, loadErr
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var normalIndexes []uint32
	if options.NeedNormal {
		normalIndexes = createNormalIndexesData(modelData)
	}
	var selectedVertexIndexes []uint32
	if options.NeedSelectedVertex {
		selectedVertexIndexes = createAllVertexIndexesData(modelData)
	}

	return &ModelRendererBufferData{
		Vertices:               vertices,
		NormalVertices:         normalVertices,
		SelectedVertexVertices: selectedVertexVertices,
		Faces:                  faces,
		NormalIndexes:          normalIndexes,
		BoneLines:              boneLines,
		BoneLineFaces:          boneLineFaces,
		BoneLineIndexes:        boneLineIndexes,
		BonePoints:             bonePoints,
		BonePointFaces:         bonePointFaces,
		BonePointIndexes:       bonePointIndexes,
		SelectedVertexIndexes:  selectedVertexIndexes,
	}, nil
}

// initializeBuffers は、モデルの頂点バッファおよび関連バッファを初期化します。
func (mr *ModelRenderer) initializeBuffers(factory *mgl.BufferFactory, modelData *model.PmxModel) {
	// 法線/選択/ボーンは必要時に生成するため、基本バッファのみ準備する。
	bufferData, err := PrepareModelRendererBufferDataWithOptions(
		context.Background(),
		modelData,
		0,
		ModelRendererBufferOptions{},
	)
	if err != nil {
		logging.DefaultLogger().Warn("描画バッファ準備に失敗しました: %v", err)
		return
	}
	mr.initializeBuffersWithData(factory, modelData, bufferData)
}

// initializeBuffersWithData は準備済みのバッファデータからOpenGLバッファを生成します。
func (mr *ModelRenderer) initializeBuffersWithData(factory *mgl.BufferFactory, modelData *model.PmxModel, bufferData *ModelRendererBufferData) {
	if mr == nil || factory == nil || modelData == nil || bufferData == nil {
		return
	}
	// ModelDrawer を初期化
	md := mr.ModelDrawer

	md.vertices = bufferData.Vertices
	md.faces = bufferData.Faces

	// GLオブジェクトの生成は並列化しない ------------
	var verticesPtr unsafe.Pointer
	if len(md.vertices) > 0 {
		verticesPtr = gl.Ptr(md.vertices)
	}
	mr.bufferHandle = factory.NewVertexBuffer(verticesPtr, len(md.vertices))

	if bufferData.NormalVertices != nil || bufferData.NormalIndexes != nil {
		md.normalVertices = bufferData.NormalVertices
		var normalPtr unsafe.Pointer
		if len(md.normalVertices) > 0 {
			normalPtr = gl.Ptr(md.normalVertices)
		}
		md.normalBufferHandle = factory.NewVertexBuffer(normalPtr, len(md.normalVertices))

		normalIndexes := bufferData.NormalIndexes
		md.normalIndexCount = len(normalIndexes)
		var normalIndexPtr unsafe.Pointer
		if len(normalIndexes) > 0 {
			normalIndexPtr = gl.Ptr(normalIndexes)
		}
		md.normalIbo = factory.NewIndexBuffer(normalIndexPtr, len(normalIndexes))
	} else {
		md.normalVertices = nil
		md.normalBufferHandle = nil
		md.normalIbo = nil
		md.normalIndexCount = 0
	}

	// ボーンラインバッファの設定
	if bufferData.BoneLines != nil || bufferData.BonePoints != nil {
		md.boneLineIndexes = bufferData.BoneLineIndexes
		md.boneLineCount = len(md.boneLineIndexes)
		if len(bufferData.BoneLines) > 0 {
			md.boneLineBufferHandle = factory.NewBoneBuffer(gl.Ptr(bufferData.BoneLines), len(bufferData.BoneLines))
		}
		if len(bufferData.BoneLineFaces) > 0 {
			md.boneLineIbo = factory.NewIndexBuffer(gl.Ptr(bufferData.BoneLineFaces), len(bufferData.BoneLineFaces))
		}

		// ボーンポイントバッファの設定
		md.bonePointIndexes = bufferData.BonePointIndexes
		md.bonePointCount = len(md.bonePointIndexes)
		if len(bufferData.BonePoints) > 0 {
			md.bonePointBufferHandle = factory.NewBoneBuffer(gl.Ptr(bufferData.BonePoints), len(bufferData.BonePoints))
		}
		if len(bufferData.BonePointFaces) > 0 {
			md.bonePointIbo = factory.NewIndexBuffer(gl.Ptr(bufferData.BonePointFaces), len(bufferData.BonePointFaces))
		}
	} else {
		md.boneLineIndexes = nil
		md.boneLineCount = 0
		md.boneLineBufferHandle = nil
		md.boneLineIbo = nil
		md.bonePointIndexes = nil
		md.bonePointCount = 0
		md.bonePointBufferHandle = nil
		md.bonePointIbo = nil
	}

	// 選択頂点バッファの設定
	if bufferData.SelectedVertexVertices != nil || bufferData.SelectedVertexIndexes != nil {
		var selectedPtr unsafe.Pointer
		if len(bufferData.SelectedVertexVertices) > 0 {
			selectedPtr = gl.Ptr(bufferData.SelectedVertexVertices)
		}
		md.selectedVertexBufferHandle = factory.NewVertexBuffer(selectedPtr, len(bufferData.SelectedVertexVertices))
		// 選択頂点用の元データ参照を保持して、VBOのベースコピーが破棄されないようにする。
		md.selectedVertexVertices = bufferData.SelectedVertexVertices

		// 選択頂点用のインデックスバッファ（すべての頂点のインデックス）
		md.selectedVertexCount = len(bufferData.SelectedVertexIndexes)
		var selectedIndexPtr unsafe.Pointer
		if len(bufferData.SelectedVertexIndexes) > 0 {
			selectedIndexPtr = gl.Ptr(bufferData.SelectedVertexIndexes)
		}
		md.selectedVertexIbo = factory.NewIndexBuffer(selectedIndexPtr, len(bufferData.SelectedVertexIndexes))
	} else {
		md.selectedVertexBufferHandle = nil
		md.selectedVertexIbo = nil
		md.selectedVertexCount = 0
		md.selectedVertexVertices = nil
	}

	// カーソル位置表示用バッファの初期化
	md.cursorPositionBufferHandle = factory.NewDebugBuffer(nil, 0)

	// SSBOの作成
	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, modelData.Vertices.Len()*4*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, ssbo)

	mr.ssbo = ssbo
}

// vertexDataTargets は頂点配列の生成対象を指定します。
type vertexDataTargets struct {
	NeedVertices bool
	NeedNormal   bool
	NeedSelected bool
}

// createAllVertexData は指定対象の頂点データを一括で生成します。
func createAllVertexData(ctx context.Context, modelData *model.PmxModel, workerCount int, targets vertexDataTargets) ([]float32, []float32, []float32, error) {
	if modelData == nil {
		return nil, nil, nil, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return nil, nil, nil, ctx.Err()
	}
	if !targets.NeedVertices && !targets.NeedNormal && !targets.NeedSelected {
		return nil, nil, nil, nil
	}

	vertexCount := modelData.Vertices.Len()
	var vertices []float32
	var normalVertices []float32
	var selectedVertices []float32
	if targets.NeedVertices {
		vertices = make([]float32, vertexCount*vertexDataSize)
	}
	if targets.NeedNormal {
		normalVertices = make([]float32, vertexCount*2*vertexDataSize)
	}
	if targets.NeedSelected {
		selectedVertices = make([]float32, vertexCount*vertexDataSize)
	}
	if vertexCount == 0 {
		return vertices, normalVertices, selectedVertices, nil
	}

	// CPU負荷が高いため、固定数ワーカーで分割する。
	workerCount = resolveWorkerCount(workerCount, vertexCount)
	batchSize := resolveBatchSize(vertexCount, workerCount, 4096)
	batches := (vertexCount + batchSize - 1) / batchSize
	values := modelData.Vertices.Values()

	var wg sync.WaitGroup
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			for b := workerIndex; b < batches; b += workerCount {
				if ctx.Err() != nil {
					return
				}
				start := b * batchSize
				end := min(start+batchSize, vertexCount)
				for i := start; i < end; i++ {
					vertex := values[i]
					if vertex == nil {
						continue
					}
					vertexBase := i * vertexDataSize
					normalBase := i * 2 * vertexDataSize
					selectedBase := i * vertexDataSize
					fillVertexData(vertices, vertexBase, normalVertices, normalBase, normalVertices, normalBase+vertexDataSize, selectedVertices, selectedBase, vertex)
				}
			}
		}(w)
	}
	wg.Wait()

	if ctx.Err() != nil {
		return nil, nil, nil, ctx.Err()
	}
	return vertices, normalVertices, selectedVertices, nil
}

// ensureNormalBuffers は法線描画用バッファを必要に応じて生成する。
func (mr *ModelRenderer) ensureNormalBuffers() bool {
	if mr == nil || mr.Model == nil {
		return false
	}
	if mr.normalBufferHandle != nil && mr.normalIbo != nil {
		return true
	}
	workerCount := max(1, runtime.GOMAXPROCS(0)-1)
	targets := vertexDataTargets{NeedNormal: true}
	_, normalVertices, _, err := createAllVertexData(context.Background(), mr.Model, workerCount, targets)
	if err != nil {
		logging.DefaultLogger().Warn("法線描画バッファ生成に失敗しました: %v", err)
		return false
	}
	normalIndexes := createNormalIndexesData(mr.Model)
	factory := mgl.NewBufferFactory()
	var normalPtr unsafe.Pointer
	if len(normalVertices) > 0 {
		normalPtr = gl.Ptr(normalVertices)
	}
	mr.normalVertices = normalVertices
	mr.normalBufferHandle = factory.NewVertexBuffer(normalPtr, len(normalVertices))
	mr.normalIndexCount = len(normalIndexes)
	var normalIndexPtr unsafe.Pointer
	if len(normalIndexes) > 0 {
		normalIndexPtr = gl.Ptr(normalIndexes)
	}
	mr.normalIbo = factory.NewIndexBuffer(normalIndexPtr, len(normalIndexes))
	return true
}

// ensureSelectedVertexBuffers は選択頂点描画用バッファを必要に応じて生成する。
func (mr *ModelRenderer) ensureSelectedVertexBuffers() bool {
	if mr == nil || mr.Model == nil {
		return false
	}
	if mr.selectedVertexBufferHandle != nil && mr.selectedVertexIbo != nil {
		return true
	}
	workerCount := max(1, runtime.GOMAXPROCS(0)-1)
	targets := vertexDataTargets{NeedSelected: true}
	_, _, selectedVertices, err := createAllVertexData(context.Background(), mr.Model, workerCount, targets)
	if err != nil {
		logging.DefaultLogger().Warn("選択頂点描画バッファ生成に失敗しました: %v", err)
		return false
	}
	selectedIndexes := createAllVertexIndexesData(mr.Model)
	factory := mgl.NewBufferFactory()
	var selectedPtr unsafe.Pointer
	if len(selectedVertices) > 0 {
		selectedPtr = gl.Ptr(selectedVertices)
	}
	mr.selectedVertexVertices = selectedVertices
	mr.selectedVertexBufferHandle = factory.NewVertexBuffer(selectedPtr, len(selectedVertices))
	mr.selectedVertexCount = len(selectedIndexes)
	var selectedIndexPtr unsafe.Pointer
	if len(selectedIndexes) > 0 {
		selectedIndexPtr = gl.Ptr(selectedIndexes)
	}
	mr.selectedVertexIbo = factory.NewIndexBuffer(selectedIndexPtr, len(selectedIndexes))
	return true
}

// ensureBoneBuffers はボーン描画用バッファを必要に応じて生成する。
func (mr *ModelRenderer) ensureBoneBuffers() bool {
	if mr == nil || mr.Model == nil {
		return false
	}
	if mr.boneLineBufferHandle != nil && mr.bonePointBufferHandle != nil {
		return true
	}
	boneLines, boneLineFaces, boneLineIndexes,
		bonePoints, bonePointFaces, bonePointIndexes := createAllBoneDebugData(mr.Model)
	factory := mgl.NewBufferFactory()
	if len(boneLines) > 0 {
		mr.boneLineBufferHandle = factory.NewBoneBuffer(gl.Ptr(boneLines), len(boneLines))
	}
	if len(boneLineFaces) > 0 {
		mr.boneLineIbo = factory.NewIndexBuffer(gl.Ptr(boneLineFaces), len(boneLineFaces))
	}
	if len(bonePoints) > 0 {
		mr.bonePointBufferHandle = factory.NewBoneBuffer(gl.Ptr(bonePoints), len(bonePoints))
	}
	if len(bonePointFaces) > 0 {
		mr.bonePointIbo = factory.NewIndexBuffer(gl.Ptr(bonePointFaces), len(bonePointFaces))
	}
	mr.boneLineIndexes = boneLineIndexes
	mr.boneLineCount = len(boneLineIndexes)
	mr.bonePointIndexes = bonePointIndexes
	mr.bonePointCount = len(bonePointIndexes)
	if mr.boneLineBufferHandle == nil || mr.bonePointBufferHandle == nil || mr.boneLineIbo == nil || mr.bonePointIbo == nil {
		return false
	}
	return true
}

// createAllBoneDebugData はボーンライン・ポイントデータを一括生成します
func createAllBoneDebugData(modelData *model.PmxModel) ([]float32, []uint32, []int, []float32, []uint32, []int) {
	// ボーン情報の並列処理
	boneLines := make([]float32, 0)
	boneLineFaces := make([]uint32, 0, modelData.Bones.Len()*2)
	boneLineIndexes := make([]int, modelData.Bones.Len()*2)

	bonePoints := make([]float32, 0)
	bonePointFaces := make([]uint32, modelData.Bones.Len())
	bonePointIndexes := make([]int, modelData.Bones.Len())

	n := 0
	for _, bone := range modelData.Bones.Values() {
		if bone == nil {
			continue
		}
		tailPos := calcTailPosition(bone, modelData.Bones)
		boneLines = append(boneLines, newBoneGl(bone)...)
		boneLines = append(boneLines, newTailBoneGl(bone, tailPos)...)
		boneLineFaces = append(boneLineFaces, uint32(n), uint32(n+1))
		boneLineIndexes[n] = bone.Index()
		boneLineIndexes[n+1] = bone.Index()

		bonePoints = append(bonePoints, newBoneGl(bone)...)
		bonePointFaces[bone.Index()] = uint32(bone.Index())
		bonePointIndexes[bone.Index()] = bone.Index()

		n += 2
	}

	return boneLines, boneLineFaces, boneLineIndexes, bonePoints, bonePointFaces, bonePointIndexes
}

// createIndexesData はインデックスデータを生成します。
func createIndexesData(ctx context.Context, modelData *model.PmxModel, workerCount int) ([]uint32, error) {
	if modelData == nil {
		return nil, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	faceCount := modelData.Faces.Len()
	faces := make([]uint32, faceCount*3)
	if faceCount == 0 {
		return faces, nil
	}

	// 面情報は大規模になりやすいため、ワーカー数を固定して処理する。
	workerCount = resolveWorkerCount(workerCount, faceCount)
	batchSize := resolveBatchSize(faceCount, workerCount, 8192)
	batches := (faceCount + batchSize - 1) / batchSize
	values := modelData.Faces.Values()

	var wg sync.WaitGroup
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			for b := workerIndex; b < batches; b += workerCount {
				if ctx.Err() != nil {
					return
				}
				start := b * batchSize
				end := min(start+batchSize, faceCount)
				for i := start; i < end; i++ {
					face := values[i]
					if face == nil {
						continue
					}
					base := i * 3
					vertices := face.VertexIndexes
					// 頂点の順序を反転（OpenGL用）
					faces[base] = uint32(vertices[2])
					faces[base+1] = uint32(vertices[1])
					faces[base+2] = uint32(vertices[0])
				}
			}
		}(w)
	}
	wg.Wait()

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return faces, nil
}

// createNormalIndexesData は法線インデックスデータを生成します
func createNormalIndexesData(modelData *model.PmxModel) []uint32 {
	normalFaces := make([]uint32, 0, modelData.Vertices.Len()*2)

	// 各頂点の開始位置と対応する法線位置を結ぶインデックスを生成
	for i := 0; i < modelData.Vertices.Len(); i++ {
		n := i * 2
		normalFaces = append(normalFaces, uint32(n), uint32(n+1))
	}

	return normalFaces
}

// createAllVertexIndexesData はすべての頂点インデックスデータを生成します
func createAllVertexIndexesData(modelData *model.PmxModel) []uint32 {
	indexes := make([]uint32, modelData.Vertices.Len())
	for i := 0; i < modelData.Vertices.Len(); i++ {
		indexes[i] = uint32(i)
	}
	return indexes
}

// fillVertexData は単一頂点のデータを各バッファへ書き込みます。
func fillVertexData(
	vertices []float32,
	vertexBase int,
	normalVertices []float32,
	normalStartBase int,
	normalEndVertices []float32,
	normalEndBase int,
	selectedVertices []float32,
	selectedBase int,
	vertex *model.Vertex,
) {
	needVertices := len(vertices) > 0
	needNormal := len(normalVertices) > 0
	needSelected := len(selectedVertices) > 0
	if !needVertices && !needNormal && !needSelected {
		return
	}
	p := mgl.NewGlVec3(&vertex.Position)
	n := mgl.NewGlVec3(&vertex.Normal)
	extUvX := float32(0)
	extUvY := float32(0)
	if len(vertex.ExtendedUvs) > 0 {
		extUvX = float32(vertex.ExtendedUvs[0].X)
		extUvY = float32(vertex.ExtendedUvs[0].Y)
	}
	deform := packDeform(vertex.Deform)
	sdefFlag := float32(mmath.BoolToInt(vertex.DeformType == model.SDEF))
	sdefC, sdefR0, sdefR1 := getSdefParams(vertex.Deform)

	if needVertices {
		writeVertexFields(vertices, vertexBase, p, n, float32(vertex.Uv.X), float32(vertex.Uv.Y), extUvX, extUvY, float32(vertex.EdgeFactor), deform, sdefFlag, sdefC, sdefR0, sdefR1)
	}
	if needNormal {
		writeVertexFields(normalVertices, normalStartBase, p, n, float32(vertex.Uv.X), float32(vertex.Uv.Y), extUvX, extUvY, float32(vertex.EdgeFactor), deform, sdefFlag, sdefC, sdefR0, sdefR1)

		normalScaled := vertex.Normal.MuledScalar(0.5)
		nScaled := mgl.NewGlVec3(&normalScaled)
		writeVertexFields(normalEndVertices, normalEndBase,
			mgl32.Vec3{p[0] + nScaled[0], p[1] + nScaled[1], p[2] + nScaled[2]},
			nScaled,
			0, 0, 0, 0, 0,
			deform, sdefFlag, sdefC, sdefR0, sdefR1,
		)
	}
	if needSelected {
		writeVertexFields(selectedVertices, selectedBase, p, n, -0.1, 0, 0, 0, 0, deform, sdefFlag, sdefC, sdefR0, sdefR1)
	}
}

// writeVertexFields は頂点配列へ共通属性を書き込みます。
func writeVertexFields(
	dst []float32,
	base int,
	position mgl32.Vec3,
	normal mgl32.Vec3,
	uvX float32,
	uvY float32,
	extUvX float32,
	extUvY float32,
	edgeFactor float32,
	deform [8]float32,
	sdefFlag float32,
	sdefC mgl32.Vec3,
	sdefR0 mgl32.Vec3,
	sdefR1 mgl32.Vec3,
) {
	if base < 0 || base+vertexDataSize > len(dst) {
		return
	}
	dst[base+0] = position[0]
	dst[base+1] = position[1]
	dst[base+2] = position[2]
	dst[base+3] = normal[0]
	dst[base+4] = normal[1]
	dst[base+5] = normal[2]
	dst[base+6] = uvX
	dst[base+7] = uvY
	dst[base+8] = extUvX
	dst[base+9] = extUvY
	dst[base+10] = edgeFactor
	dst[base+11] = deform[0]
	dst[base+12] = deform[1]
	dst[base+13] = deform[2]
	dst[base+14] = deform[3]
	dst[base+15] = deform[4]
	dst[base+16] = deform[5]
	dst[base+17] = deform[6]
	dst[base+18] = deform[7]
	dst[base+19] = sdefFlag
	dst[base+20] = sdefC[0]
	dst[base+21] = sdefC[1]
	dst[base+22] = sdefC[2]
	dst[base+23] = sdefR0[0]
	dst[base+24] = sdefR0[1]
	dst[base+25] = sdefR0[2]
	dst[base+26] = sdefR1[0]
	dst[base+27] = sdefR1[1]
	dst[base+28] = sdefR1[2]
}

// resolveWorkerCount は利用可能コア数と処理量からワーカー数を決める。
func resolveWorkerCount(workerCount int, total int) int {
	if workerCount <= 0 {
		workerCount = runtime.GOMAXPROCS(0)
	}
	if workerCount <= 0 {
		workerCount = 1
	}
	// UI応答性を確保するため、最大でも1コア分は空ける。
	maxWorkers := max(1, runtime.GOMAXPROCS(0)-1)
	if workerCount > maxWorkers {
		workerCount = maxWorkers
	}
	if total > 0 && workerCount > total {
		workerCount = total
	}
	return workerCount
}

// resolveBatchSize はバッチサイズを決定する。
func resolveBatchSize(total int, workerCount int, defaultSize int) int {
	if defaultSize <= 0 {
		defaultSize = 1024
	}
	if total <= 0 {
		return defaultSize
	}
	if workerCount <= 0 {
		workerCount = 1
	}
	// ワーカー1つあたり約4バッチになるように調整し、細かい分割を抑える。
	target := total / (workerCount * 4)
	if target < defaultSize {
		return defaultSize
	}
	return target
}
