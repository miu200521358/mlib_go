// 指示: miu200521358
package mfile

// PathService はパス操作の実装を表す。
type PathService struct{}

// NewPathService はPathServiceを生成する。
func NewPathService() *PathService {
	return &PathService{}
}

// CanSave は保存可能なパスか判定する。
func (s *PathService) CanSave(path string) bool {
	return CanSave(path)
}

// CreateOutputPath は出力パスを生成する。
func (s *PathService) CreateOutputPath(originalPath, label string) string {
	return CreateOutputPath(originalPath, label)
}

// SplitPath はパスを dir/name/ext に分割する。
func (s *PathService) SplitPath(path string) (dir, name, ext string) {
	return SplitPath(path)
}
