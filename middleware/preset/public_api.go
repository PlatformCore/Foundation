package preset

import "github.com/PlatformCore/libpackage/middleware/config"

func PublicAPI() config.Config { c := config.ProductionPreset(); c.Profile = "public_api"; return c }
