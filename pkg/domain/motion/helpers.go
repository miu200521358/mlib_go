// 指示: miu200521358
package motion

import (
	"reflect"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// copyVec3 はVec3ポインタを複製する。
func copyVec3(src *mmath.Vec3) *mmath.Vec3 {
	if src == nil {
		return nil
	}
	v := *src
	return &v
}

// copyQuaternion はQuaternionポインタを複製する。
func copyQuaternion(src *mmath.Quaternion) *mmath.Quaternion {
	if src == nil {
		return nil
	}
	q := *src
	return &q
}

// copyBoolPtr はboolポインタを複製する。
func copyBoolPtr(src *bool) *bool {
	if src == nil {
		return nil
	}
	v := *src
	return &v
}

// vec3OrZero はnilを0ベクトルとして扱う。
func vec3OrZero(src *mmath.Vec3) mmath.Vec3 {
	if src == nil {
		return mmath.Vec3{}
	}
	return *src
}

// vec3OrUnit はnilを単位ベクトル(1,1,1)として扱う。
func vec3OrUnit(src *mmath.Vec3) mmath.Vec3 {
	if src == nil {
		return vec3(1, 1, 1)
	}
	return *src
}

// quatOrIdent はnilを単位クォータニオンとして扱う。
func quatOrIdent(src *mmath.Quaternion) mmath.Quaternion {
	if src == nil {
		return mmath.NewQuaternion()
	}
	return *src
}

// vec3 はVec3を生成する。
func vec3(x, y, z float64) mmath.Vec3 {
	v := mmath.NewVec3()
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

// vec3Ptr はVec3ポインタを生成する。
func vec3Ptr(x, y, z float64) *mmath.Vec3 {
	v := vec3(x, y, z)
	return &v
}

// isNilValue はnil判定可能な型をnil判定する。
func isNilValue[T any](value T) bool {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return v.IsNil()
	default:
		return false
	}
}
