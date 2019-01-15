package tonacity

// ChordFactory The purpose of this class is to allow creating chords using indexes into the scale, without needing to worry
// about half steps, e.g., specify "third" without having to know if it's a major (2 steps) or minor (three half steps) third.
type ChordFactory struct {
	pattern Pattern    // The pattern of the scale, which must be diatonic
	root    PitchClass // The pitch class to apply the pattern from
	offset  int        // The offset into the scale of the root
}

// GetPitchClass Get the pitch that is the given interval from this factory's root. The interval should be in the range 1 (a first,
// which will return the same pitch class) to 7 (a seventh). One is zero? Yes, that's just how music works ¯\_(ツ)_/¯.
func (f *ChordFactory) GetPitchClass(interval uint8) *PitchClass {
	var halfSteps HalfSteps
	for i := 1; i < int(interval)+f.offset; i++ {
		halfSteps += f.pattern.At(i - 1)
	}
	return f.root.GetTransposedCopy(halfSteps)
}

// HasMajorThird Returns true if this factory's third is a major third, false if it is a minor third.
func (f *ChordFactory) HasMajorThird() bool {
	return f.root.GetDistanceToHigherPitchClass(*f.GetPitchClass(3)) == 4
}

// HasPerfectFourth Returns true if this factory's fourth is a perfect fourth, false if it is a minor third.
func (f *ChordFactory) HasPerfectFourth() bool {
	return f.root.GetDistanceToHigherPitchClass(*f.GetPitchClass(4)) == 5
}

// HasPerfectFifth Returns true if this factory's fifth is a perfect fifth, false if it is a minor third.
func (f *ChordFactory) HasPerfectFifth() bool {
	return f.root.GetDistanceToHigherPitchClass(*f.GetPitchClass(5)) == 7
}
