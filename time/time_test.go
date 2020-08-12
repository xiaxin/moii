package time

import (
	"github.com/xiaxin/moii/log"
	"testing"
)

func TestTimeScore(t *testing.T) {
	var score int64 = 10

	val := TimeScoreEncode(score)

	if score != TimeScoreDecode(val) {
		log.Error("val error")
	}
}
