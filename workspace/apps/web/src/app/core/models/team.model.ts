import { Player } from "./player.model";

export interface Team {
  id: number;
  name: string;
  short_name: string;
  code: string;
  founded?: number;
  country: string;
  city?: string;
  stadium?: string;
  stadium_capacity?: number;
  logo?: string;
  colors?: string;
  website?: string;
  players?: Player[];
  created_at: string;
  updated_at: string;
}

export interface TeamStatistics {
  id: number;
  team_id: number;
  team?: Team;
  season: string;
  competition: string;
  matches_played: number;
  wins: number;
  draws: number;
  losses: number;
  points: number;
  goals_scored: number;
  goals_conceded: number;
  goal_difference: number;
  goals_per_match: number;
  clean_sheets: number;
  possession: number;
  pass_accuracy: number;
  shots_per_match: number;
  shots_on_target_percentage: number;
  home_wins: number;
  home_draws: number;
  home_losses: number;
  away_wins: number;
  away_draws: number;
  away_losses: number;
  yellow_cards: number;
  red_cards: number;
  current_form: string;
  position?: number;
  created_at: string;
  updated_at: string;
}
