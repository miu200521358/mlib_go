// 指示: miu200521358
package motion

import (
	"math"
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// BoneFrame はボーンフレームを表す。
type BoneFrame struct {
	*BaseFrame
	Position           *mmath.Vec3
	Rotation           *mmath.Quaternion
	UnitRotation       *mmath.Quaternion
	Scale              *mmath.Vec3
	CancelablePosition *mmath.Vec3
	CancelableRotation *mmath.Quaternion
	CancelableScale    *mmath.Vec3
	Curves             *BoneCurves
	DisablePhysics     *bool
}

// NewBoneFrame はBoneFrameを生成する。
func NewBoneFrame(index Frame) *BoneFrame {
	return &BoneFrame{BaseFrame: NewBaseFrame(index)}
}

// Copy はフレームを複製する。
func (f *BoneFrame) Copy() (IBaseFrame, error) {
	if f == nil {
		return (*BoneFrame)(nil), nil
	}
	var curves *BoneCurves
	if f.Curves != nil {
		curves = f.Curves.Copy()
	}
	copied := &BoneFrame{
		BaseFrame:          &BaseFrame{index: f.Index(), Read: f.Read},
		Position:           copyVec3(f.Position),
		Rotation:           copyQuaternion(f.Rotation),
		UnitRotation:       copyQuaternion(f.UnitRotation),
		Scale:              copyVec3(f.Scale),
		CancelablePosition: copyVec3(f.CancelablePosition),
		CancelableRotation: copyQuaternion(f.CancelableRotation),
		CancelableScale:    copyVec3(f.CancelableScale),
		Curves:             curves,
		DisablePhysics:     copyBoolPtr(f.DisablePhysics),
	}
	return copied, nil
}

// lerpFrame は補間結果を返す。
func (next *BoneFrame) lerpFrame(prev *BoneFrame, index Frame) *BoneFrame {
	if prev == nil && next == nil {
		return NewBoneFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	if next == nil {
		return prev.copyWithIndex(index)
	}

	bf := NewBoneFrame(index)
	xy, yy, zy, ry := boneCurveT(prev.Index(), index, next.Index(), next.Curves)

	prevRot := quatOrIdent(prev.Rotation)
	nextRot := quatOrIdent(next.Rotation)
	rot := prevRot.Slerp(nextRot, ry)
	bf.Rotation = &rot

	prevUnit := quatOrIdent(prev.UnitRotation)
	nextUnit := quatOrIdent(next.UnitRotation)
	unit := prevUnit.Slerp(nextUnit, ry)
	bf.UnitRotation = &unit

	prevCancelRot := quatOrIdent(prev.CancelableRotation)
	nextCancelRot := quatOrIdent(next.CancelableRotation)
	cancelRot := prevCancelRot.Slerp(nextCancelRot, ry)
	bf.CancelableRotation = &cancelRot

	prevPos := vec3OrZero(prev.Position)
	nextPos := vec3OrZero(next.Position)
	prevCancelPos := vec3OrZero(prev.CancelablePosition)
	nextCancelPos := vec3OrZero(next.CancelablePosition)
	prevScale := vec3OrUnit(prev.Scale)
	nextScale := vec3OrUnit(next.Scale)
	prevCancelScale := vec3OrUnit(prev.CancelableScale)
	nextCancelScale := vec3OrUnit(next.CancelableScale)

	prevX := mmath.Vec4{X: prevPos.X, Y: prevCancelPos.X, Z: prevScale.X, W: prevCancelScale.X}
	nextX := mmath.Vec4{X: nextPos.X, Y: nextCancelPos.X, Z: nextScale.X, W: nextCancelScale.X}
	nowX := prevX.Lerp(nextX, xy)

	prevY := mmath.Vec4{X: prevPos.Y, Y: prevCancelPos.Y, Z: prevScale.Y, W: prevCancelScale.Y}
	nextY := mmath.Vec4{X: nextPos.Y, Y: nextCancelPos.Y, Z: nextScale.Y, W: nextCancelScale.Y}
	nowY := prevY.Lerp(nextY, yy)

	prevZ := mmath.Vec4{X: prevPos.Z, Y: prevCancelPos.Z, Z: prevScale.Z, W: prevCancelScale.Z}
	nextZ := mmath.Vec4{X: nextPos.Z, Y: nextCancelPos.Z, Z: nextScale.Z, W: nextCancelScale.Z}
	nowZ := prevZ.Lerp(nextZ, zy)

	bf.Position = vec3Ptr(nowX.X, nowY.X, nowZ.X)
	bf.CancelablePosition = vec3Ptr(nowX.Y, nowY.Y, nowZ.Y)
	bf.Scale = vec3Ptr(nowX.Z, nowY.Z, nowZ.Z)
	bf.CancelableScale = vec3Ptr(nowX.W, nowY.W, nowZ.W)

	bf.DisablePhysics = copyBoolPtr(prev.DisablePhysics)
	if bf.DisablePhysics == nil {
		bf.DisablePhysics = copyBoolPtr(next.DisablePhysics)
	}

	return bf
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *BoneFrame) copyWithIndex(index Frame) *BoneFrame {
	if f == nil {
		return nil
	}
	var curves *BoneCurves
	if f.Curves != nil {
		curves = f.Curves.Copy()
	}
	return &BoneFrame{
		BaseFrame:          &BaseFrame{index: index, Read: f.Read},
		Position:           copyVec3(f.Position),
		Rotation:           copyQuaternion(f.Rotation),
		UnitRotation:       copyQuaternion(f.UnitRotation),
		Scale:              copyVec3(f.Scale),
		CancelablePosition: copyVec3(f.CancelablePosition),
		CancelableRotation: copyQuaternion(f.CancelableRotation),
		CancelableScale:    copyVec3(f.CancelableScale),
		Curves:             curves,
		DisablePhysics:     copyBoolPtr(f.DisablePhysics),
	}
}

// splitCurve は曲線を分割する。
func (f *BoneFrame) splitCurve(prev *BoneFrame, next *BoneFrame, index Frame) {
	if f == nil || prev == nil || next == nil {
		return
	}
	if next.Curves == nil {
		f.Curves = NewBoneCurves()
		return
	}
	if f.Curves == nil {
		f.Curves = NewBoneCurves()
	}
	f.Curves.TranslateX, next.Curves.TranslateX =
		mmath.SplitCurve(next.Curves.TranslateX, float32(prev.Index()), float32(index), float32(next.Index()))
	f.Curves.TranslateY, next.Curves.TranslateY =
		mmath.SplitCurve(next.Curves.TranslateY, float32(prev.Index()), float32(index), float32(next.Index()))
	f.Curves.TranslateZ, next.Curves.TranslateZ =
		mmath.SplitCurve(next.Curves.TranslateZ, float32(prev.Index()), float32(index), float32(next.Index()))
	f.Curves.Rotate, next.Curves.Rotate =
		mmath.SplitCurve(next.Curves.Rotate, float32(prev.Index()), float32(index), float32(next.Index()))
}

// BoneNameFrames はボーン名ごとのフレーム集合を表す。
type BoneNameFrames struct {
	*BaseFrames[*BoneFrame]
	Name string
}

// NewBoneNameFrames はBoneNameFramesを生成する。
func NewBoneNameFrames(name string) *BoneNameFrames {
	return &BoneNameFrames{
		BaseFrames: NewBaseFrames(NewBoneFrame, nilBoneFrame),
		Name:       name,
	}
}

// Reduce は変曲点抽出と曲線当てはめで削減する。
func (b *BoneNameFrames) Reduce() *BoneNameFrames {
	if b == nil || b.Len() == 0 {
		return b
	}
	maxFrame := b.MaxFrame()
	maxIFrame := int(maxFrame) + 1
	if maxIFrame <= 1 {
		return b
	}

	frames := make([]Frame, 0, maxIFrame)
	xs := make([]float64, 0, maxIFrame)
	ys := make([]float64, 0, maxIFrame)
	zs := make([]float64, 0, maxIFrame)
	fixRs := make([]float64, 0, maxIFrame)
	quats := make([]mmath.Quaternion, 0, maxIFrame)

	for i := 0; i < maxIFrame; i++ {
		f := Frame(i)
		frames = append(frames, f)
		bf := b.Get(f)

		pos := vec3OrZero(nil)
		rot := mmath.NewQuaternion()
		if bf != nil {
			pos = vec3OrZero(bf.Position)
			rot = quatOrIdent(bf.Rotation)
		}
		xs = append(xs, pos.X)
		ys = append(ys, pos.Y)
		zs = append(zs, pos.Z)
		quats = append(quats, rot)
		fixRs = append(fixRs, mmath.NewQuaternion().Dot(rot))
	}

	inflectionFrames := make([]Frame, 0, b.Len())
	if !mmath.IsAllSameValues(xs) {
		inflectionFrames = append(inflectionFrames, findInflectionFrames(frames, xs, 1e-4)...)
	}
	if !mmath.IsAllSameValues(ys) {
		inflectionFrames = append(inflectionFrames, findInflectionFrames(frames, ys, 1e-4)...)
	}
	if !mmath.IsAllSameValues(zs) {
		inflectionFrames = append(inflectionFrames, findInflectionFrames(frames, zs, 1e-4)...)
	}
	if !mmath.IsAllSameValues(fixRs) {
		inflectionFrames = append(inflectionFrames, findInflectionFrames(frames, fixRs, 1e-6)...)
	}

	inflectionFrames = mmath.Unique(inflectionFrames)
	mmath.Sort(inflectionFrames)
	if len(inflectionFrames) <= 2 {
		return b
	}

	reduced := NewBoneNameFrames(b.Name)
	{
		bf := b.Get(inflectionFrames[0])
		reduceBf := NewBoneFrame(inflectionFrames[0])
		reduceBf.Position = copyVec3(bf.Position)
		reduceBf.Rotation = copyQuaternion(bf.Rotation)
		if bf.Curves != nil {
			reduceBf.Curves = bf.Curves.Copy()
		}
		reduced.Append(reduceBf)
	}

	startFrame := inflectionFrames[0]
	midFrame := inflectionFrames[1]
	endFrame := inflectionFrames[2]
	actualEnd := Frame(0)
	var i int
	for actualEnd < maxFrame {
		actualEnd = b.reduceRange(startFrame, midFrame, endFrame, xs, ys, zs, quats, reduced)

		exactI := slices.Index(inflectionFrames, actualEnd)
		if exactI == -1 {
			if actualEnd < midFrame {
				startFrame = actualEnd
				continue
			}
			startFrame = midFrame
			midFrame = actualEnd
			continue
		}
		i = exactI
		if i >= len(inflectionFrames)-1 {
			break
		}
		i += 2
		if i >= len(inflectionFrames)-1 {
			break
		}
		startFrame = actualEnd
		midFrame = inflectionFrames[i-1]
		endFrame = inflectionFrames[i]
	}

	{
		startFrame := actualEnd
		endFrame := inflectionFrames[len(inflectionFrames)-1]
		midFrame := Frame(int(startFrame+endFrame) / 2)
		actualEnd = b.reduceRange(startFrame, midFrame, endFrame, xs, ys, zs, quats, reduced)
		for actualEnd < endFrame {
			startFrame = actualEnd
			midFrame = Frame(int(actualEnd+endFrame) / 2)
			actualEnd = b.reduceRange(startFrame, midFrame, endFrame, xs, ys, zs, quats, reduced)
		}
	}

	return reduced
}

// ContainsActive は有効なキーフレが存在するか判定する。
func (b *BoneNameFrames) ContainsActive() bool {
	if b == nil || b.Len() == 0 {
		return false
	}

	active := false
	var prev *BoneFrame
	b.ForEach(func(_ Frame, bf *BoneFrame) bool {
		if bf == nil {
			return true
		}
		if !vec3OrZero(bf.Position).NearEquals(mmath.Vec3{}, 1e-2) || !quatOrIdent(bf.Rotation).NearEquals(mmath.NewQuaternion(), 1e-2) {
			active = true
			return false
		}
		if prev != nil {
			if !vec3OrZero(prev.Position).NearEquals(vec3OrZero(bf.Position), 1e-2) {
				active = true
				return false
			}
			if !quatOrIdent(prev.Rotation).NearEquals(quatOrIdent(bf.Rotation), 1e-2) {
				active = true
				return false
			}
		}
		prev = bf
		return true
	})

	return active
}

// Copy はフレーム集合を複製する。
func (b *BoneNameFrames) Copy() (*BoneNameFrames, error) {
	return deepCopy(b)
}

// BoneFrames はボーン名ごとの集合を表す。
type BoneFrames struct {
	names     []string
	nameIndex map[string]int
	values    []*BoneNameFrames
}

// NewBoneFrames はBoneFramesを生成する。
func NewBoneFrames() *BoneFrames {
	return &BoneFrames{
		names:     make([]string, 0),
		nameIndex: make(map[string]int),
		values:    make([]*BoneNameFrames, 0),
	}
}

// Names は登録順の名前一覧を返す。
func (b *BoneFrames) Names() []string {
	if b == nil {
		return nil
	}
	return append([]string(nil), b.names...)
}

// Get は名前に対応するフレーム集合を返す。
func (b *BoneFrames) Get(name string) *BoneNameFrames {
	if b == nil {
		return nil
	}
	if idx, ok := b.nameIndex[name]; ok {
		return b.values[idx]
	}
	frames := NewBoneNameFrames(name)
	b.nameIndex[name] = len(b.values)
	b.names = append(b.names, name)
	b.values = append(b.values, frames)
	return frames
}

// Update はフレーム集合を更新する。
func (b *BoneFrames) Update(frames *BoneNameFrames) {
	if b == nil || frames == nil {
		return
	}
	if idx, ok := b.nameIndex[frames.Name]; ok {
		b.values[idx] = frames
		return
	}
	b.nameIndex[frames.Name] = len(b.values)
	b.names = append(b.names, frames.Name)
	b.values = append(b.values, frames)
}

// Delete は名前を削除する。
func (b *BoneFrames) Delete(name string) {
	if b == nil {
		return
	}
	idx, ok := b.nameIndex[name]
	if !ok {
		return
	}
	b.names = append(b.names[:idx], b.names[idx+1:]...)
	b.values = append(b.values[:idx], b.values[idx+1:]...)
	delete(b.nameIndex, name)
	for i := idx; i < len(b.names); i++ {
		b.nameIndex[b.names[i]] = i
	}
}

// Indexes は全トラックのフレーム番号を返す。
func (b *BoneFrames) Indexes() []int {
	if b == nil {
		return nil
	}
	indexes := make([]int, 0)
	for _, frames := range b.values {
		frames.ForEach(func(frame Frame, _ *BoneFrame) bool {
			indexes = append(indexes, int(frame))
			return true
		})
	}
	indexes = mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

// Len は全フレーム数を返す。
func (b *BoneFrames) Len() int {
	if b == nil {
		return 0
	}
	count := 0
	for _, frames := range b.values {
		count += frames.Len()
	}
	return count
}

// MaxFrame は最大フレーム番号を返す。
func (b *BoneFrames) MaxFrame() Frame {
	if b == nil {
		return 0
	}
	maxFrame := Frame(0)
	for _, frames := range b.values {
		frame := frames.MaxFrame()
		if frame > maxFrame {
			maxFrame = frame
		}
	}
	return maxFrame
}

// MinFrame は最小フレーム番号を返す。
func (b *BoneFrames) MinFrame() Frame {
	if b == nil {
		return 0
	}
	minFrame := Frame(math.MaxFloat32)
	for _, frames := range b.values {
		frame := frames.MinFrame()
		if frame < minFrame {
			minFrame = frame
		}
	}
	if minFrame == Frame(math.MaxFloat32) {
		return 0
	}
	return minFrame
}

// Clean は無効なトラックを削除する。
func (b *BoneFrames) Clean() {
	if b == nil {
		return
	}
	keptNames := make([]string, 0, len(b.names))
	keptValues := make([]*BoneNameFrames, 0, len(b.values))
	keptIndex := make(map[string]int)
	for _, name := range b.names {
		frames := b.Get(name)
		if frames == nil || !frames.ContainsActive() {
			continue
		}
		keptIndex[name] = len(keptValues)
		keptNames = append(keptNames, name)
		keptValues = append(keptValues, frames)
	}
	b.names = keptNames
	b.values = keptValues
	b.nameIndex = keptIndex
}

// Copy はボーン集合を複製する。
func (b *BoneFrames) Copy() (*BoneFrames, error) {
	copied, err := deepCopy(b)
	if err != nil {
		return nil, err
	}
	copied.rebuildNameIndex()
	return copied, nil
}

func (b *BoneFrames) rebuildNameIndex() {
	b.nameIndex = make(map[string]int, len(b.names))
	for i, name := range b.names {
		b.nameIndex[name] = i
	}
}

func nilBoneFrame() *BoneFrame {
	return nil
}

// boneCurveT はボーン補間係数を返す。
func boneCurveT(prev, now, next Frame, curves *BoneCurves) (float64, float64, float64, float64) {
	if curves == nil {
		t := linearT(prev, now, next)
		return t, t, t, t
	}
	return curves.Evaluate(prev, now, next)
}

// linearT は線形補間の係数を返す。
func linearT(prev, now, next Frame) float64 {
	denom := float64(next - prev)
	if denom == 0 {
		return 0
	}
	return float64(now-prev) / denom
}

// reduceRange は区間ごとの曲線当てはめを行う。
func (b *BoneNameFrames) reduceRange(startFrame, midFrame, endFrame Frame, xs, ys, zs []float64, quats []mmath.Quaternion, reduced *BoneNameFrames) Frame {
	startI := int(startFrame)
	endI := int(endFrame)

	rangeXs, rangeYs, rangeZs := sliceRange(xs, ys, zs, startI, endI)
	rangeRs := make([]float64, 0, len(rangeXs))
	startQuat := quats[startI]
	endQuat := quats[endI]
	for i := startI; i <= endI; i++ {
		rangeRs = append(rangeRs, mmath.FindSlerpT(startQuat, endQuat, quats[i], 0))
	}

	xCurve, xErr := mmath.NewCurveFromValues(rangeXs, 1e-2)
	yCurve, yErr := mmath.NewCurveFromValues(rangeYs, 1e-2)
	zCurve, zErr := mmath.NewCurveFromValues(rangeZs, 1e-2)
	rCurve, rErr := mmath.NewCurveFromValues(rangeRs, 1e-4)

	if xErr == nil && yErr == nil && zErr == nil && rErr == nil && xCurve != nil && yCurve != nil && zCurve != nil && rCurve != nil {
		success := true
		for i := startI + 1; i < endI; i++ {
			if !checkCurve(
				xCurve, yCurve, zCurve, rCurve,
				xs[startI], xs[i], xs[endI],
				ys[startI], ys[i], ys[endI],
				zs[startI], zs[i], zs[endI],
				quats[startI], quats[i], quats[endI],
				Frame(startI), Frame(i), Frame(endI),
			) {
				success = false
				break
			}
		}
		if success {
			bf := b.Get(endFrame)
			reduceBf := NewBoneFrame(endFrame)
			reduceBf.Position = copyVec3(bf.Position)
			reduceBf.Rotation = copyQuaternion(bf.Rotation)
			reduceBf.Curves = &BoneCurves{
				TranslateX: xCurve,
				TranslateY: yCurve,
				TranslateZ: zCurve,
				Rotate:     rCurve,
			}
			reduced.Append(reduceBf)
			return endFrame
		}
	}

	midI := int(midFrame)
	if midI == startI || midI == endI {
		bf := b.Get(startFrame)
		reduceBf := NewBoneFrame(startFrame)
		reduceBf.Position = copyVec3(bf.Position)
		reduceBf.Rotation = copyQuaternion(bf.Rotation)
		if bf.Curves != nil {
			reduceBf.Curves = bf.Curves.Copy()
		}
		reduced.Append(reduceBf)
		return midFrame
	}

	return b.reduceRange(startFrame, Frame(int(midFrame+startFrame)/2), midFrame, xs, ys, zs, quats, reduced)
}

// checkCurve は曲線の近似一致を判定する。
func checkCurve(
	xCurve, yCurve, zCurve, rCurve *mmath.Curve,
	startX, nowX, endX float64,
	startY, nowY, endY float64,
	startZ, nowZ, endZ float64,
	startQuat, nowQuat, endQuat mmath.Quaternion,
	startFrame, nowFrame, endFrame Frame,
) bool {
	_, xy, _ := mmath.Evaluate(xCurve, float32(startFrame), float32(nowFrame), float32(endFrame))
	_, yy, _ := mmath.Evaluate(yCurve, float32(startFrame), float32(nowFrame), float32(endFrame))
	_, zy, _ := mmath.Evaluate(zCurve, float32(startFrame), float32(nowFrame), float32(endFrame))
	_, ry, _ := mmath.Evaluate(rCurve, float32(startFrame), float32(nowFrame), float32(endFrame))

	checkQuat := startQuat.Slerp(endQuat, ry)
	if !checkQuat.NearEquals(nowQuat, 1e-1) {
		return false
	}
	if !mmath.NearEquals(mmath.Lerp(startX, endX, xy), nowX, 1e-1) {
		return false
	}
	if !mmath.NearEquals(mmath.Lerp(startY, endY, yy), nowY, 1e-1) {
		return false
	}
	return mmath.NearEquals(mmath.Lerp(startZ, endZ, zy), nowZ, 1e-1)
}

// sliceRange は開始/終了に合わせて値を切り出す。
func sliceRange(xs, ys, zs []float64, startI, endI int) ([]float64, []float64, []float64) {
	if endI >= len(xs) {
		return xs[startI:], ys[startI:], zs[startI:]
	}
	return xs[startI : endI+1], ys[startI : endI+1], zs[startI : endI+1]
}

// findInflectionFrames は変曲点フレームを抽出する。
func findInflectionFrames(frames []Frame, values []float64, threshold float64) []Frame {
	if len(frames) <= 2 || len(values) <= 2 {
		return frames
	}
	out := []Frame{frames[0]}
	for i := 2; i < len(values); i++ {
		delta := values[i] - values[i-1]
		if (delta > threshold && values[i-1] < values[i-2]) || (delta < -threshold && values[i-1] > values[i-2]) {
			out = append(out, frames[i-1])
		}
	}
	first := gradient(values, 1)
	second := gradient(first, 1)
	for i := 1; i < len(second); i++ {
		d1 := mmath.Round(second[i-1], threshold)
		d2 := mmath.Round(second[i], threshold)
		if d1*d2 < 0 || (d1 == 0 && d2 < 0) || (d1 < 0 && d2 == 0) {
			out = append(out, frames[i])
		}
	}
	last := frames[len(frames)-1]
	if !mmath.Contains(out, last) {
		out = append(out, last)
	}
	out = mmath.Unique(out)
	mmath.Sort(out)
	return out
}

// gradient は一次微分を近似する。
func gradient(values []float64, dx float64) []float64 {
	n := len(values)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	if n == 1 {
		out[0] = 0
		return out
	}
	out[0] = (values[1] - values[0]) / dx
	for i := 1; i < n-1; i++ {
		out[i] = (values[i+1] - values[i-1]) / (2 * dx)
	}
	out[n-1] = (values[n-1] - values[n-2]) / dx
	return out
}
