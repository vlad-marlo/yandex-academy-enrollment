package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseWorkingHours(t *testing.T) {
	var sh, eh, sm, em uint8
	for sh = 0; sh <= maxHourValue+1; sh += 3 {
		for eh = 0; eh <= maxHourValue+1; eh += 3 {
			for sm = 0; sm <= maxMinuteValue+1; sm += 3 {
				for em = 0; sm <= maxMinuteValue+1; em += 3 {
					h, err := ParseWorkingHours(fmt.Sprintf("%d:%d-%d:%d", sh, sm, eh, em))
					if err == nil {
						if assert.NotNil(t, h) {
							assert.Equal(t, eh, h.end.hour)
							assert.Equal(t, em, h.end.minute)
							assert.Equal(t, sm, h.start.minute)
							assert.Equal(t, sh, h.start.hour)
						}
					} else {
						assert.ErrorIs(t, err, ErrBadWorkingHours)
					}
				}

			}
		}
	}
	h, err := ParseWorkingHours("2:22-22:21")
	if assert.NoError(t, err) && assert.NotNil(t, h) {
		assert.Equal(t, uint8(22), h.end.hour)
		assert.Equal(t, uint8(21), h.end.minute)
		assert.Equal(t, uint8(2), h.start.hour)
		assert.Equal(t, uint8(22), h.start.minute)
	}
}
