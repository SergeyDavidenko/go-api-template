### Base api template

First, install requirements
```console
$ pip install -r requirements.txt
```
On Ubuntu:
```console
sudo apt install cookiecutter
```

Finally, to run it based on this template, type:
```console
$ cookiecutter https://github.com/SergeyDavidenko/go-api-template.git
```


Answer the prompts with your own desired [options](). For example:
```console
full_name [Sergey Davidenko]: Sergey Davidenko
github_username [SergeyDavidenko]: SergeyDavidenko
app_name [mygolangproject]: echoserver
project_short_description [A Golang project.]: Awesome Echo Server
docker_hub_username [sergeydavidenko]: sergeydavidenko
Select docker_build_image_version:
1 - 1.24
2 - 1.22
3 - 1.21
Choose from 1, 2, 3 [1]: 1
Select db_type
    1 - postgres
    2 - mongodb
    3 - none
Choose from [1/2/3] [1]: 1
use_git [y]: y
```
