package main

type PageSettings struct {
	PageType               string `json:"pageType"`
	PageSizeAndOrientation string `json:"pageSizeAndOrientation"`
	TitleRectangle         string `json:"titleRectangle"`
	PageRectangle          string `json:"pageRectangle"`
}

type ColorModel struct {
	EnabledIcon            string `json:"enabledIcon"`
	DisabledIcon           string `json:"disabledIcon"`
	Paper                  string `json:"paper"`
	PaperBorder            string `json:"paperBorder"`
	Lassoing               string `json:"lassoing"`
	LongPressMark          string `json:"longPressMark"`
	DefaultGroupBorder     string `json:"defaultGroupBorder"`
	SelectedGroupBorder    string `json:"selectedGroupBorder"`
	HyperlinkedGroupBorder string `json:"hyperlinkedGroupBorder"`
	GroupShadow            string `json:"groupShadow"`
}

type PenModel struct {
	PenId int    `json:"penId"`
	Color string `json:"color"`
	Width int    `json:"width"`
}

type StrokeWidthModel struct {
	GroupBorder float64 `json:"groupBorder"`
	PaperBorder float64 `json:"paperBorder"`
}

type WidthModel struct {
	GroupResizeHandle   int `json:"groupResizeHandle"`
	GroupLinkHandle     int `json:"groupLinkHandle"`
	GroupRotationHandle int `json:"GroupRotationHandle"`
	EraserHandle        int `json:"eraserHandle"`
}

type ThemeSettings struct {
	ColorModel       ColorModel       `json:"colorModel"`
	PenModelList     []PenModel       `json:"penModelList"`
	StrokeWidthModel StrokeWidthModel `json:"strokeWidthModel"`
	WidthModel       WidthModel       `json:"widthModel"`
}
