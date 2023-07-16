package xyzvod

type VodUploadReq struct{
    Title string `json:"title"`
    Description string `json:"description"`
    Size int `json:"size"`
    Ext string `json:"ext"`
}

func ValidateUpload(request VodUploadReq) (bool, map[string]string) {
    messages := map[string]string{}

    if request.Title == "" {
        messages["title"] = "Title is required"
    }

    if len(request.Title) > 100 || len (request.Title) < 3 {
        messages["title"] = "Title must be between 3 and 100 characters"
    }

    if request.Description == "" {
        messages["description"] = "Description is required"
    }

    if len(request.Description) > 300 || len(request.Description) < 3 {
        messages["description"] = "Description must be between 3 and 300 characters"
    }

    if request.Size == 0 {
        messages["size"] = "Size is required"
    }

    return len(messages) == 0, messages
}