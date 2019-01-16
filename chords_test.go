package tonacity

import (
	"testing"
)

func TestChordFactory_GetPitch(t *testing.T) {
	type args struct {
		interval Interval
	}
	pf := &PitchFactory{*MiddleC()}
	scale := CreateMajorScale()
	tests := []struct {
		name string
		f    *ChordFactory
		args args
		want *Pitch
	}{
		{"C Major - C Major Third", &ChordFactory{*scale, pf.GetPitch(C(), 4), 0}, args{Third}, pf.GetPitch(E(), 4)},
		{"C Major - C Major Fifth", &ChordFactory{*scale, pf.GetPitch(C(), 4), 0}, args{Fifth}, pf.GetPitch(G(), 4)},
		{"C Major - E Minor Third", &ChordFactory{*scale, pf.GetPitch(C(), 4), 2}, args{Third}, pf.GetPitch(G(), 4)},
		{"A Major - E Major Third", &ChordFactory{*scale, pf.GetPitch(A(), 4), 4}, args{Third}, pf.GetPitch(G().Sharp(), 5)},
		{"A Major - E Major Fifth", &ChordFactory{*scale, pf.GetPitch(A(), 4), 4}, args{Fifth}, pf.GetPitch(B(), 5)},
		{"A Major - E Major Seventh", &ChordFactory{*scale, pf.GetPitch(A(), 4), 4}, args{Seventh}, pf.GetPitch(D(), 6)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.GetPitch(tt.args.interval); got.GetDistanceTo(tt.want) != 0 {
				t.Errorf("ChordFactory.GetPitchClass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChordName(t *testing.T) {
	type args struct {
		dict       *PatternDictionary
		pitchNamer *PitchNamer
		chord      []PitchClass
	}
	dict := CreateChordDictionary()
	namer := CreateSharpPitchNamer()
	tests := []struct {
		name     string
		args     args
		wantName string
		wantOk   bool
	}{
		{`Power "Chord"`, args{dict, namer, []PitchClass{*C(), *G()}}, "C5", true},
		{"G Dom7", args{dict, namer, []PitchClass{*G(), *B(), *D(), *F()}}, "G Dominant Seventh", true},
		{"G Maj", args{dict, namer, []PitchClass{*G(), *B(), *D(), *F().Sharp()}}, "G Major Seventh", true},
		{"C Major", args{dict, namer, []PitchClass{*C(), *E(), *G()}}, "C Major", true},
		{"E/CMaj", args{dict, namer, []PitchClass{*E(), *G(), *C()}}, "C Major", true},
		{"B Dim", args{dict, namer, []PitchClass{*B(), *D(), *F()}}, "B Diminished", true},
		{"C♯ Sus", args{dict, namer, []PitchClass{*F().Sharp(), *C().Sharp(), *G().Sharp()}}, "C♯ Suspended", true},
		{"D Augmented", args{dict, namer, []PitchClass{*D(), *F().Sharp(), *A().Sharp()}}, "D Augmented", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotOk := GetChordName(tt.args.dict, tt.args.pitchNamer, tt.args.chord)
			if gotName != tt.wantName {
				t.Errorf("GetChordName() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("GetChordName() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
