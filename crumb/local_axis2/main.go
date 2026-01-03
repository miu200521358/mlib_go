package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/merr"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
)

func main() {
	println("読み込み対象モデルフルパスを指定してください")

	reader := bufio.NewReader(os.Stdin)
	modelPath, err := reader.ReadString('\n')
	modelPath = strings.TrimSpace(modelPath)
	modelPath = strings.Trim(modelPath, "\"")
	fmt.Println("モデル読み込み中:", modelPath)

	// modelPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/172_桑名江/桑名江 ウメ式 ver1.00/ウメ式桑名江(武装解除)ver1.00_ローカル軸.pmx"

	data, err := repository.NewPmxRepository(false).Load(modelPath)
	if err != nil {
		fmt.Println("モデル読み込みエラー:", err)
		return
	}
	model := data.(*pmx.PmxModel)
	model.Setup()
	vertexMap := model.Vertices.GetMapByBoneIndex(0.6)
	insertBones(model.Bones, model.Vertices)

	fingerPairs := [][]pmx.StandardBoneName{
		{pmx.THUMB0, pmx.THUMB1},
		{pmx.THUMB1, pmx.THUMB2},
		{pmx.THUMB2, pmx.THUMB_TAIL},
		{pmx.INDEX1, pmx.INDEX2},
		{pmx.INDEX2, pmx.INDEX3},
		{pmx.INDEX3, pmx.INDEX_TAIL},
		{pmx.MIDDLE1, pmx.MIDDLE2},
		{pmx.MIDDLE2, pmx.MIDDLE3},
		{pmx.MIDDLE3, pmx.MIDDLE_TAIL},
		{pmx.RING1, pmx.RING2},
		{pmx.RING2, pmx.RING3},
		{pmx.RING3, pmx.RING_TAIL},
		{pmx.PINKY1, pmx.PINKY2},
		{pmx.PINKY2, pmx.PINKY3},
		{pmx.PINKY3, pmx.PINKY_TAIL},
	}

	for _, direction := range []pmx.BoneDirection{pmx.BONE_DIRECTION_LEFT, pmx.BONE_DIRECTION_RIGHT} {
		for _, pair := range fingerPairs {
			fromBoneName := pair[0].StringFromDirection(direction)
			toBoneName := pair[1].StringFromDirection(direction)

			fromBone, err := model.Bones.GetByName(fromBoneName)
			if err != nil {
				fmt.Printf("ボーン取得エラー: %s (%s): %v\n", fromBoneName, direction.String(), err)
				continue
			}

			toBone, err := model.Bones.GetByName(toBoneName)
			if err != nil {
				fmt.Printf("ボーン取得エラー: %s (%s): %v\n", toBoneName, direction.String(), err)
				continue
			}

			// ボーンの位置から頂点を取得
			vertices, found := vertexMap[fromBone.Index()]
			if !found {
				fmt.Printf("頂点が見つかりません: %s (%s)\n", fromBoneName, direction.String())
				continue
			}
			vertexIndexes := make([]string, 0, len(vertices))
			for _, v := range vertices {
				vertexIndexes = append(vertexIndexes, fmt.Sprintf("%d (%.3f)", v.Index(), v.Deform.IndexWeight(fromBone.Index())))
			}

			// ボーン軸方向を計算
			boneAxis := toBone.Position.Subed(fromBone.Position).Normalized()

			// PCAベースで★の位置を求める
			boneCenter := fromBone.Position.Added(toBone.Position).MuledScalar(0.5)
			leftStar, rightStar, horizontalDir := findStarVerticesByPCA(vertices, boneAxis, boneCenter)

			// デバッグ出力
			fmt.Printf("=== %s → %s (%s) ===\n", fromBoneName, toBoneName, direction.String())
			fmt.Printf("対象頂点Index: %v\n", vertexIndexes)
			fmt.Printf("頂点数: %d\n", len(vertices))
			if leftStar != nil && rightStar != nil {
				fmt.Printf("★左: %d (%.3f, %.3f, %.3f)\n", leftStar.Index(), leftStar.Position.X, leftStar.Position.Y, leftStar.Position.Z)
				fmt.Printf("★右: %d (%.3f, %.3f, %.3f)\n", rightStar.Index(), rightStar.Position.X, rightStar.Position.Y, rightStar.Position.Z)
				fmt.Printf("水平方向: (%.3f, %.3f, %.3f)\n", horizontalDir.X, horizontalDir.Y, horizontalDir.Z)

				// ローカル軸設定
				fromBone.BoneFlag |= pmx.BONE_FLAG_HAS_LOCAL_AXIS

				fromBone.LocalAxisX = boneAxis.Normalized()
				localAxisY := fromBone.LocalAxis.Cross(horizontalDir).Normalized()
				fromBone.LocalAxisZ = localAxisY.Cross(fromBone.LocalAxisX).Normalized()
				fmt.Printf("ローカル軸X: (%.3f, %.3f, %.3f)\n", fromBone.LocalAxisX.X, fromBone.LocalAxisX.Y, fromBone.LocalAxisX.Z)
				fmt.Printf("ローカル軸Y: (%.3f, %.3f, %.3f)\n", localAxisY.X, localAxisY.Y, localAxisY.Z)
				fmt.Printf("ローカル軸Z: (%.3f, %.3f, %.3f)\n", fromBone.LocalAxisZ.X, fromBone.LocalAxisZ.Y, fromBone.LocalAxisZ.Z)
			} else {
				fmt.Printf("★が見つかりませんでした\n")
			}
		}
	}

	fmt.Println("処理完了")
	outputPath := mfile.CreateOutputPath(modelPath, "ローカル軸")
	if err := repository.NewPmxRepository(false).Save(outputPath, model, false); err != nil {
		fmt.Println("モデル保存エラー:", err)
	}
	fmt.Println("保存完了:", outputPath)
}

func insertBones(bones *pmx.Bones, vertices *pmx.Vertices) error {

	// 左右系
	for _, direction := range []pmx.BoneDirection{pmx.BONE_DIRECTION_LEFT, pmx.BONE_DIRECTION_RIGHT} {
		for _, funcs := range [][]func(direction pmx.BoneDirection) (*pmx.Bone, error){
			// {bones.GetShoulderRoot, bones.CreateShoulderRoot},
			{bones.GetWristTail, bones.CreateWristTail},
			{bones.GetThumbTail, bones.CreateThumbTail},
			{bones.GetIndexTail, bones.CreateIndexTail},
			{bones.GetMiddleTail, bones.CreateMiddleTail},
			{bones.GetRingTail, bones.CreateRingTail},
			{bones.GetPinkyTail, bones.CreatePinkyTail},
			// {bones.GetHip, bones.CreateHip},
			// {bones.GetLegRoot, bones.CreateLegRoot},
			// {bones.GetLegD, bones.CreateLegD},
			// {bones.GetKneeD, bones.CreateKneeD},
			// {bones.GetAnkleD, bones.CreateAnkleD},
			// {bones.GetToeT, bones.CreateToeT},
			// {bones.GetToeP, bones.CreateToeP},
			// {bones.GetToeC, bones.CreateToeC},
			// {bones.GetHeel, bones.CreateHeel},
			// {bones.GetAnkleDGround, bones.CreateAnkleDGround},
			// {bones.GetToeEx, bones.CreateToeEx},
			// {bones.GetToeTD, bones.CreateToeTD},
			// {bones.GetToePD, bones.CreateToePD},
			// {bones.GetToeCD, bones.CreateToeCD},
			// {bones.GetHeelD, bones.CreateHeelD},
		} {
			getFunc := funcs[0]
			createFunc := funcs[1]

			if bone, err := getFunc(direction); err != nil && merr.IsNameNotFoundError(err) && bone == nil {
				if bone, err := createFunc(direction); err == nil && bone != nil {
					bone.IsSystem = true

					if bone.Name() == pmx.TOE_T.StringFromDirection(direction) {
						// つま先の位置は、足首・足首D・足先EXの中で最もZ値が小さい位置にする
						vertexMap := vertices.GetMapByBoneIndex(1e-1)
						if vertexMap != nil {
							for _, ankleBoneName := range []string{
								pmx.ANKLE.StringFromDirection(direction),
								pmx.ANKLE_D.StringFromDirection(direction),
								pmx.TOE_EX.StringFromDirection(direction),
							} {
								ankleBone, _ := bones.GetByName(ankleBoneName)
								if ankleBone == nil {
									continue
								}
								if boneVertices, ok := vertexMap[ankleBone.Index()]; ok && boneVertices != nil {
									for _, vertex := range boneVertices {
										if vertex.Position.Z < bone.Position.Z && vertex.Position.Y < bone.Position.Y {
											bone.Position = vertex.Position.Copy()
										}
									}
								}
							}
							bone.Position.Y = 0
						}
					}

					if err := bones.Insert(bone); err != nil {
						return err
					} else {
						bones.SetParentFromConfig(bone)

						// 再セットアップ
						bones.Setup()
					}
				} else if merr.IsParentNotFoundError(err) {
					// 何もしない
				} else {
					return err
				}
			} else if err != nil {
				return err
			}
		}

		{
			// 親指0
			if bone, err := bones.GetThumb(direction, 0); err != nil && merr.IsNameNotFoundError(err) && bone == nil {
				if thumb0, err := bones.CreateThumb0(direction); err == nil && thumb0 != nil {
					if err := bones.Insert(thumb0); err != nil {
						return err
					}
					if thumb1, err := bones.GetThumb(direction, 1); err == nil && thumb1 != nil {
						thumb1.ParentIndex = thumb0.Index()
					}
					bones.Setup()
				} else if merr.IsParentNotFoundError(err) {
					// 何もしない
				} else {
					return err
				}
			} else if err != nil {
				return err
			}
		}
	}
	return nil
}

// findStarVerticesByPCA はPCAを使って★の位置（指の左右端点）を求めます
func findStarVerticesByPCA(vertices []*pmx.Vertex, boneAxis *mmath.MVec3, boneCenter *mmath.MVec3) (*pmx.Vertex, *pmx.Vertex, *mmath.MVec3) {
	if len(vertices) < 2 {
		return nil, nil, nil
	}

	// Step 1: 各頂点をボーン軸に垂直な平面に投影
	// 投影: p' = p - (p・axis)*axis
	var projectedPoints []struct {
		vertex    *pmx.Vertex
		projected *mmath.MVec3 // ボーン中心からの相対位置（ボーン軸成分を除去）
	}

	for _, v := range vertices {
		// ボーン中心からの相対位置
		relPos := v.Position.Subed(boneCenter)
		// ボーン軸成分を除去（垂直平面に投影）
		axisComponent := boneAxis.MuledScalar(relPos.Dot(boneAxis))
		projected := relPos.Subed(axisComponent)

		projectedPoints = append(projectedPoints, struct {
			vertex    *pmx.Vertex
			projected *mmath.MVec3
		}{vertex: v, projected: projected})
	}

	// Step 2: 投影された点群の重心を計算
	var centroidX, centroidY, centroidZ float64
	for _, pp := range projectedPoints {
		centroidX += pp.projected.X
		centroidY += pp.projected.Y
		centroidZ += pp.projected.Z
	}
	n := float64(len(projectedPoints))
	centroid := &mmath.MVec3{X: centroidX / n, Y: centroidY / n, Z: centroidZ / n}

	// Step 3: 共分散行列を計算（3x3だがボーン軸方向は0になるので実質2D）
	var cxx, cxy, cxz, cyy, cyz, czz float64
	for _, pp := range projectedPoints {
		dx := pp.projected.X - centroid.X
		dy := pp.projected.Y - centroid.Y
		dz := pp.projected.Z - centroid.Z
		cxx += dx * dx
		cxy += dx * dy
		cxz += dx * dz
		cyy += dy * dy
		cyz += dy * dz
		czz += dz * dz
	}
	cxx /= n
	cxy /= n
	cxz /= n
	cyy /= n
	cyz /= n
	czz /= n

	// Step 4: 3x3共分散行列の最大固有値の固有ベクトルを求める（Power Iteration法）
	// 初期ベクトル
	eigenvector := &mmath.MVec3{X: 1, Y: 0, Z: 0}

	for iter := 0; iter < 100; iter++ {
		// 行列×ベクトル
		newX := cxx*eigenvector.X + cxy*eigenvector.Y + cxz*eigenvector.Z
		newY := cxy*eigenvector.X + cyy*eigenvector.Y + cyz*eigenvector.Z
		newZ := cxz*eigenvector.X + cyz*eigenvector.Y + czz*eigenvector.Z

		newVec := &mmath.MVec3{X: newX, Y: newY, Z: newZ}
		length := newVec.Length()
		if length < 1e-10 {
			break
		}
		eigenvector = newVec.DivedScalar(length)
	}

	// 水平方向 = 最大分散方向
	horizontalDir := eigenvector.Normalized()

	// Step 5: この方向で最も離れた2点を見つける
	var minDot, maxDot float64 = math.MaxFloat64, -math.MaxFloat64
	var leftStar, rightStar *pmx.Vertex

	for _, pp := range projectedPoints {
		// 重心からの相対位置を水平方向に投影
		relFromCentroid := pp.projected.Subed(centroid)
		dot := relFromCentroid.Dot(horizontalDir)

		if dot < minDot {
			minDot = dot
			leftStar = pp.vertex
		}
		if dot > maxDot {
			maxDot = dot
			rightStar = pp.vertex
		}
	}

	return leftStar, rightStar, horizontalDir
}
