package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type CLI struct {
	reader *bufio.Reader
	list   map[string]func()
}

func New() *CLI {
	return &CLI{
		reader: bufio.NewReader(os.Stdin),
		list:   map[string]func(){},
	}
}

func (cli *CLI) Register(name string, fn func()) {
	cli.list[name] = fn
}

func (cli *CLI) Call(command string) {
	if _, ok := cli.list[command]; ok {
		cli.list[command]()
	}
}

func (cli *CLI) Run() {
	for {
		fmt.Print("→ ")
		text, err := cli.reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		command := strings.TrimSuffix(strings.TrimSpace(text), "\n")
		if command == "help" {
			cli.DrawMenu()
		}
		if command == "exit" {
			return
		}
		if _, ok := cli.list[command]; ok {
			cli.list[command]()
		}
	}
}

func (cli *CLI) DrawMenu() {
	os.Stdout.WriteString("Добро пожаловать в консольную программу GophKeeper\n")
	os.Stdout.WriteString("Доступные команды:\n")
	os.Stdout.WriteString("- help - Вывести доступные комманды\n")
	os.Stdout.WriteString("- exit - Выход из программы\n")
	for k := range cli.list {
		os.Stdout.WriteString(fmt.Sprintf("%s\n", k))
	}
}

func (cli *CLI) Println(message string) {
	os.Stdout.WriteString(fmt.Sprintf("%s\n", message))
}

func (cli *CLI) GetUserPrint(title string) (string, error) {
	os.Stdout.WriteString(title)
	text, err := cli.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(strings.TrimSpace(text), "\n"), nil
}

func (cli *CLI) GetHideUserPrint(title string) (string, error) {
	os.Stdout.WriteString(title)
	text, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	cli.Println("")
	return strings.TrimSuffix(strings.TrimSpace(string(text)), "\n"), nil
}
