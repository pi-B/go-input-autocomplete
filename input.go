package input_autocomplete

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/term"
)

type Input struct {
	cursor           *Cursor
	fixedText        string
	currentText      string
	isCycling        bool
	cyclingPos       int
	matches          []string
	hasMultiline     bool
	lastAutocomplete string
}

func NewInput(fixedText string) *Input {
	return &Input{
		cursor:           NewCursor(),
		fixedText:        fixedText,
		currentText:      "",
		isCycling:        false,
		cyclingPos:       0,
		matches:          []string{},
		hasMultiline:     false,
		lastAutocomplete: "",
	}
}

func (i *Input) canDeleteChar() bool {
	return i.cursor.GetPosition() >= 1
}

func (i *Input) AddChar(char rune) {
	i.isCycling = false
	pos := i.cursor.GetPosition()
	c := string(char)

	if pos == len(i.currentText) {
		i.currentText += c
		fmt.Print(c)
		i.cursor.IncrementPosition()
	} else {
		aux := len(i.currentText) - pos
		i.currentText = i.currentText[:pos] + c + i.currentText[pos:]
		i.cursor.SetPosition(len(i.currentText))
		i.Print()
		i.cursor.MoveLeftNPos(aux)
	}
}

func (i *Input) RemoveChar() {
	i.isCycling = false
	if i.canDeleteChar() {
		pos := i.cursor.GetPosition()
		aux := len(i.currentText) - pos
		i.currentText = i.currentText[:pos-1] + i.currentText[pos:]
		i.cursor.SetPosition(len(i.currentText))
		i.Print()
		i.cursor.MoveLeftNPos(aux)
	}
}

func (i *Input) MoveCursorLeft() {
	i.isCycling = false
	i.cursor.MoveLeft()
}

func (i *Input) MoveCursorRight() {
	i.isCycling = false
	if i.cursor.GetPosition() < len(i.currentText) {
		i.cursor.MoveRight()
	}
}

func (i *Input) Autocomplete() {
	if !i.isCycling {
		i.isCycling = true
		i.cyclingPos = 0
		i.matches = Autocomplete(i.currentText)
		if len(i.matches) <= 1 {
			i.isCycling = false
		}
		if len(i.matches) == 0 {
			return
		}
	}

	if i.currentText != i.lastAutocomplete {
		if len(i.matches) == 1 {
			fmt.Print("\033[J\033[G\033[K")
			fmt.Print(i.fixedText + i.matches[0])
			i.currentText = i.matches[0]
			i.hasMultiline = false
		} else {
			i.PrintAllMatches()
			i.lastAutocomplete = i.currentText
		}

		i.cursor.SetPosition(len(i.currentText))

	}

}

func (i *Input) PrintAllMatches() {
	var max_len int

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Can't get the terminal size :", err)
		return
	}

	for _, match := range i.matches {
		if len(match) > max_len {
			max_len = len(match) + 1
		}
	}

	display_slice := make([]string, len(i.matches))
	res := ""
	current_line := ""
	line_nb := 1

	for j, match := range i.matches {
		display_slice[j] = match
		for k := 0; k < max_len-len(match); k++ {
			display_slice[j] += " "
		}
		if len(current_line)+len(display_slice[j]) > width {
			res += current_line + "\n"
			current_line = ""
			line_nb++
			i.hasMultiline = true
		} else {
			current_line += display_slice[j] + " "
		}
	}
	res += current_line
	fmt.Print("\033[G\033[K")
	fmt.Println("")
	fmt.Print(res)
	if i.hasMultiline {
		fmt.Print("\033[J")
	}
	fmt.Printf("\033[%vA\033[G", line_nb)
	if line_nb > 1 {
		i.hasMultiline = true
	}
	fmt.Print(i.fixedText + i.currentText)
}

func (i *Input) RemoveLastSlashIfNeeded() {
	os := runtime.GOOS
	size := len(i.currentText)
	var slash byte

	switch os {
	case "linux", "darwin":
		slash = '/'
	case "windows":
		slash = '\\'
	}

	if size > 0 && i.currentText[size-1] == slash {
		i.currentText = i.currentText[:size-1]
	}
}

func (i *Input) Print() {
	fmt.Print("\033[G\033[K")
	fmt.Print(i.fixedText + i.currentText)
}

func (i *Input) GetCurrentText() string {
	return i.currentText
}
