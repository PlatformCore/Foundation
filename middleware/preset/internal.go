package preset

import "github.com/PlatformCore/libpackage/middleware/config"

func Internal() config.Config { c := config.ProductionPreset(); c.Profile = "internal"; return c }
