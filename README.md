# ⛲️ Watermill Test Logger

Logger for testing with watermill.

## ⬇️ Installation

`go get github.com/survivorbat/go-watermill-test-logger`

## 📋 Usage

```go
package main

import (
 "github.com/stretchr/testify/require"
 "github.com/ThreeDotsLabs/watermill/message"
 "github.com/survivorbat/go-watermill-test-logger"
)

func TestWatermill(t *testing.T) {
  // This Will fail the test if any error occurred, and log to the test instance
  logger := watertestlogger.NewTestAdaptor(t, true, watermill.LogLevelInfo)

  router, err := message.NewRouter(message.RouterConfig{}, logger)
  require.NoError(t, err)

  // [...]
}
```

## 🔭 Plans

Not much yet.
