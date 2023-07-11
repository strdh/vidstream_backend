package utils

import (
    "net/http"
    "encoding/json"
)

type Response struct {
    Status int `json:"status"`
    Message string `json:"message"`
    Data interface{} `json:"data"`
}

func WriteResponse(w http.ResponseWriter, r *http.Request, status int, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    
    response := Response{
        Status: status,
        Message: message,
        Data: data,
    }

    jsonData, err := json.Marshal(response)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(status)
    w.Write(jsonData)
}