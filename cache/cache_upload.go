package cache

import (
    "time"
    "errors"
)

type UpCache struct{
    Expire int64 `json:"expire"`
    Created int64 `json:"created"`
    TotalChunk int `json:"total_chunk"`
    ChunkRemaining int `json:"chunk_remaining"`
}

var UpCacheMap = make(map[string]UpCache)

func GetUpCache(key string) (UpCache, error) {
    val, ok := UpCacheMap[key]
    if !ok || (val.Expire + val.Created) <= time.Now().Unix() {
        return val, errors.New("key not exists")
    }

    return val, nil
}

func SetUpCache(key string, val UpCache) error {
    _, ok := UpCacheMap[key]
    if ok {
        return errors.New("key already exists")
    }

    UpCacheMap[key] = val
    return nil
}

func DelUpCache(key string) error {
    _, ok := UpCacheMap[key]
    if !ok {
        return errors.New("key not exists")
    }

    delete(UpCacheMap, key)
    return nil
}