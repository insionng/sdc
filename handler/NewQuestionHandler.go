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

type NewQuestionHandler struct {
	lib.AuthHandler
}

func (self *NewQuestionHandler) Get() {
	self.TplNames = "sdc/new-question.html"
}

func (self *NewQuestionHandler) Post() {
	self.TplNames = "sdc/new-question.html"

	flash := beego.NewFlash()
	tags := template.HTMLEscapeString(strings.TrimSpace(strings.ToLower(self.GetString("tags"))))

	if tags == "" {

		flash.Error("尚未设置标签,请设定正确的标签!")
		flash.Store(&self.Controller)
		return
	} else {

		uid, _ := self.GetSession("userid").(int64)
		sess_username, _ := self.GetSession("username").(string)
		qid_title := template.HTMLEscapeString(strings.TrimSpace(self.GetString("title")))
		qid_content := template.HTMLEscapeString(strings.TrimSpace(self.GetString("content")))

		if qid_title != "" && qid_content != "" {

			qs := new(model.Question)
			qs.Title = qid_title
			qs.Tags = tags
			qs.Content = qid_content
			qs.Uid = uid
			qs.Author = sess_username
			qs.Created = time.Now()
			qs.Updated = qs.Created

			if s, e := helper.GetBannerThumbnail(qid_content); e == nil {
				qs.Attachment = s
			}

			if thumbnails, thumbnailslarge, thumbnailsmedium, thumbnailssmall, e := helper.GetThumbnails(qid_content); e == nil {
				qs.Thumbnails = thumbnails
				qs.ThumbnailsLarge = thumbnailslarge
				qs.ThumbnailsMedium = thumbnailsmedium
				qs.ThumbnailsSmall = thumbnailssmall
			}

			if qts, _, err := model.PostQuestion(qs); err == nil {
				model.SetRecordforImageOnPost(qts.Id, uid)
				self.Redirect("/"+strconv.Itoa(int(qts.Id))+"/", 302)
			} else {

				flash.Error(fmt.Sprint(err))
				flash.Store(&self.Controller)
				return
			}
		} else {
			flash.Error("问题标题或内容为空!")
			flash.Store(&self.Controller)
			return
		}
	}
}
