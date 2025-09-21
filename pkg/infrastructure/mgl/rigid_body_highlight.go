package mgl

import (
	"math"
	"time"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/physics"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
)

const (
	highlightColorR = 1.0
	highlightColorG = 1.0
	highlightColorB = 0.0
	highlightColorA = 0.35
	epsilon         = 1e-6
)

type RigidBodyHighlighter struct {
	highlightBuffer         *VertexBufferHandle     // ハイライト用頂点バッファ
	highlightVertices       []float32               // ハイライト用頂点配列
	debugHover              *DebugRigidBodyHover    // デバッグ用ホバー情報
	debugHoverRigid         *physics.RigidBodyValue // デバッグ用ホバー剛体
	debugHoverStartTime     time.Time               // ハイライト開始時刻（自動クリア用）
	prevRigidBodyDebugState bool                    // 前回の剛体デバッグ状態
}

func NewRigidBodyHighlighter() *RigidBodyHighlighter {
	return &RigidBodyHighlighter{
		highlightVertices: make([]float32, 0),
	}
}

// DebugRigidBodyHover holds information about the rigid body under the debug cursor.
type DebugRigidBodyHover struct {
	RigidBody *pmx.RigidBody
	HitPoint  *mmath.MVec3
}

func (mp *RigidBodyHighlighter) DebugHoverInfo() *DebugRigidBodyHover {
	return mp.debugHover
}

func (mp *RigidBodyHighlighter) UpdateDebugHoverByRigidBody(modelIndex int, rb *physics.RigidBodyValue, enable bool) {
	mlog.V("ハイライト開始: enable=%v, rigidBody=%v", enable, rb != nil)

	if !enable || rb == nil {
		mlog.V("ハイライト無効またはrigidBodyがnil - クリア")
		mp.clearDebugHover()
		return
	}

	mp.debugHover = &DebugRigidBodyHover{
		RigidBody: rb.PmxRigidBody,
		HitPoint:  nil, // ヒット点は不明（レイキャストしていないため）
	}
	mp.debugHoverRigid = rb
	mp.debugHoverStartTime = time.Now() // タイマー開始

	mlog.V("ハイライト設定完了 - 頂点再構築開始")
	mp.rebuildHighlightVertices(rb)
}

func (mp *RigidBodyHighlighter) DrawDebugHighlight(shader rendering.IShader, isDrawRigidBodyFront bool) {
	mlog.V("ハイライト描画開始: 頂点数=%d, rigid=%v", len(mp.highlightVertices), mp.debugHoverRigid != nil)
	if len(mp.highlightVertices) == 0 || mp.debugHoverRigid == nil {
		mlog.V("頂点またはデバッグ剛体がない - 描画スキップ")
		return
	}

	program := shader.Program(rendering.ProgramTypePhysics)
	mlog.V("シェーダープログラム取得: program=%d", program)
	gl.UseProgram(program)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	if mp.highlightBuffer == nil {
		mlog.V("ハイライトバッファ初期化: 頂点数=%d", len(mp.highlightVertices))
		mp.highlightBuffer = NewBufferFactory().CreateDebugBuffer(gl.Ptr(&mp.highlightVertices[0]), len(mp.highlightVertices))
	}

	mp.highlightBuffer.Bind()
	mp.highlightBuffer.UpdateDebugBuffer(mp.highlightVertices)

	triangleCount := int32(len(mp.highlightVertices) / 7)
	mlog.V("描画実行: DrawArrays 三角形数=%d", triangleCount)
	gl.DrawArrays(gl.TRIANGLES, 0, triangleCount)

	mp.highlightBuffer.Unbind()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	gl.UseProgram(0)
	mlog.V("ハイライト描画完了")
}

func (mp *RigidBodyHighlighter) clearDebugHover() {
	mp.debugHover = nil
	mp.debugHoverRigid = nil
	if mp.highlightVertices != nil {
		mp.highlightVertices = mp.highlightVertices[:0]
	}
}

// CheckAndClearHighlightOnDebugChange は剛体デバッグ状態変更時にハイライトをクリアします
func (mp *RigidBodyHighlighter) CheckAndClearHighlightOnDebugChange(currentDebugState bool) {
	// 状態が変更された場合にハイライトをクリア
	if mp.prevRigidBodyDebugState != currentDebugState {
		if !currentDebugState {
			// デバッグが無効になった場合はハイライトをクリア
			mp.clearDebugHover()
		}
		mp.prevRigidBodyDebugState = currentDebugState
	}
}

// CheckAndClearExpiredHighlight は2秒経過したハイライトを自動的にクリアします
func (mp *RigidBodyHighlighter) CheckAndClearExpiredHighlight() {
	if mp.debugHover == nil {
		// ハイライトが設定されていない場合は何もしない
		return
	}

	// 2秒経過をチェック
	if time.Since(mp.debugHoverStartTime) >= 2*time.Second {
		mlog.V("ハイライト自動クリア: 2秒経過しました")
		mp.clearDebugHover()
	}
}

// rebuildHighlightVertices はPMXのサイズ情報とBulletの位置・向き情報を組み合わせてハイライト頂点を構築します
func (mp *RigidBodyHighlighter) rebuildHighlightVertices(rb *physics.RigidBodyValue) {
	mlog.V("ハイライト頂点再構築開始: rb=%v", rb != nil)
	mp.highlightVertices = mp.highlightVertices[:0]
	if rb == nil || rb.BtRigidBody == nil || mp.debugHover == nil || mp.debugHover.RigidBody == nil {
		mlog.V("剛体またはbtRigidBodyまたはPMX剛体情報がnil - 頂点再構築中止")
		return
	}

	// Bulletから位置・向きを取得
	transformIface := rb.BtRigidBody.GetWorldTransform()
	transform, ok := transformIface.(bt.BtTransform)
	mlog.V("Transform取得結果: ok=%v", ok)
	if !ok {
		mlog.V("Transform取得失敗 - 頂点再構築中止")
		return
	}

	mat := mgl32.Mat4{}
	transform.GetOpenGLMatrix(&mat[0])
	worldMat := newMMat4ByMgl(&mat)

	// PMXから形状種別とサイズを取得
	pmxRigidBody := mp.debugHover.RigidBody
	shapeType := pmxRigidBody.ShapeType
	shapeSize := pmxRigidBody.Size
	mlog.V("PMX形状情報: ShapeType=%d, Size=[%.3f, %.3f, %.3f]",
		shapeType, shapeSize.X, shapeSize.Y, shapeSize.Z)

	switch shapeType {
	case 0: // SHAPE_SPHERE
		mlog.V("Shape種別: Sphere (PMX ShapeType=0)")
		radius := math.Abs(float64(shapeSize.X))
		mp.appendSphereHighlight(worldMat, radius)
	case 1: // SHAPE_BOX
		mlog.V("Shape種別: Box (PMX ShapeType=1)")
		hx := math.Abs(float64(shapeSize.X))
		hy := math.Abs(float64(shapeSize.Y))
		hz := math.Abs(float64(shapeSize.Z))
		mp.appendBoxHighlightWithSize(worldMat, hx, hy, hz)
	case 2: // SHAPE_CAPSULE
		mlog.V("Shape種別: Capsule (PMX ShapeType=2)")
		radius := math.Abs(float64(shapeSize.X))
		halfHeight := math.Abs(float64(shapeSize.Y)) / 2.0 // PMXの高さを半分にする
		mp.appendCapsuleHighlight(worldMat, radius, halfHeight)
	default:
		mlog.V("Shape種別: 未対応 - PMX ShapeType=%d", shapeType)
		// 未対応の形状の場合はデフォルトのBox形状で描画
		mp.appendGenericBoxHighlight(worldMat, 0.5)
	}

	mlog.V("頂点生成完了: 頂点数=%d", len(mp.highlightVertices))
}

// NewMMat4ByMgl OpenGL座標系からMMD座標系に変換された行列を返します
func newMMat4ByMgl(mat *mgl32.Mat4) *mmath.MMat4 {
	mm := mmath.NewMMat4ByValues(
		float64(mat[0]), float64(-mat[1]), float64(-mat[2]), float64(mat[3]),
		float64(-mat[4]), float64(mat[5]), float64(mat[6]), float64(mat[7]),
		float64(-mat[8]), float64(mat[9]), float64(mat[10]), float64(mat[11]),
		float64(-mat[12]), float64(mat[13]), float64(mat[14]), float64(mat[15]),
	)
	return mm
}

// appendBoxHighlightWithSize は指定サイズのBox形状を描画します
func (mp *RigidBodyHighlighter) appendBoxHighlightWithSize(world *mmath.MMat4, hx, hy, hz float64) {
	corners := []*mmath.MVec3{
		{X: -hx, Y: -hy, Z: -hz},
		{X: hx, Y: -hy, Z: -hz},
		{X: hx, Y: hy, Z: -hz},
		{X: -hx, Y: hy, Z: -hz},
		{X: -hx, Y: -hy, Z: hz},
		{X: hx, Y: -hy, Z: hz},
		{X: hx, Y: hy, Z: hz},
		{X: -hx, Y: hy, Z: hz},
	}

	indices := []int{
		0, 1, 2, 0, 2, 3,
		4, 5, 6, 4, 6, 7,
		0, 1, 5, 0, 5, 4,
		2, 3, 7, 2, 7, 6,
		1, 2, 6, 1, 6, 5,
		0, 3, 7, 0, 7, 4,
	}

	mp.appendVertices(world, corners, indices)
}

// appendGenericBoxHighlight はデフォルトサイズのBox形状を描画します
func (mp *RigidBodyHighlighter) appendGenericBoxHighlight(world *mmath.MMat4, size float64) {
	hx, hy, hz := size, size, size

	corners := []*mmath.MVec3{
		{X: -hx, Y: -hy, Z: -hz},
		{X: hx, Y: -hy, Z: -hz},
		{X: hx, Y: hy, Z: -hz},
		{X: -hx, Y: hy, Z: -hz},
		{X: -hx, Y: -hy, Z: hz},
		{X: hx, Y: -hy, Z: hz},
		{X: hx, Y: hy, Z: hz},
		{X: -hx, Y: hy, Z: hz},
	}

	indices := []int{
		0, 1, 2, 0, 2, 3,
		4, 5, 6, 4, 6, 7,
		0, 1, 5, 0, 5, 4,
		2, 3, 7, 2, 7, 6,
		1, 2, 6, 1, 6, 5,
		0, 3, 7, 0, 7, 4,
	}

	mp.appendVertices(world, corners, indices)
}

func (mp *RigidBodyHighlighter) appendSphereHighlight(world *mmath.MMat4, radius float64) {
	latSegments := 12
	lonSegments := 24

	var vertices []*mmath.MVec3
	var indices []int

	for lat := 0; lat < latSegments; lat++ {
		theta1 := float64(lat) / float64(latSegments) * math.Pi
		theta2 := float64(lat+1) / float64(latSegments) * math.Pi

		y1 := radius * math.Cos(theta1)
		y2 := radius * math.Cos(theta2)
		r1 := radius * math.Sin(theta1)
		r2 := radius * math.Sin(theta2)

		for lon := 0; lon < lonSegments; lon++ {
			phi1 := float64(lon) / float64(lonSegments) * 2 * math.Pi
			phi2 := float64(lon+1) / float64(lonSegments) * 2 * math.Pi

			p1 := &mmath.MVec3{X: r1 * math.Cos(phi1), Y: y1, Z: r1 * math.Sin(phi1)}
			p2 := &mmath.MVec3{X: r2 * math.Cos(phi1), Y: y2, Z: r2 * math.Sin(phi1)}
			p3 := &mmath.MVec3{X: r2 * math.Cos(phi2), Y: y2, Z: r2 * math.Sin(phi2)}
			p4 := &mmath.MVec3{X: r1 * math.Cos(phi2), Y: y1, Z: r1 * math.Sin(phi2)}

			base := len(vertices)
			vertices = append(vertices, p1, p2, p3, p4)
			indices = append(indices,
				base, base+1, base+2,
				base, base+2, base+3,
			)
		}
	}

	mp.appendVertices(world, vertices, indices)
}

func (mp *RigidBodyHighlighter) appendCapsuleHighlight(world *mmath.MMat4, radius, halfHeight float64) {
	segments := 16
	capSegments := segments / 2

	var vertices []*mmath.MVec3
	var indices []int

	// Cylinder
	for i := 0; i < segments; i++ {
		phi1 := float64(i) / float64(segments) * 2 * math.Pi
		phi2 := float64(i+1) / float64(segments) * 2 * math.Pi

		top1 := &mmath.MVec3{X: radius * math.Cos(phi1), Y: halfHeight, Z: radius * math.Sin(phi1)}
		top2 := &mmath.MVec3{X: radius * math.Cos(phi2), Y: halfHeight, Z: radius * math.Sin(phi2)}
		bottom1 := &mmath.MVec3{X: radius * math.Cos(phi1), Y: -halfHeight, Z: radius * math.Sin(phi1)}
		bottom2 := &mmath.MVec3{X: radius * math.Cos(phi2), Y: -halfHeight, Z: radius * math.Sin(phi2)}

		base := len(vertices)
		vertices = append(vertices, top1, bottom1, top2, bottom2)
		indices = append(indices,
			base, base+1, base+2,
			base+1, base+3, base+2,
		)
	}

	// Top hemisphere
	centerTop := halfHeight
	for lat := 0; lat < capSegments; lat++ {
		theta1 := float64(lat) / float64(capSegments) * (math.Pi / 2)
		theta2 := float64(lat+1) / float64(capSegments) * (math.Pi / 2)

		y1 := centerTop + radius*math.Cos(theta1)
		y2 := centerTop + radius*math.Cos(theta2)
		r1 := radius * math.Sin(theta1)
		r2 := radius * math.Sin(theta2)

		for lon := 0; lon < segments; lon++ {
			phi1 := float64(lon) / float64(segments) * 2 * math.Pi
			phi2 := float64(lon+1) / float64(segments) * 2 * math.Pi

			p1 := &mmath.MVec3{X: r1 * math.Cos(phi1), Y: y1, Z: r1 * math.Sin(phi1)}
			p2 := &mmath.MVec3{X: r2 * math.Cos(phi1), Y: y2, Z: r2 * math.Sin(phi1)}
			p3 := &mmath.MVec3{X: r2 * math.Cos(phi2), Y: y2, Z: r2 * math.Sin(phi2)}
			p4 := &mmath.MVec3{X: r1 * math.Cos(phi2), Y: y1, Z: r1 * math.Sin(phi2)}

			base := len(vertices)
			vertices = append(vertices, p1, p2, p3, p4)
			indices = append(indices,
				base, base+1, base+2,
				base, base+2, base+3,
			)
		}
	}

	// Bottom hemisphere
	centerBottom := -halfHeight
	for lat := 0; lat < capSegments; lat++ {
		theta1 := float64(lat)/float64(capSegments)*(math.Pi/2) + math.Pi/2
		theta2 := float64(lat+1)/float64(capSegments)*(math.Pi/2) + math.Pi/2

		y1 := centerBottom + radius*math.Cos(theta1)
		y2 := centerBottom + radius*math.Cos(theta2)
		r1 := radius * math.Sin(theta1)
		r2 := radius * math.Sin(theta2)

		for lon := 0; lon < segments; lon++ {
			phi1 := float64(lon) / float64(segments) * 2 * math.Pi
			phi2 := float64(lon+1) / float64(segments) * 2 * math.Pi

			p1 := &mmath.MVec3{X: r1 * math.Cos(phi1), Y: y1, Z: r1 * math.Sin(phi1)}
			p2 := &mmath.MVec3{X: r2 * math.Cos(phi1), Y: y2, Z: r2 * math.Sin(phi1)}
			p3 := &mmath.MVec3{X: r2 * math.Cos(phi2), Y: y2, Z: r2 * math.Sin(phi2)}
			p4 := &mmath.MVec3{X: r1 * math.Cos(phi2), Y: y1, Z: r1 * math.Sin(phi2)}

			base := len(vertices)
			vertices = append(vertices, p1, p2, p3, p4)
			indices = append(indices,
				base, base+1, base+2,
				base, base+2, base+3,
			)
		}
	}

	mp.appendVertices(world, vertices, indices)
}

func (mp *RigidBodyHighlighter) appendVertices(world *mmath.MMat4, vertices []*mmath.MVec3, indices []int) {
	color := []float32{float32(highlightColorR), float32(highlightColorG), float32(highlightColorB), float32(highlightColorA)}
	for _, idx := range indices {
		if idx < 0 || idx >= len(vertices) {
			continue
		}
		local := vertices[idx]
		worldPos := world.MulVec3(&mmath.MVec3{X: local.X, Y: local.Y, Z: local.Z})
		mp.highlightVertices = append(mp.highlightVertices,
			-float32(worldPos.X), // X軸反転
			float32(worldPos.Y),
			float32(worldPos.Z),
			color[0], color[1], color[2], color[3],
		)
	}
}
