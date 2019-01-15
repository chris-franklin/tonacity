package tonacity

import (
	"math"
	"reflect"
	"testing"
)

func TestPitch_Octave(t *testing.T) {
	type args struct {
		middleC *Pitch
	}
	middleC := MiddleC()
	tests := []struct {
		name string
		p    *Pitch
		args args
		want int8
	}{
		{"Bb3", middleC.GetTransposedCopy(-1), args{middleC}, 3},
		{"Middle C", middleC, args{middleC}, 4},
		{"Same Octave B flat", middleC.GetTransposedCopy(11), args{middleC}, 4},
		{"C5", middleC.GetTransposedCopy(12), args{middleC}, 5},
		{"C sharp 3", middleC.GetTransposedCopy(-11), args{middleC}, 3},
		{"C3", middleC.GetTransposedCopy(-12), args{middleC}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Octave(tt.args.middleC); got != tt.want {
				t.Errorf("Pitch.Octave() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPitchClass_GetDistanceToHigherPitchClass(t *testing.T) {
	type args struct {
		b PitchClass
	}
	tests := []struct {
		name string
		pc   *PitchClass
		args args
		want HalfSteps
	}{
		{"C to C", C(), args{*C()}, 12},
		{"C to D", C(), args{*D()}, 2},
		{"C to B", C(), args{*B()}, 11},
		{"A to C", A(), args{*C()}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pc.GetDistanceToHigherPitchClass(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PitchClass.GetDistanceToHigherPitchClass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPitchClass_GetDistanceToLowerPitchClass(t *testing.T) {
	type args struct {
		b PitchClass
	}
	tests := []struct {
		name string
		pc   *PitchClass
		args args
		want HalfSteps
	}{
		{"C to C", C(), args{*C()}, -12},
		{"D to C", D(), args{*C()}, -2},
		{"C to B", C(), args{*B()}, -1},
		{"C to A", C(), args{*A()}, -3},
		{"A to C", A(), args{*C()}, -9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pc.GetDistanceToLowerPitchClass(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PitchClass.GetDistanceToLowerPitchClass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPitch_FrequencyInHertz(t *testing.T) {
	type args struct {
		concertPitch float64
	}
	tests := []struct {
		name string
		p    *Pitch
		args args
		want float64
	}{
		{"A4", A4(), args{16}, 16},
		{"A5", A4().GetTransposedCopy(12), args{16}, 32},
		{"A3", A4().GetTransposedCopy(-12), args{16}, 8},
		{"A2", A4().GetTransposedCopy(-24), args{16}, 4},
		{"A1", A4().GetTransposedCopy(-36), args{16}, 2},
		{"A0", A4().GetTransposedCopy(-48), args{16}, 1},
		{"C4", MiddleC(), args{StandardConcertPitch}, 261.626},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.FrequencyInHertz(tt.args.concertPitch); math.Abs(got-tt.want) > 0.001 {
				t.Errorf("Pitch.FrequencyInHertz() = %v, want %v", got, tt.want)
			}
		})
	}
}
