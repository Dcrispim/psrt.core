import { CSSProperties, Dispatch, RefObject, SetStateAction } from 'react'

export interface PSRTEntry {
  x: number
  y: number
  size: number
  width: number
  style?: CSSProperties
  /** Original PSRT style JSON string (for Go adapter). */
  styleRaw?: string
  index?: number
  text: string
  /** == mask blocks: fixed height % (4th coord), not textSize. */
  maskHeight?: number
}

export type ChapterType = {
  id: number
  number: number
  name?: string
  titleId: number
  authorId: number
  psrtContent: string
  createdAt: string
  updatedAt: string
}

export interface PSRTEntryHeader {
  coordinates: { x: number; y: number; size: number; width: number }
  style: CSSProperties
  index: number
  pageLink?: string
}

export interface PSRTSection {
  title: string
  pageLink?: string
  pageStyle?: CSSProperties
  entries: PSRTEntry[]
}

export interface PSRTFile {
  sections: PSRTSection[]
}

export type subtitleContextType = {
  setSubtitle: Dispatch<SetStateAction<PSRTFile>>
  subtitle: PSRTFile
  currentPageIndex: number
  entryIndex: number
  zoom: number
  imageSize: {
    w?: number
    h?: number
  }
  currentPage: PSRTSection
  handleSubtitleChange: (
    _prop: keyof PSRTEntry | 'pageLink',
    _newSubtitleEntry: number | string | CSSProperties
  ) => void
  setZoom: Dispatch<SetStateAction<number>>
  setCurrentPageIndex: Dispatch<SetStateAction<number>>
  setEntryIndex: Dispatch<SetStateAction<number>>
  setImageSize: Dispatch<SetStateAction<{ w?: number; h?: number }>>
  currentEntry: PSRTEntry
  handleStyleChange: (_key: keyof CSSProperties, _value: string) => void
  newProperty: string
  setNewProperty: Dispatch<SetStateAction<string>>
  setNewValue: Dispatch<SetStateAction<string>>
  newValue: string
  addStyle: () => void
  popStyle: (_k: keyof CSSProperties) => void
  handleAddPage: () => void
  handleDownload: () => void
  handleRemovePage: () => void
  handleAddEntry: () => void
  handleRemoveEntry: () => void
  handleRemoveAllEntries: () => void
  style: CSSProperties | undefined
  uploadRef: RefObject<HTMLLabelElement>
  editingText: boolean
  setEditingText: Dispatch<SetStateAction<boolean>>
  fileUpload: () => void
  handleDuplicateEntry: () => void
}
