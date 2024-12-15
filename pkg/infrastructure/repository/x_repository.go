package repository

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type XRepository struct {
	*baseRepository[*pmx.PmxModel]
	tokens []token
	pos    int
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
	if isExist, err := mutils.ExistsFile(path); err != nil || !isExist {
		return false, fmt.Errorf(mi18n.T("ファイル存在エラー", map[string]interface{}{"Path": path}))
	}

	_, _, ext := mutils.SplitPath(path)
	if strings.ToLower(ext) != ".x" {
		return false, fmt.Errorf(mi18n.T("拡張子エラー", map[string]interface{}{"Path": path, "Ext": ".x"}))
	}

	return true, nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *XRepository) Load(path string) (core.IHashModel, error) {
	runtime.GOMAXPROCS(int(runtime.NumCPU()))
	defer runtime.GOMAXPROCS(max(1, int(runtime.NumCPU()/4)))

	mlog.IL(mi18n.T("読み込み開始", map[string]interface{}{"Type": "Pmx", "Path": path}))
	defer mlog.I(mi18n.T("読み込み終了", map[string]interface{}{"Type": "X"}))

	// モデルを新規作成
	model := rep.newFunc(path)

	// ファイルを開く
	err := rep.open(path)
	if err != nil {
		mlog.E("ReadByFilepath.Open error: %v", err)
		return model, err
	}

	err = rep.loadModel(model)
	if err != nil {
		mlog.E("ReadByFilepath.loadData error: %v", err)
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
		mlog.E("ReadByFilepath.loadData error: %v", err)
		return model, err
	}

	rep.close()
	model.Setup()

	return model, nil
}

func (rep *XRepository) LoadName(path string) string {
	return ""
}

func (rep *XRepository) loadModel(model *pmx.PmxModel) error {
	tok := newTokenizer(rep.reader)
	var tokens []token
	for {
		t := tok.nextToken()
		tokens = append(tokens, t)
		if t.typ == tokEOF {
			break
		}
	}
	rep.tokens = tokens

	err := rep.parseXFile(model)
	if err != nil {
		mlog.E("loadData.parseXFile error: %v", err)
		return err
	}

	// ボーンを一本追加
	bone := pmx.NewBoneByName("センター")
	bone.BoneFlag = pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_CAN_TRANSLATE | pmx.BONE_FLAG_IS_VISIBLE
	bone.IsSystem = false
	model.Bones.Append(bone)

	// モデルのハッシュを更新
	model.UpdateHash()

	return nil
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

type token struct {
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

func (t *tokenizer) nextToken() token {
	t.skipWhitespaceAndComments()
	if t.pos >= len(t.runes) {
		return token{typ: tokEOF}
	}
	c := t.runes[t.pos]

	// Punctuation
	switch c {
	case '{':
		t.pos++
		return token{typ: tokLCurly, val: "{"}
	case '}':
		t.pos++
		return token{typ: tokRCurly, val: "}"}
	case ';':
		t.pos++
		return token{typ: tokSemicolon, val: ";"}
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
func (t *tokenizer) readAngleBracketedToken() token {
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
	return token{typ: tokAngleBracketed, val: val}
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

func (t *tokenizer) readString() token {
	// we are at '"'
	start := t.pos
	t.pos++
	for t.pos < len(t.runes) && t.runes[t.pos] != '"' {
		t.pos++
	}
	val := string(t.runes[start+1 : t.pos])
	t.pos++ // skip closing "
	return token{typ: tokString, val: val}
}

func (t *tokenizer) readNumber() token {
	start := t.pos
	for t.pos < len(t.runes) && (isDigit(t.runes[t.pos]) || t.runes[t.pos] == '.' || t.runes[t.pos] == '-' || t.runes[t.pos] == '+') {
		t.pos++
	}
	val := string(t.runes[start:t.pos])
	return token{typ: tokNumber, val: val}
}

func (t *tokenizer) readIdentifier() token {
	start := t.pos
	for t.pos < len(t.runes) && isIdentChar(t.runes[t.pos]) {
		t.pos++
	}
	val := string(t.runes[start:t.pos])
	return token{typ: tokIdentifier, val: val}
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

func (rep *XRepository) peek() token {
	if rep.pos < len(rep.tokens) {
		return rep.tokens[rep.pos]
	}
	return token{typ: tokEOF}
}

func (rep *XRepository) next() token {
	t := rep.peek()
	rep.pos++
	return t
}

func (rep *XRepository) expect(typ tokenType) token {
	t := rep.next()
	if t.typ != typ {
		panic(fmt.Sprintf("expected %v got %v (%s)", typ, t.typ, t.val))
	}
	return t
}

// func (rep *XRepository) expectIdentifier(val string) {
// 	t := rep.next()
// 	if t.typ != tokIdentifier || t.val != val {
// 		panic(fmt.Sprintf("expected identifier '%s', got '%s'", val, t.val))
// 	}
// }

func (rep *XRepository) parseXFile(model *pmx.PmxModel) error {
	// Parse until EOF
	for rep.peek().typ != tokEOF {
		t := rep.peek()

		if t.typ == tokIdentifier && t.val == "template" {
			// We are encountering a template definition
			rep.next() // consume 'template'
			rep.parseTemplateDefinition()
		} else if t.typ == tokIdentifier && t.val == "Header" {
			rep.next()
			rep.parseHeader(model)
		} else if t.typ == tokIdentifier && t.val == "Mesh" {
			rep.next()
			rep.parseMesh(model)
		} else if t.typ == tokIdentifier {
			// Could be a template instance block
			templateName := t.val
			// Look ahead to see if next token is '{'
			rep.next() // consume the identifier
			if rep.peek().typ == tokLCurly {
				// Known template name followed by '{' means instance block
				// Parse as an instance of that template
				rep.parseTemplateInstance(templateName)
			} else {
				// If not '{', it's not a valid instance block,
				// possibly skip or handle error.
				rep.skipUnknownTemplate()
			}
		} else {
			// Not template keyword or known template name, skip or consume token
			rep.next()
		}
	}
	return nil
}

// parseTemplateDefinition parses a template definition block like:
//
//	template Mesh {
//	   <...GUID...>
//	   DWORD nVertices;
//	   ...
//	}
func (rep *XRepository) parseTemplateDefinition() {
	// expect template name
	nameTok := rep.expect(tokIdentifier)
	_ = nameTok.val

	// Expect '{'
	rep.expect(tokLCurly)

	// Expect GUID line: <...>
	guidTok := rep.next()
	if guidTok.typ != tokAngleBracketed || !strings.HasPrefix(guidTok.val, "<") || !strings.HasSuffix(guidTok.val, ">") {
		panic("Expected GUID in angle brackets in template definition")
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
			panic("Unexpected EOF in template definition")
		}
	}

	// At this point, we have the template defined. In a real parser,
	// you would store the template definition schema somewhere.
}

// parseTemplateInstance parses an instance of a previously defined template:
// e.g. Mesh { ... actual data ... }
func (rep *XRepository) parseTemplateInstance(templateName string) {
	// We already consumed the templateName and peeked '{'
	rep.expect(tokLCurly)

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
			panic("Unexpected EOF in template instance")
		}
	}

	// After this, the instance block is fully parsed.
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

func (rep *XRepository) parseHeader(model *pmx.PmxModel) {
	rep.expect(tokLCurly)
	majorVersion := uint16(rep.parseNumberAsFloat())
	rep.expect(tokSemicolon)
	minorVersion := uint16(rep.parseNumberAsFloat())
	rep.expect(tokSemicolon)
	flags := uint32(rep.parseNumberAsFloat())
	rep.expect(tokSemicolon)
	rep.expect(tokRCurly)

	model.Comment = fmt.Sprintf("X File Version %d.%d, flags: %d", majorVersion, minorVersion, flags)
}

func (rep *XRepository) parseNumberAsFloat() float64 {
	t := rep.expect(tokNumber)
	f, err := strconv.ParseFloat(t.val, 32)
	if err != nil {
		panic(err)
	}
	return f
}

func (rep *XRepository) parseNumberAsInt() int {
	t := rep.expect(tokNumber)
	val, err := strconv.ParseUint(t.val, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(val)
}

func (rep *XRepository) parseString() string {
	t := rep.next()
	if t.typ == tokString {
		return t.val
	}
	// In .x files, strings can sometimes appear without quotes as identifiers.
	if t.typ == tokIdentifier {
		return t.val
	}
	panic("expected string")
}

func (rep *XRepository) parseVector() (v *mmath.MVec3) {
	x := rep.parseNumberAsFloat()
	rep.expect(tokSemicolon)
	y := rep.parseNumberAsFloat()
	rep.expect(tokSemicolon)
	z := rep.parseNumberAsFloat()
	rep.expect(tokSemicolon)

	return &mmath.MVec3{X: x, Y: y, Z: z}
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

func (rep *XRepository) parseColorRGBA() *mmath.MVec4 {
	r := rep.parseNumberAsFloat()
	rep.expect(tokSemicolon)
	g := rep.parseNumberAsFloat()
	rep.expect(tokSemicolon)
	b := rep.parseNumberAsFloat()
	rep.expect(tokSemicolon)
	a := rep.parseNumberAsFloat()
	rep.expect(tokSemicolon)
	return &mmath.MVec4{X: r, Y: g, Z: b, W: a}
}

func (rep *XRepository) parseColorRGB() *mmath.MVec3 {
	r := rep.parseNumberAsFloat()
	rep.expect(tokSemicolon)
	g := rep.parseNumberAsFloat()
	rep.expect(tokSemicolon)
	b := rep.parseNumberAsFloat()
	rep.expect(tokSemicolon)
	return &mmath.MVec3{X: r, Y: g, Z: b}
}

func (rep *XRepository) parseMaterial(model *pmx.PmxModel) {
	rep.expect(tokLCurly)
	mat := pmx.NewMaterial()

	mat.Diffuse = rep.parseColorRGBA()
	rep.expect(tokSemicolon)

	// power
	power := rep.parseNumberAsFloat()
	rep.expect(tokSemicolon)

	specular := rep.parseColorRGB()
	mat.Specular = &mmath.MVec4{
		X: specular.X,
		Y: specular.Y,
		Z: specular.Z,
		W: power,
	}
	rep.expect(tokSemicolon)

	mat.Ambient = rep.parseColorRGB()
	rep.expect(tokSemicolon)

	// Optional TextureFilename
	for rep.peek().typ == tokIdentifier && rep.peek().val == "TextureFilename" {
		rep.next()
		tf := rep.parseTextureFilename()

		tex := pmx.NewTexture()
		tex.SetName(tf)
		model.Textures.Append(tex)
		mat.TextureIndex = tex.Index()
	}
	rep.expect(tokRCurly)

	mat.SetName(fmt.Sprintf("材質%02d", model.Materials.Len()+1))
	mat.Edge.W = 1.0
	mat.EdgeSize = 10.0
	if mat.TextureIndex >= 0 {
		// テクスチャがある場合、スフィアモードを無効にする
		mat.SphereMode = pmx.SPHERE_MODE_INVALID
	} else {
		// テクスチャがない場合、スフィアモードを乗算にする
		mat.SphereMode = pmx.SPHERE_MODE_MULTIPLICATION
	}
	model.Materials.Append(mat)
}

func (rep *XRepository) parseTextureFilename() string {
	rep.expect(tokLCurly)
	tf := rep.parseString()
	rep.expect(tokSemicolon)
	rep.expect(tokRCurly)
	return tf
}

func (rep *XRepository) parseMeshFace() (fs []*pmx.Face) {
	count := rep.parseNumberAsInt()
	rep.expect(tokSemicolon)
	vertexIndexes := make([]int, 0, 4)
	for i := 0; i < count; i++ {
		idx := rep.parseNumberAsInt()
		vertexIndexes = append(vertexIndexes, idx)
	}
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

	return fs
}

func (rep *XRepository) parseMeshTextureCoords(model *pmx.PmxModel) {
	rep.expect(tokLCurly)
	count := rep.parseNumberAsInt()
	rep.expect(tokSemicolon)
	for i := 0; i < count; i++ {
		u := rep.parseNumberAsFloat()
		rep.expect(tokSemicolon)
		v := rep.parseNumberAsFloat()
		rep.expect(tokSemicolon)
		model.Vertices.Get(i).Uv = &mmath.MVec2{X: u, Y: v}
	}
	rep.expect(tokSemicolon)
	rep.expect(tokRCurly)
}

func (rep *XRepository) parseMeshMaterialList(model *pmx.PmxModel, facesList [][]*pmx.Face) (faceMap map[int][]int) {
	rep.expect(tokLCurly)
	nMat := rep.parseNumberAsInt()
	rep.expect(tokSemicolon)
	nFaceIdx := rep.parseNumberAsInt()
	rep.expect(tokSemicolon)

	facesByMaterials := make(map[int][][]*pmx.Face)
	for i := 0; i < nMat; i++ {
		facesByMaterials[i] = make([][]*pmx.Face, 0, nFaceIdx)
	}

	// faceIndices
	faceIndicesByMaterials := make(map[int][]int)
	for i := 0; i < nMat; i++ {
		faceIndicesByMaterials[i] = make([]int, 0, nFaceIdx)
	}

	for i := 0; i < nFaceIdx; i++ {
		matIdx := rep.parseNumberAsInt()
		facesByMaterials[matIdx] = append(facesByMaterials[matIdx], facesList[i])
		faceIndicesByMaterials[matIdx] = append(faceIndicesByMaterials[matIdx], i)
	}
	rep.expect(tokSemicolon)
	rep.expect(tokSemicolon)

	faceMap = make(map[int][]int)

	// Materials
	// After this, we might have that many Material references or full Material templates
	for i := 0; i < nMat; i++ {
		if rep.peek().typ == tokIdentifier && rep.peek().val == "Material" {
			rep.next()
			rep.parseMaterial(model)
		} else {
			// could be a reference (string) or skip
			// Just skip if unknown
			rep.skipUnknownTemplate()
		}

		// 面を割り当てる
		for j, fs := range facesByMaterials[i] {
			fis := make([]int, 0, len(fs))
			for _, f := range fs {
				f.SetIndex(model.Faces.Len())
				model.Faces.Append(f)
				fis = append(fis, f.Index())
				model.Materials.Get(i).VerticesCount += 3
			}
			faceMap[faceIndicesByMaterials[i][j]] = fis
		}
	}

	rep.expect(tokRCurly)
	return faceMap
}

func (rep *XRepository) parseMeshNormals(model *pmx.PmxModel) {
	rep.expect(tokLCurly)
	nNorm := rep.parseNumberAsInt()
	rep.expect(tokSemicolon)
	for i := 0; i < nNorm; i++ {
		model.Vertices.Get(i).Normal = rep.parseVector()
	}
	nFaceNorm := rep.parseNumberAsInt()
	rep.expect(tokSemicolon)
	for i := 0; i < nFaceNorm; i++ {
		_ = rep.parseMeshFace()
	}
	rep.expect(tokRCurly)
}

func (rep *XRepository) parseMesh(model *pmx.PmxModel) {
	rep.expect(tokLCurly)
	nVertices := rep.parseNumberAsInt()
	rep.expect(tokSemicolon)
	for i := 0; i < nVertices; i++ {
		v := pmx.NewVertex()
		v.Position = rep.parseVector()
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
	rep.expect(tokSemicolon)

	nFaces := rep.parseNumberAsInt()
	rep.expect(tokSemicolon)
	facesList := make([][]*pmx.Face, 0, nFaces)
	faceTotalCount := 0
	for i := 0; i < nFaces; i++ {
		fs := rep.parseMeshFace()
		facesList = append(facesList, fs)
		faceTotalCount += len(fs)
	}
	rep.expect(tokSemicolon)

	// Optional sub-templates
	for rep.peek().typ == tokIdentifier {
		switch rep.peek().val {
		// case "MeshNormals":
		// 	rep.next()
		// 	rep.parseMeshNormals(model)
		case "MeshMaterialList":
			rep.next()
			rep.parseMeshMaterialList(model, facesList)
		case "MeshTextureCoords":
			rep.next()
			rep.parseMeshTextureCoords(model)
		default:
			// skip unknown sub-template
			rep.next()
			rep.skipUnknownTemplate()
		}
	}

	rep.expect(tokRCurly)
}

func (rep *XRepository) skipUnknownTemplate() {
	// すでに "templateName" のような識別子を読んだあとで呼び出されることを想定
	// 次のトークンは "{" のはず
	t := rep.next()
	if t.typ != tokLCurly {
		// もし "{" がなければスキップ対象はないので return
		return
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
			// 適宜エラー処理
			return
		}
	}
	// braceCount が0になったので対応する "}" に到達し、スキップ完了
}
