// Package llm is the canonical provider port contract for the hawk ecosystem.
//
// It is the single source of truth for the conversation DTOs and the Provider
// interface that hawk (product face) and eyrie (provider engine) speak across
// their boundary. Both sides alias to these types, so there is exactly one
// definition of eachDTO and no per-call conversion.
//
// hawk owns the product vocabulary (hence names like EyrieMessage); eyrie
// implements the port. eyrie's internal transport types stay eyrie-scoped and
// never appear here.
package llm

import "context"

// ContentPart is a provider-neutral multimodal message part.
type ContentPart struct {
	Type       string          `json:"type"`
	Text       string          `json:"text,omitempty"`
	ImageURL   *ImageURLPart   `json:"image_url,omitempty"`
	InputAudio *InputAudioPart `json:"input_audio,omitempty"`
}

// ImageURLPart describes an image URL or data URI.
type ImageURLPart struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// InputAudioPart describes base64-encoded audio content.
type InputAudioPart struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

// EyrieMessage is the provider-neutral conversation message shape.
type EyrieMessage struct {
	Role         string        `json:"role"`
	Content      string        `json:"content,omitempty"`
	Thinking     string        `json:"thinking,omitempty"`
	ContentParts []ContentPart `json:"content_parts,omitempty"`
	Images       []string      `json:"images,omitempty"`
	ToolUse      []ToolCall    `json:"tool_use,omitempty"`
	ToolResults  []ToolResult  `json:"tool_results,omitempty"`
}

// ToolCall is a tool invocation.
type ToolCall struct {
	ID        string                 `json:"id,omitempty"`
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult is a tool execution result.
type ToolResult struct {
	ToolUseID string `json:"tool_use_id"`
	Content   string `json:"content"`
	IsError   bool   `json:"is_error,omitempty"`
}

// Tool is a tool definition.
type EyrieTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ResponseFormat specifies the desired output format for a model response.
type ResponseFormat struct {
	Type   string `json:"type"`
	Schema string `json:"schema,omitempty"`
}

// ToolChoiceOption controls how the model uses tools.
type ToolChoiceOption struct {
	Type                   string `json:"type"`
	Name                   string `json:"name,omitempty"`
	DisableParallelToolUse bool   `json:"disable_parallel_tool_use,omitempty"`
}

// ChatOptions holds request options for an engine chat call.
type ChatOptions struct {
	Provider             string            `json:"provider,omitempty"`
	Model                string            `json:"model,omitempty"`
	Temperature          *float64          `json:"temperature,omitempty"`
	MaxTokens            int               `json:"max_tokens,omitempty"`
	Stream               bool              `json:"stream,omitempty"`
	Tools                []EyrieTool        `json:"tools,omitempty"`
	System               string            `json:"system,omitempty"`
	EnableCaching        bool              `json:"enable_caching,omitempty"`
	ResponseFormat       *ResponseFormat   `json:"response_format,omitempty"`
	ReasoningEffort      string            `json:"reasoning_effort,omitempty"`
	ThinkingBudgetTokens int               `json:"thinking_budget_tokens,omitempty"`
	ThinkingMode         string            `json:"thinking_mode,omitempty"`
	ThinkingDisplay      string            `json:"thinking_display,omitempty"`
	ThinkingEnabled      *bool             `json:"thinking_enabled,omitempty"`
	VirtualKeyID         string            `json:"virtual_key_id,omitempty"`
	KimiContextCacheID   string            `json:"kimi_context_cache_id,omitempty"`
	KimiCacheResetTTL    bool              `json:"kimi_cache_reset_ttl,omitempty"`
	TopP                 *float64          `json:"top_p,omitempty"`
	TopK                 *int              `json:"top_k,omitempty"`
	StopSequences        []string          `json:"stop_sequences,omitempty"`
	ToolChoice           *ToolChoiceOption `json:"tool_choice,omitempty"`
	MetadataUserID       string            `json:"metadata_user_id,omitempty"`
	ServiceTier          string            `json:"service_tier,omitempty"`
	OutputEffort         string            `json:"output_effort,omitempty"`
	OutputSchema         string            `json:"output_schema,omitempty"`
	PresencePenalty      *float64          `json:"presence_penalty,omitempty"`
	FrequencyPenalty     *float64          `json:"frequency_penalty,omitempty"`
	N                    *int              `json:"n,omitempty"`
	LogProbs             *bool             `json:"logprobs,omitempty"`
	TopLogProbs          *int              `json:"top_logprobs,omitempty"`
	Seed                 *int              `json:"seed,omitempty"`
	Store                *bool             `json:"store,omitempty"`
	Metadata             map[string]string `json:"metadata,omitempty"`
	Modalities           []string          `json:"modalities,omitempty"`
	AudioConfig          string            `json:"audio_config,omitempty"`
	Prediction           string            `json:"prediction,omitempty"`
	WebSearchOptions     string            `json:"web_search_options,omitempty"`
}

// ContinuationConfig controls output continuation behavior.
type ContinuationConfig struct {
	MaxContinuations int
	MaxTotalTokens   int
}

// Usage tracks token usage.
type EyrieUsage struct {
	PromptTokens        int `json:"prompt_tokens"`
	CompletionTokens    int `json:"completion_tokens"`
	TotalTokens         int `json:"total_tokens"`
	CacheCreationTokens int `json:"cache_creation_tokens,omitempty"`
	CacheReadTokens     int `json:"cache_read_tokens,omitempty"`
	ThinkingTokens      int `json:"thinking_tokens,omitempty"`
}

// ResolvedRoute is the concrete provider/model route selected by the engine.
type ResolvedRoute struct {
	Provider          string `json:"provider"`
	Model             string `json:"model"`
	DeploymentRouting bool   `json:"deployment_routing,omitempty"`
}

// Response is the chat response DTO.
type EyrieResponse struct {
	Content        string         `json:"content"`
	Thinking       string         `json:"thinking,omitempty"`
	Usage          *EyrieUsage    `json:"usage,omitempty"`
	ToolCalls      []ToolCall     `json:"tool_calls,omitempty"`
	FinishReason   string         `json:"finish_reason"`
	RequestID      string         `json:"request_id,omitempty"`
	OrganizationID string         `json:"organization_id,omitempty"`
	Route          *ResolvedRoute `json:"route,omitempty"`
}

// EyrieStreamEvent is a streaming event.
type EyrieStreamEvent struct {
	Type       string         `json:"type"`
	Content    string         `json:"content,omitempty"`
	ToolCall   *ToolCall      `json:"tool_call,omitempty"`
	Thinking   string         `json:"thinking,omitempty"`
	Error      string         `json:"error,omitempty"`
	RequestID  string         `json:"request_id,omitempty"`
	Usage      *EyrieUsage    `json:"usage,omitempty"`
	StopReason string         `json:"stop_reason,omitempty"`
	TTFTms     int            `json:"ttft_ms,omitempty"`
	TTFT       int            `json:"ttft,omitempty"`
	Route      *ResolvedRoute `json:"route,omitempty"`
}

// StreamResult wraps a streaming response with cleanup. Callers must call Close()
// when done reading events, or cancel the context.
type StreamResult struct {
	Events    <-chan EyrieStreamEvent
	RequestID string
	cancel    context.CancelFunc
}

// NewStreamResult constructs a stream result. The cancel function is optional
// and must be idempotent.
func NewStreamResult(events <-chan EyrieStreamEvent, requestID string, cancel context.CancelFunc) *StreamResult {
	return &StreamResult{Events: events, RequestID: requestID, cancel: cancel}
}

// Close stops the stream and releases resources.
func (sr *StreamResult) Close() {
	if sr != nil && sr.cancel != nil {
		sr.cancel()
	}
}
