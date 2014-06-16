package helper

/*
	helper 模块是纯功能性质 辅助性质的代码
	对数据库直接操作的一切代码都不能写在此
*/

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	fpath "path"
	"regexp"
	"sdc/plugin/blackfriday"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

//分页计算函数
func Pages(results_count int, page int, pagesize int) (pages int, pageout int, beginnum int, endnum int, offset int) {
	//取得记录总数，计算总页数用
	//results_count,总共有results_count条记录

	//设定每一页显示的记录数
	if pagesize < 0 || pagesize < 1 {
		pagesize = 10 //如无设置，则默认每页显示10条记录
	}

	//计算总页数
	pages = int(math.Ceil(float64(results_count) / float64(pagesize)))
	//返回pages

	//判断页数设置,否则，设置为第一页
	if page < 0 || page < 1 {
		page = 1
	}
	if page > pages {
		page = pages
	}
	//返回page

	beginnum = page - 4
	endnum = page + 5

	if page < 5 {
		beginnum = 1
		endnum = 10 //可用链接数，现在是当前页加前后两页共5页，if条件为可用链接数的一半
	}
	if page > pages-5 {
		beginnum = pages - 9
		endnum = pages
	}
	if beginnum < 1 {
		beginnum = 1
	}
	if endnum > pages {
		endnum = pages
	}
	//返回beginnum
	//返回endnum

	//计算记录偏移量
	offset = int((page - 1) * pagesize)
	return int(pages), int(page), int(beginnum), int(endnum), offset
}

func Pagesbar(url string, keyword string, results_max int, pages int, page int, beginnum int, endnum int, style int) (output template.HTML) {
	var raw string
	switch {
	case style == 0: //sdc 定制版pagesbar

		if keyword != "" {
			keyword = keyword + "/"
		}
		nextpage, prevpage := 0, 0
		pagemindle := 2
		pagewidth := pagemindle * 2
		raw = `<div class="page-nav"><ul class="pagination">`
		if results_max > 0 {
			count := int(pages + 1)
			//prev page
			if (page != beginnum) && (page > beginnum) {
				prevpage = page - 1
				raw = raw + `<li><a class="prev" href="` + url + keyword + "page-" + strconv.Itoa(prevpage) + `/">上一页</a></li>`
			}

			//current page and loop pages
			j := 0
			for i := page; i < count; i++ {
				j += 1
				if i == page {
					raw = raw + `<li class="active"><a href="` + url + keyword + "page-" + strconv.Itoa(i) + `/">` + strconv.Itoa(i) + "</a></li>"
				} else {

					raw = raw + `<li><a href="` + url + keyword + "page-" + strconv.Itoa(i) + `/">` + strconv.Itoa(i) + "</a></li>"
				}
				if j > pagewidth {
					break
				}
			}

			raw = raw + "<li><span>共" + strconv.Itoa(int(pages)) + "页</span></li>"

			//next page
			if (page != endnum) && (page < endnum) {
				nextpage = page + 1
				raw = raw + `<li><a class="next" href="` + url + keyword + "page-" + strconv.Itoa(nextpage) + `/">下一页</a></li>`
			}
		} else {
			raw = raw + "<li><span>共0页</span></li>"

		}
		raw = raw + "</ul></div>"
		/*
			if nextpage == 0 && prevpage == 0 {
				output = template.HTML(raw)
			} else {

				if nextpage > 0 || prevpage > 0 {
					raw = raw + `<div id="pagenavi-fixed">`

					if prevpage > 0 {
						raw = raw + `
								<div class="pages-prev">
									<a href="` + url + keyword + `page-` + strconv.Itoa(prevpage) + `/">上一页 &raquo;</a>
								</div>`
					}

					if nextpage > 0 {
						raw = raw + `
								<div class="pages-next">
									<a href="` + url + keyword + `page-` + strconv.Itoa(nextpage) + `/">下一页 &raquo;</a>
								</div>`
					}
					raw = raw + "</div>"
				}
				output = template.HTML(raw)
			}
		*/
		output = template.HTML(raw)
	case style == 1:
		/*
			<div class="pagination"><ul>
					<li><a href="#">&laquo;</a></li>
					<li class="active"><a href="#">1</a></li>
					<li><a href="#">2</a></li>
					<li><a href="#">3</a></li>
					<li><a href="#">4</a></li>
					<li><a href="#">&raquo;</a></li>
			</ul></div>
		*/

		/*
			<ul class="pager">
			  <li class="previous">
			    <a href="#">&larr; Older</a>
			  </li>
			  <li class="next">
			    <a href="#">Newer &rarr;</a>
			  </li>
			</ul>
		*/
		if results_max > 0 {
			raw = "<ul class='pager'>"
			count := pages + 1
			//begin page
			if (page != beginnum) && (page > beginnum) {
				raw = raw + "<li class='previous'><a href='?" + keyword + "page=" + strconv.Itoa(page-1) + "'>&laquo;</a></li>"
			}

			for i := 1; i < count; i++ {
				//current page and loop pages
				if i == page {
					raw = raw + "<li class='active'><a href='javascript:void();'>" + strconv.Itoa(i) + "</a></li>"
				} else {
					raw = raw + "<li><a href='?" + keyword + "page=" + strconv.Itoa(i) + "'>" + strconv.Itoa(i) + "</a></li>"
				}
			}

			//next page
			if (page != endnum) && (page < endnum) {
				raw = raw + "<li class='next'><a href='?" + keyword + "page=" + strconv.Itoa(page+1) + "'>&raquo;</a></li>"
			}
			raw = raw + "</ul>"
		}

		output = template.HTML(raw)

	case style == 2:
		/*
			<div class="pagination"><ul>
					<li><a href="#">&laquo;</a></li>
					<li class="active"><a href="#">1</a></li>
					<li><a href="#">2</a></li>
					<li><a href="#">3</a></li>
					<li><a href="#">4</a></li>
					<li><a href="#">&raquo;</a></li>
			</ul></div>
		*/

		if results_max > 0 {
			raw = "<div class='pagination pagination-centered'><ul>"
			count := pages + 1
			//begin page
			if (page != beginnum) && (page > beginnum) {
				raw = raw + "<li><a href='?" + keyword + "page=" + strconv.Itoa(page-1) + "'>&laquo;</a></li>"
			}
			for i := 1; i < count; i++ {
				//current page and loop pages
				if i == page {
					raw = raw + "<li class='active'><a href='javascript:void();'>" + strconv.Itoa(i) + "</a></li>"
				} else {
					raw = raw + "<li><a href='?" + keyword + "page=" + strconv.Itoa(i) + "'>" + strconv.Itoa(i) + "</a></li>"
				}
				//next page
				if (page != endnum) && (page < endnum) && (i == pages) {
					raw = raw + "<li><a href='?" + keyword + "page=" + strconv.Itoa(page+1) + "'>&raquo;</a></li>"
				}
			}
			raw = raw + "</ul></div>"
		}

		output = template.HTML(raw)

	case style == 3:
		/*
			<div class="pagenav">
				<p>
					<a href="" class="on">1</a>
					<a href="">2</a>
					<a href="">3</a>
					<a href="">4</a>
				</p>
			</div>
		*/
		raw = "<div class=\"pagenav\">"
		if results_max > 0 {
			raw = raw + "<p>"
			count := int(pages + 1)
			for i := 1; i < count; i++ {
				if i == page { //当前页
					raw = raw + "<a onclick=\"javascript:void();\" class=\"on\">" + strconv.Itoa(i) + "</a>"
				} else { //普通页码链接
					raw = raw + "<a href='?" + keyword + "page=" + strconv.Itoa(i) + "'>" + strconv.Itoa(i) + "</a>"
				}
			}
			if (page != pages) && (page < pages) { //下一页
				raw = raw + "<a class='next' href='?" + keyword + "page=" + strconv.Itoa(page+1) + "'>下一页</a>"
			}

		} else {
			raw = raw + "<h2>No Data!</h2>"
			raw = raw + "<span class='page-numbers'>共0页</span>"
		}
		raw = raw + "</p>"
		output = template.HTML(raw + "</div>")

	}

	return output
}

/** 微博时间格式化显示
 * @param timestamp，标准时间戳
 */
func TimeSince(created time.Time) string {

	//减去8小时
	//d, _ := time.ParseDuration("-8h")
	d, _ := time.ParseDuration("-0h")
	t := created.Add(d)

	since := int(time.Since(t).Minutes())
	output := ""
	switch {
	case since < 0:
		output = fmt.Sprintf("穿越了 %d 分钟..", -since)
	case since < 1:
		output = "刚刚" //"小于 1 分钟"
	case since < 60:
		output = fmt.Sprintf("%d 分钟之前", since)
	case since < 60*24:
		output = fmt.Sprintf("%d 小时之前", since/(60))
	case since < 60*24*30:
		output = fmt.Sprintf("%d 天之前", since/(60*24))
	case since < 60*24*365:
		output = fmt.Sprintf("%d 月之前", since/(60*24*30))
	default:
		output = fmt.Sprintf("%d 年之前", since/(60*24*365))
	}
	return output
}

func SmcTimeSince(timeAt time.Time) string {
	now := time.Now()
	since := math.Abs(float64(now.UTC().Unix() - timeAt.UTC().Unix()))

	output := ""
	switch {
	case since < 60:
		output = "刚刚"
	case since < 60*60:
		output = fmt.Sprintf("%v分钟前", math.Floor(since/60))
	case since < 60*60*24:
		output = fmt.Sprintf("%v小时前", math.Floor(since/3600))
	case since < 60*60*24*2:
		output = fmt.Sprintf("昨天%v", timeAt.Format("15:04"))
	case since < 60*60*24*3:
		output = fmt.Sprintf("前天%v", timeAt.Format("15:04"))
	case timeAt.Format("2006") == now.Format("2006"):
		output = timeAt.Format("1月2日 15:04")
	default:
		output = timeAt.Format("2006年1月2日 15:04")
	}
	// if math.Floor(since/3600) > 0 {
	//     if timeAt.Format("2006-01-02") == now.Format("2006-01-02") {
	//         output = "今天 "
	//         output += timeAt.Format("15:04")
	//     } else {
	//         if timeAt.Format("2006") == now.Format("2006") {
	//             output = timeAt.Format("1月2日 15:04")
	//         } else {
	//             output = timeAt.Format("2006年1月2日 15:04")
	//         }
	//     }
	// } else {
	//     m := math.Floor(since / 60)
	//     if m > 0 {
	//         output = fmt.Sprintf("%v分钟前", m)
	//     } else {
	//         output = "刚刚"
	//     }
	// }
	return output
}

//获取这个小时的开始点
func ThisHour() time.Time {
	t := time.Now()
	year, month, day := t.Date()
	hour, _, _ := t.Clock()

	return time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
}

//获取今天的开始点
func ThisDate() time.Time {
	t := time.Now()
	year, month, day := t.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

//获取这周的开始点
func ThisWeek() time.Time {
	t := time.Now()
	year, month, day := t.AddDate(0, 0, -1*int(t.Weekday())).Date()

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

//获取这月的开始点
func ThisMonth() time.Time {
	t := time.Now()
	year, month, _ := t.Date()

	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

//获取今年的开始点
func ThisYear() time.Time {
	t := time.Now()
	year, _, _ := t.Date()

	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
}

func FixedpathByNumber(n int, layer int) string {

	hash := md5.New()
	o := ""
	for i := 1; i < layer+1; i++ {

		s := strconv.Itoa(RangeRand(n^n/3+i) / 33)
		hash.Write([]byte(s))
		result := hex.EncodeToString(hash.Sum(nil))
		r := result[0:n]
		o += r + "/"
	}
	return o
}

func FixedpathByString(s string, layer int) string {

	hash := md5.New()
	output := ""
	for i := 1; i < layer+1; i++ {

		s += s + strconv.Itoa(i+i*i)
		hash.Write([]byte(s))
		result := hex.EncodeToString(hash.Sum(nil))
		r := result[0:2]
		output += r + "/"
	}
	return output
}

func StringNewRand(len int) string {

	u := make([]byte, len/2)

	// Reader is a global, shared instance of a cryptographically strong pseudo-random generator.
	// On Unix-like systems, Reader reads from /dev/urandom.
	// On Windows systems, Reader uses the CryptGenRandom API.
	_, err := io.ReadFull(crand.Reader, u)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", u)
}

// NewUUID generates a new UUID based on version 4.
func StringNewUUID() string {

	u := make([]byte, 16)

	// Reader is a global, shared instance of a cryptographically strong pseudo-random generator.
	// On Unix-like systems, Reader reads from /dev/urandom.
	// On Windows systems, Reader uses the CryptGenRandom API.
	_, err := io.ReadFull(crand.Reader, u)
	if err != nil {
		panic(err)
	}

	// Set version (4) and variant (2).
	var version byte = 4 << 4
	var variant byte = 2 << 4
	u[6] = version | (u[6] & 15)
	u[8] = variant | (u[8] & 15)

	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

//round() 函数对浮点数进行四舍五入
//语法 round(val,prec) 参数 val 规定要舍入的数字。 prec 规定小数点后的位数
func Round(val float64, prec int) float64 {
	var t float64
	f := math.Pow10(prec)
	x := val * f
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return val
	}
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.50000000001 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.50000000001 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f

	if !math.IsInf(x, 0) {
		return x
	}

	return t
}

//生成规定范围内的整数
//设置起始数字范围，0开始,n截止
func RangeRand(n int) int {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(n)

}

//标准正态分布随机整数，n为随机个数,从0开始
func Nrand(n int64) float64 {
	//sample = NormFloat64() * desiredStdDev + desiredMean
	// 默认位置参数(期望desiredMean)为0,尺度参数(标准差desiredStdDev)为1.

	var i, sample int64 = 0, 0
	desiredMean := 0.0
	desiredStdDev := 100.0

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i < n {
		rn := int64(r.NormFloat64()*desiredStdDev + desiredMean)
		sample = rn % n
		i += 1
	}

	return math.Abs(float64(sample))
}

// 对字符串进行md5哈希,
// 返回32位小写md5结果
/*
func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
*/
func MD5(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	result := hex.EncodeToString(hash.Sum(nil))
	return result
}

// 对字符串进行md5哈希,
// 返回16位小写md5结果
func MD5_16(s string) string {
	return MD5(s)[8:24]
}

// 对字符串进行sha1哈希,
// 返回42位小写sha1结果
func SHA1(s string) string {

	hasher := sha1.New()
	hasher.Write([]byte(s))

	//result := fmt.Sprintf("%x", (hasher.Sum(nil)))
	result := hex.EncodeToString(hasher.Sum(nil))
	return result
}

//AES加密
func AesEncrypt(content string, privateKey string, publicKey string) (string, error) {

	if c, err := aes.NewCipher([]byte(privateKey)); err != nil {
		fmt.Println("AesEncrypt:", err)
		return "", err
	} else {

		cfb := cipher.NewCFBEncrypter(c, []byte(publicKey))
		ciphertext := make([]byte, len(content))
		cfb.XORKeyStream(ciphertext, []byte(content))

		return string(ciphertext), err
	}

}

//AES解密
func AesDecrypt(ciphertext string, privateKey string, publicKey string) (string, error) {

	if c, err := aes.NewCipher([]byte(privateKey)); err != nil {
		return "", err
	} else {

		cipherz := []byte(ciphertext)
		cfbdec := cipher.NewCFBDecrypter(c, []byte(publicKey))
		contentCopy := make([]byte, len(cipherz))
		cfbdec.XORKeyStream(contentCopy, cipherz)

		return string(contentCopy), err
	}
}

// RSA加密
func RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(crand.Reader, pub, origData)
}

// RSA解密
func RsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(crand.Reader, priv, ciphertext)
}

func Filehash(path_or string, file_or *os.File) (string, error) {
	if (path_or != "" && file_or == nil) || (path_or == "" && file_or != nil) {
		if path_or != "" && file_or == nil {

			if file, err := os.Open(path_or); err != nil {
				return "", err
			} else {
				defer file.Close()
				h := sha1.New()

				if _, erro := io.Copy(h, file); erro != nil {
					return "", erro
				} else {

					//return fmt.Srintf("%x", h.Sum(nil))
					result := hex.EncodeToString(h.Sum(nil))
					//result := fmt.Sprintf("%d", h.Sum(nil))
					//result, _ := fmt.Printf("%d", h.Sum(nil))
					return result, nil
				}
			}
		} else {
			h := sha1.New()

			if _, erro := io.Copy(h, file_or); erro != nil {
				return "", erro
			} else {

				//return fmt.Srintf("%x", h.Sum(nil))
				result := hex.EncodeToString(h.Sum(nil))
				//result := fmt.Sprintf("%d", h.Sum(nil))
				//result, _ := fmt.Printf("%d", h.Sum(nil))
				return result, nil
			}
		}
	}
	return "", errors.New("没有参数无法生成hash,请输入文件路径 或 *os.File!")
}

func Filehash_number(path string) (int, error) {

	if file, err := os.Open(path); err != nil {
		return 0, err
	} else {

		h := sha1.New()

		if _, erro := io.Copy(h, file); erro != nil {
			return 0, erro
		} else {

			//dst, _ := strconv.Atoi(fmt.Sprintf("%d", h.Sum(nil)))
			//return fmt.Srintf("%x", h.Sum(nil))
			//result := fmt.Sprintf("%d", h.Sum(nil))
			result, _ := fmt.Printf("%d", h.Sum(nil))
			return result, nil
		}
	}

}

func Filehash_block(path string, block int64) string {
	file, err := os.Open(path)
	defer file.Close()
	hash := ""

	if err != nil {
		return ""
	}

	data := make([]byte, block)
	for {
		n, err := file.Read(data)

		if n != 0 {
			//hash = MD5(string(data))
			hash = SHA1(string(data))
		} else {
			break
		}

		if err != nil && err != io.EOF {
			//panic(err)
			return ""
		}
	}

	return hash
}

/**
* user : example@example.com login smtp server user
* password: xxxxx login smtp server password
* host: smtp.example.com:port   smtp.163.com:25
* to: example@example.com;example1@163.com;example2@sina.com.cn;...
* subject:The subject of mail
* body: The content of mail
* mailtype: mail type html or text
 */
func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func GetSensitiveInfoRemovedEmail(email string) string {
	const (
		mail_separator_sign = "@"
		min_mail_id_length  = 2
	)

	emailSepPos := strings.Index(email, mail_separator_sign)

	if emailSepPos < 0 {
		return email
	}

	mailId, mailDomain := email[:emailSepPos], email[emailSepPos+1:]

	if mailIdLength := len(mailId); mailIdLength > min_mail_id_length {
		firstChar, lastChar := string(mailId[0]), string(mailId[mailIdLength-1])
		stars := "***"
		switch mailIdLength - min_mail_id_length {
		case 1:
			stars = "*"
		case 2:
			stars = "**"
		}
		mailId = firstChar + stars + lastChar
	}

	result := mailId + mail_separator_sign + mailDomain
	return result
}

func Html2str(html string) string {
	src := string(html)

	//替换HTML的空白字符为空格
	re := regexp.MustCompile(`\s`) //ns*r
	src = re.ReplaceAllString(src, " ")

	//将HTML标签全转换成小写
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

func Str2html(raw string) template.HTML {
	return template.HTML(raw)
}

//截取字符
func Substr(str string, start, length int, symbol string) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end]) + symbol
}

func GetFile(file_url string, file_path string, useragent string, referer string) error {

	f, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("os.OpenFile errors:", err)
		return err
	}
	stat, err := f.Stat() //获取文件状态
	if err != nil {
		fmt.Println("f.Stat() errors:", err)
		return err
	}

	ss, _ := strconv.Atoi(fmt.Sprintf("%v", stat.Size))
	f.Seek(int64(ss), 0) //把文件指针指到文件末

	req, err := http.NewRequest("GET", file_url, nil)
	if err != nil {
		fmt.Println("http.NewRequest errors:", err)
		return err
	}

	if useragent == "default" {
		useragent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.64 Safari/537.31"
	}

	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	req.Header.Set("User-Agent", useragent)
	req.Header.Set("Range", "bytes="+fmt.Sprintf("%s", stat.Size)+"-")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do(req) errors:", err)
		return err
	}

	defer f.Close()
	defer resp.Body.Close()

	if written, err := io.Copy(f, resp.Body); err != nil {
		return err
	} else {

		if fs, e := os.Stat(file_path); e != nil {
			if ferr := os.Remove(file_path); ferr != nil {
				fmt.Println("Remove file error:", ferr)
			}
			return err
		} else {

			if rh, e := strconv.Atoi(resp.Header.Get("Content-Length")); e != nil || (fs.Size() != int64(rh)) {
				if rh != 0 {

					if fs.Size() != int64(rh) {

						er := errors.New(file_url + " save failed!")
						fmt.Println(er)

						if ferr := os.Remove(file_path); ferr != nil {
							fmt.Println("Remove file error:", ferr)
						}
						return er

					}
					return e
				} else {

					fmt.Println(file_url + " download success!")
					fmt.Println("written: ", written)
				}
			} else {

				fmt.Println(file_url + " download success!")
				fmt.Println("written: ", written)
			}
		}
	}
	return err
}

func PostFile(filepath string, actionurl string, fieldname string) (*http.Response, error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	// use the body_writer to write the Part headers to the buffer
	_, err := body_writer.CreateFormFile(fieldname, filepath)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	// the file data will be the second part of the body
	fh, err := os.Open(filepath)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}
	defer fh.Close()
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()
	close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(close_string)
	// use multi-reader to defer the reading of the file data until writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filepath)
		return nil, err
	}

	if req, err := http.NewRequest("POST", actionurl, request_reader); err != nil {
		return nil, err
	} else {

		// Set headers for multipart, and Content Length
		req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
		req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())

		return http.DefaultClient.Do(req)
	}

}

func WriteFile(path string, filename string, content string) error {
	//path = path[0 : len(path)-len(filename)]
	filename = path + filename
	os.MkdirAll(path, 0644)
	file, err := os.Create(filename)
	if err != nil {
		return err
	} else {
		writer := bufio.NewWriter(file)
		writer.WriteString(content)
		writer.Flush()
	}
	defer file.Close()
	return nil
}

func MoveFile(frompath string, topath string) error {

	if fromfile, err := os.Open(frompath); err != nil {
		return err
	} else {

		if tofile, err := os.OpenFile(topath, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
			return err
		} else {
			io.Copy(tofile, fromfile)
			fromfile.Close()
			tofile.Close()
			os.Remove(frompath)
			/*
				io.Copy 在一般情况下拷贝不会出错，多个携程访问的时候可能会出现“read ./data/*.png: Access is denied.”的错误，
				造成这个错误的原因很可能是由于多个协程争抢打开文件导致，然而实际情况可能报错后却又删除成功。
				如果我们根据这个错误作出判断的话就会错上加错，所以在这里不做任何判断，完全由上帝决定好了。
			*/
			return nil

		}
	}

}

func Htmlquote(text string) string {
	//HTML编码为实体符号
	/*
	   Encodes `text` for raw use in HTML.
	       >>> htmlquote("<'&\\">")
	       '&lt;&#39;&amp;&quot;&gt;'
	*/

	text = strings.Replace(text, "&", "&amp;", -1) // Must be done first!
	text = strings.Replace(text, "…", "&hellip;", -1)
	text = strings.Replace(text, "<", "&lt;", -1)
	text = strings.Replace(text, ">", "&gt;", -1)
	text = strings.Replace(text, "'", "&#39;", -1)
	text = strings.Replace(text, "\"", "&#34;", -1)
	text = strings.Replace(text, "\"", "&quot;", -1)
	text = strings.Replace(text, "“", "&ldquo;", -1)
	text = strings.Replace(text, "”", "&rdquo;", -1)
	text = strings.Replace(text, " ", "&nbsp;", -1)
	return text
}

func Htmlunquote(text string) string {
	//实体符号解释为HTML
	/*
	   Decodes `text` that's HTML quoted.
	       >>> htmlunquote('&lt;&#39;&amp;&quot;&gt;')
	       '<\\'&">'
	*/

	// strings.Replace(s, old, new, n)
	// 在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换

	text = strings.Replace(text, "&nbsp;", " ", -1)
	text = strings.Replace(text, "&rdquo;", "”", -1)
	text = strings.Replace(text, "&ldquo;", "“", -1)
	text = strings.Replace(text, "&quot;", "\"", -1)
	text = strings.Replace(text, "&#34;", "\"", -1)
	text = strings.Replace(text, "&#39;", "'", -1)
	text = strings.Replace(text, "&gt;", ">", -1)
	text = strings.Replace(text, "&lt;", "<", -1)
	text = strings.Replace(text, "&hellip;", "…", -1)
	text = strings.Replace(text, "&amp;", "&", -1) // Must be done last!
	return text
}

func CheckPassword(password string) (b bool) {
	if ok, _ := regexp.MatchString(`^[\@A-Za-z0-9\!\#\$\%\^\&\*\~\{\}\[\]\.\,\<\>\(\)\_\+\=]{4,30}$`, password); !ok {
		return false
	}
	return true
}

func CheckUsername(username string) (b bool) {
	if ok, _ := regexp.MatchString("^[\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]{4,30}$", username); !ok {
		return false
	}
	return true
}

func CheckEmail(email string) (b bool) {
	if ok, _ := regexp.MatchString(`^([a-zA-Z0-9._-])+@([a-zA-Z0-9_-])+((\.[a-zA-Z0-9_-]{2,3}){1,2})$`, email); !ok {
		return false
	}
	return true
}

/*
#gravity可用值有九个,分别是:

西北方 NorthWest：左上角为坐标原点，x轴从左到右，y轴从上到下，也是默认值。
北方   North：上部中间位置为坐标原点，x轴从左到右，y轴从上到下。
东北方 NorthEast：右上角位置为坐标原点，x轴从右到左，y轴从上到下。
西方   West：左边缘中间位置为坐标原点，x轴从左到右，y轴从上到下。
中央   Center：正中间位置为坐标原点，x轴从左到右，y轴从上到下。
东方   East：右边缘的中间位置为坐标原点，x轴从右到左，y轴从上到下。
西南方 SouthWest：左下角为坐标原点，x轴从左到右，y轴从下到上。
南方   South：下边缘的中间为坐标原点，x轴从左到右，y轴从下到上。
东南方 SouthEast：右下角坐标原点，x轴从右到左，y轴从下到上。

*/
func Thumbnail(mode string, input_file string, output_file string, output_size string, output_align string, background string) error {
	//预处理gif格式
	if strings.HasSuffix(input_file, "gif") {
		if Exist(input_file) {
			/*
				convert input_file -coalesce m_file
			*/
			cmd := exec.Command("convert", "-coalesce", input_file, input_file)
			err := cmd.Run()

			return err
		} else {
			return errors.New("需要被缩略处理的GIF图片文件并不存在!")
		}
	}

	switch {
	case mode == "resize":
		if Exist(input_file) {
			/*
				convert -resize 256x256^ -gravity center -extent 256x256  src.jpg dest.jpg
				详细使用格式 http://www.imagemagick.org/Usage/resize/
			*/
			cmd := exec.Command("convert", "-resize", output_size+"^", "-gravity", output_align, "-extent", output_size, "-background", background, input_file, output_file)
			err := cmd.Run()

			return err
		} else {
			return errors.New("需要被缩略处理的图片文件并不存在!")
		}
	case mode == "crop":
		if Exist(input_file) {
			/*
			   convert -crop 300×400 center src.jpg dest.jpg 从src.jpg坐标为x:10 y:10截取300×400的图片存为dest.jpg
			   convert -crop 300×400-10+10 src.jpg dest.jpg 从src.jpg坐标为x:0 y:10截取290×400的图片存为dest.jpg
			   详细使用格式 http://www.imagemagick.org/Usage/crop/
			*/
			cmd := exec.Command("convert", "-gravity", output_align, "-crop", output_size+"+0+0", "+repage", "-background", background, "-extent", output_size, input_file, output_file)
			err := cmd.Run()

			return err
		} else {
			return errors.New("需要被缩略处理的图片文件并不存在!")
		}
	default:
		if Exist(input_file) {

			cmd := exec.Command("convert", "-thumbnail", output_size, "-background", background, "-gravity", output_align, "-extent", output_size, input_file, output_file)
			err := cmd.Run()

			return err
		} else {
			return errors.New("需要被缩略处理的图片文件并不存在!")
		}
	}

}

func Watermark(watermark_file string, input_file string, output_file string, output_align string) error {
	//composite -gravity southeast -dissolve 30 -geometry +15%+15%  lhslogo.png input_file.jpg output_file.jpg
	cmd := exec.Command("composite", "-gravity", output_align, "-dissolve", "100", watermark_file, input_file, output_file)

	err := cmd.Run()

	if err != nil {
		return err
	} else {
		return nil
	}

}

func Rex(text string, iregexp string) (b bool) {
	if ok, _ := regexp.MatchString(iregexp, text); !ok {
		return false
	}
	return true
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//发送报文 是否加密 HTTP状态 动作URL 数据内容 RSA公匙
func SendingPackets(encrypt bool, status string, actionurl string, content string, aesKey string, aesPublicKey string, rsaPublicKey []byte) (*http.Response, error) {
	/*
	   1.对数据进行AES加密
	   2.对AES密匙KEY进行RSA加密
	   3.POST的时候,把RSA密码串放置于URL发送
	   4.POST的时候,把AES密码串放置于BODY发送
	*/
	//只有公钥则只能加密  公钥私钥都有才能解密 所以私匙不能在客户端公开  客户端获取的内容由服务端的权限控制决定
	var body_buf io.Reader
	if encrypt {
		// AES对内容进行加密
		if aes_encrypt_content, err := AesEncrypt(content, aesKey, aesPublicKey); err != nil {

			return nil, err
		} else {
			body_buf = bytes.NewBufferString(aes_encrypt_content)

			// 对AES密匙aesKey进行RSA加密

			if rsa_encrypt_content, err := RsaEncrypt([]byte(aesKey), rsaPublicKey); err != nil {

				return nil, err
			} else {
				//转换RSA密文BYTE编码为16进制字符串
				aesKey = fmt.Sprintf("%x", rsa_encrypt_content)

			}
		}

	} else {
		//无需加密
		body_buf = bytes.NewBufferString(content)

	}

	//hash就是各种内容的混合体加key的hash值,验证这个hash是否一致来保证内容不被非法更改
	createtime := strconv.Itoa(int(time.Now().UnixNano()))
	//hash+createtime+aeskey
	actionurl = actionurl + "?hash=" + Encrypt_hash(status+createtime+string(content)+string(rsaPublicKey), nil) + "-" + createtime + "-" + aesKey

	if req, err := http.NewRequest(status, actionurl, body_buf); err != nil {
		return nil, err

	} else {
		hd, err := http.DefaultClient.Do(req)
		return hd, err
	}
}

func ReceivingPackets(decrypt bool, hash string, status string, content []byte, aesPublicKey string, rsaPublicKey []byte, rsaPrivateKey []byte) ([]byte, error) {

	//防擅改校验数据
	if hash != "" {
		/*
		   1.对AES数据进行AES解密得出内容
		   2.对RSA数据进行RSA解密得出AES密匙KEY
		*/

		//分解hash+createtime+aeskey
		s := strings.Split(hash, "-")
		hash = s[0]
		createtime := s[1]
		aseKey := s[2]

		//若 decrypt为真则进行解密
		if decrypt {
			if aseKey != "" {

				//对16进制字符串aseKey进行解码
				if x, err := hex.DecodeString(aseKey); err == nil {

					//RSA解密  得出 AES KEY
					if rsa_decrypt_content, err := RsaDecrypt(x, rsaPrivateKey); err != nil {
						return nil, err
					} else {
						//还原  aseKey
						aseKey = string(rsa_decrypt_content)

						//对AES数据进行AES解密得出内容
						if aes_decrypt_content, err := AesDecrypt(string(content), aseKey, aesPublicKey); err != nil {
							return nil, err
						} else {
							content = []byte(aes_decrypt_content)
						}
					}
				} else {
					//16进制解码错误
					return nil, err
				}

			} else {
				return nil, errors.New("AES KEY为空无法进行解密")
			}
		}

		if (hash != "") && (createtime != "") {

			if Validate_hash(hash, status+createtime+string(content)+string(rsaPublicKey)) {
				//返回数据明文
				return content, nil
			} else {
				return nil, errors.New("报文无法通过数据校验")
			}
		}
	}
	return nil, errors.New("数据校验HASH值为空")
}

func GetBanner(content string) (string, error) {

	if imgs, num := GetImages(content); num > 0 {

		for _, v := range imgs {
			// 只获取本地图片,外部图片不太可靠
			if IsLocal(v) {
				if Exist(Url2local(v)) {

					return v, nil
				}
			}
			return v, errors.New("GetBanner没有图片链接")
		}
	}
	return "", errors.New("GetBanner没有图片链接")
}

func IsLocal(path string) bool {
	if path != "" {
		/*
			把本地路径的无点形式转为有点形式
			转换之后,如果之前传入的是恰好是一个网址而不是本地路径,则在后面的分拣中会把它列入并非本地路径的行列
			因为本地路径在本系统中是预设想必为当前网站项目文件夹范围内的 ./root/path  而不能跳出 到 ../root/path外,
			所以跳出到 ../root/path 外的路径必定不是本地路径!
		*/
		path = Url2local(path)

		//检查带点的本地路径
		s := strings.SplitN(path, ".", -1)
		if len(s) > 1 && len(s) < 4 {
			//通过路径的前缀是否为"/"判断是不是本地文件
			if s[1] != "" {
				if strings.HasPrefix(s[1], "/") {
					return true
				} else {
					// 第一轮次检查的时候碰上"/"开头的本地路径会判断不出来,需要再进行第2次判断"/"开头的情况
					ss := strings.SplitN("."+s[1], ".", -1)
					if len(ss) > 1 && len(ss) < 4 {
						//通过路径的前缀是否为"/"判断是不是本地文件
						if ss[1] != "" {
							if strings.HasPrefix(ss[1], "/") {
								return true
							} else {
								return false
							}
						} else {
							return false
						}
					}
					return false
				}
			} else {
				return false
			}
		}
	}
	return false
}

func Local2url(path string) string {
	if strings.HasPrefix(path, "./") {
		path = strings.Replace(path, "./", "/", -1)
	}
	return path
}

func Url2local(path string) string {
	if strings.HasPrefix(path, "/") {
		path = strings.Replace(path, "/", "./", 1)
	}
	return path
}

//设置后缀
func SetSuffix(content string, str string) string {

	content = Url2local(content)
	if content != "" {

		s := strings.SplitN(content, ".", -1)

		if len(s) > 1 && len(s) < 4 {
			// 判断是不是本地文件或本地路径
			if s[1] != "" && IsLocal(s[1]) {
				return Local2url(s[1] + str)
			} else {
				return Local2url(content)
			}
		}

	}
	return Local2url(content)
}

func GetBannerThumbnail(content string) (string, error) {
	//开始提取img
	if s, e := GetBanner(content); e == nil {

		//配置缩略图
		input_file := Url2local(s)
		output_file := Url2local(SetSuffix(s, "_banner.jpg"))
		output_size := "680x300"
		output_align := "center"
		background := "black"

		//处理缩略图
		if err := Thumbnail("crop", input_file, output_file, output_size, output_align, background); err == nil {

			return Local2url(output_file), err
		} else {

			fmt.Println("GetBannerThumbnail生成缩略图出错:", err)

			if e := os.Remove(output_file); e != nil {
				fmt.Println("GetBannerThumbnail清除残余缩略图文件出错:", e)
				return "", e
			}
			return "", err

		}
	} else {
		return "", e
	}
}

func GetThumbnails(content string) (thumbnails string, thumbnailslarge string, thumbnailsmedium string, thumbnailssmall string, err error) {
	/*
		Thumbnails        string //Original remote file
		ThumbnailsLarge   string //200x300
		ThumbnailsMedium  string //200x150
		ThumbnailsSmall   string //70x70
	*/
	//开始提取img 默认所有图片为本地文件
	if original_file, e := GetBanner(content); e == nil {

		//配置缩略图
		input_file := Url2local(original_file)
		output_file_Large := Url2local(SetSuffix(original_file, "_large.jpg"))
		output_file_Medium := Url2local(SetSuffix(original_file, "_medium.jpg"))
		output_file_Small := Url2local(SetSuffix(original_file, "_small.jpg"))
		output_size_Large := "200x300"
		output_size_Medium := "200x150"
		output_size_Small := "70x70"
		output_align := "center"
		background := "#ffffff"

		//处理缩略图
		//原始文件
		thumbnails = original_file

		//大缩略图
		if err := Thumbnail("resize", input_file, output_file_Large, output_size_Large, output_align, background); err == nil {

			thumbnailslarge = Local2url(output_file_Large)
		} else {

			fmt.Println("GetThumbnails生成thumbnailslarge缩略图出错:", err)

			if e := os.Remove(output_file_Large); e != nil {
				fmt.Println("GetThumbnails清除残余thumbnailslarge缩略图文件出错:", e)

			}
		}

		//中缩略图
		if err := Thumbnail("resize", input_file, output_file_Medium, output_size_Medium, output_align, background); err == nil {

			thumbnailsmedium = Local2url(output_file_Medium)
		} else {

			fmt.Println("GetThumbnails生成output_file_Medium缩略图出错:", err)

			if e := os.Remove(output_file_Medium); e != nil {
				fmt.Println("GetThumbnails清除残余output_file_Medium缩略图文件出错:", e)

			}
		}

		//小缩略图
		if err := Thumbnail("resize", input_file, output_file_Small, output_size_Small, output_align, background); err == nil {

			thumbnailssmall = Local2url(output_file_Small)
		} else {

			fmt.Println("GetThumbnails生成output_file_Small缩略图出错:", err)

			if e := os.Remove(output_file_Small); e != nil {
				fmt.Println("GetThumbnails清除残余output_file_Small缩略图文件出错:", e)

			}
		}
		return thumbnails, thumbnailslarge, thumbnailsmedium, thumbnailssmall, nil
	} else {
		return "", "", "", "", e
	}
}

func MakeThumbnails(localpath string) (thumbnails string, thumbnailslarge string, thumbnailsmedium string, thumbnailssmall string, err error) {
	/*
		Thumbnails        string //Original remote file
		ThumbnailsLarge   string //200x300
		ThumbnailsMedium  string //200x150
		ThumbnailsSmall   string //70x70
	*/
	//开始提取img 默认所有图片为本地文件
	if original_file := Url2local(localpath); original_file != "" {

		//配置缩略图
		input_file := Url2local(original_file)
		output_file_Large := Url2local(SetSuffix(original_file, "_large.jpg"))
		output_file_Medium := Url2local(SetSuffix(original_file, "_medium.jpg"))
		output_file_Small := Url2local(SetSuffix(original_file, "_small.jpg"))
		output_size_Large := "200x300"
		output_size_Medium := "200x150"
		output_size_Small := "70x70"
		output_align := "center"
		background := "#ffffff"

		//处理缩略图
		//原始文件 也缩略处理以适合内容框大小
		if err := Thumbnail("thumbnail", input_file, original_file, "696x", output_align, background); err == nil {
			watermark_file := "./theme/default/static/mzr/img/watermark.png"
			Watermark(watermark_file, original_file, original_file, "SouthEast")
			thumbnails = Local2url(original_file)
		} else {

			fmt.Println("MakeThumbnails生成thumbnails缩略图出错:", err)

			if e := os.Remove(original_file); e != nil {
				fmt.Println("MakeThumbnails清除残余thumbnails缩略图文件出错:", e)

			}
		}

		//大缩略图
		if err := Thumbnail("resize", input_file, output_file_Large, output_size_Large, output_align, background); err == nil {

			thumbnailslarge = Local2url(output_file_Large)
		} else {

			fmt.Println("MakeThumbnails生成thumbnailslarge缩略图出错:", err)

			if e := os.Remove(output_file_Large); e != nil {
				fmt.Println("MakeThumbnails清除残余thumbnailslarge缩略图文件出错:", e)

			}
		}

		//中缩略图
		if err := Thumbnail("resize", input_file, output_file_Medium, output_size_Medium, output_align, background); err == nil {

			thumbnailsmedium = Local2url(output_file_Medium)
		} else {

			fmt.Println("MakeThumbnails生成output_file_Medium缩略图出错:", err)

			if e := os.Remove(output_file_Medium); e != nil {
				fmt.Println("MakeThumbnails清除残余output_file_Medium缩略图文件出错:", e)

			}
		}

		//小缩略图
		if err := Thumbnail("resize", input_file, output_file_Small, output_size_Small, output_align, background); err == nil {

			thumbnailssmall = Local2url(output_file_Small)
		} else {

			fmt.Println("MakeThumbnails生成output_file_Small缩略图出错:", err)

			if e := os.Remove(output_file_Small); e != nil {
				fmt.Println("MakeThumbnails清除残余output_file_Small缩略图文件出错:", e)

			}
		}
		return thumbnails, thumbnailslarge, thumbnailsmedium, thumbnailssmall, nil
	} else {
		return "", "", "", "", errors.New("输入的图片路径为空!")
	}
}

//  返回 图片url列表集合
func GetImages(content string) (imgs []string, num int) {

	//替换HTML的空白字符为空格
	ren := regexp.MustCompile(`\s`) //ns*r
	bodystr := ren.ReplaceAllString(content, " ")

	//匹配所有图片
	//rem := regexp.MustCompile(`<img.*?src="(.+?)".*?`) //匹配最前面的图
	rem := regexp.MustCompile(`<img.+?src="(.+?)".*?`) //匹配最前面的图
	img_urls := rem.FindAllSubmatch([]byte(bodystr), -1)

	for _, bv := range img_urls {
		if m := string(bv[1]); m != "" {

			if !ContainsSets(imgs, m) {
				imgs = append(imgs, m)
			}
		}
	}

	return imgs, len(imgs)
}

func Base64Encoding(s string) string {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	defer encoder.Close()
	encoder.Write([]byte(s))
	return buf.String()
}

//返回获得的网页内容
func GetPage(url string) (string, error) {

	//ua := "Mozilla/5.0 (Windows; U; Windows NT 8.8; en-US) AppleWebKit/883.13 (KHTML, like Gecko) Chrome/88.3.13.87 Safari/883.13"
	ua := "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.92 Safari/537.1 VERYHOURSPIDER"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err

}

//PHA算法  获取图像指纹
func GetImagePha(path string) (string, error) {

	if infile, err := os.Open(path); err != nil {
		return "", err
	} else {

		// Decode picture.
		if srcImg, _, err := image.Decode(infile); err != nil {
			fmt.Println("Decode picture:", err)
		} else {

			return PHA(srcImg), err

		}

	}
	return "", errors.New("获取图片PHA值出现错误")

}

//指纹比较
func PhaCompare(path1 string, path2 string) (int, error) {

	if fg1, err := GetImagePha(path1); err != nil {
		return -1, err
	} else {
		if fg2, err := GetImagePha(path2); err != nil {
			return -1, err
		} else {

			return CompareDiff(fg1, fg2), err
		}
	}

}

//差集
func DifferenceSets(a []string, b []string) []string {

	f := make([]string, 0)

	for _, v := range a {
		//如果a集合某元素存在于b集合中
		var in bool
		for _, vv := range b {
			if v == vv {
				in = true
				break
			}
		}
		if !in {
			f = append(f, v)
		}
	}
	return f
}

//交集
func IntersectionSets(fora []string, forb []string) []string {

	i, c, d := []string{}, []string{}, []string{}
	if len(fora) > len(forb) {

		c = forb
		d = fora

	} else {

		c = fora
		d = forb
	}
	for _, v := range c {

		//如果c集合中某元素v存在于d集合中
		for _, vv := range d {
			if v == vv {
				i = append(i, v)
				break
			}
		}
	}
	return i
}

//对称差=并集-交集  即是 并集和交集的差集就是对称差
func SymmetricDifferenceSets(fora []string, forb []string) []string {

	return DifferenceSets(UnionSets(fora, forb), IntersectionSets(fora, forb))
}

//并集
func UnionSets(fora []string, forb []string) []string {
	uvalue := []string{}
	//求两个字符串数组的并集
	for _, v := range fora {
		if ContainsSets(uvalue, v) {
			continue
		} else {
			uvalue = append(uvalue, v)
		}

	}
	for _, v := range forb {
		if ContainsSets(uvalue, v) {
			continue
		} else {
			uvalue = append(uvalue, v)
		}
	}

	return uvalue
}

func ContainsSets(values []string, ivalue string) bool {
	for _, v := range values {

		if v == ivalue {
			return true
		}
	}
	return false
}

func DelLostImages(oldz string, newz string) {

	oldfiles, onum := GetImages(oldz)
	newfiles, nnum := GetImages(newz)

	//初步过滤门槛,提高效率,因为下面的操作太多循环,能避免进入则避免
	if (onum > 0 && nnum > 0) || (onum > 0 && nnum < 1) || (onum == nnum) {

		oldfiles_local := []string{}
		newfiles_local := []string{}

		for _, v := range oldfiles {
			if IsLocal(v) {
				oldfiles_local = append(oldfiles_local, v)
				//如果本地同时也存在banner缓存文件,则加入旧图集合中,等待后面一次性删除
				if p := Url2local(SetSuffix(v, "_banner.jpg")); Exist(p) {
					oldfiles_local = append(oldfiles_local, p)
				}
			}
		}
		//fmt.Println("旧图集合:", oldfiles_local)

		for _, v := range newfiles {
			if IsLocal(v) {
				newfiles_local = append(newfiles_local, v)
			}
		}
		//fmt.Println("新图集合:", newfiles_local)

		//旧图集合-新图集合 =待删图集
		for k, v := range DifferenceSets(oldfiles_local, newfiles_local) {
			if p := Url2local(v); Exist(p) { //如若文件存在,则处理,否则忽略
				fmt.Println("删除文件:", p)
				if err := os.Remove(p); err != nil {
					fmt.Println("#", k, ",DEL FILE ERROR:", err)
				}
			}
		}
	}

}

//字符串转换来unit16
func StringToUTF16(s string) []uint16 {
	return utf16.Encode([]rune(s + "\x00"))
}

func VerifyUserfile(path string, usr string) bool {
	fname := fpath.Base(path)[0:48]
	if fhashed, e := Filehash(Url2local(path), nil); e == nil {
		return Validate_hash(fname, fhashed+usr)
	} else {
		return false
	}

}

//获取文本中 @user 中的用户名集合
func AtUsers(content string) (usrs []string) {
	// 新浪微博中的用户名格式为是“4-30个字符，支持中英文、数字、"_"或减号”
	//也就是说，支持中文、字母、数字、下划线及减号，并且是4到30个字符,这里 汉字作为一个字符

	rx := regexp.MustCompile("@([\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)")
	//^[\\x{4e00}-\\x{9fa5}]+$
	//rx := regexp.MustCompile("@[^,，：:\\s@]{4,30}")
	atusr := rx.FindAllSubmatch([]byte(content), -1)
	for _, v := range atusr {
		if m := string(v[1]); m != "" {
			//usrs = append(usrs, m)
			if ContainsSets(usrs, m) {
				continue
			} else {
				usrs = append(usrs, m)
			}
		}
	}

	return usrs
}

//获取文本中 @urls 的网址集合 ###AtPages函数的作用是提取@后面的url网址,并不是提取图片,请注意!
func AtPages(content string) ([]string, string) {
	urls := []string{}
	rxs := "@([a-zA-z]+://[^\\s]*)"
	rx := regexp.MustCompile(rxs)

	aturl := rx.FindAllSubmatch([]byte(content), -1)

	if len(aturl) > 0 {

		for _, v := range aturl {
			if m := string(v[0]); m != "" {

				if !ContainsSets(urls, m[1:]) {
					//替换@link链接
					content = strings.Replace(content, m,
						"<a href='/url/?localtion="+m[1:]+"' target='_blank' rel='nofollow'><span>@</span><span>"+m[1:]+"</span></a>", -1)

					urls = append(urls, m[1:])
				}
			}
		}
	}

	return urls, content
}

func AtPagesGetImages(content string) ([]string, string) {
	imgs := []string{}
	links, content := AtPages(content)
	for _, v := range links {

		page, _ := GetPage(v)
		imgz, n := GetImages(page)

		if n > 0 {
			for _, vv := range imgz {
				//vv为单图url 相对路径的处理较为复杂,这里暂时抛弃相对路径的图片 后续再修正
				if strings.HasPrefix(vv, "/") {

					if strings.HasSuffix(v, "/") {
						vv = v + vv[1:]
					} else {
						vv = v + vv
					}

					//vv = v + vv[1:]
				} else {
					vv = Fixurl(v, vv)
				}
				if !strings.HasPrefix(vv, "../") {

					if !strings.HasSuffix(v, "js") {
						if !ContainsSets(imgs, vv) {
							imgs = append(imgs, vv)
						}
					}
				}

			}
		}
	}
	return imgs, content
}

func Fixurl(current_url, url string) string {

	re1, _ := regexp.Compile("http[s]?://[^/]+")
	destrooturl := re1.FindString(current_url)

	//当url为：//wahaha/xiaoxixi/tupian.png
	if strings.HasPrefix(url, "//") {
		url = "http:" + url
	} else if strings.HasPrefix(url, "/") {
		// re1,_ := regexp.Compile("http[s]?://[^/]+")
		// destrooturl := re1.FindString(current_url)
		url = destrooturl + url
	}

	//当url为："../wahaha/xiaoxixi/tupian.png"、"./wahaha/xiaoxixi/tupian.png"、"wahaha/xiaoxixi/tupian.png"
	if !strings.HasPrefix(url, "/") && !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		// current_url = strings.TrimSuffix(current_url, "/")
		if destrooturl == current_url {
			url = current_url + "/" + url
		} else {
			re2, _ := regexp.Compile("[^/]+?$")
			url = re2.ReplaceAllString(current_url, "") + url
		}

	}

	return url
}

func PingFile(url string) bool {
	r, e := http.Head(url)
	if e != nil {
		return false
	} else {

		if r.Status == "404" {
			return false
		}
	}
	return true
}

//分割tags
func Tags(content string, str string) []string {

	if content != "" && str != "" {

		ss := []string{}
		s := strings.SplitN(content, str, -1)
		for _, v := range s {
			if v != "" {
				ss = append(ss, v)
			}
		}

		return ss

	}
	return nil
}

//返回数字带数量级单位 千对k 百万对M  京对G
func Metric(n int64) string {

	switch {
	case n >= 1000:
		return fmt.Sprint(n/1000, "k")
	case n >= 1000000:
		return fmt.Sprint(n/1000000, "m")
	default:
		return fmt.Sprint(n)
	}
}

//根据用户邮箱显示Gravatar头像
func Gravatar(email string, height int) string {
	if email != "" && height != 0 {
		// 将邮箱转换成MD5哈希值，并设置图像的大小为height像素
		usergravatar := `http://www.gravatar.com/avatar/` + MD5(email) + `?s=` + strconv.Itoa(height)
		return usergravatar
	} else {
		return ""
	}
}

func Markdown(md string) template.HTML {
	//这句有XSS漏洞	output := Htmlunquote(string(blackfriday.MarkdownCommon([]byte(md))))

	text := strings.Replace(string(blackfriday.MarkdownCommon([]byte(md))), "&amp;#34;", "&#34;", -1) // &#34; """
	text = strings.Replace(text, "&amp;#39;", "&#39;", -1)                                            //&#39; '''
	text = strings.Replace(text, "&amp;lt;", "&lt;", -1)                                              // <
	text = strings.Replace(text, "&amp;gt;", "&gt;", -1)                                              // >
	text = strings.Replace(text, "&amp;hellip;", "&hellip;", -1)                                      // 省略号
	text = strings.Replace(text, "&hellip;", "…", -1)                                                 // 省略号
	//text = strings.Replace(text, "onerror", "<i>on</i>error", -1)

	return template.HTML(text)
}

func Markdown2Text(md string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	/*
		text := re.ReplaceAllString(string(blackfriday.MarkdownCommon([]byte(md))), "")
		text = strings.Replace(text, "&amp;#34;", "&#34;", -1) // &#34; """
		text = strings.Replace(text, "&amp;#39;", "&#39;", -1) //&#39; '''
		text = strings.Replace(text, "&amp;lt;", "&lt;", -1)   // <
		text = strings.Replace(text, "&amp;gt;", "&gt;", -1)   // >
		text = strings.Replace(text, "&hellip;", "…", -1)  // 省略号
	*/

	text := re.ReplaceAllString(string(Markdown(md)), "")

	return text
}
