package repository

import (
	"encoding/binary"
	"fmt"

	"github.com/miu200521358/win"
	"golang.org/x/sys/windows"
)

// parseCompressedBinaryXFile は、圧縮バイナリ形式の X ファイルを解凍し、
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
	prevBytesDecompressed := uint32(headerSize)
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

// decompressRFC1951Block は、圧縮データ compressed を元に展開先のバッファ decompressed に展開を行い、
// 展開したバイト数を返します。展開中、ブロックごとに最終ブロックフラグおよびブロック種別（非圧縮／固定ハフマン／動的ハフマン）に基づいて処理を行います。
func (rep *XRepository) decompressRFC1951Block(compressed []byte, decompressed []byte, prevBytesDecompressed uint32) (uint32, error) {
	var ctx decompressionContext
	rep.clearContext(&ctx)
	// 前回の展開バイト数を設定
	ctx.BytesDecompressed = uint32(prevBytesDecompressed)

	// 入力／出力バッファおよびサイズ情報の設定
	ctx.Src = compressed
	ctx.Dest = decompressed
	ctx.SizeInBytes = uint32(len(compressed))
	ctx.DestSize = uint32(len(decompressed))
	// 通常、バイト数 × 8 で全ビット数となる
	ctx.SizeInBits = int64(len(compressed)) * 8

	// 事前に初期化された固定ハフマンツリーを使用
	ctx.Fixed = globalFixedHuffmanTree

	final := false
	// 圧縮データ全体を処理するループ
	for !final && ctx.BitsRead < ctx.SizeInBits {
		// 最初の1ビットが最終ブロックフラグとなる
		wTmp := rep.getBits(&ctx, 1)
		if wTmp == 1 {
			final = true
		}
		// 次の2ビットでブロックの種類を取得
		blockType := rep.getBits(&ctx, 2)
		switch blockType {
		case 0:
			// 非圧縮ブロック
			if !rep.processUncompressed(&ctx) {
				return 0, fmt.Errorf("ProcessUncompressed に失敗")
			}
		case 1:
			// 固定ハフマン圧縮ブロック
			if !rep.processHuffmanFixed(&ctx) {
				return 0, fmt.Errorf("ProcessHuffmanFixed に失敗")
			}
		case 2:
			// 動的ハフマン圧縮ブロック
			if !rep.processHuffmanCustom(&ctx) {
				return 0, fmt.Errorf("ProcessHuffmanCustom に失敗")
			}
		default:
			return 0, fmt.Errorf("不明なブロック種別: %d", blockType)
		}
	}

	return ctx.BytesDecompressed, nil
}

// processHuffmanCustom は、動的ハフマン圧縮ブロックの復号処理を行います。
// まず、CreateHuffmanTreeTable によりカスタムハフマンツリー（リテラル／長さツリーおよび距離ツリー）を構築し、
// その後、DecodeWithHuffmanTree によりブロックを復号します。
// 最後に、DisposeCustomHuffmanTree で構築したツリーを解放します。
// 復号処理が正常に完了した場合は true を、エラーがあれば false を返します。
func (rep *XRepository) processHuffmanCustom(ctx *decompressionContext) bool {
	// カスタムハフマンツリーの構築
	if !rep.createHuffmanTreeTable(ctx) {
		return false
	}

	// ハフマンツリーを用いて復号処理
	if !rep.decodeWithHuffmanTree(ctx, &ctx.CustomLiteralLength, &ctx.CustomDistance) {
		rep.disposeCustomHuffmanTree(ctx)
		return false
	}

	// 後片付け
	rep.disposeCustomHuffmanTree(ctx)
	return true
}

// createHuffmanTreeTable は、動的ハフマン圧縮ブロックにおける
// リテラル／長さツリーおよび距離ツリーを構築します。
// 正常に構築できた場合は true を、エラーがあれば false を返します。
func (rep *XRepository) createHuffmanTreeTable(ctx *decompressionContext) bool {
	// 1. パラメータ読み出し
	numLiteralLengthCodes := uint16(rep.getBits(ctx, 5)) + 257
	numDistanceCodes := uint16(rep.getBits(ctx, 5)) + 1
	numCodeLengthCode := uint16(rep.getBits(ctx, 4)) + 4

	// 2. コード長のコード値の配列を初期化（長さ20）
	var codeLengthsOfTheCodeLength [20]byte
	for i := range codeLengthsOfTheCodeLength {
		codeLengthsOfTheCodeLength[i] = 0
	}
	itsOrder := []int{16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15}
	// numCodeLengthCode 個分のコード長を、その順序に従って取得
	for i := 0; i < int(numCodeLengthCode); i++ {
		codeLengthsOfTheCodeLength[itsOrder[i]] = byte(rep.getBits(ctx, 3))
	}

	// 3. CodeLengthsTree の作成
	var codeLengthsTree huffmanTree
	bitLengths := make([]int, 20)
	nextCodes := make([]int, 20)
	node := make([]huffmanCode, 19)
	maxBits := -1
	// node 配列（サイズ19）に、各コード長を設定し、bitLengths 配列でカウント
	for i := 0; i < 19; i++ {
		lenVal := int(codeLengthsOfTheCodeLength[i])
		node[i].Length = uint16(lenVal)
		bitLengths[lenVal]++
		if lenVal > maxBits {
			maxBits = lenVal
		}
	}
	// 次のコードの値を計算
	code := 0
	bitLengths[0] = 0
	for bits := 1; bits <= 18; bits++ {
		code = (code + bitLengths[bits-1]) << 1
		nextCodes[bits] = code
	}
	// node 配列にハフマンコードを割り当てる
	for n := 0; n < 19; n++ {
		lenVal := int(node[n].Length)
		node[n].BinCode = uint16(n)
		if lenVal != 0 {
			node[n].HuffmanCode = uint16(nextCodes[lenVal])
			nextCodes[lenVal]++
		}
	}
	// CodeLengthsTree の構築
	if !rep.setupHuffmanTree(&codeLengthsTree, node) {
		return false
	}

	// 4. リテラル／長さツリーの作成
	literalLengthTree := make([]huffmanCode, numLiteralLengthCodes)
	var prev uint16 = 0xffff
	var repeat uint16 = 0
	n := 0
	indvCount := 0
	repeatCount := 0
	maxBits = -1
	for n < int(numLiteralLengthCodes) {
		literalLengthTree[n].BinCode = uint16(n)
		if repeat > 0 {
			literalLengthTree[n].Length = prev
			n++
			repeatCount++
			repeat--
		} else {
			// GetACodeWithHuffmanTree を用いてコード長を取得
			lenVal := int(rep.getACodeWithHuffmanTree(ctx, &codeLengthsTree))
			switch lenVal {
			case 16: // 前回のコード長を 3～6 回コピー
				repeat = uint16(rep.getBits(ctx, 2)) + 3
			case 17: // 0 を 3～10 回繰り返す
				prev = 0
				repeat = uint16(rep.getBits(ctx, 3)) + 3
			case 18: // 0 を 11～138 回繰り返す
				prev = 0
				repeat = uint16(rep.getBits(ctx, 7)) + 11
			default:
				if repeat > 0 {
					return false
				}
				prev = uint16(lenVal)
				repeat = 0
				literalLengthTree[n].Length = uint16(lenVal)
				n++
				indvCount++
				if lenVal > maxBits {
					maxBits = lenVal
				}
			}
		}
	}
	if repeat > 0 {
		return false
	}
	// 次のコード値の割り当てのために再度 bitLengths, nextCodes を構築
	bitLengths = make([]int, maxBits+1)
	nextCodes = make([]int, maxBits+1)
	for i := 0; i < int(numLiteralLengthCodes); i++ {
		bitLengths[int(literalLengthTree[i].Length)]++
	}
	bitLengths[0] = 0
	code = 0
	for bits := 1; bits <= maxBits; bits++ {
		code = (code + bitLengths[bits-1]) << 1
		nextCodes[bits] = code
	}
	for n = 0; n < int(numLiteralLengthCodes); n++ {
		l := int(literalLengthTree[n].Length)
		if l != 0 {
			literalLengthTree[n].HuffmanCode = uint16(nextCodes[l])
			nextCodes[l]++
		}
	}
	bSuccess := rep.setupHuffmanTree(&ctx.CustomLiteralLength, literalLengthTree)
	if !bSuccess {
		return false
	}

	// 5. 距離ツリーの作成
	if bSuccess {
		distanceTree := make([]huffmanCode, numDistanceCodes)
		maxBits = -1
		prev = 0xffff
		repeat = 0
		n = 0
		for n < int(numDistanceCodes) {
			distanceTree[n].BinCode = uint16(n)
			if repeat > 0 {
				distanceTree[n].Length = prev
				n++
				repeat--
			} else {
				lenVal := int(rep.getACodeWithHuffmanTree(ctx, &codeLengthsTree))
				switch lenVal {
				case 16:
					repeat = uint16(rep.getBits(ctx, 2)) + 3
				case 17:
					prev = 0
					repeat = uint16(rep.getBits(ctx, 3)) + 3
				case 18:
					prev = 0
					repeat = uint16(rep.getBits(ctx, 7)) + 11
				default:
					if repeat > 0 {
						return false
					}
					prev = uint16(lenVal)
					repeat = 0
					distanceTree[n].Length = uint16(lenVal)
					n++
					if lenVal > maxBits {
						maxBits = lenVal
					}
				}
			}
		}
		if repeat > 0 {
			return false
		}
		bitLengths = make([]int, maxBits+1)
		nextCodes = make([]int, maxBits+1)
		for i := 0; i < int(numDistanceCodes); i++ {
			bitLengths[int(distanceTree[i].Length)]++
		}
		bitLengths[0] = 0
		code = 0
		for bits := 1; bits <= maxBits; bits++ {
			code = (code + bitLengths[bits-1]) << 1
			nextCodes[bits] = code
		}
		for n = 0; n < int(numDistanceCodes); n++ {
			l := int(distanceTree[n].Length)
			if l != 0 {
				distanceTree[n].HuffmanCode = uint16(nextCodes[l])
				nextCodes[l]++
			}
		}
		bSuccess = rep.setupHuffmanTree(&ctx.CustomDistance, distanceTree)
		if !bSuccess {
		}
	}

	// CodeLengthsTree のリソース解放（Go では不要なので、参照をクリア）
	codeLengthsTree.Codes = nil

	if !bSuccess {
		rep.disposeCustomHuffmanTree(ctx)
	}
	return bSuccess
}

// disposeCustomHuffmanTree は、DecompressionContext 内のカスタムハフマンツリー
// （CustomDistance と CustomLiteralLength）のリソースを解放（nil に設定）し、
// フィールドの値をリセットします。
func (rep *XRepository) disposeCustomHuffmanTree(ctx *decompressionContext) {
	// CustomDistance のリソース解放
	ctx.CustomDistance.Codes = nil
	ctx.CustomDistance.NumMaxBits = 0
	ctx.CustomDistance.NumCodes = 0

	// CustomLiteralLength のリソース解放
	ctx.CustomLiteralLength.Codes = nil
	ctx.CustomLiteralLength.NumMaxBits = 0
	ctx.CustomLiteralLength.NumCodes = 0
}

// processHuffmanFixed は、固定ハフマンブロックの復号処理を行います。
// 固定ハフマンツリー (ctx.Fixed) を用いて、DecodeWithHuffmanTree を呼び出し、
// 復号処理の成否を返します。
func (rep *XRepository) processHuffmanFixed(ctx *decompressionContext) bool {
	// 固定ハフマンツリーはすでに設定されているため、初期化不要
	result := rep.decodeWithHuffmanTree(ctx, &ctx.Fixed, nil)
	return result
}

// processUncompressed は、非圧縮ブロックのデータを処理し、
// チェックサム検証後、サイズ分のバイトを出力バッファへコピーします。
// 正常終了の場合は true を、エラー時は false を返します。
func (rep *XRepository) processUncompressed(ctx *decompressionContext) bool {
	// サイズ（16bit）を読み出す
	sizeLow := rep.getBits(ctx, 8) & 0xff
	sizeHigh := rep.getBits(ctx, 8) & 0xff
	size := uint16(sizeLow) | (uint16(sizeHigh) << 8)

	// チェックサム用の値（16bit）を読み出し、反転してサイズと一致するか検証
	tmpLow := rep.getBits(ctx, 8) & 0xff
	tmpHigh := rep.getBits(ctx, 8) & 0xff
	wTmp := uint16(tmpLow) | (uint16(tmpHigh) << 8)

	// ビット反転（Go では ^ 演算子でビットごとのNOT）
	wTmp = ^wTmp

	if wTmp != size {
		return false
	}

	// size バイト分のデータを出力先バッファにコピーする
	for i := 0; i < int(size); i++ {
		b := rep.getBits(ctx, 8) & 0xff
		if int(ctx.BytesDecompressed) >= len(ctx.Dest) {
			return false
		}
		ctx.Dest[ctx.BytesDecompressed] = byte(b)
		ctx.BytesDecompressed++
	}

	return true
}

// decodeWithHuffmanTree は、literalLengthTree および distanceTree を用いて
// 圧縮データからリテラルおよび長さ・距離ペアを復号し、ctx.Dest に展開します。
// 終端コード（256）が現れた場合に処理を終了し、正常終了なら true を、エラー時は false を返します。
func (rep *XRepository) decodeWithHuffmanTree(ctx *decompressionContext, literalLengthTree *huffmanTree, distanceTree *huffmanTree) bool {
	for {
		// リテラル/長さコードの取得
		w := rep.getACodeWithHuffmanTree(ctx, literalLengthTree)

		if w == 0xffff {
			return false
		}
		// 終端コード 256 の場合、ブロック終了
		if w == 256 {
			break
		}
		if w < 256 {
			// リテラルバイトの場合、そのまま出力
			if int(ctx.BytesDecompressed) >= len(ctx.Dest) {
				return false
			}
			ctx.Dest[ctx.BytesDecompressed] = byte(w)
			ctx.BytesDecompressed++
		} else {
			// 長さ/距離ペアの場合
			length := rep.decodeLength(ctx, w)
			var distCode uint16
			if distanceTree != nil {
				distCode = rep.getACodeWithHuffmanTree(ctx, distanceTree)
			} else {
				distCode = rep.getBits(ctx, 5)
				distCode = rep.binreverse(int(distCode), 5)
			}

			distance := rep.decodeDistance(ctx, distCode)

			// コピー元インデックスは、現在の出力位置から distance 分戻った位置
			start := int(ctx.BytesDecompressed) - distance

			if start < 0 {
				return false
			}
			// length バイト分、過去の出力からコピーする
			for i := range length {
				if int(ctx.BytesDecompressed) >= len(ctx.Dest) || start+i >= len(ctx.Dest) {
					return false
				}
				ctx.Dest[ctx.BytesDecompressed] = ctx.Dest[start+i]
				ctx.BytesDecompressed++
			}
		}
	}
	return true
}

// decodeLength は、baseCode（リテラル/長さコード）から実際の長さを復号します。
// C++実装:
//
//	if (baseCode <= 264)
//	   return baseCode - 257 + 3;
//	if (baseCode == 285)
//	   return 258;
//	if (baseCode > 285)
//	   return 0xffff;  // error
//	w = baseCode - 265;
//	x = w >> 2;
//	numExtraBits = x + 1;
//	y = (4 << numExtraBits) + 3;
//	y += (w & 3) << numExtraBits;
//	extra = GetBits(ctx, numExtraBits);
//	return y + extra;
func (rep *XRepository) decodeLength(ctx *decompressionContext, baseCode uint16) int {
	if baseCode <= 264 {
		val := int(baseCode) - 257 + 3
		return val
	}
	if baseCode == 285 {
		return 258
	}
	if baseCode > 285 {
		return 0xffff
	}
	// baseCode が 265 ～ 284 の場合
	w := baseCode - 265
	x := w >> 2
	numExtraBits := x + 1
	// y = (4 << numExtraBits) + 3
	y := uint16((4 << int(numExtraBits)) + 3)
	// y += (w & 3) << numExtraBits;
	y += uint16((int(w & 3)) << int(numExtraBits))
	extra := rep.getBits(ctx, int(numExtraBits))
	result := int(y) + int(extra)
	return result
}

// decodeDistance は、baseCode（距離コード）から実際の距離を復号します。
// C++実装:
//
//	if (baseCode <= 3)
//	   return baseCode + 1;
//	if (baseCode > 29)
//	   return 0;  // error
//	w = baseCode - 4;
//	x = w >> 1;
//	numExtraBits = x + 1;
//	y = (2 << numExtraBits) + 1;
//	y += (w & 1) << numExtraBits;
//	extra = GetBits(ctx, numExtraBits);
//	return y + extra;
func (rep *XRepository) decodeDistance(ctx *decompressionContext, baseCode uint16) int {
	if baseCode <= 3 {
		val := int(baseCode) + 1
		return val
	}
	if baseCode > 29 {
		return 0
	}
	w := baseCode - 4
	x := w >> 1
	numExtraBits := x + 1
	y := uint16((2 << int(numExtraBits)) + 1)
	y += uint16((int(w & 1)) << int(numExtraBits))
	extra := rep.getBits(ctx, int(numExtraBits))
	result := int(y) + int(extra)
	return result
}

// compareCode は、HuffmanCode a のビット長とコード値と、指定された length および code を比較します。
// a.Length が length より小さい場合は -1、大きい場合は 1、等しい場合は
// a.HuffmanCode と code の大小で比較し、等しければ 0 を返します。
func (rep *XRepository) compareCode(a huffmanCode, length int, code uint16) int {
	if int(a.Length) < length {
		return -1
	}
	if int(a.Length) > length {
		return 1
	}
	if a.HuffmanCode < code {
		return -1
	}
	if a.HuffmanCode == code {
		return 0
	}
	return 1
}

// getACodeWithHuffmanTree は、ctx.Src からビットを順次読み出し、
// tree 内のハフマンコードと照合して該当するシンボル（binCode）を返します。
// 照合処理は、入力の最小ビット長から始まり、ビット数を増加させながらバイナリサーチで行います。
// 該当が見つからなかった場合は 0xffff を返します。
func (rep *XRepository) getACodeWithHuffmanTree(ctx *decompressionContext, tree *huffmanTree) uint16 {
	if tree.Codes == nil || len(tree.Codes) == 0 {
		return 0xffff
	}
	maxBits := int(tree.NumMaxBits)
	var w uint16 = 0
	// 初期ビット数は、ツリーの先頭コードの長さとする
	minBits := int(tree.Codes[0].Length)
	// 初期の minBits 分のビットを取得
	for bits := 0; bits < minBits; bits++ {
		w <<= 1
		w |= rep.getABit(ctx)
	}

	// 現在のビット数 (minBits) から最大ビット長まで拡張しながら検索
	bits := minBits
	for bits <= maxBits {
		left := 0
		right := int(tree.NumCodes)
		// バイナリサーチで該当コードを探す
		for left < right {
			mid := (left + right) / 2
			comp := rep.compareCode(tree.Codes[mid], bits, w)
			if comp == 0 {
				return tree.Codes[mid].BinCode
			} else if comp < 0 {
				left = mid + 1
			} else {
				right = mid
			}
		}
		// 該当が見つからなかったので、1ビット追加して再検索
		w <<= 1
		addedBit := rep.getABit(ctx)
		w |= addedBit
		bits++
	}
	return 0xffff
}

// setupHuffmanTree は、入力の HuffmanCode 配列 codes から
// 有効なコードを抽出し、挿入ソートにより順序付けた結果を
// tree の Codes にセットし、NumCodes および NumMaxBits を設定します。
func (rep *XRepository) setupHuffmanTree(tree *huffmanTree, codes []huffmanCode) bool {
	// 一時的なスライスを確保（最大長は入力数と同じ）
	dest := make([]huffmanCode, len(codes))
	numCodes := 0
	maxBits := 0

	// 入力の各コードについて
	for _, codeEntry := range codes {
		length := int(codeEntry.Length)
		if length != 0 {
			j := numCodes
			codeVal := int(codeEntry.HuffmanCode)
			// 挿入位置を決定（降順にソート：長さが短いものが後ろ、同じ長さなら HuffmanCode が小さい順）
			for j > 0 {
				prev := dest[j-1]
				if int(prev.Length) < length {
					break
				}
				if int(prev.Length) == length && int(prev.HuffmanCode) <= codeVal {
					break
				}
				// シフトして挿入スペースを確保
				dest[j] = dest[j-1]
				j--
			}
			// 挿入
			dest[j].Length = uint16(length)
			dest[j].HuffmanCode = uint16(codeVal)
			dest[j].BinCode = codeEntry.BinCode
			numCodes++
			if length > maxBits {
				maxBits = length
			}
		}
	}
	tree.NumMaxBits = uint16(maxBits)
	tree.NumCodes = uint16(numCodes)
	tree.Codes = dest[:numCodes]
	return true
}

// setupFixedHuffmanTree は、固定ハフマンテーブル fixedCodeSorted を用いて
// DecompressionContext の Fixed フィールドを初期化します。
func (rep *XRepository) setupFixedHuffmanTree(ctx *decompressionContext) {
	// 事前に初期化された固定ハフマンツリーを使用
	ctx.Fixed = globalFixedHuffmanTree
}

// clearContext は、DecompressionContext のフィールドを初期状態にリセットします。
func (rep *XRepository) clearContext(ctx *decompressionContext) {
	// ソース／デスティネーション関連のフィールドをリセット
	ctx.Src = nil
	ctx.BitsRead = 0
	ctx.SizeInBits = 0
	ctx.SizeInBytes = 0
	ctx.Dest = nil
	ctx.DestSize = 0
	ctx.BytesDecompressed = 0

	// カスタム距離ツリーの初期化
	ctx.CustomDistance.NumMaxBits = 0
	ctx.CustomDistance.Codes = nil
	ctx.CustomDistance.NumCodes = 0

	// カスタムリテラル/長さツリーの初期化
	ctx.CustomLiteralLength.NumMaxBits = 0
	ctx.CustomLiteralLength.Codes = nil
	ctx.CustomLiteralLength.NumCodes = 0

	// 固定ハフマンツリーの初期化
	ctx.Fixed.NumMaxBits = 0
	ctx.Fixed.Codes = nil
	ctx.Fixed.NumCodes = 0
}

// binreverse は、引数 code の下位 length ビットを反転した結果を返します。
// 例: code=0b1011, length=4 の場合、反転結果は 0b1101 となります。
func (rep *XRepository) binreverse(code int, length int) uint16 {
	var w uint16 = 0
	for i := 0; i < length; i++ {
		w <<= 1
		currentBit := code & 1
		w |= uint16(currentBit)
		code >>= 1
	}
	return w
}

// getABit は、DecompressionContext のSrcから現在のビット位置の1ビットを読み出します。
func (rep *XRepository) getABit(ctx *decompressionContext) uint16 {
	// 現在のビット位置からバイト配列上のインデックスとビット位置を算出
	index := ctx.BitsRead >> 3
	if int(index) >= len(ctx.Src) {
		return 0
	}
	byteVal := ctx.Src[index]
	bitIndex := ctx.BitsRead & 7
	// 対象のビットを取得
	bit := (byteVal >> bitIndex) & 1
	ctx.BitsRead++ // 読み出し後にカウンタをインクリメント
	return uint16(bit)
}

// getBits は、DecompressionContext のSrcから指定されたsizeビット分を読み出します。
func (rep *XRepository) getBits(ctx *decompressionContext, size int) uint16 {
	if size <= 0 {
		return 0
	}

	// 一度に処理するためビット位置をサイズ分進める
	ctx.BitsRead += int64(size)
	current := ctx.BitsRead
	var result uint16 = 0

	// 最適化: サイズに応じて一度に複数ビットを処理
	for i := 0; i < size; i++ {
		result <<= 1
		current-- // 読み出すビット位置を後ろにずらす
		index := current >> 3
		if int(index) >= len(ctx.Src) {
			return result
		}
		byteVal := ctx.Src[index]
		bitIndex := current & 7
		bit := (byteVal >> bitIndex) & 1
		result |= uint16(bit)
	}
	return result
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
	BinCode     uint16
	Length      uint16
	HuffmanCode uint16
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
	NumMaxBits uint16
	Codes      []huffmanCode
	NumCodes   uint16
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
	NumLiteralLengthCodes          uint16
	NumDistanceCodes               uint16
	LiteralLengthTree              []uint16
	NumElementsOfLiteralLengthTree uint16
	DistanceTree                   []uint16
	NumElementsOfDistanceTree      uint16
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
	Src                 []byte
	BitsRead            int64
	SizeInBits          int64
	SizeInBytes         uint32
	Dest                []byte
	DestSize            uint32
	BytesDecompressed   uint32
	CustomLiteralLength huffmanTree
	CustomDistance      huffmanTree
	Fixed               huffmanTree
}

var fixedCodeSorted []huffmanCode
var globalFixedHuffmanTree huffmanTree

func init() {
	// 初期化時にfixedCodeSorted配列を構築する

	// グループ1: BinCode 0x100-0x117, Length 7, HuffmanCode 0x0-0x17
	for i := uint16(0x100); i <= uint16(0x117); i++ {
		fixedCodeSorted = append(fixedCodeSorted, huffmanCode{
			BinCode:     i,
			Length:      uint16(7),
			HuffmanCode: i - uint16(0x100),
		})
	}

	// グループ2: BinCode 0x0-0x8f, Length 8, HuffmanCode 0x30-0xbf
	for i := uint16(0x0); i <= uint16(0x8f); i++ {
		fixedCodeSorted = append(fixedCodeSorted, huffmanCode{
			BinCode:     i,
			Length:      uint16(8),
			HuffmanCode: i + uint16(0x30),
		})
	}

	// グループ3: BinCode 0x118-0x11f, Length 8, HuffmanCode 0xc0-0xc7
	for i := uint16(0x118); i <= uint16(0x11f); i++ {
		fixedCodeSorted = append(fixedCodeSorted, huffmanCode{
			BinCode:     i,
			Length:      uint16(8),
			HuffmanCode: i - uint16(0x118) + uint16(0xc0),
		})
	}

	// グループ4: BinCode 0x90-0xff, Length 9, HuffmanCode 0x190-0x1ff
	for i := uint16(0x90); i <= uint16(0xff); i++ {
		fixedCodeSorted = append(fixedCodeSorted, huffmanCode{
			BinCode:     i,
			Length:      uint16(9),
			HuffmanCode: i + uint16(0x100),
		})
	}

	// グローバルな固定ハフマンツリーの初期化
	globalFixedHuffmanTree.NumCodes = 288
	globalFixedHuffmanTree.NumMaxBits = 9
	globalFixedHuffmanTree.Codes = fixedCodeSorted
}
