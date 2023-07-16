package main

import (
    "os"
    "log"
    "fmt"
    "net/http"
    "xyzstream/utils"
    "xyzstream/config"
    "xyzstream/middleware"
    "xyzstream/domain/xyzvod"
    "xyzstream/domain/xyzauth"
    "github.com/joho/godotenv"
    "github.com/gorilla/mux"
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

    router.HandleFunc("/register", xyzauth.Register).Methods("POST", "OPTIONS")
    router.HandleFunc("/login", xyzauth.Login).Methods("POST", "OPTIONS")

    router.HandleFunc("/vodupload", xyzvod.VodUpload).Methods("POST", "OPTIONS")
    router.HandleFunc("/vodupload/{upulid}/chunk", xyzvod.HandleChunk).Methods("POST", "OPTIONS")
    router.HandleFunc("/vodupload/{upulid}", xyzvod.ContinueUpload).Methods("POST", "OPTIONS")
    router.HandleFunc("/vod", middleware.JWTMiddleware(xyzvod.VodList)).Methods("GET", "OPTIONS")
    router.HandleFunc("/vod/next/{id}", xyzvod.VodListNext).Methods("GET", "OPTIONS")
    router.HandleFunc("/vod/{vodulid}", xyzvod.VodDetail).Methods("GET", "OPTIONS")
    router.HandleFunc("/vod/stream/{segment}", xyzvod.VodStream).Methods("GET", "OPTIONS")

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