package types

type Prediction struct {
	JobID                   string  `json:"job_id"`
	PredictedMemoryMB       uint32  `json:"predicted_memory_mb"`
	PredictedRuntime        float64 `json:"predicted_runtime"`
	Confidence              float32 `json:"confidence"`
	InterferenceProbability float32 `json:"interference_probability"`
}

type PredictionRequest struct {
	Job           Job   `json:"job"`
	CoLocatedJobs []Job `json:"co_located_jobs"`
}

type PredictionResponse struct {
	Predictions  []Prediction `json:"predictions"`
	ModelVersion string       `json:"model_version"`
}
