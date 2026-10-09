package main

import (
	"bufio"
	"fmt"
	"homework/colors"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Введите help, чтобы увидеть список доступных команд")
	scanner := bufio.NewScanner(os.Stdin)
	mainPath, _ := os.Getwd()
	currPath := mainPath

	for {
		fmt.Print(colors.Bold + currPath + "> " + colors.Reset)
		scanner.Scan()
		if strings.TrimSpace(scanner.Text()) == "" {
			printError("Введите команду")
			continue
		}
		fields := strings.Fields(scanner.Text())
		cmd := strings.ToLower(fields[0])

		switch cmd {
		case "help":
			showHelp()

		case "cd":
			if len(fields) != 2 {
				printError("Неверно введена команда")
				fmt.Println(colors.Bold + colors.Green + "Правильная команда: " + colors.Reset + colors.Italic + colors.Bold + "cd <папка>" + colors.Reset)
				continue
			}
			currPath = moveToFolder(fields[1], currPath)

		case "create", "delete":
			if len(fields) != 3 {
				printError("Неверное использование команды")
				continue
			}
			contentType := strings.ToLower(fields[1])
			contentName := fields[2]
			if cmd == "create" {
				switch contentType {
				case "file":
					createFile(contentName, currPath)
				case "folder":
					createFolder(contentName, currPath)

				default:
					printError("Неверное указан тип")
				}
			} else {
				if strings.TrimSpace(contentName) == "" || checkInvalidSymbol(contentName) {
					printError("Недопустимое имя")
					continue
				}
				path := filepath.Join(currPath, contentName)
				_, err := os.Stat(path)
				if err != nil {
					printError("Файл или папка не найдены")
					continue
				}
				switch contentType {
				case "file":
					deleteFile(path)
				case "folder":
					fmt.Println("Вы действительно хотите удалить папку " + contentName + " со всем содержимым?")
					fmt.Println("Введите " + colors.Bold + "yes" + colors.Reset + " или " + colors.Bold + "no" + colors.Reset)
					scanner.Scan()
					isDelete := strings.TrimSpace(scanner.Text())
					switch isDelete {
					case "yes", "y":
						deleteFolder(path)
					case "no", "n":
						continue
					default:
						fmt.Println("Введите " + colors.Bold + "yes" + colors.Reset + " или " + colors.Bold + "no" + colors.Reset)
					}

				}
			}

		case "cd..":
			currPath = backToFolder(currPath)

		case "show":
			showContent(currPath)

		case "rename":
			if len(fields) != 2 {
				printError("Неправильное использование команды")
				continue
			}
			oldName := fields[1]
			filePath := filepath.Join(currPath, oldName)
			_, err := os.Stat(filePath)
			if err != nil {
				printError("Файл или папка не найдены")
				continue
			}
			fmt.Println("Введите новое название")
			scanner.Scan()
			newName := scanner.Text()
			if strings.TrimSpace(newName) == "" || checkInvalidSymbol(newName) {
				printError("Недопустимое имя")
				continue
			}
			renameContent(filePath, newName)

		case "exit":
			return
		default:
			printError("Неверно введена команда")
			fmt.Println("Чтобы посмотреть список доступных команд введите " + colors.Bold + "help" + colors.Reset)

		}
	}
}

/********* ФУНКЦИИ *********/

func createFile(fileName string, currPath string) {
	if checkInvalidSymbol(fileName) {
		printError("Название содержит недопустимые символы")
	} else {
		_, error := os.Stat(filepath.Join(currPath, fileName))
		if os.IsNotExist(error) {
			file, err := os.Create(filepath.Join(currPath, fileName))
			if err != nil {
				printError("Ошибка создания файла")
				return
			}
			defer file.Close()
			fmt.Println("Файл успешно создан")
		} else {
			printError("Файл с таким названием уже существует")
		}
	}

}

func createFolder(folderName string, currPath string) {
	if checkInvalidSymbol(folderName) {
		printError("Название содержит недопустимые символы")
	} else {
		_, error := os.Stat(filepath.Join(currPath, folderName))
		if os.IsNotExist(error) {
			err := os.Mkdir(filepath.Join(currPath, folderName), 0777)
			if err != nil {
				printError("Ошибка создания папки")
			} else {
				fmt.Println("Папка успешно создана")
			}
		} else {
			printError("Папка с таким названием уже существует")

		}

	}

}

func deleteFolder(folderAdress string) {
	err := os.RemoveAll(folderAdress)
	if err != nil {
		printError("Ошибка удаления")
		return
	}
	fmt.Println("Успешное удаление")
}

func deleteFile(fileAdress string) {
	err := os.Remove(fileAdress)
	if err != nil {
		printError("Ошибка удаления")
		return
	}
	fmt.Println("Успешное удаление")
}

func moveToFolder(folderName string, currPath string) string {
	file, err := os.Stat(filepath.Join(currPath, folderName))
	if err != nil || !file.IsDir() {
		printError("Папки с таким названием не существует")
	} else {
		currPath = filepath.Join(currPath, folderName)
	}
	return currPath
}

func backToFolder(currPath string) string {
	return filepath.Dir(currPath)
}

func showContent(currPath string) {
	files, err := os.ReadDir(currPath)
	if err != nil {
		printError("Не удалось прочитать папку")
		return
	}
	if len(files) == 0 {
		fmt.Println(colors.Bold + colors.Cyan + "Папка пустая" + colors.Reset)
	}
	for _, file := range files {
		filePath := filepath.Join(currPath, file.Name())
		fileInfo, err := os.Lstat(filePath)
		if err != nil {
			fmt.Println("Error", err)
			return
		} else {
			if fileInfo.IsDir() {
				fmt.Println(colors.Yellow + colors.Bold + file.Name() + "\\" + colors.Reset)
			} else {
				fmt.Println(colors.Bold + colors.White + file.Name() + colors.Reset)
			}
		}

	}
}

func renameContent(currPath string, newName string) {
	newPath := filepath.Join(backToFolder(currPath), strings.TrimSpace(newName))
	err := os.Rename(currPath, newPath)
	if err != nil {
		fmt.Println("Error", err)
	}
}

func showHelp() {
	fmt.Println(colors.Bold + "Список доступных команд" + colors.Reset)
	fmt.Println("1." + colors.Italic + colors.Bold + "cd <папка> " + colors.Reset + "- команда для перехода в папку внутри текущего хранилища")
	fmt.Println("2." + colors.Italic + colors.Bold + "cd.. " + colors.Reset + "- команда для возвращения к предыдущей папке")
	fmt.Println("3." + colors.Italic + colors.Bold + " create file <имя> " + colors.Reset + "/" + colors.Italic + colors.Bold + " create folder <имя> " + colors.Reset + "- команда для создания нового файла или папки в хранилище, в которым вы находитесь")
	fmt.Println("4." + colors.Italic + colors.Bold + " delete file <имя> " + colors.Reset + "/" + colors.Italic + colors.Bold + " delete folder <имя> " + colors.Reset + "- команда для удаления файла или папки в хранилище, в которым вы находитесь")
	fmt.Println("5." + colors.Italic + colors.Bold + " rename <имя> " + colors.Reset + "- команда для переименования файла или папки в хранилище, в которым вы находитесь")
	fmt.Println("6." + colors.Italic + colors.Bold + "show " + colors.Reset + "- команда для показа всех папок и файлов, находящихся в хранилище, в котором вы сейчас находитесь")
	fmt.Println("7." + colors.Italic + colors.Bold + "exit " + colors.Reset + "- команда для завершения программы")
}

func printError(err string) {
	fmt.Println(colors.Bold + colors.Red + err + colors.Reset)
}

func checkInvalidSymbol(contentName string) bool {
	invalid := "\"<>\\:/?*|"
	if strings.ContainsAny(contentName, invalid) {
		return true
	}
	return false
}
