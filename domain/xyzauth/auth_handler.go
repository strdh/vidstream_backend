package xyzauth

import (
    "os"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "xyzstream/utils"
    "golang.org/x/crypto/bcrypt"
)

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

    isValid, messages := ValidateRegister(request)
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

    isValid, messages := ValidateLogin(request)
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