package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test"
)

go_library(
    name = "go_default_library",
    srcs = [
        "expand.go",
    ],
    deps = [
        "@com_github_golang_glog//:go_default_library",
    ],
)

go_test(
    name = "go_default_tests",
    srcs = ["expand_test.go"],
    library = ":go_default_library",
)
