package handler

import (
	"fmt"
	"github.com/astaxie/beego"
	"sdc/lib"
	"sdc/model"
	"strconv"
	"time"
)

type DeleteQuestionHandler struct {
	lib.AuthHandler
}

func (self *DeleteQuestionHandler) Get() {

	flash := beego.NewFlash()
	qid, _ := self.GetInt(":qid")
	uid, _ := self.GetSession("userid").(int64)
	role, _ := self.GetSession("userrole").(int64)
	if qid > 0 {

		if e := model.DelQuestion(qid, uid, role); e != nil {

			self.TplNames = "error.html"
			flash.Error("删除 Topic id:" + strconv.Itoa(int(qid)) + "出现错误 " + fmt.Sprintf("%s", e) + "!")
			flash.Store(&self.Controller)

			return
		}
	}
	self.Redirect("/?ver="+strconv.Itoa(int(time.Now().UnixNano())), 302)
}
