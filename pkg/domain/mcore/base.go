package mcore

// IndexModel はインデックスを持つモデルの基幹structです。
// 各エンティティ（Vertex, Face等）はこのstructを埋め込んで使用します。
type IndexModel struct {
	index int
}

// NewIndexModel は新しいIndexModelを生成します。
func NewIndexModel(index int) *IndexModel {
	return &IndexModel{index: index}
}

// Index はモデルのインデックスを返します。
func (m *IndexModel) Index() int {
	return m.index
}

// SetIndex はモデルのインデックスを設定します。
func (m *IndexModel) SetIndex(index int) {
	m.index = index
}

// IsValid はモデルが有効かどうかを返します。
// 埋め込み先でオーバーライド可能です。
func (m *IndexModel) IsValid() bool {
	return m.index >= 0
}

// IndexNameModel は名前とインデックスを持つモデルの基幹structです。
// Bone, Material, Morph等はこのstructを埋め込んで使用します。
type IndexNameModel struct {
	IndexModel
	name        string
	englishName string
}

// NewIndexNameModel は新しいIndexNameModelを生成します。
func NewIndexNameModel(index int, name, englishName string) *IndexNameModel {
	return &IndexNameModel{
		IndexModel:  IndexModel{index: index},
		name:        name,
		englishName: englishName,
	}
}

// Name はモデルの名前を返します。
func (m *IndexNameModel) Name() string {
	return m.name
}

// SetName はモデルの名前を設定します。
func (m *IndexNameModel) SetName(name string) {
	m.name = name
}

// EnglishName はモデルの英語名を返します。
func (m *IndexNameModel) EnglishName() string {
	return m.englishName
}

// SetEnglishName はモデルの英語名を設定します。
func (m *IndexNameModel) SetEnglishName(englishName string) {
	m.englishName = englishName
}

// IsValid はモデルが有効かどうかを返します。
// 埋め込み先でオーバーライド可能です。
func (m *IndexNameModel) IsValid() bool {
	return m.IndexModel.IsValid()
}
