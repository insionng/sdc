package handler

import (
	"fmt"
	"html/template"
	"sdc/helper"
	"sdc/lib"
	"sdc/model"
	"strconv"
	"strings"
)

type MainHandler struct {
	lib.BaseHandler
}

func (self *MainHandler) Get() {
	//fmt.Println("im MainHandler")
	self.Data["catpage"] = "home"
	self.TplNames = "sdc/main.html"

	ipage, _ := self.GetInt(":page")
	page := int(ipage)

	tab := template.HTMLEscapeString(strings.TrimSpace(self.GetString(":tab")))

	url := "/"
	if tab == "lastest" {
		url = "/lastest/"
		tab = "id"
		self.Data["tab"] = "lastest"
	} else if tab == "hotness" {
		url = "/hotness/"
		tab = "hotness"
		self.Data["tab"] = "hotness"
	} else if tab == "unanswered" {
		url = "/unanswered/"
		tab = "unanswered"
		self.Data["tab"] = "unanswered"
	} else {
		url = "/lastest/"
		tab = "id"
		self.Data["tab"] = "lastest"
	}

	pagesize := 30
	results_count, err := model.GetQuestionsCount(0, pagesize, tab)
	if err != nil {
		return
	}
	pages, page, beginnum, endnum, offset := helper.Pages(int(results_count), page, pagesize)
	/*
	   <article class="post" id="q-1010000000444736">
	       <div class="status">
	           <span class="answer answered" title="18 个答案">18</span>
	           <span class="vote voted" title="3 个投票">3</span>
	       </div>
	       <div class="p-summary">
	           <a class="author" data-toggle="tooltip" data-placement="bottom" rel="tooltip" href="http://segmentfault.com/u/tonyx" title="TonyX &bull; 658">
	               <img class="avatar-40" src="http://sfault-avatar.b0.upaiyun.com/262/496/2624967969-1030000000395062_medium40" alt="TonyX" />
	           </a>
	           <h2>
	               <a href="http://segmentfault.com/q/1010000000444736" title="你最喜欢的开发工具是什么？">你最喜欢的开发工具是什么？</a>
	           </h2>
	           <div class="meta">
	               <span class="views"> <i class="i-view"></i>
	                   912 次浏览
	               </span>
	               &nbsp;
	               <ul class="meta-tags">
	                   <li> <i class="i-tag"></i>
	                   </li>
	                   <li>
	                       <a data-tid="1040000000090473" href="/t/ide">ide</a>
	                   </li>
	               </ul>
	               &nbsp;
	               <span class="datetime">
	                   <i class="i-time"></i>
	                   17分钟前
	               </span>
	           </div>
	       </div>
	   </article>
	*/
	if qts, err := model.GetQuestions(offset, pagesize, tab); err == nil {
		results_count := len(*qts)
		if results_count > 0 {
			i := 1
			output := ""
			for _, v := range *qts {

				i += i

				output = output + `<article class="post" id="` + strconv.Itoa(int(v.Id)) + `"><div class="status">`

				if v.Ctype == 2 {
					output = output + `<span class="answer answered-accepted" title="` + strconv.Itoa(int(v.ReplyCount)) + ` 个答案">` + strconv.Itoa(int(v.ReplyCount)) + `</span>`
				} else if v.Ctype == -1 {
					output = output + `<span class="answer answered closed" title="` + strconv.Itoa(int(v.ReplyCount)) + ` 个答案">` + strconv.Itoa(int(v.ReplyCount)) + `</span>`
				} else if v.ReplyCount == 0 {
					output = output + `<span class="answer" title="0 个答案">0</span>`
				} else {
					output = output + `<span class="answer answered" title="` + strconv.Itoa(int(v.ReplyCount)) + ` 个答案">` + strconv.Itoa(int(v.ReplyCount)) + `</span>`
				}

				if v.Hotvote == 0 {
					output = output + `<span class="vote" title="0 个投票">0</span>`
				} else {
					output = output + `<span class="vote voted" title="` + strconv.Itoa(int(v.Hotvote)) + ` 个投票">` + strconv.Itoa(int(v.Hotvote)) + `</span>`
				}
				output = output + `</div><div class="p-summary">`
				avatar := ""
				if usr, err := model.GetUser(v.Uid); err == nil && usr != nil {
					avatar = helper.Gravatar(usr.Email, 40)
				}
				if avatar == "" {
					output = output + `<a class="author" data-toggle="tooltip" data-placement="bottom" rel="tooltip" href="###" title="&bull; ` + v.Author + ` &bull;">
	               <img class="avatar-40" src="http://sfault-avatar.b0.upaiyun.com/262/496/2624967969-1030000000395062_medium40" alt="` + v.Author + `" />
	           </a>`
				} else {
					output = output + `<a class="author" data-toggle="tooltip" data-placement="bottom" rel="tooltip" href="###/user/` + v.Author + `/" title="&bull; ` + v.Author + ` &bull;">
	               <img class="avatar-40" src="` + avatar + `" alt="` + v.Author + `" />
	           </a>`
				}

				output = output + `
	           <h2>
	               <a href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `">` + v.Title + `</a>
	           </h2>
	           <div class="meta">
	               <span class="views"> <i class="i-view"></i>
	                   ` + strconv.Itoa(int(v.Views)) + ` 次浏览
	               </span>
	               &nbsp;`

				if v.Tags != "" {
					output = output + `<ul class="meta-tags"><li><i class="i-tag"></i></li>`
					tags := helper.Tags(v.Tags, ",")
					itags := len(tags)

					for _, tag := range tags {
						itags = itags - 1
						if itags == 0 {
							//等于0则不添加逗号
							output = output + `<li><a data-tid="` + strconv.Itoa(int(v.Id)) + `" href="/tag/` + tag + `/">` + tag + `</a></li>`

						} else {
							output = output + `<li><a data-tid="` + strconv.Itoa(int(v.Id)) + `" href="/tag/` + tag + `/">` + tag + `</a> , </li>`

						}
					}
					output = output + `</ul>`
				}
				output = output + `&nbsp;<span class="datetime"><i class="i-time"></i>` + " " + v.Author + " " + helper.TimeSince(v.Created) + `提问` + `</span></div></div></article>`
			}
			self.Data["questions"] = output
			self.Data["pagesbar"] = helper.Pagesbar(url, "", results_count, pages, page, beginnum, endnum, 0)
		}

	} else {
		fmt.Println("首页推荐榜单 数据查询出错", err)
	}

	self.Data["replys_8s"] = model.GetAnswersByPid(0, 1, 0, 8, "id")

	//侧栏九宫格推荐榜单
	//先行取出最热门的9个节点 然后根据节点获取该节点下最热门的话题
	/*
		if nd, err := model.GetNodes(0, 9, "hotness"); err == nil {
			if len(*nd) > 0 {
				for _, v := range *nd {

					i := 0
					output_start := `<ul class="widgets-popular widgets-similar clx">`
					output := ""
					if qts := model.GetTopicsByNid(v.Id, 0, 1, 0, "hotness"); err == nil {

						if len(*qts) > 0 {
							for _, v := range *qts {

								i += 1
								if i == 3 {
									output = output + `<li class="similar similar-third">`
									i = 0
								} else {
									output = output + `<li class="similar">`
								}
								output = output + `<a target="_blank" href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `">
													<img src="` + v.ThumbnailsSmall + `" wdith="70" height="70" />
												</a>
											</li>`
							}
						}
					}
					output_end := "</ul>"
					if len(output) > 0 {
						output = output_start + output + output_end
						self.Data["topic_hotness_9_module"] = template.HTML(output)
					} else {
						self.Data["topic_hotness_9_module"] = nil
					}

				}
			}
		} else {
			fmt.Println("节点数据查询出错", err)
		}
	*/

	//侧栏九宫格推荐榜单
	//根据最热的1个分类查找该分类下级9个话题
	/*
		if cats, err := model.GetCategorys(0, 1, "hotness"); err == nil && len(cats) > 0 {
			for _, v := range cats {
				output_start := `<ul class="widgets-popular widgets-similar clx">`
				output := ""
				i := 0
				//根据CID 获取同分类下的最热的9个话题
				if qts := model.GetTopicsByCid(v.Id, 0, 9, 0, "hotness"); len(*qts) > 0 {

					for _, v := range *qts {

						i += 1
						if i == 3 {
							output = output + `<li class="similar similar-third">`
							i = 0
						} else {
							output = output + `<li class="similar">`
						}
						output = output + `<a target="_blank" href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `">
													<img src="` + v.ThumbnailsSmall + `" wdith="70" height="70" />
												</a>
											</li>`
					}
				}
				output_end := "</ul>"

				if len(output) > 0 {
					output = output_start + output + output_end
					self.Data["topic_hotness_9_module"] = template.HTML(output)
				} else {
					self.Data["topic_hotness_9_module"] = nil
				}
			}

		}
	*/

	//根据最热的1个节点查找该节点的上级分类下的9个话题
	/*
		if nd, err := model.GetNodes(0, 1, "hotness"); err == nil && len(*nd) > 0 {
			for _, v := range *nd {

				output_start := `<ul class="widgets-popular widgets-similar clx">`
				output := ""
				i := 0
				//根据节点的上级PID 获取同分类下的最热话题
				if qts := model.GetQuestionsByCid(v.Pid, 0, 9, 0, "hotness"); len(*qts) > 0 {

					for _, v := range *qts {

						i += 1
						if i == 3 {
							output = output + `<li class="similar similar-third">`
							i = 0
						} else {
							output = output + `<li class="similar">`
						}
						output = output + `<a target="_blank" href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `">
														<img src="` + v.ThumbnailsSmall + `" wdith="70" height="70" />
													</a>
												</li>`
					}
				}
				output_end := "</ul>"

				if len(output) > 0 {
					output = output_start + output + output_end
					self.Data["topic_hotness_9_module"] = template.HTML(output)
				} else {
					self.Data["topic_hotness_9_module"] = nil
				}

			}
		}
	*/
}
