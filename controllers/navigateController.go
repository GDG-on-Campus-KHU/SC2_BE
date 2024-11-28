package controllers

import(
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "strconv"
    "io/ioutil"
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/GDG-on-Campus-KHU/SC2_BE/models"
    "github.com/GDG-on-Campus-KHU/SC2_BE/db"
)

type RouteResponse struct {
    Start []float64   `json:"start"`
    Goal  []float64   `json:"goal"`
    Path  [][]float64 `json:"path"`
}

type Response struct {
    Code            int     `json:"code"`
    Message         string  `json:"message"`
    CurrentDateTime string  `json:"currentDateTime"`
    Route           Route   `json:"route"`
}

type Route struct {
    Traoptimal []TraoptimalRoute `json:"traoptimal"`
}

type TraoptimalRoute struct {
    Summary     Summary      `json:"summary"`
    Path        [][]float64 `json:"path"`
    Section     []Section   `json:"section"`
    Guide       []Guide     `json:"guide"`
}

type Summary struct {
    Start         Location    `json:"start"`
    Goal          Goal        `json:"goal"`
    Distance      int         `json:"distance"`
    Duration      int         `json:"duration"`
    DepartureTime string      `json:"departureTime"`
    Bbox          [][]float64 `json:"bbox"`
    TollFare      int         `json:"tollFare"`
    TaxiFare      int         `json:"taxiFare"`
    FuelPrice     int         `json:"fuelPrice"`
}

type Location struct {
    Location []float64 `json:"location"`
}

type Goal struct {
    Location []float64 `json:"location"`
    Dir      int       `json:"dir"`
}

type Section struct {
    PointIndex  int    `json:"pointIndex"`
    PointCount  int    `json:"pointCount"`
    Distance    int    `json:"distance"`
    Name        string `json:"name"`
    Congestion  int    `json:"congestion"`
    Speed       int    `json:"speed"`
}

type Guide struct {
    PointIndex    int    `json:"pointIndex"`
    Type          int    `json:"type"`
    Instructions  string `json:"instructions"`
    Distance      int    `json:"distance"`
    Duration      int    `json:"duration"`
}

type NavigateClient struct {
    ClientID     string
    ClientSecret string
    BaseURL      string
}

func NewNavigateClient(clientID, clientSecret string) *NavigateClient {
    return &NavigateClient{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        BaseURL:      "https://naveropenapi.apigw.ntruss.com/map-direction/v1/driving",
    }
}

func (c *NavigateClient) GetNavigate(start, goal string) (*Response, error) {
    req, err := http.NewRequest("GET", c.BaseURL, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Add("X-NCP-APIGW-API-KEY-ID", c.ClientID)
    req.Header.Add("X-NCP-APIGW-API-KEY", c.ClientSecret)

    q := req.URL.Query()
    q.Add("start", start)          
    q.Add("goal", goal)
    q.Add("option", "traoptimal")  
    req.URL.RawQuery = q.Encode()

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
    }

    var response Response
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("error parsing response: %v", err)
    }

    return &response, nil
}

func formatCoordinate(coord string) string {
    // 문자열의 마지막 쉼표 제거
    coord = strings.TrimSuffix(coord, ",")
    
    // 문자열을 int64로 변환
    numVal, err := strconv.ParseInt(coord, 10, 64)
    if err != nil {
        return ""
    }
    // 실수로 변환 (네이버 좌표계는 10자리로 되어있음)
    return fmt.Sprintf("%.6f", float64(numVal)/10000000)
}


func GetNavigateHandler(c *gin.Context) {
    if db.GetMongoClient() == nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not initialized"})
        return
    }
    collection := db.GetMongoClient().Database("SC2_DB").Collection("placeList")

    clientID := os.Getenv("Naver_Cloud_Client_ID")
    clientSecret := os.Getenv("Naver_Cloud_Client_Secret")

    client := NewNavigateClient(clientID, clientSecret)

    // Query 파라미터에서 시작 좌표와 목적지 이름 받기
    start := c.Query("start")
    // example
    //start := "127.0584795,37.2484202"
    locationName := c.Query("location")
    fmt.Println(locationName)

    if start == "" || locationName == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "start coordinates and location name are required",
        })
        return
    }
    
    var goalPlace models.Location
    err := collection.FindOne(context.TODO(), gin.H{"title": locationName}).Decode(&goalPlace)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": fmt.Sprintf("location not found: %v", err),
        })
        return
    }

    // 좌표 변환 및 문자열 생성
    goalMapx := formatCoordinate(goalPlace.Mapx)
    goalMapy := formatCoordinate(goalPlace.Mapy)
    goal := fmt.Sprintf("%s,%s", goalMapx, goalMapy)

    // 네비게이션 API 호출
    response, err := client.GetNavigate(start, goal)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    if len(response.Route.Traoptimal) > 0 {
        routeResponse := RouteResponse{
            Start: response.Route.Traoptimal[0].Summary.Start.Location,
            Goal:  response.Route.Traoptimal[0].Summary.Goal.Location,
            Path:  response.Route.Traoptimal[0].Path,
        }
        c.JSON(http.StatusOK, routeResponse)
    } else {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "no route found",
        })
    }
}