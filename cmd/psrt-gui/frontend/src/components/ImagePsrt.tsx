import React, {
  CSSProperties,
  type HTMLAttributes,
  type Ref,
  useEffect,
  useRef,
  useState,
} from 'react';
import ReactDOMServer from 'react-dom/server';
import {
  resolveEntryStyle,
  useAdaptedEntryStyles,
} from '../lib/useAdaptedEntryStyles';
import { textBlockWidthPx, textFontSizePx } from '../lib/psrtGeometry';
import { isPresentStyleValue } from '../lib/styleValue';
import { PSRTEntry, PSRTSection } from '../types/types';
import { FormattedText } from './FormattedText';

export interface ImageReferenceMetrics {
  refWidth: number;
  refHeight: number;
  zoom: number;
}

function EntryBlock({
  containerStyle,
  textStyle,
  className = 'entry',
  content,
  consts,
  ...divProps
}: {
  containerStyle: CSSProperties;
  textStyle: CSSProperties;
  content: string;
  consts?: Record<string, string>;
  className?: string;
} & Omit<HTMLAttributes<HTMLDivElement>, 'className' | 'style' | 'children'>) {
  return (
    <div className={className} style={containerStyle} {...divProps}>
      <span className="entry-text" style={textStyle}>
        <FormattedText content={content} consts={consts} />
      </span>
    </div>
  );
}

export default function ImagePsrt({
  pageData,
  hideSub,
  alt,
  editor,
  onClickEntry,
  selectedIndex,
  pageLink,
  parentZoom = 1,
  fixedReferenceSize = false,
  getSize,
  consts,
  imageContainerRef,
}: {
  pageData: PSRTSection;
  alt: string;
  hideSub?: boolean;
  editor?: boolean;
  onClickEntry?: (_entryIndex: number, _entry: PSRTEntry) => void;
  selectedIndex?: number;
  pageLink?: string;
  parentZoom?: number;
  fixedReferenceSize?: boolean;
  getSize?: ({ w, h }: { w?: number; h?: number }) => void;
  consts?: Record<string, string>;
  imageContainerRef?: Ref<HTMLDivElement>;
}) {
  const imageRef = useRef<HTMLImageElement>(null);
  const [refSize, setRefSize] = useState<{ w: number; h: number } | null>(null);

  const imgSrc = pageLink ?? pageData.pageLink ?? '';

  useEffect(() => {
    setRefSize(null);
  }, [imgSrc]);

  const syncRefSize = () => {
    const img = imageRef.current;
    if (!img?.naturalWidth || !img.naturalHeight) return;
    setRefSize({ w: img.naturalWidth, h: img.naturalHeight });
    getSize?.({ w: img.naturalWidth, h: img.naturalHeight });
  };

  useEffect(() => {
    const img = imageRef.current;
    if (img?.complete && img.naturalWidth > 0) {
      syncRefSize();
    }
  }, [imgSrc]);

  useEffect(() => {
    if (!fixedReferenceSize) {
      getSize?.({
        w: imageRef.current?.width,
        h: imageRef.current?.height,
      });
    }
  }, [parentZoom, getSize, fixedReferenceSize, refSize]);

  const metrics: ImageReferenceMetrics | undefined =
    refSize && fixedReferenceSize
      ? { refWidth: refSize.w, refHeight: refSize.h, zoom: parentZoom }
      : refSize
        ? { refWidth: refSize.w, refHeight: refSize.h, zoom: 1 }
        : undefined;

  const entries = pageData?.entries ?? [];
  const adaptedByIndex = useAdaptedEntryStyles(
    entries,
    metrics?.refWidth ?? 0,
    metrics?.refHeight ?? 0,
    metrics?.zoom ?? parentZoom,
  );

  const pageBg =
    pageData.pageStyle?.background ??
    (pageData.pageStyle as Record<string, string | undefined> | undefined)?.backGround;

  const containerStyle: CSSProperties | undefined = (() => {
    const bgStyle: CSSProperties =
      typeof pageBg === 'string' && pageBg ? { backgroundColor: pageBg } : {};
    if (fixedReferenceSize && refSize) {
      return {
        ...bgStyle,
        width: refSize.w * parentZoom,
        height: refSize.h * parentZoom,
        flexShrink: 0,
      };
    }
    return Object.keys(bgStyle).length > 0 ? bgStyle : undefined;
  })();

  const isCanvasSized = fixedReferenceSize && refSize !== null;
  const containerClass = [
    'relative',
    fixedReferenceSize ? 'psrt-image-container--fixed-ref' : 'w-full',
    fixedReferenceSize && (isCanvasSized ? 'psrt-image-container--sized' : 'psrt-image-container--sizing'),
  ]
    .filter(Boolean)
    .join(' ');

  const showEntries =
    !hideSub &&
    (fixedReferenceSize ? refSize !== null : Boolean(parentZoom));

  const resolveEntryStyles = (entry: PSRTEntry) => {
    const idx = entry.index ?? 0;
    const resolved = resolveEntryStyle(
      entry,
      adaptedByIndex.get(idx),
      metrics,
    );
    let container = resolved.container;
    if (metrics) {
      container = applyBalloonIfNeeded(entry, container, metrics);
    }
    const merged = resolved.merged ?? { ...container, ...resolved.text };
    return {
      container,
      text: resolved.text,
      merged,
      hasStroke: resolved.hasStroke,
    };
  };

  return (
    <div
      id="psrt-image-container"
      ref={imageContainerRef}
      className={containerClass}
      style={containerStyle}
    >
      <div id="image-container" className="w-full h-full">
        {imgSrc ? (
          <img
            ref={imageRef}
            className="psrt-page-image"
            src={imgSrc}
            alt={alt}
            draggable={false}
            onDragStart={(e) => e.preventDefault()}
            onLoad={syncRefSize}
          />
        ) : (
          <div className="preview-web-empty" style={{ padding: 24 }}>
            Sem imagem de fundo
          </div>
        )}
      </div>
      {showEntries && (
        <div id="entries-container" className="absolute w-full h-[100%] top-0 left-0">
          {editor &&
            onClickEntry != null &&
            entries.map((entry: PSRTEntry) => {
              const textIndex = entry.index ?? 0;
              const { container, text } = resolveEntryStyles(entry);
              const { style, ...coordinates } = entry;
              const isSelected = selectedIndex === textIndex;

              return (
                <EntryBlock
                  key={`${textIndex}--${pageLink}--mask`}
                  containerStyle={{
                    ...container,
                    textShadow: 'none',
                    backgroundColor: 'transparent',
                    zIndex: isSelected ? 100 : 1,
                    pointerEvents: isSelected ? 'auto' : 'none',
                    ...(isSelected
                      ? {
                          cursor: 'default',
                          borderWidth: '1px',
                          borderColor: '#01fe2b',
                          borderStyle: 'solid',
                        }
                      : { cursor: 'pointer' }),
                  }}
                  textStyle={{
                    ...text,
                    color: 'transparent',
                    WebkitTextStrokeWidth: 'unset',
                    WebkitTextStrokeColor: '#000',
                  }}
                  onClick={() => onClickEntry?.(textIndex, entry)}
                  data-cordinate={JSON.stringify(coordinates)}
                  content={entry.text ?? ''}
                  consts={consts}
                />
              );
            })}

          {entries?.map((entry: PSRTEntry) => {
            const textIndex = entry.index ?? 0;
            const { container, text } = resolveEntryStyles(entry);
            const { style, ...coordinates } = entry;

            return (
              <EntryBlock
                key={`${textIndex}--${pageLink}`}
                containerStyle={{
                  ...container,
                  cursor: 'pointer',
                  zIndex: textIndex,
                  pointerEvents: onClickEntry ? 'auto' : 'none',
                }}
                textStyle={text}
                onClick={() => onClickEntry?.(textIndex, entry)}
                data-cordinate={JSON.stringify(coordinates)}
                content={entry.text ?? ''}
                consts={consts}
              />
            );
          })}

          {entries?.map((entry: PSRTEntry) => {
            const textIndex = entry.index ?? 0;
            const { container, text, hasStroke } = resolveEntryStyles(entry);
            if (!hasStroke) return null;
            const { style, ...coordinates } = entry;

            return (
              <EntryBlock
                key={`${textIndex}--${pageLink}--stroke`}
                containerStyle={{
                  ...container,
                  pointerEvents: 'none',
                  textShadow: 'none',
                  backgroundColor: 'transparent',
                }}
                textStyle={{
                  fontSize: text.fontSize,
                  WebkitTextStrokeWidth: text.WebkitTextStrokeWidth,
                  WebkitTextStrokeColor: text.WebkitTextStrokeColor,
                  WebkitTextStroke: text.WebkitTextStroke,
                  color: 'transparent',
                }}
                data-cordinate={JSON.stringify(coordinates)}
                content={entry.text ?? ''}
                consts={consts}
              />
            );
          })}
        </div>
      )}
    </div>
  );
}

function applyBalloonIfNeeded(
  entry: PSRTEntry,
  container: CSSProperties,
  metrics: ImageReferenceMetrics,
): CSSProperties {
  const style = entry.style as CSSProperties & { balloon?: string };
  if (!style?.balloon || style.balloon === 'none') return container;

  const refWidth = metrics.refWidth;
  const refHeight = metrics.refHeight;
  const zoom = metrics.zoom;
  const fontPx = textFontSizePx(entry.size, refWidth, refHeight, zoom);
  const blockWidthPx =
    refWidth > 0 ? Math.round(textBlockWidthPx(entry.width, refWidth) * zoom) : 0;

  if (style.balloon !== 'default' || fontPx <= 0) return container;

  const balloonH = isPresentStyleValue('height', style.height)
    ? (Number(style.height) / 100) * refHeight * zoom
    : blockWidthPx;
  const svgDataUrl = svgToDataUrl(
    <SpechBaloonSvg
      stroke={parseFloat(String(style.borderWidth ?? '0'))}
      w={blockWidthPx}
      h={balloonH}
      color={String(style.backgroundColor ?? '#fff')}
      borderColor={String(style.borderColor ?? '#000')}
    />,
  );
  return {
    ...container,
    backgroundSize: 'cover',
    backgroundImage: `url(${svgDataUrl})`,
  };
}

export const SpechBaloonSvg: React.FC<{
  w: number;
  stroke: number;
  h: number;
  color: string;
  borderColor: string;
}> = ({ w, stroke, h, color, borderColor }) => {
  const strokeMax = Math.min(stroke, w * 0.4, h * 0.4);
  return (
    <svg
      width={w.toString()}
      height={h.toString()}
      viewBox={`0 0 ${w} ${h}`}
      xmlns="http://www.w3.org/1999/xhtml"
      fill="blue"
    >
      <ellipse
        cx={w / 2}
        cy={h / 2}
        rx={w / 2 - strokeMax / 2}
        ry={h / 2 - strokeMax / 2}
        fill="none"
        stroke={borderColor}
        strokeWidth={strokeMax}
      />
      {stroke > 0 && (
        <polygon
          points={`${0},${h} ${w * 0.33 - Math.min(strokeMax * 4, w * 0.1)},${h - h * 0.4} ${Math.min(w * 0.62 + strokeMax * 4, w * 0.8)},${h - h * 0.4}`}
          fill={borderColor}
        />
      )}
      <polygon
        points={`${0},${h} ${w * 0.8},${h - h * 0.5} ${w * 0.4},${h / 2}`}
        fill={color}
      />
      <ellipse
        cx={w / 2}
        cy={h / 2}
        rx={w / 2 - strokeMax / 2}
        ry={h / 2 - strokeMax / 2}
        fill={color}
      />
    </svg>
  );
};

export function svgToDataUrl(svg: React.ReactElement): string {
  const svgString = ReactDOMServer.renderToStaticMarkup(svg);
  const encodedSvg = encodeURIComponent(svgString);
  return `data:image/svg+xml,${encodedSvg}`;
}
