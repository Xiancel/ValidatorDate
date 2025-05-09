package main

//Валідатор данних
//валідує данні згідно запросу користувача
//Види валідації (пример):
//- валідація айпи адреси
//- перевірка на надійність пароля

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

// меню опцій для користувача
func optionMenu() {
	fmt.Println("1. Перевірка email-адреси")
	fmt.Println("2. Перевірка надійності пароля")
	fmt.Println("3. Перевірка телефонного номера")
	fmt.Println("4. Перевірка IP-адреси")
	fmt.Println("5. Перевірка URL-адреси")
	fmt.Println("6. Перевірка Дати")
	fmt.Println("0. Вихід ")
}

// дозволенні символи для Email
func allowedCharEmail(chars rune) bool {
	return (chars >= 'a' && chars <= 'z') ||
		(chars >= 'A' && chars <= 'Z') ||
		(chars >= '0' && chars <= '9') ||
		chars == '.' || chars == '-' || chars == '_'
}

// перевірка валідності Email
func validEmail(email string) error {
	//слайс для хранение причин ошибки
	var reasons []string

	//перевірка ємейла на наявність символа @
	if strings.Count(email, "@") == 0 {
		reasons = append(reasons, "- Нема символа @")
	}

	//перевірка на наявність біль ніж 1 символа @
	if strings.Count(email, "@") != 1 {
		reasons = append(reasons, "- Не повинно бути біль ніж 1 символ @")
	}

	//перевірка на наявність пробілів
	if strings.Count(email, " ") != 0 {
		reasons = append(reasons, "- Не повинно бути пробілів ")
	}

	//розділення ємейлу на локальну та доменну частину
	parts := strings.Split(email, "@")

	//перевірка на наявність символі в локаліній частині
	if len(parts[0]) == 0 {
		reasons = append(reasons, "- В локаліній частини немає символів")
	}

	//перевірка на наявність символі в доменній частині
	if len(parts[1]) == 0 {
		reasons = append(reasons, "- В доменній частини немає символів")
	}

	//розділення домену по .
	domen := strings.Split(parts[1], ".")

	//перевірка на наявність точки в доменні частині
	if len(domen) <= 1 {
		reasons = append(reasons, "- Немає крапки в доменній частині")
	}

	//удаление пробілів в доменні верхного рівня
	last := strings.TrimSpace(domen[len(domen)-1])

	//перевірка довжини доменна верхнього рівня
	if len(last) < 2 || len(last) > 6 {
		reasons = append(reasons, "- Домен верхнього рівня повинен мати довжину від 2 до 6 символів ")
	}

	//перевірка допустимих символів у локальній частині
	local := parts[0]
	for _, aChars := range local {
		if !allowedCharEmail(aChars) {
			reasons = append(reasons, "- Локаліна частина не повинна мати недопустимий символ \n допустимі символи: літери, цифри, крапки, дефіси, підкреслення")
		}
	}

	//виведення помилки, якщо є причини помилки
	if len(reasons) > 0 {
		return errors.New("Email невалідний! Причини: \n" + strings.Join(reasons, "\n"))
	}

	return nil
}

// перевірка надійності пароля
func validPass(password string) error {
	//слайс для хранение причин ошибки
	var reasons []string

	//змінні для перевірки типів символів
	var digits, upper, lower, special bool

	//перевірка довжини пароля
	if len(password) < 8 {
		reasons = append(reasons, "- Менш ніж 8 символів")
	}

	//перевірка на наявність пробілів
	if strings.Count(password, " ") != 0 {
		reasons = append(reasons, "- Присутній Пробіл")
	}

	//перевірка символів
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

	//перевірка допустимих спеціальниз символів
	if strings.ContainsAny(password, "!@#$%^&*()-_=+[]{}|;:,.<>/?'\"") {
		special = true
	}

	//валідація змінних: чисел, великих літер, маленьких літер, спецсимволів
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

	//виведення помилки, якщо є причини помилки
	if len(reasons) > 0 {
		return errors.New("Пароль не надійний! Причини: \n" + strings.Join(reasons, "\n"))
	}

	return nil
}

// перевірка валідності Телефоного номеру
func validPhone(phonNum string) error {
	//слайс для хранение причин ошибки
	var reasons []string

	//слайс для телефоних операторів
	UAoperators := []string{"039", "050", "063", "066", "067", "068", "091", "092", "093", "094", "095", "096", "097", "098", "099"}

	//змінні для перевірки на цифри і операторів
	var onlydigits, validOp bool

	//перевірка початку номера з +
	if !strings.HasPrefix(phonNum, "+") {
		reasons = append(reasons, "- Повино щоб Номер Починався з +")
	}

	//перевірка на один плюс (на початку)
	if strings.Count(phonNum, "+") > 1 {
		reasons = append(reasons, "- Номер повинен містити + тільки з початку")
	}

	//перевірка на недопістимі символи
	phonNum = strings.TrimSpace(phonNum)
	for _, char := range phonNum {
		if !strings.ContainsAny(string(char), "1234567890 -+()") {
			reasons = append(reasons, "- Номер містить недопустимі символи")
			break
		}
	}

	//видалення все крім цифр
	lenphone := regexp.MustCompile(`\D`)
	onlyDPhone := lenphone.ReplaceAllString(phonNum, "")

	//перевірка довжини номера
	if len(onlyDPhone) > 15 {
		reasons = append(reasons, "- У номері більше символів чим треба")
	}
	if len(onlyDPhone) < 10 {
		reasons = append(reasons, "- У номері менше символів чим треба")
	}

	//перевірка на тільки цифри у номері
	for _, char := range onlyDPhone {
		if unicode.IsDigit(char) {
			onlydigits = true
		}
	}

	//перевірка коду оператора
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

	//валідація
	if !onlydigits {
		reasons = append(reasons, "- В номері не повино бути літер")
	}
	if !validOp {
		reasons = append(reasons, "- Номер повинен бути з правильним оператором ")
	}

	//виведення помилки, якщо є причини помилки
	if len(reasons) > 0 {
		return errors.New("Телефон не Валідний! Причини: \n" + strings.Join(reasons, "\n"))
	}

	return nil
}

// перевірка на валідність Ip-адреси
func validIp(ip string) error {
	//слайс для хранение причин ошибки
	var reasons []string

	//розділення айпи адреси по .
	parts := strings.Split(ip, ".")

	//перевірка на 4 частини в айпи адресі
	if len(parts) != 4 {
		reasons = append(reasons, "- Неверный формат ввода повино бути 4 части")
	}

	//перевірка на порожність частини, на пробіли,неправиліні символи, діапазону числ, на вмістимість 0 якщо це не 0.0.0.0
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

	//очищення пробілів і перевірка на довжину
	ip = strings.TrimSpace(ip)
	if len(ip) < 7 || len(ip) > 15 {
		reasons = append(reasons, "- Довжина має бути від 7 до 15 символів ")
	}

	//виведення помилки, якщо є причини помилки
	if len(reasons) > 0 {
		return errors.New("Ip-Адреса не Валідна! Причини: \n" + strings.Join(reasons, "\n"))
	}

	return nil
}

// допустимі символи в URl-адресі
func allowedCharUrl(chars rune) bool {
	return (chars >= 'a' && chars <= 'z') ||
		(chars >= 'A' && chars <= 'Z') ||
		(chars >= '0' && chars <= '9') ||
		chars == '.' || chars == '-'
}

// перевірка на валідність Url-адреси
func validUrl(url string) error {
	//слайс для хранение причин ошибки
	var reasons []string

	//перевірка на наявність протокола
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		reasons = append(reasons, "- Немає протокола http:// або https:// ")
	}

	//перевірка на пробіли
	if strings.Count(url, " ") != 0 {
		reasons = append(reasons, "- Не повинно бути пробілів у домені")
	}

	//видалення протоколу для подальщої перевірки
	var remPref string
	if strings.HasPrefix(url, "http://") {
		remPref = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		remPref = strings.TrimPrefix(url, "https://")
	}

	//перевірка на наличие точки в доменній частині
	if !strings.Contains(remPref, ".") {
		reasons = append(reasons, "- В доменній частині немає точки ")
	}

	//перевірка на недопустимі символи
	remPref = strings.TrimSpace(remPref)
	for _, chars := range remPref {
		if !allowedCharUrl(chars) {
			reasons = append(reasons, "- Доменне Ім'я має неправильний формат \n правильний формат: літери, цифри, дефіси")
			break
		}
	}

	//виведення помилки, якщо є причини помилки
	if len(reasons) > 0 {
		return errors.New("Url-адреса не Валідна! Причини: \n" + strings.Join(reasons, "\n"))
	}

	return nil
}

// перевірка на високостний рік
func leapYear(year int) bool {
	return (year%4 == 0 && year&100 != 0) || year%400 == 0
}

// перевірка на валідність Url-адреси
func validDate(date string) error {
	//слайс для хранение причин ошибки
	var reasons []string

	//перевірка на пробіли
	if strings.Count(date, " ") != 0 {
		reasons = append(reasons, "- Не повинно бути пробілів у даті")
	}

	//перевірка на недопустимі символи
	date = strings.TrimSpace(date)
	for _, char := range date {
		if !strings.ContainsAny(string(char), "1234567890-./") {
			reasons = append(reasons, "- Дата містить недопустимі символи \n допустимі символи: - . /")
			break
		}
	}

	//перевірка який стоить роздільник
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

	//розділення дати
	parts := strings.Split(date, sep)

	//перевірка на частини повино бути 3
	if len(parts) != 3 {
		reasons = append(reasons, "- Неправельний формат вводу повино бути 3 части")
	}

	//ініціалізація частин як день, місяць і рік
	day, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	month, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	year, _ := strconv.Atoi(strings.TrimSpace(parts[len(parts)-1]))

	//перевірка місяця
	if month < 1 || month > 12 {
		reasons = append(reasons, "- Невірний місяці")
	}

	//створення змінной для визначення кількості днів в місяці
	var dayInMon int

	//визначення кількості днів від місяця та року
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

	//валідація місяця
	if day < 1 || day > dayInMon {
		reasons = append(reasons, "-  Невірний день Місяця")
	}

	//валідація року
	if year < 1900 || year > 9999 {
		reasons = append(reasons, "-  Невірний рік")
	}

	//виведення помилки, якщо є причини помилки
	if len(reasons) > 0 {
		return errors.New("Дата не Валідна! Причини: \n" + strings.Join(reasons, "\n"))
	}

	return nil
}

// головна функція
func main() {
	//створення змінной вибір
	var choise int

	//створення зміних email, password, phoneNum, ip, url, date
	var email, password, phoneNum, ip, url, date string

	//створення читача для консолі
	reader := bufio.NewReader(os.Stdin)

	//цикл для программи
	for {
		fmt.Println("\nВиберіть опцію:")
		optionMenu()

		fmt.Print("\nВаш вибір: ")
		fmt.Scanln(&choise)

		//залежно від вибору користувача вибераеться підходящий кейс
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
