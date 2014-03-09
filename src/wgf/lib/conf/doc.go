/*
package conf 是一个配置文件组件.


支持key=value与include、注释语法。
如以下配置文件：

wgf.conf

	//限制最大并发请求数，默认Index <-- 这一行是注释，只要用//开头即可
	wgf.defaultAction = Index
	wgf.viewDir = /home/deploy/mwiki/public/

	//包含其它配置文件
	//include可以帮助我们更合理的组织Conf配置
	include plugin/*.conf
*/
package conf
