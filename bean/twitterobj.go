package bean

type Task struct {
	TaskType int
	TaskValue string
}

type Content struct {
	TaskType int
	ContentValue string
	TaskValue string
}

type Tweet struct {
	UserName string
	Content string
	RetCount string
	ReplyCount string
	FavCount string
}

func (tweet *Tweet)GetSql()string{
	return "insert into tweet(name,content,retcount,replycount,favcount) values (?,?,?,?,?)"
}

func (tweet *Tweet)GetVal()[]interface{}{
	return []interface{}{tweet.UserName, tweet.Content, tweet.RetCount, tweet.ReplyCount, tweet.FavCount}
}

type UserProfile struct {
	UserName string
	Introduction string
	Url string
	Following string
	Follower string
	CreateTime string
}

func (up *UserProfile)GetSql()string{
	return "insert into up(name,introduction,url,following,follower,createtime) values (?,?,?,?,?,?)"
}

func (up *UserProfile)GetVal()[]interface{}{
	return []interface{}{up.UserName, up.Introduction, up.Url, up.Following, up.Follower, up.CreateTime}
}

type UserRelation struct {
	UserA string
	UserB string
	Following string
}

func (ur *UserRelation)GetSql()string{
	return "insert into ur(namea,nameb,following) values (?,?,?)"
}

func (ur *UserRelation)GetVal()[]interface{}{
	return []interface{}{ur.UserA, ur.UserB, ur.Following}
}

type DbInstance interface{
	GetSql()string
	GetVal()[]interface{}
}

func (t *Task) GetTaskUrl() string{
	switch t.TaskType{
	case UserProfileTask:
		return "https://twitter.com/" + t.TaskValue
	case UserTweetTask:
		return "https://twitter.com/" + t.TaskValue
	case UserFollowingTask:
		return "https://twitter.com/" + t.TaskValue + "/following"
	case UserFollowerTask:
		return "https://twitter.com/" + t.TaskValue + "/followers"
	case KeywordTweetTask:
		return "https://twitter.com/search?q=" + t.TaskValue
	default:
		return ""
	}
}
