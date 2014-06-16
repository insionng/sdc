package handler

import (
	"sdc/lib"
	"sdc/model"
	"strconv"
)

type DeleteAnswerHandler struct {
	lib.AuthHandler
}

func (self *DeleteAnswerHandler) Get() {
	aid, _ := self.GetInt(":aid")
	qid, _ := self.GetInt(":qid")
	uid, _ := self.GetSession("userid").(int64)
	role, _ := self.GetSession("userrole").(int64)

	if aid > 0 && qid > 0 {
		if model.DelAnswer(aid, uid, role) == nil {

			if affected, err := model.SetReplyCountByPid(qid); err == nil && affected != 0 {

				self.Redirect("/"+strconv.Itoa(int(qid))+"/", 302)
			} else {
				self.Redirect("/", 302)

			}
		} else {
			self.Redirect("/", 302)
		}
	} else {
		self.Redirect("/", 302)
	}
}
