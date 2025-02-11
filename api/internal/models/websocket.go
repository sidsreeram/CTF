package models

type WebSocketService interface {
    BroadcastScoreUpdate(challengeID, score int)
    BroadcastEvent(event string, data interface{})
}
