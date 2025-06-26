package contracts

type CreateRoomCommand struct {
	Username string `json:"username" validate:"required"`
}
type CreateRoomResponse struct {
	RoomId string         `json:"room_id"`
	Player PlayerResponse `json:"player"`
}

type JoinRoomCommand struct {
	Username string `json:"username" validate:"required"`
	RoomId   string
}

type JoinRoomResponse struct {
	RoomId string         `json:"room_id"`
	Player PlayerResponse `json:"player"`
}

type SetCombinationCommand struct {
	PlayerId    string `json:"player_id" validate:"required"`
	Combination int    `json:"combination" validate:"required"`
	RoomId      string
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
