package handler

import (
	"sdc/lib"
	"sdc/model"
	"strconv"
)

type OpenQuestionHandler struct {
	lib.AuthHandler
}

func (self *OpenQuestionHandler) Get() {
	qid, _ := self.GetInt(":qid")
	uid, _ := self.GetSession("userid").(int64)
	role, _ := self.GetSession("userrole").(int64)
	if qid > 0 && uid > 0 {
		if err := model.SetCtypeforQuestion(qid, uid, role, 1); err == nil {
			self.Redirect("/"+strconv.Itoa(int(qid))+"/", 302)
		} else {
			self.Redirect("/", 302)
		}
	} else {
		self.Redirect("/", 302)
	}
}
