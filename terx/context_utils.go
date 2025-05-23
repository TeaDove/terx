package terx

import (
	"fmt"
	"strings"
)

func extractCommandAndText(text string, botUsername string, isChat bool) (string, string) {
	if len(text) <= 1 || text[0] != '/' || strings.HasPrefix(text, "/@") {
		return "", text
	}

	spaceIdx := strings.Index(text, " ")

	atIdx := strings.Index(text, "@")
	if atIdx == -1 && isChat {
		return "", text
	}

	if atIdx != -1 && (spaceIdx == -1 || spaceIdx > atIdx) { //nolint: nestif // too lazy to fix
		var extractedUsername string
		if spaceIdx == -1 {
			extractedUsername = text[atIdx:]
		} else {
			extractedUsername = text[atIdx:spaceIdx]
		}

		if extractedUsername == fmt.Sprintf("@%s", botUsername) {
			if spaceIdx == -1 {
				return text[1:atIdx], ""
			}

			return text[1:atIdx], text[spaceIdx+1:]
		}

		return "", text
	}

	if spaceIdx == -1 {
		return text[1:], ""
	}

	return text[1:spaceIdx], text[spaceIdx+1:]
}
