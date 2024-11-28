package service

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "strconv"
    "github.com/GDG-on-Campus-KHU/SC2_BE/models"
)

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

func (c *NavigateClient) GetNavigate(start, goal string) (*models.NavResponse, error) {
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

    var response models.NavResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("error parsing response: %v", err)
    }

    return &response, nil
}

func FormatCoordinate(coord string) string {
    coord = strings.TrimSuffix(coord, ",")
    numVal, err := strconv.ParseInt(coord, 10, 64)
    if err != nil {
        return ""
    }
    return fmt.Sprintf("%.6f", float64(numVal)/10000000)
}