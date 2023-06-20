load("@bazel_skylib//lib:paths.bzl", "paths")

# Ref: https://www.stevenengelhardt.com/2023/03/03/practical-bazel-custom-bazel-make-variables/
def _impl(ctx):
    return [
        platform_common.TemplateVariableInfo({
            "LINUX_BUILD_DIRNAME": paths.dirname(ctx.expand_location("$(location :help.h)", ctx.attr.data)),
        }),
    ]

# dirname can be imported as a "toolchain" that will set the "make" variable
# LINUX_BUILD_DIRNAME to be the path to "this directory" (the directory
# containing help.h and the other linux-build headers).
dirname = rule(
    implementation = _impl,
    attrs = {
        "data": attr.label_list(allow_files = True,
                                default = [":help.h"]),
    },
)
