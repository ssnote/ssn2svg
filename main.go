package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type CallbackLine func(string)

func readLines(callbackLine CallbackLine) error {
	bs, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var buf *bytes.Buffer = bytes.NewBuffer(bs)

	zr, err := gzip.NewReader(buf)
	if err != nil {
		return err
	}
	defer zr.Close()

	scanner := bufio.NewScanner(zr)
	for scanner.Scan() {
		line := scanner.Text()
		callbackLine(line)
	}

	return nil
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func createCurrentStyle(defaultStyle *Style) *Style {
	var styleFilePath string
	if len(os.Args) > 1 {
		styleFilePath = os.Args[1]
	}

	if Exists(styleFilePath) {
		style, err := createStyle(styleFilePath)
		if err == nil {
			return style
		} else {
			return defaultStyle
		}
	} else {
		return defaultStyle
	}
}

func toStrokeObject(addStrokeContents AddOrDelStrokeContents) *StrokeObject {
	strokeObject := new(StrokeObject)
	strokeObject.PageUuid = "0" // default value
	strokeObject.Uuid = addStrokeContents.Uuid
	strokeObject.Pts = addStrokeContents.Pts
	strokeObject.PenId = addStrokeContents.PenId
	strokeObject.LogicalStrokeWidth = 0 // default value

	return strokeObject
}

func copyStrokeObject(strokeObject StrokeObject) *StrokeObject {
	retVal := new(StrokeObject)
	retVal.PageUuid = strokeObject.PageUuid
	retVal.Uuid = strokeObject.Uuid
	retVal.Pts = strokeObject.Pts
	retVal.PenId = strokeObject.PenId
	retVal.LogicalStrokeWidth = strokeObject.LogicalStrokeWidth
	return retVal
}

func fixColor(argbColorStr string, callback func(Color)) {
	argb := strings.Split(argbColorStr, ",")
	r, err1 := strconv.ParseUint(argb[1], 10, 8)
	g, err2 := strconv.ParseUint(argb[2], 10, 8)
	b, err3 := strconv.ParseUint(argb[3], 10, 8)
	if err1 == nil && err2 == nil && err3 == nil {
		callback(Color{R: uint8(r), G: uint8(g), B: uint8(b)})
	}
}

type GroupAndStrokesMap map[string][]string

func appendStrokeUuidsIntoGroupAndStrokesMap(groupAndStrokesMap GroupAndStrokesMap, addGroupContents *AddGroupContents) {
	key := addGroupContents.Uuid
	groupAndStrokesMap[key] = []string{}
	for _, addStrokeContents := range addGroupContents.Strokes {
		groupAndStrokesMap[key] = append(groupAndStrokesMap[key], addStrokeContents.Uuid)
	}
}

func toStrokeColorList(pageContents *PageContents) []Color {
	defaultPenColor := Color{R: 8, G: 8, B: 8}

	strokeColorList := []Color{}

	colorMap := make(map[int]string)
	//widthMap := make(map[int]int) // TODO support pen width in the future.
	for _, penModel := range pageContents.ThemeSettings.PenModelList {
		colorMap[penModel.PenId] = penModel.Color
	}

	maxPenIdValue := 0
	for k, _ := range colorMap {
		if maxPenIdValue < k {
			maxPenIdValue = k
		}
	}

	for penId := 0; penId < (maxPenIdValue + 1); penId++ {
		penModelColor, ok := colorMap[penId]
		if ok {
			fixColor(penModelColor, func(color Color) {
				strokeColorList = append(strokeColorList, color)
			})
		} else {
			strokeColorList = append(strokeColorList, defaultPenColor)
		}
	}

	return strokeColorList
}

func main() {
	defaultStyle := createDefaultStyle()

	groupAndStrokesMap := make(GroupAndStrokesMap)

	db := new(Db)

	callbackLine := func(line string) {
		cmd0 := new(Cmd)
		err := json.Unmarshal([]byte(line), &cmd0)
		if err == nil {
			if cmd0.CmdName == "PUT_APP_NAME" {
				putAppNameContents := new(PutAppNameContents)
				json.Unmarshal([]byte(cmd0.Contents), &putAppNameContents)
			}

			if cmd0.CmdName == "ADD_PAGE" {
				addPageContents := new(AddPageContents)
				json.Unmarshal([]byte(cmd0.Contents), &addPageContents)

				pageContents := new(PageContents)
				json.Unmarshal([]byte(addPageContents.PageContents), &pageContents)

				defaultStyle.FillBackground = true
				fixColor(pageContents.ThemeSettings.ColorModel.Paper, func(color Color) {
					defaultStyle.BackgroundColor = color
				})

				strokeColorList := toStrokeColorList(pageContents)

				if len(strokeColorList) > 0 {
					defaultStyle.StrokeColorList = strokeColorList
				}
			}

			if cmd0.CmdName == "ADD_STROKE" {
				addStrokeContents := new(AddOrDelStrokeContents)
				json.Unmarshal([]byte(cmd0.Contents), &addStrokeContents)

				strokeObject := toStrokeObject(*addStrokeContents)
				db.Add(*strokeObject)
			}

			if cmd0.CmdName == "ADD_STROKES" {
				addStrokesContents := new(AddStrokesContents)
				json.Unmarshal([]byte(cmd0.Contents), &addStrokesContents)

				for _, addStrokeContents := range addStrokesContents.Strokes {
					strokeObject := toStrokeObject(addStrokeContents)
					db.Add(*strokeObject)
				}
			}

			if cmd0.CmdName == "ADD_GROUP" {
				addGroupContents := new(AddGroupContents)
				json.Unmarshal([]byte(cmd0.Contents), &addGroupContents)

				appendStrokeUuidsIntoGroupAndStrokesMap(groupAndStrokesMap, addGroupContents)

				for _, addStrokeContents := range addGroupContents.Strokes {
					strokeObject := toStrokeObject(addStrokeContents)
					db.Add(*strokeObject)
				}
			}

			if cmd0.CmdName == "ADD_GROUPS" {
				addGroupsContents := new(AddGroupsContents)
				json.Unmarshal([]byte(cmd0.Contents), &addGroupsContents)
				for _, addGroupContents := range addGroupsContents.Groups {
					appendStrokeUuidsIntoGroupAndStrokesMap(groupAndStrokesMap, &addGroupContents)

					for _, addStrokeContents := range addGroupContents.Strokes {
						strokeObject := toStrokeObject(addStrokeContents)
						db.Add(*strokeObject)
					}
				}
			}

			if cmd0.CmdName == "TRANSFORM_GROUP" {
				transformGroupContents := new(TransformGroupContents)
				json.Unmarshal([]byte(cmd0.Contents), &transformGroupContents)
				strokeUuids, ok := groupAndStrokesMap[transformGroupContents.Uuid]
				if ok {
					values := transformGroupContents.MatrixValues

					transformMatrix := matrix2d{
						values[0], values[1], values[2],
						values[3], values[4], values[5],
						values[6], values[7], values[8]}

					for _, strokeUuid := range strokeUuids {
						strokeObject, err := db.Get(strokeUuid)
						if err == nil {
							db.Remove(strokeUuid)

							aCopy := copyStrokeObject(strokeObject)
							aCopy.Pts = mapPoints(transformMatrix, strokeObject.Pts)
							db.Add(*aCopy)
						}
					}
				}
			}

			/*
				// For now, do not do anything for this command.
				if cmd0.CmdName == "DEL_GROUPS" {
					delGroupsContents := new(DelGroupsContents)
					json.Unmarshal([]byte(cmd0.Contents), &delGroupsContents)
				}
			*/

			if cmd0.CmdName == "DEL_STROKES" {
				delStrokesContents := new(DelStrokesContents)
				json.Unmarshal([]byte(cmd0.Contents), &delStrokesContents)

				for _, delStrokeContents := range delStrokesContents.Strokes {
					db.Remove(delStrokeContents.Uuid)
				}
			}

			/*
				// For now, do not do anything for this command.
				if cmd0.CmdName == "PUT_TITLE_TEXT" {
					putTitleTextContents := new(PutTitleTextContents)
					json.Unmarshal([]byte(cmd0.Contents), &putTitleTextContents)
				}

				// For now, do not do anything for this command.
				if cmd0.CmdName == "PUT_GROUP_TEXT" {
					putGroupTextContents := new(PutGroupTextContents)
					json.Unmarshal([]byte(cmd0.Contents), &putGroupTextContents)
				}
			*/
		}
	}

	err := readLines(callbackLine)
	if err == nil {
		svg := createSvg(db, createCurrentStyle(defaultStyle))
		fmt.Fprintln(os.Stdout, svg)
	}
}
