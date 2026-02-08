// 指示: miu200521358
package io_csv

// CsvProfile はCSV構造の検証条件を表す。
type CsvProfile struct {
	// HasHeader は先頭行をヘッダとして扱うかを示す。
	HasHeader bool
	// ExactColumns は列数完全一致条件を示す。0の場合は未指定扱い。
	ExactColumns int
	// MinColumns は列数下限条件を示す。0の場合は未指定扱い。
	MinColumns int
	// Header はヘッダ名検証用の期待列名を示す。
	Header []string
	// AllowExtraColumns は余剰列を許容するかを示す。
	AllowExtraColumns bool
}

// NewFixedColumn4Profile は4列固定用の標準プロファイルを返す。
func NewFixedColumn4Profile() CsvProfile {
	return CsvProfile{
		HasHeader:         true,
		ExactColumns:      4,
		AllowExtraColumns: false,
	}
}

// NewSimpleKeyValueProfile はキー/値2列用の標準プロファイルを返す。
func NewSimpleKeyValueProfile() CsvProfile {
	return CsvProfile{
		HasHeader:         true,
		ExactColumns:      2,
		Header:            []string{"キー", "値"},
		AllowExtraColumns: false,
	}
}

// NewFreeTableProfile は最小1列の可変表用プロファイルを返す。
func NewFreeTableProfile() CsvProfile {
	return CsvProfile{
		HasHeader:         true,
		MinColumns:        1,
		AllowExtraColumns: true,
	}
}

// cloneCsvProfile はCsvProfileを複製して返す。
func cloneCsvProfile(profile CsvProfile) CsvProfile {
	cloned := profile
	if profile.Header != nil {
		cloned.Header = make([]string, len(profile.Header))
		copy(cloned.Header, profile.Header)
	}
	return cloned
}

// cloneCsvProfilePtr はCsvProfileポインタを安全に複製して返す。
func cloneCsvProfilePtr(profile *CsvProfile) *CsvProfile {
	if profile == nil {
		return nil
	}
	cloned := cloneCsvProfile(*profile)
	return &cloned
}
