package deadletter

import (
	"context"

	"github.com/PlatformCore/libpackage/platform/eventbus/subscriber"
)

type Writer interface {
	Write(context.Context, subscriber.Message, error) error
}


