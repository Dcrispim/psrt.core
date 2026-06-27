//go:build js

package wasmbridge

import "syscall/js"

// Register binds all WASM handlers on the exports object.
func Register(exports js.Value) {
	exports.Set("parse", HandleParse())
	exports.Set("parseFast", HandleParseFast())
	exports.Set("loadSource", HandleLoadSource())
	exports.Set("convertLegacyDocument", HandleConvertLegacyDocument())
	exports.Set("stringify", HandleStringify())
	exports.Set("formatDocument", HandleFormatDocument())
	exports.Set("configure", HandleConfigure())

	exports.Set("addConst", HandleAddConst())
	exports.Set("removeConst", HandleRemoveConst())
	exports.Set("substituteConstReferences", HandleSubstituteConstReferences())
	exports.Set("revertConstReferences", HandleRevertConstReferences())

	exports.Set("addFont", HandleAddFont())
	exports.Set("removeFont", HandleRemoveFont())

	exports.Set("renamePage", HandleRenamePage())
	exports.Set("setPagePath", HandleSetPagePath())
	exports.Set("setPageStyle", HandleSetPageStyle())
	exports.Set("removePageStyleKey", HandleRemovePageStyleKey())
	exports.Set("movePage", HandleMovePage())
	exports.Set("addPage", HandleAddPage())
	exports.Set("removePage", HandleRemovePage())

	exports.Set("setTextStyle", HandleSetTextStyle())
	exports.Set("removeTextStyleKey", HandleRemoveTextStyleKey())
	exports.Set("setTextContent", HandleSetTextContent())
	exports.Set("addText", HandleAddText())
	exports.Set("removeText", HandleRemoveText())
	exports.Set("reorderTextRelative", HandleReorderTextRelative())
	exports.Set("reorderTextTo", HandleReorderTextTo())
	exports.Set("reorderTextByDelta", HandleReorderTextByDelta())
	exports.Set("setTextPosition", HandleSetTextPosition())
	exports.Set("nudgeTextPosition", HandleNudgeTextPosition())
	exports.Set("parseTextIndex", HandleParseTextIndex())

	exports.Set("setMaskPosition", HandleSetMaskPosition())
	exports.Set("addMask", HandleAddMask())
	exports.Set("removeMask", HandleRemoveMask())
	exports.Set("setMaskStyle", HandleSetMaskStyle())
	exports.Set("removeMaskStyleKey", HandleRemoveMaskStyleKey())

	exports.Set("setPathMaskPosition", HandleSetPathMaskPosition())
	exports.Set("addPathMask", HandleAddPathMask())
	exports.Set("removePathMask", HandleRemovePathMask())
	exports.Set("setPathMaskStyle", HandleSetPathMaskStyle())
	exports.Set("removePathMaskStyleKey", HandleRemovePathMaskStyleKey())
	exports.Set("setPathMaskPath", HandleSetPathMaskPath())

	exports.Set("setStyleKey", HandleSetStyleKey())
	exports.Set("removeStyleKey", HandleRemoveStyleKey())
	exports.Set("mergeStyle", HandleMergeStyle())

	exports.Set("findPage", HandleFindPage())
	exports.Set("findPageIndex", HandleFindPageIndex())
	exports.Set("findTextByIndex", HandleFindTextByIndex())
	exports.Set("findMaskByIndex", HandleFindMaskByIndex())
	exports.Set("findPathMaskByIndex", HandleFindPathMaskByIndex())

	exports.Set("compileToHtml", HandleCompileToHtml())
	exports.Set("compileToSvg", HandleCompileToSvg())

	exports.Set("adaptEntriesForWeb", HandleAdaptEntriesForWeb())
	exports.Set("mergePageDocumentPSRT", HandleMergePageDocumentPSRT())
	exports.Set("formatPageDocumentJSON", HandleFormatPageDocumentJSON())

	exports.Set("resolveDocument", HandleResolveDocument())
	exports.Set("resolveDocumentStrict", HandleResolveDocumentStrict())
}
