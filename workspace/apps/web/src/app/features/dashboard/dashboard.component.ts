import { CommonModule } from "@angular/common";
import { Component } from "@angular/core";
import { CompetitionRankingsComponent } from "../rankings/competition-rankings.component";

@Component({
  selector: "app-dashboard",
  standalone: true,
  imports: [CommonModule, CompetitionRankingsComponent],
  template: `<div class="container">
    <app-competition-rankings></app-competition-rankings>
  </div>`,
  styles: [
    `
      .container {
        padding: 0;
      }
    `,
  ],
})
export class DashboardComponent {}
