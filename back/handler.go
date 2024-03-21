package back

import (
	"github.com/gin-gonic/gin"
	"github.com/orcastor/orcas/rpc/util"
)

func list(ctx *gin.Context) {
	var req struct {
	}
	ctx.BindJSON(&req)
	util.Response(ctx, gin.H{})
}
