(function () {
  function parseFonts() {
    window.__psrtFontsParsed = {};
    if (!window.__psrtFonts || typeof opentype === 'undefined') return;
    for (const name of Object.keys(window.__psrtFonts)) {
      try {
        const bin = atob(window.__psrtFonts[name]);
        const buf = new Uint8Array(bin.length);
        for (let i = 0; i < bin.length; i++) buf[i] = bin.charCodeAt(i);
        window.__psrtFontsParsed[name] = opentype.parse(buf.buffer);
      } catch (e) {
        /* skip bad font */
      }
    }
  }

  function fontWeightNumber(w) {
    if (!w || w === 'normal') return 400;
    if (w === 'bold' || w === 'bolder') return 700;
    const n = parseInt(w, 10);
    return Number.isFinite(n) ? n : 400;
  }

  function pickFont(style) {
    const families = style.fontFamily.split(',');
    const weight = fontWeightNumber(style.fontWeight);
    const bold = weight >= 600;
    const italic = style.fontStyle === 'italic' || style.fontStyle === 'oblique';

    for (const raw of families) {
      const name = raw.replace(/['"]/g, '').trim();
      if (!name) continue;
      const parsed = window.__psrtFontsParsed || {};
      if (bold && italic) {
        for (const key of [name + ' Bold Italic', name + '-BoldItalic', name + 'BoldItalic']) {
          if (parsed[key]) return parsed[key];
        }
      }
      if (bold) {
        for (const key of [name + ' Bold', name + '-Bold', name + 'Bold', name + '-700']) {
          if (parsed[key]) return parsed[key];
        }
      }
      if (italic) {
        for (const key of [name + ' Italic', name + '-Italic', name + 'Italic']) {
          if (parsed[key]) return parsed[key];
        }
      }
      if (parsed[name]) return parsed[name];
    }
    if (window.__psrtFontsParsed && window.__psrtFontsParsed.PSRTDefault) {
      return window.__psrtFontsParsed.PSRTDefault;
    }
    const keys = Object.keys(window.__psrtFontsParsed || {});
    return keys.length ? window.__psrtFontsParsed[keys[0]] : null;
  }

  function canvasFont(style) {
    return (
      (style.fontStyle || 'normal') +
      ' ' +
      (style.fontWeight || '400') +
      ' ' +
      style.fontSize +
      ' ' +
      style.fontFamily
    );
  }

  function baselineY(rect, style, text) {
    const canvas = document.createElement('canvas').getContext('2d');
    canvas.font = canvasFont(style);
    const m = canvas.measureText(text || 'Mg');
    const ascent = m.actualBoundingBoxAscent || parseFloat(style.fontSize) * 0.8;
    return rect.top + ascent;
  }

  function pathAttrs(style) {
    const fill = style.color && style.color !== 'rgba(0, 0, 0, 0)' ? style.color : '#000000';
    const out = { fill: fill, opacity: style.opacity || '' };
    const sw = style.webkitTextStrokeWidth || '0';
    const sc = style.webkitTextStrokeColor || style.webkitTextStroke || '';
    if (sw && sw !== '0px' && sc) {
      out.stroke = sc;
      out.strokeWidth = sw;
      out.paintOrder = 'stroke fill';
    }
    return out;
  }

  function styleElementForTextNode(node) {
    let el = node.parentElement;
    while (el && el !== document.body) {
      const tag = el.tagName;
      if (tag === 'STRONG' || tag === 'EM' || tag === 'U' || tag === 'S') {
        return el;
      }
      if (el.classList) {
        for (const cls of el.classList) {
          if (cls.endsWith('-inner')) return el;
        }
      }
      el = el.parentElement;
    }
    return node.parentElement || node;
  }

  function emitPath(text, rect, style) {
    if (!text || !text.trim()) return null;
    const font = pickFont(style);
    if (!font) return null;
    const fontSize = parseFloat(style.fontSize);
    if (!fontSize || fontSize <= 0) return null;
    const y = baselineY(rect, style, text);
    let path;
    try {
      path = font.getPath(text, 0, 0, fontSize);
    } catch (e) {
      return null;
    }
    const bb = path.getBoundingBox();
    const naturalW = bb.x2 - bb.x1;
    let sx = 1;
    if (naturalW > 0 && rect.width > 0) {
      sx = rect.width / naturalW;
      if (Math.abs(sx - 1) > 0.02) {
        path.scale(sx, 1);
      }
    }
    path.translate(rect.left - bb.x1 * sx, y);
    const d = path.toPathData(2);
    if (!d) return null;
    const attrs = pathAttrs(style);
    attrs.d = d;
    return attrs;
  }

  function trimSegmentEnd(text, start, end) {
    if (end <= start + 1) return end;
    let e = end;
    while (e > start && /\s/.test(text[e - 1])) e--;
    return e;
  }

  function processTextNode(node, paths) {
    const text = node.textContent;
    if (!text) return;
    const range = document.createRange();
    let start = 0;
    while (start < text.length) {
      range.setStart(node, start);
      range.setEnd(node, text.length);
      const rects = range.getClientRects();
      if (!rects.length) break;
      const firstTop = rects[0].top;
      let end = start + 1;
      while (end <= text.length) {
        range.setEnd(node, end);
        const rs = range.getClientRects();
        if (!rs.length) break;
        const last = rs[rs.length - 1];
        if (end > start + 1 && last.top > firstTop + 1) break;
        end++;
      }
      if (end > start + 1) {
        end--;
      }
      if (end <= start) end = start + 1;
      if (end > text.length) end = text.length;
      end = trimSegmentEnd(text, start, end);
      if (end <= start) {
        start++;
        continue;
      }
      range.setStart(node, start);
      range.setEnd(node, end);
      const rect = range.getBoundingClientRect();
      if (rect.width > 0 && rect.height > 0) {
        const style = window.getComputedStyle(styleElementForTextNode(node));
        const seg = text.slice(start, end);
        const p = emitPath(seg, rect, style);
        if (p) paths.push(p);
      }
      start = end;
      while (start < text.length && /\s/.test(text[start])) start++;
    }
  }

  function walk(node, paths) {
    if (!node) return;
    if (node.nodeType === Node.TEXT_NODE) {
      processTextNode(node, paths);
      return;
    }
    if (node.nodeType !== Node.ELEMENT_NODE) return;
    if (node.nodeName === 'BR') return;
    for (const child of node.childNodes) {
      walk(child, paths);
    }
  }

  function outlineBlock(blockEl) {
    const index = parseInt(blockEl.getAttribute('data-block-index') || '0', 10);
    const inner = blockEl.querySelector('span');
    const paths = [];
    if (inner) {
      walk(inner, paths);
    }
    const plain = inner ? inner.innerText : blockEl.innerText;
    const innerWrap = blockEl.firstElementChild || blockEl;
    const style = window.getComputedStyle(innerWrap);
    const transform = style.transform && style.transform !== 'none' ? style.transform : '';
    return { index: index, paths: paths, plainText: plain || '', transform: transform };
  }

  parseFonts();

  return new Promise(function (resolve, reject) {
    document.fonts.ready
      .then(function () {
        parseFonts();
        const blocks = [];
        document.querySelectorAll('.psrt-text-block').forEach(function (el) {
          blocks.push(outlineBlock(el));
        });
        resolve(blocks);
      })
      .catch(reject);
  });
})();
