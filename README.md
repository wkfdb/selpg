# go 语言实现的selpg命令行程序


------

首先所有的参考资料基本都来自这个引导页。
[CLI命令行程序Go语言实现][1]
主要就是实现了一个selpg命令行程序。由于使用了Flag包，所以命令行采用等号的形式来传递参数。Flag包为解析命令行提供了十分方便的方法。
该代码依照[selpg程序][2]的设计方法来进行设计，实现了其中的所有功能。采用Go语言进行开发。


## 使用方式
进入到selpg.go的目录，然后运行`go build selpg.go`，然后运行生成的可运行程序就可以了。在运行时添加参数即可。下面介绍了实际的使用方式和产生的结果。


## 实验测试结果

首先为了方便起见在代码中设置了默认10行为一页。测试文档中有200行文本，每隔2行有一个换页符。

开发环境 Ubuntu 16.04

然后依次运行一下代码得到的结果如下：


1、`./selpg -s=1 -e=1 test.txt`
读取输入文件的第一页，将结果显示到屏幕上。由于输出方式为Println，自动添加换行，所以显示结果为每两行一个空行。空行实际上是换页符在ubuntu命令行中的表示形式。



2、`./selpg -s=1 -e=1 < test.txt`
以管道方式来打开输入文件，将第一页显示到屏幕上。



3、`./selpg -s=1 -e=1 test.txt >output.txt`
还是读取第一页，但是输出被shell内核重定向到output.txt文件中。
那个奇怪的东西其实是换页符。用Ubuntu的记事本打开就是那个样子。

4、`./selpg -s=1 -e=2 test.txt 2>error.txt`
selpg 将第 1 页到第 2 页写至标准输出（屏幕）；所有的错误消息被 shell／内核重定向至“error.txt”。

5、 `./selpg -s=1 -e=2 test.txt >output.txt 2>error.txt`
selpg 将第 1 页到第 2 页写至标准输出，标准输出被重定向至“output.txt”；selpg 写至标准错误的所有内容都被重定向至“error.txt”

6、 ```./selpg -s=1 -e=1 test.txt >output.txt 2>/dev/null```
selpg 将第 1 页写至标准输出，标准输出被重定向至“output.txt”；selpg 写至标准错误的所有内容都被重定向至 /dev/null（空设备），这意味着错误消息被丢弃了。设备文件 /dev/null 废弃所有写至它的输出，当从该设备文件读取时，会立即返回 EOF。

7、 ``` ./selpg -s=1 -e=1 test.txt | wc```
selpg 的标准输出透明地被 shell／内核重定向，成为“other_command”的标准输入，第 1 页被写至该标准输入。“other_command”的示例可以是 lp，它使输出在系统缺省打印机上打印。“other_command”的示例也可以 wc，它会显示选定范围的页中包含的行数、字数和字符数。“other_command”可以是任何其它能从其标准输入读取的命令。错误消息仍在屏幕显示。

8、 ```./selpg -s=2 -e=1 input.txt 2>error.txt | wc```
与上面的示例 7 相似，只有一点不同：错误消息被写至“error_file”

9、 ```./selpg -s=1 -e=2 -l=4 test.txt```
该命令将页长设置为 4 行，这样 selpg 就可以把输入当作被定界为该长度的页那样处理。第 1-2 页被写至 selpg 的标准输出（屏幕）。

10、 ```./selpg -s=1 -e=2 -f test.txt```
假定页由换页符定界。第 1-2 页被写至 selpg 的标准输出（屏幕）。

11、 ```./selpg -s=1 -e=2 -dlp1 test.txt```
第 1 页到第 2 页由管道输送至命令“lp -dlp1”，该命令将使输出在打印机 lp1 上打印。然后因为没有打印机，该命令运行时被内部替换为cat -n。试验结果见images里面的效果图。

--------

  [1]: https://pmlpml.github.io/ServiceComputingOnCloud/ex-cli-basic
  [2]: https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html


