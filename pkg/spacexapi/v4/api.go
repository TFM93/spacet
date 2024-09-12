package spacexapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// api uses mongoose paginate
type queryBody struct {
	Query   map[string]interface{} `json:"query"`
	Options map[string]interface{} `json:"options"`
}

type paginatedResponse[T any] struct {
	Docs        []T  `json:"docs"`
	TotalDocs   int  `json:"totalDocs"`
	Limit       int  `json:"limit"`
	TotalPages  int  `json:"totalPages"`
	Page        int  `json:"page"`
	HasNextPage bool `json:"hasNextPage"`
}

type FiltersWithPagination struct {
	Query map[string]interface{}
	Limit int
	Page  int
}

func get[T, V any](c *client, ctx context.Context, endpoint string, filters *FiltersWithPagination) (hasMoreData bool, _ []V, _ error) {

	// set default pagination
	if filters == nil {
		filters = &FiltersWithPagination{
			Limit: 10,
			Page:  0,
			Query: make(map[string]interface{}),
		}
	}

	body := queryBody{
		Query: filters.Query,
		Options: map[string]interface{}{
			"limit": filters.Limit,
			"page":  filters.Page,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return hasMoreData, nil, fmt.Errorf("marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return hasMoreData, nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return hasMoreData, nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return hasMoreData, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var parsedResp paginatedResponse[V]
	if err := json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		return hasMoreData, nil, fmt.Errorf("decoding response: %w", err)
	}

	//todo- compare advantages over this:
	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := json.Unmarshal(data, &target); err != nil {
	// 	return nil, err
	// }

	return parsedResp.HasNextPage, parsedResp.Docs, nil
}

func getAll[T any](ctx context.Context, getFn func(context.Context, *FiltersWithPagination) (bool, []T, error), filters map[string]interface{}) ([]T, error) {
	var target []T

	filter := FiltersWithPagination{
		Query: filters,
		Limit: 50,
		Page:  0,
	}

	for {
		hasMoreData, landPads, err := getFn(ctx, &filter)
		if err != nil {
			return nil, err
		}
		target = append(target, landPads...)

		if !hasMoreData {
			break // all fetched
		}
		filter.Page++
	}
	return target, nil
}
