// 指示: miu200521358
package vrm

// Vrm1HumanBone はVRM1 humanoid.humanBones要素を表す。
type Vrm1HumanBone struct {
	Node int
}

// Vrm1Humanoid はVRM1 humanoid要素を表す。
type Vrm1Humanoid struct {
	HumanBones map[string]Vrm1HumanBone
}
