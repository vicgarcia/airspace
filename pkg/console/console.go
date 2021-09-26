package console

import (
	"fmt"
	"regexp"
	"os"
	"os/exec"
)

const (
	ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
)

var re = regexp.MustCompile(ansi)

func SanitizeInput(str string) string {
	return re.ReplaceAllString(str, "")
}

func Render(output string) {
	fmt.Print(output)
}

func Renderln(output string) {
	fmt.Println(output)
}

func RenderPrompt() {
	fmt.Print("\n> ")
}

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
