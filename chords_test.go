package tonacity

import (
	"testing"
)

func TestChordFactory_GetPitchClass(t *testing.T) {
	type args struct {
		interval uint8
	}
	tests := []struct {
		name string
		f    *ChordFactory
		args args
		want *PitchClass
	}{
		{"C Major - C Major Third", &ChordFactory{*CreateMajorScale(), *C(), 0}, args{3}, E()},
		{"C Major - C Major Fifth", &ChordFactory{*CreateMajorScale(), *C(), 0}, args{5}, G()},
		{"C Major - E Minor Third", &ChordFactory{*CreateMajorScale(), *C(), 2}, args{3}, G()},
		{"A Major - E Major Third", &ChordFactory{*CreateMajorScale(), *A(), 4}, args{3}, G().Sharp()},
		{"A Major - E Major Fifth", &ChordFactory{*CreateMajorScale(), *A(), 4}, args{5}, B()},
		{"A Major - E Major Seventh", &ChordFactory{*CreateMajorScale(), *A(), 4}, args{7}, D()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.GetPitchClass(tt.args.interval); !got.HasSamePitchAs(tt.want) {
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
		{"G Dom7", args{dict, namer, []PitchClass{*G(), *B(), *D(), *F()}}, "G Dominant Seventh", true},
		{"G Maj", args{dict, namer, []PitchClass{*G(), *B(), *D(), *F().Sharp()}}, "G Major Seventh", true},
		{"C Major", args{dict, namer, []PitchClass{*C(), *E(), *G()}}, "C Major", true},
		{"E/CMaj", args{dict, namer, []PitchClass{*E(), *G(), *C()}}, "C Major", true},
		{"B Dim", args{dict, namer, []PitchClass{*B(), *D(), *F()}}, "B Diminished", true},
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
