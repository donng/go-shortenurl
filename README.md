`go-shortenurl` 是一个短地址服务，项目基于慕课网的免费课程 [Go开发短地址服务](https://www.imooc.com/learn/1150) 实现和优化，首先感谢创作者的无私贡献。

慕课网的路由和中间件分别使用 [mux](https://github.com/gorilla/mux) 和 [alice](https://github.com/justinas/alice)  实现，本项目选择了 [chi](https://github.com/go-chi/chi) 一个高效、简洁的扩展，同时具备了路由和中间件的功能。

