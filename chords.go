package tonacity

import (
	"fmt"
	"sort"
)

const (

	// These are intervals in the context of a key, which is made up of seven notes.

	// First For specifying the root tone. Here for completeness, not actually useful.
	First = 1
	// Second For specifying the interval of a second.
	Second = 2
	// Third For specifying the interval of a Third.
	Third = 3
	// Fourth For specifying the interval of a Fourth.
	Fourth = 4
	// Fifth For specifying the interval of a Fifth.
	Fifth = 5
	// Sixth For specifying the interval of a Sixth.
	Sixth = 6
	// Seventh For specifying the interval of a Seventh.
	Seventh = 7
	// Eighth is equivalent to first

	// MinorThird The interval of a minor third, in half steps.
	MinorThird = HalfStepValue * 3
	// MajorThird The interval of a major third, in half steps.
	MajorThird = HalfStepValue * 4
	// PerfectFourth The interval of a perfect fourth, in half steps.
	PerfectFourth = HalfStepValue * 5
	// PerfectFifth The interval of a perfect fifth, in half steps.
	PerfectFifth = HalfStepValue * 7
)

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
	return f.root.GetDistanceToHigherPitchClass(*f.GetPitchClass(Third)) == MajorThird
}

// HasPerfectFourth Returns true if this factory's fourth is a perfect fourth, false if it is a minor third.
func (f *ChordFactory) HasPerfectFourth() bool {
	return f.root.GetDistanceToHigherPitchClass(*f.GetPitchClass(Fourth)) == PerfectFourth
}

// HasPerfectFifth Returns true if this factory's fifth is a perfect fifth, false if it is a minor third.
func (f *ChordFactory) HasPerfectFifth() bool {
	return f.root.GetDistanceToHigherPitchClass(*f.GetPitchClass(Fifth)) == PerfectFifth
}

// Given a set of three or more pitch classes, how do we determine what the chord is?
// Does the key matter? Only in determining the Roman Numeral to use.
// Only the pattern of half step intervals matters.
// For each (type of) chord we know, we can iterate through all its inversions, adding them to a trie. The only remaining challenge
// is that to name the chord we need to know which position is the root.

func createTriadPattern(third1 HalfSteps, third2 HalfSteps) *Pattern {
	// Necessary to make a pattern that spans exactly an octave, as the third interval is required for chord inversions
	return MakePattern(third1, third2, OctaveValue-(third1+third2))
}

func createTetrachordPattern(third1 HalfSteps, third2 HalfSteps, third3 HalfSteps) *Pattern {
	// Necessary to make a pattern that spans exactly an octave, as the third interval is required for chord inversions
	return MakePattern(third1, third2, third3, OctaveValue-(third1+third2+third3))
}

// CreateMajorTriadPattern Creates the pattern for a Major Triad. Three of these exist in a key (I, IV, V, iii, vi, vii)
func CreateMajorTriadPattern() *Pattern {
	return createTriadPattern(MajorThird, MinorThird)
}

// CreateMinorTriadPattern Creates the pattern for a minor triad.  Three of these exist in a key (II, III, VI, i, iv, v)
func CreateMinorTriadPattern() *Pattern {
	return createTriadPattern(MinorThird, MajorThird)
}

// CreateDiminishedTriadPattern Creates the pattern for a diminished triad. Only one of these exists in a key (VII in major, ii in minor).
func CreateDiminishedTriadPattern() *Pattern {
	return createTriadPattern(MinorThird, MinorThird)
}

// CreateMajorSeventhPattern Creates the pattern for a Major Seventh Chord.
func CreateMajorSeventhPattern() *Pattern {
	return createTetrachordPattern(MajorThird, MinorThird, MajorThird)
}

// CreateDominantSeventhPattern Creates the pattern for a Dominant Seventh Chord.
func CreateDominantSeventhPattern() *Pattern {
	return createTetrachordPattern(MajorThird, MinorThird, MinorThird)
}

// CreateMinorSeventhPattern Creates the pattern for a Minor Seventh Chord.
func CreateMinorSeventhPattern() *Pattern {
	return createTetrachordPattern(MinorThird, MajorThird, MinorThird)
}

// CreateDiminishedSeventhPattern Creates the pattern for a Diminished Seventh Chord.
func CreateDiminishedSeventhPattern() *Pattern {
	return createTetrachordPattern(MinorThird, MinorThird, MinorThird)
}

type chordDictionaryEntry struct {
	name      string
	rootIndex int
}

func addChordToDict(dict *PatternDictionary, chord *Pattern, name string) {
	for i := 0; i < chord.Length(); i++ {
		dict.AddPattern(chord.Offset(-i), &chordDictionaryEntry{name, i})
	}
}

// CreateChordDictionary Creates a new dictionary for the purpose of naming chords based on the half step intervals between pitches.
func CreateChordDictionary() (dict *PatternDictionary) {
	dict = &PatternDictionary{NewTrie(12)}

	addChordToDict(dict, CreateMajorTriadPattern(), "Major")
	addChordToDict(dict, CreateMinorTriadPattern(), "Minor")
	addChordToDict(dict, CreateDiminishedTriadPattern(), "Diminished")

	addChordToDict(dict, CreateMajorSeventhPattern(), "Major Seventh")
	addChordToDict(dict, CreateDominantSeventhPattern(), "Dominant Seventh")
	addChordToDict(dict, CreateMinorSeventhPattern(), "Minor Seventh")
	addChordToDict(dict, CreateDiminishedSeventhPattern(), "Diminished Seventh")

	return
}

func GetChordName(dict *PatternDictionary, pitchNamer *PitchNamer, chord []PitchClass) (name string, ok bool) {

	// Assumption: only unique pitch classes are in chord

	if len(chord) <= 2 {
		ok = false
		return
	}

	// 1. Sort the pitches so their values are monotonically ascending
	pitchIndices := make(sort.IntSlice, len(chord), len(chord))

	for i, _ := range chord {
		pitchIndices[i] = int(chord[i].value)
	}

	pitchIndices.Sort()

	// 2. Create the pattern of the intervals between pitches, including between the last and first

	pattern := make([]HalfSteps, 0, len(pitchIndices))

	for i, v := range pitchIndices {
		interval := HalfSteps(pitchIndices[(i+1)%len(pitchIndices)] - v)
		if interval < 0 { // the last interval will be negative
			interval += OctaveValue
		}
		pattern = append(pattern, interval)
	}

	e := dict.GetEntry(&Pattern{pattern})
	entry, ok := e.(*chordDictionaryEntry)
	if !ok {
		return
	}

	root := pitchNamer.Name(PitchClass{HalfSteps(pitchIndices[entry.rootIndex])})
	name = fmt.Sprintf("%s %s", root, entry.name)
	ok = true

	return
}
