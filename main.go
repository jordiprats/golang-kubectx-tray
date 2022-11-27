package main

import (
	"fmt"
	"os"
	"strings"
	"traykubectx/icons"

	"github.com/fsnotify/fsnotify"
	"github.com/getlantern/systray"
	"gopkg.in/yaml.v3"
)

// kind: Config
// current-context: default/api-whatever:6443/cluster-admin

type KubeConfigContext struct {
	Name    string `yaml:"name"`
	Context struct {
		Cluster   string `yaml:"cluster"`
		User      string `yaml:"user"`
		Namespace string `yaml:"namespace"`
	} `yaml:"context"`
}

type KubeConfigFile struct {
	Kind           string              `yaml:"kind"`
	CurrentContext string              `yaml:"current-context"`
	Contexts       []KubeConfigContext `yaml:"contexts"`
}

type KubeCtxTrayItem struct {
	Match string `yaml:"match"`
	Title string `yaml:"title"`
	Icon  string `yaml:"icon"`
}

type KubeCtxTrayConfig struct {
	Contexts []KubeCtxTrayItem `yaml:"contexts"`
}

var kubeconfig_path string
var kctconfig_path string
var watcher *fsnotify.Watcher
var menu_context *systray.MenuItem
var menu_namespace *systray.MenuItem

func main() {
	systray.Run(onReady, onExit)
}

func chooseIcon(context string) {

	kctconfig_read, err := os.ReadFile(kctconfig_path)
	if err != nil {
		fmt.Printf("Error reading kctconfig file: %s\n", err)
	}

	kctconfigData := KubeCtxTrayConfig{}
	err = yaml.Unmarshal(kctconfig_read, &kctconfigData)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}

	for _, item := range kctconfigData.Contexts {
		if strings.Contains(context, item.Match) {
			systray.SetTitle(" " + item.Title)
			if item.Icon == "green" {
				systray.SetIcon(icons.Green)
				return
			} else if item.Icon == "yellow" {
				systray.SetIcon(icons.Yellow)
				return
			} else if item.Icon == "red" {
				systray.SetIcon(icons.Red)
				return
			} else if item.Icon == "loki" {
				systray.SetIcon(icons.Loki)
				return
			} else if item.Icon == "odin" {
				systray.SetIcon(icons.Odin)
				return
			} else {
				systray.SetIcon(icons.Kube)
				return
			}
		}
	}

	systray.SetTitle("")
	systray.SetIcon(icons.Kube)
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

	if kubeconfigData.CurrentContext == "" {
		fmt.Println("Error retrieving context")
		os.Exit(1)
	}

	fmt.Println(kubeconfigData.CurrentContext)
	chooseIcon(kubeconfigData.CurrentContext)

	if menu_context == nil {
		menu_context = systray.AddMenuItem(kubeconfigData.CurrentContext, "context")
		menu_context.Disable()
	} else {
		menu_context.SetTitle(kubeconfigData.CurrentContext)
	}

	current_namespace := ""
	for _, item := range kubeconfigData.Contexts {
		if item.Name == kubeconfigData.CurrentContext {
			current_namespace = item.Context.Namespace
			break
		}
	}

	if menu_namespace == nil {
		menu_namespace = systray.AddMenuItem(current_namespace, "context")
		menu_namespace.Disable()
		systray.AddSeparator()
	} else {
		menu_namespace.SetTitle(current_namespace)
	}

}

func onReady() {
	// systray.SetIcon(icon.Data)
	systray.SetTooltip("kubectl config current-context")

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
	kctconfig_path = home_dir + "/.kube/kct-config"

	setIcon()

	mQuit := systray.AddMenuItem("Quit", "Quit")
	// mQuit.SetIcon(icons.Kube)
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

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
	fmt.Println("onExit")
}
