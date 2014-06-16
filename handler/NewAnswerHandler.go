package handler

import (
	"fmt"
	"html/template"
	"sdc/lib"
	"sdc/model"
	"strconv"
	"strings"
)

type NewAnswerHandler struct {
	lib.BaseHandler
}

func (self *NewAnswerHandler) Post() {
	qid, _ := self.GetInt(":qid")
	suid, _ := self.GetSession("userid").(int64)

	author := template.HTMLEscapeString(strings.TrimSpace(strings.ToLower(self.GetString("author"))))
	email := template.HTMLEscapeString(strings.TrimSpace(strings.ToLower(self.GetString("email"))))
	website := template.HTMLEscapeString(strings.TrimSpace(strings.ToLower(self.GetString("website"))))
	rc := template.HTMLEscapeString(strings.TrimSpace(self.GetString("content")))

	//不等于0,即是注册用户或管理层 此时把ctype设置为1 主要是为了区分游客
	if suid != 0 {
		if qid > 0 && rc != "" {

			if usr, err := model.GetUser(suid); err == nil {
				//为安全计,先行保存回应,顺手获得aid,在后面顺手再更新替换@通知的链接
				if aid, err := model.AddAnswer(qid, suid, 1, rc, usr.Username, usr.Content, usr.Email, usr.Website); err != nil {
					fmt.Println("#", aid, ":", err)
				} else {

					//如果回应内容中有@通知 则处理以下事件
					/*
						if users := helper.AtUsers(rc); len(users) > 0 {
							if tp, err := model.GetQuestion(qid); err == nil {
								todo := []string{}
								for _, v := range users {
									//判断被通知之用户名是否真实存在
									if u, e := model.GetUserByUsername(v); e == nil && u != nil {
										//存在的则加入待操作列
										todo = append(todo, v)
										//替换被通知用户的用户名带上用户主页链接
										rc = strings.Replace(rc, "@"+v,
											"<a href='/user/"+u.Username+"/' title='"+u.Nickname+"' target='_blank'><span>@</span><span>"+u.Username+"</span></a>", -1)

										//发送通知内容到用户的 时间线
										model.AddTimeline(usr.Username+"在「"+tp.Title+"」的回应里提到了你~",
											rc+"[<a href='/"+self.GetString(":qid")+"/#answer-"+strconv.Itoa(int(aid))+"'>"+tp.Title+"</a>]",
											tp.Cid, tp.Nid, u.Id, usr.Username, usr.Content)

									}

								}
								if len(todo) > 0 {
									model.SetReplyContentByRid(aid, rc)
								}

							}
						}
					*/
					self.Redirect("/"+self.GetString(":qid")+"/#answer-"+strconv.Itoa(int(aid)), 302)
					return
				}
			}
			self.Redirect("/"+self.GetString(":qid")+"/", 302)
		} else if qid > 0 {
			self.Redirect("/"+self.GetString(":qid")+"/", 302)
		} else {
			self.Redirect("/", 302)
		}
	} else { //游客回应 此时把ctype设置为-1   游客不开放@通知功能
		if author != "" && email != "" && qid > 0 && rc != "" {
			if aid, err := model.AddAnswer(qid, suid, -1, rc, author, "", email, website); err != nil {
				fmt.Println("#", aid, ":", err)
				self.Redirect("/"+self.GetString(":qid")+"/", 302)
			} else {
				self.Redirect("/"+self.GetString(":qid")+"/#answer-"+strconv.Itoa(int(aid)), 302)
			}
		} else if qid > 0 {
			self.Redirect("/"+self.GetString(":qid")+"/", 302)
		} else {
			self.Redirect("/", 302)
		}

	}

}
