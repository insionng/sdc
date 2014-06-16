package handler

import (
	"fmt"
	"sdc/helper"
	"sdc/lib"
	"sdc/model"
	"strconv"
	"time"
)

type LikeHandler struct {
	lib.BaseHandler
}

func (self *LikeHandler) Get() {

	if helper.IsSpider(self.Ctx.Request.UserAgent()) != true {
		name := self.GetString(":name")
		id, _ := self.GetInt(":id")
		uid, _ := self.GetSession("userid").(int64)

		if name == "question" {
			if model.IsQuestionMark(uid, id) {
				//self.Abort("304") <-白痴函数 妈的 难道这货不是用来设置状态号的?居然尼玛的直接panic!
				self.Ctx.Output.SetStatus(304)
				return

			} else {
				if qs, err := model.GetQuestion(id); err == nil {

					qs.Hotup = qs.Hotup + 1
					qs.Hotscore = helper.Qhot_QScore(qs.Hotup, qs.Hotdown)
					qs.Hotvote = helper.Qhot_Vote(qs.Hotup, qs.Hotdown)
					qs.Hotness = helper.Qhot(qs.Views, qs.ReplyCount, qs.Hotscore, model.GetAScoresByPid(id), qs.Created, qs.ReplyTime)

					if _, err := model.PutQuestion(id, qs); err != nil {
						fmt.Println("PutQuestion执行错误:", err)
					} else {
						model.SetQuestionMark(uid, id)
					}
					//&hearts; 有用 ({{.article.Hotup}})
					self.Ctx.WriteString(strconv.Itoa(int(qs.Hotscore)))
				} else {
					return
				}
			}
		} else if name == "answer" {
			if model.IsAnswerMark(uid, id) {
				//self.Abort("304")
				self.Ctx.Output.SetStatus(304)
				return

			} else {
				if ans, err := model.GetAnswer(id); err == nil {

					ans.Hotup = ans.Hotup + 1
					ans.Views = ans.Views + 1
					ans.Hotscore = helper.Qhot_AScore(ans.Hotup, ans.Hotdown)
					ans.Hotvote = helper.Qhot_Vote(ans.Hotup, ans.Hotdown)
					ans.Hotness = helper.Qhot(ans.Views, ans.ReplyCount, ans.Hotscore, ans.Views, ans.Created, ans.ReplyTime)

					if _, err := model.PutAnswer(id, ans); err != nil {
						fmt.Println("PutAnswer执行错误:", err)
					} else {
						model.SetAnswerMark(uid, id)
					}
					self.Ctx.WriteString(strconv.Itoa(int(ans.Hotscore)))
				} else {
					return
				}
			}
		} else if name == "topic" {

			if tp, err := model.GetTopic(id); err == nil {

				tp.Hotup = tp.Hotup + 1
				tp.Hotscore = helper.Hotness_Score(tp.Hotup, tp.Hotdown)
				tp.Hotness = helper.Hotness(tp.Hotup, tp.Hotdown, time.Now())
				model.PutTopic(id, tp)
				//&hearts; 有用 ({{.article.Hotup}})
				self.Ctx.WriteString(strconv.Itoa(int(tp.Hotup)))
			} else {
				return
			}
		} else if name == "node" {

			if nd, err := model.GetNode(id); err == nil {

				nd.Hotup = nd.Hotup + 1
				nd.Hotscore = helper.Hotness_Score(nd.Hotup, nd.Hotdown)
				nd.Hotness = helper.Hotness(nd.Hotup, nd.Hotdown, time.Now())
				model.PutNode(id, nd)
				self.Ctx.WriteString("node liked")
			} else {
				return
			}
		} else {
			//self.Abort("304")
			self.Ctx.Output.SetStatus(304)
		}

	} else {
		//self.Abort("401")
		self.Ctx.Output.SetStatus(401)
	}

}
