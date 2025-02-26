export interface VocabItem {
  english: string;
  german: string;
  parts: string[];
}

export interface VocabResponse {
  vocabs: VocabItem[];
}
