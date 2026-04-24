package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

// Необходимо реализовать собственный простейший Unix shell.
// Встроенные команды:
//– cd <path> – смена текущей директории.
//– pwd – вывод текущей директории.
//– echo <args> – вывод аргументов.
//– kill <pid> – послать сигнал завершения процессу с заданным PID.
//– ps – вывести список запущенных процессов.

// Запуск внешних команд через exec (с помощью системных вызовов fork/exec либо стандартных функций os/exec).
// Конвейеры (pipelines): возможность объединять команды через |, чтобы вывод одной команды направлять на ввод следующей
//(как в обычном shell).
// Например: ps | grep myprocess | wc -l.

// Обработку завершения: при нажатии Ctrl+D (EOF) шелл должен завершаться; Ctrl+C — прерывание текущей запущенной
// команды, но без закрыватия самой shell.

var ErrExit = errors.New("exit")

func main() {
	// Игнорируем Ctrl+C в родительском процессе (самом шелле)
	// Дочерние процессы будут получать его по умолчанию от терминала
	signal.Ignore(os.Interrupt)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("shell> ")
		if !scanner.Scan() {
			fmt.Println()
			break
		}

		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 1. Подстановка переменных окружения
		line = os.ExpandEnv(line)

		// 2. Обработка команд с учетом логических операторов && и ||
		err := executeLogical(line)
		if errors.Is(err, ErrExit) {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}
}

// executeLogical рекурсивно парсит && и ||
func executeLogical(line string) error {
	idxAnd := strings.Index(line, "&&")
	idxOr := strings.Index(line, "||")

	// Если нашли && и он стоит раньше (или || вообще нет)
	if idxAnd != -1 && (idxOr == -1 || idxAnd < idxOr) {
		left := line[:idxAnd]
		right := line[idxAnd+2:]
		err := executePipeline(left)
		if err == nil { // && выполняет правую часть, только если левая успешна
			return executeLogical(right)
		}
		return err
	}

	// Если нашли ||
	if idxOr != -1 {
		left := line[:idxOr]
		right := line[idxOr+2:]
		err := executePipeline(left)
		if err != nil { // || выполняет правую часть, только если левая упала
			return executeLogical(right)
		}
		return nil
	}

	// Если логических операторов больше нет, запускаем конвейер
	return executePipeline(line)
}

// executePipeline обрабатывает конвейеры (|) и редиректы (>, <)
func executePipeline(line string) error {
	rawCmds := strings.Split(line, "|")
	var cmds []*exec.Cmd

	var lastReader io.ReadCloser // Читающая сторона предыдущего пайпа

	for i, cmdStr := range rawCmds {
		args, inFile, outFile := parseRedirections(cmdStr)
		if len(args) == 0 {
			continue
		}

		// По умолчанию потоки смотрят в стандартные потоки терминала
		var stdin io.Reader = os.Stdin
		var stdout io.Writer = os.Stdout

		// Если это не первая команда в пайплайне, читаем из предыдущего пайпа
		if lastReader != nil {
			stdin = lastReader
		}

		// Если есть редирект на чтение (<)
		if inFile != "" {
			f, err := os.Open(inFile)
			if err != nil {
				return err
			}
			defer f.Close()
			stdin = f
		}

		var nextReader io.ReadCloser
		// Если это не последняя команда, создаем новый пайп для вывода
		if i < len(rawCmds)-1 {
			pr, pw, err := os.Pipe()
			if err != nil {
				return err
			}
			stdout = pw
			nextReader = pr
		}

		// Если есть редирект на запись (>)
		if outFile != "" {
			f, err := os.Create(outFile)
			if err != nil {
				return err
			}
			defer f.Close()
			stdout = f
		}

		// Выполнение встроенной или внешней команды
		if isBuiltin(args[0]) {
			err := handleBuiltin(args, stdout)
			// Закрыть writer пайпа, если встроенная команда в него писала
			if stdout != os.Stdout {
				if closer, ok := stdout.(io.WriteCloser); ok {
					closer.Close()
				}
			}
			if err != nil {
				return err
			}
		} else {
			// Подготовка внешней команды
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdin = stdin
			cmd.Stdout = stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Start(); err != nil {
				return fmt.Errorf("%s: %w", args[0], err)
			}

			// закрыть пишущую сторону пайпа в родительском процессе
			// Иначе следующая команда будет бесконечно ждать EOF
			if stdout != os.Stdout {
				if closer, ok := stdout.(io.WriteCloser); ok {
					closer.Close()
				}
			}
			cmds = append(cmds, cmd)
		}

		lastReader = nextReader
	}

	// Ждем завершения всех внешних команд в пайплайне
	var lastErr error
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// parseRedirections извлекает файлы для < и > и возвращает чистые аргументы команды
func parseRedirections(cmdStr string) (args []string, inFile, outFile string) {
	rawArgs := strings.Fields(cmdStr)
	for i := 0; i < len(rawArgs); i++ {
		if rawArgs[i] == "<" && i+1 < len(rawArgs) {
			inFile = rawArgs[i+1]
			i++ // Пропускаем имя файла
		} else if rawArgs[i] == ">" && i+1 < len(rawArgs) {
			outFile = rawArgs[i+1]
			i++ // Пропускаем имя файла
		} else {
			args = append(args, rawArgs[i])
		}
	}
	return args, inFile, outFile
}

func isBuiltin(name string) bool {
	switch name {
	case "cd", "pwd", "echo", "kill", "ps", "exit":
		return true
	}
	return false
}

func handleBuiltin(args []string, out io.Writer) error {
	switch args[0] {
	case "cd":
		path := "/"
		if len(args) > 1 {
			path = args[1]
		}
		if err := os.Chdir(path); err != nil {
			return err
		}
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Fprintln(out, dir)
	case "echo":
		fmt.Fprintln(out, strings.Join(args[1:], " "))
	case "kill":
		if len(args) < 2 {
			return fmt.Errorf("usage: kill <pid>")
		}
		var pid int
		if _, err := fmt.Sscanf(args[1], "%d", &pid); err != nil {
			return fmt.Errorf("invalid pid: %s", args[1])
		}
		proc, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		return proc.Signal(syscall.SIGTERM)
	case "ps":
		cmd := exec.Command("ps")
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		return cmd.Run()
	case "exit":
		return ErrExit
	}
	return nil
}
