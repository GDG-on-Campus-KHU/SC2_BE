package controllers

import(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

type RouteResponse struct {
    Start []float64   `json:"start"`
    Goal  []float64   `json:"goal"`
    Path  [][]float64 `json:"path"`
}


type Response struct{
	Code				int		`json:"code"`
	Message				string	`json:"message"`
	CurrentDateTime		string	`json:"currentDateTime"`
	Route				Route	`json:"route"`
}

type Route struct{
	Traoptimal	[]TraoptimalRoute	`json:"traoptimal"`
}

type TraoptimalRoute struct{
	Summary		Summary		`json:"summary"`
	Path		[][]float64	`json:"path"`
	Section		[]Section	`json:"section"`
	Guide		[]Guide		`json:"guide"`
}

type Summary struct {
    Start         Location 		`json:"start"`
    Goal          Goal     		`json:"goal"`
    Distance      int      		`json:"distance"`
    Duration      int      		`json:"duration"`
    DepartureTime string   		`json:"departureTime"`
    Bbox          [][]float64	`json:"bbox"`
    TollFare      int      		`json:"tollFare"`
    TaxiFare      int      		`json:"taxiFare"`
    FuelPrice     int      		`json:"fuelPrice"`
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

func GetNavigateHandler(c *gin.Context) {
    clientID := os.Getenv("Naver_Cloud_Client_ID")
    clientSecret := os.Getenv("Naver_Cloud_Client_Secret")

    client := NewNavigateClient(clientID, clientSecret)

    // Query 파라미터에서 좌표 받기
    start := c.Query("start")
    goal := c.Query("goal")

    /*example
	start := "127.0584795,37.2484202"
	goal := "127.0675645,37.2502671"
    */

    if start == "" || goal == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "start and goal coordinates are required",
        })
        return
    }

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