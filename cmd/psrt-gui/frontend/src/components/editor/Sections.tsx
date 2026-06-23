import { useState, type ReactNode } from "react";
import s from "./sidebar.module.css";
import { IconChevron } from "./icons";

const STORAGE_PREFIX = "psrt-sidebar-section-";

function readStored(key: string, defaultOpen: boolean): boolean {
  try {
    const raw = localStorage.getItem(STORAGE_PREFIX + key);
    if (raw === null) return defaultOpen;
    return raw === "1";
  } catch {
    return defaultOpen;
  }
}

function writeStored(key: string, open: boolean): void {
  try {
    localStorage.setItem(STORAGE_PREFIX + key, open ? "1" : "0");
  } catch {
    /* ignore quota / private mode */
  }
}

export function Section({
  title,
  children,
  defaultOpen = true,
  storageKey,
}: {
  title: string;
  children: ReactNode;
  defaultOpen?: boolean;
  storageKey?: string;
}) {
  const key = storageKey ?? title.toLowerCase().replace(/\s+/g, "-");
  const [open, setOpen] = useState(() => readStored(key, defaultOpen));

  const toggle = () => {
    setOpen((prev) => {
      const next = !prev;
      writeStored(key, next);
      return next;
    });
  };

  return (
    <div className={s.section}>
      <button
        type="button"
        className={s.sectionHeader}
        onClick={toggle}
        aria-expanded={open}
      >
        <span>{title}</span>
        <span className={`${s.chev} ${open ? s.chevOpen : ""}`}>
          <IconChevron />
        </span>
      </button>
      {open && <div className={s.sectionBody}>{children}</div>}
    </div>
  );
}
