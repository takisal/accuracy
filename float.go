package accurate

import (
	"log"
	"strconv"
)

/*
An Float represents a signed floating point number that can do accurate calculations, meaning there are no inaccuracies when storing a value, unlike with traditional floats in computer science. An empty Float should not be considered a valid value and can not be used in artihmetic.

Operations always take pointer arguments (*Float) rather than Float values, and each unique Float value requires its own unique *Float pointer. To "copy" a Float value, an existing (or newly allocated) Float must be set to a new value using the Float.Set method; shallow copies of Float are not supported and may lead to errors.

Note that methods may leak the Float's value through timing side-channels. Because of this and because of the scope and complexity of the implementation, Float is not well-suited to implement cryptographic operations.

Note that arithmetic with Floats is orders of magnitude slower than native floats. Floats should only ever be used if no other solution is suitable and exact value accuracy is worth the drastically increased computation time.
*/

type Float struct {
	Value           string `json:"Value"`
	SubOnePrecision uint   `json:"SubOnePrecision"`
	NonDecRep       string `json:"NonDecRep"`
}

// Set will change the value of x to v
func (x *Float) Set(v string) {
	x.Value = v
	var decSpots uint = 0
	var decBool bool = false
	var ndc string = ""
	for i := 0; i < len(v); i++ {
		if decBool == true {
			decSpots++
		}
		if string(v[i]) == "." {
			decBool = true
		} else {
			ndc += string(v[i])
		}
	}
	x.SubOnePrecision = decSpots
	x.NonDecRep = (ndc)

}
func Maxi(a *Float, b *Float) *Float {
	if b.Cmp(a) == 1 {
		return b
	}
	return a
}
func Mini(a *Float, b *Float) *Float {
	if b.Cmp(a) == -1 {
		return b
	}
	return a
}
func RoundTo(a *Float, d uint) *Float {
	v := a.Value
	origSuf := ""
	origPref := ""
	var sufAr []uint8
	var prefAr []uint8
	var tdp bool = false
	for i := len(v) - 1; i >= 0; i-- {
		if string(v[i]) == "." {
			tdp = true
		} else {
			tempuint, _ := strconv.Atoi(string(v[i]))

			if tdp == false {
				origSuf = string(v[i]) + origSuf
				sufAr = append([]uint8{uint8(tempuint)}, sufAr...)
			} else {
				origPref = string(v[i]) + origPref
				prefAr = append([]uint8{uint8(tempuint)}, prefAr...)
			}
		}
	}
	log.Println(prefAr, sufAr, origPref, origSuf)
	if d >= uint(len(sufAr)) {
		return NewFloat(origPref + "." + origSuf)
	}
	if d == 0 {
		if string(origSuf[0]) > "4" {

			var carryover uint8 = 1
			for i := len(prefAr) - 1; i >= 0 && carryover == 1; i-- {
				prefAr[i]++
				carryover = 0
				if prefAr[i] > 9 {
					prefAr[i] = 0
					carryover = 1
				}
			}
			if carryover == 1 {
				prefAr = append([]uint8{1}, prefAr...)
			}
			origPref = ""
			for i := 0; i < len(prefAr); i++ {
				origPref += strconv.FormatUint(uint64(prefAr[i]), 10)
			}
			return NewFloat(origPref + ".0")
		} else {
			return NewFloat(origPref + ".0")
		}
	} else {
		if string(origSuf[int(d)]) > "4" {
			var carryover uint8 = 1
			for i := int(d) - 1; i >= 0 && carryover == 1; i-- {
				sufAr[i]++
				carryover = 0
				if sufAr[i] > 9 {
					sufAr[i] = 0
					carryover = 1
				}
			}
			if carryover == 1 {
				origSuf = "0"
				carryover = 1
				for i := len(prefAr) - 1; i >= 0 && carryover == 1; i-- {
					prefAr[i]++
					carryover = 0
					if prefAr[i] > 9 {
						prefAr[i] = 0
						carryover = 1
					}
				}
				if carryover == 1 {
					prefAr = append([]uint8{1}, prefAr...)
				}
				origPref = ""
				for i := 0; i < len(prefAr); i++ {
					origPref += strconv.FormatUint(uint64(prefAr[i]), 10)
				}
			} else {
				tempStr := ""
				for i := 0; i < len(sufAr)-1; i++ {
					tempStr += strconv.FormatUint(uint64(sufAr[i]), 10)
				}
				origSuf = tempStr
			}
		}
	}
	for uint(len(origSuf)) > d {
		origSuf = origSuf[:len(origSuf)-1]
	}
	for len(origSuf) > 1 && string(origSuf[len(origSuf)-1]) == "0" {
		origSuf = origSuf[:len(origSuf)-1]
	}
	for len(origPref) > 1 && string(origPref[0]) == "0" {
		origPref = origPref[1:]
	}

	var retAF *Float = NewFloat(origPref + "." + origSuf)

	return retAF
}

// NewFloat allocates and returns a pointer to a new Float set to v
func NewFloat(v string) *Float {
	//creates an empty instance of an Float
	var x Float = Float{Value: "0", SubOnePrecision: 0, NonDecRep: "0"}

	x.Value = v
	var decSpots uint = 0
	var decBool bool = false
	var ndc string = ""
	for i := 0; i < len(v); i++ {
		if decBool == true {
			decSpots++
		}
		if string(v[i]) == "." {
			decBool = true
		} else {
			ndc += string(v[i])
		}
	}
	if decBool == false {
		x.Value += ".0"
		decSpots = 1
		ndc += "0"
	}
	x.SubOnePrecision = decSpots

	x.NonDecRep = (ndc)

	return &x
}

// TODO: check 0 decimal places and 1 decimal place

/* Div computes the quotient a/b for b != 0 and returns the quotient. Div will round down to d decimal places while trimming off trailing zeroes (after the decimal place), e.g. 100/10 = 10.0 regardless of if d is 1 or 25.
 */
func (a *Float) Div(b *Float, d uint) *Float {

	decimalPlaces := d
	if d == 0 {
		decimalPlaces = 2
	}

	var negCount uint8 = 0
	var atp string = a.NonDecRep
	var btp string = b.NonDecRep
	if string(a.Value[0]) == "-" {
		atp = a.NonDecRep[1:]
		negCount++
	}
	if string(b.Value[0]) == "-" {
		btp = b.NonDecRep[1:]
		negCount++
	}
	dif := int(b.SubOnePrecision) - int(a.SubOnePrecision)
	//refers to how many more digits on left of decimal point

	prefStr, sufStr := origDivString(atp, btp, max(int(decimalPlaces), int(decimalPlaces)+dif))
	for len(sufStr) < int(decimalPlaces)+1 {
		sufStr = sufStr + "0"
	}
	if len(prefStr) == 0 {
		prefStr = "0"
	}
	for dif > 0 {
		dif--
		prefStr += string(sufStr[0])
		sufStr = sufStr[1:]
		if len(sufStr) == 0 {
			sufStr = "0"
		}
	}

	for dif < 0 {
		dif++
		sufStr = string(prefStr[len(prefStr)-1]) + sufStr
		prefStr = prefStr[:len(prefStr)-1]
		if len(prefStr) == 0 {
			prefStr = "0"
		}

	}
	for len(sufStr) > int(decimalPlaces)+1 {
		sufStr = sufStr[:len(sufStr)-1]
	}
	if len(sufStr) < int(decimalPlaces)+1 {
		log.Println("TOO SHORT SUFFIX")
	}
	if d == 0 {
		if string(sufStr[0]) > "4" {
			prefAr := make([]uint8, len(prefStr))
			for i := 0; i < len(prefStr); i++ {
				tempuint, _ := strconv.Atoi(string(prefStr[i]))
				prefAr[i] = uint8(tempuint)
			}
			var carryover uint8 = 1
			for i := len(prefAr) - 1; i >= 0 && carryover == 1; i-- {
				prefAr[i]++
				carryover = 0
				if prefAr[i] > 9 {
					prefAr[i] = 0
					carryover = 1
				}
			}
			if carryover == 1 {
				prefAr = append([]uint8{1}, prefAr...)
			}
			prefStr = ""
			for i := 0; i < len(prefAr); i++ {
				prefStr += strconv.FormatUint(uint64(prefAr[i]), 10)
			}
			return NewFloat(prefStr + ".0")
		} else {
			return NewFloat(prefStr + ".0")
		}
	}
	sufAr := make([]uint8, len(sufStr))
	for i := 0; i < len(sufStr); i++ {
		tempuint, _ := strconv.Atoi(string(sufStr[i]))
		sufAr[i] = uint8(tempuint)
	}
	if sufStr == "0" {

	} else if len(sufStr) > 0 && len(sufStr) > (int(decimalPlaces)) && string(sufStr[int(decimalPlaces)]) > "4" {
		var carryover uint8 = 1
		for i := decimalPlaces - 1; i >= 0 && carryover == 1; i-- {

			sufAr[i]++
			carryover = 0
			if sufAr[i] > 9 {
				sufAr[i] = 0
				carryover = 1
			}
		}
		if carryover == 1 {

			sufStr = "0"
			prefAr := make([]uint8, len(prefStr))
			for i := 0; i < len(prefStr); i++ {
				tempuint, _ := strconv.Atoi(string(prefStr[i]))
				prefAr[i] = uint8(tempuint)
			}
			carryover = 1
			for i := len(prefAr) - 1; i >= 0 && carryover == 1; i-- {
				prefAr[i]++
				carryover = 0
				if prefAr[i] > 9 {
					prefAr[i] = 0
					carryover = 1
				}
			}
			if carryover == 1 {
				prefAr = append([]uint8{1}, prefAr...)

			}
			prefStr = ""
			for i := 0; i < len(prefAr); i++ {
				prefStr += strconv.FormatUint(uint64(prefAr[i]), 10)
			}
		} else {
			tempStr := ""
			for i := 0; i < len(sufAr)-1; i++ {
				tempStr += strconv.FormatUint(uint64(sufAr[i]), 10)
			}
			sufStr = tempStr
		}
	} else if len(sufStr) > 0 && len(sufStr) > int(int(decimalPlaces)) {

		sufStr = sufStr[:len(sufStr)-1] + "0"

	}

	for len(sufStr) > 1 && string(sufStr[len(sufStr)-1]) == "0" {
		sufStr = sufStr[:(len(sufStr) - 1)]
	}
	for len(prefStr) > 1 && string(prefStr[0]) == "0" {
		prefStr = prefStr[1:]
	}
	if negCount == 1 {
		prefStr = "-" + prefStr
	}
	for uint(len(sufStr)) > d && d > 0 {
		sufStr = sufStr[:len(sufStr)-1]
	}
	var retAF *Float = NewFloat(prefStr + "." + sufStr)

	return retAF

}

// Add computes the sum a+b and returns the sum
func (a *Float) Add(b *Float) *Float {
	var negCount uint8 = 0
	var atp string = a.NonDecRep
	var btp string = b.NonDecRep
	bneg := false
	aneg := false
	if string(a.Value[0]) == "-" {
		atp = a.NonDecRep[1:]
		negCount++
		aneg = true
	}
	if string(b.Value[0]) == "-" {
		btp = b.NonDecRep[1:]
		negCount++
		bneg = true
	}
	astr := atp
	bstr := btp
	if negCount == 0 || negCount == 2 {
		apd := a.SubOnePrecision
		bpd := b.SubOnePrecision
		for apd < bpd {
			astr += "0"
			apd++
		}
		for apd > bpd {
			bstr += "0"
			bpd++
		}
		var decSpots uint = apd
		productStr := addstr(astr, bstr)
		sufStr := ""
		prefStr := ""
		visSpots := 0
		for i := len(productStr) - 1; i >= 0; i-- {
			if visSpots >= int(decSpots) {
				prefStr = string(productStr[i]) + prefStr
			} else {
				sufStr = string(productStr[i]) + sufStr
			}
			visSpots++
		}
		if len(sufStr) == 0 {
			sufStr += "0"
		}
		if negCount == 2 {
			nzc := 0
			for i := 0; i < len(prefStr); i++ {
				if string(prefStr[i]) != "0" {
					nzc++
				}
			}
			for i := 0; i < len(sufStr); i++ {
				if string(sufStr[i]) != "0" {
					nzc++
				}
			}
			if nzc != 0 {
				prefStr = "-" + prefStr
			}
		}
		var retAF *Float = NewFloat(prefStr + "." + sufStr)

		return retAF
	} else {
		//(a == false && b == true) || (a == true && b == false)

		if aneg == false && bneg == true {
			tempAF := NewFloat(b.Value[1:])
			return a.Sub(tempAF)
		} else {
			//(aneg == true && bneg == false)
			tempAF := NewFloat(a.Value[1:])
			newRes := tempAF.Sub(b)
			if string(newRes.Value[0]) == "-" {
				return NewFloat(newRes.Value[1:])
			} else {
				nzc := 0
				for i := 0; i < len(newRes.Value); i++ {
					if string(newRes.Value[i]) != "0" && string(newRes.Value[i]) != "." {
						nzc++
					}
				}
				if nzc > 0 {
					return NewFloat("-" + newRes.Value)
				} else {
					return newRes
				}

			}
		}
	}
}

// Sub computes the difference of a-b and returns the difference
func (a *Float) Sub(b *Float) *Float {
	var negCount uint8 = 0
	var atp string = a.NonDecRep
	var btp string = b.NonDecRep
	if a.Value == b.Value {
		return NewFloat("0.0")
	}
	bneg := false
	aneg := false
	if string(a.Value[0]) == "-" {
		atp = a.NonDecRep[1:]
		negCount++
		aneg = true
	}
	if string(b.Value[0]) == "-" {
		btp = b.NonDecRep[1:]
		negCount++
		bneg = true
	}
	if negCount == 0 || negCount == 2 {
		astr := atp
		bstr := btp
		apd := a.SubOnePrecision
		bpd := b.SubOnePrecision
		for apd < bpd {
			astr += "0"
			apd++
		}
		for apd > bpd {
			bstr += "0"
			bpd++
		}
		var decSpots uint = apd
		productStr := substr(astr, bstr)
		sufStr := ""
		prefStr := ""
		visSpots := 0
		for i := len(productStr) - 1; i >= 0; i-- {
			if visSpots >= int(decSpots) {
				prefStr = string(productStr[i]) + prefStr
			} else {
				sufStr = string(productStr[i]) + sufStr
			}
			visSpots++
		}
		for len(sufStr) < int(decSpots) {
			sufStr = "0" + sufStr
		}
		if len(sufStr) == 0 {
			sufStr = "0"
		}

		if negCount == 2 {
			if string(prefStr[0]) == "-" {
				prefStr = prefStr[1:]
			} else {
				nzc := 0
				for i := 0; i < len(prefStr); i++ {
					if string(prefStr[i]) != "0" {
						nzc++
					}
				}
				for i := 0; i < len(sufStr); i++ {
					if string(sufStr[i]) != "0" {
						nzc++
					}
				}
				if nzc != 0 {
					prefStr = "-" + prefStr
				}
			}
		}
		if len(prefStr) == 0 {
			prefStr += "0"
		}
		var retAF *Float = NewFloat(prefStr + "." + sufStr)

		return retAF
	} else {
		//if aneg == true && bneg == false || aneg == false && bneg == true
		astr := atp
		bstr := btp
		apd := a.SubOnePrecision
		bpd := b.SubOnePrecision
		for apd < bpd {
			astr += "0"
			apd++
		}
		for apd > bpd {
			bstr += "0"
			bpd++
		}
		var decSpots uint = apd
		productStr := addstr(astr, bstr)
		sufStr := ""
		prefStr := ""
		visSpots := 0
		for i := len(productStr) - 1; i >= 0; i-- {
			if visSpots >= int(decSpots) {
				prefStr = string(productStr[i]) + prefStr
			} else {
				sufStr = string(productStr[i]) + sufStr
			}
			visSpots++
		}
		if len(sufStr) == 0 {
			sufStr += "0"
		}
		if aneg == true && bneg == false {
			prefStr = "-" + prefStr
		}

		var retAF *Float = NewFloat(prefStr + "." + sufStr)

		return retAF
	}
}

// Mul computes the product a*b and returns the product
func (a *Float) Mul(b *Float) *Float {
	var negCount uint8 = 0
	var atp string = a.NonDecRep
	var btp string = b.NonDecRep
	if string(a.Value[0]) == "-" {
		atp = a.NonDecRep[1:]
		negCount++
	}
	if string(b.Value[0]) == "-" {
		btp = b.NonDecRep[1:]
		negCount++
	}
	var decSpots uint = a.SubOnePrecision + b.SubOnePrecision
	productStr := mulstr(atp, btp)
	for len(productStr) <= int(decSpots) {
		productStr = "0" + productStr
	}
	sufStr := ""
	prefStr := ""
	visSpots := 0
	for i := len(productStr) - 1; i >= 0; i-- {
		if visSpots >= int(decSpots) {
			prefStr = string(productStr[i]) + prefStr
		} else {
			sufStr = string(productStr[i]) + sufStr
		}
		visSpots++
	}
	if len(sufStr) == 0 {
		sufStr = "0"
	}
	if len(prefStr) == 0 {
		prefStr = "0"
	}
	for len(prefStr) > 1 && string(prefStr[0]) == "0" {
		prefStr = prefStr[:1]
	}
	if negCount == 1 {
		prefStr = "-" + prefStr
	}
	var retAF *Float = NewFloat(prefStr + "." + sufStr)

	return retAF

}

// Cmp compares a with b, returning -1 if a < b, 1 if a > b, and 0 if a == b
func (a *Float) Cmp(b *Float) int8 {

	amags := len(a.Value) - int(a.SubOnePrecision+1)
	bmags := len(b.Value) - int(b.SubOnePrecision+1)
	if amags > bmags {
		return 1
	}
	if amags < bmags {
		return -1
	}
	for i := 0; i < amags; i++ {
		if a.Value[i] < b.Value[i] {
			return -1
		} else if a.Value[i] > b.Value[i] {
			return 1
		}
	}
	lowestSuf := min(amags+int(a.SubOnePrecision+1), bmags+int(b.SubOnePrecision+1))
	for i := amags + 1; i < lowestSuf; i++ {
		if a.Value[i] < b.Value[i] {
			return -1
		} else if a.Value[i] > b.Value[i] {
			return 1
		}
	}

	if len(b.Value) > lowestSuf {
		for i := lowestSuf; i < len(b.Value); i++ {
			if b.Value[i] > byte('0') {
				return -1
			}
		}
	} else if len(a.Value) > lowestSuf {
		for i := lowestSuf; i < len(a.Value); i++ {
			if a.Value[i] > byte('0') {
				return 1
			}
		}
	}

	return 0
}

func findhighestbeloworequal(a string, b string) (uint8, string) {
	var cur string = "0"
	var retval uint8 = 0
	var retstr string = ""
	for i := 1; i <= 10; i++ {
		cur = addstr(cur, b)
		if strcmp(cur, a) < 1 {
			retstr = cur
			retval = uint8(i)
		}
	}
	return retval, retstr
}
func trimslice(v []uint16) []uint16 {
	for len(v) > 0 && v[0] == 0 {
		v = v[1:]
	}
	if len(v) == 0 {
		return []uint16{0}
	} else {
		return v
	}
}
func trimstring(v string) string {
	for len(v) > 0 && string(v[0]) == "0" {
		v = v[1:]
	}
	if len(v) == 0 {
		return "0"
	} else {
		return v
	}
}
func addslice(a []uint16, b []uint16) []uint16 {
	valAr := make([]uint16, max(len(a), len(b)))
	var nza bool = false
	var nzb bool = false
	for i := 0; i < len(a); i++ {
		if a[i] != 0 {
			nza = true
			break
		}
	}

	for i := 0; i < len(b); i++ {
		if b[i] != 0 {
			nzb = true
			break
		}
	}

	if nza == false && nzb == false {
		return []uint16{0}
	} else if nza == false {
		return trimslice(b)
	} else if nzb == false {
		return trimslice(a)
	}

	for len(a) != len(b) {
		if len(a) < len(b) {
			a = append([]uint16{0}, a...)
		} else {
			b = append([]uint16{0}, b...)
		}
	}
	for i := 0; i < len(a); i++ {
		valAr[i] = uint16(a[i] + b[i])
	}
	var curSum uint16 = 0
	for i := len(valAr) - 1; i >= 0; i-- {
		valAr[i] += (curSum)
		curSum = 0
		if valAr[i] > 9 {
			curSum += (valAr[i] / 10)
			valAr[i] = valAr[i] % 10
		}
	}
	for curSum > 0 {
		valAr = append([]uint16{curSum}, valAr...)
		curSum = 0
		if valAr[0] > 9 {
			curSum += (valAr[0] / 10)
			valAr[0] = valAr[01] % 10
		}
	}
	var valStr []uint16
	var fnz bool = false
	for i := 0; i < len(valAr); i++ {
		if valAr[i] != 0 {
			fnz = true
		}
		if fnz {
			valStr = append(valStr, valAr[i])
		}
	}

	return trimslice(valStr)

}
func addstr(a string, b string) string {
	valAr := make([]uint8, max(len(a), len(b)))
	var nza bool = false
	var nzb bool = false
	for i := 0; i < len(a); i++ {
		if string(a[i]) != "0" {
			nza = true
			break
		}
	}

	for i := 0; i < len(b); i++ {
		if string(b[i]) != "0" {
			nzb = true
			break
		}
	}

	if nza == false && nzb == false {
		return "0"
	} else if nza == false {
		return trimstring(b)
	} else if nzb == false {
		return trimstring(a)
	}
	for len(a) != len(b) {
		if len(a) < len(b) {
			a = "0" + a
		} else {
			b = "0" + b
		}
	}

	for i := 0; i < len(a); i++ {
		a1, _ := strconv.Atoi(string(a[i]))
		b1, _ := strconv.Atoi(string(b[i]))
		ard := a1 + b1
		valAr[i] = uint8(ard)
	}
	var curSum uint8 = 0
	for i := len(valAr) - 1; i >= 0; i-- {
		valAr[i] += (curSum)
		curSum = 0
		if valAr[i] > 9 {
			curSum += (valAr[i] / 10)
			valAr[i] = valAr[i] % 10
		}
	}
	for curSum > 0 {
		valAr = append([]uint8{curSum}, valAr...)
		curSum = 0
		if valAr[0] > 9 {
			curSum += (valAr[0] / 10)
			valAr[0] = valAr[0] % 10
		}
	}
	valStr := ""
	var fnz bool = false
	for i := 0; i < len(valAr); i++ {
		if valAr[i] != 0 {
			fnz = true
		}
		if fnz {
			valStr += strconv.FormatUint(uint64(valAr[i]), 10)
		}
	}
	return trimstring(valStr)

}
func substr(f string, s string) string {
	var a string
	var b string
	var negatory bool = false
	var eqq = true
	if len(f) == len(s) {

		for i := 0; i < len(f); i++ {
			if f[i] > s[i] {
				a = f
				b = s
				eqq = false
				break

			} else if f[i] < s[i] {
				b = f
				negatory = true
				a = s
				eqq = false
				break
			}
		}
		if eqq == true {
			return "0"
		}
	} else if len(f) < len(s) {
		b = f
		negatory = true
		a = s
	} else {
		a = f
		b = s
	}
	valAr := make([]int8, max(len(a), len(b)))
	for len(a) != len(b) {
		if len(a) < len(b) {
			a = "0" + a
		} else {
			b = "0" + b
		}
	}
	for i := 0; i < len(a); i++ {
		a1, _ := strconv.Atoi(string(a[i]))
		b1, _ := strconv.Atoi(string(b[i]))
		ard := int8(a1 - b1)
		valAr[i] = (ard)
	}
	var carrysub int = 0
	for i := len(valAr) - 1; i >= 0; i-- {
		valAr[i] -= int8(carrysub)
		carrysub = 0
		if valAr[i] < 0 {
			carrysub = 1
			valAr[i] = 10 + valAr[i]
		}
	}

	valStr := ""
	var fnz bool = false
	for i := 0; i < len(valAr); i++ {
		if valAr[i] != 0 {
			fnz = true
		}
		if fnz {
			valStr += strconv.FormatUint(uint64(valAr[i]), 10)
		}
	}
	for len(valStr) > 1 && string(valStr[0]) == "0" {
		valStr = valStr[1:]
	}
	if negatory && valStr != "0" {
		valStr = "-" + valStr
	}

	return valStr
}
func strcmp(a string, b string) int8 {
	if len(a) == 0 {
		if len(b) == 0 {
			return 0
		} else {
			return -1
		}
	}
	if len(b) == 0 {
		if len(a) == 0 {
			return 0
		} else {
			return 1
		}
	}
	for len(a) > 1 && string(a[0]) == ("0") {
		a = a[1:]
	}
	for len(b) > 1 && string(b[0]) == ("0") {
		b = b[1:]
	}
	var prefa string
	var prefb string
	var sufa string
	var sufb string
	var dechit bool = false
	for i := 0; i < len(a); i++ {
		if string(a[i]) == "." {
			dechit = true
		} else {
			if dechit == false {
				prefa += string(a[i])
			} else {
				sufa += string(a[i])
			}
		}
	}
	dechit = false
	for i := 0; i < len(b); i++ {
		if string(b[i]) == "." {
			dechit = true
		} else {
			if dechit == false {
				prefb += string(b[i])
			} else {
				sufb += string(b[i])
			}
		}
	}
	if len(prefa) == len(prefb) {

		for i := 0; i < len(prefa); i++ {
			if prefa[i] > prefb[i] {
				return 1
			} else if prefa[i] < prefb[i] {
				return -1
			}
		}

		minsuflen := min(len(sufa), len(sufb))
		for i := 0; i < minsuflen; i++ {
			if sufa[i] > sufb[i] {
				return 1
			} else if sufa[i] < sufb[i] {
				return -1
			}
		}
		if len(sufa) == len(sufb) {
			return 0
		}
		if len(sufa) > len(sufb) {
			for i := minsuflen; i < len(sufa); i++ {
				if string(sufa[i]) > "0" {
					return 1

				}
			}
		} else {
			for i := minsuflen; i < len(sufb); i++ {
				if string(sufb[i]) > "0" {
					return -1
				}
			}
		}
		return 0

	} else if len(prefa) > len(prefb) {
		return 1
	} else {
		return -1
	}
}

func origDivString(v string, b string, decimalPlaces int) (string, string) {
	valueStr := ""
	valAr := make([]uint8, len(v))
	blen := len(b)
	for len(v) >= blen {
		//how many times does it go into the first 3 digit
		newv := v[:blen]
		v = v[blen:]
		incfac := len(v)
		if strcmp(newv, b) == -1 {
			if len(v) == 0 {
				v = newv
				break
			}
			newv += string(v[0])
			v = v[1:]
			incfac--

		}
		numbelow, actualAmount := findhighestbeloworequal(newv, b)
		pta := substr(newv, actualAmount)
		v = pta + v
		valAr[incfac] += numbelow
		//subract from first 3 digits and replace
	}

	var curSum uint8 = 0
	for i := 0; i < len(valAr); i++ {
		valAr[i] += (curSum)
		curSum = 0
		if valAr[i] > 9 {
			curSum += (valAr[i] / 10)
			valAr[i] = valAr[i] % 10
		}
	}
	for curSum > 0 {
		valAr = append(valAr, curSum)
		curSum = 0
		if valAr[len(valAr)-1] > 9 {
			curSum += (valAr[len(valAr)-1] / 10)
			valAr[len(valAr)-1] = valAr[len(valAr)-1] % 10
		}
	}
	var fnz bool = false
	for i := len(valAr) - 1; i >= 0; i-- {
		if valAr[i] != 0 {
			fnz = true
		}
		if fnz == true {
			valueStr = valueStr + strconv.FormatUint(uint64(valAr[i]), 10)
		}
	}
	if len(valueStr) == 0 {
		valueStr = "0"
	}
	remmy, _ := strconv.ParseUint(v, 10, 8)
	if remmy == 0 {
		return valueStr, "0"
	} else {
		fa := divString(v+"0", b, decimalPlaces+1)
		return valueStr, fa //must round
	}
}
func divString(v string, b string, c int) string {
	if c <= 0 {
		return ""
	}
	if v == "0" || v == "00" || v == "000" {
		return "0"
	}
	numbelow, acAm := findhighestbeloworequal(v, b)
	remmy := substr(v, acAm)
	return strconv.FormatUint(uint64(numbelow), 10) + divString(remmy+"0", b, c-1)

}
func mulstr(astr string, bstr string) string {
	//more efficient if a is longer
	var tstr string
	if len(astr) < len(bstr) {
		tstr = astr
		astr = bstr
		bstr = tstr
	}
	var holder []uint16 = []uint16{0}
	var a []uint16
	var b []uint16
	for i := 0; i < len(astr); i++ {
		intpar, _ := strconv.ParseUint(string(astr[i]), 10, 8)
		a = append(a, uint16(intpar))
	}
	for i := 0; i < len(bstr); i++ {
		intpar, _ := strconv.ParseUint(string(bstr[i]), 10, 8)
		b = append(b, uint16(intpar))
	}
	z := 0
	for j := len(b) - 1; j >= 0; j-- {
		valAr := make([]uint16, len(a))
		for i := 0; i < len(a); i++ {
			valAr[i] = uint16(a[i] * b[j])
		}
		var curSum uint16 = 0
		for i := len(valAr) - 1; i >= 0; i-- {
			valAr[i] += (curSum)
			curSum = 0
			if valAr[i] > 9 {
				curSum += (valAr[i] / 10)
				valAr[i] = valAr[i] % 10
			}
		}
		for curSum > 0 {
			valAr = append([]uint16{curSum}, valAr...)
			curSum = 0
			if valAr[0] > 9 {
				curSum += (valAr[0] / 10)
				valAr[0] = valAr[0] % 10
			}
		}
		for k := 0; k < z; k++ {
			valAr = append(valAr, 0)
		}
		holder = addslice(holder, valAr)
		z++
	}
	retStr := ""
	for i := 0; i < len(holder); i++ {
		retStr += strconv.FormatUint(uint64(holder[i]), 10)
	}
	return retStr
}
