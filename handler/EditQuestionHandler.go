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

type EditQuestionHandler struct {
	lib.AuthHandler
}

func (self *EditQuestionHandler) Get() {
	self.TplNames = "sdc/edit-question.html"
	flash := beego.NewFlash()

	qid, _ := self.GetInt(":qid")

	if qid_handler, err := model.GetQuestion(qid); err == nil && qid_handler != nil {
		uid, _ := self.GetSession("userid").(int64)
		role, _ := self.GetSession("userrole").(int64)
		allow := bool(false)

		if qid_handler.Uid == uid && qid_handler.Id == qid {
			allow = true
		} else if role < 0 {
			allow = true
		}

		if allow {

			self.Data["question"] = *qid_handler
			self.Data["inode"], _ = model.GetNode(qid_handler.Nid)
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

func (self *EditQuestionHandler) Post() {
	self.TplNames = "sdc/edit-question.html"

	flash := beego.NewFlash()
	tags := template.HTMLEscapeString(strings.TrimSpace(strings.ToLower(self.GetString("tags"))))

	qid, _ := self.GetInt(":qid")

	if qid_handler, err := model.GetQuestion(qid); err == nil {
		uid, _ := self.GetSession("userid").(int64)
		role, _ := self.GetSession("userrole").(int64)
		allow := bool(false)

		if qid_handler.Uid == uid && qid_handler.Id == qid {
			allow = true
		} else if role < 0 {
			allow = true
		}

		if allow {

			self.Data["question"] = *qid_handler
			self.Data["inode"], _ = model.GetNode(qid_handler.Nid)

			if tags == "" {

				flash.Error("尚未设置标签,请设定正确的标签!")
				flash.Store(&self.Controller)
				return
			} else {
				qid_title := template.HTMLEscapeString(strings.TrimSpace(self.GetString("title")))
				qid_content := template.HTMLEscapeString(strings.TrimSpace(self.GetString("content")))

				if qid_title != "" && qid_content != "" {
					tags := template.HTMLEscapeString(strings.TrimSpace(strings.ToLower(self.GetString("tags"))))

					if qs, err := model.GetQuestion(qid); qs != nil && err == nil {

						qs.Title = qid_title

						//删去用户没再使用的图片
						helper.DelLostImages(qs.Content, qid_content)
						qs.Content = qid_content

						if s, e := helper.GetBannerThumbnail(qid_content); e == nil {
							qs.Attachment = s
						}

						if thumbnails, thumbnailslarge, thumbnailsmedium, thumbnailssmall, e := helper.GetThumbnails(qid_content); e == nil {
							qs.Thumbnails = thumbnails
							qs.ThumbnailsLarge = thumbnailslarge
							qs.ThumbnailsMedium = thumbnailsmedium
							qs.ThumbnailsSmall = thumbnailssmall
						}
						/*
							if cat, err := model.GetCategory(nd.Pid); err == nil {
								qs.Category = cat.Title
							}
						*/

						qs.Tags = tags
						qs.Updated = time.Now()

						if row, err := model.PutQuestion(qid, qs); row == 1 && err == nil {
							model.SetRecordforImageOnEdit(qid, qs.Uid)
							self.Redirect("/"+strconv.Itoa(int(qid))+"/", 302)
						} else {

							flash.Error("更新问题出现错误:", fmt.Sprint(err))
							flash.Store(&self.Controller)
							return
						}
					} else {
						flash.Error("无法获取根本不存在的问题!")
						flash.Store(&self.Controller)
						return
					}

				} else {

					flash.Error("问题标题或内容为空!")
					flash.Store(&self.Controller)
					return
				}
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
