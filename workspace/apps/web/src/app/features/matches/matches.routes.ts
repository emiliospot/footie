import { Routes } from "@angular/router";

export const MATCHES_ROUTES: Routes = [
  {
    path: "",
    loadComponent: () =>
      import("./match-list/match-list.component").then(
        (m) => m.MatchListComponent,
      ),
  },
  {
    path: ":id",
    loadComponent: () =>
      import("./match-detail/match-detail.component").then(
        (m) => m.MatchDetailComponent,
      ),
  },
];
