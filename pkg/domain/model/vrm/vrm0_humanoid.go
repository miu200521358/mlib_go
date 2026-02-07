// 指示: miu200521358
package vrm

// Vrm0HumanBone はVRM0 humanoid.humanBones要素を表す。
type Vrm0HumanBone struct {
	Bone string
	Node int
}

// Vrm0Humanoid はVRM0 humanoid要素を表す。
type Vrm0Humanoid struct {
	HumanBones []Vrm0HumanBone
}
