package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	listenAddr = flag.String("http", ":8080", "Http listen address.")
	cpus       = flag.Int("cpus", 2, "Number of CPUs. Use `nproc` on linux to find your number of cores.")
	STRING     = flag.String("string", "\b5Ὂg̀9! ℃ᾭG 👏$⌘語❉☹∳Жm f", "String to analyze")
)

const (
	getTemplate = "static/get.html"
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*cpus)

	if err := run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() error {

	url := "localhost" + *listenAddr
	browserIsRunning := startBrowser(url)

	err := runServer()

	// server runs, but browser could not be started
	if err == nil && browserIsRunning == false {
		log.Println("inspector string now listening on port", *listenAddr)
	}

	return err
}

// startBrowser tries to open the URL in a browser, and returns
// whether it succeed.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default: // linux
		args = []string{"xdg-open"}
	}

	cmd := exec.Command(args[0], append(args[1:], url)...)

	return cmd.Start() == nil
}

func inspectString(s string) string {

	var out string

	out += fmt.Sprintf("\t\t<table>\n")
	for index, c := range s {
		link := fmt.Sprintf("<td><a href=\""+getInfoPage(c)+"\"> %#U </a></td> ", c)
		out += fmt.Sprintf("\t\t\t<tr>%v <td>starts at byte position %v</td></tr>\n", link, index)
		out += fmt.Sprintf("\t\t\t<tr><td></td><td>is hex byte [% x] </td></tr>", getHexBytes(c))
		if unicode.IsControl(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is control code point</td></tr>")
		}
		if unicode.IsDigit(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is digit code point</td></tr>")
		}
		if unicode.IsGraphic(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is graphic code point</td></tr>")
		}
		if unicode.IsLetter(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is letter code point</td></tr>")
		}
		if unicode.IsLower(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is lower case code point</td></tr>")
		}
		if unicode.IsMark(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is mark code point</td></tr>")
		}
		if unicode.IsNumber(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is number code point</td></tr>")
		}
		if unicode.IsPrint(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is printable code point</td></tr>")
		}
		if !unicode.IsPrint(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is not printable code point</td></tr>")
		}
		if unicode.IsPunct(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is punct code point</td></tr>")
		}
		if unicode.IsSpace(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is space code point</td></tr>")
		}
		if unicode.IsSymbol(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is symbol code point</td></tr>")
		}
		if unicode.IsTitle(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is title case code point</td></tr>")
		}
		if unicode.IsUpper(c) {
			out += fmt.Sprintf("\n\t\t\t<tr><td></td><td>is upper case code point</td></tr>")
		}
		out += fmt.Sprintf("\n")
	}
	out += fmt.Sprintf("\t\t</table>\n")

	return out

}

func getHexBytes(r rune) []byte {
	buf := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(buf, r)
	return buf
}

func getInfoPage(c rune) string {
	codepoint := fmt.Sprintf("%U", c)
	codepoint = strings.TrimLeft(codepoint, "U+")

	return "http://www.fileformat.info/info/unicode/char/" + codepoint + "/index.htm"
}
