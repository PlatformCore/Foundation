package errorx

import "github.com/PlatformCore/libpackage/transport/core"

type Mapper struct{ DefaultStatus int }

func (m Mapper) Map(err error) *core.Error {
	e := core.ToError(err)
	if e.Status == 0 {
		if m.DefaultStatus != 0 {
			e.Status = m.DefaultStatus
		} else {
			e.Status = 500
		}
	}
	return e
}
