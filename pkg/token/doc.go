// token.go 由几个核心方法组成
// Init : 设置包级别的配置 config, config 会用于本包后面的 token 签发和解析 :
//		key: 用于签发和解析 token 的密钥；
//		identityKey: token 中用户身份的键，fastgo 中是 UserID；
//      expiration: 签发的 token 过期时间。
// Parse : 使用指定的密钥 key 解析 token，解析成功返回 token 上下文（fastgo 中是 UserID），否则报错。
// ParseRequest : 从请求头中获取令牌，并将其传递给 Parse 函数以解析令牌；
// Sign : 使用 JWT Key 签发 token，token 的 claims 中会存放用户身份（fastgo 中是 UserID）、token 生效时间、token 签发时间、token 过期时间。

package token
