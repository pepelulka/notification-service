package models

type Group struct {
	GroupId        int    `json:"group_id"`
	Name           string `json:"name"`
	ParticipantIds []int  `json:"participant_ids"`
}

type GroupCreate struct {
	Name           string `json:"name"`
	ParticipantIds []int  `json:"participant_ids"`
}
