package repository

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

// parseCompressedBinaryXFile は、圧縮バイナリ形式の X ファイルを解凍し、
// 解凍後のデータをバイナリパーサーへ渡します。
func (rep *XRepository) parseCompressedBinaryXFile(model *pmx.PmxModel) error {
	var decompressedBuffer []byte
	var err error

	// 圧縮バイナリ形式の X ファイルを解凍します。
	if decompressedBuffer, err = rep.decompressedBinaryXFile(); err != nil {
		return fmt.Errorf("failed to decompress binary X file: %w", err)
	}

	rep.buffers = decompressedBuffer

	// // 一旦解凍後のバイナリデータをファイルに出力
	// path := strings.Replace(rep.path, ".x", "_decompressed.bin", 1)
	// os.WriteFile(path, decompressedBuffer, 0644)

	// 解凍後のデータをバイナリパーサーへ渡します。
	if err := rep.parseBinaryXFile(model); err != nil {
		return fmt.Errorf("failed to parse binary X file: %w", err)
	}

	return nil
}

type binaryTokenType uint16

// バイナリXファイルのトークン定数
const (
	tokenName        binaryTokenType = 1
	tokenString      binaryTokenType = 2
	tokenInteger     binaryTokenType = 3
	tokenGuid        binaryTokenType = 5
	tokenIntegerList binaryTokenType = 6
	tokenFloatList   binaryTokenType = 7
	tokenMatrix4x4   binaryTokenType = 9
	tokenOBrace      binaryTokenType = 10
	tokenCBrace      binaryTokenType = 11
	tokenOParen      binaryTokenType = 12
	tokenCParen      binaryTokenType = 13
	tokenOBracket    binaryTokenType = 14
	tokenCBracket    binaryTokenType = 15
	tokenOAngle      binaryTokenType = 16
	tokenCAngle      binaryTokenType = 17
	tokenDot         binaryTokenType = 18
	tokenComma       binaryTokenType = 19
	tokenSemicolon   binaryTokenType = 20
	tokenTemplate    binaryTokenType = 31
	tokenWord        binaryTokenType = 40
	tokenDword       binaryTokenType = 41
	tokenFloat       binaryTokenType = 42
	tokenDouble      binaryTokenType = 43
	tokenChar        binaryTokenType = 44
	tokenUchar       binaryTokenType = 45
	tokenSword       binaryTokenType = 46
	tokenSdword      binaryTokenType = 47
	tokenVoid        binaryTokenType = 48
	tokenLpstr       binaryTokenType = 49
	tokenUnicode     binaryTokenType = 50
	tokenCstring     binaryTokenType = 51
	tokenArray       binaryTokenType = 52
	tokenEOF         binaryTokenType = 9999 // 定義にはない
)

const (
	maxLoopCount = math.MaxInt
)

// parseBinaryXFile は、バイナリ形式の X ファイルを解析します。
// https://learn.microsoft.com/ja-jp/previous-versions/direct-x/cc371722(v=msdn.10)
func (rep *XRepository) parseBinaryXFile(model *pmx.PmxModel) error {
	// ヘッダーの解析
	if err := rep.parseHeader(); err != nil {
		return fmt.Errorf("[parseBinaryXFile] failed to parse header: %w", err)
	}

	for n := range maxLoopCount {
		tokenType, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseBinaryXFile][%04d] failed to get token: %w", n, err)
		}

		switch tokenType {
		case tokenEOF:
			return nil
		case tokenTemplate:
			if err := rep.parseTemplate(); err != nil {
				return fmt.Errorf("[parseBinaryXFile][%04d] failed to parse template: %w", n, err)
			}
		case tokenName:
			// 名前トークンが出現したらオブジェクト識別子として処理
			rep.pos -= 2 // トークンを再度読めるようにポジションを戻す

			// 識別子を読み取る
			identifier, err := rep.readName()
			if err != nil {
				return fmt.Errorf("[parseBinaryXFile][%04d] failed to read object identifier: %w", n, err)
			}

			mlog.D("Processing object with identifier: %s\n", identifier)

			// オブジェクト処理
			if err := rep.parseObject(); err != nil {
				return fmt.Errorf("[parseBinaryXFile][%04d] failed to parse object '%s': %w", n, identifier, err)
			}

		case tokenWord, tokenDword, tokenFloat, tokenDouble, tokenChar, tokenUchar,
			tokenSword, tokenSdword, tokenLpstr, tokenUnicode, tokenCstring:
			// プリミティブ型をオブジェクト識別子として処理
			mlog.D("[%04d] Processing object with primitive type identifier: %d\n", n, tokenType)

			// オブジェクト処理
			if err := rep.parseObject(); err != nil {
				return fmt.Errorf("[parseBinaryXFile][%04d] failed to parse object with primitive type identifier %d: %w", n, tokenType, err)
			}

		default:
			// 予期しないトークンの場合
			return fmt.Errorf("[parseBinaryXFile][%04d] unexpected token at top level: %d", n, tokenType)
		}
	}

	return nil
}

func (rep *XRepository) parseHeader() error {
	// ヘッダーのマジックナンバーを読み取る
	if formatMagic, err := rep.readBytes(4); err != nil {
		return fmt.Errorf("[parseBinaryXFile] failed to read format magic: %w", err)
	} else {
		mlog.D("Format magic: %s\n", string(formatMagic))
	}

	// バージョン番号を読み取る
	if version, err := rep.readBytes(4); err != nil {
		return fmt.Errorf("[parseBinaryXFile] failed to read version number: %w", err)
	} else {
		mlog.D("Version number: %s\n", string(version))
	}

	// ファイルの種類を読み取る
	if fileType, err := rep.readBytes(4); err != nil {
		return fmt.Errorf("[parseBinaryXFile] failed to read file type: %w", err)
	} else {
		mlog.D("File type: %s\n", string(fileType))
	}

	// ファイルのFloatサイズを読み取る
	if floatSize, err := rep.readBytes(4); err != nil {
		return fmt.Errorf("[parseBinaryXFile] failed to read float size: %w", err)
	} else {
		mlog.D("Float size: %s\n", string(floatSize))
		if string(floatSize) == "0032" {
			rep.floatSize = 4 // 32bit
		} else {
			rep.floatSize = 8 // 64bit
		}
	}

	return nil
}

func (rep *XRepository) parseTemplate() error {
	// テンプレート名を読み取る
	name, err := rep.readName()
	if err != nil {
		return fmt.Errorf("[parseTemplate] failed to read template name: %w", err)
	}

	// 開始ブレースを確認
	token, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseTemplate] failed to read token after template name: %w", err)
	}
	if token != tokenOBrace {
		return fmt.Errorf("[parseTemplate] expected opening brace after template name, got token %d", token)
	}

	// GUID（クラスID）を読み取る
	guid, err := rep.readGuid()
	if err != nil {
		return fmt.Errorf("[parseTemplate] failed to read template GUID: %w", err)
	}

	// template_parts の解析
	if err := rep.parseTemplateParts(); err != nil {
		return fmt.Errorf("[parseTemplate] failed to parse template parts: %w", err)
	}

	// 終了ブレースを確認
	token, err = rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseTemplate] failed to read token at end of template: %w", err)
	}
	if token != tokenCBrace {
		return fmt.Errorf("[parseTemplate] expected closing brace at end of template, got token %d", token)
	}

	// テンプレート情報をログに出力
	mlog.D("Parsed template: name=%s, guid=%d\n", name, guid)
	return nil
}

// parseTemplateParts は、テンプレートの本体部分を解析します
func (rep *XRepository) parseTemplateParts() error {
	// 次のトークンを先読みして形式を判断
	peekToken, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseTemplateParts] failed to peek token for template parts: %w", err)
	}
	rep.pos -= 2 // 先読みしたトークンを戻す

	if peekToken == tokenOBracket {
		// template_members_part TOKEN_OBRACKET template_option_info TOKEN_CBRACKET の形式

		// まず template_members_part（オプショナル）を処理
		// template_members_part は空の場合と template_members_list の場合がある
		nextPeek, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseTemplateParts] failed to peek next token: %w", err)
		}
		rep.pos -= 2 // 先読みしたトークンを戻す

		// 次のトークンがOBracketでない場合は、template_members_listがある
		if nextPeek != tokenOBracket {
			if err := rep.parseTemplateMembersList(); err != nil {
				return fmt.Errorf("[parseTemplateParts] failed to parse template members part: %w", err)
			}
		}

		// TOKEN_OBRACKET
		bracketToken, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseTemplateParts] failed to read opening bracket token: %w", err)
		}
		if bracketToken != tokenOBracket {
			return fmt.Errorf("[parseTemplateParts] expected opening bracket in template parts, got token %d", bracketToken)
		}

		// template_option_info の処理
		if err := rep.parseTemplateOptionInfo(); err != nil {
			return fmt.Errorf("[parseTemplateParts] failed to parse template option info: %w", err)
		}

		// 閉じ括弧を期待
		closeBracket, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseTemplateParts] failed to read closing bracket token: %w", err)
		}
		if closeBracket != tokenCBracket {
			return fmt.Errorf("[parseTemplateParts] expected closing bracket after template option info, got token %d", closeBracket)
		}
	} else {
		// template_members_list の形式
		if err := rep.parseTemplateMembersList(); err != nil {
			return fmt.Errorf("[parseTemplateParts] failed to parse template members list: %w", err)
		}

		bracketToken, err := rep.getToken()
		rep.pos -= 2 // トークンを戻す
		if err != nil {
			return fmt.Errorf("[parseTemplateParts] failed to read opening bracket token: %w", err)
		}
		if bracketToken == tokenOBracket {
			return rep.parseTemplateParts()
		}
	}

	return nil
}

// parseTemplateMembersList は、テンプレートメンバーのリストを解析します
func (rep *XRepository) parseTemplateMembersList() error {
	for n := range maxLoopCount {
		// 先読みして次のトークンが有効なtemplate_membersの開始かどうかを確認
		peekToken, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseTemplateMembersList][%04d] failed to peek token: %w", n, err)
		}
		rep.pos -= 2 // 先読みしたトークンを戻す

		// template_membersの範囲外のトークンならリスト終了
		if peekToken == tokenOBracket || peekToken == tokenCBrace ||
			peekToken == tokenCBracket || peekToken == tokenEOF {
			return nil
		}

		// template_membersの処理
		if err := rep.parseTemplateMembers(); err != nil {
			return fmt.Errorf("[parseTemplateMembersList][%04d] failed to parse template members: %w", n, err)
		}
	}

	return nil
}

// parseTemplateMembers は、単一のテンプレートメンバーを解析します
func (rep *XRepository) parseTemplateMembers() error {
	token, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseTemplateMembers] failed to read template member token: %w", err)
	}

	switch token {
	case tokenWord, tokenDword, tokenFloat, tokenDouble, tokenChar, tokenUchar,
		tokenSword, tokenSdword, tokenLpstr, tokenUnicode, tokenCstring:
		// プリミティブ型の処理
		if err := rep.parsePrimitive(token); err != nil {
			return fmt.Errorf("[parseTemplateMembers] failed to parse primitive type: %w", err)
		}
	case tokenArray:
		// 配列の処理
		if err := rep.parseArray(); err != nil {
			return fmt.Errorf("[parseTemplateMembers] failed to parse array: %w", err)
		}
	case tokenName:
		// 名前が来た場合はテンプレート参照
		rep.pos -= 2 // 名前トークンを再度読めるように位置を戻す
		_, err := rep.readName()
		if err != nil {
			return fmt.Errorf("[parseTemplateMembers] failed to read template reference name: %w", err)
		}

		// オプショナルな名前があるかチェック
		nextToken, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseTemplateMembers] failed to read token after template reference: %w", err)
		}

		if nextToken == tokenName {
			rep.pos -= 2
			_, err = rep.readName() // オプショナル名を読み飛ばす
			if err != nil {
				return fmt.Errorf("[parseTemplateMembers] failed to read optional name: %w", err)
			}
		} else {
			rep.pos -= 2 // トークンを戻す
		}

		// セミコロンを期待
		nextToken, err = rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseTemplateMembers] failed to read token after reference: %w", err)
		}
		if nextToken != tokenSemicolon {
			return fmt.Errorf("[parseTemplateMembers] expected semicolon after template reference, got token %d", nextToken)
		}
	default:
		return fmt.Errorf("[parseTemplateMembers] unexpected token in template members: %d", token)
	}

	return nil
}

func (rep *XRepository) parseTemplateOptionInfo() error {
	token, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseTemplateOptionInfo] failed to read template option info token: %w", err)
	}

	if token == tokenDot {
		// ellipsis ("...") の処理
		// 続けて2つのドットを期待
		for range 2 {
			dotToken, err := rep.getToken()
			if err != nil {
				return fmt.Errorf("[parseTemplateOptionInfo] failed to read ellipsis dot token: %w", err)
			}
			if dotToken != tokenDot {
				return fmt.Errorf("[parseTemplateOptionInfo] expected dot token in ellipsis, got token %d", dotToken)
			}
		}
		mlog.D("Parsed ellipsis in template option info")
	} else {
		// template_option_list の処理
		rep.pos -= 2 // 先読みしたトークンを戻す
		if err := rep.parseTemplateOptionList(); err != nil {
			return fmt.Errorf("[parseTemplateOptionInfo] failed to parse template option list: %w", err)
		}
	}

	return nil
}

// データオブジェクトの解析
func (rep *XRepository) parseObject() error {
	// オブジェクト識別子は既に読み取り済み（呼び出し元で処理）

	// オプショナルな名前があるかチェック
	nextToken, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseObject] failed to read token after object identifier: %w", err)
	}

	if nextToken == tokenName {
		rep.pos -= 2
		name, err := rep.readName() // オブジェクト名を読み飛ばす
		if err != nil {
			return fmt.Errorf("[parseObject] failed to read object name: %w", err)
		}
		mlog.D("Parsed object name: %s at position %d (0x%08X)\n", name, rep.pos, rep.pos)
	} else {
		rep.pos -= 2 // トークンを戻す
	}

	// 開始ブレースを期待
	braceToken, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseObject] failed to read opening brace token: %w", err)
	}
	if braceToken != tokenOBrace {
		return fmt.Errorf("[parseObject] expected opening brace for object, got token %d", braceToken)
	}

	// オプショナルなクラスIDがあるか確認
	nextToken, err = rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseObject] failed to read token after opening brace: %w", err)
	}
	rep.pos -= 2 // トークンを戻す

	if nextToken == tokenGuid {
		// クラスIDを読み飛ばす
		_, err = rep.readGuid()
		if err != nil {
			return fmt.Errorf("[parseObject] failed to read object class ID: %w", err)
		}
	}

	// データパートリストを解析
	for n := range maxLoopCount {
		token, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseObject][%04d] failed to read token in data parts list: %w", n, err)
		}

		if token == tokenCBrace {
			// オブジェクト終了
			return nil
		}

		rep.pos -= 2 // トークンを処理メソッドで読めるように戻す

		switch token {
		case tokenOBrace:
			// データ参照
			if err := rep.parseDataReference(); err != nil {
				return fmt.Errorf("[parseObject][%04d] failed to parse data reference: %w", n, err)
			}
		case tokenName:
			// 名前で始まる場合はネストされたオブジェクトか、データ型を表す識別子
			peekToken, err := rep.getToken()
			rep.pos -= 2 // 覗いたトークンを戻す

			if peekToken == tokenName || peekToken == tokenOBrace {
				// ネストされたオブジェクト

				if peekToken == tokenName {
					// 識別子を読み取る
					identifier, err := rep.readName()
					if err != nil {
						return fmt.Errorf("[parseObject][%04d] failed to read object identifier: %w", n, err)
					}

					mlog.D("Processing object with identifier: %s\n", identifier)
				}

				if err := rep.parseObject(); err != nil {
					return fmt.Errorf("[parseObject][%04d] failed to parse nested object: %w", n, err)
				}
			} else {
				// データ識別子の後にリストがくると予想される
				_, err = rep.readName() // 識別子を読み飛ばす
				if err != nil {
					return fmt.Errorf("[parseObject][%04d] failed to read data identifier: %w", n, err)
				}

				// リスト処理
				if err := rep.parseDataList(); err != nil {
					return fmt.Errorf("[parseObject][%04d] failed to parse data list: %w", n, err)
				}
			}
		case tokenIntegerList:
			// 整数リスト
			_, err := rep.readIntegerList()
			if err != nil {
				return fmt.Errorf("[parseObject][%04d] failed to read integer list: %w", n, err)
			}
		case tokenFloatList:
			// 浮動小数点リスト
			_, err := rep.readFloatList()
			if err != nil {
				return fmt.Errorf("[parseObject][%04d] failed to read float list: %w", n, err)
			}
		case tokenString:
			// 文字列リスト
			if err := rep.parseStringList(); err != nil {
				return fmt.Errorf("[parseObject][%04d] failed to parse string list: %w", n, err)
			}
		default:
			return fmt.Errorf("[parseObject][%04d] unexpected token in data parts list: %d", n, token)
		}
	}

	return nil
}

// データ参照の解析
func (rep *XRepository) parseDataReference() error {
	// 開始ブレースは既に読み取り済み
	_, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseDataReference] failed to read opening brace: %w", err)
	}

	// 参照名を読み取る
	nameToken, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseDataReference] failed to read reference name token: %w", err)
	}
	if nameToken != tokenName {
		return fmt.Errorf("[parseDataReference] expected name token for data reference, got token %d", nameToken)
	}

	rep.pos -= 2
	_, err = rep.readName() // 参照名を読み飛ばす
	if err != nil {
		return fmt.Errorf("[parseDataReference] failed to read reference name: %w", err)
	}

	// オプショナルなクラスID
	nextToken, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseDataReference] failed to read token after reference name: %w", err)
	}

	if nextToken == tokenGuid {
		_, err = rep.readGuid() // クラスIDを読み飛ばす
		if err != nil {
			return fmt.Errorf("[parseDataReference] failed to read reference class ID: %w", err)
		}
	} else {
		rep.pos -= 2 // トークンを戻す
	}

	// 閉じブレースを期待
	closeBrace, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseDataReference] failed to read closing brace: %w", err)
	}
	if closeBrace != tokenCBrace {
		return fmt.Errorf("[parseDataReference] expected closing brace for data reference, got token %d", closeBrace)
	}

	return nil
}

// データリストの解析
func (rep *XRepository) parseDataList() error {
	// リストの種類を判断
	token, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseDataList] failed to read data list token: %w", err)
	}

	rep.pos -= 2 // トークンを処理メソッドで読めるように戻す

	switch token {
	case tokenIntegerList:
		_, err := rep.readIntegerList()
		if err != nil {
			return fmt.Errorf("[parseDataList] failed to read integer list: %w", err)
		}
	case tokenFloatList:
		_, err := rep.readFloatList()
		if err != nil {
			return fmt.Errorf("[parseDataList] failed to read float list: %w", err)
		}
	case tokenString:
		if err := rep.parseStringList(); err != nil {
			return fmt.Errorf("[parseDataList] failed to parse string list: %w", err)
		}
	default:
		return fmt.Errorf("[parseDataList] unexpected token for data list: %d", token)
	}

	return nil
}

// 文字列リストの解析
func (rep *XRepository) parseStringList() error {
	for n := range maxLoopCount {
		// 文字列を読み取る
		stringToken, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseStringList][%04d] failed to read string token: %w", n, err)
		}
		if stringToken != tokenString {
			return fmt.Errorf("[parseStringList][%04d] expected string token in string list, got token %d", n, stringToken)
		}

		rep.pos -= 2
		_, err = rep.readString() // 文字列を読み飛ばす（終了トークンも含む）
		if err != nil {
			return fmt.Errorf("[parseStringList][%04d] failed to read string: %w", n, err)
		}

		// 次のトークンを確認
		nextToken, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseStringList][%04d] failed to read token after string: %w", n, err)
		}

		if nextToken != tokenComma && nextToken != tokenSemicolon {
			rep.pos -= 2 // トークンを戻す
			break        // リスト終了
		}
	}

	return nil
}

// parseTemplateOptionList は、テンプレートオプションリストを解析します
func (rep *XRepository) parseTemplateOptionList() error {
	for n := range maxLoopCount {
		// 次のトークンをチェックして閉じ括弧なら終了
		peekToken, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseTemplateOptionList][%04d] failed to read token in template option list: %w", n, err)
		}

		if peekToken == tokenCBracket {
			rep.pos -= 2 // トークンを戻す（parseTemplateOptionInfoで処理）
			return nil
		}

		rep.pos -= 2 // トークンを戻す

		// 名前を期待
		nameToken, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseTemplateOptionList][%04d] failed to read name token in template option part: %w", n, err)
		}

		if nameToken != tokenName {
			return fmt.Errorf("[parseTemplateOptionList][%04d] expected name token in template option part, got token %d", n, nameToken)
		}

		rep.pos -= 2
		_, err = rep.readName() // オプション名を読み飛ばす
		if err != nil {
			return fmt.Errorf("[parseTemplateOptionList][%04d] failed to read option name: %w", n, err)
		}

		// オプショナルなクラスID
		nextToken, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseTemplateOptionList][%04d] failed to read token after option name: %w", n, err)
		}
		rep.pos -= 2 // トークンを戻す

		if nextToken == tokenGuid {
			_, err = rep.readGuid() // クラスIDを読み飛ばす
			if err != nil {
				return fmt.Errorf("[parseTemplateOptionList][%04d] failed to read option class ID: %w", n, err)
			}
		}
	}

	return nil
}

// プリミティブ型の処理
func (rep *XRepository) parsePrimitive(token binaryTokenType) error {
	// オプショナルな名前があるかチェック
	nextToken, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parsePrimitive] failed to read token after primitive type: %w", err)
	}

	if nextToken == tokenName {
		rep.pos -= 2
		_, err = rep.readName() // 名前を読み飛ばす
		if err != nil {
			return fmt.Errorf("[parsePrimitive] failed to read primitive name: %w", err)
		}
	} else {
		rep.pos -= 2 // トークンを戻す
	}

	// セミコロンを期待
	nextToken, err = rep.getToken()
	if err != nil {
		return fmt.Errorf("[parsePrimitive] failed to read token after primitive: %w", err)
	}
	if nextToken != tokenSemicolon {
		return fmt.Errorf("[parsePrimitive] expected semicolon after primitive, got token %d", nextToken)
	}

	return nil
}

// 配列の処理
func (rep *XRepository) parseArray() error {
	// 配列データ型を読み取る
	arrayDataType, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseArray] failed to read array data type: %w", err)
	}

	// プリミティブ型か名前型か判断
	if arrayDataType == tokenName {
		rep.pos -= 2
		_, err = rep.readName() // 型名を読み飛ばす
		if err != nil {
			return fmt.Errorf("[parseArray] failed to read array type name: %w", err)
		}
	} else if arrayDataType >= tokenWord && arrayDataType <= tokenCstring {
		// プリミティブ型は既に読み取り済み
	} else {
		return fmt.Errorf("[parseArray] unexpected array data type token: %d", arrayDataType)
	}

	// 配列名を読み取る
	nextToken, err := rep.getToken()
	if err != nil {
		return fmt.Errorf("[parseArray] failed to read token after array type: %w", err)
	}
	if nextToken != tokenName {
		return fmt.Errorf("[parseArray] expected name after array type, got token %d", nextToken)
	}

	rep.pos -= 2
	_, err = rep.readName() // 配列名を読み飛ばす
	if err != nil {
		return fmt.Errorf("[parseArray] failed to read array name: %w", err)
	}

	// 次元リストを解析
	for n := range maxLoopCount {
		nextToken, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseArray][%04d] failed to read token in dimension list: %w", n, err)
		}

		if nextToken == tokenSemicolon {
			// 配列定義の終了
			break
		}

		if nextToken != tokenOBracket {
			return fmt.Errorf("[parseArray][%04d] expected opening bracket for dimension, got token %d", n, nextToken)
		}

		// 次元サイズを読み取る
		sizeToken, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseArray][%04d] failed to read dimension size token: %w", n, err)
		}

		if sizeToken == tokenInteger {
			_, err = rep.readInteger() // サイズを読み飛ばす
			if err != nil {
				return fmt.Errorf("[parseArray][%04d] failed to read dimension size integer: %w", n, err)
			}
		} else if sizeToken == tokenName {
			rep.pos -= 2
			_, err = rep.readName() // 名前を読み飛ばす
			if err != nil {
				return fmt.Errorf("[parseArray][%04d] failed to read dimension size name: %w", n, err)
			}
		} else {
			return fmt.Errorf("[parseArray][%04d] unexpected dimension size token: %d", n, sizeToken)
		}

		// 閉じ括弧を期待
		closeBracket, err := rep.getToken()
		if err != nil {
			return fmt.Errorf("[parseArray] failed to read closing bracket token: %w", err)
		}
		if closeBracket != tokenCBracket {
			return fmt.Errorf("[parseArray] expected closing bracket after dimension size, got token %d", closeBracket)
		}

	}

	return nil
}

func (rep *XRepository) getToken() (binaryTokenType, error) {
	if rep.pos > len(rep.buffers)-2 {
		return tokenEOF, nil
	}

	// トークンタイプを読み取る
	tokenType := binary.LittleEndian.Uint16(rep.buffers[rep.pos:])
	// // トークン値のデバッグ出力
	// mlog.D("[getToken] %d (0x%04x) at position %d (0x%08X)", tokenType, uint16(tokenType), rep.pos, rep.pos)

	rep.pos += 2

	return binaryTokenType(tokenType), nil
}

// 基本的なデータ型を読み取るヘルパーメソッド
func (rep *XRepository) readDWORD() (uint32, error) {
	if rep.pos > len(rep.buffers)-4 {
		return 0, fmt.Errorf("[readDWORD] buffer overflow when reading DWORD")
	}
	value := binary.LittleEndian.Uint32(rep.buffers[rep.pos:])
	rep.pos += 4
	return value, nil
}

func (rep *XRepository) readWORD() (uint16, error) {
	if rep.pos > len(rep.buffers)-2 {
		return 0, fmt.Errorf("[readWORD] buffer overflow when reading WORD")
	}
	value := binary.LittleEndian.Uint16(rep.buffers[rep.pos:])
	rep.pos += 2
	return value, nil
}

func (rep *XRepository) readFloat() (float64, error) {
	if rep.pos > len(rep.buffers)-rep.floatSize {
		return 0, fmt.Errorf("[readFloat] buffer overflow when reading float")
	}
	if rep.floatSize == 4 {
		bits := binary.LittleEndian.Uint32(rep.buffers[rep.pos:])
		rep.pos += rep.floatSize
		return float64(math.Float32frombits(bits)), nil
	}

	bits := binary.LittleEndian.Uint64(rep.buffers[rep.pos:])
	rep.pos += rep.floatSize
	return math.Float64frombits(bits), nil
}

func (rep *XRepository) readDouble() (float64, error) {
	if rep.pos > len(rep.buffers)-8 {
		return 0, fmt.Errorf("[readDouble] buffer overflow when reading double")
	}
	bits := binary.LittleEndian.Uint64(rep.buffers[rep.pos:])
	rep.pos += 8
	return math.Float64frombits(bits), nil
}

func (rep *XRepository) readBytes(count int) ([]byte, error) {
	if rep.pos > len(rep.buffers)-count {
		return nil, fmt.Errorf("[readBytes] buffer overflow when reading %d bytes", count)
	}
	bytes := make([]byte, count)
	copy(bytes, rep.buffers[rep.pos:rep.pos+count])
	rep.pos += count
	return bytes, nil
}

// NAME トークンの処理
func (rep *XRepository) readName() (string, error) {
	tokenType, err := rep.getToken()
	if err != nil {
		return "", fmt.Errorf("[readName] failed to get token: %w", err)
	}
	if tokenType != tokenName {
		return "", fmt.Errorf("[readName] expected NAME token, got token %d", tokenType)
	}

	count, err := rep.readDWORD()
	if err != nil {
		return "", fmt.Errorf("[readName] failed to read NAME count: %w", err)
	}

	nameBytes, err := rep.readBytes(int(count))
	if err != nil {
		return "", fmt.Errorf("[readName] failed to read NAME data: %w", err)
	}

	// NULL文字を削除（ASCII名はNULL終端の可能性がある）
	nameBytes = bytes.TrimRight(nameBytes, "\x00")
	name := string(nameBytes)
	mlog.D("Read name: %s\n", name)
	return name, nil
}

// STRING トークンの処理
func (rep *XRepository) readString() (string, error) {
	tokenType, err := rep.getToken()
	if err != nil {
		return "", fmt.Errorf("[readString] failed to get token: %w", err)
	}
	if tokenType != tokenString {
		return "", fmt.Errorf("[readString] expected STRING token, got token %d", tokenType)
	}

	count, err := rep.readDWORD()
	if err != nil {
		return "", fmt.Errorf("[readString] failed to read STRING count: %w", err)
	}

	stringBytes, err := rep.readBytes(int(count))
	if err != nil {
		return "", fmt.Errorf("[readString] failed to read STRING data: %w", err)
	}

	// 終了トークン（SEMICOLON または COMMA）を読み取り
	terminatorToken, err := rep.getToken()
	if err != nil {
		return "", fmt.Errorf("[readString] failed to read STRING terminator: %w", err)
	}

	if terminatorToken != tokenSemicolon && terminatorToken != tokenComma {
		return "", fmt.Errorf("[readString] invalid STRING terminator token: %d", terminatorToken)
	}

	str := string(stringBytes)
	mlog.D("Read string: %s\n", str)
	return str, nil
}

// INTEGER トークンの処理
func (rep *XRepository) readInteger() (uint32, error) {
	val, err := rep.readDWORD()
	mlog.D("Read integer: %d\n", val)
	return val, err
}

// GUID トークンの処理
func (rep *XRepository) readGuid() ([16]byte, error) {
	mlog.D("[readGuid] position %d (0x%08X)", rep.pos, rep.pos)

	var guid [16]byte

	tokenType, err := rep.getToken()
	if err != nil {
		return guid, fmt.Errorf("[readGuid] failed to get token: %w", err)
	}
	if tokenType != tokenGuid {
		return guid, fmt.Errorf("[readGuid] expected GUID token, got token %d", tokenType)
	}

	// UUID data1フィールド（4バイト）
	data1Bytes, err := rep.readBytes(4)
	if err != nil {
		return guid, fmt.Errorf("[readGuid] failed to read GUID data1: %w", err)
	}

	// UUID data2フィールド（2バイト）
	data2Bytes, err := rep.readBytes(2)
	if err != nil {
		return guid, fmt.Errorf("[readGuid] failed to read GUID data2: %w", err)
	}

	// UUID data3フィールド（2バイト）
	data3Bytes, err := rep.readBytes(2)
	if err != nil {
		return guid, fmt.Errorf("[readGuid] failed to read GUID data3: %w", err)
	}

	// UUID data4フィールド（8バイト）
	data4Bytes, err := rep.readBytes(8)
	if err != nil {
		return guid, fmt.Errorf("[readGuid] failed to read GUID data4: %w", err)
	}

	// GUIDの各部分をコピー
	copy(guid[0:4], data1Bytes)
	copy(guid[4:6], data2Bytes)
	copy(guid[6:8], data3Bytes)
	copy(guid[8:16], data4Bytes)

	mlog.D("Read GUID: %v\n", guid)
	return guid, nil
}

// INTEGER_LIST トークンの処理
func (rep *XRepository) readIntegerList() ([]uint32, error) {
	tokenType, err := rep.getToken()
	if err != nil {
		return nil, fmt.Errorf("[readIntegerList] failed to get token: %w", err)
	}
	if tokenType != tokenIntegerList {
		return nil, fmt.Errorf("[readIntegerList] expected INTEGER_LIST token, got token %d", tokenType)
	}

	count, err := rep.readDWORD()
	if err != nil {
		return nil, fmt.Errorf("[readIntegerList] failed to read INTEGER_LIST count: %w", err)
	}

	list := make([]uint32, count)
	for i := range count {
		value, err := rep.readDWORD()
		if err != nil {
			return nil, fmt.Errorf("[readIntegerList] failed to read INTEGER_LIST item %d: %w", i, err)
		}
		list[i] = value
	}

	mlog.D("Read integer list: %v\n", list[:min(10, len(list))])
	return list, nil
}

// FLOAT_LIST トークンの処理
func (rep *XRepository) readFloatList() ([]float64, error) {
	tokenType, err := rep.getToken()
	if err != nil {
		return nil, fmt.Errorf("[readFloatList] failed to get token: %w", err)
	}
	if tokenType != tokenFloatList {
		return nil, fmt.Errorf("[readFloatList] expected FLOAT_LIST token, got token %d", tokenType)
	}

	count, err := rep.readDWORD()
	if err != nil {
		return nil, fmt.Errorf("[readFloatList] failed to read FLOAT_LIST count: %w", err)
	}

	list := make([]float64, count)
	for i := range count {
		value, err := rep.readFloat()
		if err != nil {
			return nil, fmt.Errorf("[readFloatList] failed to read FLOAT_LIST item %d: %w", i, err)
		}
		list[i] = value
	}

	mlog.D("Read float list: %v\n", list[:min(10, len(list))])
	return list, nil
}
