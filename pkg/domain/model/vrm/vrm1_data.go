// 指示: miu200521358
package vrm

// Vrm1Data はVRM1拡張情報を表す。
type Vrm1Data struct {
	SpecVersion string
	Meta        *Vrm1Meta
	Humanoid    *Vrm1Humanoid
	SpringBone  *Vrm1SpringBone
}

// NewVrm1Data はVrm1Dataを既定値で生成する。
func NewVrm1Data() *Vrm1Data {
	return &Vrm1Data{
		Humanoid: &Vrm1Humanoid{
			HumanBones: map[string]Vrm1HumanBone{},
		},
		SpringBone: NewVrm1SpringBone(),
	}
}
