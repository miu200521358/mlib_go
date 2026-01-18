// 指示: miu200521358
package x

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"gonum.org/v1/gonum/spatial/r3"
)

type binaryTokenType uint16

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
	tokenEOF         binaryTokenType = 9999
)

type binaryParser struct {
	data               []byte
	pos                int
	floatSize          int
	model              *model.PmxModel
	meshCtx            *meshContext
	materialTokenCount int
}

// newBinaryParser はXバイナリのパーサーを生成する。
func newBinaryParser(data []byte, modelData *model.PmxModel) *binaryParser {
	return &binaryParser{data: data, model: modelData}
}

// Parse はXバイナリを解析する。
func (p *binaryParser) Parse() error {
	if err := p.parseHeader(); err != nil {
		return err
	}
	for {
		tok, err := p.getToken()
		if err != nil {
			return err
		}
		switch tok {
		case tokenEOF:
			return nil
		case tokenTemplate:
			if err := p.skipTemplateDefinition(); err != nil {
				return err
			}
		case tokenName:
			p.pos -= 2
			name, err := p.readName()
			if err != nil {
				return err
			}
			if err := p.parseObject(name); err != nil {
				return err
			}
		case tokenWord, tokenDword, tokenFloat, tokenDouble, tokenChar, tokenUchar, tokenSword, tokenSdword, tokenLpstr, tokenUnicode, tokenCstring:
			if err := p.parseObject(""); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Xバイナリのトークンが不正です: %d", tok)
		}
	}
}

// parseHeader はヘッダを解析する。
func (p *binaryParser) parseHeader() error {
	if len(p.data) < 16 {
		return fmt.Errorf("Xバイナリヘッダが不足しています")
	}
	floatSize := string(p.data[12:16])
	if floatSize == "0032" {
		p.floatSize = 4
	} else {
		p.floatSize = 8
	}
	p.pos = 16
	return nil
}

// getToken は次のトークンを取得する。
func (p *binaryParser) getToken() (binaryTokenType, error) {
	if p.pos > len(p.data)-2 {
		return tokenEOF, nil
	}
	value := binary.LittleEndian.Uint16(p.data[p.pos:])
	p.pos += 2
	return binaryTokenType(value), nil
}

// readDWORD はDWORDを読み取る。
func (p *binaryParser) readDWORD() (uint32, error) {
	if p.pos > len(p.data)-4 {
		return 0, fmt.Errorf("DWORD読み取りでバッファが不足しています")
	}
	value := binary.LittleEndian.Uint32(p.data[p.pos:])
	p.pos += 4
	return value, nil
}

// readFloat は浮動小数点を読み取る。
func (p *binaryParser) readFloat() (float64, error) {
	if p.pos > len(p.data)-p.floatSize {
		return 0, fmt.Errorf("float読み取りでバッファが不足しています")
	}
	if p.floatSize == 4 {
		bits := binary.LittleEndian.Uint32(p.data[p.pos:])
		p.pos += p.floatSize
		return float64(math.Float32frombits(bits)), nil
	}
	bits := binary.LittleEndian.Uint64(p.data[p.pos:])
	p.pos += p.floatSize
	return math.Float64frombits(bits), nil
}

// readBytes は指定数のバイト列を読み取る。
func (p *binaryParser) readBytes(count int) ([]byte, error) {
	if p.pos > len(p.data)-count {
		return nil, fmt.Errorf("bytes読み取りでバッファが不足しています")
	}
	buf := make([]byte, count)
	copy(buf, p.data[p.pos:p.pos+count])
	p.pos += count
	return buf, nil
}

// readName はNAMEトークンを読み取る。
func (p *binaryParser) readName() (string, error) {
	tok, err := p.getToken()
	if err != nil {
		return "", err
	}
	if tok != tokenName {
		return "", fmt.Errorf("NAMEトークンが不正です: %d", tok)
	}
	count, err := p.readDWORD()
	if err != nil {
		return "", err
	}
	nameBytes, err := p.readBytes(int(count))
	if err != nil {
		return "", err
	}
	nameBytes = bytes.TrimRight(nameBytes, "\x00")
	decoded, err := decodeShiftJIS(nameBytes)
	if err != nil {
		return "", err
	}
	return decoded, nil
}

// readString はSTRINGトークンを読み取る。
func (p *binaryParser) readString() (string, error) {
	tok, err := p.getToken()
	if err != nil {
		return "", err
	}
	if tok != tokenString {
		return "", fmt.Errorf("STRINGトークンが不正です: %d", tok)
	}
	count, err := p.readDWORD()
	if err != nil {
		return "", err
	}
	buf, err := p.readBytes(int(count))
	if err != nil {
		return "", err
	}
	term, err := p.getToken()
	if err != nil {
		return "", err
	}
	if term != tokenSemicolon && term != tokenComma {
		return "", fmt.Errorf("STRING終端トークンが不正です: %d", term)
	}
	decoded, err := decodeShiftJIS(buf)
	if err != nil {
		return "", err
	}
	return decoded, nil
}

// parseObject はオブジェクトブロックを解析する。
func (p *binaryParser) parseObject(objectName string) error {
	prevMesh := p.meshCtx
	if objectName == "Mesh" {
		p.meshCtx = &meshContext{vertexOffset: p.model.Vertices.Len()}
	}
	defer func() {
		if objectName == "Mesh" {
			p.meshCtx = prevMesh
		}
	}()

	nextTok, err := p.getToken()
	if err != nil {
		return err
	}
	if nextTok == tokenName {
		p.pos -= 2
		if _, err := p.readName(); err != nil {
			return err
		}
	} else {
		p.pos -= 2
	}
	braceTok, err := p.getToken()
	if err != nil {
		return err
	}
	if braceTok != tokenOBrace {
		return fmt.Errorf("オブジェクト開始が不正です: %d", braceTok)
	}
	peek, err := p.getToken()
	if err != nil {
		return err
	}
	p.pos -= 2
	if peek == tokenGuid {
		if err := p.skipGuid(); err != nil {
			return err
		}
	}

	currentName := objectName
	for {
		tok, err := p.getToken()
		if err != nil {
			return err
		}
		if tok == tokenCBrace {
			return nil
		}
		p.pos -= 2

		switch tok {
		case tokenOBrace:
			if err := p.parseDataReference(); err != nil {
				return err
			}
		case tokenName:
			name, err := p.readName()
			if err != nil {
				return err
			}
			peek, err := p.getToken()
			if err != nil {
				return err
			}
			p.pos -= 2
			if peek == tokenName || peek == tokenOBrace {
				if err := p.parseObject(name); err != nil {
					return err
				}
			} else {
				currentName = name
				if err := p.parseDataList(currentName); err != nil {
					return err
				}
			}
		case tokenIntegerList:
			if _, err := p.readIntegerList(currentName); err != nil {
				return err
			}
		case tokenFloatList:
			if _, err := p.readFloatList(currentName); err != nil {
				return err
			}
		case tokenString:
			if err := p.parseStringList(currentName); err != nil {
				return err
			}
		default:
			return fmt.Errorf("オブジェクト解析中のトークンが不正です: %d", tok)
		}
	}
}

// parseDataReference は参照ブロックを読み飛ばす。
func (p *binaryParser) parseDataReference() error {
	tok, err := p.getToken()
	if err != nil {
		return err
	}
	if tok != tokenOBrace {
		return fmt.Errorf("参照開始トークンが不正です: %d", tok)
	}
	nameTok, err := p.getToken()
	if err != nil {
		return err
	}
	if nameTok == tokenName {
		p.pos -= 2
		if _, err := p.readName(); err != nil {
			return err
		}
	}
	peek, err := p.getToken()
	if err != nil {
		return err
	}
	if peek == tokenGuid {
		p.pos -= 2
		if err := p.skipGuid(); err != nil {
			return err
		}
	} else {
		p.pos -= 2
	}
	closeTok, err := p.getToken()
	if err != nil {
		return err
	}
	if closeTok != tokenCBrace {
		return fmt.Errorf("参照終了トークンが不正です: %d", closeTok)
	}
	return nil
}

// parseDataList はデータリストを解析する。
func (p *binaryParser) parseDataList(objectName string) error {
	tok, err := p.getToken()
	if err != nil {
		return err
	}
	p.pos -= 2
	switch tok {
	case tokenIntegerList:
		_, err := p.readIntegerList(objectName)
		return err
	case tokenFloatList:
		_, err := p.readFloatList(objectName)
		return err
	case tokenString:
		return p.parseStringList(objectName)
	default:
		return fmt.Errorf("データリストのトークンが不正です: %d", tok)
	}
}

// readIntegerList は整数リストを読み取る。
func (p *binaryParser) readIntegerList(objectName string) ([]uint32, error) {
	tok, err := p.getToken()
	if err != nil {
		return nil, err
	}
	if tok != tokenIntegerList {
		return nil, fmt.Errorf("INTEGER_LISTトークンが不正です: %d", tok)
	}
	count, err := p.readDWORD()
	if err != nil {
		return nil, err
	}
	list := make([]uint32, int(count))
	for i := range list {
		value, err := p.readDWORD()
		if err != nil {
			return nil, err
		}
		list[i] = value
	}

	switch objectName {
	case "Mesh":
		if p.meshCtx == nil {
			return nil, fmt.Errorf("Meshの頂点コンテキストが不正です")
		}
		if len(list) < 1 {
			return nil, fmt.Errorf("Meshの面情報が不足しています")
		}
		faceCount := int(list[0])
		base := p.meshCtx.vertexOffset
		for i := 0; i < faceCount; i++ {
			idx := 1 + i*4
			if idx+3 >= len(list) {
				return nil, fmt.Errorf("Meshの面情報が不足しています")
			}
			face := &model.Face{VertexIndexes: [3]int{int(list[idx+1]) + base, int(list[idx+2]) + base, int(list[idx+3]) + base}}
			p.model.Faces.AppendRaw(face)
		}
	case "MeshMaterialList":
		if len(list) < 2 {
			return nil, fmt.Errorf("MeshMaterialListの情報が不足しています")
		}
		matCount := int(list[0])
		faceCount := int(list[1])
		for p.model.Materials.Len() < matCount {
			material := newDefaultMaterial(p.model.Materials.Len())
			p.model.Materials.AppendRaw(material)
		}
		if 2+faceCount > len(list) {
			return nil, fmt.Errorf("MeshMaterialListの面数が不正です")
		}
		for i := 0; i < faceCount; i++ {
			matIdx := int(list[2+i])
			material, err := p.model.Materials.Get(matIdx)
			if err != nil {
				return nil, err
			}
			material.VerticesCount += 3
		}
	}

	return list, nil
}

// readFloatList は浮動小数点リストを読み取る。
func (p *binaryParser) readFloatList(objectName string) ([]float64, error) {
	tok, err := p.getToken()
	if err != nil {
		return nil, err
	}
	if tok != tokenFloatList {
		return nil, fmt.Errorf("FLOAT_LISTトークンが不正です: %d", tok)
	}
	count, err := p.readDWORD()
	if err != nil {
		return nil, err
	}
	list := make([]float64, int(count))
	for i := range list {
		value, err := p.readFloat()
		if err != nil {
			return nil, err
		}
		list[i] = value
	}

	switch objectName {
	case "Mesh":
		if p.meshCtx == nil {
			return nil, fmt.Errorf("Meshの頂点コンテキストが不正です")
		}
		for i := 0; i+2 < len(list); i += 3 {
			pos := mmath.Vec3{Vec: r3.Vec{X: list[i], Y: list[i+1], Z: list[i+2]}}
			pos.MulScalar(10)
			vertex := &model.Vertex{Position: pos}
			vertex.Normal = mmath.Vec3{Vec: r3.Vec{X: 0, Y: 1, Z: 0}}
			deform := model.NewBdef1(0)
			vertex.Deform = deform
			vertex.DeformType = deform.DeformType()
			vertex.EdgeFactor = 1
			p.model.Vertices.AppendRaw(vertex)
		}
	case "MeshNormals":
		if p.meshCtx == nil {
			return nil, fmt.Errorf("MeshNormalsの頂点コンテキストが不正です")
		}
		for i := 0; i+2 < len(list); i += 3 {
			vidx := i / 3
			vertex, err := p.model.Vertices.Get(p.meshCtx.vertexOffset + vidx)
			if err != nil {
				return nil, err
			}
			vertex.Normal = mmath.Vec3{Vec: r3.Vec{X: list[i], Y: list[i+1], Z: list[i+2]}}
		}
	case "MeshTextureCoords":
		if p.meshCtx == nil {
			return nil, fmt.Errorf("MeshTextureCoordsの頂点コンテキストが不正です")
		}
		for i := 0; i+1 < len(list); i += 2 {
			vidx := i / 2
			vertex, err := p.model.Vertices.Get(p.meshCtx.vertexOffset + vidx)
			if err != nil {
				return nil, err
			}
			vertex.Uv = mmath.Vec2{X: list[i], Y: list[i+1]}
		}
	case "Material":
		material := p.ensureMaterial(p.materialTokenCount)
		if len(list) >= 11 {
			material.Diffuse = mmath.Vec4{X: list[0], Y: list[1], Z: list[2], W: list[3]}
			material.Specular = mmath.Vec4{X: list[5], Y: list[6], Z: list[7], W: list[4]}
			material.Ambient = mmath.Vec3{Vec: r3.Vec{X: list[8], Y: list[9], Z: list[10]}}
		}
		material.Edge = mmath.UNIT_W_VEC4
		material.EdgeSize = 10
		applySphereMode(material)
		p.materialTokenCount++
	}

	return list, nil
}

// parseStringList は文字列リストを読み取る。
func (p *binaryParser) parseStringList(objectName string) error {
	texts := make([]string, 0)
	for {
		stringTok, err := p.getToken()
		if err != nil {
			return err
		}
		if stringTok != tokenString {
			return fmt.Errorf("STRINGトークンが不正です: %d", stringTok)
		}
		p.pos -= 2
		text, err := p.readString()
		if err != nil {
			return err
		}
		texts = append(texts, text)
		nextTok, err := p.getToken()
		if err != nil {
			return err
		}
		if nextTok != tokenComma && nextTok != tokenSemicolon {
			p.pos -= 2
			break
		}
	}

	switch objectName {
	case "TextureFilename":
		materialIndex := p.materialTokenCount - 1
		for _, text := range texts {
			texName, sphereName := splitTextureName(text)
			if texName != "" {
				idx := ensureTexture(p.model, texName, model.TEXTURE_TYPE_TEXTURE)
				if material, err := p.model.Materials.Get(materialIndex); err == nil {
					material.TextureIndex = idx
					applySphereMode(material)
				}
			}
			if sphereName != "" {
				idx := ensureTexture(p.model, sphereName, model.TEXTURE_TYPE_SPHERE)
				if material, err := p.model.Materials.Get(materialIndex); err == nil {
					material.SphereTextureIndex = idx
					applySphereMode(material)
				}
			}
		}
	}
	return nil
}

// ensureMaterial は指定位置の材質を確保する。
func (p *binaryParser) ensureMaterial(index int) *model.Material {
	for p.model.Materials.Len() <= index {
		material := newDefaultMaterial(p.model.Materials.Len())
		p.model.Materials.AppendRaw(material)
	}
	material, _ := p.model.Materials.Get(index)
	return material
}

// skipGuid はGUIDトークンを読み飛ばす。
func (p *binaryParser) skipGuid() error {
	tok, err := p.getToken()
	if err != nil {
		return err
	}
	if tok != tokenGuid {
		return fmt.Errorf("GUIDトークンが不正です: %d", tok)
	}
	if _, err := p.readBytes(16); err != nil {
		return err
	}
	return nil
}

// skipTemplateDefinition はテンプレート定義を読み飛ばす。
func (p *binaryParser) skipTemplateDefinition() error {
	if _, err := p.readName(); err != nil {
		return err
	}
	if tok, err := p.getToken(); err != nil || tok != tokenOBrace {
		return fmt.Errorf("テンプレート開始トークンが不正です")
	}
	if err := p.skipGuid(); err != nil {
		return err
	}
	braceCount := 1
	for braceCount > 0 {
		tok, err := p.getToken()
		if err != nil {
			return err
		}
		switch tok {
		case tokenOBrace:
			braceCount++
		case tokenCBrace:
			braceCount--
		case tokenEOF:
			return fmt.Errorf("テンプレート定義が未完です")
		}
	}
	return nil
}

// decodeShiftJIS はShift-JISをデコードする。
func decodeShiftJIS(raw []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(raw), japanese.ShiftJIS.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
