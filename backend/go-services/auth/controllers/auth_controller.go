package controllers

import (
	"net/http"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/models"   // 用户模型，返回的数据结构也在这里
	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/services" // 用户服务
	"github.com/gin-gonic/gin"
)

// controller就是用来接收我们的API请求的，然后调用service层的方法，最后返回结果给客户端

// 注册
// 注册接口在apifox上还没有搞好hh

// 感觉GPT写的差不多了，我们需要做的就是修改，然后加上我们自己需要的逻辑
func Register(c *gin.Context) {
	var user models.User

	// 把user使用ShouldBindJSON方法解析请求的json数据
	if err := c.ShouldBindJSON(&user); err != nil {
		// 这里是一些错误处理
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 调用service层的RegisterUser方法，把user传进去
	// controller要做的就是把user解析出来，然后把user传递到Service中做具体的任务
	// 待会我们再去写service层的RegisterUser方法
	if err := services.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// 上面是把user传递给RegisterUser方法，然后对于服务端的请求，我们都要返回一些东西吧
	// 然后具体要返回什么东西，就看我们文档中的要求
	// 定义了好返回的数据结构models/response.go之后，这里就进行具体的处理啦

	// 然后就把这个res发出去就好啦
	c.JSON(http.StatusOK, gin.H{})
}

// 登录
// 这个在文档中写好了，就先写这个把
func Login(c *gin.Context) {
	var user models.User // 这里定义了一个user变量，类型是models.User

	// 这里是gin框架的方法，用来解析请求的json数据，然后把解析后的数据放到user变量中
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 这里Login方法中也是和上面的Register方法一样
	// 把user从请求的JSON中解析出来，传递给service中的方法做具体的操作
	token, err := services.AuthenticateUser(&user) // 这里就是把user作为参数传递给了AuthenticateUser方法
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 刚才搞错啦，这个是login接口的代码
	// 然后具体要返回什么东西，就看我们文档中的要求
	// 定义了好返回的数据结构models/response.go之后，这里就进行具体的处理啦
	// 现在我们就是要把处理之后的User中的特定字段，放到Response中，然后返回给客户端
	// 定义一个Response变量res
	var res models.Response
	// 具体的赋值，把处理之后的user中的字段赋值给res中的字段
	res = models.Response{
		Message: "User registered successfully",
		Success: true,
		Token:   "", //这个就是经常念叨的JWT，前面生成好了这里直接赋值过来，生成的代码待会写前面
		User: models.UserResponse{
			// 这个就是根据service处理之后的user来的
			Avatar:    user.Avatar,
			ID:        user.ID,
			Ipaddress: user.Ipaddress,
			Nickname:  user.Nickname,
		},
	}
	// 然后就把这个res发出去就好啦
	c.JSON(http.StatusOK, res)
}
