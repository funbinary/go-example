package bftp

const (
	CODE_FILE_STATUS_OK              = 150 // 文件状态正常，准备打开数据连接
	CODE_COMMAND_OK                  = 200 // 命令执行成功
	CODE_COMMAND_IMPLEMENTED         = 202 // 命令未实现
	CODE_SYSTEM_STATUS               = 211 // 系统状态
	CODE_DIR_STATUS                  = 212 // 目录状态
	CODE_FILE_STATUS                 = 213 // 文件状态
	CODE_HELP_MESSAGE                = 214 // 帮助消息
	CODE_NAME_SYSTEM_TYPE            = 215 // NAME系统类型
	CODE_SERVICE_READY               = 220 // 服务已经准备就绪
	CODE_CLOSE_CONN                  = 221 // 服务关闭控制连接
	CODE_DATA_CONN_OPEN              = 225 // 数据连接打开，没有进行中的传输
	CODE_DATA_CONN_CLOSE             = 226 // 关闭数据连接。请求的文件操作已经成功执行
	CODE_ENTER_PASV                  = 227 // 进入被动模式
	CODE_USER_LOGINED                = 230 // 用户已经登录，继续执行
	CODE_TLS_AUTH_OK                 = 234 // TLS授权成功
	CODE_FILE_COMMANG_OK             = 250 // 请求的文件操作正确，已完成
	CODE_PATHNAME_CREATED            = 257 // PATHNAME已创建
	CODE_NEED_PASS                   = 331 // 用户名正确，需要密码
	CODE_NEED_ACCOUNT                = 332 // 需要用户名
	CODE_NEED_NEXT                   = 350 // 请求的文件操作需要进一步命令
	CODE_SRV_NO_AVAI                 = 421 // 服务不可用，关闭控制连接
	CODE_FAILED_OPEN_DATA_CONN       = 425 // 不能打开数据连接
	CODE_CONN_CLOSED                 = 426 // 连接被关闭，中止传输
	CODE_FILE_ACTION_NOTAKEN         = 450 // 文件操作未执行
	CODE_LOCAL_ERROR                 = 451 // 终止请求的操作：有本地错误
	CODE_STORAGE_SPACE_INSUFFICIENT  = 452 // 未执行请求的操作：存储空间不足
	CODE_COMMAND_ERROR               = 500 // 命令错误，未找到次命令的实现
	CODE_PARAM_ERROR                 = 500 // 参数错误
	CODE_CMD_IMPLEMENTED             = 502 // 命令未实现
	CODE_BAD_SEQ                     = 503 // 命令执行顺序错误
	CODE_CMD_IMPLEMENTED_WITH_PARAM  = 504 // 此参数下的命令功能未实现
	CODE_UNLOGIN                     = 530 // 未登录
	CODE_NEED_ACCOUNT_FOR_STOR_FILE  = 532 // 存储文件需要账户
	CODE_ACTION_NOTAKEN              = 550 // 未执行请求的操作
	CODE_ACTION_ABORTED              = 551 // 请求操作终止：页类型未知
	CODE_FILE_ACTION_ABORTED         = 552 // 请求文件操作终止，存储分配溢出
	CODE_ACTION_NOTAKEN_WITH_UNALLOW = 553 // 未执行请求的操作。不允许的文件名。
)
