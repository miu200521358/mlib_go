// 指示: miu200521358
package pmx

import "github.com/miu200521358/mlib_go/pkg/adapter/io_common"

// readIndex はバイトサイズに応じたインデックスを読み込む。
func readIndex(reader *io_common.BinaryReader, size byte, signed bool) (int, error) {
	switch size {
	case 1:
		if signed {
			value, err := reader.ReadInt8()
			return int(value), err
		}
		value, err := reader.ReadUint8()
		return int(value), err
	case 2:
		if signed {
			value, err := reader.ReadInt16()
			return int(value), err
		}
		value, err := reader.ReadUint16()
		return int(value), err
	case 4:
		value, err := reader.ReadInt32()
		return int(value), err
	default:
		return 0, io_common.NewIoFormatNotSupported("インデックスサイズが未対応です: %d", nil, size)
	}
}

// readVertexIndex は頂点インデックスを読み込む。
func readVertexIndex(reader *io_common.BinaryReader, size byte) (int, error) {
	return readIndex(reader, size, false)
}

// readSignedIndex は符号付きインデックスを読み込む。
func readSignedIndex(reader *io_common.BinaryReader, size byte) (int, error) {
	return readIndex(reader, size, true)
}

// writeIndex はバイトサイズに応じたインデックスを書き込む。
func writeIndex(writer *io_common.BinaryWriter, size byte, signed bool, value int) error {
	switch size {
	case 1:
		if signed {
			return writer.WriteInt8(int8(value))
		}
		return writer.WriteUint8(uint8(value))
	case 2:
		if signed {
			return writer.WriteInt16(int16(value))
		}
		return writer.WriteUint16(uint16(value))
	case 4:
		return writer.WriteInt32(int32(value))
	default:
		return io_common.NewIoFormatNotSupported("インデックスサイズが未対応です: %d", nil, size)
	}
}

// writeVertexIndex は頂点インデックスを書き込む。
func writeVertexIndex(writer *io_common.BinaryWriter, size byte, value int) error {
	return writeIndex(writer, size, false, value)
}

// writeSignedIndex は符号付きインデックスを書き込む。
func writeSignedIndex(writer *io_common.BinaryWriter, size byte, value int) error {
	return writeIndex(writer, size, true, value)
}
