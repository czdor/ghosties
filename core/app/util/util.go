package util

import (
	"core/internal/lib"
	"core/types"
	"fmt"
	"math/big"

	"crypto/rand"
)

func Contains(target []string, value string) bool {
	for _, v := range target {
		if v == value {
			return true
		}
	}

	return false
}

func Delete(target []string, value string) []string {
	for idx, v := range target {
		if v == value {
			return append(target[:idx], target[:idx+1]...)
		}
	}

	return target
}

// Function to convert UsersPositions map to a slice of Position
func ConvertMapToSlice(usersPositions []string) []lib.Position {
	positions := []lib.Position{}

	for _, key := range usersPositions {
		var row, col int
		fmt.Sscanf(key, "%d,%d", &row, &col)
		positions = append(positions, lib.Position{Row: row, Col: col})
	}

	return positions
}

func GetRandomId() (string, error) {
	rId, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("error generating random id: %v", err)
	}

	return rId.String(), nil
}

func NewRoomId(roomName string) (*types.RoomId, error) {
	randomId, err := GetRandomId()
	if err != nil {
		return nil, fmt.Errorf("error generating random id: %v", err)
	}

	id := fmt.Sprintf(types.RoomIdFormat, roomName, randomId)
	roomId := types.RoomId(id)

	return &roomId, nil
}
