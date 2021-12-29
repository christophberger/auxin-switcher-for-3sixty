package hifiberry

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

const (
	pac     = "/proc/asound/cards"
	stat    = "/proc/asound/card%d/pcm0p/sub0/status"
	berryRe = `(\d).*hifiberry`
)

// getCardNumber determines the card number of the hifiberry sound card.
func getCardNumber() (int, error) {
	cards, err := ioutil.ReadFile(pac)
	if err != nil {
		return -1, fmt.Errorf("getCardNumber: cannot read %s: %w", pac, err)
	}
	re := regexp.MustCompile(berryRe)
	matches := re.FindStringSubmatch(string(cards))
	if matches == nil || len(matches) < 2 {
		return -1, fmt.Errorf("getCardNumber: cannot find hifiberry card")
	}
	// parse string to int
	card, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, fmt.Errorf("getCardNumber: cannot parse card number: %w", err)
	}
	return card, nil

}

// getStatus reads the status of the hifiberry sound card.
// 0 = idle, 1 = playing
func GetSoundStatus() (bool, error) {
	num, err := getCardNumber()
	if err != nil {
		return false, fmt.Errorf("GetCardStatus: cannot get card number: %w", err)
	}
	status, err := ioutil.ReadFile(fmt.Sprintf(stat, num))
	if err != nil {
		return false, fmt.Errorf("getCardStatus: cannot read %s: %w", stat, err)
	}
	idx := strings.Index(string(status), "RUNNING")
	return idx > -1, nil
}
