package tui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"ya-gophkeeper-client/internal/auth/entity"

	"golang.org/x/term"
)

type TUI struct {
	reader *bufio.Reader
	list   map[string]func()
}

func New() *TUI {
	return &TUI{
		reader: bufio.NewReader(os.Stdin),
		list:   map[string]func(){},
	}
}

func (tui *TUI) Register(name string, fn func()) {
	tui.list[name] = fn
}

func (tui *TUI) Run() {
	for {
		fmt.Print("→ ")
		text, err := tui.reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		command := strings.TrimSuffix(strings.TrimSpace(text), "\n")
		if command == "help" {
			tui.DrawMenu()
		}
		if command == "exit" {
			return
		}
		if _, ok := tui.list[command]; ok {
			tui.list[command]()
		}
	}
}

func (tui *TUI) DrawMenu() {
	os.Stdout.WriteString("Добро пожаловать в консольную программу GophKeeper\n")
	os.Stdout.WriteString("Доступные команды:\n")
	os.Stdout.WriteString("- help - Вывести доступные комманды\n")
	os.Stdout.WriteString("- exit - Выход из программы\n")
	for k := range tui.list {
		os.Stdout.WriteString(fmt.Sprintf("%s\n", k))
	}
}

func (tui *TUI) Println(message string) {
	os.Stdout.WriteString(fmt.Sprintf("%s\n", message))
}

func (tui *TUI) GetUserPrint(title string) (string, error) {
	os.Stdout.WriteString(title)
	text, err := tui.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(strings.TrimSpace(text), "\n"), nil
}

func (tui *TUI) GetHideUserPrint(title string) (string, error) {
	os.Stdout.WriteString(title)
	text, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(strings.TrimSpace(string(text)), "\n"), nil
}

func (tui *TUI) GetEmailAndPassword() (*entity.LoginDto, error) {
	os.Stdout.WriteString("Email: ")
	email, err := tui.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	os.Stdout.WriteString("Password: ")
	pass, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	os.Stdout.WriteString("\n\n")
	return &entity.LoginDto{
		Email:    strings.TrimSpace(email),
		Password: strings.TrimSuffix(strings.TrimSpace(string(pass)), "\n"),
	}, nil
}
