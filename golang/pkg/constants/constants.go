package constants

import "os"

const (
	KeyCP10      = "CP10"
	KeyCP50      = "CP50"
	KeyCP5       = "CP5"
	KeyCP20      = "CP20"
	KeyDiffCpNeg = "DiffCpNeg"
	KeyDiffCpPos = "DiffCpPos"

	BaseDir                = "data"
	BaseHistoricalDataPath = BaseDir + string(os.PathSeparator) + "hd"
	BaseRTDataPath         = BaseDir + string(os.PathSeparator) + "rt"
)
