## About this folder

This folder contains some files that we use to build the `xenia-server` using `Jenkins` and other files like privacy policy and licenses.

PRs opened against the `xenia-server` repository will use the file called `Jenkinsfile.pr`

The `Dockerfile` in this folder (`Dockerfile.buildenv`) is the build environment for our current builds you can find the docker image to download [here](https://hub.docker.com/r/xenia/xenia-build-server/tags/) or build your own.



### Docker Image for building the Server

We have a docker image to build `xenia-server` and it is based on Go docker image.

In our Docker Hub Repository we have the following images:

- `xenia/xenia-build-server:dec-7-2018` which is based on Go 1.11 you can use for MM versions <= `5.9.0`
- `xenia/xenia-build-server:feb-28-2019` which is based on Go 1.12

