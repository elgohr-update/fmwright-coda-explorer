/*
 *    Copyright 2020 bitfly gmbh
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package handlers

import (
	"coda-explorer/db"
	"coda-explorer/services"
	"coda-explorer/templates"
	"coda-explorer/types"
	"coda-explorer/version"
	"encoding/json"
	"html/template"
	"net/http"
	"time"
)

var indexTemplate = template.Must(template.New("index").Funcs(templates.GetTemplateFuncs()).ParseFiles("templates/layout.html", "templates/index.html"))

// Index will return the main "index" page using a go template
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	indexPageData := services.LatestIndexPageData()
	data := &types.PageData{
		Meta: &types.Meta{
			Title:       "coda explorer",
			Description: "",
			Path:        "",
		},
		ShowSyncingMessage: indexPageData.Blocks[0].Ts.Before(time.Now().Add(time.Hour * -12)),
		Active:             "index",
		Data:               indexPageData,
		Version:            version.Version,
	}

	var stats []*types.Statistic
	err := db.DB.Select(&stats, "SELECT * FROM statistics WHERE value > 0 ORDER BY ts, indicator")
	if err != nil {
		logger.Errorf("error retrieving statistcs data for route %v: %v", r.URL.String(), err)
		http.Error(w, "Internal server error", 503)
		return
	}
	indexPageData.ChartData = stats

	err = indexTemplate.ExecuteTemplate(w, "layout", data)

	if err != nil {
		logger.Errorf("error executing template for %v route: %v", r.URL.String(), err)
		http.Error(w, "Internal server error", 503)
		return
	}
}

// IndexPageData will show the main "index" page in json format
func IndexPageData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=15, s-maxage=15")

	err := json.NewEncoder(w).Encode(services.LatestIndexPageData())

	if err != nil {
		logger.Errorf("error sending latest index page data: %v", err)
		http.Error(w, "Internal server error", 503)
		return
	}
}
