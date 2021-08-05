package config

import (
	"os"
	"os/user"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type singleton struct {
	Viper viper.Viper
}

type ZarfFile struct {
	Url        string
	Shasum     string
	Target     string
	Executable bool
}

type ZarfMetatdata struct {
	Name        string
	Description string
	Version     string
}

const K3sChartPath = "/var/lib/rancher/k3s/server/static/charts"
const K3sManifestPath = "/var/lib/rancher/k3s/server/manifests"
const K3sImagePath = "/var/lib/rancher/k3s/agent/images"
const PackageInitName = "zarf-init.tar.zst"
const PackageApplianceName = "zarf-appliance-init.tar.zst"
const ZarfLocal = "zarf.localhost"

var instance *singleton
var once sync.Once

func getInstance() *singleton {
	once.Do(func() {
		instance = &singleton{Viper: *viper.New()}
		setupViper()
	})
	return instance
}

func IsZarfInitConfig() bool {
	var kind string
	getInstance().Viper.UnmarshalKey("kind", &kind)
	return strings.ToLower(kind) == "zarfinitconfig"
}

func IsApplianceMode() bool {
	var mode string
	getInstance().Viper.UnmarshalKey("mode", &mode)
	return strings.ToLower(mode) == "appliance"
}

func GetPackageName() string {
	return "zarf-package-" + GetMetaData().Name + ".tar.zst"
}

func GetMetaData() ZarfMetatdata {
	var metatdata ZarfMetatdata
	getInstance().Viper.UnmarshalKey("metadata.name", &metatdata.Name)
	getInstance().Viper.UnmarshalKey("metatdata.description", &metatdata.Description)
	getInstance().Viper.UnmarshalKey("metatdata.version", &metatdata.Version)
	return metatdata
}

func GetLocalFiles() []ZarfFile {
	var files []ZarfFile
	getInstance().Viper.UnmarshalKey("local.files", &files)
	return files
}

func GetLocalImages() []string {
	var images []string
	getInstance().Viper.UnmarshalKey("local.images", &images)
	return images
}

func GetLocalManifests() string {
	var manifests string
	getInstance().Viper.UnmarshalKey("local.manifestFolder", &manifests)
	return manifests
}

func GetRemoteImages() []string {
	var images []string
	getInstance().Viper.UnmarshalKey("remote.images", &images)
	return images
}

func GetRemoteRepos() []string {
	var repos []string
	getInstance().Viper.UnmarshalKey("remote.repos", &repos)
	return repos
}

func DynamicConfigLoad(path string) {
	logContext := logrus.WithField("path", path)
	logContext.Info("Loading dynamic config")
	getInstance().Viper.SetConfigFile(path)
	if err := getInstance().Viper.MergeInConfig(); err != nil {
		logContext.Warn("Unable to load the config file")
	}
}

func WriteConfig(path string) {
	now := time.Now()
	currentUser, userErr := user.Current()
	hostname, hostErr := os.Hostname()

	// Record the time of package creation
	getInstance().Viper.Set("package.timestamp", now.Format(time.RFC1123Z))
	if hostErr == nil {
		// Record the hostname of the package creation terminal
		getInstance().Viper.Set("package.terminal", hostname)
	}
	if userErr == nil {
		// Record the name of the user creating the package
		getInstance().Viper.Set("package.user", currentUser.Name)
	}
	// Save the parsed output to the config path given
	if err := getInstance().Viper.WriteConfigAs(path); err != nil {
		logrus.WithField("path", path).Fatal("Unable to write the config file")
	}
}

func setupViper() {
	instance.Viper.AddConfigPath(".")
	instance.Viper.SetConfigName("config")

	// If a config file is found, read it in.
	if err := instance.Viper.ReadInConfig(); err == nil {
		logrus.WithField("path", instance.Viper.ConfigFileUsed()).Info("Config file loaded")
	}
}
