package api

import (
	"net/http"
	"strconv"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	querier repo.Querier
}

func NewMessageHandler(querier repo.Querier) *MessageHandler {
	return &MessageHandler{
		querier: querier,
	}
}

func (h *MessageHandler) WireHttpHandler() http.Handler {

	r := gin.Default()
	r.Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.String(http.StatusInternalServerError, "Internal Server Error: panic")
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.POST("/message", h.handleCreateMessage)
	r.POST("/thread", h.handleCreateThread)
	r.GET("/message/:id", h.handleGetMessage)
	r.GET("/thread/:id/messages", h.handleGetThreadMessages)
	r.PATCH("/message/", h.updateMessage)
	r.DELETE("/message/:id", h.deleteMessage)

	return r
}

func (h *MessageHandler) updateMessage(c *gin.Context) {
	var msgBody repo.UpdateMesageByIDParams
	err := c.ShouldBindBodyWithJSON(&msgBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// msgID := c.Param("id")
	// if msgID == "" {
	// 	c.JSON(http.StatusBadRequest, "invalid message id")
	// 	return
	// }

	message, err := h.querier.UpdateMesageByID(c, msgBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, message)

}

type threadParams struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

func (h *MessageHandler) handleCreateThread(c *gin.Context) {
	var req threadParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.querier.CreateThread(c, req.Topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) handleCreateMessage(c *gin.Context) {
	var req repo.CreateMessageParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.querier.CreateMessage(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) handleGetMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	message, err := h.querier.GetMessageByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) handleGetThreadMessages(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// Get pagination parameters from query string
	page, err := strconv.Atoi(c.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10")) // Default to 10 items per page
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page size"})
		return
	}

	// Calculate offset
	offset := int32((page - 1) * pageSize)

	// Prepare parameters for GetMessagesByThreadPaginated
	params := repo.GetMessagesByThreadPaginatedParams{
		ThreadID: id,
		Limit:    int32(pageSize),
		Offset:   offset,
	}

	// Fetch paginated messages
	messages, err := h.querier.GetMessagesByThreadPaginated(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch total count of messages for the thread
	totalCount, err := h.querier.GetTotalMessageCountByThread(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate total pages
	totalPages := (int(totalCount) + pageSize - 1) / pageSize

	// Respond with paginated data
	c.JSON(http.StatusOK, gin.H{
		"thread":   id,
		"topic":    "example",
		"messages": messages,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
			"total_count": totalCount,
		},
	})
}

func (h *MessageHandler) deleteMessage(c *gin.Context) {
	id := c.Param("id")
	// err := c.ShouldBindBodyWithJSON(&req)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.querier.DeleteMeageByID(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "message deleted successfully"})
}
