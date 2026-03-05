package fiberresp

import "github.com/gofiber/fiber/v2"

// 响应状态码
type statusCode = int

const (
	SUCCESS               statusCode = 0
	NO_USER_FOUND_WITH_ID statusCode = 10000 // 通过id找不到用户
	REVIEW_YOUR_INPUT     statusCode = 10006 // 请求参数错误,请检查请求参数
	GET_ID_ERROR          statusCode = 10007 // 获取id失败
	INVALID_TOKEN_ID      statusCode = 10012 // 无效token

	SAVE_DB_ERROR    statusCode = 10014 // 写入数据库失败
	NOT_VALID_USER   statusCode = 10015 // 无效用户
	NOT_FOUND_USER   statusCode = 10016 // 没有找到用户
	INVALID_PASSWORD statusCode = 10017 // 无效密码

	JWT_SIGN_ERROR           statusCode = 10018 // jwt签名失败
	MISSING_OR_MALFORMED_JWT statusCode = 10019 // jwt缺失或格式错误
	INVALID_OR_EXPIRED_JWT   statusCode = 10020 // jwt无效或过期

	CREATE_ERROR      statusCode = 10021 // 创建失败
	DELETE_ERROR      statusCode = 10022 // 删除失败
	INSERT_ERROR      statusCode = 10023 // 插入失败
	SELECT_ERROR      statusCode = 10024 // 查询失败
	RECORD_EXISTS     statusCode = 10032 // 密码记录已存在
	RECORD_NOT_EXISTS statusCode = 10033 // 密码记录不存在

	SM2_DECRYPT_ERROR statusCode = 20000 // sm2解密失败
	SM2_ENCRYPT_ERROR statusCode = 20001 // sm2加密失败
	SM4_DECRYPT_ERROR statusCode = 20002 // sm4解密失败
	SM4_ENCRYPT_ERROR statusCode = 20003 // sm4加密失败
	ACCOUNT_ENCRYPT   statusCode = 20004 // 账号加密失败
	PASSWORD_ENCRYPT  statusCode = 20005 // 密码加密失败
	ACCOUNT_DECRYPT   statusCode = 20006 // 账号解密失败
	PASSWORD_DECRYPT  statusCode = 20007 // 密码解密失败

	INTERNAL_SERVER_ERROR statusCode = 30000 // 服务内部错误
	SESSION_ID_EXPIRE     statusCode = 30001 // 请求头里面的password-session-id在cache里面没有找到
)

type statusCodeInfo struct {
	statusCode statusCode
	message    string
	httpCode   int
}

// 错误码信息map
var sCodeMap map[statusCode]*statusCodeInfo = map[statusCode]*statusCodeInfo{
	SUCCESS: {
		statusCode: SUCCESS,
		message:    "success",
		httpCode:   fiber.StatusOK,
	},
	NO_USER_FOUND_WITH_ID: {
		statusCode: NO_USER_FOUND_WITH_ID,
		message:    "No user found with ID",
		httpCode:   fiber.StatusNotFound,
	},
	REVIEW_YOUR_INPUT: {
		statusCode: REVIEW_YOUR_INPUT,
		message:    "Review your input",
		httpCode:   fiber.StatusBadRequest,
	},
	GET_ID_ERROR: {
		statusCode: GET_ID_ERROR,
		message:    "get id error",
		httpCode:   fiber.StatusInternalServerError,
	},
	INVALID_TOKEN_ID: {
		statusCode: INVALID_TOKEN_ID,
		message:    "Invalid token id",
		httpCode:   fiber.StatusUnauthorized,
	},
	SAVE_DB_ERROR: {
		statusCode: SAVE_DB_ERROR,
		message:    "save db error",
		httpCode:   fiber.StatusInternalServerError,
	},
	NOT_VALID_USER: {
		statusCode: NOT_VALID_USER,
		message:    "not valid user",
		httpCode:   fiber.StatusBadRequest,
	},
	NOT_FOUND_USER: {
		statusCode: NOT_FOUND_USER,
		message:    "not found user",
		httpCode:   fiber.StatusBadRequest,
	},
	INVALID_PASSWORD: {
		statusCode: INVALID_PASSWORD,
		message:    "invalid password",
		httpCode:   fiber.StatusBadRequest,
	},
	JWT_SIGN_ERROR: {
		statusCode: JWT_SIGN_ERROR,
		message:    "jwt sign error",
		httpCode:   fiber.StatusInternalServerError,
	},
	MISSING_OR_MALFORMED_JWT: {
		statusCode: MISSING_OR_MALFORMED_JWT,
		message:    "Missing or malformed JWT",
		httpCode:   fiber.StatusUnauthorized,
	},
	INVALID_OR_EXPIRED_JWT: {
		statusCode: INVALID_OR_EXPIRED_JWT,
		message:    "Invalid or expired JWT",
		httpCode:   fiber.StatusUnauthorized,
	},
	CREATE_ERROR: {
		statusCode: CREATE_ERROR,
		message:    "create error",
		httpCode:   fiber.StatusInternalServerError,
	},
	DELETE_ERROR: {
		statusCode: DELETE_ERROR,
		message:    "delete error",
		httpCode:   fiber.StatusInternalServerError,
	},
	INSERT_ERROR: {
		statusCode: INSERT_ERROR,
		message:    "insert error",
		httpCode:   fiber.StatusInternalServerError,
	},
	SELECT_ERROR: {
		statusCode: SELECT_ERROR,
		message:    "select error",
		httpCode:   fiber.StatusInternalServerError,
	},
	RECORD_EXISTS: {
		statusCode: RECORD_EXISTS,
		message:    "record exists",
		httpCode:   fiber.StatusBadRequest,
	},
	RECORD_NOT_EXISTS: {
		statusCode: RECORD_NOT_EXISTS,
		message:    "record not exists",
		httpCode:   fiber.StatusBadRequest,
	},
	SM2_DECRYPT_ERROR: {
		statusCode: SM2_DECRYPT_ERROR,
		message:    "sm2 decrypt error",
		httpCode:   fiber.StatusBadRequest,
	},
	SM2_ENCRYPT_ERROR: {
		statusCode: SM2_ENCRYPT_ERROR,
		message:    "sm2 encrypt error",
		httpCode:   fiber.StatusBadRequest,
	},
	SM4_DECRYPT_ERROR: {
		statusCode: SM4_DECRYPT_ERROR,
		message:    "sm4 decrypt error",
		httpCode:   fiber.StatusBadRequest,
	},
	ACCOUNT_ENCRYPT: {
		statusCode: ACCOUNT_ENCRYPT,
		message:    "account encrypt error",
		httpCode:   fiber.StatusInternalServerError,
	},
	PASSWORD_ENCRYPT: {
		statusCode: PASSWORD_ENCRYPT,
		message:    "password encrypt error",
		httpCode:   fiber.StatusInternalServerError,
	},
	ACCOUNT_DECRYPT: {
		statusCode: ACCOUNT_DECRYPT,
		message:    "account decrypt error",
		httpCode:   fiber.StatusInternalServerError,
	},
	PASSWORD_DECRYPT: {
		statusCode: PASSWORD_DECRYPT,
		message:    "password decrypt error",
		httpCode:   fiber.StatusInternalServerError,
	},
	SM4_ENCRYPT_ERROR: {
		statusCode: SM4_ENCRYPT_ERROR,
		message:    "sm4 encrypt error",
		httpCode:   fiber.StatusBadRequest,
	},
	INTERNAL_SERVER_ERROR: {
		statusCode: INTERNAL_SERVER_ERROR,
		message:    "服务内部错误",
		httpCode:   fiber.StatusInternalServerError,
	},
	SESSION_ID_EXPIRE: {
		statusCode: SESSION_ID_EXPIRE,
		message:    "session id expire",
		httpCode:   fiber.StatusBadRequest,
	},
}
