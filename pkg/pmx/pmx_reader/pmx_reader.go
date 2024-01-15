package pmx_reader

import (
	"golang.org/x/text/encoding"

	"github.com/miu200521358/mlib_go/pkg/core/reader"
	"github.com/miu200521358/mlib_go/pkg/pmx/pmx_model"
)

type PmxReader struct {
	reader.BaseReader[*pmx_model.PmxModel]
}

func (r *PmxReader) CreateModel(path string) pmx_model.PmxModel {
	return *pmx_model.NewPmxModel(path)
}

func (r *PmxReader) ReadHeader(model pmx_model.PmxModel) {
	fbytes, err := r.UnpackBytes(4)
	if err != nil {
		return ""
	}
	// implementation goes here
}

func (r *PmxReader) ReadData(model pmx_model.PmxModel) {
	// implementation goes here
}

func (r *PmxReader) DefineEncoding(model encoding.Encoding) {
	// implementation goes here
}
