package langsmithgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DatasetClient struct {
	baseClient
}

// NewDatasetClient creates a new LangSmith client
// The client requires an API key to authenticate requests.
// You can get an API key by signing up for a LangSmith account at https://smith.langchain.com
// The API key can be passed as an argument to the function or set as an environment variable LANGSMITH_API_KEY
func NewDatasetClient() (*DatasetClient, error) {
	if os.Getenv("LANGSMITH_API_KEY") == "" {
		return nil, errors.New("langsmith api key is required")
	}

	url := os.Getenv("LANGSMITH_URL")
	if url == "" {
		url = BASE_URL

	}

	return &DatasetClient{
		baseClient: baseClient{
			APIKey:  os.Getenv("LANGSMITH_API_KEY"),
			baseUrl: fmt.Sprintf("%s/datasets", url),
		},
	}, nil
}

func (d *DatasetClient) CreateDataset(input *Dataset) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}
	err = d.Do(d.baseUrl, http.MethodPost, jsonData)

	return err
}

func (d *DatasetClient) UploadCSV(input *DatasetCSV) error {
	body, contentType, err := input.ToMultiPart()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	_, err = io.Copy(&b, body)
	if err != nil {
		return err
	}

	err = d.PostForm(d.baseUrl+"/upload", &b, contentType)

	return err
}

func (d *DatasetClient) UploadExperiment(input *Experiment) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}
	err = d.Do(d.baseUrl+"/upload-experiment", http.MethodPost, jsonData)

	return err
}

func (d *DatasetClient) ReadDataset(id string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", d.baseUrl, id), nil)
	if err != nil {
		return nil, err
	}

	// Set the necessary headers
	req.Header.Set("x-api-key", d.APIKey)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = handleResponse(resp)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (d *DatasetClient) DownloadDatasetCsv(id string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/csv", d.baseUrl, id), nil)
	if err != nil {
		return nil, err
	}

	// Set the necessary headers
	req.Header.Set("x-api-key", d.APIKey)

	fmt.Println(req)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	_, err = io.Copy(&b, resp.Body)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *DatasetClient) GetExamples(datasetId string, offset int) ([]Example, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/examples?dataset=%s&offset=%d", BASE_URL, datasetId, offset), nil)
	if err != nil {
		return nil, err
	}

	// Set the necessary headers
	req.Header.Set("x-api-key", d.APIKey)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = handleResponse(resp)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	examples := []Example{}
	err = json.Unmarshal(b, &examples)
	if err != nil {
		return nil, err
	}
	return examples, nil
}

func (d *DatasetClient) CreateExample(example Example) error {
	jsonData, err := json.Marshal(example)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	err = d.Do(BASE_URL+"/examples", http.MethodPost, jsonData)

	return err
}

func (d *DatasetClient) CreateExamples(examples []Example) error {
	jsonData, err := json.Marshal(examples)
	if err != nil {
		return err
	}

	err = d.Do(BASE_URL+"/examples/bulk", http.MethodPost, jsonData)

	return err
}

func (d *DatasetClient) GetExamplesWithRuns(datasetId string) ([]Example, error) {
	body := map[string]any{
		"session_ids": []string{"58a7c5a2-c14e-42dc-936a-3af8e84777fa"},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/runs", d.baseUrl, datasetId), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Set the necessary headers
	req.Header.Set("x-api-key", d.APIKey)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = handleResponse(resp)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	examples := []Example{}
	err = json.Unmarshal(b, &examples)
	if err != nil {
		return nil, err
	}

	return examples, nil
}

func (d *DatasetClient) CreateComparativeExperiment(input *ComparativeExperimentRequest) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}
	err = d.Do(d.baseUrl+"/comparative", http.MethodPost, jsonData)

	return err
}

func (d *DatasetClient) ReadComparitiveExperiment(datasetId, experimentId string) ([]ComparativeExperiment, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/comparative?id=%s", d.baseUrl, datasetId, experimentId), nil)
	if err != nil {
		return nil, err
	}

	// Set the necessary headers
	req.Header.Set("x-api-key", d.APIKey)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = handleResponse(resp)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var experiments []ComparativeExperiment
	err = json.Unmarshal(b, &experiments)
	if err != nil {
		return nil, err
	}
	return experiments, nil
}

func (d *DatasetClient) QueryRuns(experimentIds []string, isRoot bool, queryParams *QueryParams) (*RunsResponse, error) {
	body := map[string]any{
		"session": experimentIds,
		"root":    isRoot,
		"select":  queryParams.Select,
		"filter":  queryParams.Filter,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/runs/query", BASE_URL), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Set the necessary headers
	req.Header.Set("x-api-key", d.APIKey)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = handleResponse(resp)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var runsResponse RunsResponse
	err = json.Unmarshal(b, &runsResponse)
	if err != nil {
		return nil, err
	}

	return &runsResponse, nil
}

func (d *DatasetClient) CreateFeedback(input *Feedback) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return d.Do(BASE_URL+"/feedback", http.MethodPost, jsonData)
}

func (d *DatasetClient) CreateTracerSession(input *TracerSessionRequest) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return d.Do(BASE_URL+"/sessions", http.MethodPost, jsonData)
}

func (d *DatasetClient) UpdateTracerSession(sessionId string, input *TracerSessionUpdate) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return d.Do(BASE_URL+"/sessions/"+sessionId, http.MethodPatch, jsonData)
}
