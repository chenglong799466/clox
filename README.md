# 以下是使用clox的步骤

以下是使用clox的步骤：

1. 配置 MySQL 连接：在 `./config/config.yml` 文件中配置 MySQL 连接信息（如果文件不存在，则创建一个）。你可以先使用本地的 MySQL 进行测试。
   在配置文件中，填写正确的 MySQL 主机、端口、用户名、密码和数据库名称。
   `config.yml` 文件内容示例：
   ```yaml
   # config.yml
   
   db_config:
     host: localhost
     port: 3306
     username: your_username
     password: your_password
     database: your_database
   ```

2. 运行方式：

    - 在编辑器（如 GoLand）中直接运行 `main` 函数：

      修改 `main` 函数中的 `tableName` 变量的 `default` 字段，将其设置为你想要生成代码的表名。

      ```
      tableName = flag.String("table", "default", "input generate table")
      ```

    - 通过命令行打包运行：

      使用 `go build` 命令生成可执行文件：

      ```bash
      go build -o <output_file> <entry_file_path>
      ```

      生成的可执行文件可以通过以下命令运行：

      ```bash
      ./<output_file> -table="tableName"
      ```

      将 `tableName` 替换为你想要生成代码的表名。

3. 生成的文件将保存在项目的根目录下。运行后，clox 将根据指定的表名生成相应的代码文件。

请确保在运行 clox 之前，已经正确配置了 MySQL 连接信息，并且能够成功连接到指定的数据库。