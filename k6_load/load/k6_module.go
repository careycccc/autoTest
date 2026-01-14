// load/k6_module.go
package load

import "go.k6.io/k6/js/modules"

func init() {
	// 注册模块名，后面 JS 用这个 import
	// 建议用你项目名，比如 github.com/xxx/xxx/load → k6/x/yourload
	modules.Register("k6/x/load", new(Load))
}

type Load struct{}
