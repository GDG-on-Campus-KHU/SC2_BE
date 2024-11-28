package models

type RouteResponse struct {
    Start []float64   `json:"start"`
    Goal  []float64   `json:"goal"`
    Path  [][]float64 `json:"path"`
}

type NavResponse struct {
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
    Start         NavLocation    `json:"start"`
    Goal          Goal        `json:"goal"`
    Distance      int         `json:"distance"`
    Duration      int         `json:"duration"`
    DepartureTime string      `json:"departureTime"`
    Bbox          [][]float64 `json:"bbox"`
    TollFare      int         `json:"tollFare"`
    TaxiFare      int         `json:"taxiFare"`
    FuelPrice     int         `json:"fuelPrice"`
}

type NavLocation struct {
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