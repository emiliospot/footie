import { Injectable } from "@angular/core";
import { HttpClient, HttpParams } from "@angular/common/http";
import { Observable } from "rxjs";
import { environment } from "@environments/environment";
import { Team, TeamStatistics } from "../models/team.model";
import { Player } from "../models/player.model";
import { PaginatedResponse } from "../models/api-response.model";

@Injectable({
  providedIn: "root",
})
export class TeamService {
  private readonly apiUrl = `${environment.apiUrl}/teams`;

  constructor(private readonly http: HttpClient) {}

  public getTeams(
    page: number = 1,
    limit: number = 20,
    country?: string,
  ): Observable<PaginatedResponse<Team>> {
    let params = new HttpParams()
      .set("page", page.toString())
      .set("limit", limit.toString());

    if (country) {
      params = params.set("country", country);
    }

    return this.http.get<PaginatedResponse<Team>>(this.apiUrl, { params });
  }

  public getTeam(id: number): Observable<Team> {
    return this.http.get<Team>(`${this.apiUrl}/${id}`);
  }

  public createTeam(team: Partial<Team>): Observable<Team> {
    return this.http.post<Team>(this.apiUrl, team);
  }

  public updateTeam(id: number, team: Partial<Team>): Observable<Team> {
    return this.http.put<Team>(`${this.apiUrl}/${id}`, team);
  }

  public deleteTeam(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  public getTeamPlayers(id: number): Observable<Player[]> {
    return this.http.get<Player[]>(`${this.apiUrl}/${id}/players`);
  }

  public getTeamStatistics(
    id: number,
    season?: string,
    competition?: string,
  ): Observable<TeamStatistics[]> {
    let params = new HttpParams();

    if (season) {
      params = params.set("season", season);
    }
    if (competition) {
      params = params.set("competition", competition);
    }

    return this.http.get<TeamStatistics[]>(`${this.apiUrl}/${id}/statistics`, {
      params,
    });
  }
}
