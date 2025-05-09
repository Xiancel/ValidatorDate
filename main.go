package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func allowedCharEmail(chars rune) bool {
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
	fmt.Println("6. Перевірка Дати")
	fmt.Println("0. Вихід ")
}

func validEmail(email string) error {

	var reasons []string

	if strings.Count(email, "@") == 0 {
		reasons = append(reasons, "- Нема символа @")
	}
	if strings.Count(email, "@") != 1 {
		reasons = append(reasons, "- Не повинно бути біль ніж 1 символ @")
	}
	if strings.Count(email, " ") != 0 {
		reasons = append(reasons, "- Не повинно бути пробілів ")
	}

	parts := strings.Split(email, "@")

	if len(parts[0]) == 0 {
		reasons = append(reasons, "- В локаліній частини немає символів")
	}
	if len(parts[1]) == 0 {
		reasons = append(reasons, "- В доменній частини немає символів")
	}

	domen := strings.Split(parts[1], ".")

	if len(domen) <= 1 {
		reasons = append(reasons, "- Немає крапки в доменній частині")
	}

	last := strings.TrimSpace(domen[len(domen)-1])

	if len(last) < 2 || len(last) > 6 {
		reasons = append(reasons, "- Домен верхнього рівня повинен мати довжину від 2 до 6 символів ")
	}

	local := parts[0]
	for _, aChars := range local {
		if !allowedCharEmail(aChars) {
			reasons = append(reasons, "- Доменна частина не повинна мати недопустимий символ \n допустимі символи: літери, цифри, крапки, дефіси, підкреслення")
		}
	}

	if len(reasons) > 0 {
		return errors.New("Email невалідний! Причини: \n" + strings.Join(reasons, "\n"))
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

func validPhone(phonNum string) error {
	var reasons []string
	UAoperators := []string{"039", "050", "063", "066", "067", "068", "091", "092", "093", "094", "095", "096", "097", "098", "099"}
	var onlydigits, validOp bool
	if !strings.HasPrefix(phonNum, "+") {
		reasons = append(reasons, "- Повино щоб Номер Починався з +")
	}
	if strings.Count(phonNum, "+") > 1 {
		reasons = append(reasons, "- Номер повинен містити + тільки з початку")
	}

	phonNum = strings.TrimSpace(phonNum)
	for _, char := range phonNum {
		if !strings.ContainsAny(string(char), "1234567890 -+()") {
			reasons = append(reasons, "- Номер містить недопустимі символи")
			break
		}
	}

	lenphone := regexp.MustCompile(`\D`)
	onlyDPhone := lenphone.ReplaceAllString(phonNum, "")
	if len(onlyDPhone) > 15 {
		reasons = append(reasons, "- У номері більше символів чим треба")
	}
	if len(onlyDPhone) < 10 {
		reasons = append(reasons, "- У номері менше символів чим треба")
	}
	for _, char := range onlyDPhone {
		if unicode.IsDigit(char) {
			onlydigits = true
		}
	}

	if strings.HasPrefix(onlyDPhone, "380") {
		onlyDPhone = onlyDPhone[2:]
		operCode := onlyDPhone[:3]
		for _, code := range UAoperators {
			if operCode == code {
				validOp = true
				break
			}
		}
	}

	if !onlydigits {
		reasons = append(reasons, "- В номері не повино бути літер")
	}
	if !validOp {
		reasons = append(reasons, "- Номер повинен бути з правильним оператором ")
	}
	if len(reasons) > 0 {
		return errors.New("Телефон не Валідний! Причини: \n" + strings.Join(reasons, "\n"))
	}

	return nil
}

func validIp(ip string) error {
	var reasons []string
	parts := strings.Split(ip, ".")

	if len(parts) != 4 {
		reasons = append(reasons, "- Неверный формат ввода повино бути 4 части")
	}

	for i, part := range parts {
		if part == "" {
			reasons = append(reasons, fmt.Sprintf("- Частина %d порожня ", i+1))
		}

		if strings.Count(part, " ") != 0 {
			reasons = append(reasons, "- Ip-Адреса не повина мати пробілів")
		}

		part = strings.TrimSpace(part)

		num, _ := strconv.Atoi(part)
		if !strings.ContainsAny(part, "1234567890") {
			reasons = append(reasons, fmt.Sprintf("- Частина %d содержит неправильний символ ", i+1))
			continue
		}

		if num < 0 || num > 255 {
			reasons = append(reasons, fmt.Sprintf("- Частина %d поза діапазоном 0-255 ", i+1))
		}
		if len(part) > 1 && part[0] == '0' {
			reasons = append(reasons, fmt.Sprintf("- Частина %d містить 0 перед числом: ", i+1))
		}
	}
	ip = strings.TrimSpace(ip)
	if len(ip) < 7 || len(ip) > 15 {
		reasons = append(reasons, "- Довжина має бути від 7 до 15 символів ")
	}

	if len(reasons) > 0 {
		return errors.New("Ip-Адреса не Валідна! Причини: \n" + strings.Join(reasons, "\n"))
	}
	return nil
}

func allowedCharUrl(chars rune) bool {
	return (chars >= 'a' && chars <= 'z') ||
		(chars >= 'A' && chars <= 'Z') ||
		(chars >= '0' && chars <= '9') ||
		chars == '.' || chars == '-'
}
func validUrl(url string) error {
	var reasons []string

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		reasons = append(reasons, "- Немає протокола http:// або https:// ")
	}

	if strings.Count(url, " ") != 0 {
		reasons = append(reasons, "- Не повинно бути пробілів у домені")
	}

	var remPref string
	if strings.HasPrefix(url, "http://") {
		remPref = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		remPref = strings.TrimPrefix(url, "https://")
	}

	if !strings.Contains(remPref, ".") {
		reasons = append(reasons, "- В доменній частині немає точки ")
	}
	remPref = strings.TrimSpace(remPref)
	for _, chars := range remPref {
		if !allowedCharUrl(chars) {
			reasons = append(reasons, "- Доменне Ім'я має неправильний формат \n правильний формат: літери, цифри, дефіси")
			break
		}
	}

	if len(reasons) > 0 {
		return errors.New("Url-адреса не Валідна! Причини: \n" + strings.Join(reasons, "\n"))
	}
	return nil
}

func leapYear(year int) bool {
	return (year%4 == 0 && year&100 != 0) || year%400 == 0
}
func validDate(date string) error {
	var reasons []string

	if strings.Count(date, " ") != 0 {
		reasons = append(reasons, "- Не повинно бути пробілів у даті")
	}
	date = strings.TrimSpace(date)
	for _, char := range date {
		if !strings.ContainsAny(string(char), "1234567890-./") {
			reasons = append(reasons, "- Дата містить недопустимі символи \n допустимі символи: - . /")
			break
		}
	}
	var sep string
	if strings.Contains(date, "-") {
		sep = "-"
	} else if strings.Contains(date, "/") {
		sep = "/"
	} else if strings.Contains(date, ".") {
		sep = "."
	} else {
		reasons = append(reasons, "- Неправельний формат вводу")
	}

	parts := strings.Split(date, sep)

	if len(parts) != 3 {
		reasons = append(reasons, "- Неправельний формат вводу повино бути 3 части")
	}
	day, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	year, _ := strconv.Atoi(parts[2])

	if month < 1 || month > 12 {
		reasons = append(reasons, "- Невірний місяці")
	}

	var dayInMon int

	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		dayInMon = 31
	case 4, 6, 9, 11:
		dayInMon = 30
	case 2:
		if leapYear(year) {
			dayInMon = 29
		} else {
			dayInMon = 28
		}
	}

	if day < 1 || day > dayInMon {
		reasons = append(reasons, "-  Невірний день Місяця")
	}

	if year < 1900 || year > 9999 {
		reasons = append(reasons, "-  Невірний рік")
	}

	if len(reasons) > 0 {
		return errors.New("Дата не Валідна! Причини: \n" + strings.Join(reasons, "\n"))
	}
	return nil
}

func main() {
	var choise int
	var email, password, phoneNum, ip, url, date string

	reader := bufio.NewReader(os.Stdin)
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
			email, _ = reader.ReadString('\n')
			err := validEmail(email)
			if err != nil {
				fmt.Println("Результат: ", err)
			} else {
				fmt.Println("Результат: Email валідний!")
			}
		case 2:
			fmt.Print("Введіть пароль: ")
			password, _ = reader.ReadString('\n')
			err := validPass(password)
			if err != nil {
				fmt.Println("Результат: ", err)
			} else {
				fmt.Println("Результат: Пароль надійний!")
			}
		case 3:
			fmt.Print("Введіть Номер Телефону: ")
			phoneNum, _ = reader.ReadString('\n')
			err := validPhone(phoneNum)
			if err != nil {
				fmt.Println("Результат: ", err)
			} else {
				fmt.Println("Результат: Номер валідний!")
			}
		case 4:
			fmt.Print("Введіть Ip-Адресу: ")
			ip, _ = reader.ReadString('\n')
			err := validIp(ip)
			if err != nil {
				fmt.Println("Результат: ", err)
			} else {
				fmt.Println("Результат: Ip-Адреса Валідна!")
			}
		case 5:
			fmt.Print("Введіть Url-Адресу: ")
			url, _ = reader.ReadString('\n')
			err := validUrl(url)
			if err != nil {
				fmt.Println("Результат: ", err)
			} else {
				fmt.Println("Результат: Url-Адреса Валідна!")
			}
		case 6:
			fmt.Print("Введіть Дату: ")
			date, _ = reader.ReadString('\n')
			err := validDate(date)
			if err != nil {
				fmt.Println("Результат: ", err)
			} else {
				fmt.Println("Результат: Дата Валідна!")
			}
		default:
			fmt.Println("Такого Вибору не існує")
		}
	}
}
