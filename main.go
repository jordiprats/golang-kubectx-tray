package main

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/getlantern/systray"
	"gopkg.in/yaml.v3"
)

// kind: Config
// current-context: default/api-whatever:6443/cluster-admin

type KubeConfigFile struct {
	Kind    string `yaml:"kind"`
	Context string `yaml:"current-context"`
}

// contexts:
// - context:
//     cluster: api-sand-emea-euc1-v3ff-p1-openshiftapps-com:6443
//     namespace: test-jordi
//     user: cluster-admin/api-sand-emea-euc1-v3ff-p1-openshiftapps-com:6443
//   name: cli-prk/api-sand-emea-euc1-v3ff-p1-openshiftapps-com:6443/cluster-admin

var kubeconfig_path string
var watcher *fsnotify.Watcher

func main() {
	systray.Run(onReady, onExit)
}

func setIcon() {
	kubeconfig_read, err := os.ReadFile(kubeconfig_path)
	if err != nil {
		fmt.Printf("Error reading kubeconfig file: %s\n", err)
		os.Exit(1)
	}

	kubeconfigData := KubeConfigFile{}
	err = yaml.Unmarshal(kubeconfig_read, &kubeconfigData)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}

	if kubeconfigData.Context == "" {
		fmt.Println("Error retrieving context")
		os.Exit(1)
	}

	fmt.Println(kubeconfigData.Context)
	chooseIcon(kubeconfigData.Context)
}

func onReady() {
	// systray.SetIcon(icon.Data)
	systray.SetTooltip("kubectl config current-context")

	systray.AddMenuItem("Quit", "Quit")
	// mQuit := systray.AddMenuItem("Quit", "Quit")
	// mQuit.SetIcon(icons.Kube)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error setting up fsnotify")
		os.Exit(1)
	}
	defer watcher.Close()

	home_dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error retrieving UserHomeDir")
		os.Exit(1)
	}

	kubeconfig_path = home_dir + "/.kube/config"

	setIcon()

	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Write {
					setIcon()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("error: %s\n", err.Error())
			}
		}

	}()

	err = watcher.Add(kubeconfig_path)
	if err != nil {
		fmt.Printf("Add failed: %s\n", err.Error())
		os.Exit(1)
	}
	<-done
}

func onExit() {
	// clean up here
	os.Exit(0)
}
