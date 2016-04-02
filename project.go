package dashboard

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	defaultProjectMap map[string]*Project
	defaultProjects   = []*Project{
		makeProject("bunto", "bunto/bunto", "ruby", "bunto"),
		makeProject("jemoji", "bunto/jemoji", "master", "jemoji"),
		makeProject("mercenary", "bunto/mercenary", "master", "mercenary"),
		makeProject("bunto-import", "bunto/bunto-import", "master", "bunto-import"),
		makeProject("bunto-feed", "bunto/bunto-feed", "master", "bunto-feed"),
		makeProject("bunto-sitemap", "bunto/bunto-sitemap", "master", "bunto-sitemap"),
		makeProject("bunto-mentions", "bunto/bunto-mentions", "master", "bunto-mentions"),
		makeProject("bunto-watch", "bunto/bunto-watch", "master", "bunto-watch"),
		makeProject("bunto-compose", "bunto/bunto-compose", "master", "bunto-compose"),
		makeProject("bunto-paginate", "bunto/bunto-paginate", "master", "bunto-paginate"),
		makeProject("bunto-gist", "bunto/bunto-gist", "master", "bunto-gist"),
		makeProject("bunto-coffeescript", "bunto/bunto-coffeescript", "master", "bunto-coffeescript"),
		makeProject("bunto-opal", "bunto/bunto-opal", "master", "bunto-opal"),
		makeProject("classifier-reborn", "bunto/classifier-reborn", "master", "classifier-reborn"),
		makeProject("bunto-sass-converter", "bunto/bunto-sass-converter", "master", "bunto-sass-converter"),
		makeProject("bunto-textile-converter", "bunto/bunto-textile-converter", "master", "bunto-textile-converter"),
		makeProject("bunto-redirect-from", "bunto/bunto-redirect-from", "master", "bunto-redirect-from"),
		makeProject("github-metadata", "bunto/github-metadata", "master", "bunto-github-metadata"),
		makeProject("plugins.buntorb", "bunto/plugins", "gh-pages", ""),
		makeProject("bunto docker", "bunto/docker", "", ""),
	}
)

func init() {
	go resetProjectsPeriodically()
}

func resetProjectsPeriodically() {
	for range time.Tick(time.Hour / 2) {
		log.Println("resetting projects' cache")
		resetProjects()
	}
}

func resetProjects() {
	for _, p := range defaultProjects {
		p.reset()
	}
}

type Project struct {
	Name    string `json:"name"`
	Nwo     string `json:"nwo"`
	Branch  string `json:"branch"`
	GemName string `json:"gem_name"`

	Gem     *RubyGem      `json:"gem"`
	Travis  *TravisReport `json:"travis"`
	GitHub  *GitHub       `json:"github"`
	fetched bool
}

func (p *Project) fetch() {
	if !p.fetched {
		rubyGemChan := rubygem(p.GemName)
		travisChan := travis(p.Nwo, p.Branch)
		githubChan := github(p.Nwo)
		p.Gem = <-rubyGemChan
		p.Travis = <-travisChan
		p.GitHub = <-githubChan
		p.fetched = true
	}
}

func (p *Project) reset() {
	p.fetched = false
	p.Gem = nil
	p.Travis = nil
	p.GitHub = nil
}

func buildProjectMap() {
	defaultProjectMap = map[string]*Project{}
	for _, p := range defaultProjects {
		defaultProjectMap[p.Name] = p
	}
}

func makeProject(name, nwo, branch, rubygem string) *Project {
	return &Project{
		Name:    name,
		Nwo:     nwo,
		Branch:  branch,
		GemName: rubygem,
	}
}

func getProject(name string) Project {
	if defaultProjectMap == nil {
		buildProjectMap()
	}

	if p, ok := defaultProjectMap[name]; ok {
		if !p.fetched {
			p.fetch()
		}
		return *p
	}
	panic(fmt.Sprintf("no project named '%s'", name))
}

func getAllProjects() []*Project {
	var wg sync.WaitGroup
	for _, p := range defaultProjects {
		wg.Add(1)
		go func(project *Project) {
			project.fetch()
			wg.Done()
		}(p)
	}
	wg.Wait()
	return defaultProjects
}
