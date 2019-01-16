package tonacity

import (
	"fmt"
	"sort"
)

// Interval The type to use when specifying an interval in a key.
type Interval uint8

const (

	// These are intervals in the context of a key, which is made up of seven notes.

	// First For specifying the root tone. Here for completeness, not actually useful.
	First Interval = 1
	// Second For specifying the interval of a Second.
	Second Interval = 2
	// Third For specifying the interval of a Third.
	Third Interval = 3
	// Fourth For specifying the interval of a Fourth.
	Fourth Interval = 4
	// Fifth For specifying the interval of a Fifth.
	Fifth Interval = 5
	// Sixth For specifying the interval of a Sixth.
	Sixth Interval = 6
	// Seventh For specifying the interval of a Seventh.
	Seventh Interval = 7
	// Ninth For specifying the interval of a Ninth.
	Ninth Interval = 9
	// Eleventh For specifying the interval of an Eleventh.
	Eleventh Interval = 11
	// Thirteenth For specifying the interval of a Thirteenth.
	Thirteenth Interval = 13

	// These are the half-step distances from the First

	// MinorSecond The interval of a minor second, in half steps.
	MinorSecond = HalfStepValue
	// MajorSecond The interval of a major second, in half steps.
	MajorSecond = WholeStepValue
	// MinorThird The interval of a minor third, in half steps.
	MinorThird = HalfStepValue * 3
	// MajorThird The interval of a major third, in half steps.
	MajorThird = HalfStepValue * 4
	// PerfectFourth The interval of a perfect fourth, in half steps.
	PerfectFourth = HalfStepValue * 5
	// PerfectFifth The interval of a perfect fifth, in half steps.
	PerfectFifth = HalfStepValue * 7
)

// Chord A collection of specific pitches, making a chord.
type Chord struct {
	pitches []Pitch // The pitches that make up this chord
	root    int     // The index into the pitch slice that contains the root pitch. Will be non-zero for inverted chords.
}

// ChordFactory The purpose of this class is to allow creating chords using scale intervals, without needing to worry
// about half steps, e.g., specify "third" without having to know if it's a major (2 steps) or minor (three half steps) third.
type ChordFactory struct {
	pattern Pattern // The pattern of the scale, which must be diatonic
	root    *Pitch  // The pitch to apply the pattern from
	offset  int     // The offset into the scale of the root
}

// GetPitch Get the pitch that is the given interval from this factory's root. The interval should be 1 (a first,
// which will return the same pitch) or higher. One is zero? Yes, that's just how music works ¯\_(ツ)_/¯.
func (f *ChordFactory) GetPitch(interval Interval) *Pitch {
	var halfSteps HalfSteps
	for i := 1; i < int(interval)+f.offset; i++ {
		halfSteps += f.pattern.At(i - 1)
	}
	return f.root.GetTransposedCopy(halfSteps)
}

// GetIntervalSize Gets the size of the given interval in half steps.
func (f *ChordFactory) GetIntervalSize(interval Interval) HalfSteps {
	return f.root.GetDistanceTo(f.GetPitch(interval))
}

// ContainsInterval Returns true if the pitch at the given interval is the given number of half steps from its root.
func (f *ChordFactory) ContainsInterval(interval Interval, halfSteps HalfSteps) bool {
	return f.GetIntervalSize(interval) == halfSteps
}

// HasMajorThird Returns true if this factory's third is a major third, false if it is a minor third.
func (f *ChordFactory) HasMajorThird() bool {
	return f.ContainsInterval(Third, MajorThird)
}

// HasPerfectFourth Returns true if this factory's fourth is a perfect fourth, false if it is a minor third.
func (f *ChordFactory) HasPerfectFourth() bool {
	return f.ContainsInterval(Fourth, PerfectFourth)
}

// HasPerfectFifth Returns true if this factory's fifth is a perfect fifth, false if it is a minor third.
func (f *ChordFactory) HasPerfectFifth() bool {
	return f.ContainsInterval(Fifth, PerfectFifth)
}

// CreateChord Create a chord using the specified intervals.
func (f *ChordFactory) CreateChord(intervals ...Interval) *Chord {
	pitches := make([]Pitch, len(intervals), len(intervals))
	for i, v := range intervals {
		pitches[i] = *f.GetPitch(v)
	}
	return &Chord{pitches, 0}
}

// Given a set of two or more pitch classes, how do we determine what the chord is?
// Does the key matter? Only in determining the Roman Numeral to use.
// Only the pattern of half step intervals matters.
// For each (type of) chord we know, we can iterate through all its inversions, adding them to a trie.

// TODO: It turns out some chords are ambiguous, and the right name to use depends on the lowest pitch so we need to
// provide an algorithm that uses pitches rather than classes.

func createTriadPattern(third1 HalfSteps, third2 HalfSteps) *Pattern {
	// Necessary to make a pattern that spans exactly an octave, as the third interval is required for chord inversions
	return MakePattern(third1, third2, OctaveValue-(third1+third2))
}

func createTetrachordPattern(third1 HalfSteps, third2 HalfSteps, third3 HalfSteps) *Pattern {
	// Necessary to make a pattern that spans exactly an octave, as the third interval is required for chord inversions
	return MakePattern(third1, third2, third3, OctaveValue-(third1+third2+third3))
}

// CreatePowerChordPattern Creates the pattern for a "power chord", consisting of the root and the fifth.
func CreatePowerChordPattern() *Pattern {
	return MakePattern(PerfectFifth, PerfectFourth)
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

// CreateAugmentedTriadPattern Creates the pattern for an augmented triad.
func CreateAugmentedTriadPattern() *Pattern {
	return createTriadPattern(MajorThird, MajorThird)
}

// CreateSuspendedPattern Creates the pattern for a suspended second triad, where the third is omitted, and a fourth/second is added.
// Important: A suspended second is equivalent to a suspended fourth with the fifth as the root.
func CreateSuspendedPattern() *Pattern {
	return createTriadPattern(PerfectFourth, MajorSecond)
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
	dict = &PatternDictionary{NewTrie(1, OctaveValue)}

	addChordToDict(dict, CreatePowerChordPattern(), "5")

	addChordToDict(dict, CreateMajorTriadPattern(), " Major")           // 4, 3
	addChordToDict(dict, CreateMinorTriadPattern(), " Minor")           // 3, 4
	addChordToDict(dict, CreateDiminishedTriadPattern(), " Diminished") // 3, 3
	addChordToDict(dict, CreateAugmentedTriadPattern(), " Augmented")   // 4, 4

	addChordToDict(dict, CreateSuspendedPattern(), " Suspended")

	addChordToDict(dict, CreateMajorSeventhPattern(), " Major Seventh")
	addChordToDict(dict, CreateDominantSeventhPattern(), " Dominant Seventh")
	addChordToDict(dict, CreateMinorSeventhPattern(), " Minor Seventh")
	addChordToDict(dict, CreateDiminishedSeventhPattern(), " Diminished Seventh")

	return
}

// GetChordName Given a set of unique pitch classes, returns the name of the produced chord. If the chord doesn't contain a name in
// the given dictionary, then ("", false) is returned.
// As this function takes pitch classes, it cannot determine extended chord names, as the pitches would loop back around (11th -> 4th).
func GetChordName(dict *PatternDictionary, pitchNamer *PitchNamer, chord []PitchClass) (name string, ok bool) {

	// Assumption: only unique pitch classes are in chord

	if len(chord) < 2 {
		// I've made the executive decision to allow power chords, even though, featuring only two pitches, they
		// aren't technically chords
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

	// 3. Look up the pattern in the dictionary

	// For now we're just taking the first name; at some point, getting all the possible names will need to be an option.

	entries := dict.GetEntries(&Pattern{pattern})
	for _, entry := range entries {
		e, ok := entry.(*chordDictionaryEntry)
		if ok {
			root := pitchNamer.Name(PitchClass{HalfSteps(pitchIndices[e.rootIndex])})
			name = fmt.Sprintf("%s%s", root, e.name)
			return name, true
		}
	}

	ok = false

	return
}
