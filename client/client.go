package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	for {
		var choice int
		fmt.Println("Выберите действие:")
		fmt.Println("1. Регистрация")
		fmt.Println("2. Авторизация")
		fmt.Println("3. Выход")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			register()
		case 2:
			login()
		case 3:
			fmt.Println("Выход из программы.")
			return
		default:
			fmt.Println("Неверный выбор.")
		}
	}
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func register() {
	var username, password string

	fmt.Println("Введите имя пользователя для регистрации:")
	fmt.Scanln(&username)
	fmt.Println("Введите пароль для регистрации:")
	fmt.Scanln(&password)

	registerUser(User{Username: username, Password: password})
}

func login() {
	var username, password string

	fmt.Println("Введите имя пользователя для входа:")
	fmt.Scanln(&username)
	fmt.Println("Введите пароль для входа:")
	fmt.Scanln(&password)

	loginUser(User{Username: username, Password: password})
}

func registerUser(user User) {
	jsonValue, _ := json.Marshal(user)
	response, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("Ошибка при регистрации пользователя:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusCreated {
		fmt.Println("Пользователь успешно зарегистрирован")
	} else if response.StatusCode == http.StatusConflict {
		fmt.Println("Ошибка: Такой пользователь уже существует")
	} else {
		fmt.Println("Ошибка при регистрации пользователя")
	}
}

func loginUser(user User) {
	jsonValue, _ := json.Marshal(user)
	response, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("Ошибка при входе пользователя:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Пользователь успешно вошел в систему")
	} else {
		fmt.Println("Ошибка: Неверные учетные данные")
	}
}
