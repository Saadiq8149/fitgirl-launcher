package handlers

import (
	"fitgirl-launcher/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type FitgirlScraperHandler struct {
	BaseUrl string
}

func CreateFitgirlScraperHandler() *FitgirlScraperHandler {
	return &FitgirlScraperHandler{
		BaseUrl: "https://fitgirl-repacks.site",
	}
}

func (fs *FitgirlScraperHandler) GetRepacks(query string, page int) models.FitgirlPage {
	c := colly.NewCollector()

	query = strings.ReplaceAll(query, " ", "+")

	repacks := []string{}
	totalPages := 0

	c.OnHTML("article.category-lossless-repack", func(e *colly.HTMLElement) {
		e.ForEach("h1.entry-title > a", func(_ int, e *colly.HTMLElement) {
			link := e.Attr("href")
			repacks = append(repacks, link)
		})
	})

	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		if e.Attr("class") == "page-numbers" {
			pageNum, _ := strconv.Atoi(e.Text)
			totalPages = max(totalPages, pageNum)
		}
	})

	c.Visit(fmt.Sprintf("%s/page/%d/?s=%s", fs.BaseUrl, page, query))
	return models.FitgirlPage{
		Results: repacks,
		Page:    page,
		Total:   totalPages,
	}
}

func (fs *FitgirlScraperHandler) GetRepackDetails(gameUrl string) models.FitgirlRepack {
	c := colly.NewCollector()

	repack := models.FitgirlRepack{}
	repack.Sources = []string{}

	c.OnHTML("title", func(e *colly.HTMLElement) {
		repack.Name = strings.TrimSpace(e.Text)
	})

	c.OnHTML("img.alignleft", func(e *colly.HTMLElement) {
		repack.CoverImage = e.Attr("src")
	})

	c.OnHTML("p", func(e *colly.HTMLElement) {
		style := e.Attr("style")

		if strings.Contains(style, "height") {
			genres := []string{}

			e.ForEach("a", func(_ int, e *colly.HTMLElement) {
				genres = append(genres, e.Text)
			})

			repack.Genres = genres
			sizeRegex := regexp.MustCompile("[0-9]+.?[0-9]? (GB|MB)")

			sizeMatches := sizeRegex.FindAllString(e.Text, -1)

			if len(sizeMatches) >= 2 {
				repack.OriginalSize = sizeMatches[0]
				repack.RepackSize = sizeMatches[1]
			}
		}
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if strings.Contains(link, "magnet:") {
			repack.Sources = append(repack.Sources, link)
		}
	})

	c.Visit(gameUrl)

	return repack
}

func (fs *FitgirlScraperHandler) GetPopularRepacks() models.PopularRepacks {
	c := colly.NewCollector()

	popularRepacks := models.PopularRepacks{}

	c.OnHTML("div.widget-grid-view-image", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, e *colly.HTMLElement) {
			popularRepack := models.PopularRepack{}
			popularRepack.Link = e.Attr("href")

			e.ForEach("img", func(_ int, e *colly.HTMLElement) {
				popularRepack.CoverImage = e.Attr("src")
			})

			popularRepacks.Repacks = append(popularRepacks.Repacks, popularRepack)

		})
	})

	c.Visit(fmt.Sprintf("%s/popular-repacks/", fs.BaseUrl))
	return popularRepacks
}

func (fs *FitgirlScraperHandler) GetLatestRepacks() models.PopularRepacks {
	c := colly.NewCollector()

	latestRepacks := models.PopularRepacks{}

	c.OnHTML("div.wplp-box-item", func(e *colly.HTMLElement) {
		latestRepack := models.PopularRepack{}

		e.ForEach("a", func(_ int, e *colly.HTMLElement) {
			latestRepack.Link = e.Attr("href")
		})

		e.ForEach("img", func(_ int, e *colly.HTMLElement) {
			latestRepack.CoverImage = e.Attr("src")
		})

		if latestRepack.Link != "" && latestRepack.CoverImage != "" {
			latestRepacks.Repacks = append(latestRepacks.Repacks, latestRepack)
		}
	})

	c.Visit(fs.BaseUrl)
	return latestRepacks
}
