package util

import (
	jwt "github.com/dgrijalva/jwt-go"
	"gogin/pkg/setting"
	"time"
)

//设置本地密钥

var jwtSecret = []byte(setting.JwtSecret)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录username,password字段，所以要自定义结构体，如果想要保存更多信息，都可以添加到这个结构体中

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

//生成token

func GenerateToken(username, password string) (string, error) {
	//现在时间
	nowTime := time.Now()
	//设置过期时间
	expireTime := nowTime.Add(3 * time.Hour)

	// 创建一个我们自己的声明
	claims := Claims{
		username, // 自定义字段
		password, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(), //发放时间
			Issuer:    "go-gin",          // 签发人
		},
	}

	//NewWithClaims(method SigningMethod, claims Claims)，
	//method对应着SigningMethodHMAC struct{}，其包含SigningMethodHS256、SigningMethodHS384、SigningMethodHS512三种crypto.Hash方案
	// 使用指定的签名方法(hash)创建签名对象
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//func (t *Token) SignedString(key interface{}) 该方法内部生成签名字符串，再用于获取完整、已签名的token
	// 使用指定的secret签名并获得完整的编码后的字符串token
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

//鉴权token

func ParseToken(token string) (*Claims, error) {
	//func (p *Parser) ParseWithClaims 用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		//func (m MapClaims) Valid() 验证基于时间的声明exp, iat, nbf，注意如果没有任何声明在令牌中，仍然会被认为是有效的。并且对于时区偏差没有计算方法
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
