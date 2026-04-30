package preset

import "github.com/PlatformCore/libpackage/middleware/config"

func Production() config.Config { return config.ProductionPreset() }
