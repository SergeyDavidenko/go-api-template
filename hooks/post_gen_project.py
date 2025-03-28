"""
Does the following:
1. Inits git if used
"""

from __future__ import print_function

import os
from subprocess import Popen

# Get the root project directory
PROJECT_DIRECTORY = os.path.realpath(os.path.curdir)


def remove_file(filename):
    """
    generic remove file from project dir
    """
    fullpath = os.path.join(PROJECT_DIRECTORY, filename)
    if os.path.exists(fullpath):
        os.remove(fullpath)


def init_git():
    """
    Initialises git on the new project folder
    """
    GIT_COMMANDS = [
        ["git", "init"],
        ["git", "add", "."],
        ["git", "commit", "-a", "-m", "Initial Commit."]
    ]

    for command in GIT_COMMANDS:
        git = Popen(command, cwd=PROJECT_DIRECTORY)
        git.wait()


print('{{ cookiecutter.docker_build_image_version }}')
tidy = Popen(
    ['go', 'mod', 'init', '{{cookiecutter.app_name}}'], cwd=PROJECT_DIRECTORY)
tidy.wait()
tidy = Popen(['go', 'mod', 'tidy'], cwd=PROJECT_DIRECTORY)
tidy.wait()

# Fmt golang files
fmt = Popen(['go', 'fmt', './...'], cwd=PROJECT_DIRECTORY)
fmt.wait()


# Initialize Git (should be run after all file have been modified or deleted)
if '{{ cookiecutter.use_git }}'.lower() == 'y':
    init_git()
else:
    remove_file(".gitignore")
