module example

go 1.20

//表明项目依赖于gee模块 并要求使用版本v0.0.0
require gee v0.0.0

//本项目在本地有一个gee目录 我们使用本地模板作为gee的版本 而不是从远程获取gee模块
replace gee => ./gee

