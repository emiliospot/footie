import { Component } from "@angular/core";
import { CommonModule } from "@angular/common";

@Component({
  selector: "app-player-list",
  standalone: true,
  imports: [CommonModule],
  template: `<div class="container"><h2>Players</h2></div>`,
  styles: [
    `
      .container {
        padding: 2rem;
      }
    `,
  ],
})
export class PlayerListComponent {}
