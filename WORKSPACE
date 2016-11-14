git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.3.0",
)

load("@io_bazel_rules_go//go:def.bzl", "go_repositories", "new_go_repository")

go_repositories()

# for building docker base images
debs = (
    (
        "busybox_deb",
        "51651980a993b02c8dc663a5539a4d83704e56c2fed93dd8d1b2580e61319af5",
        "http://ftp.us.debian.org/debian/pool/main/b/busybox/busybox-static_1.22.0-19_amd64.deb",
    ),
    (
        "libc_deb",
        "ee4d9dea08728e2c2bbf43d819c3c7e61798245fab4b983ae910865980f791ad",
        "http://ftp.us.debian.org/debian/pool/main/g/glibc/libc6_2.19-18+deb8u6_amd64.deb",
    )
)

[http_file(
    name = name,
    sha256 = sha256,
    url = url,
) for name, sha256, url in debs]


new_go_repository(
    name = "com_github_golang_glog",
    importpath = "github.com/golang/glog",
    commit = "23def4e6c14b4da8ac2ed8007337bc5eb5007998",
)

new_go_repository(
    name = "com_github_spf13_cobra",
    importpath = "github.com/spf13/cobra",
    commit = "dbb7c2d02e284aa24423e3e6e105969dffe6e4d6",
)

new_go_repository(
    name = "com_github_spf13_pflag",
    importpath = "github.com/spf13/pflag",
    commit = "1560c1005499d61b80f865c04d39ca7505bf7f0b",
)

new_go_repository(
    name = "com_github_ghodss_yaml",
    importpath = "github.com/ghodss/yaml",
    commit = "bea76d6a4713e18b7f5321a2b020738552def3ea",
)


new_go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    commit = "a5b47d31c556af34a302ce5d659e6fea44d90de0",
)
