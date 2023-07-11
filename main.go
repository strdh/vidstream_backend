package main

import (
    "log"
    "fmt"
    "os"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "xyzstream/config"
    "xyzstream/utils"
    "xyzstream/domain/xyzvod"
    "xyzstream/middleware"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    config.InitializeDB()
    config.InitializeS3()

    router := mux.NewRouter()
    router.Use(corsOptions)

    router.HandleFunc("/vodupload", middleware.JWTMiddleware(xyzvod.VodUpload)).Methods("POST", "OPTIONS")
    router.HandleFunc("/vod", xyzvod.VodList).Methods("GET", "OPTIONS")
    router.HandleFunc("/vod/next/{id}", xyzvod.VodListNext).Methods("GET", "OPTIONS")
    router.HandleFunc("/vod/{vodulid}", xyzvod.VodStream).Methods("GET", "OPTIONS")
    router.HandleFunc("/vod/stream/{segment}", xyzvod.VodSegment).Methods("GET", "OPTIONS")

    server := http.Server{
        Addr: os.Getenv("STREAM_ADDRESS"),
        Handler: router,
    }

    fmt.Println("Server running at: ", os.Getenv("STREAM_ADDRESS"))
    err = server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}

func corsOptions(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodOptions {
            utils.WriteResponse(w, r, http.StatusOK, "OK", nil)
            return
        }

        next.ServeHTTP(w, r)
    })
}