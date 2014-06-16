package handler

import (
	"fmt"
	"github.com/astaxie/beego"
	"html/template"
	"sdc/helper"
	"sdc/lib"
	"sdc/model"
	"strconv"
	"strings"
	"time"
)

type EditAnswerHandler struct {
	lib.AuthHandler
}

func (self *EditAnswerHandler) Get() {
	self.TplNames = "sdc/edit-answer.html"
	flash := beego.NewFlash()

	aid, _ := self.GetInt(":aid")

	if aid_handler, err := model.GetAnswer(aid); err == nil && aid_handler != nil {
		uid, _ := self.GetSession("userid").(int64)
		role, _ := self.GetSession("userrole").(int64)
		allow := bool(false)

		if aid_handler.Uid == uid && aid_handler.Id == aid {
			allow = true
		} else if role < 0 {
			allow = true
		}

		if allow {
			self.Data["answer"] = *aid_handler
		} else {
			//没有权限执行该操作则直接跳转到登录页面
			self.Redirect("/user/signin/", 302)
		}

	} else {

		flash.Error(fmt.Sprint(err))
		flash.Store(&self.Controller)
		return
	}
}

func (self *EditAnswerHandler) Post() {
	self.TplNames = "sdc/edit-answer.html"

	flash := beego.NewFlash()

	aid, _ := self.GetInt(":aid")

	if aid_handler, err := model.GetAnswer(aid); err == nil {
		uid, _ := self.GetSession("userid").(int64)
		role, _ := self.GetSession("userrole").(int64)
		allow := bool(false)

		if aid_handler.Uid == uid && aid_handler.Id == aid {
			allow = true
		} else if role < 0 {
			allow = true
		}

		if allow {
			self.Data["answer"] = *aid_handler
			aid_content := template.HTMLEscapeString(strings.TrimSpace(self.GetString("content")))

			if aid_content != "" {

				if anz, err := model.GetAnswer(aid); anz != nil && err == nil {

					//删去用户没再使用的图片
					helper.DelLostImages(anz.Content, aid_content)
					anz.Content = aid_content

					if s, e := helper.GetBannerThumbnail(aid_content); e == nil {
						anz.Attachment = s
					}

					/*
						if cat, err := model.GetCategory(nd.Pid); err == nil {
							anz.Category = cat.Title
						}
					*/

					anz.Updated = time.Now()

					if row, err := model.PutAnswer(aid, anz); row == 1 && err == nil {
						model.SetRecordforImageOnEdit(aid, anz.Uid)
						self.Redirect("/"+strconv.Itoa(int(anz.Pid))+"/#answer-"+strconv.Itoa(int(aid)), 302)
					} else {

						flash.Error("更新答案出现错误:", fmt.Sprint(err))
						flash.Store(&self.Controller)
						return
					}
				} else {
					flash.Error("无法获取根本不存在的答案!")
					flash.Store(&self.Controller)
					return
				}

			} else {

				flash.Error("答案内容为空!")
				flash.Store(&self.Controller)
				return
			}
		} else {
			//没有权限执行该操作则直接跳转到登录页面
			self.Redirect("/user/signin/", 302)
		}

	} else {

		flash.Error(fmt.Sprint(err))
		flash.Store(&self.Controller)
		return
	}

}
