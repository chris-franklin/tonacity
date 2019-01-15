package tonacity

import "testing"

func TestPatternDictionary_GetName(t *testing.T) {
	type args struct {
		pattern *Pattern
	}

	modeDict := BuildModeDictionary()
	scaleDict := BuildScaleDictionary()

	tests := []struct {
		name  string
		d     *PatternDictionary
		args  args
		want  string
		want1 bool
	}{
		{"Ionian", modeDict, args{CreateIonianMode()}, "Ionian", true},
		{"Dorian", modeDict, args{CreateDorianMode()}, "Dorian", true},
		{"Aeolian", modeDict, args{CreateAeolianMode()}, "Aeolian", true},
		{"Minor", scaleDict, args{CreateAeolianMode()}, "Minor", true},
		{"Pentatonic Minor", scaleDict, args{CreateMinorPentatonicScalePattern()}, "Pentatonic Minor", true},
		{"Unknown", modeDict, args{CreateMinorPentatonicScalePattern()}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.d.GetName(tt.args.pattern)
			if got != tt.want {
				t.Errorf("PatternDictionary.GetName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PatternDictionary.GetName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
