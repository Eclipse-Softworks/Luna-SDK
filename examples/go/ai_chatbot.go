// Package main demonstrates AI chatbot functionality using Luna SDK.
//
// Run with: go run ai_chatbot.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/eclipse-softworks/Luna-SDK-go/luna"
)

// ChatMessage represents a message in the conversation
type ChatMessage struct {
	Role    string // "system", "user", "assistant"
	Content string
}

// LunaChatbot is a conversational AI assistant
type LunaChatbot struct {
	client              *luna.Client
	model               string
	temperature         float64
	systemPrompt        string
	conversationHistory []ChatMessage
}

// NewLunaChatbot creates a new chatbot instance
func NewLunaChatbot(client *luna.Client, opts ...ChatbotOption) *LunaChatbot {
	bot := &LunaChatbot{
		client:       client,
		model:        "luna-gpt-4",
		temperature:  0.7,
		systemPrompt: "You are a helpful assistant powered by Luna SDK. Be concise and helpful.",
	}

	for _, opt := range opts {
		opt(bot)
	}

	// Initialize with system prompt
	bot.conversationHistory = []ChatMessage{
		{Role: "system", Content: bot.systemPrompt},
	}

	return bot
}

// ChatbotOption is a functional option for configuring the chatbot
type ChatbotOption func(*LunaChatbot)

// WithModel sets the model
func WithModel(model string) ChatbotOption {
	return func(bot *LunaChatbot) {
		bot.model = model
	}
}

// WithTemperature sets the temperature
func WithTemperature(temp float64) ChatbotOption {
	return func(bot *LunaChatbot) {
		bot.temperature = temp
	}
}

// WithSystemPrompt sets the system prompt
func WithSystemPrompt(prompt string) ChatbotOption {
	return func(bot *LunaChatbot) {
		bot.systemPrompt = prompt
	}
}

// Chat sends a message and gets a response
func (bot *LunaChatbot) Chat(ctx context.Context, userMessage string) (string, error) {
	// Add user message to history
	bot.conversationHistory = append(bot.conversationHistory, ChatMessage{
		Role:    "user",
		Content: userMessage,
	})

	// Build messages for API
	messages := make([]luna.Message, len(bot.conversationHistory))
	for i, msg := range bot.conversationHistory {
		messages[i] = luna.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Call the AI API
	response, err := bot.client.AI().ChatCompletions(ctx, luna.CompletionRequest{
		Model:       bot.model,
		Messages:    messages,
		Temperature: bot.temperature,
	})
	if err != nil {
		return "", fmt.Errorf("chat completion failed: %w", err)
	}

	// Extract assistant response
	assistantMessage := response.Choices[0].Message.Content

	// Add to history for context
	bot.conversationHistory = append(bot.conversationHistory, ChatMessage{
		Role:    "assistant",
		Content: assistantMessage,
	})

	return assistantMessage, nil
}

// ClearHistory clears the conversation history
func (bot *LunaChatbot) ClearHistory() {
	bot.conversationHistory = []ChatMessage{
		{Role: "system", Content: bot.systemPrompt},
	}
}

// GetHistory returns the conversation history
func (bot *LunaChatbot) GetHistory() []ChatMessage {
	return append([]ChatMessage{}, bot.conversationHistory...)
}

func main() {
	client := luna.NewClient(
		luna.WithAPIKey(os.Getenv("LUNA_API_KEY")),
	)

	ctx := context.Background()

	// Run examples
	if err := simpleQAExample(ctx, client); err != nil {
		fmt.Printf("Simple Q&A failed: %v\n", err)
	}

	if err := conversationExample(ctx, client); err != nil {
		fmt.Printf("Conversation example failed: %v\n", err)
	}

	if err := specializedAssistantExample(ctx, client); err != nil {
		fmt.Printf("Specialized assistant failed: %v\n", err)
	}

	// Uncomment to start interactive mode
	// interactiveChat(ctx, client)
}

// ============================================
// Simple Q&A Example
// ============================================

func simpleQAExample(ctx context.Context, client *luna.Client) error {
	fmt.Println("Simple Q&A Example\n")

	response, err := client.AI().ChatCompletions(ctx, luna.CompletionRequest{
		Model: "luna-gpt-4",
		Messages: []luna.Message{
			{Role: "system", Content: "You are a helpful coding assistant."},
			{Role: "user", Content: "What is a Go goroutine?"},
		},
		Temperature: 0.5,
	})
	if err != nil {
		return fmt.Errorf("chat completion failed: %w", err)
	}

	fmt.Println("Question: What is a Go goroutine?")
	fmt.Printf("\nAnswer: %s\n", response.Choices[0].Message.Content)

	return nil
}

// ============================================
// Multi-turn Conversation Example
// ============================================

func conversationExample(ctx context.Context, client *luna.Client) error {
	fmt.Println("\nMulti-turn Conversation Example\n")

	chatbot := NewLunaChatbot(client,
		WithSystemPrompt("You are a friendly coding tutor. Explain concepts simply."),
		WithTemperature(0.6),
	)

	// First message
	fmt.Println("User: Explain what an API is")
	response, err := chatbot.Chat(ctx, "Explain what an API is")
	if err != nil {
		return err
	}
	fmt.Printf("Assistant: %s\n", response)

	// Follow-up question (chatbot remembers context)
	fmt.Println("\nUser: Can you give me an example?")
	response, err = chatbot.Chat(ctx, "Can you give me an example?")
	if err != nil {
		return err
	}
	fmt.Printf("Assistant: %s\n", response)

	// Another follow-up
	fmt.Println("\nUser: How do REST APIs differ?")
	response, err = chatbot.Chat(ctx, "How do REST APIs differ?")
	if err != nil {
		return err
	}
	fmt.Printf("Assistant: %s\n", response)

	return nil
}

// ============================================
// Specialized Assistant Example
// ============================================

func specializedAssistantExample(ctx context.Context, client *luna.Client) error {
	fmt.Println("\nCode Review Assistant Example\n")

	codeReviewer := NewLunaChatbot(client,
		WithModel("luna-gpt-4"),
		WithTemperature(0.3), // Lower temperature for more focused responses
		WithSystemPrompt(`You are an expert code reviewer. 
			Analyze code for:
			- Bugs and potential issues
			- Performance improvements
			- Best practices
			- Security concerns
			Be specific and constructive.`),
	)

	codeToReview := `
func fetchUserData(userID string) interface{} {
	var data interface{}
	resp, _ := http.Get("/api/users/" + userID)
	json.NewDecoder(resp.Body).Decode(&data)
	return data
}
`

	fmt.Println("Code to review:")
	fmt.Println(codeToReview)

	review, err := codeReviewer.Chat(ctx, "Please review this Go function:\n"+codeToReview)
	if err != nil {
		return err
	}

	fmt.Printf("\nCode Review:\n%s\n", review)

	return nil
}

// ============================================
// Interactive Chat
// ============================================

func interactiveChat(ctx context.Context, client *luna.Client) {
	fmt.Println("\nInteractive Chat Mode\n")
	fmt.Println("Type your messages below. Type 'quit' to exit, 'clear' to reset.\n")

	chatbot := NewLunaChatbot(client,
		WithSystemPrompt("You are a helpful AI assistant. Be friendly and informative."),
	)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if strings.ToLower(input) == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		if strings.ToLower(input) == "clear" {
			chatbot.ClearHistory()
			fmt.Println("Conversation cleared.\n")
			continue
		}

		if input == "" {
			continue
		}

		response, err := chatbot.Chat(ctx, input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("\nAssistant: %s\n\n", response)
	}
}
