package repository

import (
	"encoding/binary"
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/win"
	"golang.org/x/sys/windows"
)

// parseCompressedBinaryXFile は、圧縮バイナリ形式の X ファイルを解凍し、
// 解凍後のデータをバイナリパーサーへ渡します。
func (rep *XRepository) parseCompressedBinaryXFile(model *pmx.PmxModel) error {
	var err error
	var fileSize uint64

	var xHandle win.HWND
	if xHandle, err = win.CreateFile(
		rep.path,
		win.GENERIC_READ,
		win.FILE_SHARE_READ,
		0,
		win.OPEN_EXISTING,
		win.FILE_ATTRIBUTE_NORMAL,
		0,
	); err != windows.DS_S_SUCCESS {
		return fmt.Errorf("CreateFile failed: %v", err)
	}

	if ret, err := win.GetFileSizeEx(xHandle, &fileSize); !ret {
		return fmt.Errorf("GetFileSizeEx failed: %v", err)
	}
	inputFileSize := uint32(fileSize)

	compressedBuffer := make([]byte, inputFileSize)
	var compressedByteRead uint32

	if _, err := win.ReadFile(
		xHandle,
		&compressedBuffer[0],
		inputFileSize,
		&compressedByteRead,
		0,
	); err != windows.DS_S_SUCCESS {
		return fmt.Errorf("ReadFile failed: %v", err)
	}

	// MSZIP 圧縮アルゴリズムでデコンプレッサを作成
	var decompressorHandle win.HWND
	if _, err := win.CreateDecompressor(win.COMPRESS_ALGORITHM_MSZIP|win.COMPRESS_RAW, &decompressorHandle); err != windows.DS_S_SUCCESS {
		return fmt.Errorf("CreateDecompressor failed: %v", err)
	}

	// MSZIP 圧縮データを解凍
	if err := decompressMSZipXFile(compressedBuffer, inputFileSize, decompressorHandle); err != nil {
		return fmt.Errorf("decompressMSZipXFile failed: %v", err)
	}

	// ハンドルを解放
	win.CloseDecompressor(decompressorHandle)

	return nil
}

// decompressMSZipXFile は、MSZIP形式の圧縮ブロックを複数ループで解凍し、
// 詳細なデバッグ情報を出力するサンプルです。
func decompressMSZipXFile(
	compressedBuffer []byte,
	inputFileSize uint32,
	decompressorHandle win.HWND,
) error {

	// 1) 先頭16バイトヘッダー + 4バイトのサイズフィールド
	const headerSize = 16
	const sizeFieldSize = 4
	totalHeaderSize := uint32(headerSize + sizeFieldSize) // 20バイト

	// ファイルが小さすぎないかチェック
	if inputFileSize < totalHeaderSize {
		return fmt.Errorf("file is too small or corrupted: inputFileSize=%d < totalHeaderSize=%d", inputFileSize, totalHeaderSize)
	}

	// 伸長後のファイルサイズをリトルエンディアンで読み取る
	finalSize := binary.LittleEndian.Uint32(compressedBuffer[headerSize : headerSize+sizeFieldSize])

	// 出力用バッファを確保
	if finalSize < headerSize {
		// 伸長後サイズがヘッダー以下はあり得ないのでチェック
		return fmt.Errorf("invalid finalSize=%d (< %d)", finalSize, headerSize)
	}
	newBuffer := make([]byte, finalSize)

	// 先頭16バイトをコピー
	copy(newBuffer[0:headerSize], compressedBuffer[0:headerSize])

	// 出力先のオフセットを16バイト進める
	destOffset := headerSize
	uncompressedSum := uint32(headerSize)

	// 圧縮データの開始位置 (20バイト目)
	srcOffset := uint32(totalHeaderSize)
	compressedBufferLen := uint32(len(compressedBuffer))
	if srcOffset > compressedBufferLen {
		return fmt.Errorf("file is too small for compressed data: srcOffset=%d > len=%d", srcOffset, compressedBufferLen)
	}

	// デバッグ出力: 先頭数バイトの圧縮データを表示
	maxBytesToShow := 16
	showBytes := int(compressedBufferLen - srcOffset)
	if showBytes > maxBytesToShow {
		showBytes = maxBytesToShow
	}

	// ブロック単位のループ
	blockIndex := 0
	prevBytesDecompressed := uint32(0)
	for srcOffset < compressedBufferLen && uncompressedSum < finalSize {
		blockIndex++
		if srcOffset+2 > compressedBufferLen {
			return fmt.Errorf("insufficient data for block uncompressed size (blockIndex=%d)", blockIndex)
		}
		blockUncompressedSize := binary.LittleEndian.Uint16(compressedBuffer[srcOffset : srcOffset+2])
		srcOffset += 2

		if srcOffset+2 > compressedBufferLen {
			return fmt.Errorf("insufficient data for block compressed size (blockIndex=%d)", blockIndex)
		}
		blockCompressedSize := binary.LittleEndian.Uint16(compressedBuffer[srcOffset : srcOffset+2])
		srcOffset += 2

		// ブロックデータが十分あるかチェック
		if srcOffset+uint32(blockCompressedSize) > compressedBufferLen {
			return fmt.Errorf("insufficient data for compressed block (blockIndex=%d) need=%d, remain=%d", blockIndex, blockCompressedSize, compressedBufferLen-srcOffset)
		}
		blockData := compressedBuffer[srcOffset : srcOffset+uint32(blockCompressedSize)]

		// デバッグ: ブロックの先頭数バイトを表示
		showBytes = int(blockCompressedSize)
		if showBytes > maxBytesToShow {
			showBytes = maxBytesToShow
		}

		// Decompress呼び出し
		// ここでエラーが出る場合は、フラグやブロックの開始位置が誤っているか、データが破損している可能性が高い
		success, err := win.Decompress(
			decompressorHandle,
			&blockData[0],
			uint32(blockCompressedSize),
			&newBuffer[destOffset],
			uint32(blockUncompressedSize),
			nil, // uncompressedDataSizeは使用しない
		)
		if !success {
			return fmt.Errorf("Decompress failed on block #%d: blockUncompressedSize=%d, blockCompressedSize=%d, err=%v",
				blockIndex, blockUncompressedSize, blockCompressedSize, err)
		}

		if BytesDecompressed, err := DecompressRFC1951Block(blockData[2:], newBuffer, prevBytesDecompressed); err != nil {
			return fmt.Errorf("DecompressRFC1951Block failed on block #%d: blockUncompressedSize=%d, blockCompressedSize=%d, ret=%d, err=%v",
				blockIndex, blockUncompressedSize, blockCompressedSize, BytesDecompressed, err)
		} else {
			prevBytesDecompressed = BytesDecompressed
		}

		uncompressedSum += uint32(blockUncompressedSize)
		destOffset += int(blockUncompressedSize)
		srcOffset += uint32(blockCompressedSize)
	}

	return nil
}

// DecompressRFC1951Block は、圧縮データ compressed を元に展開先のバッファ decompressed に展開を行い、
// 展開したバイト数を返します。展開中、ブロックごとに最終ブロックフラグおよびブロック種別（非圧縮／固定ハフマン／動的ハフマン）に基づいて処理を行います。
func DecompressRFC1951Block(compressed []byte, decompressed []byte, prevBytesDecompressed uint32) (uint32, error) {
	// mlog.D("DecompressRFC1951Block: 開始")

	var ctx DecompressionContext
	ClearContext(&ctx)
	// 前回の展開バイト数を設定
	ctx.BytesDecompressed = uint32(prevBytesDecompressed)

	// 入力／出力バッファおよびサイズ情報の設定
	ctx.Src = compressed
	ctx.Dest = decompressed
	ctx.SizeInBytes = uint32(len(compressed))
	ctx.DestSize = uint32(len(decompressed))
	// 通常、バイト数 × 8 で全ビット数となる
	ctx.SizeInBits = int64(len(compressed)) * 8
	// mlog.D("DecompressRFC1951Block: SizeInBytes=%d, SizeInBits=%d, DestSize=%d",
	// ctx.SizeInBytes, ctx.SizeInBits, ctx.DestSize)
	// mlog.D("DecompressRFC1951Block: 入力データの先頭数バイト: %v", compressed[:32])

	// 固定ハフマンツリーのセットアップ
	SetupFixedHuffmanTree(&ctx)

	final := false
	// 圧縮データ全体を処理するループ
	for !final && ctx.BitsRead < ctx.SizeInBits {
		// mlog.D("DecompressRFC1951Block: ブロック開始前の BytesDecompressed = %d", ctx.BytesDecompressed)

		// 最初の1ビットが最終ブロックフラグとなる
		wTmp := GetBits(&ctx, 1)
		// mlog.D("DecompressRFC1951Block: 最終ブロックフラグ = %d", wTmp)
		if wTmp == 1 {
			final = true
		}
		// 次の2ビットでブロックの種類を取得
		blockType := GetBits(&ctx, 2)
		// mlog.D("DecompressRFC1951Block: ブロック種別 = %d", blockType)
		switch blockType {
		case 0:
			// 非圧縮ブロック
			if !ProcessUncompressed(&ctx) {
				return 0, fmt.Errorf("ProcessUncompressed に失敗")
			}
		case 1:
			// 固定ハフマン圧縮ブロック
			if !ProcessHuffmanFixed(&ctx) {
				return 0, fmt.Errorf("ProcessHuffmanFixed に失敗")
			}
		case 2:
			// 動的ハフマン圧縮ブロック
			if !ProcessHuffmanCustom(&ctx) {
				return 0, fmt.Errorf("ProcessHuffmanCustom に失敗")
			}
		default:
			// mlog.D("DecompressRFC1951Block: 不明なブロック種別 %d", blockType)
			return 0, fmt.Errorf("不明なブロック種別: %d", blockType)
		}
	}

	// mlog.D("DecompressRFC1951Block: 完了 - BytesDecompressed=%d", ctx.BytesDecompressed)
	return ctx.BytesDecompressed, nil
}

// ProcessHuffmanCustom は、動的ハフマン圧縮ブロックの復号処理を行います。
// まず、CreateHuffmanTreeTable によりカスタムハフマンツリー（リテラル／長さツリーおよび距離ツリー）を構築し、
// その後、DecodeWithHuffmanTree によりブロックを復号します。
// 最後に、DisposeCustomHuffmanTree で構築したツリーを解放します。
// 復号処理が正常に完了した場合は true を、エラーがあれば false を返します。
func ProcessHuffmanCustom(ctx *DecompressionContext) bool {
	// mlog.D("ProcessHuffmanCustom: 開始")

	// カスタムハフマンツリーの構築
	if !CreateHuffmanTreeTable(ctx) {
		// mlog.D("ProcessHuffmanCustom: CreateHuffmanTreeTable に失敗")
		return false
	}

	// ハフマンツリーを用いて復号処理
	if !DecodeWithHuffmanTree(ctx, &ctx.CustomLiteralLength, &ctx.CustomDistance) {
		// mlog.D("ProcessHuffmanCustom: DecodeWithHuffmanTree に失敗")
		DisposeCustomHuffmanTree(ctx)
		return false
	}

	// 後片付け
	DisposeCustomHuffmanTree(ctx)
	// mlog.D("ProcessHuffmanCustom: 終了")
	return true
}

// CreateHuffmanTreeTable は、動的ハフマン圧縮ブロックにおける
// リテラル／長さツリーおよび距離ツリーを構築します。
// 正常に構築できた場合は true を、エラーがあれば false を返します。
func CreateHuffmanTreeTable(ctx *DecompressionContext) bool {
	// mlog.D("CreateHuffmanTreeTable: 開始")

	// 1. パラメータ読み出し
	numLiteralLengthCodes := uint16(GetBits(ctx, 5)) + 257
	numDistanceCodes := uint16(GetBits(ctx, 5)) + 1
	numCodeLengthCode := uint16(GetBits(ctx, 4)) + 4
	// mlog.D("CreateHuffmanTreeTable: numLiteralLengthCodes=%d, numDistanceCodes=%d, numCodeLengthCode=%d",
	// numLiteralLengthCodes, numDistanceCodes, numCodeLengthCode)

	// 2. コード長のコード値の配列を初期化（長さ20）
	var codeLengthsOfTheCodeLength [20]byte
	for i := range codeLengthsOfTheCodeLength {
		codeLengthsOfTheCodeLength[i] = 0
	}
	itsOrder := []int{16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15}
	// numCodeLengthCode 個分のコード長を、その順序に従って取得
	for i := 0; i < int(numCodeLengthCode); i++ {
		codeLengthsOfTheCodeLength[itsOrder[i]] = byte(GetBits(ctx, 3))
		// mlog.D("CreateHuffmanTreeTable: codeLengthsOfTheCodeLength[%d] = %d", itsOrder[i], codeLengthsOfTheCodeLength[itsOrder[i]])
	}

	// 3. CodeLengthsTree の作成
	var codeLengthsTree HuffmanTree
	bitLengths := make([]int, 20)
	nextCodes := make([]int, 20)
	node := make([]HuffmanCode, 19)
	maxBits := -1
	// node 配列（サイズ19）に、各コード長を設定し、bitLengths 配列でカウント
	for i := 0; i < 19; i++ {
		lenVal := int(codeLengthsOfTheCodeLength[i])
		node[i].Length = uint16(lenVal)
		bitLengths[lenVal]++
		if lenVal > maxBits {
			maxBits = lenVal
		}
		// mlog.D("CreateHuffmanTreeTable: node[%d].Length = %d", i, lenVal)
	}
	// 次のコードの値を計算
	code := 0
	bitLengths[0] = 0
	for bits := 1; bits <= 18; bits++ {
		code = (code + bitLengths[bits-1]) << 1
		nextCodes[bits] = code
		// mlog.D("CreateHuffmanTreeTable: nextCodes[%d] = %d", bits, code)
	}
	// node 配列にハフマンコードを割り当てる
	for n := 0; n < 19; n++ {
		lenVal := int(node[n].Length)
		node[n].BinCode = uint16(n)
		if lenVal != 0 {
			node[n].HuffmanCode = uint16(nextCodes[lenVal])
			nextCodes[lenVal]++
			// mlog.D("CreateHuffmanTreeTable: node[%d] HuffmanCode = 0x%x", n, node[n].HuffmanCode)
		}
	}
	// CodeLengthsTree の構築
	if !SetupHuffmanTree(&codeLengthsTree, node) {
		// mlog.D("CreateHuffmanTreeTable: SetupHuffmanTree for CodeLengthsTree failed")
		return false
	}

	// 4. リテラル／長さツリーの作成
	literalLengthTree := make([]HuffmanCode, numLiteralLengthCodes)
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
			lenVal := int(GetACodeWithHuffmanTree(ctx, &codeLengthsTree))
			switch lenVal {
			case 16: // 前回のコード長を 3～6 回コピー
				repeat = uint16(GetBits(ctx, 2)) + 3
				// mlog.D("CreateHuffmanTreeTable: Literal - code 16, repeat=%d", repeat)
			case 17: // 0 を 3～10 回繰り返す
				prev = 0
				repeat = uint16(GetBits(ctx, 3)) + 3
				// mlog.D("CreateHuffmanTreeTable: Literal - code 17, repeat=%d", repeat)
			case 18: // 0 を 11～138 回繰り返す
				prev = 0
				repeat = uint16(GetBits(ctx, 7)) + 11
				// mlog.D("CreateHuffmanTreeTable: Literal - code 18, repeat=%d", repeat)
			default:
				if repeat > 0 {
					// mlog.D("CreateHuffmanTreeTable: Literal - unexpected repeat > 0")
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
				// mlog.D("CreateHuffmanTreeTable: Literal - set literalLengthTree[%d].Length = %d", n-1, lenVal)
			}
		}
	}
	if repeat > 0 {
		// mlog.D("CreateHuffmanTreeTable: Literal - repeat remaining > 0")
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
		// mlog.D("CreateHuffmanTreeTable: Literal nextCodes[%d] = %d", bits, code)
	}
	for n = 0; n < int(numLiteralLengthCodes); n++ {
		l := int(literalLengthTree[n].Length)
		if l != 0 {
			literalLengthTree[n].HuffmanCode = uint16(nextCodes[l])
			nextCodes[l]++
			// mlog.D("CreateHuffmanTreeTable: Literal literalLengthTree[%d] HuffmanCode = 0x%x", n, literalLengthTree[n].HuffmanCode)
		}
	}
	bSuccess := SetupHuffmanTree(&ctx.CustomLiteralLength, literalLengthTree)
	if !bSuccess {
		// mlog.D("CreateHuffmanTreeTable: SetupHuffmanTree for customLiteralLength failed")
		return false
	}

	// 5. 距離ツリーの作成
	if bSuccess {
		distanceTree := make([]HuffmanCode, numDistanceCodes)
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
				lenVal := int(GetACodeWithHuffmanTree(ctx, &codeLengthsTree))
				switch lenVal {
				case 16:
					repeat = uint16(GetBits(ctx, 2)) + 3
					// mlog.D("CreateHuffmanTreeTable: Distance - code 16, repeat=%d", repeat)
				case 17:
					prev = 0
					repeat = uint16(GetBits(ctx, 3)) + 3
					// mlog.D("CreateHuffmanTreeTable: Distance - code 17, repeat=%d", repeat)
				case 18:
					prev = 0
					repeat = uint16(GetBits(ctx, 7)) + 11
					// mlog.D("CreateHuffmanTreeTable: Distance - code 18, repeat=%d", repeat)
				default:
					if repeat > 0 {
						// mlog.D("CreateHuffmanTreeTable: Distance - unexpected repeat > 0")
						return false
					}
					prev = uint16(lenVal)
					repeat = 0
					distanceTree[n].Length = uint16(lenVal)
					n++
					if lenVal > maxBits {
						maxBits = lenVal
					}
					// mlog.D("CreateHuffmanTreeTable: Distance - set distanceTree[%d].Length = %d", n-1, lenVal)
				}
			}
		}
		if repeat > 0 {
			// mlog.D("CreateHuffmanTreeTable: Distance - repeat remaining > 0")
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
			// mlog.D("CreateHuffmanTreeTable: Distance nextCodes[%d] = %d", bits, code)
		}
		for n = 0; n < int(numDistanceCodes); n++ {
			l := int(distanceTree[n].Length)
			if l != 0 {
				distanceTree[n].HuffmanCode = uint16(nextCodes[l])
				nextCodes[l]++
				// mlog.D("CreateHuffmanTreeTable: Distance distanceTree[%d] HuffmanCode = 0x%x", n, distanceTree[n].HuffmanCode)
			}
		}
		bSuccess = SetupHuffmanTree(&ctx.CustomDistance, distanceTree)
		if !bSuccess {
			// mlog.D("CreateHuffmanTreeTable: SetupHuffmanTree for customDistance failed")
		}
	}

	// CodeLengthsTree のリソース解放（Go では不要なので、参照をクリア）
	codeLengthsTree.Codes = nil

	if !bSuccess {
		DisposeCustomHuffmanTree(ctx)
	}
	// mlog.D("CreateHuffmanTreeTable: 完了 - bSuccess=%v", bSuccess)
	return bSuccess
}

// DisposeCustomHuffmanTree は、DecompressionContext 内のカスタムハフマンツリー
// （CustomDistance と CustomLiteralLength）のリソースを解放（nil に設定）し、
// フィールドの値をリセットします。
func DisposeCustomHuffmanTree(ctx *DecompressionContext) {
	// mlog.D("DisposeCustomHuffmanTree: 開始")

	// CustomDistance のリソース解放
	ctx.CustomDistance.Codes = nil
	ctx.CustomDistance.NumMaxBits = 0
	ctx.CustomDistance.NumCodes = 0
	// mlog.D("DisposeCustomHuffmanTree: CustomDistance をクリア")

	// CustomLiteralLength のリソース解放
	ctx.CustomLiteralLength.Codes = nil
	ctx.CustomLiteralLength.NumMaxBits = 0
	ctx.CustomLiteralLength.NumCodes = 0
	// mlog.D("DisposeCustomHuffmanTree: CustomLiteralLength をクリア")

	// mlog.D("DisposeCustomHuffmanTree: 完了")
}

// ProcessHuffmanFixed は、固定ハフマンブロックの復号処理を行います。
// 固定ハフマンツリー (ctx.Fixed) を用いて、DecodeWithHuffmanTree を呼び出し、
// 復号処理の成否を返します。
func ProcessHuffmanFixed(ctx *DecompressionContext) bool {
	// mlog.D("ProcessHuffmanFixed: 開始")
	result := DecodeWithHuffmanTree(ctx, &ctx.Fixed, nil)
	// mlog.D("ProcessHuffmanFixed: 終了 - 結果: %v", result)
	return result
}

// ProcessUncompressed は、非圧縮ブロックのデータを処理し、
// チェックサム検証後、サイズ分のバイトを出力バッファへコピーします。
// 正常終了の場合は true を、エラー時は false を返します。
func ProcessUncompressed(ctx *DecompressionContext) bool {
	// mlog.D("ProcessUncompressed: 開始")

	// サイズ（16bit）を読み出す
	sizeLow := GetBits(ctx, 8) & 0xff
	sizeHigh := GetBits(ctx, 8) & 0xff
	size := uint16(sizeLow) | (uint16(sizeHigh) << 8)
	// mlog.D("ProcessUncompressed: サイズ = %d (0x%x)", size, size)

	// チェックサム用の値（16bit）を読み出し、反転してサイズと一致するか検証
	tmpLow := GetBits(ctx, 8) & 0xff
	tmpHigh := GetBits(ctx, 8) & 0xff
	wTmp := uint16(tmpLow) | (uint16(tmpHigh) << 8)
	// mlog.D("ProcessUncompressed: 読み出しチェックサム = 0x%x", wTmp)

	// ビット反転（Go では ^ 演算子でビットごとのNOT）
	wTmp = ^wTmp
	// mlog.D("ProcessUncompressed: 反転チェックサム = 0x%x", wTmp)

	if wTmp != size {
		// mlog.D("ProcessUncompressed: チェックサムエラー: 期待値 0x%x, 実際値 0x%x", size, wTmp)
		return false
	}

	// size バイト分のデータを出力先バッファにコピーする
	for i := 0; i < int(size); i++ {
		b := GetBits(ctx, 8) & 0xff
		if int(ctx.BytesDecompressed) >= len(ctx.Dest) {
			// mlog.D("ProcessUncompressed: 出力バッファ不足: BytesDecompressed=%d, len(Dest)=%d", ctx.BytesDecompressed, len(ctx.Dest))
			return false
		}
		ctx.Dest[ctx.BytesDecompressed] = byte(b)
		ctx.BytesDecompressed++
		// mlog.D("ProcessUncompressed: バイト[%d] 0x%x を出力, BytesDecompressed=%d", i, b, ctx.BytesDecompressed)
	}

	// mlog.D("ProcessUncompressed: 完了, BytesDecompressed=%d", ctx.BytesDecompressed)
	return true
}

// DecodeWithHuffmanTree は、literalLengthTree および distanceTree を用いて
// 圧縮データからリテラルおよび長さ・距離ペアを復号し、ctx.Dest に展開します。
// 終端コード（256）が現れた場合に処理を終了し、正常終了なら true を、エラー時は false を返します。
func DecodeWithHuffmanTree(ctx *DecompressionContext, literalLengthTree *HuffmanTree, distanceTree *HuffmanTree) bool {
	// mlog.D("DecodeWithHuffmanTree: 復号処理開始")
	for {
		// mlog.D("DecodeWithHuffmanTree: 現在の出力位置 = %d", ctx.BytesDecompressed)

		// リテラル/長さコードの取得
		w := GetACodeWithHuffmanTree(ctx, literalLengthTree)
		// mlog.D("DecodeWithHuffmanTree: GetACodeWithHuffmanTree returned 0x%x", w)

		if w == 0xffff {
			// mlog.D("DecodeWithHuffmanTree: エラーコード 0xffff を検出")
			return false
		}
		// 終端コード 256 の場合、ブロック終了
		if w == 256 {
			// mlog.D("DecodeWithHuffmanTree: 終端コード 256 を検出、復号完了")
			break
		}
		if w < 256 {
			// リテラルバイトの場合、そのまま出力
			if int(ctx.BytesDecompressed) >= len(ctx.Dest) {
				// mlog.D("DecodeWithHuffmanTree: 出力バッファ不足、BytesDecompressed=%d, len(Dest)=%d", ctx.BytesDecompressed, len(ctx.Dest))
				return false
			}
			// リテラルの場合の処理前
			// mlog.D("DecodeWithHuffmanTree: リテラルコード 0x%x を処理開始, 現在の BytesDecompressed=%d, Destの容量=%d", w, ctx.BytesDecompressed, len(ctx.Dest))

			// リテラルを出力した直後
			ctx.Dest[ctx.BytesDecompressed] = byte(w)
			ctx.BytesDecompressed++
			// mlog.D("DecodeWithHuffmanTree: リテラル 0x%x を出力完了, 更新後 BytesDecompressed=%d", w, ctx.BytesDecompressed)
		} else {
			// mlog.D("DecodeWithHuffmanTree: 現在の出力位置 (BytesDecompressed) = %d", ctx.BytesDecompressed)

			// 長さ/距離ペアの場合
			length := DecodeLength(ctx, w)
			// mlog.D("DecodeWithHuffmanTree: 長さコード 0x%x により length=%d と復号", w, length)
			var distCode uint16
			if distanceTree != nil {
				distCode = GetACodeWithHuffmanTree(ctx, distanceTree)
			} else {
				distCode = GetBits(ctx, 5)
				distCode = binreverse(int(distCode), 5)
			}
			// mlog.D("DecodeWithHuffmanTree: 距離コードとして 0x%x を取得", distCode)
			// mlog.D("DecodeWithHuffmanTree: 長さ／距離ペア処理前, 現在の BytesDecompressed=%d", ctx.BytesDecompressed)

			distance := DecodeDistance(ctx, distCode)
			// mlog.D("DecodeWithHuffmanTree: 距離コード 0x%x により distance=%d と復号", distCode, distance)

			// コピー元インデックスは、現在の出力位置から distance 分戻った位置
			start := int(ctx.BytesDecompressed) - distance
			// mlog.D("DecodeWithHuffmanTree: 距離コード=0x%x, computed distance=%d, BytesDecompressed=%d, コピー開始位置 = %d",
			// distCode, distance, ctx.BytesDecompressed, start)

			if start < 0 {
				// mlog.D("DecodeWithHuffmanTree: エラー - コピー元開始位置が負値 (start=%d)", start)
				return false
			}
			// length バイト分、過去の出力からコピーする
			for i := range length {
				if int(ctx.BytesDecompressed) >= len(ctx.Dest) || start+i >= len(ctx.Dest) {
					// mlog.D("DecodeWithHuffmanTree: エラー - 出力バッファ範囲外 (BytesDecompressed=%d, start+i=%d, len(Dest)=%d)", ctx.BytesDecompressed, start+i, len(ctx.Dest))
					return false
				}
				ctx.Dest[ctx.BytesDecompressed] = ctx.Dest[start+i]
				ctx.BytesDecompressed++
			}
			// mlog.D("DecodeWithHuffmanTree: 距離=%d から length=%d バイトをコピー, 新 BytesDecompressed=%d", distance, length, ctx.BytesDecompressed)
		}
	}
	// mlog.D("DecodeWithHuffmanTree: 復号処理完了, BytesDecompressed=%d", ctx.BytesDecompressed)
	return true
}

// DecodeLength は、baseCode（リテラル/長さコード）から実際の長さを復号します。
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
func DecodeLength(ctx *DecompressionContext, baseCode uint16) int {
	// mlog.D("DecodeLength: baseCode = %d", baseCode)
	if baseCode <= 264 {
		val := int(baseCode) - 257 + 3
		// mlog.D("DecodeLength: baseCode<=264, returning %d", val)
		return val
	}
	if baseCode == 285 {
		// mlog.D("DecodeLength: baseCode==285, returning 258")
		return 258
	}
	if baseCode > 285 {
		// mlog.D("DecodeLength: baseCode>285, error, returning 0xffff")
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
	extra := GetBits(ctx, int(numExtraBits))
	result := int(y) + int(extra)
	// mlog.D("DecodeLength: w=%d, x=%d, numExtraBits=%d, y=%d, extra=%d, result=%d", w, x, numExtraBits, y, extra, result)
	return result
}

// DecodeDistance は、baseCode（距離コード）から実際の距離を復号します。
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
func DecodeDistance(ctx *DecompressionContext, baseCode uint16) int {
	// mlog.D("DecodeDistance: baseCode = %d", baseCode)
	if baseCode <= 3 {
		val := int(baseCode) + 1
		// mlog.D("DecodeDistance: baseCode<=3, returning %d", val)
		return val
	}
	if baseCode > 29 {
		// mlog.D("DecodeDistance: baseCode>29, error, returning 0")
		return 0
	}
	w := baseCode - 4
	x := w >> 1
	numExtraBits := x + 1
	y := uint16((2 << int(numExtraBits)) + 1)
	y += uint16((int(w & 1)) << int(numExtraBits))
	extra := GetBits(ctx, int(numExtraBits))
	result := int(y) + int(extra)
	// mlog.D("DecodeDistance: w=%d, x=%d, numExtraBits=%d, y=%d, extra=%d, result=%d", w, x, numExtraBits, y, extra, result)
	return result
}

// compareCode は、HuffmanCode a のビット長とコード値と、指定された length および code を比較します。
// a.Length が length より小さい場合は -1、大きい場合は 1、等しい場合は
// a.HuffmanCode と code の大小で比較し、等しければ 0 を返します。
func compareCode(a HuffmanCode, length int, code uint16) int {
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

// GetACodeWithHuffmanTree は、ctx.Src からビットを順次読み出し、
// tree 内のハフマンコードと照合して該当するシンボル（binCode）を返します。
// 照合処理は、入力の最小ビット長から始まり、ビット数を増加させながらバイナリサーチで行います。
// 該当が見つからなかった場合は 0xffff を返します。
func GetACodeWithHuffmanTree(ctx *DecompressionContext, tree *HuffmanTree) uint16 {
	if tree.Codes == nil || len(tree.Codes) == 0 {
		// mlog.D("GetACodeWithHuffmanTree: ツリーが未設定のため 0xffff を返却")
		return 0xffff
	}
	maxBits := int(tree.NumMaxBits)
	var w uint16 = 0
	// 初期ビット数は、ツリーの先頭コードの長さとする
	minBits := int(tree.Codes[0].Length)
	// mlog.D("GetACodeWithHuffmanTree: 開始 - minBits=%d, maxBits=%d", minBits, maxBits)
	// 初期の minBits 分のビットを取得
	for bits := 0; bits < minBits; bits++ {
		w <<= 1
		w |= GetABit(ctx)
	}
	// mlog.D("GetACodeWithHuffmanTree: 初期ビット取得後 - w=0x%x", w)

	// 現在のビット数 (minBits) から最大ビット長まで拡張しながら検索
	bits := minBits
	for bits <= maxBits {
		left := 0
		right := int(tree.NumCodes)
		// mlog.D("GetACodeWithHuffmanTree: バイナリサーチ開始 - bits=%d, w=0x%x, left=%d, right=%d", bits, w, left, right)
		// バイナリサーチで該当コードを探す
		for left < right {
			mid := (left + right) / 2
			comp := compareCode(tree.Codes[mid], bits, w)
			// mlog.D("GetACodeWithHuffmanTree: バイナリサーチ mid=%d, comp=%d", mid, comp)
			if comp == 0 {
				// mlog.D("GetACodeWithHuffmanTree: 該当コード発見 - index=%d, binCode=0x%x", mid, tree.Codes[mid].BinCode)
				return tree.Codes[mid].BinCode
			} else if comp < 0 {
				left = mid + 1
			} else {
				right = mid
			}
		}
		// 該当が見つからなかったので、1ビット追加して再検索
		w <<= 1
		addedBit := GetABit(ctx)
		w |= addedBit
		bits++
		// mlog.D("GetACodeWithHuffmanTree: ビット追加 - bits=%d, addedBit=%d, w=0x%x", bits, addedBit, w)
	}
	// mlog.D("GetACodeWithHuffmanTree: 該当コードが見つからず、0xffff を返却")
	return 0xffff
}

// SetupHuffmanTree は、入力の HuffmanCode 配列 codes から
// 有効なコードを抽出し、挿入ソートにより順序付けた結果を
// tree の Codes にセットし、NumCodes および NumMaxBits を設定します。
func SetupHuffmanTree(tree *HuffmanTree, codes []HuffmanCode) bool {
	// mlog.D("SetupHuffmanTree: 開始 - 入力コード数 = %d", len(codes))

	// 一時的なスライスを確保（最大長は入力数と同じ）
	dest := make([]HuffmanCode, len(codes))
	numCodes := 0
	maxBits := 0

	// 入力の各コードについて
	for _, codeEntry := range codes {
		length := int(codeEntry.Length)
		if length != 0 {
			j := numCodes
			codeVal := int(codeEntry.HuffmanCode)
			// mlog.D("SetupHuffmanTree: コード[%d] 処理開始 - Length=%d, HuffmanCode=0x%x, BinCode=0x%x", i, length, codeEntry.HuffmanCode, codeEntry.BinCode)
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
			// mlog.D("SetupHuffmanTree: コード[%d] を位置 %d に挿入, 現在の有効コード数 = %d, maxBits = %d", i, j, numCodes, maxBits)
		}
	}
	tree.NumMaxBits = uint16(maxBits)
	tree.NumCodes = uint16(numCodes)
	tree.Codes = dest[:numCodes]
	// mlog.D("SetupHuffmanTree: 完了 - NumCodes=%d, NumMaxBits=%d", tree.NumCodes, tree.NumMaxBits)
	return true
}

// SetupFixedHuffmanTree は、固定ハフマンテーブル fixedCodeSorted を用いて
// DecompressionContext の Fixed フィールドを初期化します。
func SetupFixedHuffmanTree(ctx *DecompressionContext) {
	// mlog.D("SetupFixedHuffmanTree: 開始 - 固定ハフマンツリーのセットアップを行います")

	// 固定テーブルの要素数は 288 に固定（C++実装に合わせる）
	ctx.Fixed.NumCodes = 288
	// 最大ビット長は 9
	ctx.Fixed.NumMaxBits = 9
	// 固定ハフマンコードテーブルを割り当てる
	ctx.Fixed.Codes = fixedCodeSorted

	// mlog.D("SetupFixedHuffmanTree: 完了 - NumCodes=%d, NumMaxBits=%d", ctx.Fixed.NumCodes, ctx.Fixed.NumMaxBits)
}

// ClearContext は、DecompressionContext のフィールドを初期状態にリセットします。
func ClearContext(ctx *DecompressionContext) {
	// mlog.D("ClearContext: context の初期化を開始")

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

	// mlog.D("ClearContext: context の初期化が完了")
}

// binreverse は、引数 code の下位 length ビットを反転した結果を返します。
// 例: code=0b1011, length=4 の場合、反転結果は 0b1101 となります。
func binreverse(code int, length int) uint16 {
	var w uint16 = 0
	// mlog.D("binreverse: start - code=%d, length=%d", code, length)
	for i := 0; i < length; i++ {
		w <<= 1
		currentBit := code & 1
		w |= uint16(currentBit)
		// mlog.D("binreverse: iter=%d, currentBit=%d, w=0x%x, remaining code=%d", i, currentBit, w, code)
		code >>= 1
	}
	// mlog.D("binreverse: end - reversed value=0x%x", w)
	return w
}

// GetABit は、DecompressionContext の Src から現在のビット位置の1ビットを読み出し、
// 読み出す前のビット位置や対象のバイト、読み出したビット値をデバッグログに出力します。
func GetABit(ctx *DecompressionContext) uint16 {
	// 現在のビット位置からバイト配列上のインデックスとビット位置を算出
	index := ctx.BitsRead >> 3
	if int(index) >= len(ctx.Src) {
		// mlog.D("GetABit: index out of range: index=%d, len(Src)=%d", index, len(ctx.Src))
		return 0
	}
	byteVal := ctx.Src[index]
	bitIndex := ctx.BitsRead & 7
	// 対象のビットを取得
	bit := (byteVal >> bitIndex) & 1
	// mlog.D("GetABit: BitsRead=%d, index=%d, bitIndex=%d, byte=0x%x, bit=%d", ctx.BitsRead, index, bitIndex, byteVal, bit)
	ctx.BitsRead++ // 読み出し後にカウンタをインクリメント
	return uint16(bit)
}

// GetBits は、DecompressionContext の Src から指定された size ビット分を読み出します。
// まず、ctx.BitsRead を size 分進めた上で、後ろから1ビットずつ取り出し、結果に左シフトして加算します。
// 各ビットの取得過程を mlog.D でデバッグ出力します。
func GetBits(ctx *DecompressionContext, size int) uint16 {
	// originalBitsRead := ctx.BitsRead
	ctx.BitsRead += int64(size)
	current := ctx.BitsRead
	var result uint16 = 0
	// mlog.D("GetBits: original BitsRead=%d, size=%d, new BitsRead=%d", originalBitsRead, size, ctx.BitsRead)

	for i := 0; i < size; i++ {
		result <<= 1
		current-- // 読み出すビット位置を後ろにずらす
		index := current >> 3
		if int(index) >= len(ctx.Src) {
			// mlog.D("GetBits: index out of range: index=%d, len(Src)=%d", index, len(ctx.Src))
			return result
		}
		byteVal := ctx.Src[index]
		bitIndex := current & 7
		bit := (byteVal >> bitIndex) & 1
		result |= uint16(bit)
		// mlog.D("GetBits: i=%d, current=%d, index=%d, bitIndex=%d, byte=0x%x, bit=%d, result=0x%x",
		// 	i, current, index, bitIndex, byteVal, bit, result)
	}
	return result
}

// HuffmanCode 構造体
// C++:
//
//	typedef struct _customhuffmantable{
//	    WORD    binCode;
//	    WORD    length;
//	    WORD    huffmanCode;
//	} HuffmanCode;
type HuffmanCode struct {
	BinCode     uint16
	Length      uint16
	HuffmanCode uint16
}

// HuffmanTree 構造体
// C++:
//
//	typedef struct{
//	    WORD numMaxBits;
//	    HuffmanCode *pCode;
//	    WORD numCodes;
//	} HuffmanTree;
type HuffmanTree struct {
	NumMaxBits uint16
	Codes      []HuffmanCode
	NumCodes   uint16
}

// CustomHuffmanTree 構造体
// C++:
//
//	typedef struct{
//	    WORD numLiteralLengthCodes;
//	    WORD numDistanceCodes;
//	    WORD *pLiteralLengthTree;
//	    WORD numElementsOfLiteralLengthTree;
//	    WORD *pDistanceTree;
//	    WORD numElementsOfDistanceTree;
//	} CustomHuffmanTree;
type CustomHuffmanTree struct {
	NumLiteralLengthCodes          uint16
	NumDistanceCodes               uint16
	LiteralLengthTree              []uint16
	NumElementsOfLiteralLengthTree uint16
	DistanceTree                   []uint16
	NumElementsOfDistanceTree      uint16
}

// DecompressionContext 構造体
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
//	} DecompressionContext;
type DecompressionContext struct {
	Src                 []byte
	BitsRead            int64
	SizeInBits          int64
	SizeInBytes         uint32
	Dest                []byte
	DestSize            uint32
	BytesDecompressed   uint32
	CustomLiteralLength HuffmanTree
	CustomDistance      HuffmanTree
	Fixed               HuffmanTree
}

var fixedCodeSorted = []HuffmanCode{
	{BinCode: 0x100, Length: 7, HuffmanCode: 0x0},
	{BinCode: 0x101, Length: 7, HuffmanCode: 0x1},
	{BinCode: 0x102, Length: 7, HuffmanCode: 0x2},
	{BinCode: 0x103, Length: 7, HuffmanCode: 0x3},
	{BinCode: 0x104, Length: 7, HuffmanCode: 0x4},
	{BinCode: 0x105, Length: 7, HuffmanCode: 0x5},
	{BinCode: 0x106, Length: 7, HuffmanCode: 0x6},
	{BinCode: 0x107, Length: 7, HuffmanCode: 0x7},
	{BinCode: 0x108, Length: 7, HuffmanCode: 0x8},
	{BinCode: 0x109, Length: 7, HuffmanCode: 0x9},
	{BinCode: 0x10a, Length: 7, HuffmanCode: 0xa},
	{BinCode: 0x10b, Length: 7, HuffmanCode: 0xb},
	{BinCode: 0x10c, Length: 7, HuffmanCode: 0xc},
	{BinCode: 0x10d, Length: 7, HuffmanCode: 0xd},
	{BinCode: 0x10e, Length: 7, HuffmanCode: 0xe},
	{BinCode: 0x10f, Length: 7, HuffmanCode: 0xf},
	{BinCode: 0x110, Length: 7, HuffmanCode: 0x10},
	{BinCode: 0x111, Length: 7, HuffmanCode: 0x11},
	{BinCode: 0x112, Length: 7, HuffmanCode: 0x12},
	{BinCode: 0x113, Length: 7, HuffmanCode: 0x13},
	{BinCode: 0x114, Length: 7, HuffmanCode: 0x14},
	{BinCode: 0x115, Length: 7, HuffmanCode: 0x15},
	{BinCode: 0x116, Length: 7, HuffmanCode: 0x16},
	{BinCode: 0x117, Length: 7, HuffmanCode: 0x17},

	{BinCode: 0x0, Length: 8, HuffmanCode: 0x30},
	{BinCode: 0x1, Length: 8, HuffmanCode: 0x31},
	{BinCode: 0x2, Length: 8, HuffmanCode: 0x32},
	{BinCode: 0x3, Length: 8, HuffmanCode: 0x33},
	{BinCode: 0x4, Length: 8, HuffmanCode: 0x34},
	{BinCode: 0x5, Length: 8, HuffmanCode: 0x35},
	{BinCode: 0x6, Length: 8, HuffmanCode: 0x36},
	{BinCode: 0x7, Length: 8, HuffmanCode: 0x37},
	{BinCode: 0x8, Length: 8, HuffmanCode: 0x38},
	{BinCode: 0x9, Length: 8, HuffmanCode: 0x39},
	{BinCode: 0xa, Length: 8, HuffmanCode: 0x3a},
	{BinCode: 0xb, Length: 8, HuffmanCode: 0x3b},
	{BinCode: 0xc, Length: 8, HuffmanCode: 0x3c},
	{BinCode: 0xd, Length: 8, HuffmanCode: 0x3d},
	{BinCode: 0xe, Length: 8, HuffmanCode: 0x3e},
	{BinCode: 0xf, Length: 8, HuffmanCode: 0x3f},
	{BinCode: 0x10, Length: 8, HuffmanCode: 0x40},
	{BinCode: 0x11, Length: 8, HuffmanCode: 0x41},
	{BinCode: 0x12, Length: 8, HuffmanCode: 0x42},
	{BinCode: 0x13, Length: 8, HuffmanCode: 0x43},
	{BinCode: 0x14, Length: 8, HuffmanCode: 0x44},
	{BinCode: 0x15, Length: 8, HuffmanCode: 0x45},
	{BinCode: 0x16, Length: 8, HuffmanCode: 0x46},
	{BinCode: 0x17, Length: 8, HuffmanCode: 0x47},
	{BinCode: 0x18, Length: 8, HuffmanCode: 0x48},
	{BinCode: 0x19, Length: 8, HuffmanCode: 0x49},
	{BinCode: 0x1a, Length: 8, HuffmanCode: 0x4a},
	{BinCode: 0x1b, Length: 8, HuffmanCode: 0x4b},
	{BinCode: 0x1c, Length: 8, HuffmanCode: 0x4c},
	{BinCode: 0x1d, Length: 8, HuffmanCode: 0x4d},
	{BinCode: 0x1e, Length: 8, HuffmanCode: 0x4e},
	{BinCode: 0x1f, Length: 8, HuffmanCode: 0x4f},
	{BinCode: 0x20, Length: 8, HuffmanCode: 0x50},
	{BinCode: 0x21, Length: 8, HuffmanCode: 0x51},
	{BinCode: 0x22, Length: 8, HuffmanCode: 0x52},
	{BinCode: 0x23, Length: 8, HuffmanCode: 0x53},
	{BinCode: 0x24, Length: 8, HuffmanCode: 0x54},
	{BinCode: 0x25, Length: 8, HuffmanCode: 0x55},
	{BinCode: 0x26, Length: 8, HuffmanCode: 0x56},
	{BinCode: 0x27, Length: 8, HuffmanCode: 0x57},
	{BinCode: 0x28, Length: 8, HuffmanCode: 0x58},
	{BinCode: 0x29, Length: 8, HuffmanCode: 0x59},
	{BinCode: 0x2a, Length: 8, HuffmanCode: 0x5a},
	{BinCode: 0x2b, Length: 8, HuffmanCode: 0x5b},
	{BinCode: 0x2c, Length: 8, HuffmanCode: 0x5c},
	{BinCode: 0x2d, Length: 8, HuffmanCode: 0x5d},
	{BinCode: 0x2e, Length: 8, HuffmanCode: 0x5e},
	{BinCode: 0x2f, Length: 8, HuffmanCode: 0x5f},
	{BinCode: 0x30, Length: 8, HuffmanCode: 0x60},
	{BinCode: 0x31, Length: 8, HuffmanCode: 0x61},
	{BinCode: 0x32, Length: 8, HuffmanCode: 0x62},
	{BinCode: 0x33, Length: 8, HuffmanCode: 0x63},
	{BinCode: 0x34, Length: 8, HuffmanCode: 0x64},
	{BinCode: 0x35, Length: 8, HuffmanCode: 0x65},
	{BinCode: 0x36, Length: 8, HuffmanCode: 0x66},
	{BinCode: 0x37, Length: 8, HuffmanCode: 0x67},
	{BinCode: 0x38, Length: 8, HuffmanCode: 0x68},
	{BinCode: 0x39, Length: 8, HuffmanCode: 0x69},
	{BinCode: 0x3a, Length: 8, HuffmanCode: 0x6a},
	{BinCode: 0x3b, Length: 8, HuffmanCode: 0x6b},
	{BinCode: 0x3c, Length: 8, HuffmanCode: 0x6c},
	{BinCode: 0x3d, Length: 8, HuffmanCode: 0x6d},
	{BinCode: 0x3e, Length: 8, HuffmanCode: 0x6e},
	{BinCode: 0x3f, Length: 8, HuffmanCode: 0x6f},
	{BinCode: 0x40, Length: 8, HuffmanCode: 0x70},
	{BinCode: 0x41, Length: 8, HuffmanCode: 0x71},
	{BinCode: 0x42, Length: 8, HuffmanCode: 0x72},
	{BinCode: 0x43, Length: 8, HuffmanCode: 0x73},
	{BinCode: 0x44, Length: 8, HuffmanCode: 0x74},
	{BinCode: 0x45, Length: 8, HuffmanCode: 0x75},
	{BinCode: 0x46, Length: 8, HuffmanCode: 0x76},
	{BinCode: 0x47, Length: 8, HuffmanCode: 0x77},
	{BinCode: 0x48, Length: 8, HuffmanCode: 0x78},
	{BinCode: 0x49, Length: 8, HuffmanCode: 0x79},
	{BinCode: 0x4a, Length: 8, HuffmanCode: 0x7a},
	{BinCode: 0x4b, Length: 8, HuffmanCode: 0x7b},
	{BinCode: 0x4c, Length: 8, HuffmanCode: 0x7c},
	{BinCode: 0x4d, Length: 8, HuffmanCode: 0x7d},
	{BinCode: 0x4e, Length: 8, HuffmanCode: 0x7e},
	{BinCode: 0x4f, Length: 8, HuffmanCode: 0x7f},
	{BinCode: 0x50, Length: 8, HuffmanCode: 0x80},
	{BinCode: 0x51, Length: 8, HuffmanCode: 0x81},
	{BinCode: 0x52, Length: 8, HuffmanCode: 0x82},
	{BinCode: 0x53, Length: 8, HuffmanCode: 0x83},
	{BinCode: 0x54, Length: 8, HuffmanCode: 0x84},
	{BinCode: 0x55, Length: 8, HuffmanCode: 0x85},
	{BinCode: 0x56, Length: 8, HuffmanCode: 0x86},
	{BinCode: 0x57, Length: 8, HuffmanCode: 0x87},
	{BinCode: 0x58, Length: 8, HuffmanCode: 0x88},
	{BinCode: 0x59, Length: 8, HuffmanCode: 0x89},
	{BinCode: 0x5a, Length: 8, HuffmanCode: 0x8a},
	{BinCode: 0x5b, Length: 8, HuffmanCode: 0x8b},
	{BinCode: 0x5c, Length: 8, HuffmanCode: 0x8c},
	{BinCode: 0x5d, Length: 8, HuffmanCode: 0x8d},
	{BinCode: 0x5e, Length: 8, HuffmanCode: 0x8e},
	{BinCode: 0x5f, Length: 8, HuffmanCode: 0x8f},
	{BinCode: 0x60, Length: 8, HuffmanCode: 0x90},
	{BinCode: 0x61, Length: 8, HuffmanCode: 0x91},
	{BinCode: 0x62, Length: 8, HuffmanCode: 0x92},
	{BinCode: 0x63, Length: 8, HuffmanCode: 0x93},
	{BinCode: 0x64, Length: 8, HuffmanCode: 0x94},
	{BinCode: 0x65, Length: 8, HuffmanCode: 0x95},
	{BinCode: 0x66, Length: 8, HuffmanCode: 0x96},
	{BinCode: 0x67, Length: 8, HuffmanCode: 0x97},
	{BinCode: 0x68, Length: 8, HuffmanCode: 0x98},
	{BinCode: 0x69, Length: 8, HuffmanCode: 0x99},
	{BinCode: 0x6a, Length: 8, HuffmanCode: 0x9a},
	{BinCode: 0x6b, Length: 8, HuffmanCode: 0x9b},
	{BinCode: 0x6c, Length: 8, HuffmanCode: 0x9c},
	{BinCode: 0x6d, Length: 8, HuffmanCode: 0x9d},
	{BinCode: 0x6e, Length: 8, HuffmanCode: 0x9e},
	{BinCode: 0x6f, Length: 8, HuffmanCode: 0x9f},
	{BinCode: 0x70, Length: 8, HuffmanCode: 0xa0},
	{BinCode: 0x71, Length: 8, HuffmanCode: 0xa1},
	{BinCode: 0x72, Length: 8, HuffmanCode: 0xa2},
	{BinCode: 0x73, Length: 8, HuffmanCode: 0xa3},
	{BinCode: 0x74, Length: 8, HuffmanCode: 0xa4},
	{BinCode: 0x75, Length: 8, HuffmanCode: 0xa5},
	{BinCode: 0x76, Length: 8, HuffmanCode: 0xa6},
	{BinCode: 0x77, Length: 8, HuffmanCode: 0xa7},
	{BinCode: 0x78, Length: 8, HuffmanCode: 0xa8},
	{BinCode: 0x79, Length: 8, HuffmanCode: 0xa9},
	{BinCode: 0x7a, Length: 8, HuffmanCode: 0xaa},
	{BinCode: 0x7b, Length: 8, HuffmanCode: 0xab},
	{BinCode: 0x7c, Length: 8, HuffmanCode: 0xac},
	{BinCode: 0x7d, Length: 8, HuffmanCode: 0xad},
	{BinCode: 0x7e, Length: 8, HuffmanCode: 0xae},
	{BinCode: 0x7f, Length: 8, HuffmanCode: 0xaf},
	{BinCode: 0x80, Length: 8, HuffmanCode: 0xb0},
	{BinCode: 0x81, Length: 8, HuffmanCode: 0xb1},
	{BinCode: 0x82, Length: 8, HuffmanCode: 0xb2},
	{BinCode: 0x83, Length: 8, HuffmanCode: 0xb3},
	{BinCode: 0x84, Length: 8, HuffmanCode: 0xb4},
	{BinCode: 0x85, Length: 8, HuffmanCode: 0xb5},
	{BinCode: 0x86, Length: 8, HuffmanCode: 0xb6},
	{BinCode: 0x87, Length: 8, HuffmanCode: 0xb7},
	{BinCode: 0x88, Length: 8, HuffmanCode: 0xb8},
	{BinCode: 0x89, Length: 8, HuffmanCode: 0xb9},
	{BinCode: 0x8a, Length: 8, HuffmanCode: 0xba},
	{BinCode: 0x8b, Length: 8, HuffmanCode: 0xbb},
	{BinCode: 0x8c, Length: 8, HuffmanCode: 0xbc},
	{BinCode: 0x8d, Length: 8, HuffmanCode: 0xbd},
	{BinCode: 0x8e, Length: 8, HuffmanCode: 0xbe},
	{BinCode: 0x8f, Length: 8, HuffmanCode: 0xbf},
	{BinCode: 0x118, Length: 8, HuffmanCode: 0xc0},
	{BinCode: 0x119, Length: 8, HuffmanCode: 0xc1},
	{BinCode: 0x11a, Length: 8, HuffmanCode: 0xc2},
	{BinCode: 0x11b, Length: 8, HuffmanCode: 0xc3},
	{BinCode: 0x11c, Length: 8, HuffmanCode: 0xc4},
	{BinCode: 0x11d, Length: 8, HuffmanCode: 0xc5},
	{BinCode: 0x11e, Length: 8, HuffmanCode: 0xc6},
	{BinCode: 0x11f, Length: 8, HuffmanCode: 0xc7},
	{BinCode: 0x90, Length: 9, HuffmanCode: 0x190},
	{BinCode: 0x91, Length: 9, HuffmanCode: 0x191},
	{BinCode: 0x92, Length: 9, HuffmanCode: 0x192},
	{BinCode: 0x93, Length: 9, HuffmanCode: 0x193},
	{BinCode: 0x94, Length: 9, HuffmanCode: 0x194},
	{BinCode: 0x95, Length: 9, HuffmanCode: 0x195},
	{BinCode: 0x96, Length: 9, HuffmanCode: 0x196},
	{BinCode: 0x97, Length: 9, HuffmanCode: 0x197},
	{BinCode: 0x98, Length: 9, HuffmanCode: 0x198},
	{BinCode: 0x99, Length: 9, HuffmanCode: 0x199},
	{BinCode: 0x9a, Length: 9, HuffmanCode: 0x19a},
	{BinCode: 0x9b, Length: 9, HuffmanCode: 0x19b},
	{BinCode: 0x9c, Length: 9, HuffmanCode: 0x19c},
	{BinCode: 0x9d, Length: 9, HuffmanCode: 0x19d},
	{BinCode: 0x9e, Length: 9, HuffmanCode: 0x19e},
	{BinCode: 0x9f, Length: 9, HuffmanCode: 0x19f},
	{BinCode: 0xa0, Length: 9, HuffmanCode: 0x1a0},
	{BinCode: 0xa1, Length: 9, HuffmanCode: 0x1a1},
	{BinCode: 0xa2, Length: 9, HuffmanCode: 0x1a2},
	{BinCode: 0xa3, Length: 9, HuffmanCode: 0x1a3},
	{BinCode: 0xa4, Length: 9, HuffmanCode: 0x1a4},
	{BinCode: 0xa5, Length: 9, HuffmanCode: 0x1a5},
	{BinCode: 0xa6, Length: 9, HuffmanCode: 0x1a6},
	{BinCode: 0xa7, Length: 9, HuffmanCode: 0x1a7},
	{BinCode: 0xa8, Length: 9, HuffmanCode: 0x1a8},
	{BinCode: 0xa9, Length: 9, HuffmanCode: 0x1a9},
	{BinCode: 0xaa, Length: 9, HuffmanCode: 0x1aa},
	{BinCode: 0xab, Length: 9, HuffmanCode: 0x1ab},
	{BinCode: 0xac, Length: 9, HuffmanCode: 0x1ac},
	{BinCode: 0xad, Length: 9, HuffmanCode: 0x1ad},
	{BinCode: 0xae, Length: 9, HuffmanCode: 0x1ae},
	{BinCode: 0xaf, Length: 9, HuffmanCode: 0x1af},
	{BinCode: 0xb0, Length: 9, HuffmanCode: 0x1b0},
	{BinCode: 0xb1, Length: 9, HuffmanCode: 0x1b1},
	{BinCode: 0xb2, Length: 9, HuffmanCode: 0x1b2},
	{BinCode: 0xb3, Length: 9, HuffmanCode: 0x1b3},
	{BinCode: 0xb4, Length: 9, HuffmanCode: 0x1b4},
	{BinCode: 0xb5, Length: 9, HuffmanCode: 0x1b5},
	{BinCode: 0xb6, Length: 9, HuffmanCode: 0x1b6},
	{BinCode: 0xb7, Length: 9, HuffmanCode: 0x1b7},
	{BinCode: 0xb8, Length: 9, HuffmanCode: 0x1b8},
	{BinCode: 0xb9, Length: 9, HuffmanCode: 0x1b9},
	{BinCode: 0xba, Length: 9, HuffmanCode: 0x1ba},
	{BinCode: 0xbb, Length: 9, HuffmanCode: 0x1bb},
	{BinCode: 0xbc, Length: 9, HuffmanCode: 0x1bc},
	{BinCode: 0xbd, Length: 9, HuffmanCode: 0x1bd},
	{BinCode: 0xbe, Length: 9, HuffmanCode: 0x1be},
	{BinCode: 0xbf, Length: 9, HuffmanCode: 0x1bf},
	{BinCode: 0xc0, Length: 9, HuffmanCode: 0x1c0},
	{BinCode: 0xc1, Length: 9, HuffmanCode: 0x1c1},
	{BinCode: 0xc2, Length: 9, HuffmanCode: 0x1c2},
	{BinCode: 0xc3, Length: 9, HuffmanCode: 0x1c3},
	{BinCode: 0xc4, Length: 9, HuffmanCode: 0x1c4},
	{BinCode: 0xc5, Length: 9, HuffmanCode: 0x1c5},
	{BinCode: 0xc6, Length: 9, HuffmanCode: 0x1c6},
	{BinCode: 0xc7, Length: 9, HuffmanCode: 0x1c7},
	{BinCode: 0xc8, Length: 9, HuffmanCode: 0x1c8},
	{BinCode: 0xc9, Length: 9, HuffmanCode: 0x1c9},
	{BinCode: 0xca, Length: 9, HuffmanCode: 0x1ca},
	{BinCode: 0xcb, Length: 9, HuffmanCode: 0x1cb},
	{BinCode: 0xcc, Length: 9, HuffmanCode: 0x1cc},
	{BinCode: 0xcd, Length: 9, HuffmanCode: 0x1cd},
	{BinCode: 0xce, Length: 9, HuffmanCode: 0x1ce},
	{BinCode: 0xcf, Length: 9, HuffmanCode: 0x1cf},
	{BinCode: 0xd0, Length: 9, HuffmanCode: 0x1d0},
	{BinCode: 0xd1, Length: 9, HuffmanCode: 0x1d1},
	{BinCode: 0xd2, Length: 9, HuffmanCode: 0x1d2},
	{BinCode: 0xd3, Length: 9, HuffmanCode: 0x1d3},
	{BinCode: 0xd4, Length: 9, HuffmanCode: 0x1d4},
	{BinCode: 0xd5, Length: 9, HuffmanCode: 0x1d5},
	{BinCode: 0xd6, Length: 9, HuffmanCode: 0x1d6},
	{BinCode: 0xd7, Length: 9, HuffmanCode: 0x1d7},
	{BinCode: 0xd8, Length: 9, HuffmanCode: 0x1d8},
	{BinCode: 0xd9, Length: 9, HuffmanCode: 0x1d9},
	{BinCode: 0xda, Length: 9, HuffmanCode: 0x1da},
	{BinCode: 0xdb, Length: 9, HuffmanCode: 0x1db},
	{BinCode: 0xdc, Length: 9, HuffmanCode: 0x1dc},
	{BinCode: 0xdd, Length: 9, HuffmanCode: 0x1dd},
	{BinCode: 0xde, Length: 9, HuffmanCode: 0x1de},
	{BinCode: 0xdf, Length: 9, HuffmanCode: 0x1df},
	{BinCode: 0xe0, Length: 9, HuffmanCode: 0x1e0},
	{BinCode: 0xe1, Length: 9, HuffmanCode: 0x1e1},
	{BinCode: 0xe2, Length: 9, HuffmanCode: 0x1e2},
	{BinCode: 0xe3, Length: 9, HuffmanCode: 0x1e3},
	{BinCode: 0xe4, Length: 9, HuffmanCode: 0x1e4},
	{BinCode: 0xe5, Length: 9, HuffmanCode: 0x1e5},
	{BinCode: 0xe6, Length: 9, HuffmanCode: 0x1e6},
	{BinCode: 0xe7, Length: 9, HuffmanCode: 0x1e7},
	{BinCode: 0xe8, Length: 9, HuffmanCode: 0x1e8},
	{BinCode: 0xe9, Length: 9, HuffmanCode: 0x1e9},
	{BinCode: 0xea, Length: 9, HuffmanCode: 0x1ea},
	{BinCode: 0xeb, Length: 9, HuffmanCode: 0x1eb},
	{BinCode: 0xec, Length: 9, HuffmanCode: 0x1ec},
	{BinCode: 0xed, Length: 9, HuffmanCode: 0x1ed},
	{BinCode: 0xee, Length: 9, HuffmanCode: 0x1ee},
	{BinCode: 0xef, Length: 9, HuffmanCode: 0x1ef},
	{BinCode: 0xf0, Length: 9, HuffmanCode: 0x1f0},
	{BinCode: 0xf1, Length: 9, HuffmanCode: 0x1f1},
	{BinCode: 0xf2, Length: 9, HuffmanCode: 0x1f2},
	{BinCode: 0xf3, Length: 9, HuffmanCode: 0x1f3},
	{BinCode: 0xf4, Length: 9, HuffmanCode: 0x1f4},
	{BinCode: 0xf5, Length: 9, HuffmanCode: 0x1f5},
	{BinCode: 0xf6, Length: 9, HuffmanCode: 0x1f6},
	{BinCode: 0xf7, Length: 9, HuffmanCode: 0x1f7},
	{BinCode: 0xf8, Length: 9, HuffmanCode: 0x1f8},
	{BinCode: 0xf9, Length: 9, HuffmanCode: 0x1f9},
	{BinCode: 0xfa, Length: 9, HuffmanCode: 0x1fa},
	{BinCode: 0xfb, Length: 9, HuffmanCode: 0x1fb},
	{BinCode: 0xfc, Length: 9, HuffmanCode: 0x1fc},
	{BinCode: 0xfd, Length: 9, HuffmanCode: 0x1fd},
	{BinCode: 0xfe, Length: 9, HuffmanCode: 0x1fe},
	{BinCode: 0xff, Length: 9, HuffmanCode: 0x1ff},
}
