package evaluator

import "github.com/PlatformCore/Foundation/platform/featureflag/types"

func Enabled(flag types.Flag, _ types.Target) bool {
	return flag.Enabled
}


