package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type Content struct {
	Parts []string `json:Parts`
	Role  string   `json:Role`
}
type Candidates struct {
	Content *Content `json:Content`
}
type ContentResponse struct {
	Candidates *[]Candidates `json:Candidates`
}

type PromptBody struct {
	Prompt string `json:prompt`
}

func welcome(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Welcome to the API",
	})
}

func generateText(c *fiber.Ctx) error {
	API_KEY := os.Getenv("GEMINI_API_KEY")

	var prompt PromptBody
	err := c.BodyParser(&prompt)
	log.Printf("Prompt: %v", prompt)

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(API_KEY))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts
	model := client.GenerativeModel("gemini-1.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(prompt.Prompt))
	if err != nil {
		log.Fatal(err)
	}

	marshalResponse, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(marshalResponse))

	var generateResponse ContentResponse
	if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
		log.Fatal(err)
	}

	//Create a string from the parts of the response
	var data string

	for _, cad := range *generateResponse.Candidates {
		if cad.Content != nil {
			for _, part := range cad.Content.Parts {
				data += part
				fmt.Print(part)
			}
		}
	}

	return c.JSON(fiber.Map{
		"response": data,
	})
}

func summarizeWhatsappConversation(c *fiber.Ctx) error {
	API_KEY := os.Getenv("GEMINI_API_KEY")

	var prompt PromptBody
	err := c.BodyParser(&prompt)
	log.Printf("Prompt: %v", prompt)

	parameter := "I'd like a summary of a WhatsApp conversation. Here's the conversation:" + prompt.Prompt
	parameter = parameter + "In the conversation, please identify: Key topics discussed, Decisions made (if any), Action items (if any), Any disagreements or arguments that arose. Additionally, please note the sentiment of the conversation (positive, negative, neutral). Here are some details to consider including in the summary: Names of participants (if appropriate), Important dates or deadlines mentioned, Links or attachments shared (mention the type of content), For disagreements or arguments, please summarize: The main points of contention, Who was involved, (Optional) The outcome of the disagreement (resolved, unresolved). Please keep the summary concise while maintaining clarity. NOTE: If you see multiple names, it is a group chat. If you see only a single name, and the other name is blank, it is a private chat, and the blank named texts are from me"

	log.Printf("Parameter: %v", parameter)

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(API_KEY))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts
	model := client.GenerativeModel("gemini-1.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(parameter))
	if err != nil {
		log.Fatal(err)
	}

	marshalResponse, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(marshalResponse))

	var generateResponse ContentResponse
	if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
		log.Fatal(err)
	}

	//Create a string from the parts of the response
	var data string

	for _, cad := range *generateResponse.Candidates {
		if cad.Content != nil {
			for _, part := range cad.Content.Parts {
				data += part
				fmt.Print(part)
			}
		}
	}

	return c.JSON(fiber.Map{
		"response": data,
	})
}

func main() {
	app := fiber.New()

	// Load the .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.Get("/api", welcome)
	app.Post("/gemini/generate-text", generateText)
	app.Post("/gemini/summarize-whatsapp-conversation", summarizeWhatsappConversation)

	log.Fatal(app.Listen(":8000"))
}
