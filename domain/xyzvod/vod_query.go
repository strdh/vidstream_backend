package xyzvod

import (
    "log"
    "xyzstream/config"
)

type Vod struct{
    Id int `json:"id"`
    VodUlid string `json:"vod_ulid"`
    Title string `json:"title"`
    Description string `json:"description"`
    Duration int `json:"duration"`
    Created int `json:"created"`
}

func Vods() ([]Vod, error) {
    var vods []Vod
    var temp Vod

    rows, err := config.DB.Query("SELECT * FROM vods ORDER BY created DESC LIMIT 30")
    if err != nil {
        log.Println(err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&temp.Id, &temp.VodUlid, &temp.Title, &temp.Description, &temp.Duration, &temp.Created)
        if err != nil {
            log.Println(err)
            return nil, err
        }

        vods = append(vods, temp)
    }

    return vods, nil
}

func VodsNext(id int) ([]Vod, error) {
    var vods []vod
    var temp Vod

    rows, err := config.DB.Query("SELECT * FROM vods WHERE id <= ? ORDER BY created DESC LIMIT 30", id)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&temp.Id, &temp.VodUlid, &temp.Title, &temp.Description, &temp.Duration, &temp.Created)
        if err != nil {
            log.Println(err)
            return nil, err
        }

        vods = append(vods, temp)
    }

    return vods, nil
}

func ByUlid(vodUlid string) (Vod, error) {
    var vod Vod
    row := config.DB.QueryRow("SELECT * FROM vods WHERE vod_ulid = ?", vodUlid)
    err := row.Scan(&vod.Id, &vod.VodUlid, &vod.Title, &vod.Description, &vod.Duration, &vod.Created)
    if err != nil {
        log.Println(err)
        return vod, err
    }

    return vod, nil
}

func Create(vod Vod) (int64, error) {
    result, err := config.DB.Exec("INSERT INTO vods (vod_ulid, title, description, duration, created) VALUES (?, ?, ?, ?, ?)", vod.VodUlid, vod.Title, vod.Description, vod.Duration, vod.Created)
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