package mmath

import "math"

// 線形補間
func LerpFloat(v1, v2 float64, t float64) float64 {
	return v1 + ((v2 - v1) * t)
}

func Sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

func NearEquals(v float64, other float64, epsilon float64) bool {
	return math.Abs(v-other) <= epsilon
}

func ToRadian(degree float64) float64 {
	return degree * math.Pi / 180
}

func ToDegree(radian float64) float64 {
	return radian * 180 / math.Pi
}

// Clamp01 ベクトルの各要素をmin～maxの範囲内にクランプします
func ClampFloat(v float64, min float64, max float64) float64 {
	if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

// Clamp01 ベクトルの各要素をmin～maxの範囲内にクランプします
func ClampFloat32(v float32, min float32, max float32) float32 {
	if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

// ボーンから見た頂点ローカル位置を求める
// vertexPositions: グローバル頂点位置
// startBonePosition: 親ボーン位置
// endBonePosition: 子ボーン位置
func GetVertexLocalPositions(vertexPositions []*MVec3, startBonePosition *MVec3, endBonePosition *MVec3) []*MVec3 {
	vertexSize := len(vertexPositions)
	boneVector := endBonePosition.Sub(startBonePosition)
	boneDirection := boneVector.Normalized()

	localPositions := make([]*MVec3, vertexSize)
	for i := 0; i < vertexSize; i++ {
		vertexPosition := vertexPositions[i]
		subedVertexPosition := vertexPosition.Subed(startBonePosition)
		projection := subedVertexPosition.Project(boneDirection)
		localPosition := endBonePosition.Added(projection)
		localPositions[i] = localPosition
	}

	return localPositions
}
