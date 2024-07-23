package failover

import (
	"GkWeiBook/webook/internal/service/sms"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestFailOverSMSService_Send(t *testing.T) {

	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) []sms.Service
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewFailOverSMSService(tt.mock(ctrl))

			err := svc.Send(nil, "", nil, "")
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
