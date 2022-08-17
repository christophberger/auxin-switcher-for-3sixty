package systemctl

import (
	"context"
	"fmt"
	"time"

	sc "github.com/taigrr/systemctl"
)

func Restart(unit string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := sc.Options{UserMode: false}
	err := sc.Restart(ctx, unit, opts)
	if err != nil {
		return fmt.Errorf("systemctl.Restart: failed to restart %s: %w", unit, err)
	}
	return nil
}
