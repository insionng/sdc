package handler

import (
	//"fmt"
	//"html/template"
	"sdc/helper"
	"sdc/lib"
	"sdc/model"
	//"strconv"
)

type QuestionHandler struct {
	lib.BaseHandler
}

func (self *QuestionHandler) Get() {
	self.TplNames = "sdc/question.html"

	qid, _ := self.GetInt(":qid")

	if qid > 0 {

		if qs, err := model.GetQuestion(qid); qs != nil && err == nil {
			avatar := ""
			if usr, err := model.GetUser(qs.Uid); err == nil && usr != nil {
				avatar = helper.Gravatar(usr.Email, 32)
			}
			self.Data["avatar"] = avatar
			self.Data["article"] = *qs
			self.Data["replys"] = *model.GetAnswersByPid(qid, 0, 0, 0, "hotness")

			/*
				if qss := model.GetTopicsByCidOnBetween(qs.Cid, qid-5, qid+5, 0, 11, 0, "asc"); qss != nil && qid != 0 && len(qss) > 0 {

					for i, v := range qss {

						if v.Id == qid {
							//两侧的翻页按钮参数 初始化 s
							prev := i - 1
							next := i + 1
							//两侧的翻页按钮参数 初始化 e

							//话题内容部位 页码  初始化 s
							ipagesbar := `<div class="link_pages">`
							h := 5
							ipagesbar_start := i - h
							ipagesbar_end := i + h
							j := 0
							//话题内容部位 页码  初始化 e
							for i, v := range qss {
								//两侧的翻页按钮 s
								if prev == i {
									self.Data["previd"] = v.Id
									self.Data["prev"] = v.Title
								}
								if next == i {
									self.Data["nexqid"] = v.Id
									self.Data["next"] = v.Title
								}
								//两侧的翻页按钮 e

								//话题内容部位 页码 s
								if ipagesbar_start == i {
									ipagesbar = ipagesbar + `<a href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `"><span><</span></a>`
								}
								if i > ipagesbar_start && i < ipagesbar_end {
									j += 1
									if v.Id == qid { // current

										ipagesbar = ipagesbar + `<span>` + strconv.Itoa(int(v.Id)) + `</span>`

									} else { //loop

										ipagesbar = ipagesbar + `<a href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `"><span>` + strconv.Itoa(int(v.Id)) + `</span></a>`

									}
									if j > (2 * h) {
										break
									}
								}
								if ipagesbar_end == i {
									ipagesbar = ipagesbar + `<a href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `"><span>></span></a>`
								}
								//话题内容部位 页码 e
							}
							self.Data["ipagesbar"] = template.HTML(ipagesbar + "</div>")
						}
					}
					// (qss []*Topic)

				}
			*/

			//侧栏你可能喜欢 推荐同节点下的最热话题
			/*
				if qss := model.GetTopicsByNid(qs.Nid, 0, 6, 0, "hotness"); *qss != nil {

					if len(*qss) > 0 {
						i := 0
						ouqsut := `<ul class="widgets-similar clx">`
						for _, v := range *qss {

							i += 1
							if i == 3 {
								ouqsut = ouqsut + `<li class="similar similar-third">`
								i = 0
							} else {
								ouqsut = ouqsut + `<li class="similar">`
							}
							ouqsut = ouqsut + `<a target="_blank" href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `">
													<img src="` + v.ThumbnailsSmall + `" wdith="70" height="70" />
												</a>
											</li>`
						}
						ouqsut = ouqsut + `</ul>`
						self.Data["topic_sidebar_hotness_6_module"] = template.HTML(ouqsut)
					}

				}
			*/

			//侧栏九宫格推荐榜单
			//先行取出最热门的一个节点 然后根据节点获取该节点下最热门的话题
			/*
				if nd, err := model.GetNodes(0, 1, "hotness"); err == nil {
					if len(*nd) == 1 {
						for _, v := range *nd {

							if qss := model.GetTopicsByNid(v.Id, 0, 9, 0, "hotness"); err == nil {

								if len(*qss) > 0 {
									i := 0
									ouqsut := `<ul class="widgets-popular widgets-similar clx">`
									for _, v := range *qss {

										i += 1
										if i == 3 {
											ouqsut = ouqsut + `<li class="similar similar-third">`
											i = 0
										} else {
											ouqsut = ouqsut + `<li class="similar">`
										}
										ouqsut = ouqsut + `<a target="_blank" href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `">
													<img src="` + v.ThumbnailsSmall + `" wdith="70" height="70" />
												</a>
											</li>`
									}
									ouqsut = ouqsut + `</ul>`
									self.Data["topic_hotness_9_module"] = template.HTML(ouqsut)
								}
							} else {
								fmt.Println("推荐榜单(9)数据查询出错", err)
							}
						}
					}
				} else {
					fmt.Println("节点数据查询出错", err)
				}
			*/

			//底部六格推荐
			//推荐同一作者的最热话题
			/*
				if qss := model.GetTopicsByUid(qs.Uid, 0, 6, 0, "hotness"); len(*qss) > 0 {
					i := 0
					ouqsut := `<ul class="widgets-similar clx">`
					for _, v := range *qss {

						i += 1
						if i == 3 {
							ouqsut = ouqsut + `<li class="likesimilar likesimilar-3">`
							i = 0
						} else {
							ouqsut = ouqsut + `<li class="likesimilar">`
						}
						ouqsut = ouqsut + `<a target="_blank" href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `" alt="` + v.Title + `" class="likeimglink">
													<img src="` + v.ThumbnailsMedium + `" wdith="150" height="150" />
													<span class="bg">` + v.Title + `</span>
												</a>
											</li>`
					}
					ouqsut = ouqsut + `</ul>`
					self.Data["topic_hotness_6_module"] = template.HTML(ouqsut)
				}
			*/
		} else {
			self.Redirect("/", 302)
		}
	} else {
		self.Redirect("/", 302)
	}

}
