import { Routes } from "@angular/router";
import { authGuard, roleGuard } from "./core/guards/auth.guard";

export const routes: Routes = [
  {
    path: "",
    redirectTo: "/dashboard",
    pathMatch: "full",
  },
  {
    path: "auth",
    loadChildren: () =>
      import("./features/auth/auth.routes").then((m) => m.AUTH_ROUTES),
  },
  {
    path: "dashboard",
    loadChildren: () =>
      import("./features/dashboard/dashboard.routes").then(
        (m) => m.DASHBOARD_ROUTES,
      ),
    // TODO: Re-enable auth guard when auth endpoints are implemented
    // canActivate: [authGuard],
  },
  {
    path: "teams",
    loadChildren: () =>
      import("./features/teams/teams.routes").then((m) => m.TEAMS_ROUTES),
    canActivate: [authGuard],
  },
  {
    path: "players",
    loadChildren: () =>
      import("./features/players/players.routes").then((m) => m.PLAYERS_ROUTES),
    canActivate: [authGuard],
  },
  {
    path: "matches",
    loadChildren: () =>
      import("./features/matches/matches.routes").then((m) => m.MATCHES_ROUTES),
    canActivate: [authGuard],
  },
  {
    path: "rankings",
    loadChildren: () =>
      import("./features/rankings/rankings.routes").then(
        (m) => m.RANKINGS_ROUTES,
      ),
    // TODO: Re-enable auth guard when auth endpoints are implemented
    // canActivate: [authGuard],
  },
  {
    path: "admin",
    loadChildren: () =>
      import("./features/admin/admin.routes").then((m) => m.ADMIN_ROUTES),
    canActivate: [authGuard, roleGuard("admin")],
  },
  {
    path: "unauthorized",
    loadComponent: () =>
      import("./shared/components/unauthorized/unauthorized.component").then(
        (m) => m.UnauthorizedComponent,
      ),
  },
  {
    path: "**",
    loadComponent: () =>
      import("./shared/components/not-found/not-found.component").then(
        (m) => m.NotFoundComponent,
      ),
  },
];
