package notifier

import (
	"testing"
	"github.com/joaquinicolas/Elca/Novelty"
)

func TestNotifyNovelty(t *testing.T) {
	novelty := &Novelty.Novelty{
		Text:"ALARM: MONITOR      MODULE ADDR 1M004   Z004   12:33A 081115 1M004 SIGNAL SILENCED   12:33A 081115 Tue SYSTEM RESET  12:34A081115 Tue ",

	}
	err := NotifyNovelty(novelty)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotifyAlive(t *testing.T) {
	err := NotifyAlive()
	if err != nil {
		t.Fatal(err)
	}

}