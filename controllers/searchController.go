package controllers

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/GDG-on-Campus-KHU/SC2_BE/models"
    "go.mongodb.org/mongo-driver/mongo"
)

var mongoClient *mongo.Client

func SetMongoClient(client *mongo.Client) {
    mongoClient = client
}

// HTML 태그 제거 함수
func removeHTMLTags(str string) string {
    return strings.ReplaceAll(strings.ReplaceAll(str, "<b>", ""), "</b>", "")
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

func NaverSearch(query string, display int) (*NaverSearchResponse, error) {
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

    var searchResponse NaverSearchResponse
    err = json.Unmarshal(body, &searchResponse)
    if err != nil {
        return nil, err
    }

    return &searchResponse, nil
}

func saveSearchResults(searchResponse *NaverSearchResponse) error{
    if mongoClient == nil {
        return fmt.Errorf("MongoDB client not initialized")
    }

    collection := mongoClient.Database("SC2_DB").Collection("placeList")

    for _, item := range searchResponse.Items {
        location := models.Location{
            Title:       removeHTMLTags(item.Title),
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

func NaverSearchHandler(c *gin.Context) {
    query := c.Query("query")
    display := 5

    // 네이버 API 검색 수행
    result, err := NaverSearch(query, display)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // MongoDB에 결과 저장 (HTML 태그가 제거된 상태로 저장됨)
    err = saveSearchResults(result)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save to database: %v", err)})
        return
    }

    // 응답에서도 HTML 태그 제거
    for i := range result.Items {
        result.Items[i].Title = removeHTMLTags(result.Items[i].Title)
    }

    // 간소화된 응답 생성
    simplifiedResults := make([]SimplifiedResponse, len(result.Items))
    for i, item := range result.Items {
        simplifiedResults[i] = SimplifiedResponse{
            Title:       removeHTMLTags(item.Title),
            RoadAddress: item.RoadAddress,
        }
    }

    c.JSON(http.StatusOK, gin.H{"items": simplifiedResults})
}

func DeleteAllPlacesHandler(c *gin.Context) {
    if mongoClient == nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not initialized"})
        return
    }

    collection := mongoClient.Database("SC2_DB").Collection("placeList")

    result, err := collection.DeleteMany(context.TODO(), struct{}{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete documents: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": fmt.Sprintf("Successfully deleted %d documents", result.DeletedCount),
    })
}