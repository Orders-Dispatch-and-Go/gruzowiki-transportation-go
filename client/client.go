package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gruzowiki/rest/models"
	"io"
	"net/http"
	"strings"
	"time"
)

type FeignClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewRoutesClient(baseURL string) *FeignClient {
	return &FeignClient{
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

func (c *FeignClient) GetPotentialTrips(cargoRequestRouteID string, tripRouteIDs []string) ([]string, error) {
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

	err = c.logFeignClientRequest(url, requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	err = c.logFeignClientResponse(resp)

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

func (c *FeignClient) MergeRoutes(tripRouteID string, cargoRequestRouteIDs []string) (uuid.UUID, error) {
	requestBody := MergeCargoRequestIntoTripRoute{
		TripRouteID:         tripRouteID,
		CargoRequestRouteID: cargoRequestRouteIDs,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.baseURL + "/routes/trips/merge"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if err = c.logFeignClientRequest(url, requestBody); err != nil {
		return uuid.Nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to make request: %w", err)
	}

	if err = c.logFeignClientResponse(resp); err != nil {
		return uuid.Nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return uuid.Nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		RouteID uuid.UUID `json:"routeId"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return uuid.Nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.RouteID, nil
}

func (c *FeignClient) CreateRouteForCargoRequest(
	request models.PostCargoRequestRequest,
	fromStationId uuid.UUID,
	toStationId uuid.UUID,
) (*uuid.UUID, error) {
	url := c.baseURL + "/routes/cargo_requests"
	response, err := c.CreateRoute(url, request.FromStation, request.ToStation, fromStationId, toStationId)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *FeignClient) CreateRouteForTrip(
	request models.CreateTripRequest,
	fromStationId uuid.UUID,
	toStationId uuid.UUID,
) (*uuid.UUID, error) {
	url := c.baseURL + "/routes/trips"
	response, err := c.CreateRoute(url, request.FromStation, request.ToStation, fromStationId, toStationId)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *FeignClient) CreateRoute(
	url string,
	fromStation models.Station,
	toStation models.Station,
	fromStationId uuid.UUID,
	toStationId uuid.UUID,
) (*uuid.UUID, error) {
	requestBody := CreateRouteRequestBody{
		FromStation: StationDTO{
			ID:      fromStationId,
			Address: fromStation.Address,
			Coords:  StationCoords{Lat: fromStation.Coords.Lat, Lon: fromStation.Coords.Lon},
		},
		ToStation: StationDTO{
			ID:      toStationId,
			Address: toStation.Address,
			Coords:  StationCoords{Lat: toStation.Coords.Lat, Lon: toStation.Coords.Lon},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	err = c.logFeignClientRequest(url, requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	err = c.logFeignClientResponse(resp)

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var response CreateRouteResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.ID, nil
}

func (c *FeignClient) GetRouteForCargoRequest(cargoRequestRouteID uuid.UUID) (*models.GetTripRouteResponse, error) {
	url := c.baseURL + "/routes/cargo_requests/"
	response, err := c.GetRoute(url, cargoRequestRouteID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *FeignClient) GetRouteForTrip(tripRouteId uuid.UUID) (*models.GetTripRouteResponse, error) {
	url := c.baseURL + "/routes/trips/"
	response, err := c.GetRoute(url, tripRouteId)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *FeignClient) GetRoute(url string, routeId uuid.UUID) (*models.GetTripRouteResponse, error) {
	url += routeId.String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	err = c.logFeignClientRequest(url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	err = c.logFeignClientResponse(resp)

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var response models.GetTripRouteResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (c *FeignClient) logFeignClientRequest(url string, request interface{}) error {
	fmt.Printf("\nFeign Client Request Url: %s", url)
	var body = "nil"
	if request != nil {
		jsonData, err := json.Marshal(request)
		if err != nil {
			fmt.Println("failed to log request in feign client request:", err)
			return err
		}
		jsonDataString := string(jsonData)
		body = strings.ReplaceAll(jsonDataString, "\n", "")
		body = strings.ReplaceAll(body, " ", "")
	}
	fmt.Printf("\nFeign Client Request Body: %s", body)
	return nil
}

func (c *FeignClient) logFeignClientResponse(response *http.Response) error {
	fmt.Printf("\nFeign Client Response Status:%s", response.Status)
	if response.Body != nil {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("failed to read feing client response body:", err)
			return err
		}
		bodyString := string(bodyBytes)
		fmt.Printf("\nFeign Client Response Body: %s", bodyString)
		response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return nil
}
