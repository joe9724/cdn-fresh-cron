package model

type PurgeOkModel struct{
	Result []ResultModel `json:"result"`
}
type ResultModel struct{
	Url string `json:"url"`
	Status string `json:"status"`
	TastID string `json:"task_id"`
	Code int64 `json:"code"`
	EventID string `json:"event_id"`
}

type ErrorResult struct{
    Type string `json:"type"`
    ErrorCode string `json:"error_code"`
    Request string `json:"request"`
    Field string `json:"field"`
    Message string `json:"message"`
}
type ProduceModel struct{
	ID int64 `json:"id"`
	EventID string `json:"event_id"`
	Url string `json:"url"`
	Status int64 `json:"status"`
	TryTimes int64 `json:"try_times"`
	TaskID string `json:"task_id"`
	Message string `json:"message"`
}
type QueryModel struct{
	Result []QueryResult `json:"result"`
}

type QueryResult struct{
	Url string `json:"url"`
	Progress int64 `json:"progress"`
	TaskID string `json:"task_id"`
}