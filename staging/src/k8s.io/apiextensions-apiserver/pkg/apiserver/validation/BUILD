package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["validation.go"],
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/go-openapi/spec:go_default_library",
        "//vendor/github.com/go-openapi/validate:go_default_library",
        "//vendor/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions:go_default_library",
    ],
)
