// Package пакет для работы с консолью
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

type CLIer interface {
	// Register регистрация функции по имени
	Register(string, func())
	// GetUserPrint запросить ввод текста от пользователя
	GetUserPrint(string) (string, error)
	// GetHideUserPrint запросить скрытый ввод текста от пользователя
	GetHideUserPrint(string) (string, error)
	// Println вывести в консоль текст
	Println(string)
	// CallFn вызов функции вручную
	CallFn(string)
	// Run запуск цикла опроса команд
	Run()
}

type CLI struct {
	reader *bufio.Reader
	list   map[string]func()
}

// New получить экземпляр структуры для работы с консолью
func New() *CLI {
	return &CLI{
		reader: bufio.NewReader(os.Stdin),
		list:   map[string]func(){},
	}
}

// Register регистрация функции по имени
func (cli *CLI) Register(name string, fn func()) {
	cli.list[name] = fn
}

// CallFn вызов функции вручную
func (cli *CLI) CallFn(command string) {
	if _, ok := cli.list[command]; ok {
		cli.list[command]()
	}
}

// Run запуск цикла опроса команд
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

// DrawMenu вывести в консоль описание и команды
func (cli *CLI) DrawMenu() {
	os.Stdout.WriteString("Добро пожаловать в консольную программу GophKeeper\n")
	os.Stdout.WriteString("Доступные команды:\n")
	os.Stdout.WriteString("- help - Вывести доступные комманды\n")
	os.Stdout.WriteString("- exit - Выход из программы\n")
	os.Stdout.WriteString("- registration - Регистрация\n")
	os.Stdout.WriteString("- login - Авторизация\n")
	os.Stdout.WriteString("- sync - Синхронизация данных с сервером\n")
	os.Stdout.WriteString("- list - Вывести все данных\n")
	os.Stdout.WriteString("- add - Добавить данные\n")
}

// Println вывести в консоль текст
func (cli *CLI) Println(message string) {
	os.Stdout.WriteString(fmt.Sprintf("%s\n", message))
}

// GetUserPrint запросить ввод текста от пользователя
func (cli *CLI) GetUserPrint(title string) (string, error) {
	os.Stdout.WriteString(title)
	text, err := cli.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(strings.TrimSpace(text), "\n"), nil
}

// GetHideUserPrint запросить скрытый ввод текста от пользователя
func (cli *CLI) GetHideUserPrint(title string) (string, error) {
	os.Stdout.WriteString(title)
	text, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	cli.Println("")
	return strings.TrimSuffix(strings.TrimSpace(string(text)), "\n"), nil
}
