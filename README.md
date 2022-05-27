# SDK

Go APIs to Nextmv's software development kit.

### Environment variables defining plugin to load

NEXTMV_LIBRARY_PATH

Library path to look for the plugin, defaults to ``~/.nextmv/lib``

NEXTMV_SDK_VERSION

Version of the plugin to use, defaults to no version 

NEXTMV_SDK_OS

Operating system binary to use, defaults to runtime.GOOS if not defined 

NEXTMV_SDK_ARCH

Architecture binary to use, defaults to runtime.GOARCH if not defined