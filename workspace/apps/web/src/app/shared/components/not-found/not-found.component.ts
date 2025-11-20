import { Component } from "@angular/core";
import { RouterLink } from "@angular/router";
import { CommonModule } from "@angular/common";

@Component({
  selector: "app-not-found",
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="not-found-container">
      <h1>404</h1>
      <h2>Page Not Found</h2>
      <p>The page you're looking for doesn't exist.</p>
      <a routerLink="/dashboard" class="btn-home">Go to Dashboard</a>
    </div>
  `,
  styles: [
    `
      .not-found-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        min-height: 100vh;
        text-align: center;
        padding: 2rem;
      }

      h1 {
        font-size: 6rem;
        margin: 0;
        color: #3f51b5;
      }

      h2 {
        font-size: 2rem;
        margin: 1rem 0;
      }

      p {
        font-size: 1.2rem;
        margin: 1rem 0 2rem;
        color: #666;
      }

      .btn-home {
        padding: 0.75rem 2rem;
        background-color: #3f51b5;
        color: white;
        text-decoration: none;
        border-radius: 4px;
        transition: background-color 0.3s;
      }

      .btn-home:hover {
        background-color: #303f9f;
      }
    `,
  ],
})
export class NotFoundComponent {}
