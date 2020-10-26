# 支持子命令命令行程序支持包开发

## 概述

本仓库包括了程序包说明文档 `README.md` ，实现程序包的过程以及测试的的文档 `specification.md` ，`API.html` 是我生成的线下的 API 文件，方便查看函数的相关用法，其他的代码文件共同组成了一个名为 mycobra 的程序包，使用 mycobra 命令行生成一个简单的带子命令的命令行程序

## 获取包

输入以下的命令即可获取我实现的 mycobra 包

`go get github.com/hupf3/mycobra`

或者在 src 的相应目录下输入以下命令

`git clone https://github.com/hupf3/mycobra.git `

`go install`

## 使用说明

**实验环境**：

- 操作系统：mac os
- golang 版本: golang 1.14及以上

为了方便了解此程序包的使用，我实现了一个 pinfo 命令，用来获取个人信息的命令，并且封装成了一个包，输入以下命令即可获取：

`go get github.com/hupf3/pinfo`

或者在相应的 src 目录结构下是用以下命令：

`git clone https://github.com/hupf3/pinfo.git`

`go install`

执行完毕上面的命令后，在 `GOPATH` 路径中的 bin 文件夹下面会多出一个 pinfo 的可执行文件，即为成功

获取完上面的包后，打开命令行，输入 pinfo 测试是否成功安装 pinfo 命令

<img src="https://img-blog.csdnimg.cn/20201026213110118.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

如果出现如上图所示的 pinfo 命令的说明即可证明成功安装 pinfo 命令

然后按照说明可以得知，该命令是获取个人信息的命令，该命令下有三个字命令，分别为：age 获取年龄， id 获取学号， name获取名字

为了实现带参数的命令，我在 name 子命令中定义了三个参数，-f 获取姓，-g 获取名，-a 获取全称

我还实现了 help 相当于是全局的参数，当命令行输入command + help 时会显示该 command 的使用用法，示例如下：

<img src="https://img-blog.csdnimg.cn/20201026213549530.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom: 33%;" />

<img src="https://img-blog.csdnimg.cn/20201026213645207.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

测试完 help 功能后可以进行测试 age 命令是否实现：

<img src="https://img-blog.csdnimg.cn/20201026213807386.png#pic_center" alt="在这里插入图片描述" style="zoom:50%;" />

测试完 age 功能后可以进行测试 id 命令是否实现：

<img src="https://img-blog.csdnimg.cn/20201026213844794.png#pic_center" alt="在这里插入图片描述" style="zoom:50%;" />

测试完 id 功能后可以进行测试 name 各个命令参数是否实现：

<img src="https://img-blog.csdnimg.cn/2020102621390017.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom:50%;" />

至此，示例展示完毕，也是比较简单的进行了实现。也可以通过查看 API 文档进行具体的学习，API 文档生成的过程在 `specification.md` 中有具体的说明！

