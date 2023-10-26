package hash

import (
	"github.com/rs/xid"
)

func GenGUID() string {
	guid := xid.New()
	return guid.String()
}
