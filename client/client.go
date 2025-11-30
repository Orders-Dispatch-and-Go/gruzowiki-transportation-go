package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type RoutesClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewRoutesClient(baseURL string) *RoutesClient {
    return &RoutesClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

type PotentialRoutesRequest struct {
    CargoRequestRouteID string   `json:"cargoRequestRouteId"`
    TripRouteIDs        []string `json:"tripRouteIds"`
}

type PotentialRoutesResponse struct {
    TripIDs []string `json:"tripIds"`
}

func (c *RoutesClient) GetPotentialTrips(cargoRequestRouteID string, tripRouteIDs []string) ([]string, error) {
    requestBody := PotentialRoutesRequest{
        CargoRequestRouteID: cargoRequestRouteID,
        TripRouteIDs:        tripRouteIDs,
    }

    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    url := c.baseURL + "/routes/trips/potential"
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
    }

    var response PotentialRoutesResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return response.TripIDs, nil
}