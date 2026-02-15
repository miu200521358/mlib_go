// 指示: miu200521358
package ui

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/miu200521358/mlib_go/cmd/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_csv"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

const (
	motionCsvBoneSuffix  = "bone"
	motionCsvMorphSuffix = "morph"
	motionCsvTimeLayout  = "20060102_150405"
	motionCsvExt         = ".csv"
	motionVmdExt         = ".vmd"
)

const (
	motionCsvCsvSaveFailedErrorID = "15607"
	motionCsvMarshalFailedErrorID = "95602"
)

var motionCsvNowFunc = time.Now

// exportMotionCsv はVMDモーションをボーン/モーフCSVとして出力する。
func exportMotionCsv(inputPath string, motionData *motion.VmdMotion) error {
	if inputPath == "" || motionData == nil {
		return nil
	}
	if !strings.EqualFold(filepath.Ext(inputPath), motionVmdExt) {
		return nil
	}

	basePath := buildMotionCsvBasePathByMotionPath(inputPath)
	return exportMotionCsvByBasePath(basePath, motionData)
}

// exportMotionCsvByOutputPath はCSV出力パスを基準にボーン/モーフCSVを出力する。
func exportMotionCsvByOutputPath(outputPath string, motionData *motion.VmdMotion) error {
	basePath := buildMotionCsvBasePathByOutputPath(outputPath)
	return exportMotionCsvByBasePath(basePath, motionData)
}

// buildMotionCsvDefaultOutputPath は入力モーションパスからCSV出力先の既定値を生成する。
func buildMotionCsvDefaultOutputPath(inputPath string) string {
	if inputPath == "" {
		return ""
	}
	return buildMotionCsvBasePathByMotionPath(inputPath) + motionCsvExt
}

// exportMotionCsvByBasePath はCSV出力ベースパスからボーン/モーフCSVを出力する。
func exportMotionCsvByBasePath(basePath string, motionData *motion.VmdMotion) error {
	if basePath == "" || motionData == nil {
		return nil
	}

	timestamp := motionCsvNowFunc().Format(motionCsvTimeLayout)
	var exportErr error

	if motionData.BoneFrames != nil && motionData.BoneFrames.Len() > 0 {
		bonePath := buildMotionCsvOutputPath(basePath, motionCsvBoneSuffix, timestamp)
		if err := saveMotionBoneCsv(bonePath, buildMotionBoneCsvRows(motionData)); err != nil {
			exportErr = errors.Join(exportErr, err)
		}
	}

	if motionData.MorphFrames != nil && motionData.MorphFrames.Len() > 0 {
		morphPath := buildMotionCsvOutputPath(basePath, motionCsvMorphSuffix, timestamp)
		if err := saveMotionMorphCsv(morphPath, buildMotionMorphCsvRows(motionData)); err != nil {
			exportErr = errors.Join(exportErr, err)
		}
	}

	return exportErr
}

// buildMotionCsvBasePathByMotionPath は入力モーションパスからCSV出力ベースパスを生成する。
func buildMotionCsvBasePathByMotionPath(inputPath string) string {
	dir := filepath.Dir(inputPath)
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	return filepath.Join(dir, base)
}

// buildMotionCsvBasePathByOutputPath はCSV出力パスからCSV出力ベースパスを生成する。
func buildMotionCsvBasePathByOutputPath(outputPath string) string {
	if outputPath == "" {
		return ""
	}
	ext := filepath.Ext(outputPath)
	if strings.EqualFold(ext, motionCsvExt) {
		return strings.TrimSuffix(outputPath, ext)
	}
	return outputPath
}

// buildMotionCsvOutputPath はCSV出力ベースパスからCSV出力パスを生成する。
func buildMotionCsvOutputPath(basePath, kind, timestamp string) string {
	return fmt.Sprintf("%s_%s_%s%s", basePath, kind, timestamp, motionCsvExt)
}

// saveMotionBoneCsv はボーンCSVを保存する。
func saveMotionBoneCsv(outputPath string, rows []motionBoneCsvRow) error {
	model, err := io_csv.Marshal(rows)
	if err != nil {
		return merr.NewCommonError(
			motionCsvMarshalFailedErrorID,
			merr.ErrorKindInternal,
			messages.MessageMotionCsvBoneMarshalFailed,
			err,
		)
	}
	if err := io_csv.NewCsvRepository().Save(outputPath, model, io_common.SaveOptions{}); err != nil {
		return merr.NewCommonError(
			motionCsvCsvSaveFailedErrorID,
			merr.ErrorKindValidate,
			messages.MessageMotionCsvBoneSaveFailed,
			err,
			filepath.Base(outputPath),
		)
	}
	return nil
}

// saveMotionMorphCsv はモーフCSVを保存する。
func saveMotionMorphCsv(outputPath string, rows []motionMorphCsvRow) error {
	model, err := io_csv.Marshal(rows)
	if err != nil {
		return merr.NewCommonError(
			motionCsvMarshalFailedErrorID,
			merr.ErrorKindInternal,
			messages.MessageMotionCsvMorphMarshalFailed,
			err,
		)
	}
	if err := io_csv.NewCsvRepository().Save(outputPath, model, io_common.SaveOptions{}); err != nil {
		return merr.NewCommonError(
			motionCsvCsvSaveFailedErrorID,
			merr.ErrorKindValidate,
			messages.MessageMotionCsvMorphSaveFailed,
			err,
			filepath.Base(outputPath),
		)
	}
	return nil
}

// buildMotionBoneCsvRows はボーンフレーム群をCSV行へ変換する。
func buildMotionBoneCsvRows(motionData *motion.VmdMotion) []motionBoneCsvRow {
	if motionData == nil || motionData.BoneFrames == nil {
		return []motionBoneCsvRow{}
	}
	rows := make([]motionBoneCsvRow, 0, motionData.BoneFrames.Len())
	for _, name := range motionData.BoneFrames.Names() {
		boneFrames := motionData.BoneFrames.Get(name)
		if boneFrames == nil {
			continue
		}
		boneFrames.ForEach(func(frame motion.Frame, value *motion.BoneFrame) bool {
			rows = append(rows, newMotionBoneCsvRow(name, frame, value))
			return true
		})
	}
	return rows
}

// buildMotionMorphCsvRows はモーフフレーム群をCSV行へ変換する。
func buildMotionMorphCsvRows(motionData *motion.VmdMotion) []motionMorphCsvRow {
	if motionData == nil || motionData.MorphFrames == nil {
		return []motionMorphCsvRow{}
	}
	rows := make([]motionMorphCsvRow, 0, motionData.MorphFrames.Len())
	for _, name := range motionData.MorphFrames.Names() {
		morphFrames := motionData.MorphFrames.Get(name)
		if morphFrames == nil {
			continue
		}
		morphFrames.ForEach(func(frame motion.Frame, value *motion.MorphFrame) bool {
			if value == nil {
				return true
			}
			rows = append(rows, motionMorphCsvRow{
				MorphName: name,
				Frame:     int(frame),
				Size:      value.Ratio,
			})
			return true
		})
	}
	return rows
}

// newMotionBoneCsvRow はボーンフレームをCSV行へ変換する。
func newMotionBoneCsvRow(name string, frameIndex motion.Frame, frame *motion.BoneFrame) motionBoneCsvRow {
	positionX := 0.0
	positionY := 0.0
	positionZ := 0.0
	if frame != nil && frame.Position != nil {
		positionX = frame.Position.X
		positionY = frame.Position.Y
		positionZ = frame.Position.Z
	}

	rotationX := 0.0
	rotationY := 0.0
	rotationZ := 0.0
	if frame != nil && frame.Rotation != nil {
		degrees := frame.Rotation.ToMMDDegrees()
		rotationX = degrees.X
		rotationY = degrees.Y
		rotationZ = degrees.Z
	}

	interpolation := buildMotionBoneInterpolation(frame)

	return motionBoneCsvRow{
		BoneName:  name,
		Frame:     int(frameIndex),
		PositionX: positionX,
		PositionY: positionY,
		PositionZ: positionZ,
		RotationX: rotationX,
		RotationY: rotationY,
		RotationZ: rotationZ,
		Interp00:  int(interpolation[0]),
		Interp01:  int(interpolation[1]),
		Interp02:  int(interpolation[2]),
		Interp03:  int(interpolation[3]),
		Interp04:  int(interpolation[4]),
		Interp05:  int(interpolation[5]),
		Interp06:  int(interpolation[6]),
		Interp07:  int(interpolation[7]),
		Interp08:  int(interpolation[8]),
		Interp09:  int(interpolation[9]),
		Interp10:  int(interpolation[10]),
		Interp11:  int(interpolation[11]),
		Interp12:  int(interpolation[12]),
		Interp13:  int(interpolation[13]),
		Interp14:  int(interpolation[14]),
		Interp15:  int(interpolation[15]),
		Interp16:  int(interpolation[16]),
		Interp17:  int(interpolation[17]),
		Interp18:  int(interpolation[18]),
		Interp19:  int(interpolation[19]),
		Interp20:  int(interpolation[20]),
		Interp21:  int(interpolation[21]),
		Interp22:  int(interpolation[22]),
		Interp23:  int(interpolation[23]),
		Interp24:  int(interpolation[24]),
		Interp25:  int(interpolation[25]),
		Interp26:  int(interpolation[26]),
		Interp27:  int(interpolation[27]),
		Interp28:  int(interpolation[28]),
		Interp29:  int(interpolation[29]),
		Interp30:  int(interpolation[30]),
		Interp31:  int(interpolation[31]),
		Interp32:  int(interpolation[32]),
		Interp33:  int(interpolation[33]),
		Interp34:  int(interpolation[34]),
		Interp35:  int(interpolation[35]),
		Interp36:  int(interpolation[36]),
		Interp37:  int(interpolation[37]),
		Interp38:  int(interpolation[38]),
		Interp39:  int(interpolation[39]),
		Interp40:  int(interpolation[40]),
		Interp41:  int(interpolation[41]),
		Interp42:  int(interpolation[42]),
		Interp43:  int(interpolation[43]),
		Interp44:  int(interpolation[44]),
		Interp45:  int(interpolation[45]),
		Interp46:  int(interpolation[46]),
		Interp47:  int(interpolation[47]),
		Interp48:  int(interpolation[48]),
		Interp49:  int(interpolation[49]),
		Interp50:  int(interpolation[50]),
		Interp51:  int(interpolation[51]),
		Interp52:  int(interpolation[52]),
		Interp53:  int(interpolation[53]),
		Interp54:  int(interpolation[54]),
		Interp55:  int(interpolation[55]),
		Interp56:  int(interpolation[56]),
		Interp57:  int(interpolation[57]),
		Interp58:  int(interpolation[58]),
		Interp59:  int(interpolation[59]),
		Interp60:  int(interpolation[60]),
		Interp61:  int(interpolation[61]),
		Interp62:  int(interpolation[62]),
		Interp63:  int(interpolation[63]),
	}
}

// buildMotionBoneInterpolation はボーン補間値64要素を返す。
func buildMotionBoneInterpolation(frame *motion.BoneFrame) [64]byte {
	var interpolation [64]byte
	copy(interpolation[:], motion.INITIAL_BONE_CURVES)
	if frame == nil || frame.Curves == nil || len(frame.Curves.Values) == 0 {
		return interpolation
	}
	copy(interpolation[:], frame.Curves.Values)
	return interpolation
}

// motionBoneCsvRow はボーンCSVの1行を表す。
type motionBoneCsvRow struct {
	BoneName  string  `csv:"ボーン名"`
	Frame     int     `csv:"フレーム"`
	PositionX float64 `csv:"位置X"`
	PositionY float64 `csv:"位置Y"`
	PositionZ float64 `csv:"位置Z"`
	RotationX float64 `csv:"回転X"`
	RotationY float64 `csv:"回転Y"`
	RotationZ float64 `csv:"回転Z"`
	Interp00  int     `csv:"【X_x1】"`
	Interp01  int     `csv:"Y_x1"`
	Interp02  int     `csv:"Z_x1"`
	Interp03  int     `csv:"R_x1"`
	Interp04  int     `csv:"【X_y1】"`
	Interp05  int     `csv:"Y_y1"`
	Interp06  int     `csv:"Z_y1"`
	Interp07  int     `csv:"R_y1"`
	Interp08  int     `csv:"【X_x2】"`
	Interp09  int     `csv:"Y_x2"`
	Interp10  int     `csv:"Z_x2"`
	Interp11  int     `csv:"R_x2"`
	Interp12  int     `csv:"【X_y2】"`
	Interp13  int     `csv:"Y_y2"`
	Interp14  int     `csv:"Z_y2"`
	Interp15  int     `csv:"R_y2"`
	Interp16  int     `csv:"【Y_x1】"`
	Interp17  int     `csv:"Z_x1"`
	Interp18  int     `csv:"R_x1"`
	Interp19  int     `csv:"X_y1"`
	Interp20  int     `csv:"【Y_y1】"`
	Interp21  int     `csv:"Z_y1"`
	Interp22  int     `csv:"R_y1"`
	Interp23  int     `csv:"X_x2"`
	Interp24  int     `csv:"【Y_x2】"`
	Interp25  int     `csv:"Z_x2"`
	Interp26  int     `csv:"R_x2"`
	Interp27  int     `csv:"X_y2"`
	Interp28  int     `csv:"【Y_y2】"`
	Interp29  int     `csv:"Z_y2"`
	Interp30  int     `csv:"R_y2"`
	Interp31  int     `csv:"1"`
	Interp32  int     `csv:"【Z_x1】"`
	Interp33  int     `csv:"R_x1"`
	Interp34  int     `csv:"X_y1"`
	Interp35  int     `csv:"Y_y1"`
	Interp36  int     `csv:"【Z_y1】"`
	Interp37  int     `csv:"R_y1"`
	Interp38  int     `csv:"X_x2"`
	Interp39  int     `csv:"Y_x2"`
	Interp40  int     `csv:"【Z_x2】"`
	Interp41  int     `csv:"R_x2"`
	Interp42  int     `csv:"X_y2"`
	Interp43  int     `csv:"Y_y2"`
	Interp44  int     `csv:"【Z_y2】"`
	Interp45  int     `csv:"R_y2"`
	Interp46  int     `csv:"1"`
	Interp47  int     `csv:"0"`
	Interp48  int     `csv:"【R_x1】"`
	Interp49  int     `csv:"X_y1"`
	Interp50  int     `csv:"Y_y1"`
	Interp51  int     `csv:"Z_y1"`
	Interp52  int     `csv:"【R_y1】"`
	Interp53  int     `csv:"X_x2"`
	Interp54  int     `csv:"Y_x2"`
	Interp55  int     `csv:"Z_x2"`
	Interp56  int     `csv:"【R_x2】"`
	Interp57  int     `csv:"X_y2"`
	Interp58  int     `csv:"Y_y2"`
	Interp59  int     `csv:"Z_y2"`
	Interp60  int     `csv:"【R_y2】"`
	Interp61  int     `csv:"01"`
	Interp62  int     `csv:"00"`
	Interp63  int     `csv:"00"`
}

// motionMorphCsvRow はモーフCSVの1行を表す。
type motionMorphCsvRow struct {
	MorphName string  `csv:"モーフ名"`
	Frame     int     `csv:"フレーム"`
	Size      float64 `csv:"大きさ"`
}
