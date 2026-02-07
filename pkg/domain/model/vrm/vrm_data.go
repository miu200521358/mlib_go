// 指示: miu200521358
package vrm

import "encoding/json"

// VrmVersion はVRMの拡張バージョン種別を表す。
type VrmVersion string

const (
	// VRM_VERSION_0 はVRM0拡張を表す。
	VRM_VERSION_0 VrmVersion = "VRM0"
	// VRM_VERSION_1 はVRM1拡張を表す。
	VRM_VERSION_1 VrmVersion = "VRM1"
)

// VrmProfile はVRM読み込み時のプロファイルを表す。
type VrmProfile string

const (
	// VRM_PROFILE_STANDARD はVRM標準プロファイルを表す。
	VRM_PROFILE_STANDARD VrmProfile = "Standard"
	// VRM_PROFILE_VROID はVRoidプロファイルを表す。
	VRM_PROFILE_VROID VrmProfile = "VRoid"
)

// VrmData はVRM固有情報のルート構造を表す。
type VrmData struct {
	Version        VrmVersion
	Profile        VrmProfile
	AssetGenerator string
	Nodes          []Node
	Vrm0           *Vrm0Data
	Vrm1           *Vrm1Data
	RawExtensions  map[string]json.RawMessage
}

// NewVrmData はVrmDataを既定値で生成する。
func NewVrmData() *VrmData {
	return &VrmData{
		Profile:       VRM_PROFILE_STANDARD,
		Nodes:         make([]Node, 0),
		RawExtensions: map[string]json.RawMessage{},
	}
}
