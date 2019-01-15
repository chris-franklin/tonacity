package tonacity

// Welcome to music theory, where nothing is unanimously agreed upon, there are multiple equivalent ways of saying the same thing, and the distance between two notes is a second.
// Evidence:
//  - What physical frequency is a particular note?
//  - Which C is middle C?
//  - A𝄫 == G == F𝄪
//  - B♭ == A♯
//  - C♭ == B
//  - In C Major the distance between C (note I in the scale), and F (note IV in the scale) is "a fourth"

// In this library presentation has been kept completely separate to the data structures. This is because the name of something typically depends entirely on context.

const (
	// WholeStepValue The numerical value of a whole step
	WholeStepValue = 2
	// HalfStepValue The numerical value of a half step
	HalfStepValue = WholeStepValue / 2
	// OctaveValue The numerical value of an entire octave
	OctaveValue = HalfStepValue * 12
)

// HalfSteps The type for specifying a number of half steps
type HalfSteps int8

// Singer Make that object SING. Something which can, given a valid starting pitch, produce an indefinite sequence of notes. For example, the G Major key signature will produce,
// given a starting pitch of B3, the pitch classes B3, C4, D4, E4, F♯4, G4, A4, B4, C5, D5, etc.
// Not all starting pitches may be valid, for example if the receiver is the key of G Major then a starting pitch of E♭, which is not in the key, is not valid.
type Singer interface {
	// Sing Generates the next note. Bool will be false if there are no more notes.
	Sing() (pitch Pitch, more bool)
}

// PitchClassProducer Something which contains a set of pitch classes, where order and starting point are irrelevant. For example, the G Major key signature will produce
// the pitch classes G, A, B, C, D, E, F♯.
type PitchClassProducer interface {
	// Generates a slice of unique pitch classes, returned in ascending order (to be useful to the caller, not because the order necessarily has a meaning).
	ProducePitchClasses() []*PitchClass
}

// Transposer Something which can be transposed to alter its tone by a number of half steps. For example: transposing the scale of G Major
// by (plus) three half steps will yield the scale of B♭ Major.
type Transposer interface {
	Transpose(halfSteps HalfSteps)
}

// PatternRepeatingSinger A singer that can produce pitches indefinitely. Unless the pattern of steps purposefully loops back, i.e. the sum of all is zero, then
// this singer does NOT loop, i.e., given infinite time it will either tend to a pitch of negative or positive infinity hertz.
type PatternRepeatingSinger struct {
	pattern   Pattern
	nextPitch Pitch
	offset    int
}

// Sing Keeps producing the next pitch in the sequence according to its underlying pattern of half-step intervals
func (singer PatternRepeatingSinger) Sing() (pitch Pitch, more bool) {
	pitch = singer.nextPitch
	singer.nextPitch.Transpose(singer.pattern.At(singer.offset))
	singer.offset = (singer.offset + 1) % singer.pattern.Length()
	// Can always produce more notes
	more = true
	return
}

type KeySignature struct {
	pitches []*PitchClass
}

type TimeSignature struct {
	noteCount int
	noteValue int
}

type Bar struct {
	time TimeSignature
}

type Note struct {
	// The value of this note is one over this value
	value int
}

type Stave struct {
	bars []Bar
}

type Clef struct {
	firstNote int
}
