#A another tool for http loading test

##说明
1.针对url进行loading test，类似于ab test;
2.支持http、https，get 和 post；
3.可以配置header和post的body发送数据，简单易用；
4.充分使用golang 的 并发模型，高效稳定。输出的test数据可直接导入excel；
5.后续可以只用画图工具，生成chart，不过目前还没这个需要。

##用法
通过配置问题进行配置设置。
sample：
```
{
	"Duration":1800,
	"ThreadNum":[40, 80, 160],
	"TestRequest":[
		{
			"URL":"http://cnpvgvb1ep015.pvgl.sap.corp:49193/cart/add.json",
			"Method":"POST",
			"Body":"skuId=16&quantity=1&securityToken=3",
			"Header":{
				"tid":"1",
				"token":"1122"
			}
		},
		{
			"URL":"http://cnpvgvb1ep015.pvgl.sap.corp:49193/",
			"Method":"GET",
			"Body":"",
			"Header":{
				"tid":"1",
				"token":"1122"
			}
		}
	]
}
```
其中：
Duration是每个URL的测试时间，以秒为单位；ThreadNum为并发访问数，会对每个URL分别进行测试；
TestRequest为具体URL的配置，可以配置GET 或者 POST，同样的header、body可自己定义。

测试结果：

```
Threads	NumReqs	TPS	AvgResp(s)	MaxResp(s)	MinResp(s)	ErrNums	ErrCodes	URL
40	241	2.008	18.17	27.231	1.743	0	map[]	http://10.128.163.72/
80	1	0.008	11.562	11.562	11.562	0	map[]	http://10.128.163.72/
120	120	1	12.067	15.181	11.822	120	map[503:120]	http://10.128.163.72/
```

以上结果可直接导入到excel。

##实现说明：
借鉴[Go Concurrency Patterns](https://talks.golang.org/2012/concurrency.slide#1)，每个Thread都是一个goroutine，
使用fan-in模型，有个单独的routine来收集测试结果数据。总是实现上较为简洁，只使用了Go原生的lib。




