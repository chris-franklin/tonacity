package tonacity

import "testing"

func TestTrieNode_AddValue(t *testing.T) {
	type args struct {
		path  []HalfSteps
		value interface{}
	}
	tests := []struct {
		name string
		node *TrieNode
		args args
		want string
	}{
		{"0-2", NewTrie(0, 2), args{[]HalfSteps{0, 1, 2}, "A"}, "A"},
		{"1-3", NewTrie(1, 3), args{[]HalfSteps{1, 2, 3}, "B"}, "B"},
		{"(-5)-(-3)", NewTrie(-5, -3), args{[]HalfSteps{-5, -4, -3}, "C"}, "C"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.AddValue(tt.args.path, tt.args.value)
			gotName := tt.node.FindValue(tt.args.path)
			if len(gotName) != 1 || gotName[0] != tt.want {
				t.Errorf("FindValue() gotName = %v, want %v", gotName, tt.want)
			}
		})
	}
}
