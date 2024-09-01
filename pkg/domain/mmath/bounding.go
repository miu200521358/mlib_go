package mmath

import (
	"math"
	"sort"

	"github.com/gonum/matrix/mat64"
)

func calculateMinMax(positions []*MVec3) (*MVec3, *MVec3) {
	min := &MVec3{X: math.MaxFloat32, Y: math.MaxFloat32, Z: math.MaxFloat32}
	max := &MVec3{X: -math.MaxFloat32, Y: -math.MaxFloat32, Z: -math.MaxFloat32}

	for _, position := range positions {
		if position.X < min.X {
			min.X = position.X
		}
		if position.Y < min.Y {
			min.Y = position.Y
		}
		if position.Z < min.Z {
			min.Z = position.Z
		}

		if position.X > max.X {
			max.X = position.X
		}
		if position.Y > max.Y {
			max.Y = position.Y
		}
		if position.Z > max.Z {
			max.Z = position.Z
		}
	}

	return min, max
}

func CalculateBoundingBox(positions []*MVec3, threshold float64) (size *MVec3, position *MVec3, radians *MVec3) {
	filteredPositions := MedianBasedOutlierFilter(positions, threshold)
	min, max := calculateMinMax(filteredPositions)
	size = max.Subed(min).MulScalar(0.5)
	position = min.Add(size)
	rotation := calculateRotation(filteredPositions)

	return size, position, rotation.ToRadians()
}

func calculateRotation(vertices []*MVec3) *MQuaternion {
	covariance := computeCovarianceMatrix(vertices)
	eigenvectors := computeEigenVectors(covariance)

	// 固有ベクトルをベクトル3として取り出し
	rotationX := &MVec3{eigenvectors.At(0, 0), eigenvectors.At(1, 0), eigenvectors.At(2, 0)}
	rotationY := &MVec3{eigenvectors.At(0, 1), eigenvectors.At(1, 1), eigenvectors.At(2, 1)}
	rotationZ := &MVec3{eigenvectors.At(0, 2), eigenvectors.At(1, 2), eigenvectors.At(2, 2)}

	// 回転行列からオイラー角（回転角度）を計算
	rotation := NewMQuaternionFromAxes(rotationX, rotationY, rotationZ)

	return rotation
}

func computeEigenVectors(covariance *mat64.Dense) *mat64.Dense {
	var eig mat64.Eigen
	eig.Factorize(covariance, true, true)
	eigenvectors := eig.Vectors()

	return eigenvectors
}

func computeCovarianceMatrix(positions []*MVec3) *mat64.Dense {
	mean := computeMean(positions)
	n := len(positions)

	covariance := mat64.NewDense(3, 3, nil)
	for _, position := range positions {
		diff := position.Subed(mean)
		diffVec := mat64.NewDense(3, 1, []float64{diff.X, diff.Y, diff.Z})
		diffVecT := mat64.NewDense(1, 3, []float64{diff.X, diff.Y, diff.Z})

		var product mat64.Dense
		product.Mul(diffVec, diffVecT)
		covariance.Add(covariance, &product)
	}
	covariance.Scale(1/float64(n-1), covariance)

	return covariance
}

func computeMean(positions []*MVec3) *MVec3 {
	sum := &MVec3{}
	for _, position := range positions {
		sum.Add(position)
	}
	return sum.MulScalar(1.0 / float64(len(positions)))
}

func CalculateBoundingSphere(positions []*MVec3, threshold float64) (size *MVec3, position *MVec3) {
	filteredPositions := MedianBasedOutlierFilter(positions, threshold)
	position = computeMean(filteredPositions)
	radius := computeRadius(filteredPositions, position)

	return &MVec3{X: radius, Y: 0, Z: 0}, position
}

func computeRadius(positions []*MVec3, center *MVec3) float64 {
	maxDistance := 0.0
	for _, position := range positions {
		distance := position.Subed(center).Length()
		if distance > maxDistance {
			maxDistance = distance
		}
	}
	return maxDistance
}

func CalculateBoundingCapsule(positions []*MVec3, threshold float64) (size *MVec3, position *MVec3, radians *MVec3) {
	filteredPositions := MedianBasedOutlierFilter(positions, threshold)
	covariance := computeCovarianceMatrix(filteredPositions)
	eigenvectors := computeEigenVectors(covariance)

	// 主成分分析によって得られた主要軸を使用
	axis := &MVec3{eigenvectors.At(0, 0), eigenvectors.At(1, 0), eigenvectors.At(2, 0)}
	axis.Normalize()

	// カプセルの両端を初期化
	minPoint := filteredPositions[0]
	maxPoint := filteredPositions[0]

	// カプセルの軸に投影された最大・最小点を探す
	for _, position := range filteredPositions {
		proj := position.Dot(axis)
		if proj < minPoint.Dot(axis) {
			minPoint = position
		}
		if proj > maxPoint.Dot(axis) {
			maxPoint = position
		}
	}

	// 中心位置と高さを計算
	position = minPoint.Added(maxPoint).MuledScalar(0.5)
	height := maxPoint.Subed(minPoint).Length()

	// 最大半径を計算
	maxRadius := 0.0
	for _, position := range filteredPositions {
		dist := distancePointToLine(position, minPoint, maxPoint)
		if dist > maxRadius {
			maxRadius = dist
		}
	}

	// カプセルの回転を計算
	rotation := calculateRotationFromAxis(axis)

	return &MVec3{X: maxRadius, Y: height, Z: 0.0}, position, rotation.ToRadians()
}

// 点から直線への距離を計算
func distancePointToLine(point, lineStart, lineEnd *MVec3) float64 {
	lineDir := lineEnd.Subed(lineStart).Normalize()
	projected := lineStart.Added(lineDir.MuledScalar(point.Subed(lineStart).Dot(lineDir)))
	return point.Subed(projected).Length()
}

// カプセルの軸からオイラー角を計算
func calculateRotationFromAxis(axis *MVec3) *MQuaternion {
	// Y軸がカプセルの軸になるようにする
	rotationAxis := MVec3UnitY.Cross(axis).Normalized()
	angle := math.Acos(float64(MVec3UnitY.Dot(axis)))
	return NewMQuaternionFromAxisAnglesRotate(rotationAxis, angle)
}

// MedianBasedOutlierFilter は、中央値を基準に外れ値をフィルタリングします
func MedianBasedOutlierFilter(positions []*MVec3, threshold float64) []*MVec3 {
	median := calculateMedian(positions)
	weights := calculateWeights(positions, median, threshold)
	filteredPositions := filterOutliers(positions, weights, threshold)
	return filteredPositions
}

func calculateMedian(vectors []*MVec3) *MVec3 {
	n := len(vectors)
	if n == 0 {
		return &MVec3{}
	}

	xValues := make([]float64, n)
	yValues := make([]float64, n)
	zValues := make([]float64, n)

	for i, vec := range vectors {
		xValues[i] = vec.X
		yValues[i] = vec.Y
		zValues[i] = vec.Z
	}

	sort.Float64s(xValues)
	sort.Float64s(yValues)
	sort.Float64s(zValues)

	medianVec := &MVec3{
		X: xValues[n/2],
		Y: yValues[n/2],
		Z: zValues[n/2],
	}

	return medianVec
}

func calculateWeights(vectors []*MVec3, median *MVec3, threshold float64) []float64 {
	weights := make([]float64, len(vectors))
	maxDistance := 0.0

	// 各ベクトルの距離を計算し、最大距離を求める
	for _, vec := range vectors {
		dist := vec.Distance(median)
		if dist > maxDistance {
			maxDistance = dist
		}
	}

	// 各ベクトルの重みを計算
	for i, vec := range vectors {
		dist := vec.Distance(median)
		normalizedDist := dist / maxDistance
		weights[i] = 1 - normalizedDist
		if weights[i] < threshold {
			weights[i] = 0
		}
	}

	return weights
}

func filterOutliers(vectors []*MVec3, weights []float64, threshold float64) []*MVec3 {
	filteredVectors := make([]*MVec3, 0)
	for i, vec := range vectors {
		if weights[i] >= threshold {
			filteredVectors = append(filteredVectors, vec)
		}
	}
	return filteredVectors
}
