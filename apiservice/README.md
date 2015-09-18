#APIService
##说明
支持oneline 和 offline 两种模式。
offline时，或者称为MockServer，顾名思义，离线使用的后台服务，即模拟真实的server，当后台服务不可达，或者在开发模式下使用。
特别在两个团队分前后台开发的时候，把后台服务直接模拟出来，两个团队之间只进行API接口编程，这样开发效率也会有较好的提升。

offline时，会让其链接真正的service，但是要起到debug的作用。（under developement
）
##offline的框架介绍
1.	RESTFul资源服务器，作为离线使用需要完成正常的所有功能。
2. 	自带存储，需要存储可能重复使用的信息。保证功能的完备性

###Third party lib
1. [httprouter](http://godoc.org/github.com/julienschmidt/httprouter)
2. [boltDB](http://godoc.org/github.com/boltdb/bolt)