package tonacity

import (
	"fmt"
	"math"
)

// StandardConcertPitch The almost-completely agreed-upon pitch of A4
const StandardConcertPitch = 440

// PitchClass All musical tones separated by an octave. A major scale has 7, the chromatic scale has 11, etc.
type PitchClass struct {
	value HalfSteps
}

// Transpose Transpose this pitch class to the pitch class the given number of half steps away. As this is a class and not a specific pitch
// this operation will "loop" around. For example: transposing C by 2, 14, 26, etc. will produce D.
func (pc *PitchClass) Transpose(halfSteps HalfSteps) {
	pc.value = (pc.value + halfSteps) % OctaveValue
}

// GetTransposedCopy Returns a copy of this pitch class transposed by the given number of half steps.
func (pc *PitchClass) GetTransposedCopy(halfSteps HalfSteps) *PitchClass {
	copy := *pc
	copy.Transpose(halfSteps)
	return &copy
}

// Sharp Returns a copy of this pitch class transposed by +1 half step.
func (pc *PitchClass) Sharp() *PitchClass {
	return pc.GetTransposedCopy(HalfStepValue)
}

// Flat Returns a copy of this pitch class transposed by -1 half step.
func (pc *PitchClass) Flat() *PitchClass {
	return pc.GetTransposedCopy(-HalfStepValue)
}

// GetDistanceToHigherPitchClass Gets the positive number of half steps to the given pitch class. Passing in the same pitch class
// will return an entire octave.
func (pc *PitchClass) GetDistanceToHigherPitchClass(b PitchClass) HalfSteps {
	if pc.value == b.value {
		return OctaveValue
	}
	if pc.value > b.value {
		return OctaveValue - (pc.value - b.value)
	}
	return b.value - pc.value
}

// GetDistanceToLowerPitchClass Gets the negative number of half steps to the given pitch class. Passing in the same pitch class
// will return an entire (negative) octave.
func (pc *PitchClass) GetDistanceToLowerPitchClass(b PitchClass) HalfSteps {
	if pc.value == b.value {
		return -OctaveValue
	}
	if pc.value > b.value {
		return b.value - pc.value
	}
	return -(OctaveValue - (b.value - pc.value))
}

// HasSamePitchAs Returns true if both receiver and argument are non-null and have the same pitch class.
func (pc *PitchClass) HasSamePitchAs(other *PitchClass) bool {
	if pc == nil || other == nil {
		return false
	}
	return (pc.value % OctaveValue) == (other.value % OctaveValue)
}

// PitchNamer Uses a lookup to return names for pitch classes. This exists because the name for a particular pitch depends entirely on context.
// In the key of B♭ (or A♯) Major, flats are used because the key includes the note A. Sharp or flat is chosen to avoid such overlaps.
// G♭ (F♯) is an edge case as it includes both B♭, B, F, and F♯, so when using flats describes B as C♭, and when using sharps describes F as E♯.
// There are also double-sharps and double-flats. A diminished 7th chord uses a double flat (𝄫), because that note is the 7th lowered a whole step,
// and the chord is a diminished *7th*. So in a C Diminished 7th chord, the A is referred to as B𝄫.
type PitchNamer struct {
	lookup map[HalfSteps]string
}

// Name Get the name of the given pitch class.
func (pn *PitchNamer) Name(pitch PitchClass) string {
	return pn.lookup[pitch.value]
}

// CreateSharpPitchNamer Creates a pitch namer that will use sharps (♯) to describe the "the black keys".
func CreateSharpPitchNamer() *PitchNamer {
	return &PitchNamer{map[HalfSteps]string{
		0:  "C",
		1:  "C♯",
		2:  "D",
		3:  "D♯",
		4:  "E",
		5:  "F",
		6:  "F♯",
		7:  "G",
		8:  "G♯",
		9:  "A",
		10: "A♯",
		11: "B",
	}}
}

// CreateFlatPitchNamer Creates a pitch namer that will use flats (♭) to describe the "the black keys".
func CreateFlatPitchNamer() *PitchNamer {
	return &PitchNamer{map[HalfSteps]string{
		0:  "C",
		1:  "D♭",
		2:  "D",
		3:  "E♭",
		4:  "E",
		5:  "F",
		6:  "G♭",
		7:  "G",
		8:  "A♭",
		9:  "A",
		10: "B♭",
		11: "B",
	}}
}

// Pitch A specific note, defined as a pitch class (e.g. C) and an octave (e.g. 4). For example: A piano goes from A0 to C8
type Pitch struct {
	class PitchClass // the pitch class, e.g., C
	value HalfSteps  // ordinal value of the pitch, for the purpose of calculating inter-pitch distances
	// Why isn't physical frequency here? Because while the standard is for A4 to be 440Hz, that's not what everyone agrees upon
}

func (p *Pitch) String() string {
	return fmt.Sprintf("class %d, value %d", p.class.value, p.value)
}

// Class The class of the pitch, e.g., C
func (p *Pitch) Class() PitchClass {
	return p.class
}

// GetDistanceTo Gets the number of half steps to the given pitch. If b is lower then a then the number will be negative.
func (p *Pitch) GetDistanceTo(b *Pitch) HalfSteps {
	return b.value - p.value
}

// Octave The octave of the pitch relative to a (Piano) Middle C (C4). Be careful when using this for presentation, as C♭4 will return 3, as
// it will be treated as a B. In such situations, first get the octave of the natural tone, then ornament it afterwards.
func (p *Pitch) Octave(middleC *Pitch) (octave int8) {
	offsetFromCZero := middleC.GetDistanceTo(p) + OctaveValue*4
	octave = int8(offsetFromCZero / OctaveValue)
	return
}

// Transpose Transposes this pitch by the given number of half steps
func (p *Pitch) Transpose(halfSteps HalfSteps) {
	p.class.Transpose(halfSteps)
	p.value += halfSteps
}

// GetTransposedCopy Returns a copy of this pitch transposed by the given number of half steps
func (p *Pitch) GetTransposedCopy(halfSteps HalfSteps) *Pitch {
	copy := *p
	copy.Transpose(halfSteps)
	return &copy
}

// RaiseToNext Raises this pitch to the next one of the given class. Guaranteed to modify this pitch.
func (p *Pitch) RaiseToNext(pitchClass PitchClass) {
	diff := p.class.GetDistanceToHigherPitchClass(pitchClass)
	p.Transpose(diff)
}

// LowerToNext Lowers this pitch to the next one of the given class. Guaranteed to modify this pitch.
func (p *Pitch) LowerToNext(pitchClass PitchClass) {
	diff := p.class.GetDistanceToLowerPitchClass(pitchClass)
	p.Transpose(diff)
}

// FrequencyInHertz Returns the physical frequency of the sound produced by the given pitch, using the given frequency as A4.
func (p *Pitch) FrequencyInHertz(concertPitch float64) float64 {
	return math.Pow(2, float64(A4().GetDistanceTo(p))/12.0) * concertPitch
}

// PitchFactory Allows creating pitches relative to C4.
type PitchFactory struct {
	c4 Pitch
}

// GetPitch Get the pitch with the given class in the given octave.
func (f *PitchFactory) GetPitch(pitchClass *PitchClass, octave int) *Pitch {

	// We're working with C4 being zero, so offset octave
	octave -= 4

	// octave == 0 => After a pitch in the same octave (4)
	// octave > 0  => After a pitch in a higher octave (5+)
	// octave < 0  => After a pitch in a lower octave  (3-)

	interval := (f.c4.class.GetDistanceToHigherPitchClass(*pitchClass) % OctaveValue)
	return &Pitch{*pitchClass, interval + HalfSteps(octave*OctaveValue)}
}
