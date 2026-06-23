import type { PsrtDocument } from '../types/document';

/** Parsed from /exemplo.psrt — bundled for web dev (no Go API on startup). */
export const exemploDocument: PsrtDocument = {
  pages: [
    {
      name: 'capa',
      style: { backGround: '#0F0F14' },
      imageUrl: 'file:///d%3A/projs/GO/psrt/grade_referencia.png',
      texts: [
        {
          x: 50,
          y: 50,
          width: 84,
          textSize: 1,
          style: { color: '#FFFFFF', fontWeight: '700' },
          index: 0,
          content: 'DailyTrack',
        },
        {
          x: 40,
          y: 18,
          width: 20,
          textSize: 1,
          style: { color: '#ff0000', backgroundColor: '#000' },
          index: 1,
          content: 'Recomendação do dia',
        },
        {
          x: 8,
          y: 78,
          width: 84,
          textSize: 1,
          style: { color: '#1DB954' },
          index: 2,
          content: 'Toque para ouvir',
        },
      ],
    },
    {
      name: 'mood-sexta',
      style: { backGround: '#1C1C26' },
      imageUrl: 'https://picsum.photos/seed/psrt-mood/1080/1920',
      texts: [
        {
          x: 50,
          y: 50,
          width: 80,
          textSize: 1,
          style: { color: '#FFFFFF' },
          index: 0,
          content: 'Energia de sexta',
        },
        {
          x: 10,
          y: 22,
          width: 80,
          textSize: 1,
          style: { color: '#9146FF' },
          index: 1,
          imageRef: 'https://picsum.photos/seed/psrt-inline/200/200',
          content: 'Artista — Nome da faixa',
        },
        {
          x: 10,
          y: 88,
          width: 80,
          textSize: 1,
          style: { color: '#A1A1AA' },
          index: 2,
          content: 'Spotify · Preview',
        },
      ],
    },
    {
      name: 'intro',
      style: {},
      imageUrl:
        'https://cdn.nexustoons.com/manga_pages/317/28583/page_1_d5bf6a1a.avif',
      texts: [
        {
          x: 22.6,
          y: 56.11,
          width: 77,
          textSize: 3,
          style: {
            color: '#f1f1f1ff',
            'font-weight': '600',
            background: '#000000ff',
            padding: '10px',
            'text-align': 'center',
          },
          index: 1,
          content: 'O Soberando supremo da eternidade',
        },
        {
          x: 12.65,
          y: 70.01,
          width: 25.5,
          textSize: 3,
          style: {
            color: '#f1f1f1ff',
            'font-weight': '600',
            background: '#000000ff',
            'text-align': 'center',
          },
          index: 2,
          content: 'Tradutor',
        },
        {
          x: 42,
          y: 69,
          width: 23,
          textSize: 4,
          style: { background: '#000', color: '#ffff' },
          index: 0,
          content: 'revisor',
        },
      ],
    },
  ],
  fonts: [
    'https://cdn.jsdelivr.net/npm/@fontsource/inter@5.0.18/files/inter-latin-400-normal.woff2',
    'https://cdn.jsdelivr.net/npm/@fontsource/roboto@5.0.8/files/roboto-latin-400-normal.woff2',
  ],
  consts: {
    accent_spotify: '#1DB954',
    shadow_card: '"boxShadow":"0 8px 24px rgba(0,0,0,0.35)"',
    text_secondary: '#A1A1AA',
  },
};
