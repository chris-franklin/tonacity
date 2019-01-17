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
	// Seventh For specifying the interval of a Seventh. Limit on a piano for average female hands.
	Seventh Interval = 7
	// Ninth For specifying the interval of a Ninth. Limit on a piano for average male hands.
	Ninth Interval = 9
	// Eleventh For specifying the interval of an Eleventh. You need large hands to reach this on a piano.
	Eleventh Interval = 11
	// Thirteenth For specifying the interval of a Thirteenth. On a piano this is 29cm on the white keys, Rachmaninov (6'6") and Liszt could manage it.
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
}

func MakeChord(pitches ...Pitch) *Chord {
	return &Chord{pitches}
}

func (c *Chord) String() string {
	return fmt.Sprintf("%v", c.pitches)
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
	return &Chord{pitches}
}

// Given a set of two or more pitch classes, how do we determine what the chord is?
// Does the key matter? Only in determining the Roman Numeral to use.
// Only the pattern of half step intervals matters.
// For each (type of) chord we know, we can iterate through all its inversions, adding them to a trie.

func createTriadPattern(third1 HalfSteps, third2 HalfSteps) *Pattern {
	// Necessary to make a pattern that spans exactly an octave, as the third interval is required for chord inversions
	return MakePattern(third1, third2)
}

func createTetrachordPattern(third1 HalfSteps, third2 HalfSteps, third3 HalfSteps) *Pattern {
	// Necessary to make a pattern that spans exactly an octave, as the third interval is required for chord inversions
	return MakePattern(third1, third2, third3)
}

// CreatePowerChordPattern Creates the pattern for a "power chord", consisting of the root and the fifth.
func CreatePowerChordPattern() *Pattern {
	return MakePattern(PerfectFifth)
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

func addChordWithInversionsToDict(dict *PatternDictionary, chord *Pattern, name string) {
	// Add the chord in root position
	dict.AddPattern(chord, &chordDictionaryEntry{name, 0})
	for i := chord.Length(); i > 0; i-- {
		// Get the next inversion of the chord
		chord.Invert()
		// dict doesn't keep reference to chord, so mutations after the call don't matter
		dict.AddPattern(chord, &chordDictionaryEntry{name, i})
	}
}

func addChordToDict(dict *PatternDictionary, chord *Pattern, name string) {
	dict.AddPattern(chord, &chordDictionaryEntry{name, 0})
}

// CreateChordDictionary Creates a new dictionary for the purpose of naming chords based on the half step intervals between pitches.
func CreateChordDictionary() (dict *PatternDictionary) {
	dict = &PatternDictionary{NewTrie(1, OctaveValue)}

	addChordWithInversionsToDict(dict, CreatePowerChordPattern(), "5")

	addChordWithInversionsToDict(dict, CreateMajorTriadPattern(), " Major")           // 4, 3
	addChordWithInversionsToDict(dict, CreateMinorTriadPattern(), " Minor")           // 3, 4
	addChordWithInversionsToDict(dict, CreateDiminishedTriadPattern(), " Diminished") // 3, 3
	// Inversion of an augmented chord is an augmented chord with a different root (because 12-(4+4)=4)
	addChordToDict(dict, CreateAugmentedTriadPattern(), " Augmented") // 4, 4

	addChordWithInversionsToDict(dict, CreateSuspendedPattern(), " Suspended")

	addChordWithInversionsToDict(dict, CreateMajorSeventhPattern(), " Major Seventh")
	addChordWithInversionsToDict(dict, CreateDominantSeventhPattern(), " Dominant Seventh")
	addChordWithInversionsToDict(dict, CreateMinorSeventhPattern(), " Minor Seventh")
	addChordWithInversionsToDict(dict, CreateDiminishedSeventhPattern(), " Diminished Seventh")

	return
}

// GetName will return the name of this chord, if its intervals are a valid pattern in the given dictionary. This function is
// specifically preferable for guitars or similar, where extended chords (those with ninths - a stretch on a piano, elevenths, and thirteenths) are used more.
func (c *Chord) GetName(dict *PatternDictionary, pitchNamer *PitchNamer) (name string, ok bool) {
	sort.Sort(ByPitch(c.pitches))

	intervals := make([]HalfSteps, len(c.pitches)-1, len(c.pitches)-1)

	for i := 0; i < len(intervals); i++ {
		intervals[i] = c.pitches[i].GetDistanceTo(&c.pitches[i+1])
	}

	entries := dict.GetEntries(&Pattern{intervals})
	for _, entry := range entries {
		e, ok := entry.(*chordDictionaryEntry)
		if ok {
			firstNote := pitchNamer.Name(c.pitches[0].Class())
			if e.rootIndex == 0 {
				// Chord is in root position
				name = fmt.Sprintf("%s%s", firstNote, e.name)
			} else {
				// Chord is inverted
				root := pitchNamer.Name(c.pitches[e.rootIndex].Class())
				name = fmt.Sprintf("%s%s/%s", root, e.name, firstNote)
			}
			return name, true
		}
	}
	ok = false
	return
}

// GetChordName Given a set of unique pitch classes, returns the name of the produced chord. If the chord doesn't contain a name in
// the given dictionary, then ("", false) is returned.
// As this function takes pitch classes, it cannot determine extended chord names, as the pitches would loop back around (11th -> 4th).
// It also cannot give the "/<low note>" modifier on an inverted chord, as pitch ordering is lost.
// This function is useful for when distinct pitches are far apart (left and right hands on piano) but do combine to make a chord.
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

	// 2. Create the pattern of the intervals between pitches

	pattern := make([]HalfSteps, 0, len(pitchIndices))

	for i, v := range pitchIndices {
		if i+1 < len(pitchIndices) {
			interval := HalfSteps(pitchIndices[(i+1)] - v)
			pattern = append(pattern, interval)
		}
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
