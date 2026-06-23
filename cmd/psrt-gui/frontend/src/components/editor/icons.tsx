type Props = { size?: number };
const s = (n = 14) => ({ width: n, height: n, viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: 2, strokeLinecap: "round" as const, strokeLinejoin: "round" as const });

export const IconPlus = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M12 5v14M5 12h14"/></svg>);
export const IconDuplicate = ({ size: n }: Props = {}) => (
  <svg {...s(n)}>
    <rect x="9" y="9" width="13" height="13" rx="2" />
    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
  </svg>
);
export const IconTrash = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M6 6l1 14a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2l1-14"/></svg>);
export const IconChevron = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M9 6l6 6-6 6"/></svg>);
export const IconChevronDown = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M6 9l6 6 6-6"/></svg>);
export const IconEdit = ({ size: n }: Props = {}) => (
  <svg {...s(n)}>
    <path d="M12 20h9" />
    <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z" />
  </svg>
);

export const IconAlignLeft = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M4 6h16M4 12h10M4 18h14"/></svg>);
export const IconAlignCenter = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M4 6h16M7 12h10M5 18h14"/></svg>);
export const IconAlignRight = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M4 6h16M10 12h10M6 18h14"/></svg>);
export const IconAlignJustify = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M4 6h16M4 12h16M4 18h16"/></svg>);

export const IconVTop = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M4 4h16M12 20V8M8 12l4-4 4 4"/></svg>);
export const IconVMid = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M4 12h4M16 12h4M12 4v16M9 7l3-3 3 3M9 17l3 3 3-3"/></svg>);
export const IconVBot = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M4 20h16M12 4v12M8 12l4 4 4-4"/></svg>);

export const IconBold = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M7 4h6a4 4 0 0 1 0 8H7zM7 12h7a4 4 0 0 1 0 8H7z"/></svg>);
export const IconItalic = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M19 4h-9M14 20H5M15 4L9 20"/></svg>);
export const IconUnderline = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M6 4v7a6 6 0 0 0 12 0V4M4 21h16"/></svg>);
export const IconStrike = ({ size: n }: Props = {}) => (<svg {...s(n)}><path d="M4 12h16M17 6.5A4 4 0 0 0 13 4h-2a4 4 0 0 0 0 8M7 17.5A4 4 0 0 0 11 20h2a4 4 0 0 0 0-8"/></svg>);
