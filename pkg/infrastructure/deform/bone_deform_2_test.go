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

	{

		fno := int(0)
		boneDeltas := DeformBone(motion.BoneFrames, fno, model, nil, true, nil, nil)
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
}

func TestVmdMotion_DeformArmIk2(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("C:/MMD/mmd_base/tests/resources/唱(ダンスのみ)_0274F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4_準標準.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		boneDeltas := DeformBone(motion.BoneFrames, fno, model, nil, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{X: 0.04952335, Y: 9.0, Z: 1.72378033}
			if !boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.04952335, Y: 7.97980869, Z: 1.72378033}
			if !boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.04952335, Y: 11.02838314, Z: 2.29172656}
			if !boneDeltas.GetByName("腰").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("腰").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("腰").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.04952335, Y: 11.9671191, Z: 1.06765032}
			if !boneDeltas.GetByName("上半身").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("上半身").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.26284261, Y: 13.14576297, Z: 0.84720008}
			if !boneDeltas.GetByName("上半身2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("上半身2").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.33636433, Y: 15.27729547, Z: 0.77435588}
			if !boneDeltas.GetByName("右肩").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.63104276, Y: 15.44542768, Z: 0.8507726}
			if !boneDeltas.GetByName("右肩C").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩C").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩C").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.63104276, Y: 15.44542768, Z: 0.8507726}
			if !boneDeltas.GetByName("右腕").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.90326269, Y: 14.53727204, Z: 0.7925801}
			if !boneDeltas.GetByName("右腕捩").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.50502977, Y: 12.52976106, Z: 0.66393998}
			if !boneDeltas.GetByName("右ひじ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじ").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.46843236, Y: 12.88476121, Z: 0.12831076}
			if !boneDeltas.GetByName("右手捩").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.36287259, Y: 13.90869981, Z: -1.41662258}
			if !boneDeltas.GetByName("右手首").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手首").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.81521586, Y: 14.00661535, Z: -1.55616424}
			if !boneDeltas.GetByName("右手先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.63104276, Y: 15.44542768, Z: 0.8507726}
			if !boneDeltas.GetByName("右腕YZ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕YZ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕YZ").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.72589296, Y: 15.12898892, Z: 0.83049645}
			if !boneDeltas.GetByName("右腕YZ先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕YZ先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕YZ先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.72589374, Y: 15.12898632, Z: 0.83049628}
			if !boneDeltas.GetByName("右腕YZIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕YZIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕YZIK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.63104276, Y: 15.44542768, Z: 0.8507726}
			if !boneDeltas.GetByName("右腕X").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕X").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕X").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.125321, Y: 15.600293, Z: 0.746130}
			if !boneDeltas.GetByName("右腕X先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕X先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕X先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.1253241, Y: 15.60029489, Z: 0.7461294}
			if !boneDeltas.GetByName("右腕XIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕XIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕XIK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.90325538, Y: 14.53727326, Z: 0.79258165}
			if !boneDeltas.GetByName("右腕捩YZ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩YZ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩YZ").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.01247534, Y: 14.17289417, Z: 0.76923367}
			if !boneDeltas.GetByName("右腕捩YZTgt").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩YZTgt").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩YZTgt").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.01248754, Y: 14.17289597, Z: 0.76923112}
			if !boneDeltas.GetByName("右腕捩YZIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩YZIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩YZIK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.90325538, Y: 14.53727326, Z: 0.79258165}
			if !boneDeltas.GetByName("右腕捩X").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩X").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩X").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.40656426, Y: 14.68386802, Z: 0.85919594}
			if !boneDeltas.GetByName("右腕捩XTgt").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩XTgt").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩XTgt").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.40657579, Y: 14.68387899, Z: 0.8591982}
			if !boneDeltas.GetByName("右腕捩XIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩XIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩XIK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.50499623, Y: 12.52974836, Z: 0.66394738}
			if !boneDeltas.GetByName("右ひじYZ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYZ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYZ").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.48334366, Y: 12.74011791, Z: 0.34655051}
			if !boneDeltas.GetByName("右ひじYZ先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYZ先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYZ先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.48334297, Y: 12.74012453, Z: 0.34654052}
			if !boneDeltas.GetByName("右ひじYZIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYZIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYZIK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.50499623, Y: 12.52974836, Z: 0.66394738}
			if !boneDeltas.GetByName("右ひじX").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじX").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじX").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -2.01179616, Y: 12.66809052, Z: 0.72106658}
			if !boneDeltas.GetByName("右ひじX先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじX先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじX先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -2.00760407, Y: 12.67958516, Z: 0.7289003}
			if !boneDeltas.GetByName("右ひじXIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじXIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじXIK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.50499623, Y: 12.52974836, Z: 0.66394738}
			if !boneDeltas.GetByName("右ひじY").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじY").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじY").FilledGlobalPosition().MMD()))
			}
		}
		{

			expectedPosition := &mmath.MVec3{X: -1.485519, Y: 12.740760, Z: 0.346835}
			if !boneDeltas.GetByName("右ひじY先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじY先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじY先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.48334297, Y: 12.74012453, Z: 0.34654052}
			if !boneDeltas.GetByName("右ひじYIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYIK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.46845628, Y: 12.88475892, Z: 0.12832214}
			if !boneDeltas.GetByName("右手捩YZ").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩YZ").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩YZ").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.41168478, Y: 13.4363328, Z: -0.7038697}
			if !boneDeltas.GetByName("右手捩YZTgt").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩YZTgt").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩YZTgt").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.41156715, Y: 13.43632015, Z: -0.70389025}
			if !boneDeltas.GetByName("右手捩YZIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩YZIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩YZIK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.46845628, Y: 12.88475892, Z: 0.12832214}
			if !boneDeltas.GetByName("右手捩X").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩X").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩X").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.5965686, Y: 12.06213832, Z: -0.42564769}
			if !boneDeltas.GetByName("右手捩XTgt").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩XTgt").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩XTgt").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.5965684, Y: 12.06214091, Z: -0.42565404}
			if !boneDeltas.GetByName("右手捩XIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩XIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩XIK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.7198605, Y: 13.98597326, Z: -1.5267472}
			if !boneDeltas.GetByName("右手YZ先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手YZ先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手YZ先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.71969424, Y: 13.98593727, Z: -1.52669587}
			if !boneDeltas.GetByName("右手YZIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手YZIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手YZIK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.36306295, Y: 13.90872698, Z: -1.41659848}
			if !boneDeltas.GetByName("右手X").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手X").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手X").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.54727182, Y: 13.56147176, Z: -1.06342964}
			if !boneDeltas.GetByName("右手X先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手X先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手X先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.54700171, Y: 13.5614545, Z: -1.0633896}
			if !boneDeltas.GetByName("右手XIK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手XIK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手XIK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.90581859, Y: 14.5370842, Z: 0.80752276}
			if !boneDeltas.GetByName("右腕捩1").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩1").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩1").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.99954005, Y: 14.2243783, Z: 0.78748743}
			if !boneDeltas.GetByName("右腕捩2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩2").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.10880907, Y: 13.85976329, Z: 0.76412793}
			if !boneDeltas.GetByName("右腕捩3").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩3").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩3").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.21298069, Y: 13.51216081, Z: 0.74185819}
			if !boneDeltas.GetByName("右腕捩4").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩4").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩4").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.5074743, Y: 12.52953348, Z: 0.67889319}
			if !boneDeltas.GetByName("右ひじsub").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじsub").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじsub").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.617075, Y: 12.131149, Z: 0.786797}
			if !boneDeltas.GetByName("右ひじsub先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじsub先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじsub先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.472866, Y: 12.872813, Z: 0.120103}
			if !boneDeltas.GetByName("右手捩1").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩1").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩1").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.458749, Y: 13.009759, Z: -0.086526}
			if !boneDeltas.GetByName("右手捩2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩2").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.440727, Y: 13.184620, Z: -0.350361}
			if !boneDeltas.GetByName("右手捩3").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩3").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩3").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.42368773, Y: 13.34980879, Z: -0.59962077}
			if !boneDeltas.GetByName("右手捩4").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩4").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩4").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.40457204, Y: 13.511055, Z: -0.84384039}
			if !boneDeltas.GetByName("右手捩5").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩5").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩5").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.39275926, Y: 13.62582429, Z: -1.01699954}
			if !boneDeltas.GetByName("右手捩6").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩6").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩6").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.36500465, Y: 13.89623575, Z: -1.42501008}
			if !boneDeltas.GetByName("右手首R").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首R").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手首R").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.36500465, Y: 13.89623575, Z: -1.42501008}
			if !boneDeltas.GetByName("右手首1").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首1").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手首1").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.472418, Y: 13.917203, Z: -1.529887}
			if !boneDeltas.GetByName("右手首2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手首2").FilledGlobalPosition().MMD()))
			}
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

	{

		fno := int(0)
		boneDeltas := DeformBone(motion.BoneFrames, fno, model, nil, true, nil, nil)
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
}

func TestVmdMotion_DeformLegIk25_Addiction_Wa_Right(t *testing.T) {
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

	{

		fno := int(0)
		boneDeltas := DeformBone(motion.BoneFrames, fno, model, []string{"右襟先"}, true, nil, nil)
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
			expectedPosition := &mmath.MVec3{X: -0.630130, Y: 13.306682, Z: 2.752505}
			if !boneDeltas.GetByName(pmx.SHOULDER_P.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Right()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.630131, Y: 13.306683, Z: 2.742505}
			if !boneDeltas.GetByName(pmx.SHOULDER.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Right()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.948004, Y: 13.753115, Z: 2.690539}
			if !boneDeltas.GetByName(pmx.ARM.Right()).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Right()).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Right()).FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.611438, Y: 12.394744, Z: 2.353463}
			if !boneDeltas.GetByName("右上半身C-A").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身C-A").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右上半身C-A").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.664344, Y: 13.835273, Z: 2.403165}
			if !boneDeltas.GetByName("右鎖骨IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右鎖骨IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右鎖骨IK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.270636, Y: 13.350624, Z: 2.258960}
			if !boneDeltas.GetByName("右鎖骨").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右鎖骨").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右鎖骨").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.963317, Y: 14.098928, Z: 2.497183}
			if !boneDeltas.GetByName("右鎖骨先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右鎖骨先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右鎖骨先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.235138, Y: 13.300934, Z: 2.666039}
			if !boneDeltas.GetByName("右肩Rz検").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Rz検").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩Rz検").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.847069, Y: 13.997178, Z: 2.886786}
			if !boneDeltas.GetByName("右肩Rz検先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Rz検先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩Rz検先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.235138, Y: 13.300934, Z: 2.666039}
			if !boneDeltas.GetByName("右肩Ry検").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Ry検").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩Ry検").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -1.172100, Y: 13.315790, Z: 2.838742}
			if !boneDeltas.GetByName("右肩Ry検先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Ry検先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩Ry検先").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.591152, Y: 12.674325, Z: 2.391185}
			if !boneDeltas.GetByName("右上半身C-B").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身C-B").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右上半身C-B").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.588046, Y: 12.954157, Z: 2.432232}
			if !boneDeltas.GetByName("右上半身C-C").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身C-C").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右上半身C-C").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.672292, Y: 10.939227, Z: 2.148515}
			if !boneDeltas.GetByName("右上半身2").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身2").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右上半身2").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.520068, Y: 14.089510, Z: 2.812157}
			if !boneDeltas.GetByName("右襟").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右襟").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右襟").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -0.491354, Y: 14.225309, Z: 2.502640}
			if !boneDeltas.GetByName("右襟先").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右襟先").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右襟先").FilledGlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformIk28_Simple(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/IKの挙動を見たい_020.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/IKの挙動を見たい.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		boneDeltas := DeformBone(motion.BoneFrames, fno, model, nil, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{X: -9.433129, Y: 1.363848, Z: 1.867427}
			if !boneDeltas.GetByName("A+tail").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("A+tail").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("A+tail").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: -9.433129, Y: 1.363847, Z: 1.867427}
			if !boneDeltas.GetByName("A+IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("A+IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("A+IK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 5.0, Y: 4.517528, Z: 2.142881}
			if !boneDeltas.GetByName("B+tail").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("B+tail").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("B+tail").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 0.566871, Y: 1.363847, Z: 1.867427}
			if !boneDeltas.GetByName("B+IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("B+IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("B+IK").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 10.0, Y: 3.020634, Z: 3.984441}
			if !boneDeltas.GetByName("C+tail").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("C+tail").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("C+tail").FilledGlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{X: 5.566871, Y: 1.363848, Z: 1.867427}
			if !boneDeltas.GetByName("C+IK").FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("C+IK").FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("C+IK").FilledGlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformIk29_Simple(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := repository.NewVmdRepository()
	motionData, err := vr.Load("../../../test_resources/IKの挙動を見たい2_040.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*vmd.VmdMotion)

	pr := repository.NewPmxRepository()
	modelData, err := pr.Load("../../../test_resources/IKの挙動を見たい2.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		boneDeltas := DeformBone(motion.BoneFrames, fno, model, nil, true, nil, nil)
		{
			boneName := "A+2"
			expectedPosition := &mmath.MVec3{X: -5.440584, Y: 2.324726, Z: 0.816799}
			if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
			}
		}
		{
			boneName := "A+2tail"
			expectedPosition := &mmath.MVec3{X: -4.671312, Y: 3.980981, Z: -0.895119}
			if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
			}
		}
		{
			boneName := "B+2"
			expectedPosition := &mmath.MVec3{X: 4.559244, Y: 2.324562, Z: 0.817174}
			if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
			}
		}
		{
			boneName := "B+2tail"
			expectedPosition := &mmath.MVec3{X: 5.328533, Y: 3.980770, Z: -0.894783}
			if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
			}
		}
		{
			boneName := "C+2"
			expectedPosition := &mmath.MVec3{X: 8.753987, Y: 2.042284, Z: -0.736314}
			if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
			}
		}
		{
			boneName := "C+2tail"
			expectedPosition := &mmath.MVec3{X: 10.328943, Y: 3.981413, Z: -0.894101}
			if !boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).FilledGlobalPosition().MMD()))
			}
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

	{
		fno := int(0)
		boneDeltas := DeformBone(motion.BoneFrames, fno, model, nil, true, nil, nil)
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
}
