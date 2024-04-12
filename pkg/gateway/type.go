package gateway

type RequestBody struct {
	Model          string            `json:"model"`
	ResponseFormat map[string]string `json:"response_format"`
	Messages       []Message         `json:"messages"`
}

type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     any    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

type Response struct {
	Choices []Choice `json:"choices"`
	Created int      `json:"created"`
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Logprobs     any     `json:"logprobs"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
