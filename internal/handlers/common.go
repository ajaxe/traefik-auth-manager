package handlers

var idParam = newAPIParam("id")
var appIDParam = newAPIParam("appId")

func newAPIParam(p string) apiParam {
	return apiParam(p)
}

type apiParam string

func (p *apiParam) String() string {
	return string(*p)
}
func (p *apiParam) Param() string {
	return ":" + p.String()
}
