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
		fmt.Printf("chooseIcon: Error reading kctconfig file: %s\n", err)
	}

	kctconfigData := KubeCtxTrayConfig{}
	err = yaml.Unmarshal(kctconfig_read, &kctconfigData)
	if err != nil {
		fmt.Printf("chooseIcon: error: %s\n", err.Error())
	}

	for _, item := range kctconfigData.Contexts {
		if strings.Contains(context, item.Match) {
			systray.SetTitle(" " + item.Title)

			trimed_icon := strings.TrimSpace(item.Icon)

			switch trimed_icon {
			case "green":
				systray.SetIcon(icons.Green)
				return
			case "yellow":
				systray.SetIcon(icons.Yellow)
				return
			case "red":
				systray.SetIcon(icons.Red)
				return
			case "loki":
				systray.SetIcon(icons.Loki)
				return
			case "odin":
				systray.SetIcon(icons.Odin)
				return
			case "greenproc":
				systray.SetIcon(icons.ProcGreen)
				return
			case "yellowproc":
				systray.SetIcon(icons.ProcYellow)
				return
			case "redproc":
				systray.SetIcon(icons.ProcRed)
				return
			default:
				systray.SetIcon(icons.Kube)
				return
			}
		}
	}

	systray.SetTitle("")
	systray.SetIcon(icons.Kube)
}

func setDisconnectedIcon() {
	systray.SetTitle("")
	systray.SetIcon(icons.KubeDisconnected)

	if menu_context == nil {
		menu_context = systray.AddMenuItem("", "context")
		menu_context.Disable()
	} else {
		menu_context.SetTitle("")
	}

	if menu_namespace == nil {
		menu_namespace = systray.AddMenuItem("", "context")
		menu_namespace.Disable()
		systray.AddSeparator()
	} else {
		menu_namespace.SetTitle("")
	}
}

func setIcon() {

	kubeconfig_read, err := os.ReadFile(kubeconfig_path)
	if err != nil {
		fmt.Printf("setIcon: Error reading kubeconfig file: %s\n", err)
		setDisconnectedIcon()
		return
	}

	kubeconfigData := KubeConfigFile{}
	err = yaml.Unmarshal(kubeconfig_read, &kubeconfigData)
	if err != nil {
		fmt.Printf("setIcon: error: %s\n", err.Error())
		setDisconnectedIcon()
		return
	}

	if kubeconfigData.CurrentContext == "" {
		setDisconnectedIcon()
		return
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
		menu_namespace = systray.AddMenuItem(current_namespace, "namespace")
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
		fmt.Println("onReady: Error setting up fsnotify")
		os.Exit(1)
	}
	defer watcher.Close()

	home_dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("onReady: Error retrieving UserHomeDir")
		os.Exit(1)
	}

	kubeconfig_path = home_dir + "/.kube/config"
	kctconfig_path = home_dir + "/.kube/kct-config"

	setIcon()

	systray.AddSeparator()

	mRefresh := systray.AddMenuItem("Refresh", "Refresh")
	go func() {
		<-mRefresh.ClickedCh
		fmt.Println("Refresh: Starting refresh")
		setIcon()
		fmt.Println("Refresh: Finished refreshing")
	}()

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit")
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("onReady: Requesting quit")
		systray.Quit()
		fmt.Println("onReady: Finished quitting")
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
				fmt.Printf("event: %s\n", event.Op.String())
				if event.Op == fsnotify.Write || event.Op == fsnotify.Rename {
					setIcon()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("onReady: error: %s\n", err.Error())
			}
		}

	}()

	err = watcher.Add(kubeconfig_path)
	if err != nil {
		fmt.Printf("onReady: Add failed: %s\n", err.Error())
		os.Exit(1)
	}
	err = watcher.Add(kctconfig_path)
	if err != nil {
		fmt.Printf("onReady: Add kct failed: %s\n", err.Error())
	}
	<-done
}

func onExit() {
	// clean up here
	fmt.Println("onExit")
}
