package mmath

import (
	"math"

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

func CalculateBoundingBox(positions []*MVec3) (size *MVec3, position *MVec3, radians *MVec3) {
	min, max := calculateMinMax(positions)
	size = max.Subed(min).MulScalar(0.5)
	position = min.Add(size)
	rotation := calculateRotation(positions)

	return &MVec3{X: size.X, Y: size.Y / 2, Z: size.Z}, position, rotation.ToRadians()
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
