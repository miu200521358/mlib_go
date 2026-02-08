// 指示: miu200521358
package io_csv

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
)

// csvFieldMeta はCSVタグ対象フィールドのメタ情報を表す。
type csvFieldMeta struct {
	name  string
	index int
}

// Marshal はCSVタグ付き構造体スライスをCsvModelへ変換する。
func Marshal(data any) (*CsvModel, error) {
	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return nil, io_common.NewIoEncodeFailed("CSV変換対象が不正です", nil)
	}

	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, io_common.NewIoEncodeFailed("CSV変換対象がnilです", nil)
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return nil, io_common.NewIoEncodeFailed("CSV変換対象はスライス配列で指定してください", nil)
	}

	structType, pointerElement, err := resolveStructType(value.Type().Elem())
	if err != nil {
		return nil, io_common.NewIoEncodeFailed("CSV変換対象の要素型が不正です", err)
	}

	fields, err := collectCsvFields(structType)
	if err != nil {
		return nil, io_common.NewIoEncodeFailed("CSVタグ解析に失敗しました", err)
	}
	if len(fields) == 0 {
		return nil, io_common.NewIoEncodeFailed("csvタグが定義されたフィールドがありません", nil)
	}

	records := make([][]string, 0, value.Len()+1)
	records = append(records, csvHeaderRow(fields))

	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		structValue, resolveErr := resolveStructValue(elem, pointerElement)
		if resolveErr != nil {
			return nil, io_common.NewIoEncodeFailed("CSV変換対象の要素取得に失敗しました(行:%d)", resolveErr, i+2)
		}

		row, rowErr := marshalCsvRow(structValue, fields)
		if rowErr != nil {
			return nil, io_common.NewIoEncodeFailed("CSV行変換に失敗しました(行:%d)", rowErr, i+2)
		}
		records = append(records, row)
	}

	model := NewCsvModel(records)
	model.UpdateHash()
	return model, nil
}

// Unmarshal はCsvModelをCSVタグ付き構造体スライスへ変換する。
func Unmarshal(model *CsvModel, out any) error {
	if model == nil {
		return io_common.NewIoParseFailed("CSVモデルがnilです", nil)
	}

	outValue := reflect.ValueOf(out)
	if !outValue.IsValid() || outValue.Kind() != reflect.Pointer || outValue.IsNil() {
		return io_common.NewIoParseFailed("出力先は非nilポインタを指定してください", nil)
	}
	sliceValue := outValue.Elem()
	if sliceValue.Kind() != reflect.Slice {
		return io_common.NewIoParseFailed("出力先はスライスポインタを指定してください", nil)
	}

	structType, pointerElement, err := resolveStructType(sliceValue.Type().Elem())
	if err != nil {
		return io_common.NewIoParseFailed("出力先要素型が不正です", err)
	}

	fields, err := collectCsvFields(structType)
	if err != nil {
		return io_common.NewIoParseFailed("CSVタグ解析に失敗しました", err)
	}
	if len(fields) == 0 {
		return io_common.NewIoParseFailed("csvタグが定義されたフィールドがありません", nil)
	}

	records := model.Records()
	if len(records) == 0 {
		sliceValue.Set(reflect.MakeSlice(sliceValue.Type(), 0, 0))
		return nil
	}

	fieldIndexByName := csvFieldIndexMap(fields)
	columnToField := mapCsvColumns(records[0], fieldIndexByName)

	result := reflect.MakeSlice(sliceValue.Type(), 0, max(0, len(records)-1))
	for i := 1; i < len(records); i++ {
		row := records[i]
		if isEmptyCsvRow(row) {
			continue
		}

		structValue := reflect.New(structType).Elem()
		if err := unmarshalCsvRow(structValue, row, columnToField, fields, i+1); err != nil {
			return err
		}

		if pointerElement {
			result = reflect.Append(result, structValue.Addr())
		} else {
			result = reflect.Append(result, structValue)
		}
	}

	sliceValue.Set(result)
	return nil
}

// resolveStructType は要素型から構造体型を解決し、ポインタ要素かを返す。
func resolveStructType(elemType reflect.Type) (reflect.Type, bool, error) {
	pointerElement := false
	if elemType.Kind() == reflect.Pointer {
		pointerElement = true
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		return nil, false, fmt.Errorf("要素型が構造体ではありません: %s", elemType.Kind())
	}
	return elemType, pointerElement, nil
}

// collectCsvFields はcsvタグ付きフィールドを抽出する。
func collectCsvFields(structType reflect.Type) ([]csvFieldMeta, error) {
	fields := make([]csvFieldMeta, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := parseCsvTag(field.Tag.Get("csv"))
		if tag == "" {
			continue
		}
		if field.PkgPath != "" {
			return nil, fmt.Errorf("csvタグ付きフィールドは公開フィールドである必要があります: %s", field.Name)
		}
		fields = append(fields, csvFieldMeta{
			name:  tag,
			index: i,
		})
	}
	return fields, nil
}

// parseCsvTag はcsvタグ値をヘッダ名として解釈する。
func parseCsvTag(raw string) string {
	if raw == "" || raw == "-" {
		return ""
	}
	parts := strings.Split(raw, ",")
	return strings.TrimSpace(parts[0])
}

// csvHeaderRow はフィールド定義からヘッダ行を生成する。
func csvHeaderRow(fields []csvFieldMeta) []string {
	header := make([]string, len(fields))
	for i, field := range fields {
		header[i] = field.name
	}
	return header
}

// resolveStructValue は行要素を構造体値として解決する。
func resolveStructValue(elem reflect.Value, pointerElement bool) (reflect.Value, error) {
	if pointerElement {
		if elem.IsNil() {
			return reflect.Value{}, fmt.Errorf("ポインタ要素がnilです")
		}
		elem = elem.Elem()
	}
	if elem.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("要素が構造体ではありません: %s", elem.Kind())
	}
	return elem, nil
}

// marshalCsvRow は構造体1件をCSV行へ変換する。
func marshalCsvRow(structValue reflect.Value, fields []csvFieldMeta) ([]string, error) {
	row := make([]string, len(fields))
	for i, field := range fields {
		cellValue, err := marshalCsvCell(structValue.Field(field.index))
		if err != nil {
			return nil, err
		}
		row[i] = cellValue
	}
	return row, nil
}

// marshalCsvCell はフィールド値をCSVセル文字列へ変換する。
func marshalCsvCell(value reflect.Value) (string, error) {
	if !value.IsValid() {
		return "", nil
	}

	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return "", nil
		}
		return marshalCsvCell(value.Elem())
	}

	if value.CanInterface() {
		if marshaler, ok := value.Interface().(encoding.TextMarshaler); ok {
			text, err := marshaler.MarshalText()
			if err != nil {
				return "", err
			}
			return string(text), nil
		}
	}

	switch value.Kind() {
	case reflect.String:
		return value.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(value.Uint(), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(value.Float(), 'f', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64), nil
	default:
		if value.CanInterface() {
			return fmt.Sprint(value.Interface()), nil
		}
		return "", fmt.Errorf("未対応型です: %s", value.Kind())
	}
}

// csvFieldIndexMap はヘッダ名からフィールド番号を引けるマップを作成する。
func csvFieldIndexMap(fields []csvFieldMeta) map[string]int {
	indexMap := make(map[string]int, len(fields))
	for i, field := range fields {
		indexMap[field.name] = i
	}
	return indexMap
}

// mapCsvColumns はヘッダ行から列番号とフィールド番号の対応表を作成する。
func mapCsvColumns(header []string, fieldIndexByName map[string]int) map[int]int {
	columnToField := make(map[int]int, len(header))
	for columnIndex, name := range header {
		fieldIndex, ok := fieldIndexByName[strings.TrimSpace(name)]
		if !ok {
			continue
		}
		columnToField[columnIndex] = fieldIndex
	}
	return columnToField
}

// unmarshalCsvRow はCSV行を構造体へ展開する。
func unmarshalCsvRow(structValue reflect.Value, row []string, columnToField map[int]int, fields []csvFieldMeta, rowNumber int) error {
	for columnIndex, cell := range row {
		fieldIndex, ok := columnToField[columnIndex]
		if !ok || fieldIndex >= len(fields) {
			continue
		}
		targetField := structValue.Field(fields[fieldIndex].index)
		if err := unmarshalCsvCell(targetField, cell); err != nil {
			return io_common.NewIoParseFailed(
				"CSV行変換に失敗しました(行:%d 列:%d 値:%s)",
				err,
				rowNumber,
				columnIndex+1,
				cell,
			)
		}
	}
	return nil
}

// unmarshalCsvCell はCSVセル文字列をフィールドへ変換設定する。
func unmarshalCsvCell(targetField reflect.Value, cell string) error {
	if !targetField.CanSet() {
		return fmt.Errorf("フィールドに設定できません")
	}

	if targetField.Kind() == reflect.Pointer {
		if strings.TrimSpace(cell) == "" {
			return nil
		}
		value := reflect.New(targetField.Type().Elem())
		if err := unmarshalCsvCell(value.Elem(), cell); err != nil {
			return err
		}
		targetField.Set(value)
		return nil
	}

	if targetField.CanAddr() {
		if unmarshaler, ok := targetField.Addr().Interface().(encoding.TextUnmarshaler); ok {
			return unmarshaler.UnmarshalText([]byte(cell))
		}
	}

	switch targetField.Kind() {
	case reflect.String:
		targetField.SetString(cell)
		return nil
	case reflect.Bool:
		value, err := strconv.ParseBool(strings.TrimSpace(cell))
		if err != nil {
			return err
		}
		targetField.SetBool(value)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(strings.TrimSpace(cell), 10, targetField.Type().Bits())
		if err != nil {
			return err
		}
		targetField.SetInt(value)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		value, err := strconv.ParseUint(strings.TrimSpace(cell), 10, targetField.Type().Bits())
		if err != nil {
			return err
		}
		targetField.SetUint(value)
		return nil
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(strings.TrimSpace(cell), targetField.Type().Bits())
		if err != nil {
			return err
		}
		targetField.SetFloat(value)
		return nil
	default:
		return fmt.Errorf("未対応型です: %s", targetField.Kind())
	}
}

// max は2値の大きい方を返す。
func max(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}
