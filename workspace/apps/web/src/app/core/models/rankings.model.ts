export interface RankingEntry {
  rank: number;
  name: string;
  team?: string;
  value: number;
  logo?: string;
  initials?: string;
  avatarColor?: string;
}

export interface RankingCategory {
  title: string;
  unit: string;
  rankings: RankingEntry[];
}

export interface RankingsResponse {
  type: "team" | "player";
  category: string;
  categories: RankingCategory[];
}

