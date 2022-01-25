# Pluto
Pluto is a work in progress package manger. Unlike most package managers, It downloads the source of a package and builds from the source.

## Using
To install a package, run `pluto install <package>` you can find all availible packages with `pluto list`
If you want to remove a package, run `pluto remove <package>`

## Adding a package to the main package list
If you want to add a project to the main package list, follow the example format below and specify name, a working git url for the official source repository, and commands to build it. You are required to specify `build` and `remove`. If your project needs something to build, such as make, you can add that to `needs`. Pluto currently ignores dependencies listed in `needs`, but that will change when dependency management is implemented.

When specifying build commands, if your project can be installed and removed with just `make install` and `make uninstall`, you can just use those. If you want a more complicated installation or have a more complicated build process, you can put everything in a script in your project's repository and add something like `bash ./install.sh` 

Package info, using cowsay as an example
```json
{
  "name": "cowsay",
  "git": "https://github.com/cowsay-org/cowsay.git",
  "authors": [
    "Andrew Janke",
    "Tony Monroe"
  ],
  "needs": [
    "make"
  ],
  "build": [
    "sudo make install"
  ],
  "remove": [
    "sudo make uninstall"
  ]
}
```
