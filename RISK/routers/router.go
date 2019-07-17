package routers

import (
	"github.com/oganator/RISK/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	beego.Router("/Metrics", &controllers.MetricsController{})
}
