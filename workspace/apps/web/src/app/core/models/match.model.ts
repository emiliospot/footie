import { Team } from "./team.model";
import { Player } from "./player.model";

export interface Match {
  id: number;
  match_date: string;
  competition: string;
  season: string;
  round?: string;
  stadium?: string;
  attendance?: number;
  status: MatchStatus;
  referee?: string;
  home_team_id: number;
  home_team?: Team;
  home_team_score: number;
  away_team_id: number;
  away_team?: Team;
  away_team_score: number;
  events?: MatchEvent[];
  created_at: string;
  updated_at: string;
}

export enum MatchStatus {
  SCHEDULED = "scheduled",
  LIVE = "live",
  FINISHED = "finished",
  POSTPONED = "postponed",
  CANCELLED = "cancelled",
}

export interface MatchEvent {
  id: number;
  match_id: number;
  match?: Match;
  minute: number;
  extra_minute?: number;
  event_type: EventType;
  description?: string;
  player_id?: number;
  player?: Player;
  team_id?: number;
  team?: Team;
  secondary_player_id?: number;
  secondary_player?: Player;
  position_x?: number;
  position_y?: number;
  metadata?: string;
  created_at: string;
  updated_at: string;
}

export enum EventType {
  GOAL = "goal",
  YELLOW_CARD = "yellow_card",
  RED_CARD = "red_card",
  SUBSTITUTION = "substitution",
  VAR_REVIEW = "var_review",
  PENALTY = "penalty",
  OWN_GOAL = "own_goal",
}
