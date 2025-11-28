package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RankingsHandler handles competition rankings endpoints.
type RankingsHandler struct {
	*BaseHandler
}

// NewRankingsHandler creates a new rankings handler.
func NewRankingsHandler(base *BaseHandler) *RankingsHandler {
	return &RankingsHandler{BaseHandler: base}
}

// RankingEntry represents a single ranking entry (team or player).
// Note: Field order optimized for memory alignment (pointers grouped together).
type RankingEntry struct {
	Rank        int     `json:"rank"`
	Name        string  `json:"name"`
	Team        string  `json:"team,omitempty"` // For player rankings
	Value       float64 `json:"value"`
	Logo        *string `json:"logo,omitempty"`        // Team logo URL
	Initials    *string `json:"initials,omitempty"`    // Player initials for avatar
	AvatarColor *string `json:"avatarColor,omitempty"` // Color for player avatar
}

// RankingCategory represents a ranking category (e.g., "xG", "Shots").
type RankingCategory struct {
	Title    string         `json:"title"`
	Unit     string         `json:"unit"` // e.g., "/90'"
	Rankings []RankingEntry `json:"rankings"`
}

// RankingsResponse represents the response for rankings.
type RankingsResponse struct {
	Type       string            `json:"type"`     // "team" or "player"
	Category   string            `json:"category"` // "attacking", "defending", etc.
	Categories []RankingCategory `json:"categories"`
}

// GetCompetitionRankings handles GET /api/v1/rankings.
// @Summary Get competition rankings.
// @Description Get team or player rankings for a competition.
// @Tags rankings
// @Accept json
// @Produce json
// @Param type query string false "Ranking type: team or player" default(team)
// @Param category query string false "Category: attacking, defending, distribution, goalkeeper, insights" default(attacking)
// @Param championship query string false "Championship name" default(Cyprus U19 League Division 1)
// @Param season query string false "Season" default(2025/2026)
// @Success 200 {object} RankingsResponse
// @Router /rankings [get]
func (h *RankingsHandler) GetCompetitionRankings(c *gin.Context) {
	rankingType := c.DefaultQuery("type", "team")
	category := c.DefaultQuery("category", "attacking")
	// Note: championship and season parameters are accepted but not used in mock data
	// They will be used when connecting to real database
	_ = c.DefaultQuery("championship", "Cyprus U19 League Division 1")
	_ = c.DefaultQuery("season", "2025/2026")

	var response RankingsResponse
	response.Type = rankingType
	response.Category = category

	if rankingType == "team" {
		response.Categories = h.getTeamRankings(category)
	} else {
		response.Categories = h.getPlayerRankings(category)
	}

	c.JSON(http.StatusOK, response)
}

// getTeamRankings returns mock team rankings data.
// Note: Magic numbers and string literals are intentional for mock data.
func (h *RankingsHandler) getTeamRankings(category string) []RankingCategory {
	switch category {
	case "attacking":
		return []RankingCategory{
			{
				Title: "xG - Expected Goals",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Anorthosis U19", Value: 2.42, Logo: stringPtr("/assets/teams/anorthosis.png")},
					{Rank: 2, Name: "Pafos U19", Value: 2.11, Logo: stringPtr("/assets/teams/pafos.png")},
					{Rank: 3, Name: "Olympiakos U19", Value: 2.04, Logo: stringPtr("/assets/teams/olympiakos.png")},
					{Rank: 4, Name: "Omonoia FC U19", Value: 2.03, Logo: stringPtr("/assets/teams/omonoia.png")},
					{Rank: 5, Name: "AEK U19", Value: 1.97, Logo: stringPtr("/assets/teams/aek.png")},
				},
			},
			{
				Title: "Shots",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "AEK U19", Value: 17.17, Logo: stringPtr("/assets/teams/aek.png")},
					{Rank: 2, Name: "Anorthosis U19", Value: 16.13, Logo: stringPtr("/assets/teams/anorthosis.png")},
					{Rank: 3, Name: "Pafos U19", Value: 15.63, Logo: stringPtr("/assets/teams/pafos.png")},
					{Rank: 4, Name: "Olympiakos U19", Value: 15.33, Logo: stringPtr("/assets/teams/olympiakos.png")},
					{Rank: 5, Name: "APOEL U19", Value: 15.0, Logo: stringPtr("/assets/teams/apoel.png")},
				},
			},
			{
				Title: "Crosses",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "AEK U19", Value: 14.33, Logo: stringPtr("/assets/teams/aek.png")},
					{Rank: 2, Name: "Anorthosis U19", Value: 9.88, Logo: stringPtr("/assets/teams/anorthosis.png")},
					{Rank: 3, Name: "Pafos U19", Value: 9.75, Logo: stringPtr("/assets/teams/pafos.png")},
					{Rank: 4, Name: "Omonoia FC U19", Value: 8.75, Logo: stringPtr("/assets/teams/omonoia.png")},
					{Rank: 5, Name: "APOEL U19", Value: 8.20, Logo: stringPtr("/assets/teams/apoel.png")},
				},
			},
			{
				Title: "1v1 Dribbles",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "AEL U19", Value: 22.71, Logo: stringPtr("/assets/teams/ael.png")},
					{Rank: 2, Name: "Karmiotissa U19", Value: 16.20, Logo: stringPtr("/assets/teams/karmiotissa.png")},
					{Rank: 3, Name: "Anorthosis U19", Value: 15.88, Logo: stringPtr("/assets/teams/anorthosis.png")},
					{Rank: 4, Name: "Olympiakos U19", Value: 14.67, Logo: stringPtr("/assets/teams/olympiakos.png")},
					{Rank: 5, Name: "Aris U19", Value: 13.71, Logo: stringPtr("/assets/teams/aris.png")},
				},
			},
			{
				Title: "Ball Carries",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Pafos U19", Value: 137.38, Logo: stringPtr("/assets/teams/pafos.png")},
					{Rank: 2, Name: "Olympiakos U19", Value: 124.67, Logo: stringPtr("/assets/teams/olympiakos.png")},
					{Rank: 3, Name: "AEK U19", Value: 123.0, Logo: stringPtr("/assets/teams/aek.png")},
					{Rank: 4, Name: "Anorthosis U19", Value: 112.13, Logo: stringPtr("/assets/teams/anorthosis.png")},
					{Rank: 5, Name: "APOEL U19", Value: 104.40, Logo: stringPtr("/assets/teams/apoel.png")},
				},
			},
			{
				Title: "Box Penetrations",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Olympiakos U19", Value: 14.67, Logo: stringPtr("/assets/teams/olympiakos.png")},
					{Rank: 2, Name: "Omonoia FC U19", Value: 13.13, Logo: stringPtr("/assets/teams/omonoia.png")},
					{Rank: 3, Name: "AEK U19", Value: 12.83, Logo: stringPtr("/assets/teams/aek.png")},
					{Rank: 4, Name: "Pafos U19", Value: 12.13, Logo: stringPtr("/assets/teams/pafos.png")},
					{Rank: 5, Name: "Anorthosis U19", Value: 11.25, Logo: stringPtr("/assets/teams/anorthosis.png")},
				},
			},
		}
	case "defending":
		return []RankingCategory{
			{
				Title: "Tackles Won",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "APOEL U19", Value: 18.5, Logo: stringPtr("/assets/teams/apoel.png")},
					{Rank: 2, Name: "Anorthosis U19", Value: 17.2, Logo: stringPtr("/assets/teams/anorthosis.png")},
					{Rank: 3, Name: "AEK U19", Value: 16.8, Logo: stringPtr("/assets/teams/aek.png")},
					{Rank: 4, Name: "Pafos U19", Value: 15.9, Logo: stringPtr("/assets/teams/pafos.png")},
					{Rank: 5, Name: "Olympiakos U19", Value: 15.1, Logo: stringPtr("/assets/teams/olympiakos.png")},
				},
			},
			{
				Title: "Interceptions",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Anorthosis U19", Value: 12.3, Logo: stringPtr("/assets/teams/anorthosis.png")},
					{Rank: 2, Name: "APOEL U19", Value: 11.8, Logo: stringPtr("/assets/teams/apoel.png")},
					{Rank: 3, Name: "AEK U19", Value: 11.2, Logo: stringPtr("/assets/teams/aek.png")},
					{Rank: 4, Name: "Pafos U19", Value: 10.9, Logo: stringPtr("/assets/teams/pafos.png")},
					{Rank: 5, Name: "Omonoia FC U19", Value: 10.5, Logo: stringPtr("/assets/teams/omonoia.png")},
				},
			},
		}
	case "distribution":
		return []RankingCategory{
			{
				Title: "Passes Completed",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Pafos U19", Value: 485.2, Logo: stringPtr("/assets/teams/pafos.png")},
					{Rank: 2, Name: "Olympiakos U19", Value: 472.8, Logo: stringPtr("/assets/teams/olympiakos.png")},
					{Rank: 3, Name: "AEK U19", Value: 468.5, Logo: stringPtr("/assets/teams/aek.png")},
					{Rank: 4, Name: "Anorthosis U19", Value: 455.3, Logo: stringPtr("/assets/teams/anorthosis.png")},
					{Rank: 5, Name: "APOEL U19", Value: 442.1, Logo: stringPtr("/assets/teams/apoel.png")},
				},
			},
		}
	case "goalkeeper":
		return []RankingCategory{
			{
				Title: "Saves",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "AEL U19", Value: 4.8, Logo: stringPtr("/assets/teams/ael.png")},
					{Rank: 2, Name: "Karmiotissa U19", Value: 4.5, Logo: stringPtr("/assets/teams/karmiotissa.png")},
					{Rank: 3, Name: "Aris U19", Value: 4.2, Logo: stringPtr("/assets/teams/aris.png")},
					{Rank: 4, Name: "Nea Salamina U19", Value: 4.0, Logo: stringPtr("/assets/teams/nea-salamina.png")},
					{Rank: 5, Name: "Ayia Napa U19", Value: 3.8, Logo: stringPtr("/assets/teams/ayia-napa.png")},
				},
			},
		}
	case "insights":
		return []RankingCategory{
			{
				Title: "Possession",
				Unit:  "%",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Pafos U19", Value: 58.3, Logo: stringPtr("/assets/teams/pafos.png")},
					{Rank: 2, Name: "Olympiakos U19", Value: 55.7, Logo: stringPtr("/assets/teams/olympiakos.png")},
					{Rank: 3, Name: "AEK U19", Value: 54.2, Logo: stringPtr("/assets/teams/aek.png")},
					{Rank: 4, Name: "Anorthosis U19", Value: 52.8, Logo: stringPtr("/assets/teams/anorthosis.png")},
					{Rank: 5, Name: "APOEL U19", Value: 51.5, Logo: stringPtr("/assets/teams/apoel.png")},
				},
			},
		}
	default:
		return []RankingCategory{}
	}
}

// getPlayerRankings returns mock player rankings data.
// Note: Magic numbers and string literals are intentional for mock data.
func (h *RankingsHandler) getPlayerRankings(category string) []RankingCategory {
	// Player-specific avatar colors (matching Figma design)
	playerColors := map[string]string{
		"Petros Ioannou":           "#1f2937",
		"Artemis Spanos":           "#9C27B0",
		"Kyriakos Epifaniou":      "#069669",
		"Antonis Kosionou":         "#9C27B0",
		"Marinos Petrou":           "#c2410c",
		"Konstantinos Poursaitidis": "#dc2626",
		"Christos Loukaidis":        "#c2410c",
		"Dimitris Ioannou":         "#c2410c",
		"Simonas Christofi":        "#dc2626",
		"Sotiris Panagi":           "#dc2727",
		"Glaukos Chatzimitsis":     "#c2410c",
		"Panagiotis Tsivikos":      "#069669",
		"Sotiris Panaghi":          "#9333ea",
		"Alexandros Efstathiou":    "#7c3aed",
		"Giorgos Lamprou":          "#c2410c",
		"Ioannis Efraimidis":       "#d97706",
		"Kyriakos Epifanou":        "#9C27B0",
		"Panagiotis Siderenios":    "#d97706",
		"Kosmas Ioannou":           "#9333ea",
		"Kyriakos Strouthou":       "#d9790a",
		"Frixos Michailidis":       "#4CAF50",
		"Orestis Hatzivassiliou":   "#2965eb",
		"Andreas Avraam":           "#db811d",
		"Curtis Junior Makosso":    "#c2410c",
		"Dimitris Petrou":          "#F44336",
		"Andreas Georgiou":         "#4CAF50",
		"Michalis Ioannou":         "#9C27B0",
		"Petros Christou":          "#2196F3",
		"Georgios Panayi":          "#FF9800",
		"Nikos Petrou":             "#F44336",
	}

	// Helper to get player color
	getColor := func(name string) string {
		return playerColors[name]
	}

	switch category {
	case "attacking":
		return []RankingCategory{
			{
				Title: "xG - Expected Goals",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Petros Ioannou", Team: "AEK U19", Value: 0.78, Initials: stringPtr("PI"), AvatarColor: stringPtr(getColor("Petros Ioannou"))},
					{Rank: 2, Name: "Artemis Spanos", Team: "Karmiotissa U19", Value: 0.72, Initials: stringPtr("AS"), AvatarColor: stringPtr(getColor("Artemis Spanos"))},
					{Rank: 3, Name: "Kyriakos Epifaniou", Team: "Nea Salamina U19", Value: 0.72, Initials: stringPtr("KE"), AvatarColor: stringPtr(getColor("Kyriakos Epifaniou"))},
					{Rank: 4, Name: "Antonis Kosionou", Team: "Ayia Napa U19", Value: 0.69, Initials: stringPtr("AK"), AvatarColor: stringPtr(getColor("Antonis Kosionou"))},
					{Rank: 5, Name: "Marinos Petrou", Team: "Anorthosis U19", Value: 0.62, Initials: stringPtr("MP"), AvatarColor: stringPtr(getColor("Marinos Petrou"))},
				},
			},
			{
				Title: "Shots",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Konstantinos Poursaitidis", Team: "APOEL U19", Value: 5.0, Initials: stringPtr("KP"), AvatarColor: stringPtr(getColor("Konstantinos Poursaitidis"))},
					{Rank: 2, Name: "Christos Loukaidis", Team: "AEK U19", Value: 5.0, Initials: stringPtr("CL"), AvatarColor: stringPtr(getColor("Christos Loukaidis"))},
					{Rank: 3, Name: "Marinos Petrou", Team: "Anorthosis U19", Value: 4.38, Initials: stringPtr("MP"), AvatarColor: stringPtr(getColor("Marinos Petrou"))},
					{Rank: 4, Name: "Dimitris Ioannou", Team: "APOEL U19", Value: 4.0, Initials: stringPtr("DI"), AvatarColor: stringPtr(getColor("Dimitris Ioannou"))},
					{Rank: 5, Name: "Simonas Christofi", Team: "AEL U19", Value: 3.67, Initials: stringPtr("SC"), AvatarColor: stringPtr(getColor("Simonas Christofi"))},
				},
			},
			{
				Title: "Crosses",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Konstantinos Poursaitidis", Team: "APOEL U19", Value: 6.0, Initials: stringPtr("KP"), AvatarColor: stringPtr(getColor("Konstantinos Poursaitidis"))},
					{Rank: 2, Name: "Sotiris Panagi", Team: "Anorthosis U19", Value: 6.0, Initials: stringPtr("SP"), AvatarColor: stringPtr(getColor("Sotiris Panagi"))},
					{Rank: 3, Name: "Glaukos Chatzimitsis", Team: "Pafos U19", Value: 5.0, Initials: stringPtr("GC"), AvatarColor: stringPtr(getColor("Glaukos Chatzimitsis"))},
					{Rank: 4, Name: "Panagiotis Tsivikos", Team: "Pafos U19", Value: 4.50, Initials: stringPtr("PT"), AvatarColor: stringPtr(getColor("Panagiotis Tsivikos"))},
					{Rank: 5, Name: "Sotiris Panaghi", Team: "Anorthosis U19", Value: 4.0, Initials: stringPtr("SP"), AvatarColor: stringPtr(getColor("Sotiris Panaghi"))},
				},
			},
			{
				Title: "1v1 Dribbles",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Alexandros Efstathiou", Team: "AEL U19", Value: 7.17, Initials: stringPtr("AE"), AvatarColor: stringPtr(getColor("Alexandros Efstathiou"))},
					{Rank: 2, Name: "Giorgos Lamprou", Team: "Karmiotissa U19", Value: 7.0, Initials: stringPtr("GL"), AvatarColor: stringPtr(getColor("Giorgos Lamprou"))},
					{Rank: 3, Name: "Ioannis Efraimidis", Team: "Aris U19", Value: 6.0, Initials: stringPtr("IE"), AvatarColor: stringPtr(getColor("Ioannis Efraimidis"))},
					{Rank: 4, Name: "Marinos Petrou", Team: "Anorthosis U19", Value: 5.83, Initials: stringPtr("MP"), AvatarColor: stringPtr(getColor("Marinos Petrou"))},
					{Rank: 5, Name: "Kyriakos Epifanou", Team: "Nea Salamina U19", Value: 5.5, Initials: stringPtr("KE"), AvatarColor: stringPtr(getColor("Kyriakos Epifanou"))},
				},
			},
			{
				Title: "Ball Carries",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Panagiotis Siderenios", Team: "Pafos U19", Value: 30.0, Initials: stringPtr("PS"), AvatarColor: stringPtr(getColor("Panagiotis Siderenios"))},
					{Rank: 2, Name: "Kosmas Ioannou", Team: "Pafos U19", Value: 26.33, Initials: stringPtr("KI"), AvatarColor: stringPtr(getColor("Kosmas Ioannou"))},
					{Rank: 3, Name: "Kosmas Ioannou", Team: "Pafos U19", Value: 25.25, Initials: stringPtr("KI"), AvatarColor: stringPtr(getColor("Kosmas Ioannou"))},
					{Rank: 4, Name: "Kyriakos Strouthou", Team: "AEK U19", Value: 23.50, Initials: stringPtr("KS"), AvatarColor: stringPtr(getColor("Kyriakos Strouthou"))},
					{Rank: 5, Name: "Frixos Michailidis", Team: "Olympiakos U19", Value: 22.0, Initials: stringPtr("FM"), AvatarColor: stringPtr(getColor("Frixos Michailidis"))},
				},
			},
			{
				Title: "Box Penetrations",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Christos Loukaidis", Team: "AEK U19", Value: 5.33, Initials: stringPtr("CL"), AvatarColor: stringPtr(getColor("Christos Loukaidis"))},
					{Rank: 2, Name: "Orestis Hatzivassiliou", Team: "Omonoia 29M U19", Value: 5.0, Initials: stringPtr("OH"), AvatarColor: stringPtr(getColor("Orestis Hatzivassiliou"))},
					{Rank: 3, Name: "Petros Ioannou", Team: "AEK U19", Value: 4.25, Initials: stringPtr("PI"), AvatarColor: stringPtr(getColor("Petros Ioannou"))},
					{Rank: 4, Name: "Andreas Avraam", Team: "Anorthosis U19", Value: 4.25, Initials: stringPtr("AA"), AvatarColor: stringPtr(getColor("Andreas Avraam"))},
					{Rank: 5, Name: "Curtis Junior Makosso", Team: "Pafos U19", Value: 4.0, Initials: stringPtr("CJ"), AvatarColor: stringPtr(getColor("Curtis Junior Makosso"))},
				},
			},
		}
	case "defending":
		return []RankingCategory{
			{
				Title: "Tackles Won",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Dimitris Petrou", Team: "APOEL U19", Value: 4.2, Initials: stringPtr("DP"), AvatarColor: stringPtr(getColor("Dimitris Petrou"))},
					{Rank: 2, Name: "Andreas Georgiou", Team: "Anorthosis U19", Value: 3.9, Initials: stringPtr("AG"), AvatarColor: stringPtr(getColor("Andreas Georgiou"))},
					{Rank: 3, Name: "Michalis Ioannou", Team: "AEK U19", Value: 3.7, Initials: stringPtr("MI"), AvatarColor: stringPtr(getColor("Michalis Ioannou"))},
					{Rank: 4, Name: "Petros Christou", Team: "Pafos U19", Value: 3.5, Initials: stringPtr("PC"), AvatarColor: stringPtr(getColor("Petros Christou"))},
					{Rank: 5, Name: "Georgios Panayi", Team: "Olympiakos U19", Value: 3.3, Initials: stringPtr("GP"), AvatarColor: stringPtr(getColor("Georgios Panayi"))},
				},
			},
		}
	case "distribution":
		return []RankingCategory{
			{
				Title: "Passes Completed",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Panagiotis Siderenios", Team: "Pafos U19", Value: 65.2, Initials: stringPtr("PS"), AvatarColor: stringPtr(getColor("Panagiotis Siderenios"))},
					{Rank: 2, Name: "Kosmas Ioannou", Team: "Pafos U19", Value: 62.8, Initials: stringPtr("KI"), AvatarColor: stringPtr(getColor("Kosmas Ioannou"))},
					{Rank: 3, Name: "Kyriakos Strouthou", Team: "AEK U19", Value: 58.5, Initials: stringPtr("KS"), AvatarColor: stringPtr(getColor("Kyriakos Strouthou"))},
					{Rank: 4, Name: "Frixos Michailidis", Team: "Olympiakos U19", Value: 55.3, Initials: stringPtr("FM"), AvatarColor: stringPtr(getColor("Frixos Michailidis"))},
					{Rank: 5, Name: "Andreas Avraam", Team: "Anorthosis U19", Value: 52.1, Initials: stringPtr("AA"), AvatarColor: stringPtr(getColor("Andreas Avraam"))},
				},
			},
		}
	case "goalkeeper":
		return []RankingCategory{
			{
				Title: "Saves",
				Unit:  "/90'",
				Rankings: []RankingEntry{
					{Rank: 1, Name: "Nikos Petrou", Team: "AEL U19", Value: 4.8, Initials: stringPtr("NP"), AvatarColor: stringPtr(getColor("Nikos Petrou"))},
					{Rank: 2, Name: "Andreas Georgiou", Team: "Karmiotissa U19", Value: 4.5, Initials: stringPtr("AG"), AvatarColor: stringPtr(getColor("Andreas Georgiou"))},
					{Rank: 3, Name: "Michalis Ioannou", Team: "Aris U19", Value: 4.2, Initials: stringPtr("MI"), AvatarColor: stringPtr(getColor("Michalis Ioannou"))},
					{Rank: 4, Name: "Petros Christou", Team: "Nea Salamina U19", Value: 4.0, Initials: stringPtr("PC"), AvatarColor: stringPtr(getColor("Petros Christou"))},
					{Rank: 5, Name: "Georgios Panayi", Team: "Ayia Napa U19", Value: 3.8, Initials: stringPtr("GP"), AvatarColor: stringPtr(getColor("Georgios Panayi"))},
				},
			},
		}
	default:
		return []RankingCategory{}
	}
}

// stringPtr is a helper function to create a string pointer from a string value.
func stringPtr(s string) *string {
	return &s
}
