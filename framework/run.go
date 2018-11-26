package framework

import (
	"database/sql"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
	"twitter/bean"
	"twitter/utils"
)

var taskQueue = make(chan bean.Task, bean.TASK_QUEUE_NUM)
var contentQueue = make(chan bean.Content, bean.CONTENT_QUEUE_NUM)
var storeQueue = make(chan bean.DbInstance, bean.STORE_QUEUE_NUM)
var dbConnQueue = make(chan *sql.DB, bean.MAX_DB_CONNS)

var taskSet = make(map[string]bool)

func init(){
	GenTaskByName("FCBarcelona")
	InitDbConnPool()
}

func InitDbConnPool(){
	for i := 0; i < bean.MAX_DB_CONNS; i++{
		db, err := sql.Open("mysql", "root:scenderut1201@tcp(127.0.0.1:3306)/test?charset=utf8")
		utils.CheckErr(err)
		dbConnQueue <- db
	}
}

func GenTaskByName(taskName string){
	if _, exist := taskSet[taskName]; !exist {
		taskQueue <- bean.Task{bean.UserProfileTask, taskName}
		taskQueue <- bean.Task{bean.UserTweetTask, taskName}
		taskQueue <- bean.Task{bean.UserFollowingTask, taskName}
		taskQueue <- bean.Task{bean.UserFollowerTask, taskName}
		taskSet[taskName] = true
	}
}

func GenTaskByWord(taskName string){
	taskQueue <- bean.Task{bean.KeywordTweetTask, taskName}
}

func Crawl(){
	log.Printf("Crawl-------------Start")
	httpClient := bean.GetHttpClientFromFile()
	for {
		log.Printf("Crawl-------------Fettching Task")
		task := <- taskQueue
		taskUrl := task.GetTaskUrl()
		req, err := http.NewRequest("GET", taskUrl, nil)
		log.Printf("Crawl-------------Fettching Task Success: ", task.TaskValue, task.TaskType, taskUrl)
		utils.CheckErr(err)
		req.Header.Set("Referer", taskUrl)
		resp, err := httpClient.Do(req)
		utils.CheckErr(err)
		data, err := ioutil.ReadAll(resp.Body)
		utils.CheckErr(err)
		contentQueue <- bean.Content{task.TaskType, string(data), task.TaskValue}
		log.Printf("Crawl-------------Get Data Success: ", task.TaskValue, task.TaskType)
	}
}

func Parse(){
	log.Printf("Parse-------------Start")
	for {
		log.Printf("Parse-------------Fettching Task")
		content := <- contentQueue
		log.Printf("Parse-------------Fettching Task Success: ", content.TaskValue, content.TaskType)
		switch content.TaskType{
		case bean.UserProfileTask:
			ParseUserProfileContent(content)
		case bean.UserTweetTask:
			ParseUserTweetContent(content)
		case bean.UserFollowingTask:
			ParseUserFollowingContent(content)
		case bean.UserFollowerTask:
			ParseUserFollowerContent(content)
		case bean.KeywordTweetTask:
			ParseKeywordTweetContent(content)
		}
		log.Printf("Parse-------------Parsing Task Success: ", content.TaskValue, content.TaskType)
	}
}

func ParseUserProfileContent(content bean.Content){
	log.Printf("Parse-------------ParseUserProfileContent")
	//玩家profile数据
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content.ContentValue))
	utils.CheckErr(err)
	var bio, url, joindate, bir, following, followers string
	dom.Find(".ProfileHeaderCard-bio.u-dir").Each(func(i int, selection *goquery.Selection) {
		bio = strings.Trim(selection.Text(), " \n")
	})
	dom.Find(".ProfileHeaderCard-urlText.u-dir>a").Each(func(i int, selection *goquery.Selection) {
		url = strings.Trim(selection.Text(), " \n")
	})
	dom.Find(".ProfileHeaderCard-joinDateText.js-tooltip.u-dir").Each(func(i int, selection *goquery.Selection) {
		joindate = strings.Trim(selection.Text(), " \n")
	})
	dom.Find(".ProfileHeaderCard-birthdateText.u-dir").Each(func(i int, selection *goquery.Selection) {
		bir = strings.Trim(selection.Text(), "\n")
	})
	dom.Find("a[data-nav=followers]>.ProfileNav-value").Each(func(i int, selection *goquery.Selection) {
		v, exist := selection.Attr("data-count")
		if exist{
			followers = strings.Trim(v, " \n")
		}
	})
	dom.Find("a[data-nav=following]>.ProfileNav-value").Each(func(i int, selection *goquery.Selection) {
		v, exist := selection.Attr("data-count")
		if exist{
			following = strings.Trim(v, " \n")
		}
	})
	storeQueue <- &bean.UserProfile{content.TaskValue, bio, url, following, followers,joindate}
}

func ParseUserTweetContent(content bean.Content){
	log.Printf("Parse-------------ParseUserTweetContent")
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content.ContentValue))
	utils.CheckErr(err)
	var tc, rpc, rec, fc []string
	dom.Find(".TweetTextSize.TweetTextSize--normal.js-tweet-text.tweet-text").Each(func(i int, selection *goquery.Selection) {
		tc = append(tc, selection.Text())
	})
	dom.Find(".ProfileTweet-action--reply.u-hiddenVisually>.ProfileTweet-actionCount").Each(func(i int, selection *goquery.Selection) {
		c, e := selection.Attr("data-tweet-stat-count")
		if e{
			rpc = append(rpc, c)
		}
	})
	dom.Find(".ProfileTweet-action--retweet.u-hiddenVisually>.ProfileTweet-actionCount").Each(func(i int, selection *goquery.Selection) {
		c, e := selection.Attr("data-tweet-stat-count")
		if e{
			rec = append(rec, c)
		}
	})
	dom.Find(".ProfileTweet-action--favorite.u-hiddenVisually>.ProfileTweet-actionCount").Each(func(i int, selection *goquery.Selection) {
		c, e := selection.Attr("data-tweet-stat-count")
		if e{
			fc = append(fc, c)
		}
	})
	if len(tc) == len(rpc) && len(rpc) == len(rec) && len(rec) == len(fc){
		for i , con := range tc{
			storeQueue <- &bean.Tweet{content.TaskValue, con, rpc[i], rec[i], fc[i]}
		}
	}
}

func ParseUserFollowingContent(content bean.Content){
	log.Printf("Parse-------------ParseUserFollowingContent")
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content.ContentValue))
	utils.CheckErr(err)
	dom.Find(".user-actions.btn-group.not-following.not-muting").Each(func(i int, selection *goquery.Selection){
		name, exist := selection.Attr("data-screen-name")
		if exist{
			storeQueue <- &bean.UserRelation{content.TaskValue, name, "1"}
			GenTaskByName(name)
		}
	})
}

func ParseUserFollowerContent(content bean.Content){
	log.Printf("Parse-------------ParseUserFollowerContent")
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content.ContentValue))
	utils.CheckErr(err)
	dom.Find(".user-actions.btn-group.not-following.not-muting").Each(func(i int, selection *goquery.Selection){
		name, exist := selection.Attr("data-screen-name")
		if exist{
			storeQueue <- &bean.UserRelation{content.TaskValue, name, "0"}
			GenTaskByName(name)
		}
	})
}

func ParseKeywordTweetContent(content bean.Content){
	log.Printf("Parse-------------ParseKeywordTweetContent")
}

func Store(){
	log.Printf("Store-------------Start")
	log.Printf("Store-------------Getting DB conn first time")
	db := <-dbConnQueue
	log.Printf("Store-------------Getting DB conn first time success")
	for {
		if db == nil{
			log.Printf("Store-------------Getting DB conn")
			db = <- dbConnQueue
			log.Printf("Store-------------Getting DB conn success")
		}
		select {
		case ins := <-storeQueue:
			log.Printf("Store-------------Getting StoreQueue object success: ", reflect.TypeOf(ins).Name())
			s := ins.GetSql()
			v := ins.GetVal()
			log.Printf("Store------------: ", v)
			stmt, err := db.Prepare(s)
			utils.CheckErr(err)
			_, err = stmt.Exec(v...)
			utils.CheckErr(err)
		case <- time.After(1 * time.Second):
			log.Printf("Store-------------Getting StoreQueue object overtime")
			if db != nil{
				dbConnQueue <- db
				db = nil
			}
			time.Sleep(2 * time.Second)
			log.Printf("Store-------------awake from sleep")
		}
	}
}
