package models

type PlaceLocation struct {
    Title       string `bson:"title"`
    RoadAddress string `bson:"roadAddress"`
    Mapx        string `bson:"mapx"`
    Mapy        string `bson:"mapy"`
}

type LocationResponse struct {
    Items []struct {
        Title       string `json:"title"`
        RoadAddress string `json:"roadAddress"`
        Mapx        string `json:"mapx"`
        Mapy        string `json:"mapy"`
    } `json:"items"`
}

type NaverSearchResponse struct {
    Items []SearchItem `json:"items"`
}

type SimplifiedResponse struct {
    Title       string `json:"title"`
    RoadAddress string `json:"roadAddress"`
}

type SearchItem struct {
    Title       string `json:"title"`
    Category    string `json:"category"`
    Description string `json:"description"`
    Telephone   string `json:"telephone"`
    RoadAddress string `json:"roadAddress"`
    Mapx        string `json:"mapx"`
    Mapy        string `json:"mapy"`
}