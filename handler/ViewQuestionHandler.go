package handler

import (
	"fmt"
	"sdc/lib"
	"sdc/model"
	"strconv"
)

type ViewQuestionHandler struct {
	lib.BaseHandler
}

func (self *ViewQuestionHandler) Get() {
	name := self.GetString(":name")
	id, _ := self.GetInt(":id")

	if name != "" && id > 0 {
		if name == "question" {

			if qs, err := model.GetQuestion(id); qs != nil && err == nil {
				qs.Views = qs.Views + 1
				//qs.Hotup = qs.Hotup + 1
				//qs.Hotscore = helper.Qhot_QScore(qs.Hotup, qs.Hotdown)
				//qs.Hotness = helper.Qhot(qs.Views, qs.ReplyCount, qs.Hotscore, qs.Hotscore, qs.Created, qs.ReplyTime)
				//qs.Hotvote = helper.Qhot_Vote(qs.Hotup, qs.Hotdown)
				if row, e := model.PutQuestion(id, qs); e != nil {
					fmt.Println("ViewQuestionHandler更新话题ID", id, "访问次数数据错误,row:", row, e)
					self.Abort("500")
				} else {
					self.Ctx.Output.Context.WriteString(strconv.Itoa(int(qs.Views)))
				}
			}
		} else {
			self.Abort("501")
		}

	} else {
		self.Abort("501")
	}

}
