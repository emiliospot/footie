import { Component } from "@angular/core";
import { CommonModule } from "@angular/common";

@Component({
  selector: "app-team-list",
  standalone: true,
  imports: [CommonModule],
  template: `<div class="container"><h2>Teams</h2></div>`,
  styles: [
    `
      .container {
        padding: 2rem;
      }
    `,
  ],
})
export class TeamListComponent {}
