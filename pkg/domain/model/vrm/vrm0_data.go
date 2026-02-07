// 指示: miu200521358
package vrm

// Vrm0Data はVRM0拡張情報を表す。
type Vrm0Data struct {
	ExporterVersion string
	Meta            *Vrm0Meta
	Humanoid        *Vrm0Humanoid
}

// NewVrm0Data はVrm0Dataを既定値で生成する。
func NewVrm0Data() *Vrm0Data {
	return &Vrm0Data{
		Humanoid: &Vrm0Humanoid{
			HumanBones: make([]Vrm0HumanBone, 0),
		},
	}
}
