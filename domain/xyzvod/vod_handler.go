package xyzvod

import (
    "fmt"
    "io"
    "os"
    "log"
    "time"
    "bytes"
    "strconv"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "path/filepath"
    "xyzstream/cache"
    "xyzstream/utils"
    "github.com/gorilla/mux"
    "github.com/oklog/ulid/v2"
)

func VodUpload(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        utils.WriteResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        log.Println(err)
        return
    }
    defer r.Body.Close()

    var request VodUploadReq
    err = json.Unmarshal(body, &request)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusBadRequest, "Invalid request body", nil)
        log.Println(err)
        return
    }

    isValid, messages := ValidateUpload(request)
    if !isValid {
        utils.WriteResponse(w, r, http.StatusBadRequest, "Invalid request", messages)
        return
    }

    ulid := ulid.Make()
    finalUlid := ulid.String()

    upload := Upload{
        IdUser: 1,
        UpUlid: finalUlid,
        Title: request.Title,
        Description: request.Description,
        Size: request.Size,
        Ext: request.Ext,
        Progress: 0,
        Created: time.Now().Unix(),
        LastUpdated: time.Now().Unix(),
        Status: 0,
    }

    query := &VUQuery{}
    _, err = query.Create(upload)
    if err != nil {
        log.Println(err)
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        log.Println(err)
        return
    }

    tempFilePath := filepath.Join(os.Getenv("XYZ1_TEMP"), finalUlid + request.Ext)
    tempFile, err := os.Create(tempFilePath)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }
    defer tempFile.Close()

    upCache := cache.UpCache{
        Expire: 3600,
        Created: time.Now().Unix(),
        TotalChunk: request.TotalChunks,
        ChunkRemaining: request.TotalChunks,
    }

    err = cache.SetUpCache(finalUlid, upCache)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }

    utils.WriteResponse(w, r, http.StatusOK, "OK", finalUlid)
}

func ContinueUpload(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        utils.WriteResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    vars := mux.Vars(r)
    upUlid := vars["upulid"]

    query := &VUQuery{}
    id, err := query.CheckUlid(upUlid)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusNotFound, "Not found", nil)
        return
    }

    utils.WriteResponse(w, r, http.StatusOK, "OK", id)
}

func HandleChunk(w http.ResponseWriter, r *http.Request) {
    fmt.Println("handle chunk")
    if r.Method != http.MethodPost {
        utils.WriteResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    vars := mux.Vars(r)
    upUlid := vars["upulid"]

    //check cache
    chunkStatus, err := cache.GetUpCache(upUlid)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusBadRequest, "Bad Request", nil)
        return
    }

    //check chunk remaining
    chunkRemaining := chunkStatus.ChunkRemaining
    if chunkRemaining == 0 {
        utils.WriteResponse(w, r, http.StatusBadRequest, "Bad Request", nil)
        return
    }

    //extract chunk data from request
    chunkData, err := ioutil.ReadAll(r.Body)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }
    defer r.Body.Close()

    //open temp file
    tempFilePath := filepath.Join(os.Getenv("XYZ1_TEMP"), upUlid + ".temp")
    tempFile, err := os.OpenFile(tempFilePath, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }

    _, err = io.Copy(tempFile, bytes.NewReader(chunkData))
    if err != nil {
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }

    //update cache
    chunkStatus.ChunkRemaining = chunkRemaining - 1
    cache.UpCacheMap[upUlid] = chunkStatus

    // if chunkRemaining == 0 {
    //     //process the video using ffmpeg and send to object storage
    // }

    utils.WriteResponse(w, r, http.StatusOK, "OK", nil)
}

func VodList(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        utils.WriteResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    query := &VodQuery{}
    vods, err := query.Vods()
    if err != nil {
        log.Println(err)
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }

    utils.WriteResponse(w, r, http.StatusOK, "OK", vods)
}

func VodListNext(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        utils.WriteResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Println(err)
        utils.WriteResponse(w, r, http.StatusBadRequest, "Bad request", nil)
        return
    }

    query := &VodQuery{}
    vods, err := query.VodsNext(id)
    if err != nil {
        log.Println(err)
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }

    utils.WriteResponse(w, r, http.StatusOK, "OK", vods)
}

func VodDetail(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        utils.WriteResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    vars := mux.Vars(r)
    vodulid := vars["vodulid"]

    query := &VodQuery{}
    vodDetail, err := query.ByUlid(vodulid)
    if err != nil {
        log.Println(err)
        utils.WriteResponse(w, r, http.StatusNotFound, "Not found", nil)
        return
    }

    utils.WriteResponse(w, r, http.StatusOK, "OK", vodDetail)
    return
}

func VodStream(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        utils.WriteResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    vars := mux.Vars(r)
    segment := vars["segment"]

    //read the data from s3 object storage
    segmentData, err := utils.ObjRead(os.Getenv("XYZ1_BUCKET"), segment)
    if err != nil {
        log.Println(err)
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }

    contentType := "video/mp2t"
    w.Header().Set("Content-Type", contentType)
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

    reader := bytes.NewReader(segmentData)
    _, err = io.Copy(w, reader)
    if err != nil {
        log.Println(err)
        utils.WriteResponse(w, r, http.StatusInternalServerError, "Internal server error", nil)
        return
    }
}

// func getContentType(filename string) string {
//     extension := filepath.Ext(filename)
//     switch extension {
//     case ".ts":
//         return "video/mp2t"
//     case ".m3u8":
//         return "application/x-mpegURL"
//     default:
//         return "application/octet-stream"
//     }
// }