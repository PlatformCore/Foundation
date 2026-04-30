package errorx

import "github.com/PlatformCore/libpackage/transport/core"

func Write(ctx *core.Context, err error) error {
	if ctx == nil {
		return err
	}
	e := Mapper{DefaultStatus: 500}.Map(err)
	ctx.Result = core.Failure(e)
	if ctx.Response == nil {
		ctx.Response = core.NewResponse()
	}
	ctx.Response.Status(e.Status)
	_ = ctx.Response.JSON(map[string]any{"ok": false, "code": e.Code, "message": e.Message})
	return err
}
