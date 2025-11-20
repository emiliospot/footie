import { Component } from "@angular/core";
import { CommonModule } from "@angular/common";

@Component({
  selector: "app-admin-dashboard",
  standalone: true,
  imports: [CommonModule],
  template: `<div class="container"><h2>Admin Dashboard</h2></div>`,
  styles: [
    `
      .container {
        padding: 2rem;
      }
    `,
  ],
})
export class AdminDashboardComponent {}
