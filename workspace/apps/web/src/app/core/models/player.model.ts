import { Team } from "./team.model";

export interface Player {
  id: number;
  first_name: string;
  last_name: string;
  full_name: string;
  date_of_birth?: string;
  nationality?: string;
  height?: number;
  weight?: number;
  position: PlayerPosition;
  shirt_number?: number;
  preferred_foot?: PreferredFoot;
  photo?: string;
  team_id: number;
  team?: Team;
  created_at: string;
  updated_at: string;
}

export enum PlayerPosition {
  GK = "GK",
  DEF = "DEF",
  MID = "MID",
  FWD = "FWD",
}

export enum PreferredFoot {
  LEFT = "left",
  RIGHT = "right",
  BOTH = "both",
}

export interface PlayerStatistics {
  id: number;
  player_id: number;
  player?: Player;
  season: string;
  competition: string;
  matches_played: number;
  minutes_played: number;
  matches_started: number;
  sub_on: number;
  sub_off: number;
  goals: number;
  assists: number;
  shots_total: number;
  shots_on_target: number;
  shot_accuracy: number;
  goal_conversion: number;
  passes_total: number;
  passes_completed: number;
  pass_accuracy: number;
  key_passes: number;
  crosses: number;
  tackles: number;
  tackles_won: number;
  interceptions: number;
  clearances: number;
  blocked_shots: number;
  duels: number;
  duels_won: number;
  aerial_duels: number;
  aerial_duels_won: number;
  yellow_cards: number;
  red_cards: number;
  fouls: number;
  fouls_drawn: number;
  clean_sheets?: number;
  goals_conceded?: number;
  saves_total?: number;
  save_percentage?: number;
  penalties_saved?: number;
  created_at: string;
  updated_at: string;
}
