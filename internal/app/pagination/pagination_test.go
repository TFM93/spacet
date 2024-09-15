package pagination

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_DecodeCursor(t *testing.T) {
	type args struct {
		encodedCursor string
	}
	tests := []struct {
		name    string
		args    args
		wantTs  string
		wantID  string
		wantErr error
	}{
		{
			name: "success",
			args: args{
				encodedCursor: "MjAyNC0wOC0yMlQyMDowOToxMS45MzgyMiswMTowMHxjMTJlMjNmMy1mNWUzLTQxYmMtYWVjYS05ZDY2YmQwYjk2YTM=",
			},
			wantTs:  "2024-08-22 20:09:11.93822 +0100 WEST",
			wantID:  "c12e23f3-f5e3-41bc-aeca-9d66bd0b96a3",
			wantErr: nil,
		},
		{
			name: "invalid cursor",
			args: args{
				encodedCursor: "MjAyNC0wDowOToxMS45MzgyMiswMTowMHxjMTJlMjNmMy1mNWUzLTQxYmMtYWVjYS05ZDY2YmQwYjk2YTM=",
			},
			wantErr: fmt.Errorf("illegal base64 data at input byte 83"),
		},
		{
			name: "invalid cursor content",
			args: args{
				encodedCursor: base64.StdEncoding.EncodeToString([]byte("a|b|c")),
			},
			wantErr: fmt.Errorf("cursor is invalid"),
		},
		{
			name: "invalid timestamp",
			args: args{
				encodedCursor: base64.StdEncoding.EncodeToString([]byte("a|b")),
			},
			wantErr: fmt.Errorf("cursor is invalid: timestamp"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTs, gotID, err := DecodeCursor(tt.args.encodedCursor)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}

			if gotTs.String() != tt.wantTs {
				t.Errorf("DecodeCursor() gotTs = %v, want %v", gotTs, tt.wantTs)
			}
			if gotID != tt.wantID {
				t.Errorf("DecodeCursor() gotID = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func Test_encodeCursor(t *testing.T) {
	type args struct {
		ts time.Time
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				ts: time.Date(2024, 8, 22, 20, 30, 3, 0, time.UTC),
				id: "abc",
			},
			want: "MjAyNC0wOC0yMlQyMDozMDowM1p8YWJj",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeCursor(tt.args.ts, tt.args.id); got != tt.want {
				t.Errorf("EncodeCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}
