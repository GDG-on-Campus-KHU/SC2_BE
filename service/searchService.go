package service

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "strings"
    "github.com/GDG-on-Campus-KHU/SC2_BE/models"
    "github.com/GDG-on-Campus-KHU/SC2_BE/db"
)

// HTML 태그 제거 함수
func RemoveHTMLTags(str string) string {
    return strings.ReplaceAll(strings.ReplaceAll(str, "<b>", ""), "</b>", "")
}

func NaverSearch(query string, display int) (*models.NaverSearchResponse, error) {
    ClientID := os.Getenv("Naver_Client_ID")
    ClientSecret := os.Getenv("Naver_Secret")

    encodedQuery := url.QueryEscape(query)
    url := fmt.Sprintf("https://openapi.naver.com/v1/search/local.json?query=%s&display=%d", encodedQuery, display)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("X-Naver-Client-Id", ClientID)
    req.Header.Add("X-Naver-Client-Secret", ClientSecret)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var searchResponse models.NaverSearchResponse
    err = json.Unmarshal(body, &searchResponse)
    if err != nil {
        return nil, err
    }

    return &searchResponse, nil
}

func SaveSearchResults(searchResponse *models.NaverSearchResponse) error {
    if db.GetMongoClient() == nil {
        return fmt.Errorf("MongoDB client not initialized")
    }
    collection := db.GetMongoClient().Database("SC2_DB").Collection("placeList")

    for _, item := range searchResponse.Items {
        location := models.PlaceLocation{
            Title:       RemoveHTMLTags(item.Title),
            RoadAddress: item.RoadAddress,
            Mapx:        item.Mapx,
            Mapy:        item.Mapy,
        }

        _, err := collection.InsertOne(context.TODO(), location)
        if err != nil {
            return fmt.Errorf("failed to save location: %v", err)
        }
    }

    return nil
}
