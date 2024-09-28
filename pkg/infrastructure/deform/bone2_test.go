package deform

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

func TestVmdMotion_DeformArmIk4_DMF(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/nac_dmf_601.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, nil)
	{
		expectedPosition := &mmath.MVec3{X: 6.210230, Y: 8.439670, Z: 0.496305}
		if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 6.210230, Y: 8.849669, Z: 0.496305}
		if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 6.210230, Y: 12.836980, Z: -0.159825}
		if !boneDeltas.GetByName("上半身").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("上半身").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 6.261481, Y: 13.968025, Z: 0.288966}
		if !boneDeltas.GetByName("上半身2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("上半身2").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 6.541666, Y: 15.754716, Z: 1.421828}
		if !boneDeltas.GetByName("左肩").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左肩").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.451898, Y: 16.031992, Z: 1.675949}
		if !boneDeltas.GetByName("左腕").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.135534, Y: 15.373729, Z: 1.715530}
		if !boneDeltas.GetByName("左腕捩").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.646749, Y: 13.918620, Z: 1.803021}
		if !boneDeltas.GetByName("左ひじ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじ").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.164164, Y: 13.503792, Z: 1.706635}
		if !boneDeltas.GetByName("左手捩").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.772219, Y: 12.307291, Z: 1.428628}
		if !boneDeltas.GetByName("左手首").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手首").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.390504, Y: 12.011601, Z: 1.405503}
		if !boneDeltas.GetByName("左手先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.451900, Y: 16.031990, Z: 1.675949}
		if !boneDeltas.GetByName("左腕YZ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕YZ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕YZ").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.690105, Y: 15.802624, Z: 1.689741}
		if !boneDeltas.GetByName("左腕YZ先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕YZ先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕YZ先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.690105, Y: 15.802622, Z: 1.689740}
		if !boneDeltas.GetByName("左腕YZIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕YZIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕YZIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.451899, Y: 16.031988, Z: 1.675950}
		if !boneDeltas.GetByName("左腕X").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕X").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕X").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.816861, Y: 16.406412, Z: 1.599419}
		if !boneDeltas.GetByName("左腕X先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕X先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕X先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.816858, Y: 16.406418, Z: 1.599418}
		if !boneDeltas.GetByName("左腕XIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕XIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕XIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.135530, Y: 15.373726, Z: 1.715530}
		if !boneDeltas.GetByName("左腕捩YZ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩YZ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩YZ").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.409824, Y: 15.109610, Z: 1.731412}
		if !boneDeltas.GetByName("左腕捩YZTgt").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩YZTgt").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩YZTgt").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.409830, Y: 15.109617, Z: 1.731411}
		if !boneDeltas.GetByName("左腕捩YZIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩YZIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩YZIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.135530, Y: 15.373725, Z: 1.715531}
		if !boneDeltas.GetByName("左腕捩X").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩X").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩X").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.500528, Y: 15.748149, Z: 1.639511}
		if !boneDeltas.GetByName("左腕捩XTgt").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩XTgt").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩XTgt").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.500531, Y: 15.748233, Z: 1.639508}
		if !boneDeltas.GetByName("左腕捩XIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩XIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩XIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.646743, Y: 13.918595, Z: 1.803029}
		if !boneDeltas.GetByName("左ひじYZ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYZ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYZ").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.360763, Y: 13.672787, Z: 1.745903}
		if !boneDeltas.GetByName("左ひじYZ先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYZ先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYZ先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.360781, Y: 13.672805, Z: 1.745905}
		if !boneDeltas.GetByName("左ひじYZIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYZIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYZIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.646734, Y: 13.918593, Z: 1.803028}
		if !boneDeltas.GetByName("左ひじX").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじX").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじX").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.944283, Y: 13.652989, Z: 1.456379}
		if !boneDeltas.GetByName("左ひじX先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじX先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじX先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.944304, Y: 13.653007, Z: 1.456381}
		if !boneDeltas.GetByName("左ひじXIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじXIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじXIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.646734, Y: 13.918596, Z: 1.803028}
		if !boneDeltas.GetByName("左ひじY").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじY").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじY").FilledGlobalPosition().MMD()))
		}
	}
	{
		// FIXME
		expectedPosition := &mmath.MVec3{X: 9.560862, Y: 13.926876, Z: 1.431514}
		if !boneDeltas.GetByName("左ひじY先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじY先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじY先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.360781, Y: 13.672805, Z: 1.745905}
		if !boneDeltas.GetByName("左ひじYIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.164141, Y: 13.503780, Z: 1.706625}
		if !boneDeltas.GetByName("左手捩YZ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩YZ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩YZ").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.414344, Y: 12.859288, Z: 1.556843}
		if !boneDeltas.GetByName("左手捩YZTgt").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩YZTgt").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩YZTgt").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.414370, Y: 12.859282, Z: 1.556885}
		if !boneDeltas.GetByName("左手捩YZIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩YZIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩YZIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.164142, Y: 13.503780, Z: 1.706624}
		if !boneDeltas.GetByName("左手捩X").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩X").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩X").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.511073, Y: 12.928087, Z: 2.447041}
		if !boneDeltas.GetByName("左手捩XTgt").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩XTgt").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩XTgt").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.511120, Y: 12.928122, Z: 2.447057}
		if !boneDeltas.GetByName("左手捩XIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩XIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩XIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.471097, Y: 12.074032, Z: 1.410383}
		if !boneDeltas.GetByName("左手YZ先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手YZ先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手YZ先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.471111, Y: 12.074042, Z: 1.410384}
		if !boneDeltas.GetByName("左手YZIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手YZIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手YZIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.772183, Y: 12.307314, Z: 1.428564}
		if !boneDeltas.GetByName("左手X").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手X").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手X").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.802912, Y: 12.308764, Z: 0.901022}
		if !boneDeltas.GetByName("左手X先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手X先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手X先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.802991, Y: 12.308830, Z: 0.901079}
		if !boneDeltas.GetByName("左手XIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手XIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手XIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.130125, Y: 15.368912, Z: 1.728851}
		if !boneDeltas.GetByName("左腕捩1").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩1").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩1").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.365511, Y: 15.142246, Z: 1.742475}
		if !boneDeltas.GetByName("左腕捩2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩2").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.639965, Y: 14.877952, Z: 1.758356}
		if !boneDeltas.GetByName("左腕捩3").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩3").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩3").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.901615, Y: 14.625986, Z: 1.773497}
		if !boneDeltas.GetByName("左腕捩4").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩4").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩4").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.641270, Y: 13.913721, Z: 1.816324}
		if !boneDeltas.GetByName("左ひじsub").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじsub").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじsub").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.907782, Y: 13.661371, Z: 2.034630}
		if !boneDeltas.GetByName("左ひじsub先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじsub先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじsub先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 9.165060, Y: 13.499348, Z: 1.721094}
		if !boneDeltas.GetByName("左手捩1").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩1").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩1").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.978877, Y: 13.339340, Z: 1.683909}
		if !boneDeltas.GetByName("左手捩2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩2").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.741154, Y: 13.135028, Z: 1.636428}
		if !boneDeltas.GetByName("左手捩3").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩3").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩3").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.516553, Y: 12.942023, Z: 1.591578}
		if !boneDeltas.GetByName("左手捩4").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩4").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩4").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.301016, Y: 12.748707, Z: 1.544439}
		if !boneDeltas.GetByName("左手捩5").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩5").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩5").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 8.145000, Y: 12.614601, Z: 1.513277}
		if !boneDeltas.GetByName("左手捩6").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩6").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩6").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.777408, Y: 12.298634, Z: 1.439762}
		if !boneDeltas.GetByName("左手首R").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首R").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手首R").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.777408, Y: 12.298635, Z: 1.439762}
		if !boneDeltas.GetByName("左手首1").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首1").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手首1").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 7.670320, Y: 12.202144, Z: 1.486689}
		if !boneDeltas.GetByName("左手首2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手首2").FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk25_Addiction_Wa_Left(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/[A]ddiction_和洋_0126F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	boneDeltas := DeformBone(model, motion, true, 0, []string{"左襟先"})

	{
		expectedPosition := &mmath.MVec3{X: -0.225006, Y: 9.705784, Z: 2.033072}
		if !boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.237383, Y: 10.769137, Z: 2.039952}
		if !boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER2.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.460140, Y: 13.290816, Z: 2.531440}
		if !boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.460140, Y: 13.290816, Z: 2.531440}
		if !boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.784452, Y: 13.728909, Z: 2.608527}
		if !boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Left()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.272067, Y: 12.381887, Z: 2.182425}
		if !boneDeltas.GetByName("左上半身C-A").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身C-A").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左上半身C-A").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.406217, Y: 13.797803, Z: 2.460243}
		if !boneDeltas.GetByName("左鎖骨IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左鎖骨IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左鎖骨IK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.052427, Y: 13.347448, Z: 2.216718}
		if !boneDeltas.GetByName("左鎖骨").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左鎖骨").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左鎖骨").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.659554, Y: 14.017852, Z: 2.591099}
		if !boneDeltas.GetByName("左鎖骨先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左鎖骨先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左鎖骨先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.065147, Y: 13.296564, Z: 2.607907}
		if !boneDeltas.GetByName("左肩Rz検").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Rz検").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左肩Rz検").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.517776, Y: 14.134196, Z: 2.645912}
		if !boneDeltas.GetByName("左肩Rz検先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Rz検先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左肩Rz検先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.065148, Y: 13.296564, Z: 2.607907}
		if !boneDeltas.GetByName("左肩Ry検").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Ry検").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左肩Ry検").FilledGlobalPosition().MMD()))
		}
	}
	{
		// FIXME
		expectedPosition := &mmath.MVec3{X: 0.860159, Y: 13.190875, Z: 3.122428}
		if !boneDeltas.GetByName("左肩Ry検先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Ry検先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左肩Ry検先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.195053, Y: 12.648546, Z: 2.236849}
		if !boneDeltas.GetByName("左上半身C-B").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身C-B").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左上半身C-B").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.294257, Y: 12.912640, Z: 2.257159}
		if !boneDeltas.GetByName("左上半身C-C").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身C-C").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左上半身C-C").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.210011, Y: 10.897711, Z: 1.973442}
		if !boneDeltas.GetByName("左上半身2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左上半身2").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.320589, Y: 14.049745, Z: 2.637018}
		if !boneDeltas.GetByName("左襟").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左襟").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左襟").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.297636, Y: 14.263302, Z: 2.374467}
		if !boneDeltas.GetByName("左襟先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左襟先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左襟先").FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformLegIk30_Addiction_Shoes(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/[A]ddiction_和洋_1037F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	boneDeltas := DeformBone(model, motion, true, 0, nil)

	{
		expectedPosition := &mmath.MVec3{X: 0, Y: 0, Z: 0}
		if !boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.748025, Y: 1.683590, Z: 0.556993}
		if !boneDeltas.GetByName(pmx.LEG_IK.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.111190, Y: 4.955496, Z: 1.070225}
		if !boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.724965, Y: 3.674735, Z: 0.810759}
		if !boneDeltas.GetByName(pmx.LEG.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.924419, Y: 3.943549, Z: -1.420897}
		if !boneDeltas.GetByName(pmx.KNEE.Left()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.743128, Y: 1.672784, Z: 0.551317}
		if !boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.743128, Y: 1.672784, Z: 0.551317}
		if !boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.118480, Y: 3.432016, Z: -1.657329}
		if !boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D先").FilledGlobalPosition().MMD()))
		}
	}
	{
		// FIXME
		expectedPosition := &mmath.MVec3{X: 1.763123, Y: 3.708842, Z: 0.369619}
		if !boneDeltas.GetByName("左足Dw").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足Dw").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足Dw").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.698137, Y: 0.674393, Z: 0.043128}
		if !boneDeltas.GetByName("左足先EX").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足先EX").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足先EX").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.712729, Y: 0.919505, Z: 0.174835}
		if !boneDeltas.GetByName("左素足先A").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.494684, Y: 0.328695, Z: -0.500715}
		if !boneDeltas.GetByName("左素足先A先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.494684, Y: 0.328695, Z: -0.500715}
		if !boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先AIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.224127, Y: 1.738605, Z: 0.240794}
		if !boneDeltas.GetByName("左素足先B").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.582481, Y: 0.627567, Z: -0.222556}
		if !boneDeltas.GetByName("左素足先B先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.572457, Y: 0.658645, Z: -0.209595}
		if !boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先BIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.448783, Y: 1.760351, Z: 0.405702}
		if !boneDeltas.GetByName("左靴調節").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴調節").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴調節").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.875421, Y: 0.917357, Z: 0.652144}
		if !boneDeltas.GetByName("左靴追従").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.404960, Y: 1.846940, Z: 0.380388}
		if !boneDeltas.GetByName("左靴追従先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.448783, Y: 1.760351, Z: 0.405702}
		if !boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従IK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.330700, Y: 3.021906, Z: 0.153679}
		if !boneDeltas.GetByName("左足補D").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補D").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.625984, Y: 3.444340, Z: -1.272553}
		if !boneDeltas.GetByName("左足補D先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補D先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.625984, Y: 3.444340, Z: -1.272553}
		if !boneDeltas.GetByName("左足補DIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補DIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補DIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.724964, Y: 3.674735, Z: 0.810759}
		if !boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.924420, Y: 3.943550, Z: -1.420897}
		if !boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.924419, Y: 3.943550, Z: -1.420896}
		if !boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検AIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.724965, Y: 3.674735, Z: 0.810760}
		if !boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向-").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 1.566813, Y: 3.645544, Z: 0.177956}
		if !boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足w").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -1.879506, Y: 3.670895, Z: -2.715526}
		if !boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足w先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.631964, Y: 3.647012, Z: -1.211210}
		if !boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.83923, Y: 2.876222, Z: -1.494900}
		if !boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.686400, Y: 3.444961, Z: -1.285575}
		if !boneDeltas.GetByName("左膝補IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補IK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.924420, Y: 3.943550, Z: -1.420896}
		if !boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.952927, Y: 7.388766, Z: -0.972059}
		if !boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B先").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: -0.952927, Y: 7.388766, Z: -0.972060}
		if !boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検BIK").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.400272, Y: 3.809143, Z: -0.305068}
		if !boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩").FilledGlobalPosition().MMD()))
		}
	}
	{
		expectedPosition := &mmath.MVec3{X: 0.371963, Y: 4.256704, Z: -0.267830}
		if !boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
			t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩先").FilledGlobalPosition().MMD()))
		}
	}
}

func TestVmdMotion_DeformIk_Down(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/センター下げる.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式/MEIKO準標準_400.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)
	DeformBone(model, motion, true, 0, nil)
}
