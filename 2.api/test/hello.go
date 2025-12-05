package test

import (
	request "batchLog/0.core/commonResReq/req"
)

func Hello(ctx request.RequestContext) {
	// logafa.Debug("say hello", "user", "unknown")
	// time.Sleep(10 * time.Second)
	ctx.Success("Helllllllo")
}
