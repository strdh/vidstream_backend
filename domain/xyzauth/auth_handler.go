package xyzauth

import (
    "os"
    "log"
    "strings"
    "unicode"
    "regexp"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "xyzstream/utils"
    "golang.org/x/crypto/bcrypt"
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

func validateRegister(request RegisterReq) (bool, map[string]string) {
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

func validateLogin(request LoginReq) (bool, map[string]string) {
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

func hashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Println(err)
        return "", err
    }

    return string(hash), nil
}

func comparePassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

func Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        utils.WriteResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }
    defer r.Body.Close()

    var request RegisterReq
    err = json.Unmarshal(body, &request)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusBadRequest, "Invalid request body", nil)
        return
    }

    isValid, messages := validateRegister(request)
    if !isValid {
        utils.WriteResponse(w, r, http.StatusBadRequest, "invalid request", messages)
        return
    }

    hashedPassword, err := hashPassword(request.Password)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }

    user := User{
        Username: request.Username,
        Email: request.Email,
        Password: hashedPassword,
    }

    id, err := Create(user)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }

    utils.WriteResponse(w, r, http.StatusOK, "User created successfully", id)
}

func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        utils.WriteResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }

    var request LoginReq
    err = json.Unmarshal(body, &request)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusBadRequest, "Invalid request body", nil)
        return
    }

    isValid, messages := validateLogin(request)
    if !isValid {
        utils.WriteResponse(w, r, http.StatusBadRequest, "invalid request", messages)
        return
    }

    user, err := ByUsername(request.Username)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusUnauthorized, "Username or password invalid", nil)
        return
    }

    if !comparePassword(user.Password, request.Password) {
        utils.WriteResponse(w, r, http.StatusUnauthorized, "Username or password invalid", nil)
        return
    }

    jwtToken := utils.GenerateToken(user.Id, os.Getenv("JWT_KEY"))

    loginResponse := loginResp{
        Token: jwtToken,
    }

    utils.WriteResponse(w, r, http.StatusOK, "Login successful", loginResponse)
}