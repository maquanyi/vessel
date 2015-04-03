package cmd

import (
	"fmt"
	"net/http"

	"github.com/Unknwon/macaron"
	"github.com/codegangsta/cli"
	"github.com/macaron-contrib/binding"

	api "github.com/containerops/anchor"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/modules/log"
	"github.com/containerops/vessel/modules/setting"
	"github.com/containerops/vessel/modules/web"
)

var CmdWeb = cli.Command{
	Name:   "web",
	Usage:  "Start backend API server",
	Action: runWeb,
	Flags: []cli.Flag{
		cli.IntFlag{"port, p", 3000, "Port number to listen on", "VESSEL_WEB_PORT"},
	},
}

func runWeb(c *cli.Context) {
	if err := models.InitDb(); err != nil {
		log.Fatal("Fail to init DB: %v", err)
	}

	if c.IsSet("port") {
		setting.HTTPPort = c.Int("port")
	}

	bindIgnErr := binding.BindIgnErr

	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Renderer(macaron.RenderOptions{
		IndentJSON: !setting.ProdMode,
	}))
	m.Use(web.Contexter())

	group := func() {
		m.Group("/flows", func() {
			m.Combo("").Get(web.Flows).
				Post(bindIgnErr(api.CreateFlowOptions{}), web.CreateFlow)
			m.Combo("/:uuid").Get(web.GetFlow).
				Post(bindIgnErr(api.CreateFlowOptions{}), web.UpdateFlow).
				Delete(web.DeleteFlow)
		})

		m.Group("/pipelines", func() {
			m.Combo("").Get(web.Pipelines).
				Post(bindIgnErr(api.CreatePipelineOptions{}), web.CreatePipeline)
			m.Combo("/:uuid").Get(web.GetPipeline).
				Post(bindIgnErr(api.CreatePipelineOptions{}), web.UpdatePipeline).
				Delete(web.DeletePipeline)
		})

		m.Group("/stages", func() {
			m.Combo("").Get(web.Stages).
				Post(bindIgnErr(api.CreateStageOptions{}), web.CreateStage)
			m.Combo("/:uuid").Get(web.GetStage).
				Post(bindIgnErr(api.CreateStageOptions{}), web.UpdateStage).
				Delete(web.DeleteStage)
		})

		m.Group("/jobs", func() {
			m.Combo("").Get(web.Jobs).
				Post(bindIgnErr(api.CreateJobOptions{}), web.CreateJob)
			m.Combo("/:uuid").Get(web.GetJob).
				Post(bindIgnErr(api.CreateJobOptions{}), web.UpdateJob).
				Delete(web.DeleteJob)
		})

		m.Post("/build", web.Build)
	}
	m.Group("", group)
	m.Group("/v1", group)

	listenAddr := fmt.Sprintf("0.0.0.0:%d", setting.HTTPPort)
	log.Info("Vessel %s %s", setting.AppVer, listenAddr)
	if err := http.ListenAndServe(listenAddr, m); err != nil {
		log.Fatal("Fail to start web server: %v", err)
	}
}
