package repository

import (
	"encoding/binary"
	"fmt"

	"github.com/miu200521358/win"
	"golang.org/x/sys/windows"
)

// パッケージ変数として固定ハフマンツリーを保持
var fixedHuffmanTreeInstance huffmanTree

func init() {
	// マップの初期化
	fixedHuffmanMap = make(map[uint32]uint16, 288)

	for _, code := range fixedCodeSorted {
		key := (uint32(code.length) << 16) | uint32(code.huffmanCode)
		fixedHuffmanMap[key] = code.binCode
	}

	// 固定ハフマンツリーを一度だけ初期化
	fixedHuffmanTreeInstance.numCodes = 288
	fixedHuffmanTreeInstance.numMaxBits = 9
	fixedHuffmanTreeInstance.codes = fixedCodeSorted
}

// decompressedBinaryXFile は、圧縮バイナリ形式の X ファイルを解凍し、
// 解凍後のデータをバイナリパーサーへ渡します。
func (rep *XRepository) decompressedBinaryXFile() ([]byte, error) {
	var err error
	var xHandle win.HWND

	// ファイルを開く
	xHandle, err = win.CreateFile(
		rep.path,
		win.GENERIC_READ,
		win.FILE_SHARE_READ,
		0,
		win.OPEN_EXISTING,
		win.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != windows.DS_S_SUCCESS {
		return nil, fmt.Errorf("CreateFile failed: %v", err)
	}
	defer win.CloseHandle(win.HANDLE(xHandle))

	// ファイルサイズを取得
	var fileSize uint64
	if ret, err := win.GetFileSizeEx(xHandle, &fileSize); !ret {
		return nil, fmt.Errorf("GetFileSizeEx failed: %v", err)
	}
	inputFileSize := uint32(fileSize)

	// 圧縮データを読み込む
	compressedBuffer := make([]byte, inputFileSize)
	var bytesRead uint32
	if _, err := win.ReadFile(
		xHandle,
		&compressedBuffer[0],
		inputFileSize,
		&bytesRead,
		0,
	); err != windows.DS_S_SUCCESS {
		return nil, fmt.Errorf("ReadFile failed: %v", err)
	}

	// MSZIP 圧縮アルゴリズムでデコンプレッサを作成
	var decompressorHandle win.HWND
	if _, err := win.CreateDecompressor(win.COMPRESS_ALGORITHM_MSZIP|win.COMPRESS_RAW, &decompressorHandle); err != windows.DS_S_SUCCESS {
		return nil, fmt.Errorf("CreateDecompressor failed: %v", err)
	}
	defer win.CloseDecompressor(decompressorHandle)

	// MSZIP 圧縮データを解凍
	var decompressedBuffer []byte
	if decompressedBuffer, err = rep.decompressMSZipXFile(compressedBuffer, bytesRead, decompressorHandle); err != nil {
		return nil, fmt.Errorf("decompressMSZipXFile failed: %v", err)
	}

	return decompressedBuffer, nil
}

// decompressMSZipXFile は、MSZIP形式の圧縮ブロックを解凍します
func (rep *XRepository) decompressMSZipXFile(
	compressedBuffer []byte,
	inputFileSize uint32,
	decompressorHandle win.HWND,
) ([]byte, error) {
	// ヘッダーサイズ定数
	const (
		headerSize      = 16
		sizeFieldSize   = 4
		totalHeaderSize = headerSize + sizeFieldSize // 20バイト
	)

	// ファイルサイズ検証
	if inputFileSize < totalHeaderSize {
		return nil, fmt.Errorf("ファイルが小さすぎるか破損しています: 入力サイズ=%d < 必要最小サイズ=%d", inputFileSize, totalHeaderSize)
	}

	// 伸長後のファイルサイズを取得
	finalSize := binary.LittleEndian.Uint32(compressedBuffer[headerSize : headerSize+sizeFieldSize])
	if finalSize < headerSize {
		return nil, fmt.Errorf("無効な最終サイズ: %d (最小値 %d より小さい)", finalSize, headerSize)
	}

	// 出力バッファの準備
	newBuffer := make([]byte, finalSize)

	// ヘッダーをコピー
	copy(newBuffer[:headerSize], compressedBuffer[:headerSize])

	// 初期オフセット設定
	destOffset := headerSize
	uncompressedSum := uint32(headerSize)
	srcOffset := totalHeaderSize

	// バッファ長チェック
	compressedLen := len(compressedBuffer)
	if srcOffset > compressedLen {
		return nil, fmt.Errorf("圧縮データ開始位置がバッファサイズを超えています: 開始位置=%d > バッファ長=%d",
			srcOffset, compressedLen)
	}

	// ブロック単位のループ処理
	blockIndex := 0
	prevBytesDecompressed := uint32(0)
	for srcOffset < compressedLen && uncompressedSum < finalSize {
		blockIndex++

		// ブロック非圧縮サイズの読み取り
		if srcOffset+2 > compressedLen {
			return nil, fmt.Errorf("ブロック#%d: 非圧縮サイズを読み取るのに十分なデータがありません", blockIndex)
		}
		blockUncompressedSize := binary.LittleEndian.Uint16(compressedBuffer[srcOffset : srcOffset+2])
		srcOffset += 2

		// ブロック圧縮サイズの読み取り
		if srcOffset+2 > compressedLen {
			return nil, fmt.Errorf("ブロック#%d: 圧縮サイズを読み取るのに十分なデータがありません", blockIndex)
		}
		blockCompressedSize := binary.LittleEndian.Uint16(compressedBuffer[srcOffset : srcOffset+2])
		srcOffset += 2

		// ブロックデータ範囲の検証
		if srcOffset+int(blockCompressedSize) > compressedLen {
			return nil, fmt.Errorf("ブロック#%d: 圧縮データが不足しています (必要=%d, 残り=%d)",
				blockIndex, blockCompressedSize, compressedLen-srcOffset)
		}

		// ブロックデータの取得
		blockData := compressedBuffer[srcOffset : srcOffset+int(blockCompressedSize)]

		// Windowsの解凍API呼び出し
		success, err := win.Decompress(
			decompressorHandle,
			&blockData[0],
			uint32(blockCompressedSize),
			&newBuffer[destOffset],
			uint32(blockUncompressedSize),
			nil,
		)
		if !success {
			return nil, fmt.Errorf("ブロック#%d の解凍に失敗: 非圧縮サイズ=%d, 圧縮サイズ=%d, エラー=%v",
				blockIndex, blockUncompressedSize, blockCompressedSize, err)
		}

		// RFC1951形式のブロック解凍
		bytesDecompressed, err := rep.decompressRFC1951Block(blockData[2:], newBuffer, prevBytesDecompressed)
		if err != nil {
			return nil, fmt.Errorf("RFC1951ブロック#%d の解凍に失敗: 非圧縮サイズ=%d, 圧縮サイズ=%d, 解凍バイト数=%d, エラー=%v",
				blockIndex, blockUncompressedSize, blockCompressedSize, bytesDecompressed, err)
		}
		prevBytesDecompressed = bytesDecompressed

		// オフセット更新
		uncompressedSum += uint32(blockUncompressedSize)
		destOffset += int(blockUncompressedSize)
		srcOffset += int(blockCompressedSize)
	}

	return newBuffer, nil
}

// decompressRFC1951Block は圧縮データを復号化します。
// 圧縮データ compressed を展開先バッファ decompressed に展開し、展開後のバイト数を返します。
// 最終ブロックフラグとブロック種別（非圧縮/固定ハフマン/動的ハフマン）に基づいて処理を行います。
func (rep *XRepository) decompressRFC1951Block(compressed []byte, decompressed []byte, prevBytesDecompressed uint32) (uint32, error) {
	// 展開コンテキストの初期化
	ctx := newContext()

	// 前回の展開バイト数を設定
	ctx.bytesDecompressed = prevBytesDecompressed

	// 入出力バッファとサイズの設定
	ctx.src = compressed
	ctx.dest = decompressed
	ctx.sizeInBytes = uint32(len(compressed))
	ctx.destSize = uint32(len(decompressed))
	ctx.sizeInBits = int64(len(compressed)) * 8

	// 固定ハフマンツリーを準備
	rep.setupFixedHuffmanTree(ctx)

	// 圧縮データの復号化ループ
	final := false
	for !final && ctx.bitsRead < ctx.sizeInBits {
		// 最終ブロックフラグを取得
		if rep.getBits(ctx, 1) == 1 {
			final = true
		}

		// ブロックタイプを取得 (0=非圧縮, 1=固定ハフマン, 2=動的ハフマン)
		blockType := rep.getBits(ctx, 2)

		switch blockType {
		case 0: // 非圧縮ブロック
			if !rep.processUncompressed(ctx) {
				return 0, fmt.Errorf("非圧縮ブロックの処理に失敗しました")
			}
		case 1: // 固定ハフマン圧縮ブロック
			if !rep.processHuffmanFixed(ctx) {
				return 0, fmt.Errorf("固定ハフマンブロックの処理に失敗しました")
			}
		case 2: // 動的ハフマン圧縮ブロック
			if !rep.processHuffmanCustom(ctx) {
				return 0, fmt.Errorf("動的ハフマンブロックの処理に失敗しました")
			}
		default:
			return 0, fmt.Errorf("不明なブロック種別: %d", blockType)
		}
	}

	return ctx.bytesDecompressed, nil
}

// processHuffmanCustom は動的ハフマン圧縮ブロックを復号します。
// 動的ハフマン圧縮では、データに埋め込まれたコード長情報から
// カスタムハフマンツリーを構築し、それを用いてデータを復号します。
func (rep *XRepository) processHuffmanCustom(ctx *decompressionContext) bool {
	// カスタムハフマンツリーを構築
	if !rep.createHuffmanTreeTable(ctx) {
		return false
	}

	// リソース解放を保証
	defer rep.disposeCustomHuffmanTree(ctx)

	// ハフマンツリーを用いてデータを復号
	return rep.decodeWithHuffmanTree(ctx, &ctx.customLiteralLength, &ctx.customDistance)
}

// createHuffmanTreeTable は動的ハフマン圧縮ブロックで使用されるハフマンツリーを構築します。
// データストリームからコード長情報を読み取り、リテラル/長さツリーと距離ツリーを生成します。
func (rep *XRepository) createHuffmanTreeTable(ctx *decompressionContext) bool {
	// 1. パラメータを読み取る
	numLiteralLengthCodes := uint16(rep.getBits(ctx, 5)) + 257
	numDistanceCodes := uint16(rep.getBits(ctx, 5)) + 1
	numCodeLengthCodes := uint16(rep.getBits(ctx, 4)) + 4

	// 2. コード長のコード値配列を準備
	codeLengthsTree, ok := rep.buildCodeLengthsTree(ctx, numCodeLengthCodes)
	if !ok {
		return false
	}
	defer func() { codeLengthsTree.codes = nil }() // リソース解放

	// 3. リテラル/長さツリーの構築
	if !rep.buildSymbolTree(ctx, &ctx.customLiteralLength, codeLengthsTree,
		int(numLiteralLengthCodes), "リテラル/長さ") {
		return false
	}

	// 4. 距離ツリーの構築
	if !rep.buildSymbolTree(ctx, &ctx.customDistance, codeLengthsTree,
		int(numDistanceCodes), "距離") {
		rep.disposeCustomHuffmanTree(ctx)
		return false
	}

	return true
}

// buildCodeLengthsTree はコード長のコードを読み取り、コード長ツリーを構築します
func (rep *XRepository) buildCodeLengthsTree(ctx *decompressionContext, numCodeLengthCodes uint16) (*huffmanTree, bool) {
	// RFC1951仕様では、コード長コードは最大19個まで
	if numCodeLengthCodes > 19 {
		return nil, false
	}

	// RFC1951仕様で定義された順序
	canonicalOrder := []int{16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15}

	// コード長のコード値を初期化（19個まで）
	codeLengths := make([]byte, 19)
	for i := 0; i < int(numCodeLengthCodes); i++ {
		if i >= len(canonicalOrder) {
			// インデックス範囲外アクセスを防止
			return nil, false
		}
		codeLengths[canonicalOrder[i]] = byte(rep.getBits(ctx, 3))
	}

	// コード長からノードを生成
	nodes := make([]huffmanCode, 19)
	maxBits := 0
	for i := range 19 {
		nodes[i].length = uint16(codeLengths[i])
		nodes[i].binCode = uint16(i)
		if int(codeLengths[i]) > maxBits {
			maxBits = int(codeLengths[i])
		}
	}

	// ハフマンコードを割り当て
	if !rep.assignHuffmanCodes(nodes, maxBits) {
		return nil, false
	}

	// ツリーを構築
	var tree huffmanTree
	if !rep.setupHuffmanTree(&tree, nodes) {
		return nil, false
	}

	return &tree, true
}

// buildSymbolTree は指定されたコード長ツリーを使用してシンボルツリーを構築します
func (rep *XRepository) buildSymbolTree(ctx *decompressionContext, tree *huffmanTree,
	codeLengthsTree *huffmanTree, numCodes int, treeType string) bool {

	// シンボルコード配列を初期化
	symbolCodes := make([]huffmanCode, numCodes)

	// コード長の展開と読み取り
	var prev uint16 = 0xffff // 無効な値で初期化
	var repeat uint16 = 0
	maxBits := 0

	for i := range numCodes {
		symbolCodes[i].binCode = uint16(i)

		if repeat > 0 {
			// 繰り返しコードの処理
			symbolCodes[i].length = prev
			repeat--
		} else {
			// 新しいコード長を取得
			code := rep.getACodeWithHuffmanTree(ctx, codeLengthsTree)

			switch code {
			case 16: // 直前のコード長を繰り返し
				if prev == 0xffff {
					return false // 直前の値が無効
				}
				repeat = rep.getBits(ctx, 2) + 3
				symbolCodes[i].length = prev

			case 17: // 0を繰り返し (3-10回)
				repeat = rep.getBits(ctx, 3) + 2 // 1つ目は現在のiで処理するので+2
				prev = 0
				symbolCodes[i].length = 0

			case 18: // 0を繰り返し (11-138回)
				repeat = rep.getBits(ctx, 7) + 10 // 同上
				prev = 0
				symbolCodes[i].length = 0

			default:
				// 通常のコード長
				prev = code
				symbolCodes[i].length = code
				if int(code) > maxBits {
					maxBits = int(code)
				}
			}
		}
	}

	// 繰り返し処理が終わる前にループが終了した場合はエラー
	if repeat > 0 {
		return false
	}

	// ハフマンコードを割り当て
	if !rep.assignHuffmanCodes(symbolCodes, maxBits) {
		return false
	}

	// ツリーを設定
	return rep.setupHuffmanTree(tree, symbolCodes)
}

// assignHuffmanCodes は与えられたコード長に基づいてハフマンコードを割り当てます
func (rep *XRepository) assignHuffmanCodes(codes []huffmanCode, maxBits int) bool {
	if maxBits <= 0 {
		return true // コードが存在しない場合は成功扱い
	}

	// ビット長ごとのコード数をカウント
	bitLengths := make([]int, maxBits+1)
	for _, code := range codes {
		if code.length > 0 {
			bitLengths[code.length]++
		}
	}

	// 次のコード値を計算
	nextCode := make([]int, maxBits+1)
	code := 0
	bitLengths[0] = 0

	for bits := 1; bits <= maxBits; bits++ {
		code = (code + bitLengths[bits-1]) << 1
		nextCode[bits] = code
	}

	// コードに値を割り当て
	for i := range codes {
		length := int(codes[i].length)
		if length > 0 {
			codes[i].huffmanCode = uint16(nextCode[length])
			nextCode[length]++
		}
	}

	return true
}

// disposeCustomHuffmanTree は、DecompressionContext 内のカスタムハフマンツリー
// （CustomDistance と CustomLiteralLength）のリソースを解放（nil に設定）し、
// フィールドの値をリセットします。
func (rep *XRepository) disposeCustomHuffmanTree(ctx *decompressionContext) {
	// CustomDistance のリソース解放
	ctx.customDistance.codes = nil
	ctx.customDistance.numMaxBits = 0
	ctx.customDistance.numCodes = 0

	// CustomLiteralLength のリソース解放
	ctx.customLiteralLength.codes = nil
	ctx.customLiteralLength.numMaxBits = 0
	ctx.customLiteralLength.numCodes = 0
}

// processHuffmanFixed は、固定ハフマンブロックの復号処理を行います。
// 固定ハフマンツリー (ctx.Fixed) を用いて、DecodeWithHuffmanTree を呼び出し、
// 復号処理の成否を返します。
func (rep *XRepository) processHuffmanFixed(ctx *decompressionContext) bool {
	result := rep.decodeWithHuffmanTree(ctx, &ctx.fixed, nil)
	return result
}

// processUncompressed は非圧縮ブロックを処理します。
// サイズ情報の読み取り、チェックサム検証、データの出力バッファへのコピーを行います。
func (rep *XRepository) processUncompressed(ctx *decompressionContext) bool {
	// 1. ブロックサイズを読み取る (2バイト、リトルエンディアン)
	blockSize := uint16(rep.getBits(ctx, 8)) | (uint16(rep.getBits(ctx, 8)) << 8)

	// 2. チェックサムを検証（サイズのビット反転値）
	checksum := uint16(rep.getBits(ctx, 8)) | (uint16(rep.getBits(ctx, 8)) << 8)
	if ^checksum != blockSize {
		return false
	}

	// 3. 出力バッファの容量確認
	if int(ctx.bytesDecompressed)+int(blockSize) > len(ctx.dest) {
		return false
	}

	// 4. データを読み取り、出力バッファにコピー
	for range int(blockSize) {
		byteValue := byte(rep.getBits(ctx, 8))
		ctx.dest[ctx.bytesDecompressed] = byteValue
		ctx.bytesDecompressed++
	}

	return true
}

// decodeWithHuffmanTree はハフマンツリーを使用して圧縮データを復号します。
// literalLengthTree でリテラル値と長さを、distanceTree で距離を復号します。
// 出力はctx.Destバッファに書き込まれ、終端コード (256) が検出されるまで処理を続けます。
func (rep *XRepository) decodeWithHuffmanTree(ctx *decompressionContext, literalLengthTree *huffmanTree, distanceTree *huffmanTree) bool {
	for {
		// コードを取得
		code := rep.getACodeWithHuffmanTree(ctx, literalLengthTree)

		// エラーチェック
		if code == 0xffff {
			return false
		}

		// 終端コード
		if code == 256 {
			return true
		}

		// リテラル値の処理
		if code < 256 {
			if int(ctx.bytesDecompressed) >= len(ctx.dest) {
				return false // 出力バッファ不足
			}

			// バイト値を出力
			ctx.dest[ctx.bytesDecompressed] = byte(code)
			ctx.bytesDecompressed++
			continue
		}

		// 長さ/距離ペアの処理（バックリファレンス）
		if !rep.copyBackReference(ctx, code, distanceTree) {
			return false
		}
	}
}

// copyBackReference は長さ/距離ペアを処理して過去のデータをコピーします。
// code から長さを、distanceTree から距離を復号し、過去データを現在位置にコピーします。
func (rep *XRepository) copyBackReference(ctx *decompressionContext, lengthCode uint16, distanceTree *huffmanTree) bool {
	// 長さを復号
	length := rep.decodeLength(ctx, lengthCode)
	if length == 0xffff { // エラー値
		return false
	}

	// 距離コードを取得
	var distCode uint16
	if distanceTree != nil {
		distCode = rep.getACodeWithHuffmanTree(ctx, distanceTree)
	} else {
		distCode = rep.reverseBits(int(rep.getBits(ctx, 5)), 5)
	}

	// 距離を復号
	distance := rep.decodeDistance(ctx, distCode)
	if distance == 0 { // エラー値
		return false
	}

	// コピー元位置を計算
	sourcePos := int(ctx.bytesDecompressed) - distance
	if sourcePos < 0 {
		return false // 不正な距離
	}

	// 出力バッファの境界チェック
	if int(ctx.bytesDecompressed)+length > len(ctx.dest) ||
		sourcePos+length > len(ctx.dest) {
		return false
	}

	// 最適化：コピーの効率化
	destPos := int(ctx.bytesDecompressed)

	// 短いコピーはインライン展開
	if length <= 16 {
		// 4バイト単位でアンロール
		i := 0
		for ; i+4 <= length; i += 4 {
			ctx.dest[destPos+i] = ctx.dest[sourcePos+i]
			ctx.dest[destPos+i+1] = ctx.dest[sourcePos+i+1]
			ctx.dest[destPos+i+2] = ctx.dest[sourcePos+i+2]
			ctx.dest[destPos+i+3] = ctx.dest[sourcePos+i+3]
		}
		// 残りのバイト
		for ; i < length; i++ {
			ctx.dest[destPos+i] = ctx.dest[sourcePos+i]
		}
	} else if sourcePos+length <= destPos {
		// ソースとデスティネーションが重ならない場合、組み込み関数を使用
		copy(ctx.dest[destPos:destPos+length], ctx.dest[sourcePos:sourcePos+length])
	} else {
		// 重なる可能性がある場合は1バイトずつ安全にコピー
		for i := range length {
			ctx.dest[destPos+i] = ctx.dest[sourcePos+i]
		}
	}

	// バイト数を更新
	ctx.bytesDecompressed += uint32(length)

	return true
}

// decodeLength はRFC1951規格に基づき、リテラル/長さコード(baseCode)から実際の長さ値を算出します。
// baseCodeの値に応じて以下のように計算します:
// - 257-264: 直接計算式(baseCode - 257 + 3)で長さを求める
// - 265-284: 追加ビットを使用した計算で長さを求める
// - 285: 固定値258を返す
// - 286以上: 不正な値として0xffffを返す
func (rep *XRepository) decodeLength(ctx *decompressionContext, baseCode uint16) int {
	// 基本チェック
	if baseCode < 257 || baseCode > 285 {
		return 0xffff // エラー値
	}

	// テーブル参照のために調整したインデックス
	idx := baseCode - 257

	// 基本値を取得
	baseValue := int(lengthBaseValues[idx])

	// 追加ビット数を取得
	extraBits := int(lengthBitOffsets[idx])

	// 追加ビットがある場合は読み取って加算
	if extraBits > 0 {
		baseValue += int(rep.getBits(ctx, extraBits))
	}

	return baseValue
}

// decodeDistance はRFC1951規格に基づき、距離コード(baseCode)から実際の距離値を算出します。
// 距離コードは以下の3つのケースで処理されます:
// - 0-3: 直接変換 (baseCode + 1)
// - 4-29: 追加ビットを使用した複雑な計算
// - 30以上: 無効なコード (0を返す)
func (rep *XRepository) decodeDistance(ctx *decompressionContext, baseCode uint16) int {
	// 範囲チェック
	if baseCode > 29 {
		return 0 // エラー値
	}

	// 基本距離値を取得
	baseValue := int(distanceBaseValues[baseCode])

	// 追加ビット数を取得
	extraBits := int(distanceBitOffsets[baseCode])

	// 追加ビットがある場合は読み取って加算
	if extraBits > 0 {
		baseValue += int(rep.getBits(ctx, extraBits))
	}

	return baseValue
}

// compareHuffmanCode は、ハフマンコードのバイナリサーチのための比較関数です。
// まず長さ(Length)を比較し、等しい場合はコード値(HuffmanCode)を比較します。
// 戻り値:
//
//	-1: aが小さい（長さが短いか、同じ長さでコード値が小さい）
//	 0: aとパラメータが等しい
//	 1: aが大きい（長さが長いか、同じ長さでコード値が大きい）
func (rep *XRepository) compareHuffmanCode(a huffmanCode, length int, code uint16) int {
	// 長さの比較
	if aLen := int(a.length); aLen != length {
		if aLen < length {
			return -1
		}
		return 1
	}

	// 長さが同じ場合、コード値を比較
	if a.huffmanCode < code {
		return -1
	}
	if a.huffmanCode > code {
		return 1
	}

	// 両方とも一致
	return 0
}

// getACodeWithHuffmanTree はハフマン木からビットストリームに対応するシンボルを検索します。
// ビットを次々と読み出しながら、該当するコードを見つけるまで検索を続けます。
// 見つからなかった場合は0xffffを返します。
func (rep *XRepository) getACodeWithHuffmanTree(ctx *decompressionContext, tree *huffmanTree) uint16 {
	// ハフマン木の有効性チェック
	if tree.codes == nil || len(tree.codes) == 0 {
		return 0xffff
	}

	// 定数と初期値の設定
	maxBits := int(tree.numMaxBits)
	minBits := int(tree.codes[0].length)

	// 初期ビットパターンの読み込み
	bitPattern := uint16(0)
	for i := 0; i < minBits; i++ {
		bitPattern = (bitPattern << 1) | rep.readNextBit(ctx)
	}

	// 現在のビット長から最大ビット長まで検索
	for bitLength := minBits; bitLength <= maxBits; bitLength++ {
		// バイナリサーチでコードを検索
		result := rep.searchHuffmanCode(tree, bitLength, bitPattern)
		if result != 0xffff {
			return result
		}

		// 見つからなかった場合、1ビット追加して再検索
		bitPattern = (bitPattern << 1) | rep.readNextBit(ctx)
	}

	// 該当するコードが見つからなかった
	return 0xffff
}

// searchHuffmanCode はハフマン木から特定のビット長とパターンに一致するコードをバイナリサーチで探します。
// 見つかった場合はそのシンボル（BinCode）を返し、見つからなかった場合は0xffffを返します。
func (rep *XRepository) searchHuffmanCode(tree *huffmanTree, bitLength int, bitPattern uint16) uint16 {
	// 固定ハフマンコードの場合はマップを使用して高速検索
	if tree.numCodes == 288 && tree.numMaxBits == 9 && len(tree.codes) > 0 && tree.codes[0].binCode == 0x100 {
		key := (uint32(bitLength) << 16) | uint32(bitPattern)
		if binCode, found := fixedHuffmanMap[key]; found {
			return binCode
		}
		return 0xffff
	}

	// 以下は既存のバイナリサーチ
	left := 0
	right := int(tree.numCodes)

	for left < right {
		mid := (left + right) / 2
		comp := rep.compareHuffmanCode(tree.codes[mid], bitLength, bitPattern)

		if comp == 0 {
			return tree.codes[mid].binCode
		} else if comp < 0 {
			left = mid + 1
		} else {
			right = mid
		}
	}

	return 0xffff
}

// setupHuffmanTree は有効なハフマンコードを抽出し、特定の順序でソートしてツリーを構築します。
// コード長の降順（同じ長さの場合はコード値の昇順）でソートします。
// 戻り値は常にtrueです（エラーケースはありません）。
func (rep *XRepository) setupHuffmanTree(tree *huffmanTree, codes []huffmanCode) bool {
	// 作業用のコード配列を準備
	sortedCodes := make([]huffmanCode, len(codes))
	validCount := 0
	maxBitLength := 0

	// 有効なコード（Length > 0）のみを抽出して挿入ソート
	for _, code := range codes {
		bitLength := int(code.length)
		if bitLength == 0 {
			continue // 無効なコードはスキップ
		}

		// 挿入位置を決定（降順にソート）
		insertPos := validCount
		huffValue := int(code.huffmanCode)

		// 挿入ソート: 長さが長い順、同じ長さならハフマンコード値の昇順
		for insertPos > 0 {
			prev := sortedCodes[insertPos-1]
			prevLength := int(prev.length)

			if prevLength < bitLength {
				break // 現在のコードの方が長い
			}
			if prevLength == bitLength && int(prev.huffmanCode) <= huffValue {
				break // 同じ長さで前のコードが小さいか等しい
			}

			// 要素をシフト
			sortedCodes[insertPos] = sortedCodes[insertPos-1]
			insertPos--
		}

		// 適切な位置に挿入
		sortedCodes[insertPos].length = uint16(bitLength)
		sortedCodes[insertPos].huffmanCode = uint16(huffValue)
		sortedCodes[insertPos].binCode = code.binCode

		// カウンタと最大ビット長を更新
		validCount++
		if bitLength > maxBitLength {
			maxBitLength = bitLength
		}
	}

	// ツリー情報を設定
	tree.numMaxBits = uint16(maxBitLength)
	tree.numCodes = uint16(validCount)
	tree.codes = sortedCodes[:validCount]

	return true
}

// setupFixedHuffmanTree を最適化
func (rep *XRepository) setupFixedHuffmanTree(ctx *decompressionContext) {
	// 既に初期化済みのグローバルインスタンスを参照するだけ
	ctx.fixed = fixedHuffmanTreeInstance
}

// reverseBits は入力値 input の下位 n ビットを反転した結果を返します。
// 例: input=0b1011, n=4 の場合、反転結果は 0b1101 となります。
// この関数はRFC1951圧縮形式で使用されるビット反転操作を実装しています。
func (rep *XRepository) reverseBits(input int, n int) uint16 {
	// 小さな値のための高速パス
	if n <= 8 {
		// 8ビットまでのルックアップテーブルを使用
		switch n {
		case 1:
			return uint16(input & 1)
		case 2:
			return uint16(((input & 1) << 1) | ((input >> 1) & 1))
		case 3:
			return uint16(((input & 1) << 2) | ((input & 2) << 0) | ((input & 4) >> 2))
		case 5: // 距離コードで特に使用される
			// 5ビット専用の高速ルックアップ
			return uint16((input&1)<<4 | (input&2)<<2 | (input&4)<<0 | (input&8)>>2 | (input&16)>>4)
		}
	}

	// 通常の反転処理
	var result uint16 = 0
	for i := 0; i < n; i++ {
		result = (result << 1) | uint16(input&1)
		input >>= 1
	}
	return result
}

// decodeLength および decodeDistance 関数をテーブル参照で高速化
// readNextBit はバッファから1ビットを読み取ります。
// バッファが空の場合は新たにバイトを読み込みます。
func (rep *XRepository) readNextBit(ctx *decompressionContext) uint16 {
	// バッファが空の場合、32/16/8ビット単位で効率的に読み込む
	if ctx.bitCount == 0 {
		byteIndex := ctx.bitsRead >> 3
		if int(byteIndex) >= len(ctx.src) {
			return 0
		}

		// 残りバイト数に応じて最適な読み込み方を選択
		remainingBytes := len(ctx.src) - int(byteIndex)
		if remainingBytes >= 4 {
			// 4バイト一括読み込み - アラインメント考慮なし（速度優先）
			ctx.bitBuffer = uint32(ctx.src[byteIndex]) |
				uint32(ctx.src[byteIndex+1])<<8 |
				uint32(ctx.src[byteIndex+2])<<16 |
				uint32(ctx.src[byteIndex+3])<<24
			ctx.bitCount = 32
		} else if remainingBytes >= 2 {
			ctx.bitBuffer = uint32(ctx.src[byteIndex]) | uint32(ctx.src[byteIndex+1])<<8
			ctx.bitCount = 16
		} else {
			ctx.bitBuffer = uint32(ctx.src[byteIndex])
			ctx.bitCount = 8
		}
	}

	// LSBを抽出（最適化：ビット操作を最小化）
	bit := ctx.bitBuffer & 1
	ctx.bitBuffer >>= 1
	ctx.bitCount--
	ctx.bitsRead++

	return uint16(bit)
}

func (rep *XRepository) getBits(ctx *decompressionContext, size int) uint16 {
	// サイズが0または無効な場合は早期リターン
	if size <= 0 || size > 16 {
		return 0
	}

	// 最適化：全てのビットがバッファにある場合はマスク1回で処理
	if ctx.bitCount >= size {
		result := uint16(ctx.bitBuffer & bitMasks[size])
		ctx.bitBuffer >>= size
		ctx.bitCount -= size
		ctx.bitsRead += int64(size)
		return result
	}

	// 部分的なビットを取得して結合する必要がある場合
	result := uint16(0)

	// バッファに残っているビットを使用
	if ctx.bitCount > 0 {
		result = uint16(ctx.bitBuffer & bitMasks[ctx.bitCount])
		size -= ctx.bitCount
		ctx.bitsRead += int64(ctx.bitCount)
		ctx.bitCount = 0

		// バイト境界からの読み込み
		byteIndex := ctx.bitsRead >> 3
		if int(byteIndex) >= len(ctx.src) {
			return result
		}

		// 最適化: サイズに応じて一度に読み込むバイト数を調整
		if int(byteIndex)+4 <= len(ctx.src) && size > 8 {
			// 4バイト読み込み
			ctx.bitBuffer = uint32(ctx.src[byteIndex]) |
				uint32(ctx.src[byteIndex+1])<<8 |
				uint32(ctx.src[byteIndex+2])<<16 |
				uint32(ctx.src[byteIndex+3])<<24

			// 必要なビット数だけ取得
			additionalBits := uint16(ctx.bitBuffer & bitMasks[size])
			result |= additionalBits << uint16(ctx.bitCount)

			// バッファを更新
			ctx.bitBuffer >>= size
			ctx.bitCount = 32 - size
			ctx.bitsRead += int64(size)
			return result
		} else if int(byteIndex)+2 <= len(ctx.src) && size > 8 {
			// 2バイト読み込み
			ctx.bitBuffer = uint32(ctx.src[byteIndex]) | uint32(ctx.src[byteIndex+1])<<8
			if size <= 16 {
				additionalBits := uint16(ctx.bitBuffer & bitMasks[size])
				result |= additionalBits << uint16(ctx.bitCount)
				ctx.bitBuffer >>= size
				ctx.bitCount = 16 - size
				ctx.bitsRead += int64(size)
				return result
			}
		} else if int(byteIndex) < len(ctx.src) {
			// 1バイト読み込み
			ctx.bitBuffer = uint32(ctx.src[byteIndex])
			if size <= 8 {
				additionalBits := uint16(ctx.bitBuffer & bitMasks[size])
				result |= additionalBits << uint16(ctx.bitCount)
				ctx.bitBuffer >>= size
				ctx.bitCount = 8 - size
				ctx.bitsRead += int64(size)
				return result
			}
		}
	}

	// バイト単位での読み込みが必要な場合
	for size >= 8 {
		byteIndex := ctx.bitsRead >> 3
		if int(byteIndex) >= len(ctx.src) {
			return result
		}

		byteValue := ctx.src[byteIndex]
		result = (result << 8) | uint16(byteValue)
		size -= 8
		ctx.bitsRead += 8
	}

	// 残りのビット
	if size > 0 {
		byteIndex := ctx.bitsRead >> 3
		if int(byteIndex) >= len(ctx.src) {
			return result
		}

		ctx.bitBuffer = uint32(ctx.src[byteIndex])
		ctx.bitCount = 8

		mask := bitMasks[size]
		result = (result << size) | uint16(ctx.bitBuffer&mask)

		ctx.bitBuffer >>= size
		ctx.bitCount -= size
		ctx.bitsRead += int64(size)
	}

	return result
}

// -----------------------------------------------------------------------------

// newContext は、初期化された新しいDecompressionContextを返します。
// 全てのフィールドは初期値（ゼロ値）で設定されます。
func newContext() *decompressionContext {
	// 新しいコンテキストを作成
	return &decompressionContext{
		src:               nil,
		bitsRead:          0,
		sizeInBits:        0,
		sizeInBytes:       0,
		dest:              nil,
		destSize:          0,
		bytesDecompressed: 0,
		bitBuffer:         0,
		bitCount:          0,
		customDistance: huffmanTree{
			numMaxBits: 0,
			codes:      nil,
			numCodes:   0,
		},
		customLiteralLength: huffmanTree{
			numMaxBits: 0,
			codes:      nil,
			numCodes:   0,
		},
		fixed: huffmanTree{
			numMaxBits: 0,
			codes:      nil,
			numCodes:   0,
		},
	}
}

// huffmanCode 構造体
// C++:
//
//	typedef struct _customhuffmantable{
//	    WORD    binCode;
//	    WORD    length;
//	    WORD    huffmanCode;
//	} huffmanCode;
type huffmanCode struct {
	binCode     uint16
	length      uint16
	huffmanCode uint16
}

// huffmanTree 構造体
// C++:
//
//	typedef struct{
//	    WORD numMaxBits;
//	    HuffmanCode *pCode;
//	    WORD numCodes;
//	} huffmanTree;
type huffmanTree struct {
	numMaxBits uint16
	codes      []huffmanCode
	numCodes   uint16
}

// customHuffmanTree 構造体
// C++:
//
//	typedef struct{
//	    WORD numLiteralLengthCodes;
//	    WORD numDistanceCodes;
//	    WORD *pLiteralLengthTree;
//	    WORD numElementsOfLiteralLengthTree;
//	    WORD *pDistanceTree;
//	    WORD numElementsOfDistanceTree;
//	} customHuffmanTree;
type customHuffmanTree struct {
	numLiteralLengthCodes          uint16
	numDistanceCodes               uint16
	literalLengthTree              []uint16
	numElementsOfLiteralLengthTree uint16
	distanceTree                   []uint16
	numElementsOfDistanceTree      uint16
}

// decompressionContext 構造体
// C++:
//
//	typedef struct _context{
//	    BYTE    *pSrc;
//	    __int64 i64BitsRead;
//	    __int64 i64SizeInBits;
//	    DWORD   dwSizeInBytes;
//	    BYTE    *pDest;
//	    DWORD   dwDestSize;
//	    DWORD   dwBytesDecompressed;
//	    HuffmanTree customLiteralLength;
//	    HuffmanTree customDistance;
//	    HuffmanTree fixed;
//	} decompressionContext;
type decompressionContext struct {
	src                 []byte
	bitsRead            int64
	sizeInBits          int64
	sizeInBytes         uint32
	dest                []byte
	destSize            uint32
	bytesDecompressed   uint32
	customLiteralLength huffmanTree
	customDistance      huffmanTree
	fixed               huffmanTree
	bitBuffer           uint32 // 現在のビットバッファ
	bitCount            int    // バッファ内の有効ビット数
}

// 使用頻度の高いビットオペレーション定数をプリコンパイル
var (
	// ビットマスクテーブル (上記で定義済み)
	bitMasks = [33]uint32{
		0, 0x1, 0x3, 0x7, 0xF, 0x1F, 0x3F, 0x7F, 0xFF,
		0x1FF, 0x3FF, 0x7FF, 0xFFF, 0x1FFF, 0x3FFF, 0x7FFF, 0xFFFF,
		0x1FFFF, 0x3FFFF, 0x7FFFF, 0xFFFFF, 0x1FFFFF, 0x3FFFFF, 0x7FFFFF, 0xFFFFFF,
		0x1FFFFFF, 0x3FFFFFF, 0x7FFFFFF, 0xFFFFFFF, 0x1FFFFFFF, 0x3FFFFFFF, 0x7FFFFFFF, 0xFFFFFFFF,
	}
	// オペレーションごとの共通パターン用のテーブル
	lengthBitOffsets   = [29]uint16{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 0}
	distanceBitOffsets = [30]uint16{0, 0, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11, 12, 12, 13, 13}
	lengthBaseValues   = [29]uint16{3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 15, 17, 19, 23, 27, 31, 35, 43, 51, 59, 67, 83, 99, 115, 131, 163, 195, 227, 258}
	distanceBaseValues = [30]uint16{1, 2, 3, 4, 5, 7, 9, 13, 17, 25, 33, 49, 65, 97, 129, 193, 257, 385, 513, 769, 1025, 1537, 2049, 3073, 4097, 6145, 8193, 12289, 16385, 24577}
)

// fixedHuffmanMap は固定ハフマンコードの高速検索用マップ
// キーは (length << 16 | huffmanCode) の形式で、値はbinCode
var fixedHuffmanMap map[uint32]uint16

var fixedCodeSorted = []huffmanCode{
	{binCode: 0x100, length: 7, huffmanCode: 0x0},
	{binCode: 0x101, length: 7, huffmanCode: 0x1},
	{binCode: 0x102, length: 7, huffmanCode: 0x2},
	{binCode: 0x103, length: 7, huffmanCode: 0x3},
	{binCode: 0x104, length: 7, huffmanCode: 0x4},
	{binCode: 0x105, length: 7, huffmanCode: 0x5},
	{binCode: 0x106, length: 7, huffmanCode: 0x6},
	{binCode: 0x107, length: 7, huffmanCode: 0x7},
	{binCode: 0x108, length: 7, huffmanCode: 0x8},
	{binCode: 0x109, length: 7, huffmanCode: 0x9},
	{binCode: 0x10a, length: 7, huffmanCode: 0xa},
	{binCode: 0x10b, length: 7, huffmanCode: 0xb},
	{binCode: 0x10c, length: 7, huffmanCode: 0xc},
	{binCode: 0x10d, length: 7, huffmanCode: 0xd},
	{binCode: 0x10e, length: 7, huffmanCode: 0xe},
	{binCode: 0x10f, length: 7, huffmanCode: 0xf},
	{binCode: 0x110, length: 7, huffmanCode: 0x10},
	{binCode: 0x111, length: 7, huffmanCode: 0x11},
	{binCode: 0x112, length: 7, huffmanCode: 0x12},
	{binCode: 0x113, length: 7, huffmanCode: 0x13},
	{binCode: 0x114, length: 7, huffmanCode: 0x14},
	{binCode: 0x115, length: 7, huffmanCode: 0x15},
	{binCode: 0x116, length: 7, huffmanCode: 0x16},
	{binCode: 0x117, length: 7, huffmanCode: 0x17},

	{binCode: 0x0, length: 8, huffmanCode: 0x30},
	{binCode: 0x1, length: 8, huffmanCode: 0x31},
	{binCode: 0x2, length: 8, huffmanCode: 0x32},
	{binCode: 0x3, length: 8, huffmanCode: 0x33},
	{binCode: 0x4, length: 8, huffmanCode: 0x34},
	{binCode: 0x5, length: 8, huffmanCode: 0x35},
	{binCode: 0x6, length: 8, huffmanCode: 0x36},
	{binCode: 0x7, length: 8, huffmanCode: 0x37},
	{binCode: 0x8, length: 8, huffmanCode: 0x38},
	{binCode: 0x9, length: 8, huffmanCode: 0x39},
	{binCode: 0xa, length: 8, huffmanCode: 0x3a},
	{binCode: 0xb, length: 8, huffmanCode: 0x3b},
	{binCode: 0xc, length: 8, huffmanCode: 0x3c},
	{binCode: 0xd, length: 8, huffmanCode: 0x3d},
	{binCode: 0xe, length: 8, huffmanCode: 0x3e},
	{binCode: 0xf, length: 8, huffmanCode: 0x3f},
	{binCode: 0x10, length: 8, huffmanCode: 0x40},
	{binCode: 0x11, length: 8, huffmanCode: 0x41},
	{binCode: 0x12, length: 8, huffmanCode: 0x42},
	{binCode: 0x13, length: 8, huffmanCode: 0x43},
	{binCode: 0x14, length: 8, huffmanCode: 0x44},
	{binCode: 0x15, length: 8, huffmanCode: 0x45},
	{binCode: 0x16, length: 8, huffmanCode: 0x46},
	{binCode: 0x17, length: 8, huffmanCode: 0x47},
	{binCode: 0x18, length: 8, huffmanCode: 0x48},
	{binCode: 0x19, length: 8, huffmanCode: 0x49},
	{binCode: 0x1a, length: 8, huffmanCode: 0x4a},
	{binCode: 0x1b, length: 8, huffmanCode: 0x4b},
	{binCode: 0x1c, length: 8, huffmanCode: 0x4c},
	{binCode: 0x1d, length: 8, huffmanCode: 0x4d},
	{binCode: 0x1e, length: 8, huffmanCode: 0x4e},
	{binCode: 0x1f, length: 8, huffmanCode: 0x4f},
	{binCode: 0x20, length: 8, huffmanCode: 0x50},
	{binCode: 0x21, length: 8, huffmanCode: 0x51},
	{binCode: 0x22, length: 8, huffmanCode: 0x52},
	{binCode: 0x23, length: 8, huffmanCode: 0x53},
	{binCode: 0x24, length: 8, huffmanCode: 0x54},
	{binCode: 0x25, length: 8, huffmanCode: 0x55},
	{binCode: 0x26, length: 8, huffmanCode: 0x56},
	{binCode: 0x27, length: 8, huffmanCode: 0x57},
	{binCode: 0x28, length: 8, huffmanCode: 0x58},
	{binCode: 0x29, length: 8, huffmanCode: 0x59},
	{binCode: 0x2a, length: 8, huffmanCode: 0x5a},
	{binCode: 0x2b, length: 8, huffmanCode: 0x5b},
	{binCode: 0x2c, length: 8, huffmanCode: 0x5c},
	{binCode: 0x2d, length: 8, huffmanCode: 0x5d},
	{binCode: 0x2e, length: 8, huffmanCode: 0x5e},
	{binCode: 0x2f, length: 8, huffmanCode: 0x5f},
	{binCode: 0x30, length: 8, huffmanCode: 0x60},
	{binCode: 0x31, length: 8, huffmanCode: 0x61},
	{binCode: 0x32, length: 8, huffmanCode: 0x62},
	{binCode: 0x33, length: 8, huffmanCode: 0x63},
	{binCode: 0x34, length: 8, huffmanCode: 0x64},
	{binCode: 0x35, length: 8, huffmanCode: 0x65},
	{binCode: 0x36, length: 8, huffmanCode: 0x66},
	{binCode: 0x37, length: 8, huffmanCode: 0x67},
	{binCode: 0x38, length: 8, huffmanCode: 0x68},
	{binCode: 0x39, length: 8, huffmanCode: 0x69},
	{binCode: 0x3a, length: 8, huffmanCode: 0x6a},
	{binCode: 0x3b, length: 8, huffmanCode: 0x6b},
	{binCode: 0x3c, length: 8, huffmanCode: 0x6c},
	{binCode: 0x3d, length: 8, huffmanCode: 0x6d},
	{binCode: 0x3e, length: 8, huffmanCode: 0x6e},
	{binCode: 0x3f, length: 8, huffmanCode: 0x6f},
	{binCode: 0x40, length: 8, huffmanCode: 0x70},
	{binCode: 0x41, length: 8, huffmanCode: 0x71},
	{binCode: 0x42, length: 8, huffmanCode: 0x72},
	{binCode: 0x43, length: 8, huffmanCode: 0x73},
	{binCode: 0x44, length: 8, huffmanCode: 0x74},
	{binCode: 0x45, length: 8, huffmanCode: 0x75},
	{binCode: 0x46, length: 8, huffmanCode: 0x76},
	{binCode: 0x47, length: 8, huffmanCode: 0x77},
	{binCode: 0x48, length: 8, huffmanCode: 0x78},
	{binCode: 0x49, length: 8, huffmanCode: 0x79},
	{binCode: 0x4a, length: 8, huffmanCode: 0x7a},
	{binCode: 0x4b, length: 8, huffmanCode: 0x7b},
	{binCode: 0x4c, length: 8, huffmanCode: 0x7c},
	{binCode: 0x4d, length: 8, huffmanCode: 0x7d},
	{binCode: 0x4e, length: 8, huffmanCode: 0x7e},
	{binCode: 0x4f, length: 8, huffmanCode: 0x7f},
	{binCode: 0x50, length: 8, huffmanCode: 0x80},
	{binCode: 0x51, length: 8, huffmanCode: 0x81},
	{binCode: 0x52, length: 8, huffmanCode: 0x82},
	{binCode: 0x53, length: 8, huffmanCode: 0x83},
	{binCode: 0x54, length: 8, huffmanCode: 0x84},
	{binCode: 0x55, length: 8, huffmanCode: 0x85},
	{binCode: 0x56, length: 8, huffmanCode: 0x86},
	{binCode: 0x57, length: 8, huffmanCode: 0x87},
	{binCode: 0x58, length: 8, huffmanCode: 0x88},
	{binCode: 0x59, length: 8, huffmanCode: 0x89},
	{binCode: 0x5a, length: 8, huffmanCode: 0x8a},
	{binCode: 0x5b, length: 8, huffmanCode: 0x8b},
	{binCode: 0x5c, length: 8, huffmanCode: 0x8c},
	{binCode: 0x5d, length: 8, huffmanCode: 0x8d},
	{binCode: 0x5e, length: 8, huffmanCode: 0x8e},
	{binCode: 0x5f, length: 8, huffmanCode: 0x8f},
	{binCode: 0x60, length: 8, huffmanCode: 0x90},
	{binCode: 0x61, length: 8, huffmanCode: 0x91},
	{binCode: 0x62, length: 8, huffmanCode: 0x92},
	{binCode: 0x63, length: 8, huffmanCode: 0x93},
	{binCode: 0x64, length: 8, huffmanCode: 0x94},
	{binCode: 0x65, length: 8, huffmanCode: 0x95},
	{binCode: 0x66, length: 8, huffmanCode: 0x96},
	{binCode: 0x67, length: 8, huffmanCode: 0x97},
	{binCode: 0x68, length: 8, huffmanCode: 0x98},
	{binCode: 0x69, length: 8, huffmanCode: 0x99},
	{binCode: 0x6a, length: 8, huffmanCode: 0x9a},
	{binCode: 0x6b, length: 8, huffmanCode: 0x9b},
	{binCode: 0x6c, length: 8, huffmanCode: 0x9c},
	{binCode: 0x6d, length: 8, huffmanCode: 0x9d},
	{binCode: 0x6e, length: 8, huffmanCode: 0x9e},
	{binCode: 0x6f, length: 8, huffmanCode: 0x9f},
	{binCode: 0x70, length: 8, huffmanCode: 0xa0},
	{binCode: 0x71, length: 8, huffmanCode: 0xa1},
	{binCode: 0x72, length: 8, huffmanCode: 0xa2},
	{binCode: 0x73, length: 8, huffmanCode: 0xa3},
	{binCode: 0x74, length: 8, huffmanCode: 0xa4},
	{binCode: 0x75, length: 8, huffmanCode: 0xa5},
	{binCode: 0x76, length: 8, huffmanCode: 0xa6},
	{binCode: 0x77, length: 8, huffmanCode: 0xa7},
	{binCode: 0x78, length: 8, huffmanCode: 0xa8},
	{binCode: 0x79, length: 8, huffmanCode: 0xa9},
	{binCode: 0x7a, length: 8, huffmanCode: 0xaa},
	{binCode: 0x7b, length: 8, huffmanCode: 0xab},
	{binCode: 0x7c, length: 8, huffmanCode: 0xac},
	{binCode: 0x7d, length: 8, huffmanCode: 0xad},
	{binCode: 0x7e, length: 8, huffmanCode: 0xae},
	{binCode: 0x7f, length: 8, huffmanCode: 0xaf},
	{binCode: 0x80, length: 8, huffmanCode: 0xb0},
	{binCode: 0x81, length: 8, huffmanCode: 0xb1},
	{binCode: 0x82, length: 8, huffmanCode: 0xb2},
	{binCode: 0x83, length: 8, huffmanCode: 0xb3},
	{binCode: 0x84, length: 8, huffmanCode: 0xb4},
	{binCode: 0x85, length: 8, huffmanCode: 0xb5},
	{binCode: 0x86, length: 8, huffmanCode: 0xb6},
	{binCode: 0x87, length: 8, huffmanCode: 0xb7},
	{binCode: 0x88, length: 8, huffmanCode: 0xb8},
	{binCode: 0x89, length: 8, huffmanCode: 0xb9},
	{binCode: 0x8a, length: 8, huffmanCode: 0xba},
	{binCode: 0x8b, length: 8, huffmanCode: 0xbb},
	{binCode: 0x8c, length: 8, huffmanCode: 0xbc},
	{binCode: 0x8d, length: 8, huffmanCode: 0xbd},
	{binCode: 0x8e, length: 8, huffmanCode: 0xbe},
	{binCode: 0x8f, length: 8, huffmanCode: 0xbf},
	{binCode: 0x118, length: 8, huffmanCode: 0xc0},
	{binCode: 0x119, length: 8, huffmanCode: 0xc1},
	{binCode: 0x11a, length: 8, huffmanCode: 0xc2},
	{binCode: 0x11b, length: 8, huffmanCode: 0xc3},
	{binCode: 0x11c, length: 8, huffmanCode: 0xc4},
	{binCode: 0x11d, length: 8, huffmanCode: 0xc5},
	{binCode: 0x11e, length: 8, huffmanCode: 0xc6},
	{binCode: 0x11f, length: 8, huffmanCode: 0xc7},
	{binCode: 0x90, length: 9, huffmanCode: 0x190},
	{binCode: 0x91, length: 9, huffmanCode: 0x191},
	{binCode: 0x92, length: 9, huffmanCode: 0x192},
	{binCode: 0x93, length: 9, huffmanCode: 0x193},
	{binCode: 0x94, length: 9, huffmanCode: 0x194},
	{binCode: 0x95, length: 9, huffmanCode: 0x195},
	{binCode: 0x96, length: 9, huffmanCode: 0x196},
	{binCode: 0x97, length: 9, huffmanCode: 0x197},
	{binCode: 0x98, length: 9, huffmanCode: 0x198},
	{binCode: 0x99, length: 9, huffmanCode: 0x199},
	{binCode: 0x9a, length: 9, huffmanCode: 0x19a},
	{binCode: 0x9b, length: 9, huffmanCode: 0x19b},
	{binCode: 0x9c, length: 9, huffmanCode: 0x19c},
	{binCode: 0x9d, length: 9, huffmanCode: 0x19d},
	{binCode: 0x9e, length: 9, huffmanCode: 0x19e},
	{binCode: 0x9f, length: 9, huffmanCode: 0x19f},
	{binCode: 0xa0, length: 9, huffmanCode: 0x1a0},
	{binCode: 0xa1, length: 9, huffmanCode: 0x1a1},
	{binCode: 0xa2, length: 9, huffmanCode: 0x1a2},
	{binCode: 0xa3, length: 9, huffmanCode: 0x1a3},
	{binCode: 0xa4, length: 9, huffmanCode: 0x1a4},
	{binCode: 0xa5, length: 9, huffmanCode: 0x1a5},
	{binCode: 0xa6, length: 9, huffmanCode: 0x1a6},
	{binCode: 0xa7, length: 9, huffmanCode: 0x1a7},
	{binCode: 0xa8, length: 9, huffmanCode: 0x1a8},
	{binCode: 0xa9, length: 9, huffmanCode: 0x1a9},
	{binCode: 0xaa, length: 9, huffmanCode: 0x1aa},
	{binCode: 0xab, length: 9, huffmanCode: 0x1ab},
	{binCode: 0xac, length: 9, huffmanCode: 0x1ac},
	{binCode: 0xad, length: 9, huffmanCode: 0x1ad},
	{binCode: 0xae, length: 9, huffmanCode: 0x1ae},
	{binCode: 0xaf, length: 9, huffmanCode: 0x1af},
	{binCode: 0xb0, length: 9, huffmanCode: 0x1b0},
	{binCode: 0xb1, length: 9, huffmanCode: 0x1b1},
	{binCode: 0xb2, length: 9, huffmanCode: 0x1b2},
	{binCode: 0xb3, length: 9, huffmanCode: 0x1b3},
	{binCode: 0xb4, length: 9, huffmanCode: 0x1b4},
	{binCode: 0xb5, length: 9, huffmanCode: 0x1b5},
	{binCode: 0xb6, length: 9, huffmanCode: 0x1b6},
	{binCode: 0xb7, length: 9, huffmanCode: 0x1b7},
	{binCode: 0xb8, length: 9, huffmanCode: 0x1b8},
	{binCode: 0xb9, length: 9, huffmanCode: 0x1b9},
	{binCode: 0xba, length: 9, huffmanCode: 0x1ba},
	{binCode: 0xbb, length: 9, huffmanCode: 0x1bb},
	{binCode: 0xbc, length: 9, huffmanCode: 0x1bc},
	{binCode: 0xbd, length: 9, huffmanCode: 0x1bd},
	{binCode: 0xbe, length: 9, huffmanCode: 0x1be},
	{binCode: 0xbf, length: 9, huffmanCode: 0x1bf},
	{binCode: 0xc0, length: 9, huffmanCode: 0x1c0},
	{binCode: 0xc1, length: 9, huffmanCode: 0x1c1},
	{binCode: 0xc2, length: 9, huffmanCode: 0x1c2},
	{binCode: 0xc3, length: 9, huffmanCode: 0x1c3},
	{binCode: 0xc4, length: 9, huffmanCode: 0x1c4},
	{binCode: 0xc5, length: 9, huffmanCode: 0x1c5},
	{binCode: 0xc6, length: 9, huffmanCode: 0x1c6},
	{binCode: 0xc7, length: 9, huffmanCode: 0x1c7},
	{binCode: 0xc8, length: 9, huffmanCode: 0x1c8},
	{binCode: 0xc9, length: 9, huffmanCode: 0x1c9},
	{binCode: 0xca, length: 9, huffmanCode: 0x1ca},
	{binCode: 0xcb, length: 9, huffmanCode: 0x1cb},
	{binCode: 0xcc, length: 9, huffmanCode: 0x1cc},
	{binCode: 0xcd, length: 9, huffmanCode: 0x1cd},
	{binCode: 0xce, length: 9, huffmanCode: 0x1ce},
	{binCode: 0xcf, length: 9, huffmanCode: 0x1cf},
	{binCode: 0xd0, length: 9, huffmanCode: 0x1d0},
	{binCode: 0xd1, length: 9, huffmanCode: 0x1d1},
	{binCode: 0xd2, length: 9, huffmanCode: 0x1d2},
	{binCode: 0xd3, length: 9, huffmanCode: 0x1d3},
	{binCode: 0xd4, length: 9, huffmanCode: 0x1d4},
	{binCode: 0xd5, length: 9, huffmanCode: 0x1d5},
	{binCode: 0xd6, length: 9, huffmanCode: 0x1d6},
	{binCode: 0xd7, length: 9, huffmanCode: 0x1d7},
	{binCode: 0xd8, length: 9, huffmanCode: 0x1d8},
	{binCode: 0xd9, length: 9, huffmanCode: 0x1d9},
	{binCode: 0xda, length: 9, huffmanCode: 0x1da},
	{binCode: 0xdb, length: 9, huffmanCode: 0x1db},
	{binCode: 0xdc, length: 9, huffmanCode: 0x1dc},
	{binCode: 0xdd, length: 9, huffmanCode: 0x1dd},
	{binCode: 0xde, length: 9, huffmanCode: 0x1de},
	{binCode: 0xdf, length: 9, huffmanCode: 0x1df},
	{binCode: 0xe0, length: 9, huffmanCode: 0x1e0},
	{binCode: 0xe1, length: 9, huffmanCode: 0x1e1},
	{binCode: 0xe2, length: 9, huffmanCode: 0x1e2},
	{binCode: 0xe3, length: 9, huffmanCode: 0x1e3},
	{binCode: 0xe4, length: 9, huffmanCode: 0x1e4},
	{binCode: 0xe5, length: 9, huffmanCode: 0x1e5},
	{binCode: 0xe6, length: 9, huffmanCode: 0x1e6},
	{binCode: 0xe7, length: 9, huffmanCode: 0x1e7},
	{binCode: 0xe8, length: 9, huffmanCode: 0x1e8},
	{binCode: 0xe9, length: 9, huffmanCode: 0x1e9},
	{binCode: 0xea, length: 9, huffmanCode: 0x1ea},
	{binCode: 0xeb, length: 9, huffmanCode: 0x1eb},
	{binCode: 0xec, length: 9, huffmanCode: 0x1ec},
	{binCode: 0xed, length: 9, huffmanCode: 0x1ed},
	{binCode: 0xee, length: 9, huffmanCode: 0x1ee},
	{binCode: 0xef, length: 9, huffmanCode: 0x1ef},
	{binCode: 0xf0, length: 9, huffmanCode: 0x1f0},
	{binCode: 0xf1, length: 9, huffmanCode: 0x1f1},
	{binCode: 0xf2, length: 9, huffmanCode: 0x1f2},
	{binCode: 0xf3, length: 9, huffmanCode: 0x1f3},
	{binCode: 0xf4, length: 9, huffmanCode: 0x1f4},
	{binCode: 0xf5, length: 9, huffmanCode: 0x1f5},
	{binCode: 0xf6, length: 9, huffmanCode: 0x1f6},
	{binCode: 0xf7, length: 9, huffmanCode: 0x1f7},
	{binCode: 0xf8, length: 9, huffmanCode: 0x1f8},
	{binCode: 0xf9, length: 9, huffmanCode: 0x1f9},
	{binCode: 0xfa, length: 9, huffmanCode: 0x1fa},
	{binCode: 0xfb, length: 9, huffmanCode: 0x1fb},
	{binCode: 0xfc, length: 9, huffmanCode: 0x1fc},
	{binCode: 0xfd, length: 9, huffmanCode: 0x1fd},
	{binCode: 0xfe, length: 9, huffmanCode: 0x1fe},
	{binCode: 0xff, length: 9, huffmanCode: 0x1ff},
}
