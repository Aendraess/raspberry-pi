package controllers

import (
	"api/database"
	"api/dtos"
	"api/models"
	"api/services"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
)

type ChatController struct {
	chatService *services.ChatService
}

func NewChatController(chatService *services.ChatService) *ChatController {
	return &ChatController{chatService: chatService}
}

func (cc *ChatController) RegisterRoutes(app fiber.Router) {
	app.Post("/chat", cc.Chat)
	threads := app.Group("/chat/threads")
	threads.Get("/", cc.ListThreads)
	threads.Post("/", cc.CreateThread)
	threads.Get("/:id", cc.GetThread)
	threads.Post("/:id/messages", cc.AddMessage)
}

// chatMessagesToOpenAI converts stored messages to OpenAI format (user/assistant only).
func chatMessagesToOpenAI(msgs []models.ChatMessage) []openai.ChatCompletionMessage {
	out := make([]openai.ChatCompletionMessage, 0, len(msgs))
	for _, m := range msgs {
		role := openai.ChatMessageRoleUser
		switch m.Role {
		case "assistant":
			role = openai.ChatMessageRoleAssistant
		case "system":
			role = openai.ChatMessageRoleSystem
		}
		out = append(out, openai.ChatCompletionMessage{Role: role, Content: m.Content})
	}
	return out
}

// Chat sends a single message (no thread) and returns the reply.
// @Summary Chat with LLM using MCP tools (single turn)
// @Description Sends a message to OpenAI with MCP tools; no thread history.
// @Accept json
// @Produce json
// @Tags Chat
// @Param body body dtos.ChatRequest true "Chat message"
// @Success 200 {object} map[string]string "reply"
// @Router /api/chat [post]
func (cc *ChatController) Chat(c *fiber.Ctx) error {
	var req dtos.ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "message is required"})
	}
	ctx, cancel := context.WithTimeout(c.Context(), 60*time.Second)
	defer cancel()
	reply, err := cc.chatService.Chat(ctx, req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"reply": reply})
}

// ListThreads returns all chat threads (id, title, timestamps).
// @Summary List chat threads
// @Produce json
// @Tags Chat
// @Success 200 {array} dtos.ChatThreadResponse
// @Router /api/chat/threads [get]
func (cc *ChatController) ListThreads(c *fiber.Ctx) error {
	var threads []models.ChatThread
	if err := database.DB.Order("updated_at DESC").Find(&threads).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	list := make([]dtos.ChatThreadResponse, 0, len(threads))
	for _, t := range threads {
		list = append(list, dtos.ChatThreadResponse{
			ID:        t.ID,
			Title:     t.Title,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}
	return c.JSON(list)
}

// CreateThread creates a new chat thread.
// @Summary Create a chat thread
// @Accept json
// @Tags Chat
// @Produce json
// @Param body body dtos.CreateThreadRequest true "Optional title"
// @Success 201 {object} dtos.ChatThreadResponse
// @Router /api/chat/threads [post]
func (cc *ChatController) CreateThread(c *fiber.Ctx) error {
	var req dtos.CreateThreadRequest
	_ = c.BodyParser(&req)
	thread := models.ChatThread{Title: req.Title}
	if err := database.DB.Create(&thread).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(dtos.ChatThreadResponse{
		ID:        thread.ID,
		Title:     thread.Title,
		CreatedAt: thread.CreatedAt,
		UpdatedAt: thread.UpdatedAt,
	})
}

// GetThread returns a thread by id with all messages.
// @Summary Get a chat thread with messages
// @Produce json
// @Tags Chat
// @Param id path int true "Thread ID"
// @Success 200 {object} dtos.ChatThreadResponse
// @Router /api/chat/threads/{id} [get]
func (cc *ChatController) GetThread(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid thread id"})
	}
	var thread models.ChatThread
	if err := database.DB.First(&thread, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "thread not found"})
	}
	var messages []models.ChatMessage
	if err := database.DB.Where("thread_id = ?", id).Order("created_at ASC").Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	msgResp := make([]dtos.ChatMessageResponse, 0, len(messages))
	for _, m := range messages {
		msgResp = append(msgResp, dtos.ChatMessageResponse{
			ID:        m.ID,
			Role:      m.Role,
			Content:   m.Content,
			CreatedAt: m.CreatedAt,
		})
	}
	return c.JSON(dtos.ChatThreadResponse{
		ID:        thread.ID,
		Title:     thread.Title,
		CreatedAt: thread.CreatedAt,
		UpdatedAt: thread.UpdatedAt,
		Messages:  msgResp,
	})
}

// AddMessage adds a user message to the thread, gets the assistant reply, persists both, returns the reply.
// @Summary Send a message in a thread and get reply
// @Accept json
// @Produce json
// @Tags Chat
// @Param id path int true "Thread ID"
// @Param body body dtos.AddMessageRequest true "Message"
// @Success 200 {object} dtos.AddMessageResponse
// @Router /api/chat/threads/{id}/messages [post]
func (cc *ChatController) AddMessage(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid thread id"})
	}
	var req dtos.AddMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "message is required"})
	}

	var thread models.ChatThread
	if err := database.DB.First(&thread, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "thread not found"})
	}
	var existing []models.ChatMessage
	if err := database.DB.Where("thread_id = ?", id).Order("created_at ASC").Find(&existing).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	history := chatMessagesToOpenAI(existing)

	ctx, cancel := context.WithTimeout(c.Context(), 60*time.Second)
	defer cancel()
	reply, err := cc.chatService.ChatWithHistory(ctx, history, req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Persist user message and assistant reply
	userMsg := models.ChatMessage{ThreadID: id, Role: "user", Content: req.Message}
	asstMsg := models.ChatMessage{ThreadID: id, Role: "assistant", Content: reply}
	if err := database.DB.Create(&userMsg).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := database.DB.Create(&asstMsg).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Touch thread updated_at
	database.DB.Model(&thread).Update("UpdatedAt", time.Now())

	return c.JSON(dtos.AddMessageResponse{Reply: reply})
}
