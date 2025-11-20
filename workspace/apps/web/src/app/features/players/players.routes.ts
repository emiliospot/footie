import { Routes } from "@angular/router";

export const PLAYERS_ROUTES: Routes = [
  {
    path: "",
    loadComponent: () =>
      import("./player-list/player-list.component").then(
        (m) => m.PlayerListComponent,
      ),
  },
  {
    path: ":id",
    loadComponent: () =>
      import("./player-detail/player-detail.component").then(
        (m) => m.PlayerDetailComponent,
      ),
  },
];
