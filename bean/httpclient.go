package bean

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"twitter/utils"
)

func GetHttpClientFromFile() *http.Client{
	cookieJar, _ := cookiejar.New(nil)
	client := http.Client{Jar:cookieJar}
	cookies := GetCookieInfo()
	u, err := url.Parse("https://twitter.com/")
	utils.CheckErr(err)
	client.Jar.SetCookies(u, cookies)
	return &client
}

func GetHttpClient() *http.Client{
	cookieJar, _ := cookiejar.New(nil)
	client := http.Client{Jar:cookieJar}

	req, err := http.NewRequest("GET", "https://twitter.com", nil)
	utils.CheckErr(err)
	resp, err := client.Do(req)
	utils.CheckErr(err)
	data, err := ioutil.ReadAll(resp.Body)
	utils.CheckErr(err)

	defer resp.Body.Close()

	myExp := regexp.MustCompile(`value="(?P<first>(.*?))" name="authenticity_token"`)
	res := myExp.FindStringSubmatch(string(data))
	l := len(res)
	if l == 0{
		return nil
	}
	t := res[l - 1]

	v := url.Values{}
	v.Set("session[username_or_email]", "kaijianzheng@gmail.com")
	v.Set("session[password]", "zkj13902713073")
	v.Set("return_to_ssl", "true")
	v.Set("redirect_after_login", "/")
	v.Set("authenticity_token", t)
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	req, _ = http.NewRequest("POST", "https://twitter.com/sessions", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, err = client.Do(req)
	utils.CheckErr(err)

	data, err = ioutil.ReadAll(resp.Body)
	u, err := url.Parse("https://twitter.com/")
	utils.CheckErr(err)

	var cookies []*http.Cookie
	c := &http.Cookie{Name:"_ga", Value:"GA1.2.1280886672.1541124270"}
	cookies = append(cookies, c)
	client.Jar.SetCookies(u, cookies)
	for _, cookie := range client.Jar.Cookies(u) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
	SaveCookieInfo(u, &client)

	return &client
}

func SaveCookieInfo(url *url.URL, client *http.Client) error{
	f, err := os.Create("cookies.txt")
	utils.CheckErr(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, cookie := range client.Jar.Cookies(url){
		lineStr := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)
		fmt.Fprintln(w, lineStr)
	}
	return w.Flush()
}

func GetCookieInfo()[]*http.Cookie{
	f, err := os.Open("cookies.txt")
	utils.CheckErr(err)
	defer f.Close()

	buf := bufio.NewReader(f)
	var cookies []*http.Cookie
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF{
			break
		}
		utils.CheckErr(err)
		line = strings.TrimRight(line, "\n")
		cs := strings.Split(line, ":")
		if len(cs) == 2{
			c := &http.Cookie{Name:cs[0], Value:cs[1]}
			cookies = append(cookies, c)
		}
	}
	return cookies
}