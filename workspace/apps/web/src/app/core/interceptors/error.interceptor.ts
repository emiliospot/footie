import { HttpInterceptorFn, HttpErrorResponse } from "@angular/common/http";
import { inject } from "@angular/core";
import { Router } from "@angular/router";
import { catchError, throwError } from "rxjs";
import { AuthService } from "../services/auth.service";

export const errorInterceptor: HttpInterceptorFn = (req, next) => {
  const router = inject(Router);
  const authService = inject(AuthService);

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      if (error.status === 401) {
        // Unauthorized - logout and redirect to login
        authService.logout();
        router.navigate(["/auth/login"]).catch(() => {
          console.error("Navigation to login failed");
        });
      } else if (error.status === 403) {
        // Forbidden
        router.navigate(["/unauthorized"]).catch(() => {
          console.error("Navigation to unauthorized failed");
        });
      }

      return throwError(() => error);
    }),
  );
};
