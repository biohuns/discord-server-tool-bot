package util

import "fmt"

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

func InstanceStatusText(instanceName, instanceStatus string) string {
	return fmt.Sprintf("```<%s: %s>```", instanceName, instanceStatus)
}
