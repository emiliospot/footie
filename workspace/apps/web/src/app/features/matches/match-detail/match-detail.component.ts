import { Component } from "@angular/core";
import { CommonModule } from "@angular/common";

@Component({
  selector: "app-match-detail",
  standalone: true,
  imports: [CommonModule],
  template: `<div class="container"><h2>Match Detail</h2></div>`,
  styles: [
    `
      .container {
        padding: 2rem;
      }
    `,
  ],
})
export class MatchDetailComponent {}
