package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketMessageType string

const (
	WebsocketMessageTypeStart      WebsocketMessageType = "start"
	WebsocketMessageTypeInfo       WebsocketMessageType = "info"
	WebsocketMessageTypeWarning    WebsocketMessageType = "warning"
	WebsocketMessageTypeError      WebsocketMessageType = "error"
	WebsocketMessageTypeUpdate     WebsocketMessageType = "update"
	WebsocketMessageTypeLastUpdate WebsocketMessageType = "last_update"
)

type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type WebSocketMessage struct {
	Type    WebsocketMessageType `json:"type"`
	Message string               `json:"message,omitempty"`
	Payload any                  `json:"payload,omitempty"`
}

func (e *ErrorResponse) Error() string {
	if e.Code != "" && e.Message != "" {
		return e.Code + ": " + e.Message
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Code != "" {
		return e.Code
	}
	return "unknown error"
}

const (
	APISuccessCode    = "00000"
	APISuccessMessage = "success"
)

type SuccessResp struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result,omitempty"`
}

func SendApiResponseV1(ctx *gin.Context, apiResp interface{}, appErr *ApplicationError) {
	if appErr != nil {
		ctx.JSON(appErr.HttpCode, &ErrorResponse{
			Code:    string(appErr.ErrorCode),
			Message: appErr.ErrorMessage,
		})
		return
	}

	if apiResp != nil {
		ctx.JSON(http.StatusOK, apiResp)
		return
	}

	ctx.JSON(http.StatusOK, SuccessResp{
		Code:    APISuccessCode,
		Message: APISuccessMessage,
	})
}

func SendWebSocketMessage(conn *websocket.Conn, msgType WebsocketMessageType, message string, payload any) {
	wsMsg := WebSocketMessage{
		Type:    msgType,
		Message: message,
		Payload: payload,
	}
	if err := conn.WriteJSON(wsMsg); err != nil {
		fmt.Println("[SendWebSocketMessage] Error sending message:", err)
	}
}
