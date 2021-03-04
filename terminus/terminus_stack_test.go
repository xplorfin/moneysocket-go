package terminus

import (
	"testing"

	"github.com/xplorfin/moneysocket-go/moneysocket/config"
)

func TestTerminusStack(t *testing.T) {
	configuration := config.NewConfig()
	stack := NewTerminusStack(configuration)
	_ = stack
}
