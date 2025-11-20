import { Injectable } from "@angular/core";
import { HttpClient, HttpParams } from "@angular/common/http";
import { Observable } from "rxjs";
import { environment } from "@environments/environment";
import { Match, MatchEvent } from "../models/match.model";
import { PaginatedResponse } from "../models/api-response.model";

@Injectable({
  providedIn: "root",
})
export class MatchService {
  private readonly apiUrl = `${environment.apiUrl}/matches`;

  constructor(private readonly http: HttpClient) {}

  public getMatches(
    page: number = 1,
    limit: number = 20,
    teamId?: number,
    season?: string,
    competition?: string,
    status?: string,
  ): Observable<PaginatedResponse<Match>> {
    let params = new HttpParams()
      .set("page", page.toString())
      .set("limit", limit.toString());

    if (teamId) {
      params = params.set("team_id", teamId.toString());
    }
    if (season) {
      params = params.set("season", season);
    }
    if (competition) {
      params = params.set("competition", competition);
    }
    if (status) {
      params = params.set("status", status);
    }

    return this.http.get<PaginatedResponse<Match>>(this.apiUrl, { params });
  }

  public getMatch(id: number): Observable<Match> {
    return this.http.get<Match>(`${this.apiUrl}/${id}`);
  }

  public createMatch(match: Partial<Match>): Observable<Match> {
    return this.http.post<Match>(this.apiUrl, match);
  }

  public updateMatch(id: number, match: Partial<Match>): Observable<Match> {
    return this.http.put<Match>(`${this.apiUrl}/${id}`, match);
  }

  public deleteMatch(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  public getMatchEvents(id: number): Observable<MatchEvent[]> {
    return this.http.get<MatchEvent[]>(`${this.apiUrl}/${id}/events`);
  }

  public createMatchEvent(
    matchId: number,
    event: Partial<MatchEvent>,
  ): Observable<MatchEvent> {
    return this.http.post<MatchEvent>(
      `${this.apiUrl}/${matchId}/events`,
      event,
    );
  }
}
