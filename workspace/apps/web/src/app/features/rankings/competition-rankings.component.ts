import { CommonModule } from "@angular/common";
import { AfterViewInit, Component, ElementRef, OnDestroy, OnInit, ViewChild, inject } from "@angular/core";
import { RankingsResponse } from "@core/models/rankings.model";
import { RankingsService } from "@core/services/rankings.service";

@Component({
  selector: "app-competition-rankings",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./competition-rankings.component.html",
  styleUrl: "./competition-rankings.component.scss",
})
export class CompetitionRankingsComponent implements OnInit, AfterViewInit, OnDestroy {
  @ViewChild('categoryTabsContainer', { static: false }) private categoryTabsContainer!: ElementRef<HTMLDivElement>;

  // Tab states
  public rankingType: "team" | "player" = "team";
  public selectedCategory = "attacking";

  // Data
  public rankingsData: RankingsResponse | null = null;
  public loading = false;
  public error: string | null = null;

  // Scroll state
  public canScrollLeft = false;
  public canScrollRight = false;

  // Category options
  public readonly categories = [
    { value: "attacking", label: "Attacking" },
    { value: "defending", label: "Defending" },
    { value: "distribution", label: "Distribution" },
    { value: "goalkeeper", label: "Goalkeeper" },
    { value: "insights", label: "Insights" },
  ];

  private readonly rankingsService = inject(RankingsService);
  private scrollCheckInterval?: ReturnType<typeof setInterval>;

  public ngOnInit(): void {
    this.loadRankings();
  }

  public ngAfterViewInit(): void {
    this.checkScrollButtons();
    // Check scroll buttons on window resize
    globalThis.addEventListener('resize', this.checkScrollButtons);
    // Check scroll buttons on scroll
    if (this.categoryTabsContainer) {
      this.categoryTabsContainer.nativeElement.addEventListener('scroll', this.checkScrollButtons);
    }
    // Periodically check scroll buttons (for dynamic content changes)
    this.scrollCheckInterval = globalThis.setInterval(() => {
      this.checkScrollButtons();
    }, 500);
  }

  public ngOnDestroy(): void {
    globalThis.removeEventListener('resize', this.checkScrollButtons);
    if (this.categoryTabsContainer) {
      this.categoryTabsContainer.nativeElement.removeEventListener('scroll', this.checkScrollButtons);
    }
    if (this.scrollCheckInterval) {
      clearInterval(this.scrollCheckInterval);
    }
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

  public loadRankings(): void {
    this.loading = true;
    this.error = null;
    // Keep previous data visible during loading to avoid showing "Loading..." text
    const previousData = this.rankingsData;

    this.rankingsService
      .getCompetitionRankings(this.rankingType, this.selectedCategory)
      .subscribe({
        next: (data) => {
          this.rankingsData = data;
          this.loading = false;
        },
        error: (err) => {
          this.error = "Failed to load rankings. Please try again.";
          this.loading = false;
          // Restore previous data on error
          if (previousData) {
            this.rankingsData = previousData;
          }
          console.error("Error loading rankings:", err);
        },
      });
  }

  public getInitials(playerName: string): string {
    const parts = playerName.split(" ");
    if (parts.length >= 2) {
      return (parts[0][0] + parts.at(-1)?.[0]).toUpperCase();
    }
    return playerName.substring(0, 2).toUpperCase();
  }

  public getAvatarColor(_name: string, index: number): string {
    const colors = ["#4CAF50", "#9C27B0", "#2196F3", "#FF9800", "#F44336"];
    return colors[index % colors.length];
  }

  public trackByCategory(
    _index: number,
    category: { value: string; label: string },
  ): string {
    return category.value;
  }

  public trackByRankingEntry(
    _index: number,
    entry: { rank: number; name: string },
  ): number {
    return entry.rank;
  }

  public trackByCategoryCard(
    _index: number,
    category: { title: string },
  ): string {
    return category.title;
  }

  public scrollTabs(direction: 'left' | 'right'): void {
    if (!this.categoryTabsContainer) {
      return;
    }

    const container = this.categoryTabsContainer.nativeElement;
    const scrollAmount = 200; // pixels to scroll

    if (direction === 'left') {
      container.scrollBy({ left: -scrollAmount, behavior: 'smooth' });
    } else {
      container.scrollBy({ left: scrollAmount, behavior: 'smooth' });
    }

    // Update button visibility after scroll
    setTimeout(() => this.checkScrollButtons(), 300);
  }

  private readonly checkScrollButtons = (): void => {
    if (!this.categoryTabsContainer) {
      return;
    }

    const container = this.categoryTabsContainer.nativeElement;
    const { scrollLeft, scrollWidth, clientWidth } = container;

    this.canScrollLeft = scrollLeft > 0;
    this.canScrollRight = scrollLeft < scrollWidth - clientWidth - 1; // -1 for rounding errors
  };
}
