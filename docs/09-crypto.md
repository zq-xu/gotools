# 09 加解密

`cryptox` 只做一件事：AES-CBC + PKCS#7 的对称加解密，用于「密码、token 这类短字符串在库里以密文存放」。它不管理密钥，也不做密钥轮换。

## 算法与格式

| 项 | 值 |
| --- | --- |
| 算法 | AES-CBC |
| 填充 | PKCS#7 |
| IV | 每次加密随机生成，拼在密文前面 |
| 输出 | `base64(iv + ciphertext)` |
| 密钥长度 | 16 / 24 / 32 字节（AES-128 / 192 / 256） |

密文里没有版本号或 key id，**换密钥后旧数据无法解密**——这是设计上的取舍，需要轮换就得自己存版本信息。

## 装配

配置段（base64 编码的密钥）：

```yaml
crypto:
  aesKey: "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4"   # 24 字节 → AES-192
```

三种装配方式，按场景选一种：

```go
// 1. 由配置驱动（推荐）：Init 会解码、校验长度、构造实例
configx.Register("crypto", &cryptox.CryptoConfig, func() error {
	return cryptox.Init(cryptox.CryptoConfig)
})

// 2. 已经有密钥字节（例如从 KMS / 环境变量读）
cryptox.SetCrypto(cryptox.NewAesPasswordCrypto(key))

// 3. 换成自己的实现（HSM、远程加密服务）
cryptox.SetCrypto(&myCrypto{})
```

> 引用 `storex/account` 的项目**不需要**写第 1 段：该包在自己的 `init()` 里已经注册了 `crypto`。只有单独使用 `cryptox` 时才会缺这一段。

配置驱动的 `Init` 行为：

| 输入 | 结果 |
| --- | --- |
| `aesKey` 为空 | 直接返回 nil，不初始化（启动不会失败） |
| base64 非法 | 返回 error，`Setup` 中断 |
| 解码后不是 16/24/32 字节 | 返回 error，错误信息里带实际长度 |
| 合法 | 设置包级 `Crypto` |

## 使用

```go
cipher, err := cryptox.Encrypt([]byte("plain password"))
if err != nil { /* 未初始化 / 加密失败 */ }

plain, err := cryptox.Decrypt(cipher)
```

未初始化时**不会 panic**，而是返回：

```
cryptox is not initialized, call cryptox.Init or cryptox.SetCrypto first
```

这是刻意的：让「忘了配密钥」在第一次调用时暴露，而不是在启动时把整个服务拦下来（无密钥的应用也能正常启动）。

## 与账号密码的关系

`account.Account.Password` 存的是 `Encrypt` 的产物；`account.LoadLoginAccount` 用 `Decrypt` 解出明文再做字符串比较。因此：

- 建账号时**必须**先 `cryptox.Encrypt`，不能直接存明文（示例见 [02](02-rest-api-service.md) 的 `createAccount`）；
- 迁移旧数据（明文密码）时需要一次性加密回填；
- 比较用的是普通 `!=`，没有做恒定时间比较——密码校验场景可接受，但别把它当 MAC 用。

## 自定义实现

只要满足接口就能替换，适合接 KMS：

```go
type PasswordCrypto interface {
	Encrypt(plaintext []byte) (string, error)
	Decrypt(ciphertext string) (string, error)
}

cryptox.SetCrypto(kmsCrypto{})
```

注意 `SetCrypto` 直接替换包级变量，**并发装配不安全**；只在启动阶段调用一次。

相关：[05 配置与日志](05-config-and-logging.md) 讲注册段怎么参与装配。
