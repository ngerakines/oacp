package client

import (
	"fmt"
	"github.com/urfave/cli"
)

func userAgent(cliCtx *cli.Context) string {
	return fmt.Sprintf("aocp-client/%s (%s)", cliCtx.App.Version, cliCtx.App.Compiled)
}
