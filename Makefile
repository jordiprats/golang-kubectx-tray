build:
	go build -o KubeCtxTray.app/Contents/MacOS/kubetray main.go

icons:
	bash buildicon.sh