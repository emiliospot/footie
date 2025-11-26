import { Routes } from "@angular/router";

export const RANKINGS_ROUTES: Routes = [
  {
    path: "",
    loadComponent: () =>
      import("./competition-rankings.component").then(
        (m) => m.CompetitionRankingsComponent,
      ),
  },
];
