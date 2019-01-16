package tonacity

const (
	// NotesInMode The number of notes that are in a Mode
	NotesInMode = 7
)

type Pattern struct {
	intervals []HalfSteps
}

func MakePattern(intervals ...HalfSteps) *Pattern {
	return &Pattern{intervals}
}

func (p Pattern) Intervals() []HalfSteps {
	return p.Copy().intervals
}

func (p *Pattern) At(index int) HalfSteps {
	return p.intervals[index%len(p.intervals)]
}

func (p *Pattern) Length() int {
	return len(p.intervals)
}

func (p Pattern) Copy() *Pattern {
	c := make([]HalfSteps, len(p.intervals), len(p.intervals))
	copy(c, p.intervals)
	return &Pattern{c}
}

// Offset Returns a copy of this pattern offset so it starts at the interval at the given offset
func (p Pattern) Offset(o int) *Pattern {
	l := len(p.intervals)
	if o < 0 {
		o = l + o
	}
	offsetPattern := make([]HalfSteps, l, l)
	for i := 0; i < l; i++ {
		offsetPattern[i] = p.intervals[(i+o)%l]
	}
	return &Pattern{offsetPattern}
}

// Reverse Reverses this pattern so that it produces the reverse of the pattern passed in.
// For example: a major scale pattern in will produce a descending major scale pattern.
// Super simple example: {3, 2, 1} will produce {-1, -2, -3}.
func (p *Pattern) Reverse() *Pattern {
	l := len(p.intervals)
	reverse := make([]HalfSteps, l, l)
	for i := 0; i < l; i++ {
		reverse[i] = -p.intervals[l-i]
	}
	return &Pattern{reverse}
}

// ScaleDegree Used for specifying the degree of a scale
type ScaleDegree uint8

const (
	// Tonic The key note of a scale
	Tonic ScaleDegree = 1
	// Supertonic A second above the tonic
	Supertonic ScaleDegree = 2
	// Mediant A third above the tonic
	Mediant ScaleDegree = 3
	// Subdominant A fifth below the tonic
	Subdominant ScaleDegree = 4
	// Dominant A fifth above the tonic
	Dominant ScaleDegree = 5
	// Submediant A third below the tonic
	Submediant ScaleDegree = 6
	// LeadingTone A second below the tonic
	LeadingTone ScaleDegree = 7
)

// PatternRepeatsAtOctave Returns true if the given pattern will repeat at the next octave up from where it starts.
func (p *Pattern) PatternRepeatsAtOctave() bool {
	var sum HalfSteps
	for _, v := range p.intervals {
		sum += v
	}
	return sum == OctaveValue
}

// CreateAscendingSinger Creates a singer that applies this pattern from the given pitch.
func (p Pattern) CreateAscendingSinger(pitch Pitch) Singer {
	return PatternRepeatingSinger{p, pitch, 0}
}

// CreateDescendingSinger Creates a singer that applies this pattern in reverse from the given pitch.
func (p Pattern) CreateDescendingSinger(pitch Pitch) Singer {
	return PatternRepeatingSinger{*p.Reverse(), pitch, 0}
}

// ionianModePattern The pattern of the Ionian (I) Mode, repeated twice to allow slicing it to create on of the other modes.
var ionianModePattern = MakePattern()

// CreateIonianMode Creates the pattern of the Ionian Mode.
func CreateIonianMode() *Pattern {
	return MakePattern(
		WholeStepValue,
		WholeStepValue,
		HalfStepValue,
		WholeStepValue,
		WholeStepValue,
		WholeStepValue,
		HalfStepValue)
}

// CreateDorianMode Creates the pattern of the Dorian Mode.
func CreateDorianMode() *Pattern {
	return CreateIonianMode().Offset(1)
}

// CreatePhrygianMode Creates the pattern of the Phrygian Mode.
func CreatePhrygianMode() *Pattern {
	return CreateIonianMode().Offset(2)
}

// CreateLydianMode Creates the pattern of the Lydian Mode.
func CreateLydianMode() *Pattern {
	return CreateIonianMode().Offset(3)
}

// CreateMixolydianMode Creates the pattern of the Mixolydian Mode.
func CreateMixolydianMode() *Pattern {
	return CreateIonianMode().Offset(4)
}

// CreateAeolianMode Creates the pattern of the Aeolian Mode.
func CreateAeolianMode() *Pattern {
	return CreateIonianMode().Offset(5)
}

// CreateLocrianMode Creates the pattern of the Locrian Mode.
func CreateLocrianMode() *Pattern {
	return CreateIonianMode().Offset(6)
}

// CreateMajorScale Get the pattern of the Major Scale
func CreateMajorScale() *Pattern {
	return CreateIonianMode()
}

// CreateMinorScale Get the pattern of the Minor Scale
func CreateMinorScale() *Pattern {
	return CreateAeolianMode()
}

// CreateModes Returns all seven mode patterns, in order from I (Ionian) to VII (Locrian)
func CreateModes() []*Pattern {
	patterns := []*Pattern{
		CreateIonianMode(),
		CreateDorianMode(),
		CreatePhrygianMode(),
		CreateLydianMode(),
		CreateMixolydianMode(),
		CreateAeolianMode(),
		CreateLocrianMode(),
	}
	return patterns
}

var harmonicMinorScalePattern = MakePattern(
	WholeStepValue,
	HalfStepValue,
	WholeStepValue,
	WholeStepValue,
	HalfStepValue,
	WholeStepValue+HalfStepValue,
	HalfStepValue,
)

// CreateHarmonicMinorScalePattern Get the pattern of the Harmonic Minor Scale.
func CreateHarmonicMinorScalePattern() *Pattern {
	return CreateMajorScale().Offset(9)
}

var melodicMinorScalePattern = MakePattern(
	WholeStepValue,
	HalfStepValue,
	WholeStepValue,
	WholeStepValue,
	WholeStepValue,
	WholeStepValue,
	HalfStepValue,
	-WholeStepValue,
	-WholeStepValue,
	-HalfStepValue,
	-WholeStepValue,
	-WholeStepValue,
	-HalfStepValue,
	-WholeStepValue,
)

// CreateMelodicMinorAscendingScalePattern Get the pattern of the ascending Melodic Minor Scale.
func CreateMelodicMinorAscendingScalePattern() *Pattern {
	return MakePattern(WholeStepValue,
		HalfStepValue,
		WholeStepValue,
		WholeStepValue,
		WholeStepValue,
		WholeStepValue,
		HalfStepValue)
}

// CreateMelodicMinorDescendingScalePattern Get the pattern of the descending Melodic Minor Scale.
func CreateMelodicMinorDescendingScalePattern() *Pattern {
	return MakePattern(-WholeStepValue,
		-WholeStepValue,
		-HalfStepValue,
		-WholeStepValue,
		-WholeStepValue,
		-HalfStepValue,
		-WholeStepValue,
	)
}

// CreateChromaticScalePattern Get the pattern of the Chromatic Scale.
func CreateChromaticScalePattern() *Pattern {
	return MakePattern(
		HalfStepValue,
		HalfStepValue,
		HalfStepValue,
		HalfStepValue,
		HalfStepValue,
		HalfStepValue,
		HalfStepValue,
		HalfStepValue,
		HalfStepValue,
		HalfStepValue,
		HalfStepValue,
		HalfStepValue,
	)
}

// SymmetricScalePattern The Symmetric Scale Pattern
var symmetricScalePattern = MakePattern(
	WholeStepValue,
	HalfStepValue,
	WholeStepValue,
	HalfStepValue,
	WholeStepValue,
	HalfStepValue,
	WholeStepValue,
	HalfStepValue,
)

// CreateSymmetricScalePattern Get the pattern of the Symmetric Scale.
func CreateSymmetricScalePattern() *Pattern {
	return MakePattern(
		WholeStepValue,
		HalfStepValue,
		WholeStepValue,
		HalfStepValue,
		WholeStepValue,
		HalfStepValue,
		WholeStepValue,
		HalfStepValue,
	)
}

// CreateMinorPentatonicScalePattern Get the pattern of the Pentatonic Minor Scale.
func CreateMinorPentatonicScalePattern() *Pattern {
	return MakePattern(
		WholeStepValue+HalfStepValue,
		WholeStepValue,
		WholeStepValue,
		WholeStepValue+HalfStepValue,
		WholeStepValue,
	)
}

// CreateMajorPentatonicScalePattern Get the pattern of the Pentatonic Major Scale.
func CreateMajorPentatonicScalePattern() *Pattern {
	return MakePattern(
		WholeStepValue,
		WholeStepValue,
		WholeStepValue+HalfStepValue,
		WholeStepValue,
		WholeStepValue+HalfStepValue,
	)
}

// PatternDictionary A dictionary for looking up the name of a praticular pattern of intervals.
type PatternDictionary struct {
	searchTree *TrieNode
}

// GetName If this dictionary contains the given pattern, then its name will be returned. The bool will be false if the pattern is
// not in the dictionary.
func (d *PatternDictionary) GetName(pattern *Pattern) (string, bool) {
	value := d.searchTree.FindValue(pattern.intervals)
	s, ok := value.(string)
	return s, ok
}

// GetEntry If this dictionary contains the given pattern, then its entry will be returned, otherwise nil.
func (d *PatternDictionary) GetEntry(pattern *Pattern) interface{} {
	value := d.searchTree.FindValue(pattern.intervals)
	return value
}

// AddPattern Add a pattern to the dictionary.
func (d *PatternDictionary) AddPattern(pattern *Pattern, entry interface{}) {
	d.searchTree.AddValue(pattern.Intervals(), entry)
}

// BuildModeDictionary Builds a dictionary containing the seven modes.
func BuildModeDictionary() *PatternDictionary {
	dict := &PatternDictionary{NewTrie(3)}
	dict.AddPattern(CreateIonianMode(), "Ionian")
	dict.AddPattern(CreateDorianMode(), "Dorian")
	dict.AddPattern(CreatePhrygianMode(), "Phrygian")
	dict.AddPattern(CreateLydianMode(), "Lydian")
	dict.AddPattern(CreateMixolydianMode(), "Mixolydian")
	dict.AddPattern(CreateAeolianMode(), "Aeolian")
	dict.AddPattern(CreateLocrianMode(), "Locrian")
	return dict
}

// BuildScaleDictionary Builds a dictionary containing the standard scales.
func BuildScaleDictionary() *PatternDictionary {
	dict := &PatternDictionary{NewTrie(4)}
	dict.AddPattern(CreateMajorScale(), "Major")
	dict.AddPattern(CreateMinorScale(), "Minor")
	dict.AddPattern(CreateHarmonicMinorScalePattern(), "Harmonic Minor")
	dict.AddPattern(CreateMelodicMinorAscendingScalePattern(), "Ascending Melodic Minor")
	dict.AddPattern(CreateMajorPentatonicScalePattern(), "Pentatonic Major")
	dict.AddPattern(CreateMinorPentatonicScalePattern(), "Pentatonic Minor")
	return dict
}

// RootedPattern A pattern that is rooted at a specific pitch.
type RootedPattern struct {
	pattern Pattern
	root    Pitch
}

// Root Returns the root note of this pattern.
func (rp *RootedPattern) Root() Pitch {
	return rp.root
}

// Transpose Transposes the root note by the given number of half steps.
func (rp *RootedPattern) Transpose(halfSteps HalfSteps) {
	rp.root.Transpose(halfSteps)
}

// CreateSinger Create and return a new singer that will sing this pattern, starting at its root.
func (rp *RootedPattern) CreateSinger() *PatternRepeatingSinger {
	return &PatternRepeatingSinger{*rp.pattern.Copy(), rp.root, 0}
}

// CreateReverseSinger Create and return a new singer that will sing this pattern, starting at its root, and applying the pattern in reverse.
func (rp *RootedPattern) CreateReverseSinger() *PatternRepeatingSinger {
	return &PatternRepeatingSinger{*rp.pattern.Reverse(), rp.root, 0}
}
