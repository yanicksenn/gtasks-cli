load("@rules_python//python:defs.bzl", "py_runtime_pair")

py_runtime(
    name = "py3_runtime",
    interpreter_path = "/usr/local/bin/python3",
    python_version = "PY3",
)

py_runtime_pair(
    name = "py_runtime_pair",
    py2_runtime = None,
    py3_runtime = ":py3_runtime",
)

toolchain(
    name = "py_toolchain",
    toolchain = ":py_runtime_pair",
    toolchain_type = "@bazel_tools//tools/python:toolchain_type",
)