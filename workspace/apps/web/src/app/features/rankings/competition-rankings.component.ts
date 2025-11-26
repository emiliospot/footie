import { CommonModule } from "@angular/common";
import { Component, OnInit } from "@angular/core";
import { FormsModule } from "@angular/forms";
import { RankingsResponse } from "@core/models/rankings.model";
import { RankingsService } from "@core/services/rankings.service";

@Component({
  selector: "app-competition-rankings",
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: "./competition-rankings.component.html",
  styleUrl: "./competition-rankings.component.scss",
})
export class CompetitionRankingsComponent implements OnInit {
  // Filter options
  public selectedChampionship = "Cyprus U19 League Division 1";
  public selectedSeason = "2025/2026";

  // Tab states
  public rankingType: "team" | "player" = "team";
  public selectedCategory = "attacking";

  // Data
  public rankingsData: RankingsResponse | null = null;
  public loading = false;
  public error: string | null = null;

  // Category options
  public readonly categories = [
    { value: "attacking", label: "Attacking" },
    { value: "defending", label: "Defending" },
    { value: "distribution", label: "Distribution" },
    { value: "goalkeeper", label: "Goalkeeper" },
    { value: "insights", label: "Insights" },
  ];

  constructor(private readonly rankingsService: RankingsService) {}

  public ngOnInit(): void {
    this.loadRankings();
  }

  public onRankingTypeChange(type: "team" | "player"): void {
    this.rankingType = type;
    this.selectedCategory = "attacking"; // Reset to default category
    this.loadRankings();
  }

  public onCategoryChange(category: string): void {
    this.selectedCategory = category;
    this.loadRankings();
  }

  public onChampionshipChange(): void {
    this.loadRankings();
  }

  public onSeasonChange(): void {
    this.loadRankings();
  }

  public loadRankings(): void {
    this.loading = true;
    this.error = null;

    this.rankingsService
      .getCompetitionRankings(
        this.rankingType,
        this.selectedCategory,
        this.selectedChampionship,
        this.selectedSeason,
      )
      .subscribe({
        next: (data) => {
          this.rankingsData = data;
          this.loading = false;
        },
        error: (err) => {
          this.error = "Failed to load rankings. Please try again.";
          this.loading = false;
          console.error("Error loading rankings:", err);
        },
      });
  }

  public getInitials(playerName: string): string {
    const parts = playerName.split(" ");
    if (parts.length >= 2) {
      return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
    }
    return playerName.substring(0, 2).toUpperCase();
  }

  public getAvatarColor(_name: string, index: number): string {
    const colors = ["#4CAF50", "#9C27B0", "#2196F3", "#FF9800", "#F44336"];
    return colors[index % colors.length];
  }

  public trackByCategory(_index: number, category: { value: string; label: string }): string {
    return category.value;
  }

  public trackByRankingEntry(_index: number, entry: { rank: number; name: string }): number {
    return entry.rank;
  }

  public trackByCategoryCard(_index: number, category: { title: string }): string {
    return category.title;
  }
}
