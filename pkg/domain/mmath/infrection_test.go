package mmath

import (
	"reflect"
	"testing"
)

func TestFindInflectionFrames(t *testing.T) {
	tests := []struct {
		name   string
		frames []float32
		values []float64
		want   []float32
	}{
		{
			name:   "単調増加（変曲点なし）",
			frames: []float32{1, 2, 3, 4, 5},
			values: []float64{1, 2, 3, 4, 5},
			want:   []float32{1, 5},
		},
		{
			name:   "同一値",
			frames: []float32{1, 2, 3, 4, 5},
			values: []float64{1, 1, 1, 1, 1},
			want:   []float32{1, 5},
		},
		{
			name:   "変曲点あり",
			frames: []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			values: []float64{1, 1, 1, 2, 2, 3, 3, 4, 2, 1},
			want:   []float32{1, 8, 10},
		},
		{
			name:   "変曲点あり2",
			frames: []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			values: []float64{1, 1, 5, 2, 2, 3, 8, 4, 2, 1},
			want:   []float32{1, 4, 5, 8, 10},
		},
		{
			name:   "実値1",
			frames: []float32{80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99},
			values: []float64{
				0.00000000000001654822503314023,
				0.00000000000003264784940820382,
				0.0006195860461806958,
				0.003702762582628298,
				0.007285729066060613,
				0.010962595412529252,
				0.014994659166345459,
				0.018886739886832225,
				0.023014777566295438,
				0.02686273068667996,
				0.030373480273382013,
				0.03413745687672446,
				0.037860423851525424,
				0.041409676484680386,
				0.044506058787641196,
				0.046220958956129404,
				0.04400085151306907,
				0.03996070456395067,
				0.03416680567852725,
				0.02600231409467793,
			},
			want: []float32{80, 82, 95, 99},
		},
	}

	for n, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindInflectionFrames(tt.frames, tt.values)
			got = UniqueFloat32s(got)
			SortFloat32s(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[%d:%s] FindInflectionFrames() = %v, want %v", n, tt.name, got, tt.want)
			}
		})
	}
}
