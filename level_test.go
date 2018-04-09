package log

import "testing"

func TestLevel_String(t *testing.T) {
	s := Level(-1).String()
	for l := LevelAll + 1; l < LevelNone; l++ {
		if l.String() == s {
			t.Errorf("Level(%d).String() returned %s", int(l), s)
		}
	}
}
