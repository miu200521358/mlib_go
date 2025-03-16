package repository

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mstring"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type xFormat int

const (
	xFormatText             xFormat = iota // テキスト形式
	xFormatBinary                          // 未圧縮バイナリ形式
	xFormatCompressedBinary                // 圧縮バイナリ形式
	xFormatInvalid                         // 無効な形式
)

type XRepository struct {
	*baseRepository[*pmx.PmxModel]
	tokens []textToken // テキスト形式のトークン
	pos    int

	buffers            []byte // バイナリ形式のバッファ
	floatSize          int    // バイナリ形式の浮動小数点数のサイズ
	materialTokenCount int    // マテリアルトークンの出現数
}

func NewXRepository() *XRepository {
	return &XRepository{
		baseRepository: &baseRepository[*pmx.PmxModel]{
			newFunc: func(path string) *pmx.PmxModel {
				return pmx.NewPmxModel(path)
			},
		},
	}
}

func (rep *XRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	return nil
}

func (rep *XRepository) CanLoad(path string) (bool, error) {
	if isExist, err := mfile.ExistsFile(path); err != nil || !isExist {
		return false, fmt.Errorf("%s", mi18n.T("ファイル存在エラー", map[string]interface{}{"Path": path}))
	}

	_, _, ext := mfile.SplitPath(path)
	if strings.ToLower(ext) != ".x" {
		return false, fmt.Errorf("%s", mi18n.T("拡張子エラー", map[string]interface{}{"Path": path, "Ext": ".x"}))
	}

	return true, nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *XRepository) Load(path string) (core.IHashModel, error) {
	runtime.GOMAXPROCS(int(runtime.NumCPU()))
	defer runtime.GOMAXPROCS(max(1, int(runtime.NumCPU()/4)))

	mlog.IL("%s", mi18n.T("読み込み開始", map[string]interface{}{"Type": "Pmx", "Path": path}))
	defer mlog.I("%s", mi18n.T("読み込み終了", map[string]interface{}{"Type": "X"}))

	// 変数を初期化
	rep.pos = 0
	rep.buffers = nil
	rep.floatSize = 0
	rep.materialTokenCount = 0

	// パスを設定
	rep.path = path

	// モデルを新規作成
	model := rep.newFunc(path)

	// ファイルを開く
	err := rep.open(path)
	if err != nil {
		mlog.E("loadData.Open error: %v", err)
		return model, err
	}

	err = rep.loadModel(model)
	if err != nil {
		mlog.E("loadData.loadModel error: %v", err)
		return model, err
	}

	rep.close()
	model.Setup()

	return model, nil
}

// 指定されたファイルオブジェクトからデータを読み込む
func (rep *XRepository) LoadByFile(file fs.File) (core.IHashModel, error) {
	// モデルを新規作成
	model := rep.newFunc("")

	// ファイルを開く
	rep.file = file
	rep.reader = bufio.NewReader(rep.file)

	err := rep.loadModel(model)
	if err != nil {
		mlog.E("loadData.LoadByFile error: %v", err)
		return model, err
	}

	rep.close()
	model.Setup()

	return model, nil
}

func (rep *XRepository) LoadName(path string) string {
	return ""
}

// loadModel はファイル形式に応じたパース処理を行います
func (rep *XRepository) loadModel(model *pmx.PmxModel) error {
	// 形式判定：ファイル先頭16バイトから "txt" / "bin" / "zip" を判別
	var format xFormat
	var err error
	if format, err = rep.detectFormat(); err != nil {
		return err
	}

	switch format {
	case xFormatText:
		// テキスト形式：既存のトークナイザーを使用
		tok := newTokenizer(rep.reader)
		var tokens []textToken
		for {
			t := tok.nextToken()
			tokens = append(tokens, t)
			if t.typ == tokEOF {
				break
			}
		}
		rep.tokens = tokens
		if err := rep.parseTextXFile(model); err != nil {
			return err
		}
	case xFormatBinary:
		// 未圧縮バイナリ形式：バイナリデータをそのまま parseBinaryXFile に渡す
		// ファイルからバイナリデータを読み込む
		if err := binary.Read(rep.reader, binary.LittleEndian, &rep.buffers); err != nil {
			return err
		}
		if err := rep.parseBinaryXFile(model); err != nil {
			return err
		}
	case xFormatCompressedBinary:
		// 圧縮バイナリ形式：ここで MSZip のデータブロックを順次解凍 → 解凍したバイト列を parseBinaryXFile に渡す
		if err := rep.parseCompressedBinaryXFile(model); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported X file format: %v", format)
	}

	// ボーン「センター」を追加し、モデルのハッシュ更新
	bone := pmx.NewBoneByName("センター")
	bone.BoneFlag = pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_CAN_TRANSLATE | pmx.BONE_FLAG_IS_VISIBLE
	bone.IsSystem = false
	model.Bones.Append(bone)
	model.UpdateHash()

	return nil
}

// detectFormat は rep.reader の先頭16バイトからファイル形式を判定し、
// rep.format に "txt", "bin", "zip" のいずれかを設定します。
// ※ テキスト形式の場合はヘッダーをそのまま残し、バイナリ系は後続処理で破棄します。
func (rep *XRepository) detectFormat() (xFormat, error) {
	// 先頭16バイトをPeek
	header, err := rep.reader.Peek(16)
	if err != nil {
		return xFormatInvalid, err
	}
	if len(header) < 12 {
		return xFormatInvalid, fmt.Errorf("Invalid X file header")
	}
	headerStr := string(header)
	// 例: "xof 0303txt 0032", "xof 0303bin 0032", "xof 0303zip 0032"
	formatIndicator := headerStr[8:12]
	switch formatIndicator {
	case "txt ":
		return xFormatText, nil
	case "tzip":
		return xFormatBinary, nil
	case "bzip":
		return xFormatCompressedBinary, nil
	default:
		return xFormatInvalid, fmt.Errorf("Unknown X file format: %s", formatIndicator)
	}
}

type tokenType int

const (
	tokIdentifier tokenType = iota
	tokNumber
	tokString
	tokLCurly
	tokRCurly
	tokSemicolon
	tokEOF
	tokAngleBracketed
)

type textToken struct {
	typ tokenType
	val string
}

type tokenizer struct {
	runes []rune
	pos   int
}

func newTokenizer(r io.Reader) *tokenizer {
	sjisReader := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	scanner := bufio.NewScanner(sjisReader)
	lines := make([]string, 0)
	for scanner.Scan() {
		txt := scanner.Text()
		lines = append(lines, txt)
	}

	return &tokenizer{runes: []rune(strings.Join(lines, "\n"))}
}

func (t *tokenizer) nextToken() textToken {
	t.skipWhitespaceAndComments()
	if t.pos >= len(t.runes) {
		return textToken{typ: tokEOF}
	}
	c := t.runes[t.pos]

	// Punctuation
	switch c {
	case '{':
		t.pos++
		return textToken{typ: tokLCurly, val: "{"}
	case '}':
		t.pos++
		return textToken{typ: tokRCurly, val: "}"}
	case ';':
		t.pos++
		return textToken{typ: tokSemicolon, val: ";"}
	case '<':
		// Parse GUID or bracketed token
		return t.readAngleBracketedToken()
	case '"':
		return t.readString()
	}

	// Number or Identifier
	if isDigit(c) || c == '-' || c == '+' || c == '.' {
		return t.readNumber()
	}

	if isIdentStart(c) {
		return t.readIdentifier()
	}

	// If nothing matches:
	t.pos++
	return t.nextToken()
}

// This new function handles tokens enclosed by < and >
func (t *tokenizer) readAngleBracketedToken() textToken {
	// consume '<'
	start := t.pos
	t.pos++
	for t.pos < len(t.runes) && t.runes[t.pos] != '>' {
		t.pos++
	}
	if t.pos >= len(t.runes) {
		panic("Unmatched '<' in input")
	}
	// now t.runes[t.pos] should be '>'
	val := string(t.runes[start : t.pos+1]) // include '>'
	t.pos++                                 // skip the closing '>'

	// We can treat this as a GUID token or similar
	return textToken{typ: tokAngleBracketed, val: val}
}

func (t *tokenizer) skipWhitespaceAndComments() {
	for t.pos < len(t.runes) {
		c := t.runes[t.pos]
		if unicode.IsSpace(c) {
			t.pos++
		} else if c == '/' && t.pos+1 < len(t.runes) && t.runes[t.pos+1] == '/' {
			// line comment
			t.pos += 2
			for t.pos < len(t.runes) && t.runes[t.pos] != '\n' {
				t.pos++
			}
		} else if c == '/' && t.pos+1 < len(t.runes) && t.runes[t.pos+1] == '*' {
			// block comment
			t.pos += 2
			for t.pos < len(t.runes)-1 {
				if t.runes[t.pos] == '*' && t.runes[t.pos+1] == '/' {
					t.pos += 2
					break
				}
				t.pos++
			}
		} else {
			break
		}
	}
}

func (t *tokenizer) readString() textToken {
	// we are at '"'
	start := t.pos
	t.pos++
	for t.pos < len(t.runes) && t.runes[t.pos] != '"' {
		t.pos++
	}
	val := string(t.runes[start+1 : t.pos])
	t.pos++ // skip closing "
	return textToken{typ: tokString, val: val}
}

func (t *tokenizer) readNumber() textToken {
	start := t.pos
	for t.pos < len(t.runes) && (isDigit(t.runes[t.pos]) || t.runes[t.pos] == '.' || t.runes[t.pos] == '-' || t.runes[t.pos] == '+') {
		t.pos++
	}
	val := string(t.runes[start:t.pos])
	return textToken{typ: tokNumber, val: val}
}

func (t *tokenizer) readIdentifier() textToken {
	start := t.pos
	for t.pos < len(t.runes) && isIdentChar(t.runes[t.pos]) {
		t.pos++
	}
	val := string(t.runes[start:t.pos])
	return textToken{typ: tokIdentifier, val: val}
}

func isDigit(c rune) bool {
	return (c >= '0' && c <= '9')
}

func isIdentStart(c rune) bool {
	return unicode.IsLetter(c) || c == '_'
}

func isIdentChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'
}

// -------------------- PARSER --------------------

func (rep *XRepository) peek() textToken {
	if rep.pos < len(rep.tokens) {
		return rep.tokens[rep.pos]
	}
	return textToken{typ: tokEOF}
}

func (rep *XRepository) next() textToken {
	t := rep.peek()
	rep.pos++
	return t
}

func (rep *XRepository) expect(typ tokenType) (textToken, error) {
	t := rep.next()
	if t.typ != typ {
		return t, fmt.Errorf("expected %v got %v (%s)\n\n%v", typ, t.typ, t.val, mstring.GetStackTrace())
	}
	return t, nil
}

// func (rep *XRepository) expectIdentifier(val string) {
// 	t := rep.next()
// 	if t.typ != tokIdentifier || t.val != val {
// 		panic(fmt.Sprintf("expected identifier '%s', got '%s'", val, t.val))
// 	}
// }

func (rep *XRepository) parseTextXFile(model *pmx.PmxModel) error {
	// Parse until EOF
	for rep.peek().typ != tokEOF {
		t := rep.peek()

		if t.typ == tokIdentifier && t.val == "template" {
			// We are encountering a template definition
			rep.next() // consume 'template'
			if err := rep.parseTextTemplateDefinition(); err != nil {
				return err
			}
		} else if t.typ == tokIdentifier && t.val == "Header" {
			rep.next()
			if err := rep.parseTextHeader(model); err != nil {
				return err
			}
		} else if t.typ == tokIdentifier && t.val == "Mesh" {
			rep.next()
			if err := rep.parseTextMesh(model); err != nil {
				return err
			}
		} else if t.typ == tokIdentifier {
			// Look ahead to see if next token is '{'
			rep.next() // consume the identifier
			if rep.peek().typ == tokLCurly {
				// Known template name followed by '{' means instance block
				// Parse as an instance of that template
				if err := rep.parseTextTemplateInstance(); err != nil {
					return err
				}
			} else {
				// If not '{', it's not a valid instance block,
				// possibly skip or handle error.
				if err := rep.skipTextUnknownTemplate(); err != nil {
					return err
				}
			}
		} else {
			// Not template keyword or known template name, skip or consume token
			rep.next()
		}
	}
	return nil
}

// parseTextTemplateDefinition parses a template definition block like:
//
//	template Mesh {
//	   <...GUID...>
//	   DWORD nVertices;
//	   ...
//	}
func (rep *XRepository) parseTextTemplateDefinition() error {
	// expect template name
	nameTok, err := rep.expect(tokIdentifier)
	if err != nil {
		return err
	}

	_ = nameTok.val

	// Expect '{'
	if _, err := rep.expect(tokLCurly); err != nil {
		return err
	}

	// Expect GUID line: <...>
	guidTok := rep.next()
	if guidTok.typ != tokAngleBracketed || !strings.HasPrefix(guidTok.val, "<") || !strings.HasSuffix(guidTok.val, ">") {
		return fmt.Errorf("expected GUID in angle brackets in template definition\n\n%v", mstring.GetStackTrace())
	}

	// Now parse the fields until '}' is found
	// Typically these are lines like: DWORD something; array Vector ...;
	// For simplicity, we can just skip until '}' since we are only defining schema.
	// In a real implementation, you might store the schema info.

	// Skip until matching '}'
	braceCount := 1
	for braceCount > 0 {
		tok := rep.next()
		switch tok.typ {
		case tokLCurly:
			braceCount++
		case tokRCurly:
			braceCount--
		case tokEOF:
			return fmt.Errorf("unexpected EOF in template instance\n\n%v", mstring.GetStackTrace())
		}
	}

	// At this point, we have the template defined. In a real parser,
	// you would store the template definition schema somewhere.

	return nil
}

// parseTextTemplateInstance parses an instance of a previously defined template:
// e.g. Mesh { ... actual data ... }
func (rep *XRepository) parseTextTemplateInstance() error {
	// We already consumed the templateName and peeked '{'
	if _, err := rep.expect(tokLCurly); err != nil {
		return err
	}

	// Here you would parse the template instance data according to the known schema
	// For demonstration, we'll just skip until the closing '}':
	braceCount := 1
	for braceCount > 0 {
		tok := rep.next()
		switch tok.typ {
		case tokLCurly:
			braceCount++
		case tokRCurly:
			braceCount--
		case tokEOF:
			return fmt.Errorf("unexpected EOF in template instance\n\n%v", mstring.GetStackTrace())
		}
	}

	// After this, the instance block is fully parsed.
	return nil
}

// // This function just shows how you might skip unknown templates if encountered
// func (rep *XRepository) skipUnknownTemplate() {
// 	// If current token is not '{', just return
// 	if rep.peek().typ != tokLCurly {
// 		return
// 	}

// 	rep.next() // consume '{'
// 	braceCount := 1
// 	for braceCount > 0 {
// 		tok := rep.next()
// 		switch tok.typ {
// 		case tokLCurly:
// 			braceCount++
// 		case tokRCurly:
// 			braceCount--
// 		case tokEOF:
// 			// Reached EOF without closing the template properly
// 			return
// 		}
// 	}
// }

func (rep *XRepository) parseTextHeader(model *pmx.PmxModel) error {
	if _, err := rep.expect(tokLCurly); err != nil {
		return err
	}
	majorVersion, err := rep.parseNumberAsFloat()
	if err != nil {
		return err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}
	minorVersion, err := rep.parseNumberAsFloat()
	if err != nil {
		return err
	}

	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}
	flags, err := rep.parseNumberAsFloat()
	if err != nil {
		return err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}
	if _, err := rep.expect(tokRCurly); err != nil {
		return err
	}

	model.Comment = fmt.Sprintf("X File Version %.0f.%.0f, flags: %.0f", majorVersion, minorVersion, flags)

	return nil
}

func (rep *XRepository) parseNumberAsFloat() (float64, error) {
	t, err := rep.expect(tokNumber)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(t.val, 32)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func (rep *XRepository) parseNumberAsInt() (int, error) {
	t, err := rep.expect(tokNumber)
	if err != nil {
		return 0, err
	}
	val, err := strconv.ParseUint(t.val, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func (rep *XRepository) parseString() (string, error) {
	t := rep.next()
	if t.typ == tokString {
		return t.val, nil
	}
	// In .x files, strings can sometimes appear without quotes as identifiers.
	if t.typ == tokIdentifier {
		return t.val, nil
	}
	return "", fmt.Errorf("expected string\n\n%v", mstring.GetStackTrace())
}

func (rep *XRepository) parseVector() (*mmath.MVec3, error) {
	x, err := rep.parseNumberAsFloat()
	if err != nil {
		return nil, err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return nil, err
	}
	y, err := rep.parseNumberAsFloat()
	if err != nil {
		return nil, err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return nil, err
	}
	z, err := rep.parseNumberAsFloat()
	if err != nil {
		return nil, err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return nil, err
	}

	return &mmath.MVec3{X: x, Y: y, Z: z}, nil
}

// func (rep *XRepository) parseCoords2d() *mmath.MVec2 {
// 	rep.expect(tokLCurly)
// 	u := rep.parseNumberAsFloat()
// 	rep.expect(tokSemicolon)
// 	v := rep.parseNumberAsFloat()
// 	rep.expect(tokSemicolon)
// 	rep.expect(tokRCurly)
// 	return &mmath.MVec2{X: u, Y: v}
// }

func (rep *XRepository) parseColorRGBA() (*mmath.MVec4, error) {
	r, err := rep.parseNumberAsFloat()
	if err != nil {
		return nil, err
	}

	if _, err := rep.expect(tokSemicolon); err != nil {
		return nil, err
	}
	g, err := rep.parseNumberAsFloat()
	if err != nil {
		return nil, err
	}

	if _, err := rep.expect(tokSemicolon); err != nil {
		return nil, err
	}
	b, err := rep.parseNumberAsFloat()
	if err != nil {
		return nil, err
	}

	if _, err := rep.expect(tokSemicolon); err != nil {
		return nil, err
	}
	a, err := rep.parseNumberAsFloat()
	if err != nil {
		return nil, err
	}

	if _, err := rep.expect(tokSemicolon); err != nil {
		return nil, err
	}
	return &mmath.MVec4{X: r, Y: g, Z: b, W: a}, nil
}

func (rep *XRepository) parseColorRGB() (*mmath.MVec3, error) {
	r, err := rep.parseNumberAsFloat()
	if err != nil {
		return nil, err
	}

	if _, err := rep.expect(tokSemicolon); err != nil {
		return nil, err
	}
	g, err := rep.parseNumberAsFloat()
	if err != nil {
		return nil, err
	}

	if _, err := rep.expect(tokSemicolon); err != nil {
		return nil, err
	}
	b, err := rep.parseNumberAsFloat()
	if err != nil {
		return nil, err
	}

	if _, err := rep.expect(tokSemicolon); err != nil {
		return nil, err
	}
	return &mmath.MVec3{X: r, Y: g, Z: b}, nil
}

func (rep *XRepository) parseMaterialText(model *pmx.PmxModel) error {
	var err error

	// Material定義の中に文字列（材質名）がある場合があるが、スルー
	if _, err := rep.parseString(); err == nil {
		if _, err := rep.expect(tokLCurly); err != nil {
			return err
		}
	}

	mat := pmx.NewMaterial()

	mat.Diffuse, err = rep.parseColorRGBA()
	if err != nil {
		return err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}

	// power
	power, err := rep.parseNumberAsFloat()
	if err != nil {
		return err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}

	specular, err := rep.parseColorRGB()
	if err != nil {
		return err
	}
	mat.Specular = &mmath.MVec4{
		X: specular.X,
		Y: specular.Y,
		Z: specular.Z,
		W: power,
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}

	mat.Ambient, err = rep.parseColorRGB()
	if err != nil {
		return err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}

	// Optional TextureFilename
	for rep.peek().typ == tokIdentifier && rep.peek().val == "TextureFilename" {
		rep.next()
		texturePath, spherePath, err := rep.parseTextureFilename()
		if err != nil {
			return err
		}

		tex := pmx.NewTexture()
		tex.SetName(texturePath)
		model.Textures.Append(tex)
		mat.TextureIndex = tex.Index()

		if spherePath != "" {
			sphere := pmx.NewTexture()
			sphere.SetName(spherePath)
			model.Textures.Append(sphere)
			mat.SphereTextureIndex = sphere.Index()
		}
	}
	if _, err := rep.expect(tokRCurly); err != nil {
		return err
	}

	mat.SetName(fmt.Sprintf("材質%02d", model.Materials.Length()+1))
	mat.Edge.W = 1.0
	mat.EdgeSize = 10.0
	if mat.TextureIndex >= 0 && mat.SphereTextureIndex < 0 {
		// テクスチャがあり、スフィアがない場合、スフィアモードを無効にする
		mat.SphereMode = pmx.SPHERE_MODE_INVALID
	} else {
		// テクスチャがない、もしくはスフィアがある場合、スフィアモードを乗算にする
		mat.SphereMode = pmx.SPHERE_MODE_MULTIPLICATION
	}
	model.Materials.Append(mat)

	return nil
}

func (rep *XRepository) parseTextureFilename() (texturePath, spherePath string, err error) {
	if _, err := rep.expect(tokLCurly); err != nil {
		return "", "", err
	}
	tf, err := rep.parseString()
	if err != nil {
		return "", "", err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return "", "", err
	}
	if _, err := rep.expect(tokRCurly); err != nil {
		return "", "", err
	}
	// * で区切る
	tfs := strings.Split(tf, "*")
	if len(tfs) > 1 {
		return tfs[0], tfs[1], nil
	}
	return tfs[0], "", nil
}

func (rep *XRepository) parseMeshFace() (fs []*pmx.Face, err error) {
	count, err := rep.parseNumberAsInt()
	if err != nil {
		return nil, err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return nil, err
	}
	vertexIndexes := make([]int, 0, 4)
	for i := 0; i < count; i++ {
		idx, err := rep.parseNumberAsInt()
		if err != nil {
			return nil, err
		}
		vertexIndexes = append(vertexIndexes, idx)
	}
	// これはなかったりする
	rep.expect(tokSemicolon)

	// 4つの頂点を持つ場合、三角面2つに分解する
	if count == 4 {
		f1 := pmx.NewFace()
		f1.VertexIndexes[0] = vertexIndexes[0]
		f1.VertexIndexes[1] = vertexIndexes[1]
		f1.VertexIndexes[2] = vertexIndexes[2]

		f2 := pmx.NewFace()
		f2.VertexIndexes[0] = vertexIndexes[0]
		f2.VertexIndexes[1] = vertexIndexes[2]
		f2.VertexIndexes[2] = vertexIndexes[3]

		fs = append(fs, f1, f2)
	} else {
		f := pmx.NewFace()
		f.VertexIndexes[0] = vertexIndexes[0]
		f.VertexIndexes[1] = vertexIndexes[1]
		f.VertexIndexes[2] = vertexIndexes[2]
		fs = append(fs, f)
	}

	return fs, nil
}

func (rep *XRepository) parseMeshTextureCoords(model *pmx.PmxModel) error {
	if _, err := rep.expect(tokLCurly); err != nil {
		return err
	}
	count, err := rep.parseNumberAsInt()
	if err != nil {
		return err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}
	for i := 0; i < count; i++ {
		u, err := rep.parseNumberAsFloat()
		if err != nil {
			return err
		}
		if _, err := rep.expect(tokSemicolon); err != nil {
			return err
		}
		v, err := rep.parseNumberAsFloat()
		if err != nil {
			return err
		}
		if _, err := rep.expect(tokSemicolon); err != nil {
			return err
		}
		vt, err := model.Vertices.Get(i)
		if err != nil {
			return err
		}
		vt.Uv = &mmath.MVec2{X: u, Y: v}
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}
	if _, err := rep.expect(tokRCurly); err != nil {
		return err
	}

	return nil
}

func (rep *XRepository) parseMeshMaterialList(
	model *pmx.PmxModel, facesList [][]*pmx.Face,
) error {
	if _, err := rep.expect(tokLCurly); err != nil {
		return err
	}
	nMat, err := rep.parseNumberAsInt()
	if err != nil {
		return err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}
	nFaceIdx, err := rep.parseNumberAsInt()
	if err != nil {
		return err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}

	facesByMaterials := make(map[int][][]*pmx.Face)
	for i := 0; i < nMat; i++ {
		facesByMaterials[i] = make([][]*pmx.Face, 0, nFaceIdx)
	}

	// faceIndexes
	faceIndexesByMaterials := make(map[int][]int)
	for i := 0; i < nMat; i++ {
		faceIndexesByMaterials[i] = make([]int, 0, nFaceIdx)
	}

	for i := 0; i < nFaceIdx; i++ {
		matIdx, err := rep.parseNumberAsInt()
		if err != nil {
			return err
		}
		facesByMaterials[matIdx] = append(facesByMaterials[matIdx], facesList[i])
		faceIndexesByMaterials[matIdx] = append(faceIndexesByMaterials[matIdx], i)
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}
	if rep.peek().typ == tokSemicolon {
		rep.expect(tokSemicolon)
	}

	faceMap := make(map[int][]int)

	// Materials
	// After this, we might have that many Material references or full Material templates
	for i := 0; i < nMat; i++ {
		if rep.peek().typ == tokIdentifier && rep.peek().val == "Material" {
			rep.next()
			rep.parseMaterialText(model)
		} else {
			// could be a reference (string) or skip
			// Just skip if unknown
			rep.skipTextUnknownTemplate()
			continue
		}

		// 面を割り当てる
		for j, fs := range facesByMaterials[i] {
			fis := make([]int, 0, len(fs))
			for _, f := range fs {
				f.SetIndex(model.Faces.Length())
				model.Faces.Append(f)
				fis = append(fis, f.Index())
				m, err := model.Materials.Get(i)
				if err != nil {
					return err
				}
				m.VerticesCount += 3
			}
			faceMap[faceIndexesByMaterials[i][j]] = fis
		}
	}

	if _, err := rep.expect(tokRCurly); err != nil {
		return err
	}

	return nil
}

func (rep *XRepository) parseTextMesh(model *pmx.PmxModel) error {

	// 定義の中に文字列（メッシュ名）がある場合があるが、スルー
	if _, err := rep.parseString(); err == nil {
		if _, err := rep.expect(tokLCurly); err != nil {
			return err
		}
	}

	nVertices, err := rep.parseNumberAsInt()
	if err != nil {
		return err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}
	for i := 0; i < nVertices; i++ {
		v := pmx.NewVertex()
		v.Position, err = rep.parseVector()
		if err != nil {
			return err
		}
		// 頂点位置を10倍にする
		v.Position.MulScalar(10)
		// BDEF1
		v.Deform = pmx.NewBdef1(0)
		// エッジ倍率1
		v.EdgeFactor = 1
		// 法線
		v.Normal = &mmath.MVec3{X: 0, Y: 1, Z: 0}
		model.Vertices.Append(v)
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}

	nFaces, err := rep.parseNumberAsInt()
	if err != nil {
		return err
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}
	facesList := make([][]*pmx.Face, 0, nFaces)
	faceTotalCount := 0
	for range nFaces {
		fs, err := rep.parseMeshFace()
		if err != nil {
			return err
		}
		facesList = append(facesList, fs)
		faceTotalCount += len(fs)
	}
	if _, err := rep.expect(tokSemicolon); err != nil {
		return err
	}

	// Optional sub-templates
	for rep.peek().typ == tokIdentifier {
		switch rep.peek().val {
		// case "MeshNormals":
		// 	rep.next()
		// 	rep.parseMeshNormals(model)
		case "MeshMaterialList":
			rep.next()
			if err := rep.parseMeshMaterialList(model, facesList); err != nil {
				return err
			}
		case "MeshTextureCoords":
			rep.next()
			rep.parseMeshTextureCoords(model)
		default:
			// skip unknown sub-template
			rep.next()
			rep.skipTextUnknownTemplate()
		}
	}

	if _, err := rep.expect(tokRCurly); err != nil {
		return err
	}

	return nil
}

func (rep *XRepository) skipTextUnknownTemplate() error {
	// すでに "templateName" のような識別子を読んだあとで呼び出されることを想定
	// 次のトークンは "{" のはず
	t := rep.next()
	if t.typ != tokLCurly {
		// もし "{" がなければスキップ対象はないので return
		return nil
	}

	braceCount := 1
	for braceCount > 0 {
		tok := rep.next()
		switch tok.typ {
		case tokLCurly:
			braceCount++
		case tokRCurly:
			braceCount--
		case tokEOF:
			// ファイル終端まで来てしまった場合、テンプレートが不正である可能性あり
			return fmt.Errorf("unexpected EOF in template instance\n\n%v", mstring.GetStackTrace())
		}
	}
	// braceCount が0になったので対応する "}" に到達し、スキップ完了

	return nil
}
