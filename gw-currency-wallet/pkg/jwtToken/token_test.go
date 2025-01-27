package jwttoken

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	accessSecret = "H0VkSLIXCblklfssEpG7mpeCoEUY2nNLLN2qAB6310VzXTHlxeHp9+cZcwraiwBMwVWnx17C2yGd1yctgYT8kg=="

	refreshSecret = "F6gB9bpL4Lr5KQXVRWmd1NjaOtSqH3u3D1b7xxokUE4h8lSzBKBebB8Haulbx2GWRu/Zt8C0w8jADoOkju9HNQ=="
)

func TestValidate(t *testing.T) {
	tokenOpts := &TokensOption{
		UserID:        1,
		AccessExp:     time.Minute * 5,
		RefreshExp:    time.Hour * 24,
		SecretAccess:  accessSecret,
		SecretRefresh: refreshSecret,
	}

	access, refresh, err := GenerateTokens(tokenOpts)

	assert.Nil(t, err)

	err = Validate(access, []byte(accessSecret))
	assert.Nil(t, err)

	err = Validate(refresh, []byte(refreshSecret))
	assert.Nil(t, err)
}
