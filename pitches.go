package tonacity

// The entirety of music is based on picking a somewhat arbitrary frequency as a starting point, then going from there, so it seems like doing the same with this library might be a good idea

// C The natural tone of C.
func C() *PitchClass {
	return &PitchClass{0}
}

// D The natural tone of D.
func D() *PitchClass {
	return &PitchClass{2}
}

// E The natural tone of E.
func E() *PitchClass {
	return &PitchClass{4}
}

// F The natural tone of F.
func F() *PitchClass {
	return &PitchClass{5}
}

// G The natural tone of G.
func G() *PitchClass {
	return &PitchClass{7}
}

// A The natural tone of A.
func A() *PitchClass {
	return &PitchClass{9}
}

// B The natural tone of B.
func B() *PitchClass {
	return &PitchClass{11}
}

// Why aren't I using Middle C as the starting point? Because music science takes A4 as 0, and all physical note frequencies are calculated
// using their offset from A4

// A4 Returns the first A above middle C (but actually the fifth A on a piano because of A0)
func A4() *Pitch {
	return &Pitch{*A(), 0}
}

// MiddleC Returns the middle C of a piano, C4.
func MiddleC() *Pitch {
	a4 := A4()
	a4.LowerToNext(*C())
	return a4
}
