package config

import (
	"github.com/heroku/docker-registry-client/registry"
	"fmt"
	"encoding/json"
	"os"
)

const (
	URL string = "https://registry.hub.docker.com/";
	UserName string = "";
	Password string = "";

	lastV int = 5;
)

func LastVersions() {
	config := NewConfig();
	for _, br := range []string{"firefox", "chrome"} {

		browser := Browser{
			Image: "selenoid/" + br,
			Port:"4444",
			Path: "",
			Tmpfs:map[string]string{"/tmp":"size=512m"},
		}

		versions := versionsByBrowser(br, &browser);

		config.Browsers[br] = versions;

		fmt.Println(config)

	}
	writeConfig(config)

}

func browserVersions(browser string) []string {
	hub, err := registry.New(URL, UserName, Password)
	if err != nil {
		fmt.Println(err)
	}
	tags, err := hub.Tags("selenoid/" + browser)

	if err != nil {
		fmt.Println(err)
	}
	return tags;
}

func versionsByBrowser(browser string, br *Browser) Versions {
	tags := browserVersions(browser);
	var lastVersions []string;
	if len(tags) < lastV {
		lastVersions = tags;
	} else {
		lastVersions = tags[len(tags) - lastV:]
	}
	versions := Versions{
		lastVersions[len(lastVersions) - 1], //устанавливаем версию по дефолту
		make(map[string]*Browser),
	}

	for _, version := range lastVersions {
		versions.Versions[version] = br;
	}

	return versions;
}

func writeConfig(config *Config) {
	js, er := json.MarshalIndent(config, " ", "	");
	if er != nil {
		fmt.Println("Ошибка маршалинга")
	}

	file, er := os.Create("client/config.json");
	if er != nil {
		fmt.Println(er)
	}

	file.Write(js);
}