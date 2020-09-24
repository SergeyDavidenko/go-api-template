"""
Does the following:
1. Inits git if used
"""

from __future__ import print_function
import os
import shutil
from subprocess import Popen

# Get the root project directory
PROJECT_DIRECTORY = os.path.realpath(os.path.curdir)


def remove_file(filename):
    """
    generic remove file from project dir
    """
    shutil.rmtree(os.path.join(
        PROJECT_DIRECTORY, "storage"
    ))


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


def remove_storage_files():
    """
    Removes files needed for storage
    """
    for filename in ["storage",]:
        os.remove(os.path.join(
            PROJECT_DIRECTORY, filename
        ))


# Initialize Git (should be run after all file have been modified or deleted)
if '{{ cookiecutter.use_git }}'.lower() == 'y':
    init_git()
else:
    remove_file(".gitignore")

# Remove storage
if '{{ cookiecutter.use_postgresql }}'.lower() != 'y':
    remove_storage_files()