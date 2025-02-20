package models

import rs "notification_service/internal/rabbit_senders"

type TgSendRequestToSpecificUsernamesBody rs.TelegramSend

type TgSendRequestToSpecificGroupsBody struct {
	Content    string   `json:"content"`
	GroupNames []string `json:"group_names"`
}
