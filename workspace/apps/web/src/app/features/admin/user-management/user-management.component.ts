import { Component } from "@angular/core";
import { CommonModule } from "@angular/common";

@Component({
  selector: "app-user-management",
  standalone: true,
  imports: [CommonModule],
  template: `<div class="container"><h2>User Management</h2></div>`,
  styles: [
    `
      .container {
        padding: 2rem;
      }
    `,
  ],
})
export class UserManagementComponent {}
