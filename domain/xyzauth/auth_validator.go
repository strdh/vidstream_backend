package xyzauth

import (
    "unicode"
    "regexp"
    "strings"
)

type loginResp struct {
    Token string `json:"token"`
}

type RegisterReq struct {
    Email string `json:"email"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func ValidateRegister(request RegisterReq) (bool, map[string]string) {
    messages := map[string]string{}

    if request.Username == "" || len(request.Username) < 3 {
        messages["username"] = "Username is required and must be at least 3 characters"
    } else {
        // Validate username format using regular expression
        usernameRegex := regexp.MustCompile("^[a-z0-9_]+$")
        if !usernameRegex.MatchString(request.Username) {
            messages["username"] = "Username must contain only lowercase letters, numbers[0-9], and underscores"
        } else {
            usernameExists, _ := CheckUsername(request.Username)
            if usernameExists {
                messages["username"] = "Username is already taken"
            }
        }
    }

    if request.Email != "" {
        pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	    regex := regexp.MustCompile(pattern)
	    isValid := regex.MatchString(request.Email)

        if !isValid {
            messages["email"] = "Email is not valid"
        } else {
            emailExists, _ := CheckEmail(request.Email)
            if emailExists {
                messages["email"] = "Email is already taken"
            }
        }
    } else {
        messages["email"] = "Email is required"
    }

    if request.Password == "" || len(request.Password) < 6 {
        messages["password"] = "Password is required and must be at least 6 characters"
    } else {
        hasLower := false
        hasUpper := false
        hasSpecial := false
        hasNumber := false

        for _, char := range request.Password {
            if unicode.IsLower(char) {
                hasLower = true
            } else if unicode.IsUpper(char) {
                hasUpper = true
            } else if strings.ContainsAny(string(char), "!@#$%^&*()") {
                hasSpecial = true
            } else if unicode.IsNumber(char) {
                hasNumber = true
            }
        }

        if !hasLower || !hasUpper || !hasSpecial || !hasNumber {
            messages["password"] = "Password must contain at least one lowercase letter, uppercase letter, number, and special character"
        }
    }

    if len(messages) > 0 {
        return false, messages
    }

    return true, messages
}

func ValidateLogin(request LoginReq) (bool, map[string]string) {
    messages := map[string]string{}

    if request.Username == "" {
        messages["username"] = "Username is required"
    }

    if request.Password == "" {
        messages["password"] = "Password is required"
    }

    if len(messages) > 0 {
        return false, messages
    }

    return true, messages
}