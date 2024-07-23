package failover

import (
	"GkWeiBook/webook/internal/service/sms"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/context"
	"testing"
)

func TestTimeoutFailoverSMSService_Send(t1 *testing.T) {

	tests := []struct {
		name string
		mock func(ctrl *gomock.Controller) []sms.Service

		threshold uint64
		idx       uint64
		cnt       uint64

		wantErr error
		wandIdx uint64
		wandCnt uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			ctrl := gomock.NewController(t1)
			defer ctrl.Finish()
			svc := NewTimeoutFailoverSMSService(tt.mock(ctrl), tt.threshold)
			//svc.idx = tt.idx
			// 有争议，搁置
			err := svc.Send(context.Background(), "", nil, "")
			assert.Equal(t1, tt.wantErr, err)
		})
	}
}
