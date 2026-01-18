// 指示: miu200521358
package x

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"gonum.org/v1/gonum/spatial/r3"
)

func TestXRepository_Load_Text(t *testing.T) {
	path := writeTempFile(t, "sample_text.x", []byte(buildTextX()))

	r := NewXRepository()
	data, err := r.Load(path)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		t.Fatalf("Expected model type to be *PmxModel, got %T", data)
	}
	assertXModel(t, modelData, 4, mmath.Vec3{Vec: r3.Vec{X: 0, Y: 1, Z: 0}}, 9, 3)
}

func TestXRepository_Load_Binary(t *testing.T) {
	path := writeTempFile(t, "sample_bin.x", buildBinaryX())

	r := NewXRepository()
	data, err := r.Load(path)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		t.Fatalf("Expected model type to be *PmxModel, got %T", data)
	}
	assertXModel(t, modelData, 3, mmath.Vec3{Vec: r3.Vec{X: 0, Y: 0, Z: 1}}, 6, 3)
}

func TestXRepository_Load_CompressedBinary(t *testing.T) {
	compressed := buildCompressedX(buildBinaryX())
	path := writeTempFile(t, "sample_bzip.x", compressed)

	r := NewXRepository()
	data, err := r.Load(path)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		t.Fatalf("Expected model type to be *PmxModel, got %T", data)
	}
	assertXModel(t, modelData, 3, mmath.Vec3{Vec: r3.Vec{X: 0, Y: 0, Z: 1}}, 6, 3)
}

func assertXModel(t *testing.T, modelData *model.PmxModel, faceCount int, expectedNormal mmath.Vec3, mat0Count, mat1Count int) {
	t.Helper()
	if modelData == nil {
		t.Fatalf("Expected modelData to be not nil")
	}

	if modelData.Vertices.Len() != 5 {
		t.Fatalf("Expected vertex count to be %d, got %d", 5, modelData.Vertices.Len())
	}
	v1, _ := modelData.Vertices.Get(1)
	expectedPos := mmath.Vec3{Vec: r3.Vec{X: 10, Y: 0, Z: 0}}
	if !v1.Position.MMD().NearEquals(expectedPos, 1e-8) {
		t.Errorf("Expected Position to be %v, got %v", expectedPos, v1.Position.MMD())
	}
	v0, _ := modelData.Vertices.Get(0)
	if !v0.Normal.MMD().NearEquals(expectedNormal, 1e-8) {
		t.Errorf("Expected Normal to be %v, got %v", expectedNormal, v0.Normal.MMD())
	}
	v4, _ := modelData.Vertices.Get(4)
	if !v4.Uv.NearEquals(mmath.Vec2{X: 0.5, Y: 0.5}, 1e-8) {
		t.Errorf("Expected UV to be %v, got %v", mmath.Vec2{X: 0.5, Y: 0.5}, v4.Uv)
	}

	if modelData.Faces.Len() != faceCount {
		t.Fatalf("Expected face count to be %d, got %d", faceCount, modelData.Faces.Len())
	}

	if modelData.Materials.Len() != 2 {
		t.Fatalf("Expected material count to be %d, got %d", 2, modelData.Materials.Len())
	}
	mat0, _ := modelData.Materials.Get(0)
	mat1, _ := modelData.Materials.Get(1)
	if mat0.VerticesCount != mat0Count {
		t.Errorf("Expected material0 VerticesCount to be %d, got %d", mat0Count, mat0.VerticesCount)
	}
	if mat1.VerticesCount != mat1Count {
		t.Errorf("Expected material1 VerticesCount to be %d, got %d", mat1Count, mat1.VerticesCount)
	}
	if mat0.SphereMode != model.SPHERE_MODE_INVALID {
		t.Errorf("Expected material0 SphereMode to be Invalid")
	}
	if mat1.SphereMode != model.SPHERE_MODE_MULTIPLICATION {
		t.Errorf("Expected material1 SphereMode to be Multiplication")
	}

	if _, err := modelData.Bones.GetByName("センター"); err != nil {
		t.Fatalf("Expected center bone to exist, got %q", err)
	}
}

func buildTextX() string {
	lines := []string{
		"xof 0303txt 0032",
		"Header {",
		"1;",
		"0;",
		"1;",
		"}",
		"Mesh {",
		"5;",
		"0;0;0;,",
		"1;0;0;,",
		"0;1;0;,",
		"0;0;1;,",
		"1;1;1;;",
		"3;",
		"4;0,1,2,3;,",
		"3;0,2,4;,",
		"5;1,2,3,4,0;;",
		"MeshMaterialList {",
		"2;",
		"3;",
		"0,1,0;;",
		"Material {",
		"1;0;0;1;;",
		"5;",
		"0.1;0.2;0.3;;",
		"0.4;0.5;0.6;;",
		"TextureFilename { \"tex.png\"; }",
		"}",
		"Material {",
		"0.2;0.3;0.4;0.5;;",
		"10;",
		"0.5;0.4;0.3;;",
		"0.1;0.2;0.3;;",
		"TextureFilename { \"sphere.sph\"; }",
		"}",
		"}",
		"MeshTextureCoords {",
		"5;",
		"0;0;,",
		"1;0;,",
		"1;1;,",
		"0;1;,",
		"0.5;0.5;;",
		"}",
		"}",
	}
	return strings.Join(lines, "\n")
}

func buildBinaryX() []byte {
	payload := &bytes.Buffer{}
	writeName(payload, "Mesh")
	writeToken(payload, tokenOBrace)
	writeFloatList(payload, []float64{
		0, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
		1, 1, 1,
	})
	writeIntegerList(payload, []uint32{
		3,
		4, 0, 1, 2,
		3, 0, 2, 4,
		3, 1, 3, 4,
	})

	writeName(payload, "MeshNormals")
	writeToken(payload, tokenOBrace)
	writeFloatList(payload, []float64{
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	})
	writeToken(payload, tokenCBrace)

	writeName(payload, "MeshMaterialList")
	writeToken(payload, tokenOBrace)
	writeIntegerList(payload, []uint32{2, 3, 0, 1, 0})

	writeName(payload, "Material")
	writeToken(payload, tokenOBrace)
	writeFloatList(payload, []float64{
		1, 0, 0, 1,
		5,
		0.1, 0.2, 0.3,
		0.4, 0.5, 0.6,
	})
	writeName(payload, "TextureFilename")
	writeToken(payload, tokenOBrace)
	writeStringList(payload, []string{"tex.png"})
	writeToken(payload, tokenCBrace)
	writeToken(payload, tokenCBrace)

	writeName(payload, "Material")
	writeToken(payload, tokenOBrace)
	writeFloatList(payload, []float64{
		0.2, 0.3, 0.4, 0.5,
		10,
		0.5, 0.4, 0.3,
		0.1, 0.2, 0.3,
	})
	writeName(payload, "TextureFilename")
	writeToken(payload, tokenOBrace)
	writeStringList(payload, []string{"sphere.sph"})
	writeToken(payload, tokenCBrace)
	writeToken(payload, tokenCBrace)

	writeToken(payload, tokenCBrace)

	writeName(payload, "MeshTextureCoords")
	writeToken(payload, tokenOBrace)
	writeFloatList(payload, []float64{
		0, 0,
		1, 0,
		1, 1,
		0, 1,
		0.5, 0.5,
	})
	writeToken(payload, tokenCBrace)

	writeToken(payload, tokenCBrace)

	header := []byte("xof 0303tzip0032")
	return append(header, payload.Bytes()...)
}

func buildCompressedX(binaryData []byte) []byte {
	if len(binaryData) < xHeaderSize {
		return binaryData
	}
	head := append([]byte{}, binaryData[:xHeaderSize]...)
	copy(head[8:12], []byte("bzip"))
	payload := binaryData[xHeaderSize:]

	var compressed bytes.Buffer
	compressed.Write(head)
	_ = binary.Write(&compressed, binary.LittleEndian, uint32(len(head)+len(payload)))

	var block bytes.Buffer
	block.Write([]byte{'C', 'K'})
	w, _ := flate.NewWriter(&block, flate.DefaultCompression)
	_, _ = w.Write(payload)
	_ = w.Close()

	_ = binary.Write(&compressed, binary.LittleEndian, uint16(len(payload)))
	_ = binary.Write(&compressed, binary.LittleEndian, uint16(block.Len()))
	compressed.Write(block.Bytes())
	return compressed.Bytes()
}

func writeTempFile(t *testing.T, name string, data []byte) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("Expected write to succeed, got %q", err)
	}
	return path
}

func writeToken(buf *bytes.Buffer, tok binaryTokenType) {
	_ = binary.Write(buf, binary.LittleEndian, uint16(tok))
}

func writeName(buf *bytes.Buffer, name string) {
	writeToken(buf, tokenName)
	b := []byte(name)
	_ = binary.Write(buf, binary.LittleEndian, uint32(len(b)))
	_, _ = buf.Write(b)
}

func writeStringList(buf *bytes.Buffer, values []string) {
	for _, value := range values {
		writeToken(buf, tokenString)
		b := []byte(value)
		_ = binary.Write(buf, binary.LittleEndian, uint32(len(b)))
		_, _ = buf.Write(b)
		writeToken(buf, tokenSemicolon)
	}
}

func writeIntegerList(buf *bytes.Buffer, values []uint32) {
	writeToken(buf, tokenIntegerList)
	_ = binary.Write(buf, binary.LittleEndian, uint32(len(values)))
	for _, v := range values {
		_ = binary.Write(buf, binary.LittleEndian, v)
	}
}

func writeFloatList(buf *bytes.Buffer, values []float64) {
	writeToken(buf, tokenFloatList)
	_ = binary.Write(buf, binary.LittleEndian, uint32(len(values)))
	for _, v := range values {
		_ = binary.Write(buf, binary.LittleEndian, math.Float32bits(float32(v)))
	}
}
