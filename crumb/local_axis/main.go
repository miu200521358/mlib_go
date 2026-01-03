package main

import (
	"fmt"
	"runtime"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func init() {
	runtime.LockOSThread()
}

// exe を実行すると、ローカル軸を求めるための必要情報を標準入力から受け取り、標準出力にローカル軸情報を出力します。
func main() {
	// ローカル軸を設定したいボーン位置をカンマ区切りでもらうよう、printfで促す
	println("ローカル軸を設定したいボーン位置をカンマ区切りで入力してください。例: 0.0,1.0,0.0")

	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println("入力エラー:", err)
		return
	}

	bonePos, err := parsePosition(input)
	if err != nil {
		fmt.Println("ボーン位置取得エラー:", err)
		return
	}

	// ボーンの先位置をグローバル座標でカンマ区切りで受け取る
	println("ボーンの先位置をグローバル座標でカンマ区切りで入力してください。例: 0.0,2.0,0.0")
	_, err = fmt.Scanln(&input)
	if err != nil {
		fmt.Println("入力エラー:", err)
		return
	}

	boneTipPos, err := parsePosition(input)
	if err != nil {
		fmt.Println("ボーン先位置取得エラー:", err)
		return
	}

	// 水平ベクトルの始点位置をグローバル座標でカンマ区切りで受け取る
	println("水平ベクトルの始点位置をグローバル座標でカンマ区切りで入力してください。例: 1.0,1.0,0.0")
	_, err = fmt.Scanln(&input)
	if err != nil {
		fmt.Println("入力エラー:", err)
		return
	}

	horizontalVecStartPos, err := parsePosition(input)
	if err != nil {
		fmt.Println("水平ベクトル始点位置取得エラー:", err)
		return
	}

	// 水平ベクトルの終点位置をグローバル座標でカンマ区切りで受け取る
	println("水平ベクトルの終点位置をグローバル座標でカンマ区切りで入力してください。例: 1.0,1.0,1.0")
	_, err = fmt.Scanln(&input)
	if err != nil {
		fmt.Println("入力エラー:", err)
		return
	}

	horizontalVecEndPos, err := parsePosition(input)
	if err != nil {
		fmt.Println("水平ベクトル終点位置取得エラー:", err)
		return
	}

	// // 水平ベクトルを X, Y, Z のどのローカル軸で曲げたいかを、X, Y, Z のいずれかで受け取る
	// println("水平ベクトルを X, Y, Z のどのローカル軸で曲げたいかを、X, Y, Z のいずれかで入力してください。例: Y")
	// _, err = fmt.Scanln(&input)
	// if err != nil {
	// 	fmt.Println("入力エラー:", err)
	// 	return
	// }
	// horizontalAxisStr := input
	horizontalAxisStr := "Z"

	// // 垂直ベクトルを X, Y, Z のどのローカル軸で曲げたいかを、X, Y, Z のいずれかで受け取る
	// println("垂直ベクトルを X, Y, Z のどのローカル軸で曲げたいかを、X, Y, Z のいずれかで入力してください。例: Z")
	// _, err = fmt.Scanln(&input)
	// if err != nil {
	// 	fmt.Println("入力エラー:", err)
	// 	return
	// }
	// verticalAxisStr := input
	verticalAxisStr := "Y"

	// ローカル軸を求める
	localXAxis, localZAxis, err := calculateLocalAxis(bonePos, boneTipPos, horizontalVecStartPos, horizontalVecEndPos, horizontalAxisStr, verticalAxisStr)
	if err != nil {
		fmt.Println("ローカル軸計算エラー:", err)
		return
	}

	// 結果を出力
	fmt.Printf("ローカルX軸: %.6f, %.6f, %.6f\n", localXAxis.X, localXAxis.Y, localXAxis.Z)
	fmt.Printf("ローカルZ軸: %.6f, %.6f, %.6f\n", localZAxis.X, localZAxis.Y, localZAxis.Z)

	// Enterが押されるまで待機
	fmt.Println("Enterキーを押して終了してください...")
	fmt.Scanln()
}

// ローカル軸を計算する関数
func calculateLocalAxis(bonePos, boneTipPos, horizontalVecStartPos, horizontalVecEndPos mmath.MVec3,
	horizontalAxisStr, verticalAxisStr string) (*mmath.MVec3, *mmath.MVec3, error) {
	// ボーンのベクトルを計算
	boneVec := boneTipPos.Sub(&bonePos).Normalized()

	// 水平ベクトルを計算
	horizontalVec := horizontalVecEndPos.Sub(&horizontalVecStartPos).Normalized()

	// ローカル軸水平ベクトルを決定
	var localAxisX *mmath.MVec3
	switch horizontalAxisStr {
	case "X":
		localAxisX = horizontalVec
	case "Y":
		localAxisX = boneVec
	case "Z":
		localAxisX = horizontalVec.Cross(boneVec).Normalized()
	default:
		return nil, nil, fmt.Errorf("無効なローカル軸指定: %s", verticalAxisStr)
	}

	// ローカル軸垂直ベクトルを決定
	var localAxisZ *mmath.MVec3
	switch verticalAxisStr {
	case "X":
		localAxisZ = horizontalVec
	case "Y":
		localAxisZ = horizontalVec
	case "Z":
		localAxisZ = horizontalVec
	default:
		return nil, nil, fmt.Errorf("無効なローカル軸指定: %s", horizontalAxisStr)
	}

	return localAxisX, localAxisZ, nil
}

func parsePosition(input string) (mmath.MVec3, error) {
	// 入力をパースして、mmath.MVec3に変換する
	var position mmath.MVec3
	_, err := fmt.Sscanf(input, "%f,%f,%f", &position.X, &position.Y, &position.Z)
	if err != nil {
		return mmath.MVec3{}, fmt.Errorf("入力のパースエラー: %v", err)
	}
	// MVec3 に入れる
	return position, nil
}
