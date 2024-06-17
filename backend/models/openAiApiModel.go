package models

type OpenAiRequest struct {
	Message string `json:"message" validate:"required"`
}

type OpenAiResponseMessageContent struct {
	Content string `json:"content"`
}

type OpenAiResponseMessage struct {
	Message OpenAiResponseMessageContent `json:"message"`
}

type OpenAiResponse struct {
	Choices       []OpenAiResponseMessage `json:"choices"`
	Metric_Values map[string]string       `json:"metric_values"`
	Error         map[string]string       `json:"error"`
}
