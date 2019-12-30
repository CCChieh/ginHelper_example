# ginHelper_example
ginHepler的示例，主要实现三个简单的API：

* http://localhost:8080	任何用户都可以访问
* http://localhost/admin:8080	只有name为admin的用户可以访问
* http://localhost/user:8080	除了admin其他用户都可以访问

每个API都有参数name，例如http://localhost:8080?name=username 返回成功的的结果为：

```json
{
    "message": "Hello username!"
}
```

下面开始构建：

1. 定义参数

   这里定义了Hello结构体：

   ```go
   // \service\hello_service.go
   type Hello struct {
   	ginHelper.Param        //内嵌ginHelper的基本Param
   	Name            string `form:"name" binding:"required"`
   }
   ```

   可以看到和gin的结构体绑定参数的方式是一样的，唯一的区别就是这里多了一行`ginHelper.Param`，这是调用的ginHelper的Param，是在`ginHelper`中实现的一个最基本的`parameter`接口的`Param`结构体。

2. 定义Service

   定义Service很简单只需要一行让其返回一个Hello信息：

   ```go
   // \service\hello_service.go
   func (param *Hello) Service() {
   	param.Ret = gin.H{"message": "Hello " + param.Name + "!"}
   }
   ```

3. 完成中间件

   我们需要对两种身份进行验证，一种是name为admin的还有一种身份是非admin：

   ```go
   // \middleware\auth_middleware.go
   func AdminMiddleware() gin.HandlerFunc {
   	return func(c *gin.Context) {
   		name := c.Query("name")
   		if name=="admin"{
   			c.Next()
   		}else {
   			c.JSON(http.StatusUnauthorized,gin.H{"message": "只能admin访问!"})
   			c.Abort()
   		}
   	}
   }
   
   func UnAdminMiddleware() gin.HandlerFunc {
   	return func(c *gin.Context) {
   		name := c.Query("name")
   		if name=="admin"{
   			c.JSON(http.StatusUnauthorized,gin.H{"message": "admin不能访问!"})
   			c.Abort()
   		}else {
   			c.Next()
   		}
   	}
   }
   
   ```

   还得在实现一个日志输出的中间件：

   ```go
   // \middleware\logger_middleware.go
   func LoggerMiddleware() gin.HandlerFunc {
   	return func(c *gin.Context) {
   		start := time.Now()
   
   		c.Next()
   
   		fmt.Printf("%s %s %s %d\n",
   			c.Request.Method,
   			c.Request.RequestURI,
   			time.Since(start),
   			c.Writer.Status(),
   		)
   	}
   }
   
   ```

   中间件的实现方式与gin的一样。

4. 定义handler

   我们定义一个简单的handler，为了引入ginHelper使得所有的Handler自动加入到gin的路由中，这里在handler包中建立一个Helper结构体：

   ```go
   // \handler\main.go
   type Helper struct {
   }
   ```

   现在就可以添加handler：

   ```go
   // \handler\hello_handler.go
   func (h *Helper) HelloHandler() (r *ginHelper.Router) {
   	return &ginHelper.Router{
   		Param:  new(service.Hello), //所需要的参数
   		Path:   "/",                //路由路径
   		Method: "GET",              //方法
   	}
   }
   
   func (h *Helper) AdminHandler() (r *ginHelper.Router) {
   	return &ginHelper.Router{
   		Param:  new(service.Hello),
   		Path:   "/admin",
   		Method: "GET",
   		Handlers: []gin.HandlerFunc{
   			middleware.AdminMiddleware(),
   			ginHelper.GenHandlerFunc,
   		},
   	}
   }
   
   func (h *Helper) UnAdminHandler() (r *ginHelper.Router) {
   	return &ginHelper.Router{
   		Param:  new(service.Hello),
   		Path:   "/user",
   		Method: "GET",
   		Handlers: []gin.HandlerFunc{
   			middleware.UnAdminMiddleware(),
   		},
   	}
   }
   
   ```

   HelloHandler是最基本的方法，只需要定义Param、Path、Method就可以了，AdminHandler则是引入中间件的写法，可以发现就是定义了 :

   ```go
   Handlers: []gin.HandlerFunc{
   			middleware.AdminMiddleware(),
   			ginHelper.GenHandlerFunc,
   		},
   ```

   `ginHelper.GenHandlerFunc`是一个占位的标记，用来设置当前自动生成handler的运行顺序。

   在UnAdminHandler中只有：

   ```go
   Handlers: []gin.HandlerFunc{
   			middleware.UnAdminMiddleware(),
   		},
   ```

   省略了`ginHelper.GenHandlerFunc`，当其省略的时候，自动生成的handler默认自动加到尾部。

5. 启动服务

   使用gin启动服务：

   ```go
   func main() {
   	r := gin.New()
   	r.Use(middleware.LoggerMiddleware())
   	ginHelper.Build(new(handler.Helper), r)
   
   	addr := fmt.Sprintf("%s:%d", "127.0.0.1", 8080)
   
   	s := &http.Server{
   		Addr:           addr,
   		Handler:        r,
   		ReadTimeout:    10 * time.Second,
   		WriteTimeout:   10 * time.Second,
   		MaxHeaderBytes: 1 << 20,
   	}
   	fmt.Println("Service run in http://", addr)
   	if err := s.ListenAndServe(); err != nil {
   		fmt.Println(err)
   	}
   }
   ```

6. 测试api

   访问：http://127.0.0.1:8080?name=username

   返回：

   ```json
   {
       "message": "Hello username!"
   }
   ```

   访问：http://127.0.0.1:8080/user?name=username

   返回：

   ```json
   {
       "message": "Hello username!"
   }
   ```

   访问：http://127.0.0.1:8080/admin?name=admin

   返回：

   ```json
   {
       "message": "Hello admin!"
   }
   ```

   访问：http://127.0.0.1:8080/admin?name=username

   返回：

   ```json
   {
       "message": "只能admin访问!"
   }
   ```

   