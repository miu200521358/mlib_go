// 指示: miu200521358
package x

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"gonum.org/v1/gonum/spatial/r3"
)

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

// newTokenizer はXテキストのトークナイザーを生成する。
func newTokenizer(r io.Reader) (*tokenizer, error) {
	sjisReader := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	scanner := bufio.NewScanner(sjisReader)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &tokenizer{runes: []rune(strings.Join(lines, "\n"))}, nil
}

// nextToken は次のトークンを返す。
func (t *tokenizer) nextToken() textToken {
	t.skipWhitespaceAndComments()
	if t.pos >= len(t.runes) {
		return textToken{typ: tokEOF}
	}
	c := t.runes[t.pos]
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
		return t.readAngleBracketedToken()
	case '"':
		return t.readString()
	}
	if isDigit(c) || c == '-' || c == '+' || c == '.' {
		return t.readNumber()
	}
	if isIdentStart(c) {
		return t.readIdentifier()
	}
	t.pos++
	return t.nextToken()
}

// readAngleBracketedToken は<...>をトークン化する。
func (t *tokenizer) readAngleBracketedToken() textToken {
	start := t.pos
	t.pos++
	for t.pos < len(t.runes) && t.runes[t.pos] != '>' {
		t.pos++
	}
	if t.pos >= len(t.runes) {
		return textToken{typ: tokEOF}
	}
	val := string(t.runes[start : t.pos+1])
	t.pos++
	return textToken{typ: tokAngleBracketed, val: val}
}

// skipWhitespaceAndComments は空白とコメントを読み飛ばす。
func (t *tokenizer) skipWhitespaceAndComments() {
	for t.pos < len(t.runes) {
		c := t.runes[t.pos]
		switch {
		case unicode.IsSpace(c):
			t.pos++
		case c == '/' && t.pos+1 < len(t.runes) && t.runes[t.pos+1] == '/':
			t.pos += 2
			for t.pos < len(t.runes) && t.runes[t.pos] != '\n' {
				t.pos++
			}
		case c == '/' && t.pos+1 < len(t.runes) && t.runes[t.pos+1] == '*':
			t.pos += 2
			for t.pos < len(t.runes)-1 {
				if t.runes[t.pos] == '*' && t.runes[t.pos+1] == '/' {
					t.pos += 2
					break
				}
				t.pos++
			}
		default:
			return
		}
	}
}

// readString は文字列トークンを読み取る。
func (t *tokenizer) readString() textToken {
	start := t.pos
	t.pos++
	for t.pos < len(t.runes) && t.runes[t.pos] != '"' {
		t.pos++
	}
	val := string(t.runes[start+1 : t.pos])
	if t.pos < len(t.runes) {
		t.pos++
	}
	return textToken{typ: tokString, val: val}
}

// readNumber は数値トークンを読み取る。
func (t *tokenizer) readNumber() textToken {
	start := t.pos
	for t.pos < len(t.runes) {
		r := t.runes[t.pos]
		if !isDigit(r) && r != '.' && r != '-' && r != '+' {
			break
		}
		t.pos++
	}
	return textToken{typ: tokNumber, val: string(t.runes[start:t.pos])}
}

// readIdentifier は識別子トークンを読み取る。
func (t *tokenizer) readIdentifier() textToken {
	start := t.pos
	for t.pos < len(t.runes) && isIdentChar(t.runes[t.pos]) {
		t.pos++
	}
	return textToken{typ: tokIdentifier, val: string(t.runes[start:t.pos])}
}

// isDigit は数字か判定する。
func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// isIdentStart は識別子の先頭か判定する。
func isIdentStart(c rune) bool {
	return unicode.IsLetter(c) || c == '_'
}

// isIdentChar は識別子文字か判定する。
func isIdentChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'
}

type textParser struct {
	tokens  []textToken
	pos     int
	model   *model.PmxModel
	meshCtx *meshContext
}

// newTextParser はXテキストのパーサーを生成する。
func newTextParser(r io.Reader) (*textParser, error) {
	tok, err := newTokenizer(r)
	if err != nil {
		return nil, err
	}
	tokens := make([]textToken, 0)
	for {
		t := tok.nextToken()
		tokens = append(tokens, t)
		if t.typ == tokEOF {
			break
		}
	}
	return &textParser{tokens: tokens}, nil
}

// Parse はXテキストを解析する。
func (p *textParser) Parse(modelData *model.PmxModel) error {
	if modelData == nil {
		return fmt.Errorf("Xモデルがnilです")
	}
	p.model = modelData
	return p.parseTextXFile()
}

// peek は現在のトークンを返す。
func (p *textParser) peek() textToken {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return textToken{typ: tokEOF}
}

// next は次のトークンへ進めて返す。
func (p *textParser) next() textToken {
	t := p.peek()
	p.pos++
	return t
}

// expect は指定トークンを期待する。
func (p *textParser) expect(typ tokenType) (textToken, error) {
	t := p.next()
	if t.typ != typ {
		return t, fmt.Errorf("expected %v got %v (%s)", typ, t.typ, t.val)
	}
	return t, nil
}

// parseTextXFile はXテキスト全体を解析する。
func (p *textParser) parseTextXFile() error {
	for p.peek().typ != tokEOF {
		t := p.peek()
		switch {
		case t.typ == tokIdentifier && t.val == "template":
			p.next()
			if err := p.parseTextTemplateDefinition(); err != nil {
				return err
			}
		case t.typ == tokIdentifier && t.val == "Header":
			p.next()
			if err := p.parseTextHeader(); err != nil {
				return err
			}
		case t.typ == tokIdentifier && t.val == "Mesh":
			p.next()
			if err := p.parseTextMesh(); err != nil {
				return err
			}
		case t.typ == tokIdentifier:
			p.next()
			if p.peek().typ == tokLCurly {
				if err := p.skipTextUnknownTemplate(); err != nil {
					return err
				}
			}
		default:
			p.next()
		}
	}
	return nil
}

// parseTextTemplateDefinition はテンプレート定義を読み飛ばす。
func (p *textParser) parseTextTemplateDefinition() error {
	if _, err := p.expect(tokIdentifier); err != nil {
		return err
	}
	if _, err := p.expect(tokLCurly); err != nil {
		return err
	}
	guidTok := p.next()
	if guidTok.typ != tokAngleBracketed {
		return fmt.Errorf("テンプレート定義のGUIDが不正です")
	}
	braceCount := 1
	for braceCount > 0 {
		tok := p.next()
		switch tok.typ {
		case tokLCurly:
			braceCount++
		case tokRCurly:
			braceCount--
		case tokEOF:
			return fmt.Errorf("テンプレート定義が未完です")
		}
	}
	return nil
}

// parseTextHeader はHeaderブロックを解析する。
func (p *textParser) parseTextHeader() error {
	if _, err := p.expect(tokLCurly); err != nil {
		return err
	}
	major, err := p.parseNumberAsFloat()
	if err != nil {
		return err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	minor, err := p.parseNumberAsFloat()
	if err != nil {
		return err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	flags, err := p.parseNumberAsFloat()
	if err != nil {
		return err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	if _, err := p.expect(tokRCurly); err != nil {
		return err
	}
	p.model.Comment = fmt.Sprintf("X File Version %.0f.%.0f, flags: %.0f", major, minor, flags)
	return nil
}

// parseNumberAsFloat は数値トークンをfloat64として読み取る。
func (p *textParser) parseNumberAsFloat() (float64, error) {
	t, err := p.expect(tokNumber)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(t.val, 64)
}

// parseNumberAsInt は数値トークンをintとして読み取る。
func (p *textParser) parseNumberAsInt() (int, error) {
	t, err := p.expect(tokNumber)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(t.val, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

// parseString は文字列トークンを読み取る。
func (p *textParser) parseString() (string, error) {
	t := p.next()
	switch t.typ {
	case tokString, tokIdentifier:
		return t.val, nil
	default:
		return "", fmt.Errorf("文字列の読み取りに失敗しました")
	}
}

// parseVector はVec3を読み取る。
func (p *textParser) parseVector() (mmath.Vec3, error) {
	x, err := p.parseNumberAsFloat()
	if err != nil {
		return mmath.Vec3{}, err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return mmath.Vec3{}, err
	}
	y, err := p.parseNumberAsFloat()
	if err != nil {
		return mmath.Vec3{}, err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return mmath.Vec3{}, err
	}
	z, err := p.parseNumberAsFloat()
	if err != nil {
		return mmath.Vec3{}, err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return mmath.Vec3{}, err
	}
	return mmath.Vec3{Vec: r3.Vec{X: x, Y: y, Z: z}}, nil
}

// parseColorRGBA はRGBAを読み取る。
func (p *textParser) parseColorRGBA() (mmath.Vec4, error) {
	r, err := p.parseNumberAsFloat()
	if err != nil {
		return mmath.Vec4{}, err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return mmath.Vec4{}, err
	}
	g, err := p.parseNumberAsFloat()
	if err != nil {
		return mmath.Vec4{}, err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return mmath.Vec4{}, err
	}
	b, err := p.parseNumberAsFloat()
	if err != nil {
		return mmath.Vec4{}, err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return mmath.Vec4{}, err
	}
	a, err := p.parseNumberAsFloat()
	if err != nil {
		return mmath.Vec4{}, err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return mmath.Vec4{}, err
	}
	return mmath.Vec4{X: r, Y: g, Z: b, W: a}, nil
}

// parseColorRGB はRGBを読み取る。
func (p *textParser) parseColorRGB() (mmath.Vec3, error) {
	r, err := p.parseNumberAsFloat()
	if err != nil {
		return mmath.Vec3{}, err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return mmath.Vec3{}, err
	}
	g, err := p.parseNumberAsFloat()
	if err != nil {
		return mmath.Vec3{}, err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return mmath.Vec3{}, err
	}
	b, err := p.parseNumberAsFloat()
	if err != nil {
		return mmath.Vec3{}, err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return mmath.Vec3{}, err
	}
	return mmath.Vec3{Vec: r3.Vec{X: r, Y: g, Z: b}}, nil
}

// parseTextMesh はMeshブロックを解析する。
func (p *textParser) parseTextMesh() error {
	if p.peek().typ != tokLCurly {
		_, _ = p.parseString()
	}
	if _, err := p.expect(tokLCurly); err != nil {
		return err
	}
	p.meshCtx = &meshContext{vertexOffset: p.model.Vertices.Len(), faceGroups: make([][]*model.Face, 0)}

	vertexCount, err := p.parseNumberAsInt()
	if err != nil {
		return err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	for i := 0; i < vertexCount; i++ {
		pos, err := p.parseVector()
		if err != nil {
			return err
		}
		pos.MulScalar(10)
		vertex := &model.Vertex{Position: pos}
		vertex.Normal = mmath.Vec3{Vec: r3.Vec{X: 0, Y: 1, Z: 0}}
		deform := model.NewBdef1(0)
		vertex.Deform = deform
		vertex.DeformType = deform.DeformType()
		vertex.EdgeFactor = 1
		p.model.Vertices.AppendRaw(vertex)
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}

	faceCount, err := p.parseNumberAsInt()
	if err != nil {
		return err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	p.meshCtx.faceGroups = make([][]*model.Face, 0, faceCount)
	for i := 0; i < faceCount; i++ {
		faces, err := p.parseMeshFace()
		if err != nil {
			return err
		}
		p.meshCtx.faceGroups = append(p.meshCtx.faceGroups, faces)
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}

	hasMaterialList := false
	for p.peek().typ == tokIdentifier {
		switch p.peek().val {
		case "MeshMaterialList":
			hasMaterialList = true
			p.next()
			if err := p.parseMeshMaterialList(); err != nil {
				return err
			}
		case "MeshTextureCoords":
			p.next()
			if err := p.parseMeshTextureCoords(); err != nil {
				return err
			}
		case "MeshNormals":
			p.next()
			if err := p.skipTextUnknownTemplate(); err != nil {
				return err
			}
		default:
			p.next()
			if err := p.skipTextUnknownTemplate(); err != nil {
				return err
			}
		}
	}

	if !hasMaterialList {
		if err := p.appendMeshFacesWithDefault(); err != nil {
			return err
		}
	}
	if _, err := p.expect(tokRCurly); err != nil {
		return err
	}
	p.meshCtx = nil
	return nil
}

// parseMeshFace は面定義を読み取る。
func (p *textParser) parseMeshFace() ([]*model.Face, error) {
	count, err := p.parseNumberAsInt()
	if err != nil {
		return nil, err
	}
	if count < 3 {
		return nil, fmt.Errorf("Meshの面頂点数が不足しています")
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return nil, err
	}
	indexes := make([]int, 0, count)
	for i := 0; i < count; i++ {
		idx, err := p.parseNumberAsInt()
		if err != nil {
			return nil, err
		}
		indexes = append(indexes, idx)
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return nil, err
	}

	base := p.meshCtx.vertexOffset
	switch {
	case count == 4:
		f1 := &model.Face{VertexIndexes: [3]int{indexes[0] + base, indexes[1] + base, indexes[2] + base}}
		f2 := &model.Face{VertexIndexes: [3]int{indexes[0] + base, indexes[2] + base, indexes[3] + base}}
		return []*model.Face{f1, f2}, nil
	default:
		f := &model.Face{VertexIndexes: [3]int{indexes[0] + base, indexes[1] + base, indexes[2] + base}}
		return []*model.Face{f}, nil
	}
}

// parseMeshMaterialList はMeshMaterialListを解析する。
func (p *textParser) parseMeshMaterialList() error {
	if _, err := p.expect(tokLCurly); err != nil {
		return err
	}
	nMaterials, err := p.parseNumberAsInt()
	if err != nil {
		return err
	}
	if nMaterials < 0 {
		return fmt.Errorf("MeshMaterialListの材質数が不正です")
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	nFaceIdx, err := p.parseNumberAsInt()
	if err != nil {
		return err
	}
	if nFaceIdx < 0 {
		return fmt.Errorf("MeshMaterialListの面数が不正です")
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	if p.meshCtx == nil {
		return fmt.Errorf("MeshMaterialListがMesh外で検出されました")
	}
	if nFaceIdx != len(p.meshCtx.faceGroups) {
		return fmt.Errorf("MeshMaterialListの面数が不正です")
	}
	facesByMaterial := make([][][]*model.Face, nMaterials)
	for i := 0; i < nMaterials; i++ {
		facesByMaterial[i] = make([][]*model.Face, 0)
	}
	for i := 0; i < nFaceIdx; i++ {
		matIdx, err := p.parseNumberAsInt()
		if err != nil {
			return err
		}
		if matIdx < 0 || matIdx >= nMaterials {
			return fmt.Errorf("MeshMaterialListの材質番号が不正です")
		}
		facesByMaterial[matIdx] = append(facesByMaterial[matIdx], p.meshCtx.faceGroups[i])
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	if p.peek().typ == tokSemicolon {
		_, _ = p.expect(tokSemicolon)
	}

	for i := 0; i < nMaterials; i++ {
		if p.peek().typ == tokIdentifier && p.peek().val == "Material" {
			p.next()
			if err := p.parseMaterialText(); err != nil {
				return err
			}
		} else {
			if p.peek().typ == tokIdentifier {
				p.next()
				if err := p.skipTextUnknownTemplate(); err != nil {
					return err
				}
			}
			material := newDefaultMaterial(p.model.Materials.Len())
			p.model.Materials.AppendRaw(material)
		}
	}

	for i := 0; i < nMaterials; i++ {
		mat, err := p.model.Materials.Get(i)
		if err != nil {
			return err
		}
		for _, faces := range facesByMaterial[i] {
			for _, face := range faces {
				p.model.Faces.AppendRaw(face)
				mat.VerticesCount += 3
			}
		}
	}

	if _, err := p.expect(tokRCurly); err != nil {
		return err
	}
	return nil
}

// parseMaterialText はMaterialブロックを解析する。
func (p *textParser) parseMaterialText() error {
	if p.peek().typ != tokLCurly {
		_, _ = p.parseString()
	}
	if _, err := p.expect(tokLCurly); err != nil {
		return err
	}
	material := model.NewMaterial()
	diffuse, err := p.parseColorRGBA()
	if err != nil {
		return err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	power, err := p.parseNumberAsFloat()
	if err != nil {
		return err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	specular, err := p.parseColorRGB()
	if err != nil {
		return err
	}
	material.Specular = mmath.Vec4{X: specular.X, Y: specular.Y, Z: specular.Z, W: power}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	ambient, err := p.parseColorRGB()
	if err != nil {
		return err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	material.Diffuse = diffuse
	material.Ambient = ambient
	material.Edge = mmath.UNIT_W_VEC4
	material.EdgeSize = 10

	for p.peek().typ == tokIdentifier && p.peek().val == "TextureFilename" {
		p.next()
		texName, sphereName, err := p.parseTextureFilename()
		if err != nil {
			return err
		}
		if texName != "" {
			material.TextureIndex = ensureTexture(p.model, texName, model.TEXTURE_TYPE_TEXTURE)
		}
		if sphereName != "" {
			material.SphereTextureIndex = ensureTexture(p.model, sphereName, model.TEXTURE_TYPE_SPHERE)
		}
	}
	if _, err := p.expect(tokRCurly); err != nil {
		return err
	}

	material.SetName(fmt.Sprintf("材質%02d", p.model.Materials.Len()+1))
	applySphereMode(material)
	p.model.Materials.AppendRaw(material)
	return nil
}

// parseTextureFilename はTextureFilenameを解析する。
func (p *textParser) parseTextureFilename() (string, string, error) {
	if _, err := p.expect(tokLCurly); err != nil {
		return "", "", err
	}
	name, err := p.parseString()
	if err != nil {
		return "", "", err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return "", "", err
	}
	if _, err := p.expect(tokRCurly); err != nil {
		return "", "", err
	}
	texName, sphereName := splitTextureName(name)
	return texName, sphereName, nil
}

// parseMeshTextureCoords はUVを読み取る。
func (p *textParser) parseMeshTextureCoords() error {
	if _, err := p.expect(tokLCurly); err != nil {
		return err
	}
	count, err := p.parseNumberAsInt()
	if err != nil {
		return err
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	for i := 0; i < count; i++ {
		u, err := p.parseNumberAsFloat()
		if err != nil {
			return err
		}
		if _, err := p.expect(tokSemicolon); err != nil {
			return err
		}
		v, err := p.parseNumberAsFloat()
		if err != nil {
			return err
		}
		if _, err := p.expect(tokSemicolon); err != nil {
			return err
		}
		vertex, err := p.model.Vertices.Get(p.meshCtx.vertexOffset + i)
		if err != nil {
			return err
		}
		vertex.Uv = mmath.Vec2{X: u, Y: v}
	}
	if _, err := p.expect(tokSemicolon); err != nil {
		return err
	}
	if _, err := p.expect(tokRCurly); err != nil {
		return err
	}
	return nil
}

// skipTextUnknownTemplate は未知ブロックを読み飛ばす。
func (p *textParser) skipTextUnknownTemplate() error {
	tok := p.next()
	if tok.typ != tokLCurly {
		return nil
	}
	braceCount := 1
	for braceCount > 0 {
		cur := p.next()
		switch cur.typ {
		case tokLCurly:
			braceCount++
		case tokRCurly:
			braceCount--
		case tokEOF:
			return fmt.Errorf("テンプレートの終端が不正です")
		}
	}
	return nil
}

// appendMeshFacesWithDefault は材質なしの面を追加する。
func (p *textParser) appendMeshFacesWithDefault() error {
	if p.meshCtx == nil {
		return nil
	}
	material := newDefaultMaterial(p.model.Materials.Len())
	p.model.Materials.AppendRaw(material)
	for _, faces := range p.meshCtx.faceGroups {
		for _, face := range faces {
			p.model.Faces.AppendRaw(face)
			material.VerticesCount += 3
		}
	}
	return nil
}

// newDefaultMaterial は既定値の材質を生成する。
func newDefaultMaterial(index int) *model.Material {
	material := model.NewMaterial()
	material.SetName(fmt.Sprintf("材質%02d", index+1))
	material.Edge = mmath.UNIT_W_VEC4
	material.EdgeSize = 10
	applySphereMode(material)
	return material
}

// ensureTexture はテクスチャを追加または取得する。
func ensureTexture(modelData *model.PmxModel, name string, texType model.TextureType) int {
	if modelData == nil || modelData.Textures == nil {
		return -1
	}
	tex, err := modelData.Textures.GetByName(name)
	if err == nil && tex != nil {
		return tex.Index()
	}
	newTex := model.NewTexture()
	newTex.SetName(name)
	newTex.TextureType = texType
	newTex.SetValid(true)
	modelData.Textures.AppendRaw(newTex)
	return newTex.Index()
}

// applySphereMode はスフィアモードを設定する。
func applySphereMode(material *model.Material) {
	if material == nil {
		return
	}
	if material.TextureIndex >= 0 && material.SphereTextureIndex < 0 {
		material.SphereMode = model.SPHERE_MODE_INVALID
		return
	}
	material.SphereMode = model.SPHERE_MODE_MULTIPLICATION
}

// splitTextureName はテクスチャ名とスフィア名へ分割する。
func splitTextureName(value string) (string, string) {
	parts := strings.Split(value, "*")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	lower := strings.ToLower(strings.TrimSpace(value))
	if strings.HasSuffix(lower, ".sph") {
		return "", strings.TrimSpace(value)
	}
	return strings.TrimSpace(value), ""
}
