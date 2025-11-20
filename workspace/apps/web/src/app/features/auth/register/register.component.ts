import { Component } from "@angular/core";
import { CommonModule } from "@angular/common";
import { RouterLink } from "@angular/router";

@Component({
  selector: "app-register",
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `<div class="container">
    <h2>Register</h2>
    <p><a routerLink="/auth/login">Back to Login</a></p>
  </div>`,
  styles: [
    `
      .container {
        padding: 2rem;
        text-align: center;
      }
    `,
  ],
})
export class RegisterComponent {}
