package models

type Location struct {
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