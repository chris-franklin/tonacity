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
