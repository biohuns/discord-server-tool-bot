package util

import "fmt"

// ServerStatusText サーバーステータスのテキスト形式
func ServerStatusText(
	isOnline bool,
	gameName string,
	playerCount int,
	maxPlayerCount int,
	mapText string,
) string {
	statusText := "ONLINE"
	if !isOnline {
		statusText = "OFFLINE"
		gameName = "None"
		mapText = "None"
	}
	return fmt.Sprintf(
		"```[%s]: %s\n  Player: %d / %d\n  Map:    %s```",
		gameName,
		statusText,
		playerCount,
		maxPlayerCount,
		mapText,
	)
}

// InstanceStatusText インスタンスステータスのテキスト形式
func InstanceStatusText(name, status string) string {
	return fmt.Sprintf("```<%s: %s>```", name, status)
}
