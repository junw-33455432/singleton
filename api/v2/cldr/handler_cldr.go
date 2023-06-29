/*
 * Copyright 2022-2023 VMware, Inc.
 * SPDX-License-Identifier: EPL-2.0
 */

package cldr

import (
	"strings"

	"sgtnserver/api"
	"sgtnserver/internal/common"
	"sgtnserver/internal/logger"
	"sgtnserver/modules/cldr"
	"sgtnserver/modules/cldr/cldrservice"
	"sgtnserver/modules/cldr/coreutil"
	"sgtnserver/modules/cldr/localeutil"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigDefault

// GetPatternByLocale godoc
// @Summary Get pattern data
// @Description Get pattern data in a specified locale
// @tags formatting-pattern-api
// @Produce json
// @Param locale path string true "locale String. e.g. 'en-US'"
// @Param scope query string true "pattern category string, separated by commas. e.g. 'dates,numbers,currencies,plurals,measurements,dateFields'"
// @Param scopeFilter query string false "a string for filtering the pattern data, separated by comma and underline. e.g. 'dates_eras,dates_dateTimeFormats'"
// @Success 200 {object} api.Response "OK"
// @Success 206 {object} api.Response "Successful Partially"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /formatting/patterns/locales/{locale} [get]
func GetPatternByLocale(c *gin.Context) {
	locale := struct {
		Locale string `uri:"locale" binding:"locale"`
	}{}
	scope := PatternScope{}
	if err := api.ExtractParameters(c, &locale, &scope); err != nil {
		return
	}

	ctx := logger.NewContext(c, c.MustGet(api.LoggerKey))
	cldrLocale, dataMap, err := cldrservice.GetPatternByLocale(ctx, locale.Locale, scope.Scope, scope.ScopeFilter)
	var data interface{}
	if len(dataMap) > 0 {
		parts := strings.Split(cldrLocale, cldr.LocalePartSep)
		region := coreutil.ParseRegion(parts)
		if region == "" {
			region, _ = localeutil.GetLocaleDefaultRegion(ctx, cldrLocale)
		}
		data = PatternData{
			LocaleID:   cldrLocale,
			Language:   parts[0],
			Region:     region,
			Categories: dataMap,
		}
	}

	api.HandleResponse(c, data, err)
}

// GetPatternDataByLangReg godoc
// @Summary Get pattern data
// @Description Get pattern data with language, region
// @tags formatting-pattern-api
// @Produce json
// @Param language query string true "a string which represents language, e.g. en,en-US,pt,pt-BR,zh-Hans"
// @Param region query string true "a string which represents region, e.g. US,PT,CN"
// @Param scope query string true "pattern category string, separated by commas. e.g. 'dates,numbers,currencies,plurals,measurements,dateFields'"
// @Param scopeFilter query string false "a string for filtering the pattern data, separated by comma and underline. e.g. 'dates_eras,dates_dateTimeFormats'"
// @Success 200 {object} api.Response "OK"
// @Success 206 {object} api.Response "Successful Partially"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /formatting/patterns [get]
func GetPatternDataByLangReg(c *gin.Context) {
	params := PatternByLangRegReq{}
	if err := api.ExtractParameters(c, nil, &params); err != nil {
		return
	}

	dataMap, localeToSet, mError := cldrservice.GetPatternByLangReg(logger.NewContext(c, c.MustGet(api.LoggerKey)), params.Language, params.Region, params.Scope, params.ScopeFilter)
	var data interface{}
	if len(dataMap) > 0 {
		data = PatternData{
			LocaleID:   localeToSet,
			Language:   params.Language,
			Region:     params.Region,
			Categories: dataMap,
		}
	}

	api.HandleResponse(c, data, mError)
}

// GetRegionListOfLanguages godoc
// @Summary Get region names
// @Description Get region names in a specified locale
// @Tags locale-api
// @Produce json
// @Param supportedLanguageList query string true "the supported language list, separated by commas. e.g. 'en,zh,ja' "
// @Success 200 {object} api.Response "OK"
// @Success 206 {object} api.Response "Successful Partially"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /locale/regionList [get]
func GetRegionListOfLanguages(c *gin.Context) {
	params := LocaleRegionsReq{}
	if err := api.ExtractParameters(c, nil, &params); err != nil {
		return
	}

	data, multiErr := localeutil.GetTerritoriesOfMultipleLocales(logger.NewContext(c, c.MustGet(api.LoggerKey)), strings.Split(params.Locales, common.ParamSep))
	for _, d := range data {
		unsortedMap := map[string]interface{}{}
		d.Territories.ToVal(&unsortedMap)
		sortedMap := treemap.NewWithStringComparator()
		for k, v := range unsortedMap {
			sortedMap.Put(k, v)
		}
		bts, _ := sortedMap.ToJSON()
		d.Territories = json.Get(bts)
	}
	api.HandleResponse(c, data, multiErr)
}

// GetTimeZoneNames godoc
// @Summary Get timezone names
// @Description Get time zone names in a specified locale
// @Tags formatting-api
// @Produce json
// @Param displayLanguage query string true "the display language. e.g. 'en'"
// @Param defaultTerritory query boolean false "a boolean value to get the default territory timezone name or not, default is true" default(true)
// @Success 200 {object} api.Response "OK"
// @Success 206 {object} api.Response "Successful Partially"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /formatting/date/timezoneNameList [get]
func GetTimeZoneNames(c *gin.Context) {
	params := struct {
		DisplayLanguage  string `form:"displayLanguage" binding:"language"`
		DefaultTerritory bool   `form:"defaultTerritory"`
	}{DefaultTerritory: true}
	if err := api.ExtractParameters(c, nil, &params); err != nil {
		return
	}

	data, err := localeutil.GetTimeZoneNames(logger.NewContext(c, c.MustGet(api.LoggerKey)), params.DisplayLanguage, params.DefaultTerritory)
	api.HandleResponse(c, data, err)
}