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
$ cookiecutter https://github.com/SergeyDavidenko/base-api-template.git
```


Answer the prompts with your own desired [options](). For example:
```console
full_name [Sergey Davidenko]: Sergey Davidenko
github_username [SergeyDavidenko]: SergeyDavidenko
app_name [mygolangproject]: echoserver
project_short_description [A Golang project.]: Awesome Echo Server
docker_hub_username [sergeydavidenko]: sergeydavidenko
use_git [y]: y
Select docker_build_image_version:
1 - 1.15
2 - 1.14
3 - none
Choose from 1, 2, 3 [1]: 1
```