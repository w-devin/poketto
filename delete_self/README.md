# 问题记录

## 1. golang的32位exe， 在win10_RS1及以上版本系统无法自删除

这是由于32位exe在windows上运行时， 开启了文件虚拟化导致的(失去NTFS环境)， 要关闭文件虚拟化， [可以通过给exe的应用程序清单指定 requestedExecutionLevel 节点实现](https://learn.microsoft.com/zh-cn/windows/win32/sbscs/application-manifests#:~:text=%E6%8C%87%E5%AE%9A%20requestedExecutionLevel%20%E8%8A%82%E7%82%B9%E5%B0%86%E7%A6%81%E7%94%A8%E6%96%87%E4%BB%B6%E5%92%8C%E6%B3%A8%E5%86%8C%E8%A1%A8%E8%99%9A%E6%8B%9F%E5%8C%96%E3%80%82%20%E5%A6%82%E6%9E%9C%E8%A6%81%E5%88%A9%E7%94%A8%E6%96%87%E4%BB%B6%E5%92%8C%E6%B3%A8%E5%86%8C%E8%A1%A8%E8%99%9A%E6%8B%9F%E5%8C%96%E5%AE%9E%E7%8E%B0%E5%90%91%E5%90%8E%E5%85%BC%E5%AE%B9%E6%80%A7%EF%BC%8C%E5%88%99%E7%9C%81%E7%95%A5%20requestedExecutionLevel%20%E8%8A%82%E7%82%B9%E3%80%82)

manifest.xml 内容
```xml
<?xml version='1.0' encoding='UTF-8' standalone='yes'?>
<assembly xmlns='urn:schemas-microsoft-com:asm.v1' manifestVersion='1.0'>
    <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
        <security>
	        <requestedPrivileges>
		        <requestedExecutionLevel level='asInvoker' uiAccess='false' />
			</requestedPrivileges>
		</security>
	</trustInfo>
</assembly>
```

golang可以使用这个项目 https://github.com/akavel/rsrc 把xml转为 `.syso` 文件, 放到项目根目录后， 即可在`go build` 时， 把应用程序清单编译到pe的 `.rsrc` 

```bash
# install rsrc
go get github.com/akavel/rsrc
go install github.com/akavel/rsrc

# use rsrc generate manifest.syso, 其实这里任意名称即可
rsrc  -manifest .\manifest.xml -o manifest.syso

# 这里要编译项目, 指定 .go 文件的话， syso不会添加. https://github.com/akavel/rsrc/issues/39
go build .
```