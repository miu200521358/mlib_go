package mmath

import (
	"math"
	"math/rand"
	"sort"
	"time"

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

// ---------------------------------------

// Sphere represents a Sphere with a center and a radius.
type Sphere struct {
	Center *MVec3
	Radius float64
}

// average computes the centroid of a slice of Vector3.
func average(points []*MVec3) *MVec3 {
	if len(points) == 0 {
		return NewMVec3()
	}
	sum := NewMVec3()
	for _, p := range points {
		sum.Add(p)
	}
	return sum.MulScalar(1.0 / float64(len(points)))
}

// computeBoundingSphere computes a simple bounding sphere
// by taking the centroid as the center and the max distance
// to any point as the radius.
func computeBoundingSphere(points []*MVec3) *Sphere {
	if len(points) == 0 {
		return &Sphere{Center: NewMVec3(), Radius: 0}
	}

	c := average(points)
	maxR := 0.0
	for _, p := range points {
		d := p.Distance(c)
		if d > maxR {
			maxR = d
		}
	}
	return &Sphere{Center: c, Radius: maxR}
}

// kMeans2Partitions splits points into 2 clusters using a simple k-means(k=2).
// Returns two sets of points as sub clusters.
func kMeans2Partitions(points []*MVec3, maxIter int) ([]*MVec3, []*MVec3) {
	if len(points) < 2 {
		// Not enough points to split meaningfully
		return points, nil
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize two centers randomly from the dataset
	idx1 := rand.Intn(len(points))
	idx2 := rand.Intn(len(points))
	for idx2 == idx1 {
		idx2 = rand.Intn(len(points))
	}
	center1 := points[idx1]
	center2 := points[idx2]

	for iter := 0; iter < maxIter; iter++ {
		cluster1 := []*MVec3{}
		cluster2 := []*MVec3{}

		// Assign each point to the nearest center
		for _, p := range points {
			d1 := p.Distance(center1)
			d2 := p.Distance(center2)
			if d1 <= d2 {
				cluster1 = append(cluster1, p)
			} else {
				cluster2 = append(cluster2, p)
			}
		}

		// Recompute centers
		newCenter1 := average(cluster1)
		newCenter2 := average(cluster2)

		// Check convergence
		move1 := center1.Distance(newCenter1)
		move2 := center2.Distance(newCenter2)
		center1 = newCenter1
		center2 = newCenter2

		if move1 < 1e-6 && move2 < 1e-6 {
			break
		}
	}

	// Final clustering
	cluster1 := []*MVec3{}
	cluster2 := []*MVec3{}
	for _, p := range points {
		d1 := p.Distance(center1)
		d2 := p.Distance(center2)
		if d1 <= d2 {
			cluster1 = append(cluster1, p)
		} else {
			cluster2 = append(cluster2, p)
		}
	}

	return cluster1, cluster2
}

// AdaptiveCoverPointsWithSpheres recursively splits the point set
// until each bounding sphere is smaller than the given threshold.
func AdaptiveCoverPointsWithSpheres(points []*MVec3, radiusThreshold float64, maxIter int) []*Sphere {
	if len(points) == 0 {
		return nil
	}

	// Compute bounding sphere
	bounding := computeBoundingSphere(points)

	// If the radius is smaller than threshold or we have too few points, stop splitting
	if bounding.Radius <= radiusThreshold || len(points) < 2 {
		return []*Sphere{bounding}
	}

	// Otherwise, split the points into two clusters
	c1, c2 := kMeans2Partitions(points, maxIter)

	// Recursively cover each sub cluster
	spheres := []*Sphere{}
	if len(c1) > 0 {
		for _, s := range AdaptiveCoverPointsWithSpheres(c1, radiusThreshold, maxIter) {
			if s.Radius > 0.0 {
				spheres = append(spheres, s)
			}
		}
	}
	if len(c2) > 0 {
		for _, s := range AdaptiveCoverPointsWithSpheres(c2, radiusThreshold, maxIter) {
			if s.Radius > 0.0 {
				spheres = append(spheres, s)
			}
		}
	}

	return spheres
}
