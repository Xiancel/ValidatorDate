package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// допустимые символы для Email
func allowedChar(chars rune) bool {
	return (chars >= 'a' && chars <= 'z') ||
		(chars >= 'A' && chars <= 'Z') ||
		(chars >= '0' && chars <= '9') ||
		chars == '.' || chars == '-' || chars == '_'
}
func optionMenu() {
	fmt.Println("1. Перевірка email-адреси")
	fmt.Println("2. Перевірка надійності пароля")
	fmt.Println("3. Перевірка телефонного номера")
	fmt.Println("4. Перевірка IP-адреси")
	fmt.Println("5. Перевірка URL-адреси")
	fmt.Println("0. Вихід ")
}

func validEmail(email string) error {
	// створення помилки
	err := errors.New("Email не валідний!")

	//перевірка введеного email на наявність символа "@" чи пробілів
	if strings.Count(email, "@") != 1 || strings.Count(email, " ") != 0 {
		return err
	}
	//розділення ємейлу на дві частини
	parts := strings.Split(email, "@")

	//перевірка на наявність символів в локальній та доменній частинах
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return err
	}

	//розділення доменної частини
	domen := strings.Split(parts[1], ".")

	//перевірка доменной частини на наявність крапки
	if len(domen) <= 1 {
		return err
	}

	//останній єлемент в доменній частини пілся крапки
	last := domen[len(domen)-1]
	//перевірка на довжину останного елемента
	if len(last) < 2 || len(last) > 6 {
		return err
	}
	//перевірка на допустимість символів для Локаліної частини
	local := parts[0]
	for _, aChars := range local {
		if !allowedChar(aChars) {
			return err
		}
	}

	return nil
}

func validPass(password string) error {
	var reasons []string

	var digits, upper, lower, special bool
	if len(password) < 8 {
		reasons = append(reasons, "- Менш ніж 8 символів")
	}
	if strings.Count(password, " ") != 0 {
		reasons = append(reasons, "- Присутній Пробіл")
	}

	for _, r := range password {
		if unicode.IsDigit(r) {
			digits = true
		}
		if unicode.IsUpper(r) {
			upper = true
		}
		if unicode.IsLower(r) {
			lower = true
		}
	}

	if strings.ContainsAny(password, "!@#$%^&*()-_=+[]{}|;:,.<>/?'\"") {
		special = true
	}

	if !digits {
		reasons = append(reasons, "- Відсутні Цифри")
	}
	if !upper {
		reasons = append(reasons, "- Відсутні великі літери")
	}
	if !lower {
		reasons = append(reasons, "- Відсутні маленьки літери")
	}
	if !special {
		reasons = append(reasons, "- Відсутні спеціальні символи")
	}

	if len(reasons) > 0 {
		return errors.New("Пароль не надійний! Причини: \n" + strings.Join(reasons, "\n"))
	}

	return nil
}
func main() {
	var choise int
	var email, password string
	for {
		fmt.Println("\nВиберіть опцію:")
		optionMenu()

		fmt.Print("\nВаш вибір: ")
		fmt.Scanln(&choise)

		switch choise {
		case 0:
			return
		case 1:
			fmt.Print("Введіть email-адресу: ")
			fmt.Scanln(&email)
			err := validEmail(email)
			if err != nil {
				fmt.Println("Результат: ", err)
			} else {
				fmt.Println("Результат: Email валідний!")
			}
		case 2:
			fmt.Print("Введіть пароль: ")
			fmt.Scanln(&password)
			err := validPass(password)
			if err != nil {
				fmt.Println("Результат: ", err)
			} else {
				fmt.Println("Результат: Пароль надійний!")
			}
		case 3:
			fmt.Println("Ваш вибір: ", choise)
		case 4:
			fmt.Println("Ваш вибір: ", choise)
		case 5:
			fmt.Println("Ваш вибір: ", choise)
		default:
			fmt.Println("Такого Виору не існує")
		}
	}
}
