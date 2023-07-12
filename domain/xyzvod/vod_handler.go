package xyzvod

import (
    "io"
    "os"
    "log"
    "bytes"
    "strconv"
    "net/http"
    "xyzstream/utils"
    "github.com/gorilla/mux"
)

func VodUpload(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        utils.WriteResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }
}

func VodList(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        utils.WriteResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    vods, err := Vods()
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

    vods, err := VodsNext(id)
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

    vodDetail, err := ByUlid(vodulid)
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