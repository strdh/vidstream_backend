package xyzauth

import (
    "log"
    "xyzstream/config"
)

type User struct{
    Id int `json:"id"`
    Email string `json:"email"`
    Username string `json:"username"`
    Password string `json:"password"`
    Created int `json:"created"`
}

func Create(user User) (int64, error) {
    result, err := config.DB.Exec("INSERT INTO users (email, username, password, created) VALUES (?, ?, ?, ?)", user.Email, user.Username, user.Password, user.Created)
    if err != nil {
        log.Println(err)
        return 0, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        log.Println(err)
        return 0, err
    }

    return id, nil
}

func ByUsername(username string) (User, error) {
    var user User
    row := config.DB.QueryRow("SELECT * FROM users WHERE username = ?", username)
    err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.Created)
    if err != nil {
        log.Println(err)
        return user, err
    }

    return user, nil
}

func CheckUsername(username string) (bool, error) {
    var count int
    err := config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
    if err != nil {
        log.Println(err)
        return false, err
    }

    return count > 0, nil
}

func CheckEmail(email string) (bool, error) {
    var count int
    err := config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
    if err != nil {
        log.Println(err)
        return false, err
    }

    return count > 0, nil
}