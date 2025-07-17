load("@rules_proto//proto:defs.bzl", "proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "datetime-proto",
    srcs = ["datetime.proto"],
)

java_proto_library(
    name = "datetime-proto-java",
    deps = [":datetime-proto"],
)
