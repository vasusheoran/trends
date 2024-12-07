package contracts

import "gorm.io/gorm"

type CSVHeaders struct {
	gorm.Model
	A  string `json:"A" description:"Trading Symbol"`
	B  string `json:"B" description:"LTP"`
	C  string `json:"C" description:"Bid Qty"`
	D  string `json:"D" description:"Bid Rate"`
	E  string `json:"E" description:"Ask Rate"`
	F  string `json:"F" description:"Ask Qty"`
	G  string `json:"G" description:"LTQ"`
	H  string `json:"H" description:"Open"`
	I  string `json:"I" description:"High"`
	J  string `json:"J" description:"Low"`
	K  string `json:"K" description:"Prev Close"`
	L  string `json:"L" description:"Volume Traded Today"`
	M  string `json:"M" description:"Open Interest"`
	N  string `json:"N" description:"ATP"`
	O  string `json:"O" description:"Total Bid Qty"`
	P  string `json:"P" description:"Total Ask Qty"`
	Q  string `json:"Q" description:"-"`
	R  string `json:"R" description:"Exchange"`
	S  string `json:"S" description:"Upper Ckt Price"`
	T  string `json:"T" description:"Lower Ckt Price"`
	U  string `json:"U" description:"LTT"`
	V  string `json:"V" description:"LUT"`
	W  string `json:"W" description:"-"`
	X  string `json:"X" description:"-"`
	Y  string `json:"Y" description:"-"`
	Z  string `json:"Z" description:"-"`
	AA string `json:"AA" description:"-"`
	AB string `json:"AB" description:"-"`
	AC string `json:"AC" description:"-"`
	AD string `json:"AD" description:"-"`
	AE string `json:"AE" description:"-"`
	AF string `json:"AF" description:"-"`
	AG string `json:"AG" description:"-"`
	AH string `json:"AH" description:"-"`
	AI string `json:"AI" description:"-"`
	AJ string `json:"AJ" description:"-"`
	AK string `json:"AK" description:"-"`
	AL string `json:"AL" description:"-"`
	AM string `json:"AM" description:"-"`
	AN string `json:"AN" description:"-"`
	AO string `json:"AO" description:"-"`
	AP string `json:"AP" description:"-"`
	AQ string `json:"AQ" description:"-"`
	AR string `json:"AR" description:"-"`
	AS string `json:"AS" description:"-"`
	AT string `json:"AT" description:"-"`
	AU string `json:"AU" description:"-"`
	AV string `json:"AV" description:"-"`
	AW string `json:"AW" description:"-"`
	AX string `json:"AX" description:"-"`
	AY string `json:"AY" description:"-"`
	AZ string `json:"AZ" description:"-"`
	BA string `json:"BA" description:"-"`
	BB string `json:"BB" description:"-"`
	BC string `json:"BC" description:"-"`
	BD string `json:"BD" description:"-"`
	BE string `json:"BE" description:"-"`
	BF string `json:"BF" description:"-"`
	BG string `json:"BG" description:"-"`
	BH string `json:"BH" description:"-"`
	BI string `json:"BI" description:"-"`
	BJ string `json:"BJ" description:"-"`
	BK string `json:"BK" description:"-"`
	BL string `json:"BL" description:"-"`
	BM string `json:"BM" description:"-"`
	BN string `json:"BN" description:"-"`
	BO string `json:"BO" description:"-"`
	BP string `json:"BP" description:"-"`
	BQ string `json:"BQ" description:"-"`
	BR string `json:"BR" description:"-"`
	BS string `json:"BS" description:"-"`
	BT string `json:"BT" description:"-"`
	BU string `json:"BU" description:"-"`
	BV string `json:"BV" description:"-"`
	BW string `json:"BW" description:"-"`
	BX string `json:"BX" description:"-"`
	BY string `json:"BY" description:"-"`
	BZ string `json:"BZ" description:"-"`
	CA string `json:"CA" description:"-"`
	CB string `json:"CB" description:"-"`
	CC string `json:"CC" description:"-"`
	CD string `json:"CD" description:"-"`
	CE string `json:"CE" description:"-"`
	CF string `json:"CF" description:"-"`
	CG string `json:"CG" description:"-"`
	CH string `json:"CH" description:"-"`
	CI string `json:"CI" description:"-"`
	CJ string `json:"CJ" description:"-"`
	CK string `json:"CK" description:"-"`
	CL string `json:"CL" description:"-"`
	CM string `json:"CM" description:"-"`
	CN string `json:"CN" description:"-"`
	CO string `json:"CO" description:"-"`
	CP string `json:"CP" description:"-"`
	CQ string `json:"CQ" description:"-"`
	CR string `json:"CR" description:"-"`
	CS string `json:"CS" description:"-"`
	CT string `json:"CT" description:"-"`
	CU string `json:"CU" description:"-"`
	CV string `json:"CV" description:"-"`
	CW string `json:"CW" description:"-"`
	CX string `json:"CX" description:"-"`
	CY string `json:"CY" description:"-"`
	CZ string `json:"CZ" description:"-"`
	DA string `json:"DA" description:"-"`
	DB string `json:"DB" description:"-"`
	DC string `json:"DC" description:"-"`
	DD string `json:"DD" description:"-"`
	DE string `json:"DE" description:"-"`
	DF string `json:"DF" description:"-"`
	DG string `json:"DG" description:"-"`
	DH string `json:"DH" description:"-"`
	DI string `json:"DI" description:"-"`
	DJ string `json:"DJ" description:"-"`
	DK string `json:"DK" description:"-"`
	DL string `json:"DL" description:"-"`
	DM string `json:"DM" description:"Bid Qty"`
	DN string `json:"DN" description:"Bid Rate"`
	DO string `json:"DO" description:"Ask Rate"`
	DP string `json:"DP" description:"Open Interest"`
	DQ string `json:"DQ" description:"Total Bid Qty"`
	DR string `json:"DR" description:"-"`
	DS string `json:"DS" description:"-"`
	DT string `json:"DT" description:"-"`
	DU string `json:"DU" description:"-"`
	DV string `json:"DV" description:"-"`
	DW string `json:"DW" description:"-"`
	DX string `json:"DX" description:"-"`
	DY string `json:"DY" description:"-"`
	DZ string `json:"DZ" description:"-"`
	EA string `json:"EA" description:"-"`
	EB string `json:"EB" description:"-"`
	EC string `json:"EC" description:"-"`
	ED string `json:"ED" description:"-"`
	EE string `json:"EE" description:"-"`
	EF string `json:"EF" description:"-"`
}
