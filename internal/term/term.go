package term

import (
	"os"

	"github.com/mattn/go-isatty"
)

// Plain diz se a saída deve ser sem cor/estilo (pipe, CI, NO_COLOR ou não-TTY).
func Plain() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return true
	}
	if os.Getenv("CI") != "" || os.Getenv("TERM") == "dumb" {
		return true
	}
	return !isatty.IsTerminal(os.Stdout.Fd())
}

// NoEmoji diz se emojis devem ser suprimidos (acessibilidade, fonte Windows, CI).
func NoEmoji() bool {
	if Plain() {
		return true
	}
	if _, ok := os.LookupEnv("DEVBOX_NO_EMOJI"); ok {
		return true
	}
	return false
}
