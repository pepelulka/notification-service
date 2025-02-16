package models

type Group struct {
	GroupId      int        `json:"group_id"`
	Name         NullString `json:"name"`
	Participants []Person   `json:"participants"`
}

type GroupCreate struct {
	Name           NullString `json:"name"`
	ParticipantIds []int      `json:"participant_ids"`
}
