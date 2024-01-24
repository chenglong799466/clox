# colx使用方法

以下是使用colx的步骤：

1. 配置MySQL连接：在`./config/config.yml`文件中配置MySQL连接信息。你可以先使用本地MySQL进行测试。  
在配置文件中，填写正确的MySQL主机、端口、用户名、密码和数据库名称。

2. 运行方式：

  - 在编辑器（如Goland）中直接运行`main`函数：

    修改`main`函数中的`tableName`变量的`default`字段，将其设置为你想要生成代码的表名。

    ```
    tableName = flag.String("table", "default", "input generate table")
    ```

  - 通过命令行打包运行：

    使用`go build`命令生成可执行文件：

    ```bash
    go build -o <output_file> <entry_file_path>
    ```

    生成的可执行文件可以通过以下命令运行：

    ```bash
    ./<output_file> -table="tableName"
    ```

    将`tableName`替换为你想要生成代码的表名。

3. 生成的文件在项目根目录下。运行后，colx将根据指定的表名生成相应的代码文件。 生成的文件将保存在项目的根目录下。