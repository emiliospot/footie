import { CommonModule } from "@angular/common";
import { Component } from "@angular/core";
import { FormsModule } from "@angular/forms";
import { RouterLink } from "@angular/router";

@Component({
  selector: "app-login",
  standalone: true,
  imports: [CommonModule, RouterLink, FormsModule],
  template: `
    <div class="auth-container">
      <div class="auth-card">
        <h1>Welcome to Footie</h1>
        <p class="subtitle">Football Analytics Platform</p>
        <form (ngSubmit)="onSubmit()" class="auth-form">
          <div class="form-group">
            <label for="email">Email</label>
            <input
              type="email"
              id="email"
              [(ngModel)]="email"
              name="email"
              required
            />
          </div>
          <div class="form-group">
            <label for="password">Password</label>
            <input
              type="password"
              id="password"
              [(ngModel)]="password"
              name="password"
              required
            />
          </div>
          <button type="submit" class="btn-primary">Login</button>
        </form>
        <p class="auth-link">
          Don't have an account?
          <a routerLink="/auth/register">Register here</a>
        </p>
      </div>
    </div>
  `,
  styles: [
    `
      .auth-container {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      }
      .auth-card {
        background: white;
        padding: 3rem;
        border-radius: 10px;
        box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
        width: 100%;
        max-width: 400px;
      }
      h1 {
        margin: 0 0 0.5rem;
        text-align: center;
        color: #333;
      }
      .subtitle {
        text-align: center;
        color: #666;
        margin-bottom: 2rem;
      }
      .auth-form {
        display: flex;
        flex-direction: column;
      }
      .form-group {
        margin-bottom: 1.5rem;
      }
      label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
      }
      input {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 1rem;
      }
      .btn-primary {
        padding: 0.75rem;
        background: #667eea;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 1rem;
      }
      .btn-primary:hover {
        background: #5568d3;
      }
      .auth-link {
        text-align: center;
        margin-top: 1.5rem;
      }
      a {
        color: #667eea;
        text-decoration: none;
      }
    `,
  ],
})
export class LoginComponent {
  public email: string = "";
  public password: string = "";

  public onSubmit(): void {
    // TODO: Implement actual login when auth endpoints are ready
    // After successful login, redirect to dashboard
    console.warn("Login functionality to be implemented");
  }
}
