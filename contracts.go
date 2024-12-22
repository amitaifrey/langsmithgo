package langsmithgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"time"
)

const (
	BASE_URL = "https://api.smith.langchain.com/api/v1"
)

type Response struct {
	Detail string `json:"detail"`
}

type Event struct {
	EventName string `json:"event_name"`
	Reason    string `json:"reason,omitempty"`
	Value     string `json:"value,omitempty"`
}

type RunPayload struct {
	RunID              string                 `json:"id"`
	Name               string                 `json:"name"`
	RunType            RunType                `json:"run_type"`
	StartTime          time.Time              `json:"start_time"`
	Inputs             map[string]interface{} `json:"inputs"`
	ParentID           string                 `json:"parent_run_id,omitempty"`
	SessionID          string                 `json:"session_id,omitempty"`
	SessionName        string                 `json:"session_name,omitempty"`
	Tags               []string               `json:"tags,omitempty"`
	Outputs            map[string]interface{} `json:"outputs,omitempty"`
	EndTime            time.Time              `json:"end_time,omitempty"`
	Extras             map[string]interface{} `json:"extra,omitempty"`
	Events             []Event                `json:"events,omitempty"`
	Error              string                 `json:"error,omitempty"`
	ReferenceExampleID string                 `json:"reference_example_id,omitempty"`
}

type Client struct {
	baseClient
	projectName string // project name in LangSmith
}

type baseClient struct {
	APIKey  string // API key for LangSmith
	baseUrl string // base url for the LangSmith API
}

type SimplePayload struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	RunType     RunType                `json:"run_type"`
	StartTime   time.Time              `json:"start_time"`
	Inputs      map[string]interface{} `json:"inputs"`
	SessionID   string                 `json:"session_id"`
	SessionName string                 `json:"session_name"`
	Tags        []string               `json:"tags,omitempty"`
	ParentId    string                 `json:"parent_run_id,omitempty"`
	Extras      map[string]interface{} `json:"extra,omitempty"`
	Events      []Event                `json:"events,omitempty"`
	Outputs     map[string]interface{} `json:"outputs"`
	EndTime     time.Time              `json:"end_time"`
}

type PostPayload struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	RunType            RunType                `json:"run_type"`
	StartTime          time.Time              `json:"start_time"`
	Inputs             map[string]interface{} `json:"inputs"`
	SessionID          string                 `json:"session_id,omitempty"`
	SessionName        string                 `json:"session_name,omitempty"`
	Tags               []string               `json:"tags,omitempty"`
	ParentId           string                 `json:"parent_run_id,omitempty"`
	Extras             map[string]interface{} `json:"extra,omitempty"`
	Events             []Event                `json:"events,omitempty"`
	ReferenceExampleID string                 `json:"reference_example_id,omitempty"`
}

type PatchPayload struct {
	Outputs   map[string]interface{} `json:"outputs"`
	EndTime   time.Time              `json:"end_time"`
	Events    []Event                `json:"events,omitempty"`
	Extras    map[string]interface{} `json:"extra,omitempty"`
	Error     string                 `json:"error,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
}

type RunType string

// Enum values using iota
const (
	Tool      RunType = "tool"
	Chain     RunType = "chain"
	LLM       RunType = "llm"
	Retriever RunType = "retriever"
	Embedding RunType = "embedding"
	Prompt    RunType = "prompt"
	Parser    RunType = "parser"
)

// Dataset represents the Dataset schema
type Dataset struct {
	ID                      string                  `json:"id" binding:"required,uuid"`
	Name                    string                  `json:"name" binding:"required"`
	Description             *string                 `json:"description,omitempty"`
	CreatedAt               time.Time               `json:"created_at"`
	InputsSchemaDefinition  map[string]interface{}  `json:"inputs_schema_definition,omitempty"`
	OutputsSchemaDefinition map[string]interface{}  `json:"outputs_schema_definition,omitempty"`
	ExternallyManaged       *bool                   `json:"externally_managed,omitempty" default:"false"`
	Transformations         []DatasetTransformation `json:"transformations,omitempty"`
	DataType                *DataType               `json:"data_type,omitempty" default:"kv"`
	TenantID                string                  `json:"tenant_id" binding:"required,uuid"`
	ExampleCount            int                     `json:"example_count" binding:"required"`
	SessionCount            int                     `json:"session_count" binding:"required"`
	ModifiedAt              time.Time               `json:"modified_at" binding:"required"`
	LastSessionStartTime    *time.Time              `json:"last_session_start_time,omitempty"`
}

type DatasetTransformation struct {
	Path               []string                  `json:"path"`
	TransformationType DatasetTransformationType `json:"transformation_type"`
}

type DatasetTransformationType string

const (
	RemoveSystemMessages   DatasetTransformationType = "remove_system_messages"
	ConvertToOpenAIMessage DatasetTransformationType = "convert_to_openai_message"
	ConvertToOpenAITool    DatasetTransformationType = "convert_to_openai_tool"
	RemoveExtraFields      DatasetTransformationType = "remove_extra_fields"
	ExtractToolsFromRun    DatasetTransformationType = "extract_tools_from_run"
)

type DataType string

const (
	DataTypeKV  DataType = "kv"
	DataTypeLLM DataType = "llm"
	DataTypeCSV DataType = "csv"
)

type DatasetCSV struct {
	File        string   `json:"file"`
	InputKeys   []string `json:"input_keys"`
	Name        string   `json:"name,omitempty"`
	DataType    DataType `json:"data_type"`
	OutputKeys  []string `json:"output_keys,omitempty"`
	Description string   `json:"description,omitempty"`
}

func (d *DatasetCSV) ToMultiPart() (io.Reader, string, error) {
	b := bytes.NewBuffer(nil)
	mw := multipart.NewWriter(b)

	if d.File != "" {
		fileWriter, err := mw.CreateFormFile("file", "dataset.csv")
		if err != nil {
			return nil, "", err
		}
		_, err = fileWriter.Write([]byte(d.File))
		if err != nil {
			return nil, "", err
		}
	}

	for _, key := range d.InputKeys {
		err := mw.WriteField("input_keys", key)
		if err != nil {
			return nil, "", err
		}
	}

	if d.Name != "" {
		err := mw.WriteField("name", d.Name)
		if err != nil {
			return nil, "", err
		}
	}

	err := mw.WriteField("data_type", string(d.DataType))
	if err != nil {
		return nil, "", err
	}

	for _, key := range d.OutputKeys {
		err := mw.WriteField("output_keys", key)
		if err != nil {
			return nil, "", err
		}
	}

	if d.Description != "" {
		err := mw.WriteField("description", d.Description)
		if err != nil {
			return nil, "", err
		}
	}

	err = mw.Close()
	if err != nil {
		return nil, "", err
	}
	return b, mw.FormDataContentType(), nil
}

type Experiment struct {
	ExperimentName          string             `json:"experiment_name"`
	ExperimentDescription   string             `json:"experiment_description,omitempty"`
	DatasetID               string             `json:"dataset_id,omitempty"`
	DatasetName             string             `json:"dataset_name,omitempty"`
	DatasetDescription      string             `json:"dataset_description,omitempty"`
	SummaryExperimentScores []*ExperimentScore `json:"summary_experiment_scores,omitempty"`
	Results                 []*Result          `json:"results"`
	ExperimentStartTime     LangsmithTime      `json:"experiment_start_time"`
	ExperimentEndTime       LangsmithTime      `json:"experiment_end_time"`
	ExperimentMetadata      map[string]any     `json:"experiment_metadata,omitempty"`
}

type ExperimentScore struct {
	CreatedAt               LangsmithTime   `json:"created_at,omitempty"`
	ModifiedAt              LangsmithTime   `json:"modified_at,omitempty"`
	Key                     string          `json:"key"`
	Score                   float64         `json:"score,omitempty"`
	Value                   float64         `json:"value,omitempty"`
	Comment                 string          `json:"comment,omitempty"`
	Correction              map[string]any  `json:"correction,omitempty"`
	FeedbackGroupID         string          `json:"feedback_group_id,omitempty"`
	ComparativeExperimentID string          `json:"comparative_experiment_id,omitempty"`
	ID                      string          `json:"id,omitempty"`
	FeedbackSource          *FeedbackSource `json:"feedback_source,omitempty"`
	FeedbackConfig          *FeedbackConfig `json:"feedback_config,omitempty"`
	Extra                   map[string]any  `json:"extra,omitempty"`
}

type FeedbackSource struct {
	Type     string         `json:"type"`
	Metadata map[string]any `json:"metadata"`
}

type FeedbackConfig struct {
	Type       FeedbackType       `json:"type"`
	Min        float64            `json:"min"`
	Max        float64            `json:"max"`
	Categories []FeedbackCategory `json:"categories"`
}

type FeedbackType string

const (
	FeedbackTypeContinuous  FeedbackType = "continuous"
	FeedbackTypeCategorical FeedbackType = "categorical"
	FeedbackTypeFreeform    FeedbackType = "freeform"
)

type FeedbackCategory struct {
	Value float64 `json:"value"`
	Label string  `json:"label"`
}

type Result struct {
	RowID            string             `json:"row_id,omitempty"`
	Inputs           map[string]any     `json:"inputs"`
	ExpectedOutputs  map[string]any     `json:"expected_outputs,omitempty"`
	ActualOutputs    map[string]any     `json:"actual_outputs,omitempty"`
	EvaluationScores []*ExperimentScore `json:"evaluation_scores,omitempty"`
	StartTime        LangsmithTime      `json:"start_time"`
	EndTime          LangsmithTime      `json:"end_time"`
	RunName          string             `json:"run_name,omitempty"`
	Error            string             `json:"error,omitempty"`
	RunMetadata      map[string]any     `json:"run_metadata,omitempty"`
}

const langsmithTimeFormat = "2006-01-02T15:04:05"

type LangsmithTime time.Time

func (t LangsmithTime) MarshalJSON() ([]byte, error) {
	// Format the time using the custom format
	formattedTime := fmt.Sprintf("\"%s\"", time.Time(t).Format(langsmithTimeFormat))
	return []byte(formattedTime), nil
}

func (t *LangsmithTime) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	t2, err := time.Parse(langsmithTimeFormat, s[:19])
	if err != nil {
		return err
	}
	*t = LangsmithTime(t2)
	return nil
}

type Example struct {
	ID             string                 `json:"id"`
	CreatedAt      *LangsmithTime         `json:"created_at,omitempty"`
	ModifiedAt     string                 `json:"modified_at,omitempty"`
	Name           string                 `json:"name,omitempty"`
	DatasetID      string                 `json:"dataset_id"`
	SourceRunID    string                 `json:"source_run_id,omitempty"`
	Metadata       map[string]any         `json:"metadata,omitempty"`
	Inputs         map[string]any         `json:"inputs"`
	Outputs        map[string]any         `json:"outputs,omitempty"`
	AttachmentURLs map[string]interface{} `json:"attachment_urls,omitempty"`
	Runs           []Run                  `json:"runs,omitempty"`
}

type Run struct {
	Name               string         `json:"name"`
	Inputs             map[string]any `json:"inputs,omitempty"`
	InputsPreview      string         `json:"inputs_preview,omitempty"`
	RunType            RunType        `json:"run_type"`
	StartTime          *LangsmithTime `json:"start_time,omitempty"`
	EndTime            *LangsmithTime `json:"end_time,omitempty"`
	Extra              map[string]any `json:"extra,omitempty"`
	Error              string         `json:"error,omitempty"`
	ExecutionOrder     int            `json:"execution_order,omitempty"`
	Serialized         map[string]any `json:"serialized,omitempty"`
	Outputs            map[string]any `json:"outputs,omitempty"`
	OutputsPreview     string         `json:"outputs_preview,omitempty"`
	ParentRunID        string         `json:"parent_run_id,omitempty"`
	ManifestID         string         `json:"manifest_id,omitempty"`
	ManifestS3ID       string         `json:"manifest_s3_id,omitempty"`
	Events             []Event        `json:"events,omitempty"`
	Tags               []string       `json:"tags,omitempty"`
	InputsS3URLs       map[string]any `json:"inputs_s3_urls,omitempty"`
	OutputsS3URLs      map[string]any `json:"outputs_s3_urls,omitempty"`
	S3URLs             map[string]any `json:"s3_urls,omitempty"`
	TraceID            string         `json:"trace_id"`
	DottedOrder        string         `json:"dotted_order,omitempty"`
	ID                 string         `json:"id"`
	SessionID          string         `json:"session_id"`
	ReferenceExampleID string         `json:"reference_example_id,omitempty"`
	TotalTokens        int            `json:"total_tokens,omitempty"`
	PromptTokens       int            `json:"prompt_tokens,omitempty"`
	CompletionTokens   int            `json:"completion_tokens,omitempty"`
	TotalCost          string         `json:"total_cost,omitempty"`
	PromptCost         string         `json:"prompt_cost,omitempty"`
	CompletionCost     string         `json:"completion_cost,omitempty"`
	Status             string         `json:"status"`
	FeedbackStats      map[string]any `json:"feedback_stats,omitempty"`
	AppPath            string         `json:"app_path,omitempty"`
}

// QueryParams represents the query parameters for the API.
type QueryParams struct {
	ID               []string       `form:"id"`
	AsOf             string         `form:"as_of"`
	Metadata         map[string]any `form:"metadata"`
	FullTextContains []string       `form:"full_text_contains"`
	Splits           []string       `form:"splits"`
	Dataset          string         `form:"dataset"`
	Offset           int            `form:"offset,default=0"`
	Limit            int            `form:"limit,default=100"`
	Order            string         `form:"order,default=recent"`
	RandomSeed       *float64       `form:"random_seed"`
	Select           []string       `form:"select,default=id,created_at,modified_at,name,dataset_id,source_run_id,metadata,inputs,outputs"`
	Filter           string         `form:"filter"`
}

type ComparativeExperimentRequest struct {
	ID                 string         `json:"id,omitempty"`
	ExperimentIDs      []string       `json:"experiment_ids"`
	Name               *string        `json:"name,omitempty"`
	Description        *string        `json:"description,omitempty"`
	ReferenceDatasetId string         `json:"reference_dataset_id"`
	CreatedAt          LangsmithTime  `json:"created_at"`
	ModifiedAt         LangsmithTime  `json:"modified_at"`
	Extra              map[string]any `json:"extra,omitempty"`
}

type ComparativeExperiment struct {
	ID                 string            `json:"id,omitempty"`
	Name               *string           `json:"name,omitempty"`
	Description        *string           `json:"description,omitempty"`
	ReferenceDatasetId string            `json:"reference_dataset_id"`
	CreatedAt          LangsmithTime     `json:"created_at"`
	ModifiedAt         LangsmithTime     `json:"modified_at"`
	Extra              map[string]any    `json:"extra,omitempty"`
	ExperimentsInfo    []*ExperimentInfo `json:"experiments_info"`
	FeedbackStats      map[string]any    `json:"feedback_stats,omitempty"`
}

type ExperimentInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Feedback struct {
	ID                      string          `json:"id,omitempty"`
	CreatedAt               LangsmithTime   `json:"created_at,omitempty"`
	ModifiedAt              LangsmithTime   `json:"modified_at,omitempty"`
	Key                     string          `json:"key"`
	Score                   any             `json:"score,omitempty"`
	Value                   any             `json:"value,omitempty"`
	Comment                 string          `json:"comment,omitempty"`
	Correction              any             `json:"correction,omitempty"`
	FeedbackGroupID         string          `json:"feedback_group_id,omitempty"`
	ComparativeExperimentID string          `json:"comparative_experiment_id,omitempty"`
	RunID                   string          `json:"run_id,omitempty"`
	SessionID               string          `json:"session_id,omitempty"`
	TraceID                 string          `json:"trace_id,omitempty"`
	FeedbackSource          *FeedbackSource `json:"feedback_source,omitempty"`
	FeedbackConfig          *FeedbackConfig `json:"feedback_config,omitempty"`
}

type RunsResponse struct {
	Runs        []Run          `json:"runs"`
	Cursors     map[string]any `json:"cursors"`
	ParsedQuery *string        `json:"parsed_query,omitempty"`
}

type TracerSessionRequest struct {
	ID                 string         `json:"id,omitempty"`
	StartTime          LangsmithTime  `json:"start_time,omitempty"`
	EndTime            LangsmithTime  `json:"end_time,omitempty"`
	Extra              map[string]any `json:"extra,omitempty"`
	Name               string         `json:"name,omitempty"`
	Description        string         `json:"description,omitempty"`
	DefaultDatasetID   string         `json:"default_dataset_id,omitempty"`
	ReferenceDatasetID string         `json:"reference_dataset_id,omitempty"`
	TraceTier          string         `json:"trace_tier,omitempty"`
}

type TracerSessionUpdate struct {
	EndTime          LangsmithTime  `json:"end_time,omitempty"`
	Extra            map[string]any `json:"extra,omitempty"`
	Name             string         `json:"name,omitempty"`
	Description      string         `json:"description,omitempty"`
	DefaultDatasetID string         `json:"default_dataset_id,omitempty"`
	TraceTier        string         `json:"trace_tier,omitempty"`
}
