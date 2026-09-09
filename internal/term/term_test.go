package term

import "testing"

func TestPlainRespectsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if !Plain() {
		t.Fatal("esperava plain com NO_COLOR=1")
	}
}

func TestPlainRespectsCI(t *testing.T) {
	t.Setenv("CI", "true")
	if !Plain() {
		t.Fatal("esperava plain com CI=true")
	}
}

func TestNoEmojiRespectsFlag(t *testing.T) {
	t.Setenv("DEVBOX_NO_EMOJI", "1")
	if !NoEmoji() {
		t.Fatal("esperava NoEmoji com DEVBOX_NO_EMOJI=1")
	}
}
