import { useEffect, useId, useRef, useState } from "react";
import type { InlineWrapKind } from "../../lib/inlineMarkup";
import s from "./sidebar.module.css";
import type { BorderRadius, TextBlock } from "./types";
import { ZERO_BORDER_RADIUS } from "./types";
import { isUniformBorderRadius, POSITION_PANEL_PROP_KEYS } from "../../lib/textBlockAdapter";

import { SliderField, ColorField, ToggleGroup, MultiToggle, SelectField, NumberField } from "./Fields";
import {
  IconAlignLeft, IconAlignCenter, IconAlignRight, IconAlignJustify,
  IconVTop, IconVMid, IconVBot,
  IconBold, IconItalic, IconUnderline, IconStrike,
  IconPlus, IconTrash, IconEdit, IconChevronDown,
  IconDuplicate,
} from "./icons";
import { Section } from "./Sections";
import { CSS_PROP_KEYS, resolveCssPropHandler } from "./cssPropHandlers";
import {
  getStyleHeightPercent,
  getStylePaddingPercent,
  removeStyleHeight,
  setStyleHeightPercent,
  setStylePaddingPercent,
} from "../../lib/textBoxHeight";

const DEFAULT_FONTS = [
  { value: "Inter", label: "Inter" },
  { value: "Roboto", label: "Roboto" },
  { value: "system-ui", label: "Sistema" },
  { value: "Georgia", label: "Georgia" },
  { value: "monospace", label: "Monospace" },
];

interface Props {
  blocks: TextBlock[];
  activeId: string;
  textContent: string;
  onTextContentChange: (value: string) => void;
  onSelect: (id: string) => void;
  onAdd: () => void;
  onDuplicate: () => void;
  onRemove: () => void;
  onChange: (updater: (b: TextBlock) => TextBlock) => void;
  onTypographyWrap?: (kind: InlineWrapKind, textarea: HTMLTextAreaElement) => boolean;
  fontOptions?: { value: string; label: string }[];
  emptyHint?: string;
}

export function PropertiesPanel({
  blocks,
  activeId,
  textContent,
  onTextContentChange,
  onSelect,
  onAdd,
  onDuplicate,
  onRemove,
  onChange,
  onTypographyWrap,
  fontOptions,
  emptyHint,
}: Props) {
  const contentTextareaRef = useRef<HTMLTextAreaElement>(null);
  const block = blocks.find((b) => b.id === activeId);
  const [draftProps, setDraftProps] = useState<{ key: string; value: string }[]>([]);
  const [expandedPropIndex, setExpandedPropIndex] = useState<number | null>(null);
  const [editingPropIndex, setEditingPropIndex] = useState<number | null>(null);
  const [radiusLinked, setRadiusLinked] = useState(true);
  const cssPropKeyListId = useId();



  const blockPropsKey = block ? JSON.stringify(block.props) : "";

  useEffect(() => {
    if (block) setDraftProps(block.props);
  }, [block?.id, blockPropsKey]);

  useEffect(() => {
    setExpandedPropIndex(null);
    setEditingPropIndex(null);
  }, [block?.id]);

  useEffect(() => {
    if (!block) return;
    setRadiusLinked(isUniformBorderRadius(block.borderRadius));
  }, [block?.id]);

  const setField = <K extends keyof TextBlock>(key: K, value: TextBlock[K]) =>
    onChange((b) => ({ ...b, [key]: value }));
  const setFont = <K extends keyof TextBlock["font"]>(key: K, value: TextBlock["font"][K]) =>
    onChange((b) => ({ ...b, font: { ...b.font, [key]: value } }));

  const tryInlineOrBlock = (kind: InlineWrapKind, blockToggle: () => void) => {
    const ta = contentTextareaRef.current;
    if (ta && ta.selectionStart !== ta.selectionEnd && onTypographyWrap?.(kind, ta)) {
      return;
    }
    blockToggle();
  };
  const setAlign = <K extends keyof TextBlock["align"]>(key: K, value: TextBlock["align"][K]) =>
    onChange((b) => ({ ...b, align: { ...b.align, [key]: value } }));

  const persistProps = (rows: { key: string; value: string }[]) => {
    setDraftProps(rows);
    onChange((b) => ({
      ...b,
      props: rows.filter((p) => p.key.trim() !== ""),
    }));
  };

  const updatePropRowDraft = (index: number, patch: Partial<{ key: string; value: string }>) => {
    const next = [...draftProps];
    next[index] = { ...next[index], ...patch };
    setDraftProps(next);
  };

  const commitPropKeys = () => {
    onChange((b) => ({
      ...b,
      props: draftProps.filter((p) => p.key.trim() !== ""),
    }));
  };

  const setBorderRadius = (patch: Partial<BorderRadius>) => {
    onChange((b) => ({
      ...b,
      borderRadius: { ...b.borderRadius, ...patch },
    }));
  };

  const setUniformRadius = (value: number) => {
    onChange((b) => ({
      ...b,
      borderRadius: {
        topLeft: value,
        topRight: value,
        bottomRight: value,
        bottomLeft: value,
      },
    }));
  };

  const resetBorderRadius = () => {
    onChange((b) => ({ ...b, borderRadius: { ...ZERO_BORDER_RADIUS } }));
  };

  const enableUnifiedRadius = () => {
    setRadiusLinked(true);
    if (!block) return;
    const { topLeft, topRight, bottomRight, bottomLeft } = block.borderRadius;
    const v = Math.max(topLeft, topRight, bottomRight, bottomLeft);
    if (!isUniformBorderRadius(block.borderRadius)) {
      setUniformRadius(v);
    }
  };

  return (
    <>
      <div className={s.textBar}>
        <select
          className={s.select}
          value={activeId}
          aria-label="Texto ativo"
          onChange={(e) => onSelect(e.target.value)}
        >
          <option value="">— selecionar —</option>
          {blocks.map((b, i) => (
            <option key={b.id} value={b.id}>#{i + 1} — {b.name}</option>
          ))}
        </select>
        <button className={`${s.iconBtn} ${s.primary}`} onClick={onAdd} title="Adicionar texto">
          <IconPlus />
        </button>
        <button className={`${s.iconBtn} ${s.primary}`} onClick={onDuplicate} title="Duplicar texto">
          <IconDuplicate />
        </button>
        <button className={`${s.iconBtn} ${s.danger}`} onClick={onRemove} title="Remover texto">
          <IconTrash />
        </button>
      </div>
      {!block ? (
        <p className={s.emptyHint}>
          {emptyHint ?? "Selecione um texto na lista ou clique no bloco na imagem."}
        </p>
      ) : null}

      {block ? (
      <>
      <Section title="Posição & Tamanho" storageKey="position">
        <div className={s.grid2}>
          <SliderField step={0.01} label="X" value={block.x} min={0} max={100} onChange={(v) => setField("x", v)} />
          <SliderField step={0.01} label="Y" value={block.y} min={0} max={100} onChange={(v) => setField("y", v)} />
          <SliderField step={0.01} label="Largura" value={block.width} min={1} max={100} onChange={(v) => setField("width", v)} />
          <SliderField
            step={0.01}
            label="Tam. base"
            value={block.size}
            min={0.5}
            max={20}
            onChange={(v) => setField("size", v)}
          />
          {(() => {
            const heightAuto = getStyleHeightPercent(block.props) === null;
            const heightVal = getStyleHeightPercent(block.props) ?? block.size;
            return (
              <div className={s.field}>
                <div className={s.labelRow}>
                  <label className={s.label}>Altura</label>
                  <label className={s.autoCheck}>
                    <input
                      type="checkbox"
                      checked={heightAuto}
                      onChange={(e) => {
                        const auto = e.target.checked;
                        onChange((b) => ({
                          ...b,
                          props: auto
                            ? removeStyleHeight(b.props)
                            : setStyleHeightPercent(
                                b.props,
                                getStyleHeightPercent(b.props) ?? b.size,
                              ),
                        }));
                      }}
                    />
                    auto
                  </label>
                </div>
                <div className={s.sliderRow}>
                  <input
                    type="range"
                    className={s.slider}
                    min={0}
                    max={100}
                    step={0.01}
                    value={heightVal}
                    disabled={heightAuto}
                    onChange={(e) =>
                      onChange((b) => ({
                        ...b,
                        props: setStyleHeightPercent(b.props, Number(e.target.value)),
                      }))
                    }
                  />
                  <input
                    type="number"
                    className={s.numInput}
                    min={0}
                    max={100}
                    step={0.01}
                    value={heightVal}
                    disabled={heightAuto}
                    onChange={(e) =>
                      onChange((b) => ({
                        ...b,
                        props: setStyleHeightPercent(b.props, Number(e.target.value)),
                      }))
                    }
                  />
                </div>
              </div>
            );
          })()}
          <SliderField
            step={0.01}
            label="Padding"
            value={getStylePaddingPercent(block.props) ?? 0}
            min={0}
            max={20}
            onChange={(v) =>
              onChange((b) => ({
                ...b,
                props: setStylePaddingPercent(b.props, v),
              }))
            }
          />
        </div>
      </Section>

      <Section title="Border radius" storageKey="border-radius">
        <div className={s.field}>
          <label className={s.label}>Cantos</label>
          <div className={s.toggleGroup} role="group">
            <button
              type="button"
              className={`${s.toggle} ${radiusLinked ? s.toggleActive : ""}`}
              aria-pressed={radiusLinked}
              onClick={enableUnifiedRadius}
            >
              Unificado
            </button>
            <button
              type="button"
              className={`${s.toggle} ${!radiusLinked ? s.toggleActive : ""}`}
              aria-pressed={!radiusLinked}
              onClick={() => setRadiusLinked(false)}
            >
              Individual
            </button>
          </div>
        </div>
        {radiusLinked ? (
          <SliderField
            label="Raio"
            value={block.borderRadius.topLeft}
            min={0}
            max={100}
            step={1}
            onChange={setUniformRadius}
          />
        ) : (
          <div className={s.grid2}>
            <SliderField
              label="Sup. esq."
              value={block.borderRadius.topLeft}
              min={0}
              max={100}
              step={1}
              onChange={(v) => setBorderRadius({ topLeft: v })}
            />
            <SliderField
              label="Sup. dir."
              value={block.borderRadius.topRight}
              min={0}
              max={100}
              step={1}
              onChange={(v) => setBorderRadius({ topRight: v })}
            />
            <SliderField
              label="Inf. dir."
              value={block.borderRadius.bottomRight}
              min={0}
              max={100}
              step={1}
              onChange={(v) => setBorderRadius({ bottomRight: v })}
            />
            <SliderField
              label="Inf. esq."
              value={block.borderRadius.bottomLeft}
              min={0}
              max={100}
              step={1}
              onChange={(v) => setBorderRadius({ bottomLeft: v })}
            />

          </div>
        )}
        <button type="button" className={s.smallBtn} onClick={resetBorderRadius}>
          Zerar border radius
        </button>
      </Section>

      <Section title="Conteúdo" storageKey="content">
        <textarea
          ref={contentTextareaRef}
          className={s.textarea}
          value={textContent}
          onChange={(e) => onTextContentChange(e.target.value)}
          placeholder="Texto do bloco…"
        />
      </Section>

      <Section title="Tipografia" storageKey="typography">
        <SelectField
          label="Família"
          value={block.font.family}
          onChange={(v) => setFont("family", v)}
          options={fontOptions ?? DEFAULT_FONTS}
        />
        <SliderField
          label="Peso"
          value={block.font.weight}
          min={100} max={900} step={100}
          onChange={(v) => setFont("weight", v)}
        />
        <div className={s.grid2}>
          <NumberField
            label="Tamanho"
            suffix="px"
            value={block.font.size}
            onChange={(v) =>
              onChange((b) => ({
                ...b,
                font: { ...b.font, size: v, sizeOverride: true },
              }))
            }
          />
          <NumberField label="Altura linha" value={block.font.lineHeight} step={0.1} onChange={(v) => setFont("lineHeight", v)} />
        </div>
        <NumberField label="Espaçamento" suffix="em" value={block.font.letterSpacing} step={0.01} onChange={(v) => setFont("letterSpacing", v)} />
        <MultiToggle
          label="Estilo"
          options={[
            {
              active: block.font.bold,
              label: <IconBold />,
              title: "Negrito",
              onToggle: () =>
                tryInlineOrBlock("bold", () =>
                  onChange((b) => ({
                    ...b,
                    font: {
                      ...b.font,
                      bold: !b.font.bold,
                      weight: !b.font.bold ? 700 : b.font.weight < 700 ? b.font.weight : 400,
                    },
                  })),
                ),
            },
            {
              active: block.font.italic,
              label: <IconItalic />,
              title: "Itálico",
              onToggle: () => tryInlineOrBlock("italic", () => setFont("italic", !block.font.italic)),
            },
            {
              active: block.font.underline,
              label: <IconUnderline />,
              title: "Sublinhado",
              onToggle: () => tryInlineOrBlock("underline", () => setFont("underline", !block.font.underline)),
            },
            {
              active: block.font.strike,
              label: <IconStrike />,
              title: "Tachado",
              onToggle: () => tryInlineOrBlock("strike", () => setFont("strike", !block.font.strike)),
            },
          ]}
        />
        <ToggleGroup
          label="Capitalização"
          value={block.font.transform}
          onChange={(v) => setFont("transform", v)}
          options={[
            { value: "none", label: "Aa", title: "Normal" },
            { value: "uppercase", label: "AA", title: "Maiúsculas" },
            { value: "lowercase", label: "aa", title: "Minúsculas" },
            { value: "capitalize", label: "Aa.", title: "Capitalizar" },
          ]}
        />
      </Section>

      <Section title="Alinhamento" storageKey="alignment">
        <ToggleGroup
          label="Horizontal (text-align)"
          value={block.align.horizontal}
          onChange={(v) => setAlign("horizontal", v)}
          options={[
            { value: "left", label: <IconAlignLeft />, title: "Esquerda" },
            { value: "center", label: <IconAlignCenter />, title: "Centro" },
            { value: "right", label: <IconAlignRight />, title: "Direita" },
            { value: "justify", label: <IconAlignJustify />, title: "Justificado" },
          ]}
        />
        <ToggleGroup
          label="Vertical (align-items)"
          value={block.align.vertical}
          onChange={(v) => setAlign("vertical", v)}
          options={[
            { value: "flex-start", label: <IconVTop />, title: "Topo" },
            { value: "center", label: <IconVMid />, title: "Centro" },
            { value: "flex-end", label: <IconVBot />, title: "Base" },
          ]}
        />
        <ToggleGroup
          label="Direção"
          value={block.align.direction}
          onChange={(v) => setAlign("direction", v)}
          options={[
            { value: "ltr", label: "LTR", title: "Esquerda → Direita" },
            { value: "rtl", label: "RTL", title: "Direita → Esquerda" },
          ]}
        />
        <div className={s.grid2}>
          <SelectField
            label="Quebra"
            value={block.align.whiteSpace}
            onChange={(v) => setAlign("whiteSpace", v)}
            options={[
              { value: "normal", label: "Normal" },
              { value: "nowrap", label: "Nowrap" },
              { value: "pre", label: "Pre" },
              { value: "pre-wrap", label: "Pre-wrap" },
            ]}
          />
          <SelectField
            label="Overflow"
            value={block.align.overflow}
            onChange={(v) => setAlign("overflow", v)}
            options={[
              { value: "visible", label: "Visible" },
              { value: "hidden", label: "Hidden" },
              { value: "clip", label: "Clip" },
              { value: "ellipsis", label: "Ellipsis" },
            ]}
          />
        </div>
      </Section>

      <Section title="Cores" storageKey="colors">
        <div className={s.colorFieldRow}>
          <ColorField
            label="Texto"
            value={block.color}
            onChange={(v) =>
              onChange((b) => ({ ...b, color: v, colorSet: true }))
            }
          />
          <button
            type="button"
            className={s.removeBtn}
            title="Remover cor do estilo"
            disabled={!block.colorSet}
            onClick={() => onChange((b) => ({ ...b, colorSet: false }))}
          >
            <IconTrash />
          </button>
        </div>
        <div className={s.colorFieldRow}>
          <ColorField
            label="Fundo"
            value={block.background}
            onChange={(v) =>
              onChange((b) => ({ ...b, background: v, backgroundSet: true }))
            }
          />
          <button
            type="button"
            className={s.removeBtn}
            title="Remover fundo do estilo"
            disabled={!block.backgroundSet}
            onClick={() => onChange((b) => ({ ...b, backgroundSet: false }))}
          >
            <IconTrash />
          </button>
        </div>
      </Section>

      <Section title="Propriedades CSS" defaultOpen={false} storageKey="css-props">
        <datalist id={cssPropKeyListId}>
          {CSS_PROP_KEYS.map((key) => (
            <option key={key} value={key} />
          ))}
        </datalist>
        <ul className={s.cssPropList}>
          {draftProps.map((p, i) => {
            if (POSITION_PANEL_PROP_KEYS.has(p.key.trim().toLowerCase())) {
              return null;
            }
            const valueHandler = resolveCssPropHandler(p.key);
            const isOpen = expandedPropIndex === i;
            const isEditingKey = editingPropIndex === i;
            const displayLabel = p.key.trim() || "nova propriedade";

            return (
              <li
                key={i}
                className={`${s.cssPropItem} ${isOpen ? s.cssPropItemOpen : ""}`}
              >
                <div className={s.cssPropHead}>
                  {isEditingKey ? (
                    <input
                      className={`${s.input} ${s.cssPropKeyInput}`}
                      value={p.key}
                      autoFocus
                      placeholder="propriedade"
                      list={cssPropKeyListId}
                      onChange={(e) =>
                        updatePropRowDraft(i, { key: e.target.value })
                      }
                      onBlur={() => {
                        setEditingPropIndex(null);
                        commitPropKeys();
                      }}
                      onKeyDown={(e) => {
                        if (e.key === "Enter") e.currentTarget.blur();
                        if (e.key === "Escape") {
                          setEditingPropIndex(null);
                          e.currentTarget.blur();
                        }
                      }}
                    />
                  ) : (
                    <button
                      type="button"
                      className={s.cssPropHeadToggle}
                      onClick={() => setExpandedPropIndex(isOpen ? null : i)}
                      aria-expanded={isOpen}
                    >
                      <span className={s.cssPropLabel}>{displayLabel}</span>
                      <span
                        className={`${s.cssPropChev} ${isOpen ? s.cssPropChevOpen : ""}`}
                        aria-hidden
                      >
                        <IconChevronDown />
                      </span>
                    </button>
                  )}
                  {!isEditingKey ? (
                    <button
                      type="button"
                      className={s.cssPropEditBtn}
                      onClick={() => {
                        setExpandedPropIndex(i);
                        setEditingPropIndex(i);
                      }}
                      aria-label={`Editar nome de ${displayLabel}`}
                      title="Editar propriedade"
                    >
                      <IconEdit />
                    </button>
                  ) : null}
                </div>

                {isOpen ? (
                  <div className={s.cssPropBody}>
                    <div className={s.propValueCell}>
                      {valueHandler ? (
                        valueHandler({
                          value: p.value,
                          onChange: (value) => {
                            const next = [...draftProps];
                            next[i] = { ...next[i], value };
                            persistProps(next);
                          },
                        })
                      ) : (
                        <input
                          className={s.input}
                          value={p.value}
                          placeholder="valor"
                          onChange={(e) =>
                            updatePropRowDraft(i, { value: e.target.value })
                          }
                          onBlur={() => {
                            if (p.key.trim()) commitPropKeys();
                          }}
                        />
                      )}
                    </div>
                    <div className={s.cssPropActions}>
                      <button
                        type="button"
                        className={s.cssPropRemoveBtn}
                        title="Remover"
                        onClick={() => {
                          const next = draftProps.filter((_, j) => j !== i);
                          persistProps(next);
                          setExpandedPropIndex(null);
                          setEditingPropIndex(null);
                        }}
                      >
                        <IconTrash />
                        <span>Remover</span>
                      </button>
                    </div>
                  </div>
                ) : null}
              </li>
            );
          })}
        </ul>
        <button
          type="button"
          className={s.smallBtn}
          onClick={() => {
            const next = [...draftProps, { key: "", value: "" }];
            setDraftProps(next);
            const index = next.length - 1;
            setExpandedPropIndex(index);
            setEditingPropIndex(index);
          }}
        >
          + adicionar propriedade
        </button>
      </Section>
      </>
      ) : null}
    </>
  );
}

