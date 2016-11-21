package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_test",
)

go_binary(
    name = "kexpand",
    srcs = ["expand.go", "main.go", "root.go", "version.go"],
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@com_github_ghodss_yaml//:go_default_library",
        "@com_github_spf13_pflag//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_a8m_envsubst//:go_default_library",
        "//pkg/expand:go_default_library",
    ],
)
