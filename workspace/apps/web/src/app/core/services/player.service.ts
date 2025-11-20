import { Injectable } from "@angular/core";
import { HttpClient, HttpParams } from "@angular/common/http";
import { Observable } from "rxjs";
import { environment } from "@environments/environment";
import { Player, PlayerStatistics } from "../models/player.model";
import { PaginatedResponse } from "../models/api-response.model";

@Injectable({
  providedIn: "root",
})
export class PlayerService {
  private readonly apiUrl = `${environment.apiUrl}/players`;

  constructor(private readonly http: HttpClient) {}

  public getPlayers(
    page: number = 1,
    limit: number = 20,
    teamId?: number,
    position?: string,
  ): Observable<PaginatedResponse<Player>> {
    let params = new HttpParams()
      .set("page", page.toString())
      .set("limit", limit.toString());

    if (teamId) {
      params = params.set("team_id", teamId.toString());
    }
    if (position) {
      params = params.set("position", position);
    }

    return this.http.get<PaginatedResponse<Player>>(this.apiUrl, { params });
  }

  public getPlayer(id: number): Observable<Player> {
    return this.http.get<Player>(`${this.apiUrl}/${id}`);
  }

  public createPlayer(player: Partial<Player>): Observable<Player> {
    return this.http.post<Player>(this.apiUrl, player);
  }

  public updatePlayer(id: number, player: Partial<Player>): Observable<Player> {
    return this.http.put<Player>(`${this.apiUrl}/${id}`, player);
  }

  public deletePlayer(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  public getPlayerStatistics(
    id: number,
    season?: string,
    competition?: string,
  ): Observable<PlayerStatistics[]> {
    let params = new HttpParams();

    if (season) {
      params = params.set("season", season);
    }
    if (competition) {
      params = params.set("competition", competition);
    }

    return this.http.get<PlayerStatistics[]>(
      `${this.apiUrl}/${id}/statistics`,
      { params },
    );
  }
}
