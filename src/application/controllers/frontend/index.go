package frontend

import (
	"os"
	"path/filepath"

	"github.com/xiusin/pine"
	"github.com/xiusin/pine/render/engine/pjet"
	"github.com/xiusin/pinecms/src/application/controllers"
)

func (c *IndexController) Index() {
	c.setTemplateData()
	indexPage := "editor.tpl"
	pageFilePath := GetStaticFile(indexPage)
	_ = os.MkdirAll(filepath.Dir(pageFilePath), os.ModePerm)
	f, err := os.OpenFile(pageFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		c.Logger().Error(err)
		return
	}
	defer f.Close()
	jet := pine.Make(controllers.ServiceJetEngine).(*pjet.PineJet)
	temp, err := jet.GetTemplate(template("index.jet"))
	if err != nil {
		c.Logger().Error(err)
		return
	}
	err = temp.Execute(f, viewDataToJetMap(c.Render().GetViewData()), nil)
	if err != nil {
		c.Logger().Error(err)
		return
	}
	data, _ := os.ReadFile(pageFilePath)

	_ = c.Ctx().WriteHTMLBytes(data)
}
