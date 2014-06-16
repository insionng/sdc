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

type SearchHandler struct {
	lib.BaseHandler
}

func (self *SearchHandler) Get() {
	self.TplNames = "sdc/search.html"
	keyword := template.HTMLEscapeString(strings.TrimSpace(self.GetString(":keyword")))

	if keyword == "" {
		keyword = template.HTMLEscapeString(strings.TrimSpace(self.GetString("keyword")))
	}

	ipage, _ := self.GetInt(":page")
	if ipage <= 0 {
		ipage, _ = self.GetInt("page")
	}
	page := int(ipage)

	limit := 9 //每页显示数目

	//如果已经登录登录
	sess_username, _ := self.GetSession("username").(string)
	if sess_username != "" {
		limit = 30
	}

	if keyword != "" {
		if qts, err := model.SearchQuestion(keyword, 0, 0, "id"); err == nil {

			rcs := len(*qts)
			pages, pageout, beginnum, endnum, offset := helper.Pages(rcs, page, limit)

			if st, err := model.SearchQuestion(keyword, offset, limit, "hotness"); err == nil {
				results_count := len(*st)
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

				}

			}

			if k := self.GetString("keyword"); k != "" {
				self.Data["search_keyword"] = k
			} else {
				self.Data["search_keyword"] = keyword
			}

			self.Data["pagesbar"] = helper.Pagesbar("/search/", keyword, rcs, pages, pageout, beginnum, endnum, 0)
		} else {
			fmt.Println("SearchQuestion errors:", err)
			return
		}

		self.Data["replys_8s"] = model.GetAnswersByPid(0, 1, 0, 8, "id")

		//侧栏九宫格推荐榜单
		//根据用户的关键词推荐
		/*
			nds, ndserr := model.SearchNode(keyword, 0, 9, "hotness")
			cats, catserr := model.SearchCategory(keyword, 0, 9, "hotness")
			//如果在节点找到关键词
			if (ndserr == catserr) && (ndserr == nil) {

				output_start := `<ul class="widgets-popular widgets-similar clx">`
				output := ""
				i := 0
				if len(*nds) >= len(*cats) && len(*nds) > 0 {
					for _, v := range *nds {

						if tps := model.GetTopicsByNid(v.Id, 0, 1, 0, "hotness"); len(*tps) > 0 {

							for _, v := range *tps {

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
				} else if len(*cats) > len(*nds) && len(*cats) > 0 {
					for _, v := range *cats {

						if tps := model.GetTopicsByCid(v.Id, 0, 1, 0, "hotness"); len(*tps) > 0 {

							for _, v := range *tps {

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
				}
				output_end := "</ul>"
				if len(output) > 0 {
					output = output_start + output + output_end
					self.Data["topic_hotness_9_module"] = template.HTML(output)
				} else {
					self.Data["topic_hotness_9_module"] = nil
				}
			}
		*/

		//侧栏九宫格推荐榜单
		//先行取出最热门的9个节点 然后根据节点获取该节点下最热门的话题
		/*
			if nd, err := model.GetNodes(0, 9, "hotness"); err == nil {
				if len(*nd) > 0 {
					for _, v := range *nd {

						i := 0
						output_start := `<ul class="widgets-popular widgets-similar clx">`
						output := ""
						if tps := model.GetTopicsByNid(v.Id, 0, 1, 0, "hotness"); err == nil {

							if len(*tps) > 0 {
								for _, v := range *tps {

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
	} else {
		self.Redirect("/", 302)
	}

}
