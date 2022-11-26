build:
	go build -o KubeCtxTray.app/Contents/MacOS/kubetray main.go
	bash buildicon.sh