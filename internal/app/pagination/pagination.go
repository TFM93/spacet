package pagination

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// DecodeCursor takes a base64 encoded string and returns the timestamp and id
func DecodeCursor(encodedCursor string) (ts time.Time, id string, err error) {
	cursor, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return
	}

	splitted := strings.Split(string(cursor), "|")
	if len(splitted) != 2 {
		err = fmt.Errorf("cursor is invalid")
		return
	}

	ts, err = time.Parse(time.RFC3339Nano, splitted[0])
	if err != nil {
		err = fmt.Errorf("cursor is invalid: timestamp")
		return
	}
	id = splitted[1]
	return
}

// EncodeCursor returns a base64 encoded string based on the provided timestamp and id
func EncodeCursor(ts time.Time, id string) string {
	key := fmt.Sprintf("%s|%s", ts.Format(time.RFC3339Nano), id)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
