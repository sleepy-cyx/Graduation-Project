package common

// 这里定义错误码及错误信息宏，
// 错误码（四位数）设计如下：
// MySQL相关：1xxx，Redis相关：2xxx，rpc相关：3xxx,业务逻辑相关：4xxx，输入合法性相关：5xxx成功：0000 内部出错对外返回：9999
const (
	PARAMETER_INVALID              = "Parameter Invalid"
	SUCCESS                        = "Success"
	RPC_CALL_ERROR                 = "rpc call error"
	MYSQL_CONN_ERROR               = "MySQL connection error"
	DB_DATA_NOT_FOUND              = "DB Data not found"
	MYSQL_QUERY_ERROR              = "MySQL Query error"
	USERNAME_OR_PASSWORD_INCORRECT = "Username or Password Incorrect"
	SAVE_FILE_ERROR                = "Save File Error"
	RECEIVE_FILE_ERROR             = "Receive File Error"
	SERVER_ERROR                   = "Server Error"
	UNAUTHORIZED                   = "Unauthorized"
)

const (
	ERRCODE_MYSQL_CONN_ERROR               int32 = 1000
	ERRCODE_DB_DATA_NOT_FOUND              int32 = 1001
	ERRCODE_MYSQL_QUERY_ERROR              int32 = 1002
	ERRCODE_USERNAME_OR_PASSWORD_INCORRECT int32 = 4000
	ERRCODE_SUCCESS                        int32 = 0000
	ERRCODE_RPC_CALL_ERROR                 int32 = 3000
	ERRCODE_PARAMETER_INVALID              int32 = 5000
	ERRCODE_UNAUTHORIZED                   int32 = 4003
	ERRCODE_SERVER_ERROR                   int32 = 9999
)
const DEFAULT_USER_PICTURE_URL = ""

var TypeMap = map[uint]string{
	0: "会议",
	1: "娱乐",
	2: "工作",
	3: "学习",
	4: "运动",
	5: "医疗",
	6: "家庭",
	7: "旅行",
}
