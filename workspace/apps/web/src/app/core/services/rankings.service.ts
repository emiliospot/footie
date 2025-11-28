import { HttpClient, HttpParams } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import { environment } from "@environments/environment";
import { Observable } from "rxjs";
import { RankingsResponse } from "../models/rankings.model";

@Injectable({
  providedIn: "root",
})
export class RankingsService {
  private readonly apiUrl = `${environment.apiUrl}/rankings`;
  private readonly http = inject(HttpClient);

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
