package xyzvod

import (
    "log"
    "xyzstream/config"
)

type Upload struct{
    Id int `json:"id"`
    IdUser int `json:"id_user"`
    UpUlid string `json:"up_ulid"`
    Title string `json:"title"`
    Description string `json:"description"`
    Size int `json:"size"`
    Progress int `json:"progress"`
    Created int64 `json:"created"`
    LastUpdated int64 `json:"last_update"`
    Status int `json:"status"`
}

type VUQuery struct{}

func (vu *VUQuery) Create(up Upload) (int64, error) {
    result, err := config.DB.Exec("INSERT INTO uploads (id_user, upulid, title, description, size, progress, created, lastupdate, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", up.IdUser, up.UpUlid, up.Title, up.Description, up.Size, up.Progress, up.Created, up.LastUpdated, up.Status)
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

func (vu *VUQuery) CheckUlid(upulid string) (bool, error) {
    var count int
    row := config.DB.QueryRow("SELECT COUNT(upulid) FROM uploads WHERE upulid = ?", upulid)
    err := row.Scan(&count)
    if err != nil {
        log.Println(err)
        return false, err
    }

    if count > 0 {
        return true, nil
    }

    return false, nil
}