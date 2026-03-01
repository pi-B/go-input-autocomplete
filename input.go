package input_autocomplete

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/term"
)

type Input struct {
	cursor      *Cursor
	fixedText   string
	currentText string
	isCycling   bool
	cyclingPos  int
	matches     []string
}

func NewInput(fixedText string) *Input {
	return &Input{
		cursor:      NewCursor(),
		fixedText:   fixedText,
		currentText: "",
		isCycling:   false,
		cyclingPos:  0,
		matches:     []string{},
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

	i.PrintAllMatches()
	// i.currentText = i.matches[i.cyclingPos]
	// i.cyclingPos = (i.cyclingPos + 1) % len(i.matches)
	// i.cursor.SetPosition(len(i.currentText))
	// i.Print()
}

func (i *Input) PrintAllMatches() {
	var max_len int

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Impossible de récupérer la taille du terminal :", err)
		return
	}

	for _, match := range i.matches {
		if len(match) > max_len {
			max_len = len(match)
		}
	}

	display_slice := make([]string, max_len)
	res := ""
	fmt.Print("\033[G\033[K")

	for i, match := range i.matches {
		display_slice[i] = match
		for j := 0; j < max_len-len(match); j++ {
			display_slice[i] += " "
		}
		if len(res)+len(display_slice[i]) > width {
			fmt.Println(res)
			res = ""
		} else {
			res += display_slice[i]
		}
	}
	fmt.Println(res)
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
