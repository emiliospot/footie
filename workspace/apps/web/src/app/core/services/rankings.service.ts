import { Injectable } from "@angular/core";
import { HttpClient, HttpParams } from "@angular/common/http";
import { Observable } from "rxjs";
import { environment } from "@environments/environment";
import { RankingsResponse } from "../models/rankings.model";

@Injectable({
  providedIn: "root",
})
export class RankingsService {
  private readonly apiUrl = `${environment.apiUrl}/rankings`;

  constructor(private readonly http: HttpClient) {}

  public getCompetitionRankings(
    type: "team" | "player" = "team",
    category: string = "attacking",
    championship?: string,
    season?: string,
  ): Observable<RankingsResponse> {
    let params = new HttpParams()
      .set("type", type)
      .set("category", category);

    if (championship) {
      params = params.set("championship", championship);
    }
    if (season) {
      params = params.set("season", season);
    }

    return this.http.get<RankingsResponse>(this.apiUrl, { params });
  }
}

