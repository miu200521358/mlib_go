package vmd

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/reader"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

func TestVmdMotion_DeformArmIk4_DMF(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../../test_resources/nac_dmf_601.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &reader.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		boneDeltas := motion.BoneFrames.Deform(fno, model, nil, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{6.210230, 8.439670, 0.496305}
			if !boneDeltas.GetByName(pmx.CENTER.String()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{6.210230, 8.849669, 0.496305}
			if !boneDeltas.GetByName(pmx.GROOVE.String()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{6.210230, 12.836980, -0.159825}
			if !boneDeltas.GetByName("上半身").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("上半身").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{6.261481, 13.968025, 0.288966}
			if !boneDeltas.GetByName("上半身2").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身2").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("上半身2").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{6.541666, 15.754716, 1.421828}
			if !boneDeltas.GetByName("左肩").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左肩").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.451898, 16.031992, 1.675949}
			if !boneDeltas.GetByName("左腕").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.135534, 15.373729, 1.715530}
			if !boneDeltas.GetByName("左腕捩").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.646749, 13.918620, 1.803021}
			if !boneDeltas.GetByName("左ひじ").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじ").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじ").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.164164, 13.503792, 1.706635}
			if !boneDeltas.GetByName("左手捩").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.772219, 12.307291, 1.428628}
			if !boneDeltas.GetByName("左手首").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手首").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.390504, 12.011601, 1.405503}
			if !boneDeltas.GetByName("左手先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.451900, 16.031990, 1.675949}
			if !boneDeltas.GetByName("左腕YZ").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕YZ").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕YZ").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.690105, 15.802624, 1.689741}
			if !boneDeltas.GetByName("左腕YZ先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕YZ先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕YZ先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.690105, 15.802622, 1.689740}
			if !boneDeltas.GetByName("左腕YZIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕YZIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕YZIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.451899, 16.031988, 1.675950}
			if !boneDeltas.GetByName("左腕X").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕X").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕X").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.816861, 16.406412, 1.599419}
			if !boneDeltas.GetByName("左腕X先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕X先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕X先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.816858, 16.406418, 1.599418}
			if !boneDeltas.GetByName("左腕XIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕XIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕XIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.135530, 15.373726, 1.715530}
			if !boneDeltas.GetByName("左腕捩YZ").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩YZ").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩YZ").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.409824, 15.109610, 1.731412}
			if !boneDeltas.GetByName("左腕捩YZTgt").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩YZTgt").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩YZTgt").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.409830, 15.109617, 1.731411}
			if !boneDeltas.GetByName("左腕捩YZIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩YZIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩YZIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.135530, 15.373725, 1.715531}
			if !boneDeltas.GetByName("左腕捩X").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩X").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩X").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.500528, 15.748149, 1.639511}
			if !boneDeltas.GetByName("左腕捩XTgt").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩XTgt").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩XTgt").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.500531, 15.748233, 1.639508}
			if !boneDeltas.GetByName("左腕捩XIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩XIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩XIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.646743, 13.918595, 1.803029}
			if !boneDeltas.GetByName("左ひじYZ").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYZ").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYZ").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.360763, 13.672787, 1.745903}
			if !boneDeltas.GetByName("左ひじYZ先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYZ先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYZ先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.360781, 13.672805, 1.745905}
			if !boneDeltas.GetByName("左ひじYZIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYZIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYZIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.646734, 13.918593, 1.803028}
			if !boneDeltas.GetByName("左ひじX").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじX").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじX").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.944283, 13.652989, 1.456379}
			if !boneDeltas.GetByName("左ひじX先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじX先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじX先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.944304, 13.653007, 1.456381}
			if !boneDeltas.GetByName("左ひじXIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじXIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじXIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.646734, 13.918596, 1.803028}
			if !boneDeltas.GetByName("左ひじY").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじY").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじY").GlobalPosition().MMD()))
			}
		}
		{
			// FIXME
			expectedPosition := &mmath.MVec3{9.560862, 13.926876, 1.431514}
			if !boneDeltas.GetByName("左ひじY先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじY先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじY先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.360781, 13.672805, 1.745905}
			if !boneDeltas.GetByName("左ひじYIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじYIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじYIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.164141, 13.503780, 1.706625}
			if !boneDeltas.GetByName("左手捩YZ").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩YZ").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩YZ").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.414344, 12.859288, 1.556843}
			if !boneDeltas.GetByName("左手捩YZTgt").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩YZTgt").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩YZTgt").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.414370, 12.859282, 1.556885}
			if !boneDeltas.GetByName("左手捩YZIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩YZIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩YZIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.164142, 13.503780, 1.706624}
			if !boneDeltas.GetByName("左手捩X").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩X").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩X").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.511073, 12.928087, 2.447041}
			if !boneDeltas.GetByName("左手捩XTgt").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩XTgt").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩XTgt").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.511120, 12.928122, 2.447057}
			if !boneDeltas.GetByName("左手捩XIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩XIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩XIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.471097, 12.074032, 1.410383}
			if !boneDeltas.GetByName("左手YZ先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手YZ先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手YZ先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.471111, 12.074042, 1.410384}
			if !boneDeltas.GetByName("左手YZIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手YZIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手YZIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.772183, 12.307314, 1.428564}
			if !boneDeltas.GetByName("左手X").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手X").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手X").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.802912, 12.308764, 0.901022}
			if !boneDeltas.GetByName("左手X先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手X先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手X先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.802991, 12.308830, 0.901079}
			if !boneDeltas.GetByName("左手XIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手XIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手XIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.130125, 15.368912, 1.728851}
			if !boneDeltas.GetByName("左腕捩1").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩1").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩1").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.365511, 15.142246, 1.742475}
			if !boneDeltas.GetByName("左腕捩2").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩2").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩2").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.639965, 14.877952, 1.758356}
			if !boneDeltas.GetByName("左腕捩3").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩3").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩3").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.901615, 14.625986, 1.773497}
			if !boneDeltas.GetByName("左腕捩4").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左腕捩4").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左腕捩4").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.641270, 13.913721, 1.816324}
			if !boneDeltas.GetByName("左ひじsub").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじsub").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじsub").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.907782, 13.661371, 2.034630}
			if !boneDeltas.GetByName("左ひじsub先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左ひじsub先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左ひじsub先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{9.165060, 13.499348, 1.721094}
			if !boneDeltas.GetByName("左手捩1").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩1").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩1").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.978877, 13.339340, 1.683909}
			if !boneDeltas.GetByName("左手捩2").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩2").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩2").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.741154, 13.135028, 1.636428}
			if !boneDeltas.GetByName("左手捩3").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩3").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩3").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.516553, 12.942023, 1.591578}
			if !boneDeltas.GetByName("左手捩4").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩4").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩4").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.301016, 12.748707, 1.544439}
			if !boneDeltas.GetByName("左手捩5").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩5").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩5").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{8.145000, 12.614601, 1.513277}
			if !boneDeltas.GetByName("左手捩6").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手捩6").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手捩6").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.777408, 12.298634, 1.439762}
			if !boneDeltas.GetByName("左手首R").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首R").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手首R").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.777408, 12.298635, 1.439762}
			if !boneDeltas.GetByName("左手首1").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首1").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手首1").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{7.670320, 12.202144, 1.486689}
			if !boneDeltas.GetByName("左手首2").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左手首2").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左手首2").GlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformArmIk2(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("C:/MMD/mmd_base/tests/resources/唱(ダンスのみ)_0274F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &reader.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/VOCALOID/初音ミク/ISAO式ミク/I_ミクv4/Miku_V4_準標準.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		boneDeltas := motion.BoneFrames.Deform(fno, model, nil, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{0.04952335, 9.0, 1.72378033}
			if !boneDeltas.GetByName(pmx.CENTER.String()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.CENTER.String()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.CENTER.String()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 7.97980869, 1.72378033}
			if !boneDeltas.GetByName(pmx.GROOVE.String()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.GROOVE.String()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.GROOVE.String()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 11.02838314, 2.29172656}
			if !boneDeltas.GetByName("腰").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("腰").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("腰").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.04952335, 11.9671191, 1.06765032}
			if !boneDeltas.GetByName("上半身").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("上半身").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.26284261, 13.14576297, 0.84720008}
			if !boneDeltas.GetByName("上半身2").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("上半身2").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("上半身2").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.33636433, 15.27729547, 0.77435588}
			if !boneDeltas.GetByName("右肩").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !boneDeltas.GetByName("右肩C").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩C").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩C").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !boneDeltas.GetByName("右腕").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90326269, 14.53727204, 0.7925801}
			if !boneDeltas.GetByName("右腕捩").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50502977, 12.52976106, 0.66393998}
			if !boneDeltas.GetByName("右ひじ").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじ").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじ").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.46843236, 12.88476121, 0.12831076}
			if !boneDeltas.GetByName("右手捩").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36287259, 13.90869981, -1.41662258}
			if !boneDeltas.GetByName("右手首").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手首").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.81521586, 14.00661535, -1.55616424}
			if !boneDeltas.GetByName("右手先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !boneDeltas.GetByName("右腕YZ").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕YZ").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕YZ").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.72589296, 15.12898892, 0.83049645}
			if !boneDeltas.GetByName("右腕YZ先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕YZ先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕YZ先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.72589374, 15.12898632, 0.83049628}
			if !boneDeltas.GetByName("右腕YZIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕YZIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕YZIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.63104276, 15.44542768, 0.8507726}
			if !boneDeltas.GetByName("右腕X").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕X").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕X").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.125321, 15.600293, 0.746130}
			if !boneDeltas.GetByName("右腕X先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕X先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕X先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.1253241, 15.60029489, 0.7461294}
			if !boneDeltas.GetByName("右腕XIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕XIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕XIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90325538, 14.53727326, 0.79258165}
			if !boneDeltas.GetByName("右腕捩YZ").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩YZ").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩YZ").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.01247534, 14.17289417, 0.76923367}
			if !boneDeltas.GetByName("右腕捩YZTgt").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩YZTgt").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩YZTgt").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.01248754, 14.17289597, 0.76923112}
			if !boneDeltas.GetByName("右腕捩YZIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩YZIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩YZIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90325538, 14.53727326, 0.79258165}
			if !boneDeltas.GetByName("右腕捩X").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩X").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩X").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.40656426, 14.68386802, 0.85919594}
			if !boneDeltas.GetByName("右腕捩XTgt").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩XTgt").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩XTgt").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.40657579, 14.68387899, 0.8591982}
			if !boneDeltas.GetByName("右腕捩XIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩XIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩XIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50499623, 12.52974836, 0.66394738}
			if !boneDeltas.GetByName("右ひじYZ").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYZ").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYZ").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.48334366, 12.74011791, 0.34655051}
			if !boneDeltas.GetByName("右ひじYZ先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYZ先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYZ先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.48334297, 12.74012453, 0.34654052}
			if !boneDeltas.GetByName("右ひじYZIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYZIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYZIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50499623, 12.52974836, 0.66394738}
			if !boneDeltas.GetByName("右ひじX").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじX").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじX").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.01179616, 12.66809052, 0.72106658}
			if !boneDeltas.GetByName("右ひじX先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじX先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじX先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-2.00760407, 12.67958516, 0.7289003}
			if !boneDeltas.GetByName("右ひじXIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじXIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじXIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.50499623, 12.52974836, 0.66394738}
			if !boneDeltas.GetByName("右ひじY").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじY").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじY").GlobalPosition().MMD()))
			}
		}
		{

			expectedPosition := &mmath.MVec3{-1.485519, 12.740760, 0.346835}
			if !boneDeltas.GetByName("右ひじY先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじY先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじY先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.48334297, 12.74012453, 0.34654052}
			if !boneDeltas.GetByName("右ひじYIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじYIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじYIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.46845628, 12.88475892, 0.12832214}
			if !boneDeltas.GetByName("右手捩YZ").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩YZ").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩YZ").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.41168478, 13.4363328, -0.7038697}
			if !boneDeltas.GetByName("右手捩YZTgt").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩YZTgt").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩YZTgt").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.41156715, 13.43632015, -0.70389025}
			if !boneDeltas.GetByName("右手捩YZIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩YZIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩YZIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.46845628, 12.88475892, 0.12832214}
			if !boneDeltas.GetByName("右手捩X").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩X").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩X").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.5965686, 12.06213832, -0.42564769}
			if !boneDeltas.GetByName("右手捩XTgt").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩XTgt").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩XTgt").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.5965684, 12.06214091, -0.42565404}
			if !boneDeltas.GetByName("右手捩XIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩XIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩XIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.7198605, 13.98597326, -1.5267472}
			if !boneDeltas.GetByName("右手YZ先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手YZ先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手YZ先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.71969424, 13.98593727, -1.52669587}
			if !boneDeltas.GetByName("右手YZIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手YZIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手YZIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36306295, 13.90872698, -1.41659848}
			if !boneDeltas.GetByName("右手X").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手X").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手X").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.54727182, 13.56147176, -1.06342964}
			if !boneDeltas.GetByName("右手X先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手X先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手X先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.54700171, 13.5614545, -1.0633896}
			if !boneDeltas.GetByName("右手XIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手XIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手XIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.90581859, 14.5370842, 0.80752276}
			if !boneDeltas.GetByName("右腕捩1").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩1").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩1").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.99954005, 14.2243783, 0.78748743}
			if !boneDeltas.GetByName("右腕捩2").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩2").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩2").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.10880907, 13.85976329, 0.76412793}
			if !boneDeltas.GetByName("右腕捩3").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩3").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩3").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.21298069, 13.51216081, 0.74185819}
			if !boneDeltas.GetByName("右腕捩4").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右腕捩4").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右腕捩4").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.5074743, 12.52953348, 0.67889319}
			if !boneDeltas.GetByName("右ひじsub").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじsub").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじsub").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.617075, 12.131149, 0.786797}
			if !boneDeltas.GetByName("右ひじsub先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右ひじsub先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右ひじsub先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.472866, 12.872813, 0.120103}
			if !boneDeltas.GetByName("右手捩1").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩1").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩1").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.458749, 13.009759, -0.086526}
			if !boneDeltas.GetByName("右手捩2").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩2").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩2").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.440727, 13.184620, -0.350361}
			if !boneDeltas.GetByName("右手捩3").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩3").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩3").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.42368773, 13.34980879, -0.59962077}
			if !boneDeltas.GetByName("右手捩4").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩4").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩4").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.40457204, 13.511055, -0.84384039}
			if !boneDeltas.GetByName("右手捩5").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩5").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩5").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.39275926, 13.62582429, -1.01699954}
			if !boneDeltas.GetByName("右手捩6").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手捩6").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手捩6").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36500465, 13.89623575, -1.42501008}
			if !boneDeltas.GetByName("右手首R").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首R").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手首R").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.36500465, 13.89623575, -1.42501008}
			if !boneDeltas.GetByName("右手首1").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首1").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手首1").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.472418, 13.917203, -1.529887}
			if !boneDeltas.GetByName("右手首2").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右手首2").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右手首2").GlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformLegIk25_Addiction_Wa_Left(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../../test_resources/[A]ddiction_和洋_0126F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &reader.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		boneDeltas := motion.BoneFrames.Deform(fno, model, nil, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{-0.225006, 9.705784, 2.033072}
			if !boneDeltas.GetByName(pmx.UPPER.String()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER.String()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.237383, 10.769137, 2.039952}
			if !boneDeltas.GetByName(pmx.UPPER2.String()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER2.String()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.460140, 13.290816, 2.531440}
			if !boneDeltas.GetByName(pmx.SHOULDER_P.Left()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Left()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Left()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.460140, 13.290816, 2.531440}
			if !boneDeltas.GetByName(pmx.SHOULDER.Left()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Left()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Left()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.784452, 13.728909, 2.608527}
			if !boneDeltas.GetByName(pmx.ARM.Left()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Left()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Left()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.272067, 12.381887, 2.182425}
			if !boneDeltas.GetByName("左上半身C-A").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身C-A").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左上半身C-A").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.406217, 13.797803, 2.460243}
			if !boneDeltas.GetByName("左鎖骨IK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左鎖骨IK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左鎖骨IK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.052427, 13.347448, 2.216718}
			if !boneDeltas.GetByName("左鎖骨").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左鎖骨").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左鎖骨").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.659554, 14.017852, 2.591099}
			if !boneDeltas.GetByName("左鎖骨先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左鎖骨先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左鎖骨先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.065147, 13.296564, 2.607907}
			if !boneDeltas.GetByName("左肩Rz検").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Rz検").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左肩Rz検").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.517776, 14.134196, 2.645912}
			if !boneDeltas.GetByName("左肩Rz検先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Rz検先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左肩Rz検先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.065148, 13.296564, 2.607907}
			if !boneDeltas.GetByName("左肩Ry検").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Ry検").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左肩Ry検").GlobalPosition().MMD()))
			}
		}
		{
			// FIXME
			expectedPosition := &mmath.MVec3{0.860159, 13.190875, 3.122428}
			if !boneDeltas.GetByName("左肩Ry検先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左肩Ry検先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左肩Ry検先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.195053, 12.648546, 2.236849}
			if !boneDeltas.GetByName("左上半身C-B").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身C-B").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左上半身C-B").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.294257, 12.912640, 2.257159}
			if !boneDeltas.GetByName("左上半身C-C").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身C-C").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左上半身C-C").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.210011, 10.897711, 1.973442}
			if !boneDeltas.GetByName("左上半身2").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左上半身2").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左上半身2").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.320589, 14.049745, 2.637018}
			if !boneDeltas.GetByName("左襟").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左襟").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左襟").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.297636, 14.263302, 2.374467}
			if !boneDeltas.GetByName("左襟先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左襟先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左襟先").GlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformLegIk25_Addiction_Wa_Right(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../../test_resources/[A]ddiction_和洋_0126F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &reader.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		boneDeltas := motion.BoneFrames.Deform(fno, model, []string{"右襟先"}, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{-0.225006, 9.705784, 2.033072}
			if !boneDeltas.GetByName(pmx.UPPER.String()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER.String()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER.String()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.237383, 10.769137, 2.039952}
			if !boneDeltas.GetByName(pmx.UPPER2.String()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.UPPER2.String()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.UPPER2.String()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.630130, 13.306682, 2.752505}
			if !boneDeltas.GetByName(pmx.SHOULDER_P.Right()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER_P.Right()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER_P.Right()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.630131, 13.306683, 2.742505}
			if !boneDeltas.GetByName(pmx.SHOULDER.Right()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.SHOULDER.Right()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.SHOULDER.Right()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.948004, 13.753115, 2.690539}
			if !boneDeltas.GetByName(pmx.ARM.Right()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ARM.Right()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ARM.Right()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.611438, 12.394744, 2.353463}
			if !boneDeltas.GetByName("右上半身C-A").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身C-A").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右上半身C-A").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.664344, 13.835273, 2.403165}
			if !boneDeltas.GetByName("右鎖骨IK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右鎖骨IK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右鎖骨IK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.270636, 13.350624, 2.258960}
			if !boneDeltas.GetByName("右鎖骨").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右鎖骨").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右鎖骨").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.963317, 14.098928, 2.497183}
			if !boneDeltas.GetByName("右鎖骨先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右鎖骨先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右鎖骨先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.235138, 13.300934, 2.666039}
			if !boneDeltas.GetByName("右肩Rz検").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Rz検").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩Rz検").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.847069, 13.997178, 2.886786}
			if !boneDeltas.GetByName("右肩Rz検先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Rz検先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩Rz検先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.235138, 13.300934, 2.666039}
			if !boneDeltas.GetByName("右肩Ry検").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Ry検").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩Ry検").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.172100, 13.315790, 2.838742}
			if !boneDeltas.GetByName("右肩Ry検先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右肩Ry検先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右肩Ry検先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.591152, 12.674325, 2.391185}
			if !boneDeltas.GetByName("右上半身C-B").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身C-B").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右上半身C-B").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.588046, 12.954157, 2.432232}
			if !boneDeltas.GetByName("右上半身C-C").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身C-C").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右上半身C-C").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.672292, 10.939227, 2.148515}
			if !boneDeltas.GetByName("右上半身2").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右上半身2").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右上半身2").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.520068, 14.089510, 2.812157}
			if !boneDeltas.GetByName("右襟").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右襟").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右襟").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.491354, 14.225309, 2.502640}
			if !boneDeltas.GetByName("右襟先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("右襟先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("右襟先").GlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformIk28_Simple(t *testing.T) {
	mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../../test_resources/IKの挙動を見たい_020.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &reader.PmxReader{}
	modelData, err := pr.ReadByFilepath("../../../test_resources/IKの挙動を見たい.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		boneDeltas := motion.BoneFrames.Deform(fno, model, nil, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{-9.433129, 1.363848, 1.867427}
			if !boneDeltas.GetByName("A+tail").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("A+tail").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("A+tail").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-9.433129, 1.363847, 1.867427}
			if !boneDeltas.GetByName("A+IK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("A+IK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("A+IK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{5.0, 4.517528, 2.142881}
			if !boneDeltas.GetByName("B+tail").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("B+tail").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("B+tail").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.566871, 1.363847, 1.867427}
			if !boneDeltas.GetByName("B+IK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("B+IK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("B+IK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{10.0, 3.020634, 3.984441}
			if !boneDeltas.GetByName("C+tail").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("C+tail").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("C+tail").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{5.566871, 1.363848, 1.867427}
			if !boneDeltas.GetByName("C+IK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("C+IK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("C+IK").GlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformIk29_Simple(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../../test_resources/IKの挙動を見たい2_040.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &reader.PmxReader{}
	modelData, err := pr.ReadByFilepath("../../../test_resources/IKの挙動を見たい2.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{

		fno := int(0)
		boneDeltas := motion.BoneFrames.Deform(fno, model, nil, true, nil, nil)
		{
			boneName := "A+2"
			expectedPosition := &mmath.MVec3{-5.440584, 2.324726, 0.816799}
			if !boneDeltas.GetByName(boneName).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).GlobalPosition().MMD()))
			}
		}
		{
			boneName := "A+2tail"
			expectedPosition := &mmath.MVec3{-4.671312, 3.980981, -0.895119}
			if !boneDeltas.GetByName(boneName).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).GlobalPosition().MMD()))
			}
		}
		{
			boneName := "B+2"
			expectedPosition := &mmath.MVec3{4.559244, 2.324562, 0.817174}
			if !boneDeltas.GetByName(boneName).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).GlobalPosition().MMD()))
			}
		}
		{
			boneName := "B+2tail"
			expectedPosition := &mmath.MVec3{5.328533, 3.980770, -0.894783}
			if !boneDeltas.GetByName(boneName).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).GlobalPosition().MMD()))
			}
		}
		{
			boneName := "C+2"
			expectedPosition := &mmath.MVec3{8.753987, 2.042284, -0.736314}
			if !boneDeltas.GetByName(boneName).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).GlobalPosition().MMD()))
			}
		}
		{
			boneName := "C+2tail"
			expectedPosition := &mmath.MVec3{10.328943, 3.981413, -0.894101}
			if !boneDeltas.GetByName(boneName).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(boneName).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(boneName).GlobalPosition().MMD()))
			}
		}
	}
}

func TestVmdMotion_DeformLegIk30_Addiction_Shoes(t *testing.T) {
	// mlog.SetLevel(mlog.IK_VERBOSE)

	vr := &VmdMotionReader{}
	motionData, err := vr.ReadByFilepath("../../../test_resources/[A]ddiction_和洋_1037F.vmd")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	motion := motionData.(*VmdMotion)

	pr := &reader.PmxReader{}
	modelData, err := pr.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20231028/wa_129cm_20240406.pmx")

	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	model := modelData.(*pmx.PmxModel)

	{
		fno := int(0)
		boneDeltas := motion.BoneFrames.Deform(fno, model, nil, true, nil, nil)
		{
			expectedPosition := &mmath.MVec3{0, 0, 0}
			if !boneDeltas.GetByName(pmx.ROOT.String()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.ROOT.String()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.ROOT.String()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.748025, 1.683590, 0.556993}
			if !boneDeltas.GetByName(pmx.LEG_IK.Left()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG_IK.Right()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG_IK.Right()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.111190, 4.955496, 1.070225}
			if !boneDeltas.GetByName(pmx.LOWER.String()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LOWER.String()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LOWER.String()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.724965, 3.674735, 0.810759}
			if !boneDeltas.GetByName(pmx.LEG.Left()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.LEG.Right()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.LEG.Right()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.924419, 3.943549, -1.420897}
			if !boneDeltas.GetByName(pmx.KNEE.Left()).GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName(pmx.KNEE.Right()).GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName(pmx.KNEE.Right()).GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.743128, 1.672784, 0.551317}
			if !boneDeltas.GetByName("左脛骨").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.743128, 1.672784, 0.551317}
			if !boneDeltas.GetByName("左脛骨D").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.118480, 3.432016, -1.657329}
			if !boneDeltas.GetByName("左脛骨D先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左脛骨D先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左脛骨D先").GlobalPosition().MMD()))
			}
		}
		{
			// FIXME
			expectedPosition := &mmath.MVec3{1.763123, 3.708842, 0.369619}
			if !boneDeltas.GetByName("左足Dw").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足Dw").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足Dw").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.698137, 0.674393, 0.043128}
			if !boneDeltas.GetByName("左足先EX").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足先EX").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足先EX").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.712729, 0.919505, 0.174835}
			if !boneDeltas.GetByName("左素足先A").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.494684, 0.328695, -0.500715}
			if !boneDeltas.GetByName("左素足先A先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先A先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先A先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.494684, 0.328695, -0.500715}
			if !boneDeltas.GetByName("左素足先AIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先AIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先AIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.224127, 1.738605, 0.240794}
			if !boneDeltas.GetByName("左素足先B").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.582481, 0.627567, -0.222556}
			if !boneDeltas.GetByName("左素足先B先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先B先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先B先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.572457, 0.658645, -0.209595}
			if !boneDeltas.GetByName("左素足先BIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左素足先BIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左素足先BIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.448783, 1.760351, 0.405702}
			if !boneDeltas.GetByName("左靴調節").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴調節").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴調節").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.875421, 0.917357, 0.652144}
			if !boneDeltas.GetByName("左靴追従").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.404960, 1.846940, 0.380388}
			if !boneDeltas.GetByName("左靴追従先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.448783, 1.760351, 0.405702}
			if !boneDeltas.GetByName("左靴追従IK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左靴追従IK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左靴追従IK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.330700, 3.021906, 0.153679}
			if !boneDeltas.GetByName("左足補D").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補D").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.625984, 3.444340, -1.272553}
			if !boneDeltas.GetByName("左足補D先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補D先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補D先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.625984, 3.444340, -1.272553}
			if !boneDeltas.GetByName("左足補DIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足補DIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足補DIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.724964, 3.674735, 0.810759}
			if !boneDeltas.GetByName("左足向検A").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.924420, 3.943550, -1.420897}
			if !boneDeltas.GetByName("左足向検A先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検A先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検A先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.924419, 3.943550, -1.420896}
			if !boneDeltas.GetByName("左足向検AIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向検AIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向検AIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.724965, 3.674735, 0.810760}
			if !boneDeltas.GetByName("左足向-").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足向-").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足向-").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{1.566813, 3.645544, 0.177956}
			if !boneDeltas.GetByName("左足w").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足w").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-1.879506, 3.670895, -2.715526}
			if !boneDeltas.GetByName("左足w先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足w先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足w先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.631964, 3.647012, -1.211210}
			if !boneDeltas.GetByName("左膝補").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.83923, 2.876222, -1.494900}
			if !boneDeltas.GetByName("左膝補先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.686400, 3.444961, -1.285575}
			if !boneDeltas.GetByName("左膝補IK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左膝補IK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左膝補IK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.924420, 3.943550, -1.420896}
			if !boneDeltas.GetByName("左足捩検B").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.952927, 7.388766, -0.972059}
			if !boneDeltas.GetByName("左足捩検B先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検B先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検B先").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{-0.952927, 7.388766, -0.972060}
			if !boneDeltas.GetByName("左足捩検BIK").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩検BIK").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩検BIK").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.400272, 3.809143, -0.305068}
			if !boneDeltas.GetByName("左足捩").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩").GlobalPosition().MMD()))
			}
		}
		{
			expectedPosition := &mmath.MVec3{0.371963, 4.256704, -0.267830}
			if !boneDeltas.GetByName("左足捩先").GlobalPosition().MMD().NearEquals(expectedPosition, 0.03) {
				t.Errorf("Expected %v, got %v (%.3f)", expectedPosition, boneDeltas.GetByName("左足捩先").GlobalPosition().MMD(), expectedPosition.Distance(boneDeltas.GetByName("左足捩先").GlobalPosition().MMD()))
			}
		}
	}
}
