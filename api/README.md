# Proto编译

- `-I`：指定 import 路径，可以指定多个参数，编译时按顺序查找，不指定时默认查找当前目录。
- `-go_out`：golang编译支持，指定输出文件路径
- `-go_opt`：指定参数，比如`--go_opt=paths=source_relative`就是表明生成文件输出使用相对路径