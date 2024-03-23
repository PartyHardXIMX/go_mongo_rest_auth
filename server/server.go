package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func hashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

var client *mongo.Client

func connectToDatabase() (*mongo.Client, error) {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return nil, err
    }
    return client, nil
}

func main() {
    var err error
    client, err = connectToDatabase()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(context.Background())

    // Маршруты
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/login", loginHandler)

    // Запуск сервера
    fmt.Println("Сервер запущен на http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Сервер работает")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Хешируем пароль
    hashedPassword, err := hashPassword(user.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Заменяем оригинальный пароль хешем
    user.Password = hashedPassword

    collection := client.Database("auth").Collection("users")

    // Проверяем, существует ли пользователь с таким именем
    existingUser := collection.FindOne(context.Background(), bson.M{"username": user.Username})
    if existingUser.Err() == nil {
        // Пользователь существует, возвращаем ошибку
        http.Error(w, "Такой пользователь уже существует в системе", http.StatusConflict)
        return
    }

    // Пользователь не существует, регистрируем его
    _, err = collection.InsertOne(context.Background(), user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    collection := client.Database("auth").Collection("users")

    // Получаем пользователя из базы данных по имени пользователя
    result := collection.FindOne(context.Background(), bson.M{"username": user.Username})
    if result.Err() != nil {
        http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
        return
    }

    // Декодируем хеш пароля из базы данных
    var storedUser User
    if err := result.Decode(&storedUser); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Сравниваем хешированный пароль с хешем из базы данных
    err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
    if err != nil {
        http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
        return
    }

    // Если пароль совпадает, возвращаем успешный статус
    w.WriteHeader(http.StatusOK)
}
