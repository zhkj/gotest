package bean

const(
	//任务类型
	UserProfileTask = 1
	UserTweetTask = 2
	UserFollowingTask = 3
	UserFollowerTask = 4
	KeywordTweetTask = 5

	//任务队列大小
	TASK_QUEUE_NUM = 300
	//原始内容队列
	CONTENT_QUEUE_NUM = 600
	//存储队列大小
	STORE_QUEUE_NUM = 1000
	//数据库连接池大小
	MAX_DB_CONNS = 4
)