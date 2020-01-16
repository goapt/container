# Golang容器注入API设计

<a href="https://github.com/goapt/container/actions"><img src="https://github.com/goapt/container/workflows/container/badge.svg" alt="Build Status"></a>


## API 设计

### 注册容器
```golang
di := container.New()
// 绑定合约
di.Register(func() contract.User {
	 return &repostiory.User{}
})

// 绑定实例而非接口,并且同时依赖参数的注入，而注入的参数必须已经被注册
di.Register(func(user contract.User) *service.User {
    return &service.User{
		User:user
	}
})
```

### 注入容器
```golang
//注入变量
var user *UserService
di.Make(&user)


//注入函数参数
var user *UserService
di.Make(function(user *UserService){
    //somthing
})
```