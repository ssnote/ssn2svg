package main

type Cmd struct {
	CmdKind     int    `json:"cmdKind"`
	CmdName     string `json:"cmdName"`
	Contents    string `json:"contents"`
	CreatedTime int64  `json:"createdTime"`
}

type PutAppNameContents struct {
	AppName string `json:"appName"`
}

type AddPageContents struct {
	PageContents string `json:"pageContents"`
}

type PageContents struct {
	PageSettings  PageSettings  `json:"pageSettings"`
	ThemeSettings ThemeSettings `json:"themeSettings"`
}

type AddOrDelStrokeContents struct {
	Uuid       string    `json:"uuid"`
	OriginUuid string    `json:"originUuid"`
	Pts        []float64 `json:"pts"`
	PenId      int       `json:"penId"`
}

type AddStrokesContents struct {
	Strokes []AddOrDelStrokeContents `json:"strokes"`
}

type AddGroupContents struct {
	Uuid    string                   `json:"uuid"`
	Strokes []AddOrDelStrokeContents `json:"strokes"`
}

type AddGroupsContents struct {
	Groups []AddGroupContents `json:"groups"`
}

type DelStrokesContents struct {
	Strokes []AddOrDelStrokeContents `json:"strokes"`
}
type DelGroupsContents struct {
	Groups []DelStrokesContents `json:"groups"`
}

type TransformGroupContents struct {
	Uuid         string    `json:"uuid"`
	MatrixValues []float64 `json:"matrixValues"`
}

type PutTitleTextContents struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

type PutGroupTextContents struct {
	GroupUuid string `json:"groupUuid"`
	Lang      string `json:"lang"`
	Text      string `json:"text"`
}
