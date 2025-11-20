import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { Observable, BehaviorSubject, tap } from "rxjs";
import { environment } from "@environments/environment";
import {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
} from "../models/user.model";

@Injectable({
  providedIn: "root",
})
export class AuthService {
  public readonly currentUser$: Observable<User | null>;

  private readonly apiUrl = environment.apiUrl;
  private readonly currentUserSubject = new BehaviorSubject<User | null>(
    this.getUserFromStorage(),
  );

  constructor(private readonly http: HttpClient) {
    this.currentUser$ = this.currentUserSubject.asObservable();
  }

  public register(request: RegisterRequest): Observable<AuthResponse> {
    return this.http
      .post<AuthResponse>(`${this.apiUrl}/auth/register`, request)
      .pipe(tap((response: AuthResponse) => this.handleAuthResponse(response)));
  }

  public login(request: LoginRequest): Observable<AuthResponse> {
    return this.http
      .post<AuthResponse>(`${this.apiUrl}/auth/login`, request)
      .pipe(tap((response: AuthResponse) => this.handleAuthResponse(response)));
  }

  public logout(): void {
    localStorage.removeItem("token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user");
    this.currentUserSubject.next(null);
  }

  public refreshToken(): Observable<{ token: string }> {
    const refreshToken = this.getRefreshToken();
    return this.http
      .post<{ token: string }>(
        `${this.apiUrl}/auth/refresh`,
        {},
        {
          headers: { Authorization: `Bearer ${refreshToken}` },
        },
      )
      .pipe(
        tap((response: { token: string }) => {
          localStorage.setItem("token", response.token);
        }),
      );
  }

  public getToken(): string | null {
    return localStorage.getItem("token");
  }

  public getRefreshToken(): string | null {
    return localStorage.getItem("refresh_token");
  }

  public getCurrentUser(): User | null {
    return this.currentUserSubject.value;
  }

  public isAuthenticated(): boolean {
    return !!this.getToken();
  }

  public hasRole(role: string): boolean {
    const user = this.getCurrentUser();
    return user?.role === role || user?.role === "admin";
  }

  private handleAuthResponse(response: AuthResponse): void {
    localStorage.setItem("token", response.token);
    localStorage.setItem("refresh_token", response.refresh_token);
    localStorage.setItem("user", JSON.stringify(response.user));
    this.currentUserSubject.next(response.user);
  }

  private getUserFromStorage(): User | null {
    const userStr = localStorage.getItem("user");
    return userStr ? (JSON.parse(userStr) as User) : null;
  }
}
