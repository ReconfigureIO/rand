math: a library of math functions for FPGAs
===========================================

[![Documentation](https://godoc.org/github.com/ReconfigureIO/math?status.svg)](http://godoc.org/github.com/ReconfigureIO/math)

This is a library of math functions, optimized for FPGAs running on the Reconfigure.io platform.

It currently provides random number generation.

Using in your kernels
---------------------

Reconfigure.io supports including vendor packages in your kernels. You can use your favorite Go dependency manager to add it to your kernel. We use [glide](https://github.com/Masterminds/glide) for our code.

```
$ glide create --non-interactive
[INFO]  Generating a YAML configuration file and guessing the dependencies
[INFO]  Attempting to import from other package managers (use --skip-import to skip)
[INFO]  Scanning code to look for dependencies
[INFO]  Writing configuration file (glide.yaml)
[INFO]  You can now edit the glide.yaml file. Consider:
[INFO]  --> Using versions and ranges. See https://glide.sh/docs/versions/
[INFO]  --> Adding additional metadata. See https://glide.sh/docs/glide.yaml/
[INFO]  --> Running the config-wizard command to improve the versions in your configuration
$ glide get github.com/ReconfigureIO/math
[INFO]  Preparing to install 1 package.
[INFO]  Attempting to get package github.com/ReconfigureIO/math
[INFO]  --> Gathering release information for github.com/ReconfigureIO/math
[INFO]  --> Adding github.com/ReconfigureIO/math to your configuration
[INFO]  Downloading dependencies. Please wait...
[INFO]  --> Fetching updates for github.com/ReconfigureIO/math
[INFO]  Resolving imports
[INFO]  Downloading dependencies. Please wait...
[INFO]  Exporting resolved dependencies...
[INFO]  --> Exporting github.com/ReconfigureIO/math
[INFO]  Replacing existing vendor dependencies
```

You should now see it in your `vendor` directory.

```
$ tree vendor
vendor
└── github.com
    └── ReconfigureIO
        └── math
            ├── CODE_OF_CONDUCT.md
            ├── glide.yaml
            ├── LICENSE
            ├── rand
            │   ├── cmd
            │   │   └── tables
            │   │       └── main.go
            │   ├── normal.go
            │   ├── normal_test.go
            │   └── rand.go
            └── README.md
```

Contributing
------------

Pull requests & issues are enthusiastically accepted!

By participating in this project you agree to follow our [Code of Conduct](CODE_OF_CONDUCT.md).
