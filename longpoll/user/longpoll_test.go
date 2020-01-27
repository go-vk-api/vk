package longpoll

import (
	"testing"
)

func TestNewWithOptions(t *testing.T) {
	wantMode := ReceiveAttachments + ReturnFriendOnlineExtraField

	lp, err := NewWithOptions(
		nil,
		WithMode(wantMode),
	)
	if err != nil {
		t.Error(err)
	}

	if lp.Wait != DefaultWait {
		t.Errorf("lp.Wait == %d, want %d", lp.Wait, DefaultWait)
	}

	if lp.Mode != wantMode {
		t.Errorf("lp.Mode == %d, want %d", lp.Mode, wantMode)
	}

	if lp.Version != DefaultVersion {
		t.Errorf("lp.Version == %d, want %d", lp.Version, DefaultVersion)
	}

	if lp.client != nil {
		t.Errorf("lp.client = %v, want nil", lp.client)
	}
}
